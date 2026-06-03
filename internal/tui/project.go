package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

	"github.com/programmersd21/mint/internal/api"
	"github.com/programmersd21/mint/internal/cache"
	"github.com/programmersd21/mint/internal/models"
)

const bodyVisibleLines = 20

type projectModel struct {
	client  *api.Client
	cache   *cache.Cache
	project *models.Project
	loading bool
	err     error
	page    projectPage

	showBody      bool
	bodyScroll    int
	renderedBody  string
	renderedLines []string
	frame         int
	width         int
	renderWidth   int
	styles        Styles

	urls      []urlEntry
	urlCursor int
}

type urlEntry struct {
	label string
	url   string
}

func newProjectModel(client *api.Client, c *cache.Cache) projectModel {
	return projectModel{
		client: client,
		cache:  c,
		page:   projectContent,
		styles: glob,
	}
}

func (m *projectModel) setStyles(s Styles) {
	m.styles = s
}

func (m projectModel) Init() tea.Cmd { return nil }

func (m *projectModel) load(projectID string) tea.Cmd {
	m.loading = true
	m.err = nil
	m.page = projectContent
	m.showBody = false
	m.bodyScroll = 0
	m.renderedBody = ""
	m.renderedLines = nil
	m.urls = nil
	m.urlCursor = 0

	cached, err := m.cache.GetProject(projectID)
	if err == nil && cached != nil {
		m.project = cached
		m.loading = false
		m.collectURLs()
		return nil
	}

	return func() tea.Msg {
		project, err := m.client.GetProject(projectID)
		if err != nil {
			return projectErrorMsg{err: err}
		}

		_ = m.cache.CacheProject(project)
		m.cache.AddRecentlyViewed("project", project.ID, project.Title)

		return projectLoadedMsg{project: project}
	}
}

type projectLoadedMsg struct {
	project *models.Project
}

type projectErrorMsg struct {
	err error
}

func (m *projectModel) collectURLs() {
	m.urls = nil
	if m.project == nil {
		return
	}
	if m.project.SourceURL != nil {
		m.urls = append(m.urls, urlEntry{label: "source", url: *m.project.SourceURL})
	}
	if m.project.IssuesURL != nil {
		m.urls = append(m.urls, urlEntry{label: "issues", url: *m.project.IssuesURL})
	}
	if m.project.WikiURL != nil {
		m.urls = append(m.urls, urlEntry{label: "wiki", url: *m.project.WikiURL})
	}
	if m.project.DiscordURL != nil {
		m.urls = append(m.urls, urlEntry{label: "discord", url: *m.project.DiscordURL})
	}
	for _, d := range m.project.DonationURLs {
		m.urls = append(m.urls, urlEntry{label: d.Platform, url: d.URL})
	}
}

func glamourStyle(s Styles) glamour.TermRendererOption {
	style := "dark"
	if s.T != nil {
		bg := s.T.Background
		if len(bg) >= 6 {
			h := strings.TrimPrefix(bg, "#")
			r, _ := strconv.ParseInt(h[0:2], 16, 64)
			g, _ := strconv.ParseInt(h[2:4], 16, 64)
			b, _ := strconv.ParseInt(h[4:6], 16, 64)
			if r+g+b > 382 {
				style = "light"
			}
		}
	}
	return glamour.WithStandardStyle(style)
}

func (m *projectModel) renderBody() {
	if m.project == nil || m.project.Body == "" {
		m.renderedBody = ""
		m.renderedLines = nil
		return
	}
	wrapWidth := m.width
	if wrapWidth < 40 {
		wrapWidth = 40
	}
	renderer, err := glamour.NewTermRenderer(
		glamourStyle(m.styles),
		glamour.WithWordWrap(wrapWidth),
	)
	if err != nil {
		m.renderedBody = ""
		m.renderedLines = nil
		return
	}
	rendered, err := renderer.Render(m.project.Body)
	if err != nil {
		m.renderedBody = ""
		m.renderedLines = nil
		return
	}
	m.renderedBody = rendered
	m.renderedLines = strings.Split(rendered, "\n")
	m.renderWidth = m.width
}

