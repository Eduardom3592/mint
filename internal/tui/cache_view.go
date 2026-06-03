package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"

	"github.com/programmersd21/mint/internal/cache"
)

type recentItem struct {
	EntityType string
	EntityID   string
	Title      string
	ViewedAt   time.Time
}

type cacheViewModel struct {
	cache    *cache.Cache
	cursor   int
	recent   []recentItem
	loaded   bool
	styles   Styles
	askClear bool
	clearMsg string
}

func newCacheViewModel(c *cache.Cache) cacheViewModel {
	return cacheViewModel{
		cache:  c,
		cursor: 0,
		styles: glob,
	}
}

func (m *cacheViewModel) setStyles(s Styles) {
	m.styles = s
}

func (m cacheViewModel) Init() tea.Cmd { return m.load() }

func (m cacheViewModel) load() tea.Cmd {
	return func() tea.Msg {
		raw, err := m.cache.GetRecentlyViewed("project", 20)
		if err != nil {
			return cacheLoadedMsg{err: err}
		}
		items := make([]recentItem, len(raw))
		for i, r := range raw {
			items[i] = recentItem{
				EntityType: r.EntityType,
				EntityID:   r.EntityID,
				Title:      r.Title,
				ViewedAt:   r.ViewedAt,
			}
		}
		return cacheLoadedMsg{recent: items}
	}
}

type cacheLoadedMsg struct {
	recent []recentItem
	err    error
}

func (m cacheViewModel) Update(msg tea.Msg) (cacheViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case cacheLoadedMsg:
		m.loaded = true
		if msg.err == nil {
			m.recent = msg.recent
			if m.cursor >= len(m.recent) {
				m.cursor = 0
			}
		}
		return m, nil

	case tea.KeyMsg:
		if m.askClear {
			switch {
			case key.Matches(msg, keys.Enter) || key.Matches(msg, keys.Yes):
				m.askClear = false
				m.clearMsg = ""
				if m.cache != nil {
					_ = m.cache.ClearCache()
					m.loaded = false
					m.recent = nil
					m.cursor = 0
					return m, m.load()
				}
				return m, nil
			default:
				m.askClear = false
				m.clearMsg = ""
				return m, nil
			}
		}

		switch {
		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.recent)-1 {
				m.cursor++
			}
			return m, nil

		case key.Matches(msg, keys.Enter):
			if m.cursor >= 0 && m.cursor < len(m.recent) {
				item := m.recent[m.cursor]
				return m, func() tea.Msg {
					return openProjectMsg{id: item.EntityID}
				}
			}
			return m, nil

		case key.Matches(msg, keys.Clear):
			m.askClear = true
			m.clearMsg = "clear all cached data? (enter=yes, any key=no)"
			return m, nil
		}
	}
	return m, nil
}

func (m cacheViewModel) View() string {
	var sb strings.Builder
	s := m.styles

	sb.WriteString(s.Title.Render("cache"))
	sb.WriteString("\n\n")

	projects, versions, cacheSize, err := m.cache.Stats()
	if err != nil {
		return fmt.Sprintf("\n  %s", s.Error.Render(fmt.Sprintf("error: %s", err)))
	}

	sb.WriteString(s.Section.Render("stats"))
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "  %d projects\n", projects)
	fmt.Fprintf(&sb, "  %d versions\n", versions)
	fmt.Fprintf(&sb, "  %s\n", formatBytes(cacheSize))
	sb.WriteString("\n")

	if m.loaded && len(m.recent) > 0 {
		sb.WriteString(s.Section.Render("recently viewed"))
		sb.WriteString("\n\n")

		for i, item := range m.recent {
			prefix := "  "
			if i == m.cursor {
				prefix = s.Accent.Render("> ")
			}

			title := item.Title
			if len(title) > 50 {
				title = title[:47] + "..."
			}

			if i == m.cursor {
				sb.WriteString(s.Selected.Render(prefix + title))
			} else {
				sb.WriteString(prefix + title)
			}
			sb.WriteString("  " + s.Dimmed.Render(item.ViewedAt.Format("2006-01-02 15:04")) + "\n")
		}

		sb.WriteString("\n")
		sb.WriteString(s.Dimmed.Render("j/k navigate - enter open - C clear cache"))
		sb.WriteString("\n\n")
	} else if !m.loaded {
		sb.WriteString(s.Info.Render("  loading..."))
		sb.WriteString("\n\n")
	}

	if m.askClear {
		fmt.Fprintf(&sb, "  %s\n\n", s.Warning.Render(m.clearMsg))
	}

	sb.WriteString(s.Dimmed.Render("note: cache is stored in sqlite for offline access"))
	sb.WriteString("\n")
	sb.WriteString(s.Dimmed.Render("all data is read from modrinth and cached locally"))

	return sb.String()
}
