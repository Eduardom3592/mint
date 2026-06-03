package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"

	"github.com/programmersd21/mint/internal/downloads"
	"github.com/programmersd21/mint/internal/mrpack"
	"github.com/programmersd21/mint/internal/platform"
)

type downloadTab int

const (
	tabActive downloadTab = iota
	tabHistory
	tabInstalledModpacks
	tabFailed
)

type downloadViewModel struct {
	dlmgr       *downloads.Manager
	tab         downloadTab
	cursor      int
	scrollTop   int
	styles      Styles
	askClearAll bool
	clearAllMsg string
}

func newDownloadViewModel(dlmgr *downloads.Manager) downloadViewModel {
	return downloadViewModel{
		dlmgr:  dlmgr,
		tab:    tabActive,
		cursor: 0,
	}
}

func (m *downloadViewModel) setStyles(s Styles) {
	m.styles = s
}

func (m downloadViewModel) Init() tea.Cmd { return nil }

func (m downloadViewModel) items() []*downloads.Item {
	all := m.dlmgr.List()
	switch m.tab {
	case tabActive:
		var active []*downloads.Item
		for _, it := range all {
			if it.Status == downloads.StatusQueued || it.Status == downloads.StatusPreparing || it.Status == downloads.StatusDownloading || it.Status == downloads.StatusVerifying {
				active = append(active, it)
			}
		}
		return active
	case tabHistory:
		var done []*downloads.Item
		for _, it := range all {
			if it.Status == downloads.StatusCompleted {
				done = append(done, it)
			}
		}
		return done
	case tabInstalledModpacks:
		return nil
	case tabFailed:
		var failed []*downloads.Item
		for _, it := range all {
			if it.Status == downloads.StatusFailed || it.Status == downloads.StatusCancelled {
				failed = append(failed, it)
			}
		}
		return failed
	}
	return nil
}

func (m downloadViewModel) Update(msg tea.Msg) (downloadViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.askClearAll {
			switch {
			case key.Matches(msg, keys.Enter) || key.Matches(msg, keys.Yes):
				m.clearAll()
				m.askClearAll = false
				m.clearAllMsg = ""
				return m, nil
			default:
				m.askClearAll = false
				m.clearAllMsg = ""
				return m, nil
			}
		}

		switch {
		case key.Matches(msg, keys.Left):
			if m.tab > 0 {
				m.tab--
				m.cursor = 0
				m.scrollTop = 0
			}
			return m, nil

		case key.Matches(msg, keys.Right):
			if m.tab < tabFailed {
				m.tab++
				m.cursor = 0
				m.scrollTop = 0
			}
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case key.Matches(msg, keys.Down):
			items := m.items()
			if m.cursor < len(items)-1 {
				m.cursor++
			}
			return m, nil

		case key.Matches(msg, keys.Cancel):
			items := m.items()
			if m.cursor >= 0 && m.cursor < len(items) {
				m.dlmgr.Cancel(items[m.cursor].ID)
			}
			return m, nil

		case key.Matches(msg, keys.Retry):
			items := m.items()
			if m.cursor >= 0 && m.cursor < len(items) {
				m.dlmgr.Retry(items[m.cursor].ID)
			}
			return m, nil

		case key.Matches(msg, keys.Delete):
			items := m.items()
			if m.cursor >= 0 && m.cursor < len(items) {
				id := int64(items[m.cursor].ID)
				if p := m.dlmgr.Persist(); p != nil {
					_ = p.DeleteDownload(id)
				}
				items[m.cursor].Status = downloads.StatusCancelled
			}
			return m, nil

		case key.Matches(msg, keys.Open):
			items := m.items()
			if m.cursor >= 0 && m.cursor < len(items) {
				item := items[m.cursor]
				if item.Status == downloads.StatusCompleted {
					if p := item.DestPath(); p != "" {
						platform.OpenFileInDir(p)
						return m, nil
					}
				}
			}
			platform.OpenFolder(m.dlmgr.Dir())
			return m, nil

		case msg.String() == "X":
			m.askClearAll = true
			m.clearAllMsg = "delete all downloaded files and records? (enter=yes, any key=no)"
			return m, nil
		}
	}

	return m, nil
}

