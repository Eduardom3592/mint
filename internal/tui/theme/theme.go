package theme

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name string `json:"name"`

	// Core
	Background string `json:"background"`
	Foreground string `json:"foreground"`
	Surface    string `json:"surface"`

	// Accent hierarchy
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Accent    string `json:"accent"`

	// Semantic
	Success string `json:"success"`
	Warning string `json:"warning"`
	Error   string `json:"error"`
	Info    string `json:"info"`

	// UI chrome
	Border    string `json:"border"`
	Highlight string `json:"highlight"`
	Selection string `json:"selection"`
	StatusBar string `json:"status_bar"`
	Sidebar   string `json:"sidebar"`

	// Text
	Text      string `json:"text"`
	TextDim   string `json:"text_dim"`
	TextMuted string `json:"text_muted"`
}

func (t *Theme) lip(c string) lipgloss.Color {
	return lipgloss.Color(c)
}

func (t *Theme) BgColor() lipgloss.Color        { return t.lip(t.Background) }
func (t *Theme) FgColor() lipgloss.Color        { return t.lip(t.Foreground) }
func (t *Theme) SurfaceColor() lipgloss.Color   { return t.lip(t.Surface) }
func (t *Theme) PrimaryColor() lipgloss.Color   { return t.lip(t.Primary) }
func (t *Theme) SecondaryColor() lipgloss.Color { return t.lip(t.Secondary) }
func (t *Theme) AccentColor() lipgloss.Color    { return t.lip(t.Accent) }
func (t *Theme) SuccessColor() lipgloss.Color   { return t.lip(t.Success) }
func (t *Theme) WarningColor() lipgloss.Color   { return t.lip(t.Warning) }
func (t *Theme) ErrorColor() lipgloss.Color     { return t.lip(t.Error) }
func (t *Theme) InfoColor() lipgloss.Color      { return t.lip(t.Info) }
func (t *Theme) BorderColor() lipgloss.Color    { return t.lip(t.Border) }
func (t *Theme) HighlightColor() lipgloss.Color { return t.lip(t.Highlight) }
func (t *Theme) SelectionColor() lipgloss.Color { return t.lip(t.Selection) }
func (t *Theme) StatusBarColor() lipgloss.Color { return t.lip(t.StatusBar) }
func (t *Theme) SidebarColor() lipgloss.Color   { return t.lip(t.Sidebar) }
func (t *Theme) TextColor() lipgloss.Color      { return t.lip(t.Text) }
func (t *Theme) TextDimColor() lipgloss.Color   { return t.lip(t.TextDim) }
func (t *Theme) TextMutedColor() lipgloss.Color { return t.lip(t.TextMuted) }

func (t *Theme) AppStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (t *Theme) SidebarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		MarginRight(1).
		Padding(0, 1, 0, 1)
}

func (t *Theme) ContentStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 1)
}

func (t *Theme) StatusBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Height(1).
		Foreground(t.TextDimColor()).
		Padding(0, 1)
}

func (t *Theme) TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.PrimaryColor())
}

func (t *Theme) ActiveItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.PrimaryColor()).
		Bold(true).
		Padding(0, 1)
}

func (t *Theme) InactiveItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.TextDimColor()).
		Padding(0, 1)
}

func (t *Theme) SelectedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.PrimaryColor()).
		Bold(true).
		Padding(0, 1)
}

func (t *Theme) InfoStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.TextDimColor())
}

func (t *Theme) AccentStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.AccentColor()).Bold(true)
}

func (t *Theme) DimmedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.TextMutedColor())
}

func (t *Theme) SuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.SuccessColor()).Bold(true)
}

func (t *Theme) ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.ErrorColor()).Bold(true)
}

func (t *Theme) WarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.WarningColor()).Bold(true)
}

func (t *Theme) BoldStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true)
}

func (t *Theme) LabelStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.TextMutedColor()).Width(18)
}

func (t *Theme) ValueStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.TextColor())
}

func (t *Theme) SectionStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(t.SecondaryColor()).MarginTop(1).MarginBottom(1)
}

func (t *Theme) SecondaryStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.SecondaryColor())
}

func (t *Theme) PrimaryStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.PrimaryColor())
}

func (t *Theme) CyanStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.lip(t.Info))
}

func (t *Theme) SeparatorText() string {
	return lipgloss.NewStyle().Foreground(t.BorderColor()).Render("─────")
}

func (t *Theme) ProgressBar(progress float64, width int) string {
	if width < 1 {
		width = 30
	}
	if progress > 1 {
		progress = 1
	}
	if progress < 0 {
		progress = 0
	}
	filled := int(float64(width) * progress)
	empty := width - filled

	bar := lipgloss.NewStyle().Foreground(t.PrimaryColor()).Render(
		fmt.Sprintf("%s%s",
			strings.Repeat("█", filled),
			strings.Repeat("░", empty),
		),
	)

	percent := int(progress * 100)
	var pctColor lipgloss.Color
	if progress >= 1.0 {
		pctColor = t.SuccessColor()
	} else if progress > 0.5 {
		pctColor = t.WarningColor()
	} else {
		pctColor = t.PrimaryColor()
	}

	return fmt.Sprintf("%s %s", bar, lipgloss.NewStyle().Foreground(pctColor).Render(fmt.Sprintf("%3d%%", percent)))
}
