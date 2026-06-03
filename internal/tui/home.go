package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/programmersd21/mint/internal/api"
	"github.com/programmersd21/mint/internal/models"
)

type homeModel struct {
	client            *api.Client
	hits              []models.SearchHit
	loading           bool
	err               error
	frame             int
	width             int
	cursor            int
	selectedProjectID string
	styles            Styles
}

func newHomeModel(client *api.Client) homeModel {
	return homeModel{
		client: client,
		styles: glob,
	}
}

func (m *homeModel) setStyles(s Styles) {
	m.styles = s
}

func (m homeModel) Init() tea.Cmd {
	return m.fetchTop()
}

func (m *homeModel) fetchTop() tea.Cmd {
	m.loading = true
	m.err = nil
	return func() tea.Msg {
		filter := api.SearchFilter{
			Sort:        api.SortDownloads,
			Limit:       4,
			ProjectType: "mod",
		}
		result, err := m.client.Search(filter)
		if err != nil {
			return homeErrorMsg{err: err}
		}
		return homeLoadedMsg{hits: result.Hits}
	}
}

type homeLoadedMsg struct {
	hits []models.SearchHit
}

type homeErrorMsg struct {
	err error
}

func (m homeModel) Update(msg tea.Msg) (homeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case homeLoadedMsg:
		m.hits = msg.hits
		m.loading = false
		return m, nil

	case homeErrorMsg:
		m.err = msg.err
		m.loading = false
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			if len(m.hits) > 0 && m.cursor >= 0 && m.cursor < len(m.hits) {
				m.selectedProjectID = m.hits[m.cursor].ProjectID
			}
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.hits)-1 {
				m.cursor++
			}
			return m, nil

		case key.Matches(msg, keys.Home):
			m.cursor = 0
			return m, nil

		case key.Matches(msg, keys.End):
			m.cursor = len(m.hits) - 1
			if m.cursor < 0 {
				m.cursor = 0
			}
			return m, nil
		}
	}

	return m, nil
}

func (m homeModel) View() string {
	var sb strings.Builder
	s := m.styles

	if m.loading {
		sp := spinnerFrame(m.frame)
		return fmt.Sprintf("\n  %s", s.Info.Render(fmt.Sprintf("%s loading top mods...", sp)))
	}

	if m.err != nil {
		return fmt.Sprintf("\n  %s", s.Error.Render(fmt.Sprintf("error: %s", m.err)))
	}

	if len(m.hits) == 0 {
		return fmt.Sprintf("\n  %s", s.Info.Render("no mods found"))
	}

	sb.WriteString(s.Title.Render("top downloads"))
	sb.WriteString("\n\n")

	cw := m.width
	cardW := cw - 4
	if cardW < 20 {
		cardW = 20
	}

	borderColor := lipgloss.Color("#7c3aed")
	accentColor := lipgloss.Color("#a78bfa")
	if s.T != nil {
		borderColor = s.T.PrimaryColor()
		accentColor = s.T.AccentColor()
	}

	baseStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 2).
		Width(cardW)

	cardStyle := baseStyle.BorderForeground(borderColor)
	selectedStyle := baseStyle.BorderForeground(accentColor)

	for i, hit := range m.hits {
		style := cardStyle
		if i == m.cursor {
			style = selectedStyle
		}

		title := hit.Title
		dls := fmt.Sprintf("#%d  ↓ %s", i+1, formatNumber(hit.Downloads))
		desc := truncateText(hit.Description, cardW-8)
		innerW := cardW - 6
		pad := innerW - len(title) - len(dls)
		if pad < 1 {
			pad = 1
		}

		body := fmt.Sprintf("%s%s%s\n%s",
			s.Bold.Render(title),
			strings.Repeat(" ", pad),
			s.Info.Render(dls),
			s.Dimmed.Render(desc),
		)

		sb.WriteString(style.Render(body))
		sb.WriteString("\n\n")
	}

	sb.WriteString(s.Dimmed.Render("enter: open - d: download - D: quick dl - j/k: navigate - refreshes on startup"))

	return sb.String()
}
