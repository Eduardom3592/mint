package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"

	"github.com/programmersd21/mint/internal/api"
	"github.com/programmersd21/mint/internal/cache"
	"github.com/programmersd21/mint/internal/models"
	"github.com/programmersd21/mint/internal/mrpack"
)

type versionModel struct {
	client    *api.Client
	cache     *cache.Cache
	projectID string
	versions  []models.Version
	selected  int
	loading   bool
	err       error
	page      versionPage
	detail    *models.Version
	frame     int
	styles    Styles

	fileCursor int
	mrpackInfo *mrpack.PackIndex
}

func newVersionModel(client *api.Client, c *cache.Cache) versionModel {
	return versionModel{
		client: client,
		cache:  c,
		page:   versionList,
		styles: glob,
	}
}

func (m *versionModel) setStyles(s Styles) {
	m.styles = s
}

func (m versionModel) Init() tea.Cmd { return nil }

func (m *versionModel) load(projectID string) tea.Cmd {
	m.projectID = projectID
	m.loading = true
	m.err = nil
	m.selected = 0
	m.page = versionList

	cached, err := m.cache.GetVersions(projectID)
	if err == nil && len(cached) > 0 {
		m.versions = cached
		m.loading = false
		return nil
	}

	return func() tea.Msg {
		versions, err := m.client.GetProjectVersions(projectID, nil, nil)
		if err != nil {
			return versionErrorMsg{err: err}
		}

		_ = m.cache.CacheVersions(projectID, versions)

		return versionsLoadedMsg{versions: versions}
	}
}

type versionsLoadedMsg struct {
	versions []models.Version
}

type versionErrorMsg struct {
	err error
}

func (m versionModel) Update(msg tea.Msg) (versionModel, tea.Cmd) {
	switch msg := msg.(type) {
	case versionsLoadedMsg:
		m.versions = msg.versions
		m.loading = false

	case versionErrorMsg:
		m.err = msg.err
		m.loading = false

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Escape):
			if m.page == versionDetail {
				m.page = versionList
				m.detail = nil
			} else {
				m.page = versionSearch
			}
			return m, nil

		case key.Matches(msg, keys.Enter):
			if m.page == versionList && len(m.versions) > 0 {
				m.detail = &m.versions[m.selected]
				m.page = versionDetail
				m.fileCursor = 0
			}
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.page == versionList && m.selected > 0 {
				m.selected--
			}
			if m.page == versionDetail && m.fileCursor > 0 {
				m.fileCursor--
			}
			return m, nil

		case key.Matches(msg, keys.Down):
			if m.page == versionList && m.selected < len(m.versions)-1 {
				m.selected++
			}
			if m.page == versionDetail && m.detail != nil && m.fileCursor < len(m.detail.Files)-1 {
				m.fileCursor++
			}
			return m, nil

		case key.Matches(msg, keys.Open):
			if m.page == versionDetail && m.detail != nil && len(m.detail.Files) > 0 && m.fileCursor >= 0 && m.fileCursor < len(m.detail.Files) {
				openURL(m.detail.Files[m.fileCursor].URL)
			}
			return m, nil

		case key.Matches(msg, keys.Home):
			if m.page == versionList {
				m.selected = 0
			}
			return m, nil

		case key.Matches(msg, keys.End):
			if m.page == versionList {
				m.selected = len(m.versions) - 1
				if m.selected < 0 {
					m.selected = 0
				}
			}
			return m, nil
		}
	}

	return m, nil
}

