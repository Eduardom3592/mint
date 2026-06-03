package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"

	"github.com/programmersd21/mint/internal/api"
	"github.com/programmersd21/mint/internal/downloads"
	"github.com/programmersd21/mint/internal/models"
)

type pickerFilter struct {
	mcVersion string
	loader    string
	channel   string
}

type versionPickerModel struct {
	client       *api.Client
	dlmgr        *downloads.Manager
	projectID    string
	projectTitle string
	allVersions  []models.Version
	filtered     []models.Version
	loading      bool
	err          error
	cursor       int
	fileCursor   int
	selected     bool
	quick        bool
	filter       pickerFilter
	filterCursor int
	mcVersions   []string
	loaders      []string
	channels     []string
	showFilters  bool
	width        int
	styles       Styles
}

func newVersionPickerModel(client *api.Client, dlmgr *downloads.Manager) versionPickerModel {
	return versionPickerModel{
		client:   client,
		dlmgr:    dlmgr,
		styles:   glob,
		channels: []string{"all", "release", "beta", "alpha"},
	}
}

func (m *versionPickerModel) setStyles(s Styles) {
	m.styles = s
}

func (m *versionPickerModel) loadVersions(projectID, projectTitle string, quick bool) tea.Cmd {
	m.projectID = projectID
	m.projectTitle = projectTitle
	m.quick = quick
	m.loading = true
	m.err = nil
	m.filtered = nil
	m.cursor = 0
	m.fileCursor = 0
	m.selected = false
	m.showFilters = false

	return func() tea.Msg {
		versions, err := m.client.GetProjectVersions(projectID, nil, nil)
		if err != nil {
			return pickerErrorMsg{err: err}
		}
		return pickerVersionsMsg{projectID: projectID, projectTitle: projectTitle, versions: versions}
	}
}

type pickerVersionsMsg struct {
	projectID    string
	projectTitle string
	versions     []models.Version
}

type pickerErrorMsg struct {
	err error
}

func (m versionPickerModel) Update(msg tea.Msg) (versionPickerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case pickerVersionsMsg:
		m.allVersions = msg.versions
		m.loading = false

		mcSet := make(map[string]bool)
		loaderSet := make(map[string]bool)
		for _, v := range msg.versions {
			for _, gv := range v.GameVersions {
				mcSet[gv] = true
			}
			for _, l := range v.Loaders {
				loaderSet[l] = true
			}
		}
		m.mcVersions = nil
		for v := range mcSet {
			m.mcVersions = append(m.mcVersions, v)
		}
		m.loaders = nil
		for l := range loaderSet {
			m.loaders = append(m.loaders, l)
		}

		m.applyFilter()

		if m.quick && len(m.filtered) > 0 {
			v := m.filtered[0]
			for _, f := range v.Files {
				if f.Primary {
					var hash *downloads.HashInfo
					if h, ok := f.Hashes["sha1"]; ok {
						hash = &downloads.HashInfo{Type: downloads.HashSHA1, Value: h}
					} else if h, ok := f.Hashes["sha512"]; ok {
						hash = &downloads.HashInfo{Type: downloads.HashSHA512, Value: h}
					}
					m.dlmgr.Enqueue(m.projectID, m.projectTitle, v.ID, v.VersionNumber, f.URL, f.Filename, f.Size, hash)
					break
				}
			}
			m.allVersions = nil
		}
		return m, nil

	case pickerErrorMsg:
		m.err = msg.err
		m.loading = false
		return m, nil

	case tea.KeyMsg:
		km := tea.KeyMsg(msg)
		if m.selected {
			return m.handleFileSelection(km)
		}
		if m.showFilters {
			return m.handleFilterInput(km)
		}
		return m.handleVersionList(km)
	}

	return m, nil
}

func (m versionPickerModel) handleFilterInput(msg tea.KeyMsg) (versionPickerModel, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Escape):
		m.showFilters = false
		return m, nil

	case key.Matches(msg, keys.Up):
		if m.filterCursor > 0 {
			m.filterCursor--
		}
		return m, nil

	case key.Matches(msg, keys.Down):
		if m.filterCursor < 2 {
			m.filterCursor++
		}
		return m, nil

	case key.Matches(msg, keys.Enter):
		return m, nil

	case key.Matches(msg, keys.Tab), key.Matches(msg, keys.ShiftTab):
		m.showFilters = false
		return m, nil
	}

	switch m.filterCursor {
	case 0:
		if key.Matches(msg, keys.Right) || key.Matches(msg, keys.Down) {
			m.filter.mcVersion = nextOption(m.filter.mcVersion, m.mcVersions)
			m.applyFilter()
		} else if key.Matches(msg, keys.Left) || key.Matches(msg, keys.Up) {
			m.filter.mcVersion = prevOption(m.filter.mcVersion, m.mcVersions)
			m.applyFilter()
		}
	case 1:
		if key.Matches(msg, keys.Right) || key.Matches(msg, keys.Down) {
			m.filter.loader = nextOption(m.filter.loader, m.loaders)
			m.applyFilter()
		} else if key.Matches(msg, keys.Left) || key.Matches(msg, keys.Up) {
			m.filter.loader = prevOption(m.filter.loader, m.loaders)
			m.applyFilter()
		}
	case 2:
		if key.Matches(msg, keys.Right) || key.Matches(msg, keys.Down) {
			m.filter.channel = nextOption(m.filter.channel, m.channels)
			m.applyFilter()
		} else if key.Matches(msg, keys.Left) || key.Matches(msg, keys.Up) {
			m.filter.channel = prevOption(m.filter.channel, m.channels)
			m.applyFilter()
		}
	}

	return m, nil
}

