package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"

	"github.com/programmersd21/mint/internal/cache"
)

type settingsViewModel struct {
	cache    *cache.Cache
	cursor   int
	inputs   []textinput.Model
	focus    int
	styles   Styles
	askReset bool
	resetMsg string
}

func newSettingsViewModel(c *cache.Cache) settingsViewModel {
	m := settingsViewModel{
		cache:  c,
		cursor: 0,
		focus:  -1,
		styles: glob,
	}

	dir, key, workers := m.loadSettings()

	dirInput := textinput.New()
	dirInput.Placeholder = "download directory"
	dirInput.SetValue(dir)
	dirInput.Prompt = ""
	dirInput.CharLimit = 200

	keyInput := textinput.New()
	keyInput.Placeholder = "modrinth api key"
	keyInput.SetValue(key)
	keyInput.Prompt = ""
	keyInput.CharLimit = 100
	keyInput.EchoMode = textinput.EchoPassword
	keyInput.EchoCharacter = '*'

	workersInput := textinput.New()
	workersInput.Placeholder = "max concurrent downloads (1-10)"
	workersInput.SetValue(workers)
	workersInput.Prompt = ""
	workersInput.CharLimit = 2

	m.inputs = []textinput.Model{dirInput, workersInput, keyInput}

	return m
}

func (m *settingsViewModel) loadSettings() (dir, key, workers string) {
	dir = "downloads"
	workers = "3"
	if m.cache != nil {
		if v, _ := m.cache.GetSetting("download_dir"); v != "" {
			dir = v
		}
		if v, _ := m.cache.GetSetting("api_key"); v != "" {
			key = v
		}
		if v, _ := m.cache.GetSetting("max_workers"); v != "" {
			workers = v
		}
	}
	return
}

func (m *settingsViewModel) setStyles(s Styles) {
	m.styles = s
}

func (m settingsViewModel) saveAll() {
	if m.cache == nil {
		return
	}
	_ = m.cache.SetSetting("download_dir", m.inputs[0].Value())
	_ = m.cache.SetSetting("max_workers", m.inputs[1].Value())
	if v := m.inputs[2].Value(); v != "" {
		_ = m.cache.SetSetting("api_key", v)
	}
}

func (m *settingsViewModel) resetDefaults() {
	_ = m.cache.ResetSettings()
	dir, key, workers := m.loadSettings()
	m.inputs[0].SetValue(dir)
	m.inputs[1].SetValue(workers)
	m.inputs[2].SetValue(key)
}

func (m settingsViewModel) Init() tea.Cmd { return nil }

func (m settingsViewModel) Update(msg tea.Msg) (settingsViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.askReset {
			switch {
			case key.Matches(msg, keys.Enter) || key.Matches(msg, keys.Yes):
				m.askReset = false
				m.resetMsg = ""
				m.resetDefaults()
				return m, nil
			default:
				m.askReset = false
				m.resetMsg = ""
				return m, nil
			}
		}

		if m.focus >= 0 {
			switch {
			case key.Matches(msg, keys.Escape):
				m.saveAll()
				m.inputs[m.focus].Blur()
				m.focus = -1
				return m, nil

			case key.Matches(msg, keys.Enter):
				m.saveAll()
				m.inputs[m.focus].Blur()
				m.focus = -1
				return m, nil
			}

			var cmd tea.Cmd
			m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
			return m, cmd
		}

		switch {
		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case key.Matches(msg, keys.Down):
			if m.cursor < 3 {
				m.cursor++
			}
			return m, nil

		case key.Matches(msg, keys.Enter):
			if m.cursor == 3 {
				m.askReset = true
				m.resetMsg = "reset all settings to defaults? (enter=yes, any key=no)"
				return m, nil
			}
			m.focus = m.cursor
			m.inputs[m.focus].Focus()
			m.inputs[m.focus].SetValue(m.inputs[m.focus].Value())
			m.inputs[m.focus].SetCursor(len(m.inputs[m.focus].Value()))
			return m, nil
		}
	}
	return m, nil
}

func (m settingsViewModel) View() string {
	var sb strings.Builder
	s := m.styles

	sb.WriteString(s.Title.Render("settings"))
	sb.WriteString("\n\n")

	labels := []string{"download directory", "max concurrent downloads", "api key"}

	for i, label := range labels {
		prefix := "  "
		if i == m.cursor && m.focus < 0 && !m.askReset {
			prefix = s.Accent.Render("> ")
		}

		if m.focus == i {
			fmt.Fprintf(&sb, "%s%s: %s\n", prefix, s.Accent.Render(label), m.inputs[i].View())
			sb.WriteString("\n")
		} else {
			val := m.inputs[i].Value()
			if i == 2 && len(val) > 0 {
				val = maskAPIKey(val)
			}
			fmt.Fprintf(&sb, "%s%s: %s\n", prefix, s.Accent.Render(label), val)
		}
	}

	sb.WriteString("\n")

	resetPrefix := "  "
	if m.cursor == 3 && m.focus < 0 && !m.askReset {
		resetPrefix = s.Accent.Render("> ")
	}
	fmt.Fprintf(&sb, "%s%s\n", resetPrefix, s.Warning.Render("reset all settings"))

	if m.askReset {
		fmt.Fprintf(&sb, "  %s\n", s.Warning.Render(m.resetMsg))
	}
	sb.WriteString("\n")

	if m.askReset {
		sb.WriteString(s.Info.Render("enter to confirm / any key to cancel"))
	} else if m.focus >= 0 {
		sb.WriteString(s.Info.Render("enter to confirm / esc to cancel"))
	} else {
		sb.WriteString(s.Info.Render("j/k navigate - enter to edit"))
	}
	sb.WriteString("\n\n")

	sb.WriteString(s.Section.Render("about"))
	sb.WriteString("\n\n")
	fmt.Fprintf(&sb, "  %s v%s\n", s.Accent.Render("mint"), "0.1.0")
	fmt.Fprintf(&sb, "  %s\n", "~/.local/share/mint")
	sb.WriteString("\n")
	sb.WriteString(s.Info.Render("a terminal client for modrinth"))

	return sb.String()
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return key
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}
