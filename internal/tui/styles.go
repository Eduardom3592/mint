package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/programmersd21/mint/internal/tui/theme"
)

type Styles struct {
	App        lipgloss.Style
	Sidebar    lipgloss.Style
	Content    lipgloss.Style
	StatusBar  lipgloss.Style
	Title      lipgloss.Style
	ActiveItem lipgloss.Style
	Inactive   lipgloss.Style
	Selected   lipgloss.Style
	Info       lipgloss.Style
	Accent     lipgloss.Style
	Dimmed     lipgloss.Style
	Success    lipgloss.Style
	Error      lipgloss.Style
	Warning    lipgloss.Style
	Bold       lipgloss.Style
	Label      lipgloss.Style
	Value      lipgloss.Style
	Section    lipgloss.Style
	Secondary  lipgloss.Style
	Primary    lipgloss.Style
	Cyan       lipgloss.Style
	Separator  string
	T          *theme.Theme
}

func NewStyles(t *theme.Theme) Styles {
	if t == nil {
		t = theme.Default()
	}
	return Styles{
		App:        t.AppStyle(),
		Sidebar:    t.SidebarStyle(),
		Content:    t.ContentStyle(),
		StatusBar:  t.StatusBarStyle(),
		Title:      t.TitleStyle(),
		ActiveItem: t.ActiveItemStyle(),
		Inactive:   t.InactiveItemStyle(),
		Selected:   t.SelectedItemStyle(),
		Info:       t.InfoStyle(),
		Accent:     t.AccentStyle(),
		Dimmed:     t.DimmedStyle(),
		Success:    t.SuccessStyle(),
		Error:      t.ErrorStyle(),
		Warning:    t.WarningStyle(),
		Bold:       t.BoldStyle(),
		Label:      t.LabelStyle(),
		Value:      t.ValueStyle(),
		Section:    t.SectionStyle(),
		Secondary:  t.SecondaryStyle(),
		Primary:    t.PrimaryStyle(),
		Cyan:       t.CyanStyle(),
		Separator:  t.SeparatorText(),
		T:          t,
	}
}

var glob Styles

func init() {
	glob = NewStyles(theme.Default())
}

func setGlobalTheme(t *theme.Theme) {
	glob = NewStyles(t)
}

func ProgressBar(progress float64, width int) string {
	return glob.T.ProgressBar(progress, width)
}

func formatNumber(n int) string {
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%.1fk", float64(n)/1_000)
	}
	return fmt.Sprintf("%d", n)
}

func formatBytes(b int64) string {
	if b < 1024 {
		return fmt.Sprintf("%d B", b)
	}
	if b < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(b)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024))
}

func truncateText(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