func (m downloadViewModel) clearAll() {
	for _, it := range m.dlmgr.List() {
		if p := it.DestPath(); p != "" {
			os.Remove(p)
		}
	}
	installedDir := filepath.Join(m.dlmgr.Dir(), "installed")
	os.RemoveAll(installedDir)
	m.dlmgr.ClearAll()
}

func (m downloadViewModel) View() string {
	var sb strings.Builder
	s := m.styles

	sb.WriteString(s.Title.Render("downloads"))
	sb.WriteString("\n\n")

	tabNames := []string{"active", "history", "installed modpacks", "failed"}
	for i, name := range tabNames {
		if i == int(m.tab) {
			sb.WriteString(s.Accent.Render(fmt.Sprintf(" [%s]", name)))
		} else {
			sb.WriteString(s.Dimmed.Render(fmt.Sprintf(" %s ", name)))
		}
	}
	sb.WriteString("\n\n")

	if m.tab == tabInstalledModpacks {
		sb.WriteString(m.renderInstalled())
		if m.askClearAll {
			fmt.Fprintf(&sb, "\n  %s\n", s.Warning.Render(m.clearAllMsg))
		}
		return sb.String()
	}

	items := m.items()
	if len(items) == 0 {
		emptyMsg := "no active downloads"
		switch m.tab {
		case tabHistory:
			emptyMsg = "no completed downloads"
		case tabFailed:
			emptyMsg = "no failed downloads"
		}
		fmt.Fprintf(&sb, "  %s\n\n", s.Info.Render(emptyMsg))
		if m.tab == tabActive {
			fmt.Fprintf(&sb, "  %s", s.Dimmed.Render("press d on a project or home card to start"))
		}
		if m.askClearAll {
			fmt.Fprintf(&sb, "\n  %s\n", s.Warning.Render(m.clearAllMsg))
		}
		return sb.String()
	}

	height := 25
	start := m.cursor - height/2
	if start < 0 {
		start = 0
	}
	end := start + height
	if end > len(items) {
		end = len(items)
		start = end - height
		if start < 0 {
			start = 0
		}
	}

	for i := start; i < end; i++ {
		item := items[i]

		prefix := "  "
		if i == m.cursor {
			prefix = s.Accent.Render("> ")
		}

		nameStr := item.Filename
		if len(nameStr) > 40 {
			nameStr = nameStr[:37] + "..."
		}

		statusColor := s.Dimmed
		statusIcon := "."
		switch item.Status {
		case downloads.StatusQueued:
			statusColor = s.Dimmed
			statusIcon = "."
		case downloads.StatusPreparing:
			statusColor = s.Primary
			statusIcon = "~"
		case downloads.StatusDownloading:
			statusColor = s.Primary
			statusIcon = "↓"
		case downloads.StatusVerifying:
			statusColor = s.Accent
			statusIcon = "⟳"
		case downloads.StatusCompleted:
			statusColor = s.Success
			statusIcon = "✓"
		case downloads.StatusFailed:
			statusColor = s.Error
			statusIcon = "✖"
		case downloads.StatusCancelled:
			statusColor = s.Warning
			statusIcon = "−"
		}

		fmt.Fprintf(&sb, "%s%s", prefix, statusColor.Render(statusIcon))
		fmt.Fprintf(&sb, " %s", nameStr)

		if item.Status == downloads.StatusDownloading {
			bar := ProgressBar(item.Progress, 20)
			fmt.Fprintf(&sb, " %s", bar)
		} else if item.Status == downloads.StatusCompleted {
			fmt.Fprintf(&sb, " %s", s.Success.Render("done"))
		} else if item.Status == downloads.StatusFailed && i == m.cursor {
			errStr := item.Error
			if len(errStr) > 20 {
				errStr = errStr[:17] + "..."
			}
			fmt.Fprintf(&sb, " %s", s.Error.Render(errStr))
		} else if item.Status == downloads.StatusVerifying {
			sb.WriteString(s.Accent.Render(" verifying..."))
		}

		sb.WriteString("\n")

		if i == m.cursor {
			sb.WriteString(m.renderDetail(items[i]))
		}
	}

	sb.WriteString("\n")

	if m.askClearAll {
		fmt.Fprintf(&sb, "  %s\n\n", s.Warning.Render(m.clearAllMsg))
	}

	help := "h/l: tabs | j/k: navigate"
	if m.tab == tabActive {
		help += " | c: cancel"
	}
	if m.tab == tabFailed {
		help += " | r: retry"
	}
	help += " | o: open | del: remove | X: clear all"
	sb.WriteString(s.Dimmed.Render(help))

	return sb.String()
}

