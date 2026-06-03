package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type keymap struct {
	Quit     key.Binding
	Help     key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Enter    key.Binding
	Escape   key.Binding
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Home     key.Binding
	End      key.Binding
	Slash    key.Binding
	Backspce key.Binding
	Delete   key.Binding
	Right    key.Binding
	Left     key.Binding
	Space    key.Binding
	Theme    key.Binding
	Open     key.Binding
	Download key.Binding
	QuickDL  key.Binding
	Retry    key.Binding
	Cancel   key.Binding
	Inspect  key.Binding
	Install  key.Binding
	Clear    key.Binding
	Yes      key.Binding
}

var keys = keymap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next page"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev page"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup", "ctrl+u"),
		key.WithHelp("pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown", "ctrl+d"),
		key.WithHelp("pgdn", "page down"),
	),
	Home: key.NewBinding(
		key.WithKeys("home", "g"),
		key.WithHelp("g", "top"),
	),
	End: key.NewBinding(
		key.WithKeys("end", "G"),
		key.WithHelp("G", "bottom"),
	),
	Slash: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Backspce: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("bs", "delete"),
	),
	Delete: key.NewBinding(
		key.WithKeys("delete"),
		key.WithHelp("del", "delete"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→", "next option"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←", "prev option"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "retry"),
	),
	Theme: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "theme"),
	),
	Open: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open url / folder"),
	),
	Download: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "download"),
	),
	QuickDL: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "quick download latest"),
	),
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "cancel"),
	),
	Inspect: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "inspect mrpack"),
	),
	Install: key.NewBinding(
		key.WithKeys("I"),
		key.WithHelp("I", "install mrpack"),
	),
	Clear: key.NewBinding(
		key.WithKeys("C"),
		key.WithHelp("C", "clear cache"),
	),
	Yes: key.NewBinding(
		key.WithKeys("y", "Y"),
	),
}