func (m projectModel) Update(msg tea.Msg) (projectModel, tea.Cmd) {
	switch msg := msg.(type) {
	case projectLoadedMsg:
		m.project = msg.project
		m.loading = false
		m.collectURLs()

	case projectErrorMsg:
		m.err = msg.err
		m.loading = false

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Escape):
			m.page = projectSearch
			m.project = nil
			return m, nil

		case key.Matches(msg, keys.Tab):
			m.showBody = !m.showBody
			m.bodyScroll = 0
			m.urlCursor = 0
			if m.showBody && m.renderedLines == nil && m.project != nil {
				m.renderBody()
			}
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.showBody && m.bodyScroll > 0 {
				m.bodyScroll--
			} else if !m.showBody && m.urlCursor > 0 {
				m.urlCursor--
			}
			return m, nil

		case key.Matches(msg, keys.Down):
			if m.showBody && m.bodyScroll < len(m.renderedLines)-bodyVisibleLines {
				m.bodyScroll++
			} else if !m.showBody && m.urlCursor < len(m.urls)-1 {
				m.urlCursor++
			}
			return m, nil

		case key.Matches(msg, keys.Home):
			if m.showBody {
				m.bodyScroll = 0
			} else {
				m.urlCursor = 0
			}
			return m, nil

		case key.Matches(msg, keys.End):
			if m.showBody {
				m.bodyScroll = len(m.renderedLines) - bodyVisibleLines
				if m.bodyScroll < 0 {
					m.bodyScroll = 0
				}
			} else {
				m.urlCursor = len(m.urls) - 1
				if m.urlCursor < 0 {
					m.urlCursor = 0
				}
			}
			return m, nil

		case key.Matches(msg, keys.Open):
			if !m.showBody && m.urlCursor >= 0 && m.urlCursor < len(m.urls) {
				openURL(m.urls[m.urlCursor].url)
			}
			return m, nil
		}
	}

	return m, nil
}

func (m projectModel) metadataView() string {
	var sb strings.Builder
	s := m.styles
	p := m.project

	sb.WriteString(s.Title.Render(p.Title))
	sb.WriteString("\n\n")

	sb.WriteString(s.Info.Render(p.Description))
	sb.WriteString("\n\n")

	sb.WriteString(s.Accent.Render(fmt.Sprintf("↓ %s  * %d",
		formatNumber(p.Downloads), p.Followers)))
	sb.WriteString("\n\n")

	if len(p.Loaders) > 0 {
		fmt.Fprintf(&sb, "loaders: %s\n", s.Primary.Render(strings.Join(p.Loaders, ", ")))
	}

	if len(p.GameVersions) > 0 {
		versions := p.GameVersions
		if len(versions) > 10 {
			versions = versions[:10]
		}
		fmt.Fprintf(&sb, "versions: %s\n", s.Cyan.Render(strings.Join(versions, ", ")))
		if len(p.GameVersions) > 10 {
			fmt.Fprintf(&sb, "  +%d more\n", len(p.GameVersions)-10)
		}
	}

	if len(p.Categories) > 0 {
		fmt.Fprintf(&sb, "categories: %s\n", s.Info.Render(strings.Join(p.Categories, ", ")))
	}

	fmt.Fprintf(&sb, "license: %s\n", s.Info.Render(p.License.Name))
	fmt.Fprintf(&sb, "sides: %s / %s\n", p.ClientSide, p.ServerSide)

	for i, u := range m.urls {
		urlText := hyperlink(u.url, u.url)
		prefix := "  "
		if i == m.urlCursor {
			prefix = s.Accent.Render("> ")
		}
		fmt.Fprintf(&sb, "%s%s: %s\n", prefix, s.Info.Render(u.label), urlText)
	}

	return sb.String()
}

func (m projectModel) bodyView() string {
	var sb strings.Builder
	s := m.styles

	if !m.showBody {
		sb.WriteString(s.Info.Render("tab to toggle description"))
		return sb.String()
	}

	if m.renderedLines == nil {
		return ""
	}

	sb.WriteString(s.Section.Render("description"))
	sb.WriteString("\n\n")

	start := m.bodyScroll
	end := start + bodyVisibleLines
	if end > len(m.renderedLines) {
		end = len(m.renderedLines)
	}

	for _, line := range m.renderedLines[start:end] {
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	if len(m.renderedLines) > bodyVisibleLines {
		sb.WriteString(s.Dimmed.Render(scrollLine(m.bodyScroll, len(m.renderedLines), bodyVisibleLines)))
	}

	return sb.String()
}

func (m projectModel) View() string {
	s := m.styles

	if m.loading {
		sp := spinnerFrame(m.frame)
		return fmt.Sprintf("\n  %s", s.Info.Render(fmt.Sprintf("%s loading project...", sp)))
	}

	if m.err != nil {
		return fmt.Sprintf("\n  %s", s.Error.Render(fmt.Sprintf("error: %s", m.err)))
	}

	if m.project == nil {
		return fmt.Sprintf("\n  %s", s.Info.Render("select a project to view"))
	}

	var sb strings.Builder
	sb.WriteString(m.metadataView())

	countStr := ""
	if len(m.urls) > 1 && !m.showBody {
		countStr = fmt.Sprintf(" (%d/%d)", m.urlCursor+1, len(m.urls))
	}

	if m.project.Body != "" {
		sb.WriteString("\n")
		sb.WriteString(m.bodyView())
	}

	sb.WriteString("\n\n")
	help := "tab: toggle body"
	if !m.showBody && len(m.urls) > 0 {
		help += fmt.Sprintf(" - o: open url%s", countStr)
	}
	if m.showBody {
		help += " - up/down: scroll"
	}
	help += " - esc: back"
	sb.WriteString(s.Dimmed.Render(help))

	return sb.String()
}