func (m versionModel) View() string {
	var sb strings.Builder
	s := m.styles

	if m.loading {
		sp := spinnerFrame(m.frame)
		return fmt.Sprintf("\n  %s", s.Info.Render(fmt.Sprintf("%s loading versions...", sp)))
	}

	if m.err != nil {
		return fmt.Sprintf("\n  %s", s.Error.Render(fmt.Sprintf("error: %s", m.err)))
	}

	if m.page == versionDetail && m.detail != nil {
		return m.renderDetail()
	}

	sb.WriteString(s.Title.Render("versions"))
	sb.WriteString("\n\n")

	if len(m.versions) == 0 {
		fmt.Fprintf(&sb, "  %s", s.Info.Render("no versions found"))
		return sb.String()
	}

	maxResults := len(m.versions)
	height := 20
	start := m.selected - height/2
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
		v := m.versions[i]

		prefix := "  "
		if i == m.selected {
			prefix = "> "
		}

		typeStr := string(v.VersionType)
		if v.Featured {
			typeStr = "* " + typeStr
		}

		if i == m.selected {
			sb.WriteString(s.Selected.Render(fmt.Sprintf("%s%s", prefix, v.VersionNumber)))
			fmt.Fprintf(&sb, "  %s  %s  %s downloads\n",
				s.Accent.Render(typeStr),
				s.Info.Render(v.DatePublished.Format("2006-01-02")),
				s.Accent.Render(formatNumber(v.Downloads)))
			fmt.Fprintf(&sb, "  loaders: %s\n", s.Info.Render(strings.Join(v.Loaders, ", ")))
			gv := v.GameVersions
			if len(gv) > 6 {
				gv = append(gv[:6], "...")
			}
			fmt.Fprintf(&sb, "  versions: %s\n", s.Info.Render(strings.Join(gv, ", ")))
		} else {
			fmt.Fprintf(&sb, "%s%s", prefix, v.VersionNumber)
			fmt.Fprintf(&sb, " [%s]  %s  %s downloads\n",
				typeStr,
				v.DatePublished.Format("2006-01-02"),
				formatNumber(v.Downloads))
		}

		if i < end-1 {
			sb.WriteString("\n")
		}
	}

	if len(m.versions) > height {
		sb.WriteString(s.Dimmed.Render(scrollLine(m.selected, len(m.versions), height)))
	}

	sb.WriteString("\n\n")
	sb.WriteString(s.Dimmed.Render("j/k navigate - enter detail - esc back"))

	return sb.String()
}

func (m versionModel) renderDetail() string {
	var sb strings.Builder
	v := m.detail
	s := m.styles

	sb.WriteString(s.Title.Render(v.VersionNumber))
	sb.WriteString("\n\n")

	fmt.Fprintf(&sb, "  name: %s\n", s.Accent.Render(v.Name))
	fmt.Fprintf(&sb, "  type: %s\n", string(v.VersionType))
	if v.Featured {
		fmt.Fprintf(&sb, "  %s\n", s.Success.Render("* featured"))
	}
	fmt.Fprintf(&sb, "  date: %s\n", v.DatePublished.Format("2006-01-02 15:04"))
	fmt.Fprintf(&sb, "  downloads: %s\n", formatNumber(v.Downloads))
	fmt.Fprintf(&sb, "  loaders: %s\n", strings.Join(v.Loaders, ", "))
	fmt.Fprintf(&sb, "  versions: %s\n", strings.Join(v.GameVersions, ", "))
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
			fmt.Fprintf(&sb, "  %s: %s\n", dep.DependencyType, target)
		}
		sb.WriteString("\n")
	}

	if len(v.Files) > 0 {
		sb.WriteString(s.Section.Render("files"))
		sb.WriteString("\n")
		for i, f := range v.Files {
			primary := ""
			if f.Primary {
				primary = s.Accent.Render(" *")
			}
			urlText := hyperlink(f.URL, f.Filename)
			prefix := "  "
			if i == m.fileCursor {
				prefix = s.Accent.Render("> ")
			}
			fmt.Fprintf(&sb, "%s%s (%s)%s\n", prefix, urlText, formatBytes(f.Size), primary)
		}
		sb.WriteString("\n")
	}

	if m.mrpackInfo != nil {
		sb.WriteString(s.Section.Render("mrpack info"))
		sb.WriteString("\n")
		fmt.Fprintf(&sb, "  name: %s\n", m.mrpackInfo.Name)
		fmt.Fprintf(&sb, "  version id: %s\n", m.mrpackInfo.VersionID)
		if m.mrpackInfo.Summary != "" {
			fmt.Fprintf(&sb, "  summary: %s\n", m.mrpackInfo.Summary)
		}
		if mc, ok := m.mrpackInfo.Dependencies["minecraft"]; ok {
			fmt.Fprintf(&sb, "  game: %s\n", mc)
		}
		for k, v := range m.mrpackInfo.Dependencies {
			if k != "minecraft" {
				fmt.Fprintf(&sb, "  %s: %s\n", k, v)
			}
		}
		fmt.Fprintf(&sb, "  files: %d\n", len(m.mrpackInfo.Files))
		sb.WriteString("\n")
	}

	if v.Changelog != "" {
		sb.WriteString(s.Section.Render("changelog"))
		sb.WriteString("\n")
		if len(v.Changelog) > 500 {
			sb.WriteString(v.Changelog[:500] + "...\n")
		} else {
			sb.WriteString(v.Changelog + "\n")
		}
	}

	help := "esc back"
	if len(v.Files) > 0 {
		help = "up/down: select file - d: download - i: inspect mrpack - I: install - o: open url - " + help
	}
	sb.WriteString("\n")
	sb.WriteString(s.Dimmed.Render(help))

	return sb.String()
}