func nextOption(current string, options []string) string {
	for i, o := range options {
		if o == current && i < len(options)-1 {
			return options[i+1]
		}
	}
	if len(options) > 0 {
		return options[0]
	}
	return current
}

func prevOption(current string, options []string) string {
	for i, o := range options {
		if o == current && i > 0 {
			return options[i-1]
		}
	}
	if len(options) > 0 {
		return options[len(options)-1]
	}
	return current
}

func (m versionPickerModel) handleVersionList(msg tea.KeyMsg) (versionPickerModel, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Escape):
		m.allVersions = nil
		return m, nil

	case key.Matches(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil

	case key.Matches(msg, keys.Down):
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
		return m, nil

	case key.Matches(msg, keys.Enter):
		if len(m.filtered) > 0 && m.cursor >= 0 && m.cursor < len(m.filtered) {
			m.selected = true
			m.fileCursor = 0
		}
		return m, nil

	case key.Matches(msg, keys.Home):
		m.cursor = 0
		return m, nil

	case key.Matches(msg, keys.End):
		m.cursor = len(m.filtered) - 1
		if m.cursor < 0 {
			m.cursor = 0
		}
		return m, nil

	case key.Matches(msg, keys.Tab):
		if !m.showFilters {
			m.showFilters = true
			m.filterCursor = 0
		}
		return m, nil
	}
	return m, nil
}

func (m versionPickerModel) handleFileSelection(msg tea.KeyMsg) (versionPickerModel, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Escape):
		m.selected = false
		return m, nil

	case key.Matches(msg, keys.Up):
		if m.fileCursor > 0 {
			m.fileCursor--
		}
		return m, nil

	case key.Matches(msg, keys.Down):
		idx := m.cursor
		if idx >= 0 && idx < len(m.filtered) {
			if m.fileCursor < len(m.filtered[idx].Files)-1 {
				m.fileCursor++
			}
		}
		return m, nil

	case key.Matches(msg, keys.Enter):
		idx := m.cursor
		if idx >= 0 && idx < len(m.filtered) && m.fileCursor >= 0 && m.fileCursor < len(m.filtered[idx].Files) {
			v := m.filtered[idx]
			f := v.Files[m.fileCursor]
			var hash *downloads.HashInfo
			if h, ok := f.Hashes["sha1"]; ok {
				hash = &downloads.HashInfo{Type: downloads.HashSHA1, Value: h}
			} else if h, ok := f.Hashes["sha512"]; ok {
				hash = &downloads.HashInfo{Type: downloads.HashSHA512, Value: h}
			}
			m.dlmgr.Enqueue(m.projectID, m.projectTitle, v.ID, v.VersionNumber, f.URL, f.Filename, f.Size, hash)
			m.allVersions = nil
		}
		return m, nil
	}
	return m, nil
}