type installedModpack struct {
	name    string
	version string
	loader  string
}

func (m downloadViewModel) renderInstalled() string {
	var sb strings.Builder
	s := m.styles

	modpacks := m.scanInstalledModpacks()
	if len(modpacks) == 0 {
		fmt.Fprintf(&sb, "  %s\n\n", s.Info.Render("no installed modpacks"))
		fmt.Fprintf(&sb, "  %s", s.Dimmed.Render("install a modpack with I"))
		return sb.String()
	}
	for _, mp := range modpacks {
		fmt.Fprintf(&sb, "  %s %s\n", s.Success.Render("✓"), mp.name)
		if mp.version != "" {
			fmt.Fprintf(&sb, "    %s: %s\n", s.Info.Render("version"), mp.version)
		}
		if mp.loader != "" {
			fmt.Fprintf(&sb, "    %s: %s\n", s.Info.Render("loader"), mp.loader)
		}
	}

	return sb.String()
}

func (m downloadViewModel) scanInstalledModpacks() []installedModpack {
	installedDir := filepath.Join(m.dlmgr.Dir(), "installed")
	entries, err := os.ReadDir(installedDir)
	if err != nil {
		return nil
	}
	var modpacks []installedModpack
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		indexPath := filepath.Join(installedDir, entry.Name(), "modrinth.index.json")
		idx, _, err := mrpack.Parse(indexPath)
		if err != nil {
			modpacks = append(modpacks, installedModpack{name: entry.Name()})
			continue
		}
		ver := idx.VersionID
		loader := ""
		for depID := range idx.Dependencies {
			if depID == "fabric-loader" || depID == "quilt-loader" || depID == "forge" || depID == "neoforge" {
				loader = depID
				break
			}
		}
		modpacks = append(modpacks, installedModpack{
			name:    idx.Name,
			version: ver,
			loader:  loader,
		})
	}
	return modpacks
}

func (m downloadViewModel) renderDetail(item *downloads.Item) string {
	var sb strings.Builder
	s := m.styles

	if item.ProjectTitle != "" {
		fmt.Fprintf(&sb, "    %s: %s\n", s.Info.Render("project"), item.ProjectTitle)
	}
	if item.VersionNumber != "" {
		fmt.Fprintf(&sb, "    %s: %s\n", s.Info.Render("version"), item.VersionNumber)
	}

	fmt.Fprintf(&sb, "    %s: %s\n", s.Info.Render("file"), item.Filename)

	if item.TotalSize > 0 {
		pct := int(item.Progress * 100)
		fmt.Fprintf(&sb, "    %s: %s / %s (%d%%)\n",
			s.Info.Render("size"),
			formatBytes(item.DownloadedSize),
			formatBytes(item.TotalSize),
			pct)
	} else if item.DownloadedSize > 0 {
		fmt.Fprintf(&sb, "    %s: %s\n", s.Info.Render("downloaded"), formatBytes(item.DownloadedSize))
	}

	if item.Status == downloads.StatusDownloading {
		if item.Speed > 0 {
			fmt.Fprintf(&sb, "    %s: %s/s\n", s.Info.Render("speed"), formatBytes(int64(item.Speed)))
		}
		if item.ETA > 0 && item.ETA < 24*time.Hour {
			eta := item.ETA.Round(time.Second)
			fmt.Fprintf(&sb, "    %s: %s\n", s.Info.Render("eta"), eta)
		}
	}

	if item.Hash != nil {
		status := s.Success.Render("verified")
		if !item.Verified && item.Status == downloads.StatusCompleted {
			status = s.Error.Render("verification failed")
		} else if item.Status == downloads.StatusDownloading || item.Status == downloads.StatusVerifying {
			status = s.Info.Render("pending...")
		}
		fmt.Fprintf(&sb, "    %s: %s\n", s.Info.Render("hash"), status)
	}

	return sb.String()
}
