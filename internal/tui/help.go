package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
)

type helpModel struct {
	styles Styles
}

func newHelpModel() helpModel {
	return helpModel{styles: glob}
}

func (h *helpModel) setStyles(s Styles) {
	h.styles = s
}

func (h helpModel) Init() tea.Cmd { return nil }

func (h helpModel) Update(msg tea.Msg) (helpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Escape), key.Matches(msg, keys.Quit), key.Matches(msg, keys.Help):
			return h, nil
		}
	}
	return h, nil
}

func (h helpModel) View() string {
	var sb strings.Builder
	s := h.styles

	sb.WriteString(s.Title.Render("help"))
	sb.WriteString("\n\n")
	sb.WriteString(s.Info.Render("mint - modrinth terminal client"))
	sb.WriteString("\n\n")

	sb.WriteString(s.Section.Render("navigation"))
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("j/k  or  up/down"), "navigate lists")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("g / G"), "top / bottom")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("/"), "focus search")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("tab / S-tab"), "switch pages")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("t"), "theme switcher")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("enter"), "select")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("esc"), "back")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("?"), "help")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("q / C-c"), "quit")

	sb.WriteString("\n")
	sb.WriteString(s.Section.Render("actions"))
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("d"), "download (version picker)")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("D"), "quick download latest")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("i"), "inspect mrpack")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("I"), "download & install mrpack")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("r"), "retry failed download")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("c"), "cancel download")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("o"), "open download folder")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("del"), "remove from history")

	sb.WriteString("\n")
	sb.WriteString(s.Section.Render("settings"))
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("j/k  or  up/down"), "select a setting")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("enter"), "edit the value")
	fmt.Fprintf(&sb, "  %-20s %s\n", s.Accent.Render("enter / esc"), "confirm and save")

	sb.WriteString("\n")
	sb.WriteString(s.Section.Render("pages"))
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "  %s  %s\n", s.Dimmed.Render("home"), "discover the top 4 downloaded mods")
	fmt.Fprintf(&sb, "  %s  %s\n", s.Dimmed.Render("search"), "browse & search modrinth")
	fmt.Fprintf(&sb, "  %s  %s\n", s.Dimmed.Render("downloads"), "active / history / installed modpacks / failed downloads")
	fmt.Fprintf(&sb, "  %s  %s\n", s.Dimmed.Render("cache"), "browse cached data + reopen projects")
	fmt.Fprintf(&sb, "  %s  %s\n", s.Dimmed.Render("settings"), "configure app (editable)")

	sb.WriteString("\n")
	sb.WriteString(s.Info.Render("esc / ? to close"))

	return sb.String()
}