func (m *versionPickerModel) applyFilter() {
	m.filtered = nil
	for _, v := range m.allVersions {
		if m.filter.channel != "" && m.filter.channel != "all" && string(v.VersionType) != m.filter.channel {
			continue
		}
		if m.filter.mcVersion != "" {
			found := false
			for _, gv := range v.GameVersions {
				if gv == m.filter.mcVersion {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if m.filter.loader != "" {
			found := false
			for _, l := range v.Loaders {
				if l == m.filter.loader {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		m.filtered = append(m.filtered, v)
	}
	if m.cursor >= len(m.filtered) {
		m.cursor = 0
	}
}

func (m versionPickerModel) View() string {
	s := m.styles

	if m.loading {
		return fmt.Sprintf("\n  %s", s.Info.Render("loading versions..."))
	}

	if m.err != nil {
		return fmt.Sprintf("\n  %s", s.Error.Render(fmt.Sprintf("error: %s", m.err)))
	}

	if m.selected {
		return m.fileSelectionView()
	}

	return m.versionListView()
}

func (m versionPickerModel) versionListView() string {
	var sb strings.Builder
	s := m.styles

	sb.WriteString(s.Title.Render(fmt.Sprintf("select version: %s", m.projectTitle)))
	sb.WriteString("\n\n")

	if m.showFilters {
		sb.WriteString(s.Section.Render("filters"))
		sb.WriteString("\n")

		filters := []struct {
			label string
			value string
		}{
			{"mc version", m.filter.mcVersion},
			{"loader", m.filter.loader},
			{"channel", m.filter.channel},
		}
		for i, f := range filters {
			val := f.value
			if val == "" {
				val = "all"
			}
			if i == m.filterCursor {
				sb.WriteString(s.Accent.Render(fmt.Sprintf(" > %s: %s", f.label, val)))
			} else {
				fmt.Fprintf(&sb, "   %s: %s", f.label, val)
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
		sb.WriteString(s.Dimmed.Render("tab: done - esc: cancel"))
		sb.WriteString("\n\n")
	}

	fmt.Fprintf(&sb, "  %s", s.Info.Render(fmt.Sprintf("%d versions", len(m.filtered))))
	sb.WriteString("\n\n")

	maxResults := len(m.filtered)
	height := 15
	start := m.cursor - height/2
	if start < 0 {
		start = 0
	}
	end := start + height
	if end > maxResults {
		end = maxResults
		start = end - height
		if start < 0 {
			start = 0
		}
	}

	for i := start; i < end; i++ {
		v := m.filtered[i]
		prefix := "  "
		if i == m.cursor {
			prefix = s.Accent.Render("> ")
		}
		typeStr := string(v.VersionType)
		if v.Featured {
			typeStr = "*" + typeStr
		}
		loaders := strings.Join(v.Loaders, ",")
		versions := strings.Join(v.GameVersions, ",")

		fmt.Fprintf(&sb, "%s%s", prefix, v.VersionNumber)
		fmt.Fprintf(&sb, "  [%s]", typeStr)
		fmt.Fprintf(&sb, "  %s\n", v.DatePublished.Format("2006-01-02"))
		if i == m.cursor {
			fmt.Fprintf(&sb, "    loaders: %s\n", s.Info.Render(loaders))
			fmt.Fprintf(&sb, "    mc: %s\n", s.Info.Render(versions))
			fmt.Fprintf(&sb, "    files: %d  downloads: %s\n", len(v.Files), formatNumber(v.Downloads))
		}
	}

	if len(m.filtered) > height {
		sb.WriteString("\n")
		sb.WriteString(s.Dimmed.Render(scrollLine(m.cursor, len(m.filtered), height)))
	}

	sb.WriteString("\n\n")
	if !m.showFilters {
		sb.WriteString(s.Dimmed.Render("j/k: navigate - enter: select version - tab: filter - esc: back"))
	}

	return sb.String()
}

func (m versionPickerModel) fileSelectionView() string {
	var sb strings.Builder
	s := m.styles

	if m.cursor < 0 || m.cursor >= len(m.filtered) {
		return ""
	}

	v := m.filtered[m.cursor]
	sb.WriteString(s.Title.Render(fmt.Sprintf("select file: %s %s", v.VersionNumber, v.Name)))
	sb.WriteString("\n\n")

	fmt.Fprintf(&sb, "  %s\n", s.Info.Render(fmt.Sprintf("type: %s  date: %s  downloads: %s",
		v.VersionType, v.DatePublished.Format("2006-01-02"), formatNumber(v.Downloads))))
	fmt.Fprintf(&sb, "  loaders: %s\n", strings.Join(v.Loaders, ", "))
	fmt.Fprintf(&sb, "  mc: %s\n", strings.Join(v.GameVersions, ", "))
	sb.WriteString("\n")

	if len(v.Dependencies) > 0 {
		sb.WriteString(s.Section.Render("dependencies"))
		sb.WriteString("\n")
		for _, dep := range v.Dependencies {
			var target string
			if dep.ProjectID != nil {
				target = *dep.ProjectID
			} else if dep.VersionID != nil {
				target = *dep.VersionID
			}
			var sym string
			switch dep.DependencyType {
			case "required":
				sym = s.Success.Render("  ✓")
			case "optional":
				sym = s.Info.Render("  ○")
			case "incompatible":
				sym = s.Error.Render("  ✖")
			default:
				sym = s.Info.Render(fmt.Sprintf("  %s", dep.DependencyType[:1]))
			}
			fmt.Fprintf(&sb, "%s %s\n", sym, target)
		}
		sb.WriteString("\n")
	}

	if len(v.Files) > 0 {
		sb.WriteString(s.Section.Render("files"))
		sb.WriteString("\n")
		for i, f := range v.Files {
			primary := ""
			if f.Primary {
				primary = s.Accent.Render(" *primary")
			}
			prefix := "  "
			if i == m.fileCursor {
				prefix = s.Accent.Render("> ")
			}
			fmt.Fprintf(&sb, "%s%s (%s)%s\n", prefix, f.Filename, formatBytes(f.Size), primary)
			if i == m.fileCursor {
				if h, ok := f.Hashes["sha1"]; ok {
					fmt.Fprintf(&sb, "    sha1: %s\n", s.Dimmed.Render(truncateText(h, 16)))
				}
			}
		}
	}

	sb.WriteString("\n")
	sb.WriteString(s.Dimmed.Render("enter: download - up/down: select file - esc: back"))

	return sb.String()
}
