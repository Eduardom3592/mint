package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/programmersd21/mint/internal/tui/theme"
)

type themeSwitcherModel struct {
	manager *theme.Manager
	themes  []theme.Theme
	cursor  int
	visible bool
	styles  Styles
}

func newThemeSwitcherModel(mgr *theme.Manager) themeSwitcherModel {
	return themeSwitcherModel{
		manager: mgr,
		themes:  mgr.List(),
		cursor:  0,
		visible: false,
		styles:  glob,
	}
}

func (m *themeSwitcherModel) setStyles(s Styles) {
	m.styles = s
}

func (m themeSwitcherModel) Init() tea.Cmd { return nil }

func (m themeSwitcherModel) Update(msg tea.Msg) (themeSwitcherModel, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Escape), key.Matches(msg, keys.Theme):
			m.visible = false
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < len(m.themes) {
					t := &m.themes[m.cursor]
					setGlobalTheme(t)
					return m, m.broadcastThemeChange(t)
				}
			}
			return m, nil

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.themes)-1 {
				m.cursor++
				if m.cursor < len(m.themes) {
					t := &m.themes[m.cursor]
					setGlobalTheme(t)
					return m, m.broadcastThemeChange(t)
				}
			}
			return m, nil

		case key.Matches(msg, keys.Enter):
			if m.cursor < len(m.themes) {
				t := &m.themes[m.cursor]
				_, _ = m.manager.Set(t.Name)
				setGlobalTheme(t)
				m.visible = false
				return m, m.broadcastThemeChange(t)
			}
			return m, nil
		}
	}

	return m, nil
}

func (m themeSwitcherModel) broadcastThemeChange(t *theme.Theme) tea.Cmd {
	return func() tea.Msg {
		return themeChangedMsg{theme: t}
	}
}

type themeChangedMsg struct {
	theme *theme.Theme
}

func (m themeSwitcherModel) View() string {
	if !m.visible {
		return ""
	}

	s := m.styles
	var sb strings.Builder

	panelW := 56

	sb.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.T.BorderColor()).
		Padding(1, 2).
		Width(panelW).
		Render(m.renderPanelContent()))

	return sb.String()
}

func (m themeSwitcherModel) renderPanelContent() string {
	s := m.styles
	var sb strings.Builder

	sb.WriteString(s.T.PrimaryStyle().Render("theme switcher"))
	sb.WriteString("\n\n")

	for i, t := range m.themes {
		current := m.manager.Current()
		isCurrent := current != nil && current.Name == t.Name

		prefix := "  "
		if i == m.cursor {
			prefix = s.T.ActiveItemStyle().Render(">")
		}

		mark := " "
		if isCurrent {
			mark = s.T.SuccessStyle().Render("o")
		}

		line := fmt.Sprintf("%s %s %s", prefix, mark, t.Name)

		if i == m.cursor {
			sb.WriteString(s.T.SelectedItemStyle().Render(line))
		} else if isCurrent {
			sb.WriteString(s.T.PrimaryStyle().Render(line))
		} else {
			sb.WriteString(line)
		}
		sb.WriteString("\n")
	}

	if m.cursor < len(m.themes) {
		t := &m.themes[m.cursor]
		sb.WriteString("\n")
		sb.WriteString(s.T.SeparatorText())
		sb.WriteString("\n\n")
		sb.WriteString(m.renderColorSwatches(t))
	}

	sb.WriteString("\n")
	sb.WriteString(s.T.DimmedStyle().Render("up/down browse - enter confirm - esc cancel"))

	return sb.String()
}

func (m themeSwitcherModel) renderColorSwatches(t *theme.Theme) string {
	var sb strings.Builder

	swatches := []struct {
		name  string
		color string
	}{
		{"primary", t.Primary},
		{"secondary", t.Secondary},
		{"accent", t.Accent},
		{"success", t.Success},
		{"warning", t.Warning},
		{"error", t.Error},
		{"surface", t.Surface},
	}

	barW := 20

	sb.WriteString(t.PrimaryStyle().Render("  colors"))
	sb.WriteString("\n\n")

	for _, sw := range swatches {
		bar := lipgloss.NewStyle().
			Foreground(lipgloss.Color(sw.color)).
			Render(strings.Repeat("-", barW))

		name := lipgloss.NewStyle().Foreground(lipgloss.Color(sw.color)).Render(sw.name)
		fmt.Fprintf(&sb, "    %-12s %s\n", name, bar)
	}

	return sb.String()
}
