package theme

func Default() *Theme {
	t := Mint()
	return &t
}

func BuiltInThemes() []Theme {
	return []Theme{
		Mint(),
		CatppuccinMocha(),
		CatppuccinLatte(),
		TokyoNight(),
		GruvboxDark(),
		GruvboxLight(),
		Nord(),
		Dracula(),
		Kanagawa(),
		RosePine(),
	}
}

func Mint() Theme {
	return Theme{
		Name: "Mint",

		Background: "#1A1A2E",
		Foreground: "#F1F1F6",
		Surface:    "#16213E",

		Primary:   "#FF6B9D",
		Secondary: "#C084FC",
		Accent:    "#FFB84D",

		Success: "#4ADE80",
		Warning: "#FFB84D",
		Error:   "#FF4D6D",
		Info:    "#67E8F9",

		Border:    "#2D2B55",
		Highlight: "#FF8FAB",
		Selection: "#1F2A47",
		StatusBar: "#16213E",
		Sidebar:   "#16213E",

		Text:      "#F1F1F6",
		TextDim:   "#8B8FA3",
		TextMuted: "#5C6080",
	}
}

func CatppuccinMocha() Theme {
	return Theme{
		Name: "Catppuccin Mocha",

		Background: "#1E1E2E",
		Foreground: "#CDD6F4",
		Surface:    "#181825",

		Primary:   "#F5C2E7",
		Secondary: "#CBA6F7",
		Accent:    "#FAB387",

		Success: "#A6E3A1",
		Warning: "#F9E2AF",
		Error:   "#F38BA8",
		Info:    "#89DCEB",

		Border:    "#45475A",
		Highlight: "#F5C2E7",
		Selection: "#313244",
		StatusBar: "#181825",
		Sidebar:   "#181825",

		Text:      "#CDD6F4",
		TextDim:   "#A6ADC8",
		TextMuted: "#6C7086",
	}
}

func CatppuccinLatte() Theme {
	return Theme{
		Name: "Catppuccin Latte",

		Background: "#EFF1F5",
		Foreground: "#4C4F69",
		Surface:    "#E6E9EF",

		Primary:   "#DD7878",
		Secondary: "#8839EF",
		Accent:    "#FE640B",

		Success: "#40A02B",
		Warning: "#DF8E1D",
		Error:   "#D20F39",
		Info:    "#04A5E5",

		Border:    "#BCC0CC",
		Highlight: "#DD7878",
		Selection: "#CCD0DA",
		StatusBar: "#E6E9EF",
		Sidebar:   "#E6E9EF",

		Text:      "#4C4F69",
		TextDim:   "#5C5F77",
		TextMuted: "#8C8FA1",
	}
}

func TokyoNight() Theme {
	return Theme{
		Name: "Tokyo Night",

		Background: "#1A1B26",
		Foreground: "#C0CAF5",
		Surface:    "#1F2335",

		Primary:   "#7AA2F7",
		Secondary: "#BB9AF7",
		Accent:    "#E0AF68",

		Success: "#9ECE6A",
		Warning: "#E0AF68",
		Error:   "#F7768E",
		Info:    "#7DCFFF",

		Border:    "#3B4261",
		Highlight: "#7AA2F7",
		Selection: "#24283B",
		StatusBar: "#1F2335",
		Sidebar:   "#1F2335",

		Text:      "#C0CAF5",
		TextDim:   "#A9B1D6",
		TextMuted: "#565F89",
	}
}

func GruvboxDark() Theme {
	return Theme{
		Name: "Gruvbox Dark",

		Background: "#282828",
		Foreground: "#EBDBB2",
		Surface:    "#1D2021",

		Primary:   "#FB4934",
		Secondary: "#D3869B",
		Accent:    "#D79921",

		Success: "#98971A",
		Warning: "#D79921",
		Error:   "#CC241D",
		Info:    "#458588",

		Border:    "#504945",
		Highlight: "#FB4934",
		Selection: "#3C3836",
		StatusBar: "#1D2021",
		Sidebar:   "#1D2021",

		Text:      "#EBDBB2",
		TextDim:   "#D5C4A1",
		TextMuted: "#7C6F64",
	}
}

func GruvboxLight() Theme {
	return Theme{
		Name: "Gruvbox Light",

		Background: "#FBF1C7",
		Foreground: "#3C3836",
		Surface:    "#F2E5BC",

		Primary:   "#9D0006",
		Secondary: "#8F3F71",
		Accent:    "#B57614",

		Success: "#79740E",
		Warning: "#B57614",
		Error:   "#9D0006",
		Info:    "#076678",

		Border:    "#D5C4A1",
		Highlight: "#9D0006",
		Selection: "#EBDAB2",
		StatusBar: "#F2E5BC",
		Sidebar:   "#F2E5BC",

		Text:      "#3C3836",
		TextDim:   "#504945",
		TextMuted: "#928374",
	}
}

func Nord() Theme {
	return Theme{
		Name: "Nord",

		Background: "#2E3440",
		Foreground: "#ECEFF4",
		Surface:    "#3B4252",

		Primary:   "#88C0D0",
		Secondary: "#B48EAD",
		Accent:    "#EBCB8B",

		Success: "#A3BE8C",
		Warning: "#EBCB8B",
		Error:   "#BF616A",
		Info:    "#81A1C1",

		Border:    "#4C566A",
		Highlight: "#88C0D0",
		Selection: "#434C5E",
		StatusBar: "#3B4252",
		Sidebar:   "#3B4252",

		Text:      "#ECEFF4",
		TextDim:   "#D8DEE9",
		TextMuted: "#7B88A1",
	}
}

func Dracula() Theme {
	return Theme{
		Name: "Dracula",

		Background: "#282A36",
		Foreground: "#F8F8F2",
		Surface:    "#21222C",

		Primary:   "#FF79C6",
		Secondary: "#BD93F9",
		Accent:    "#F1FA8C",

		Success: "#50FA7B",
		Warning: "#F1FA8C",
		Error:   "#FF5555",
		Info:    "#8BE9FD",

		Border:    "#44475A",
		Highlight: "#FF79C6",
		Selection: "#3A3C4E",
		StatusBar: "#21222C",
		Sidebar:   "#21222C",

		Text:      "#F8F8F2",
		TextDim:   "#D8D8E8",
		TextMuted: "#6C6F88",
	}
}

func Kanagawa() Theme {
	return Theme{
		Name: "Kanagawa",

		Background: "#1F1F28",
		Foreground: "#DCD7BA",
		Surface:    "#1A1A22",

		Primary:   "#E46876",
		Secondary: "#938AA9",
		Accent:    "#DCA561",

		Success: "#98BB6C",
		Warning: "#DCA561",
		Error:   "#C34043",
		Info:    "#7E9CD8",

		Border:    "#363646",
		Highlight: "#E46876",
		Selection: "#2D2D3D",
		StatusBar: "#1A1A22",
		Sidebar:   "#1A1A22",

		Text:      "#DCD7BA",
		TextDim:   "#C8C4A8",
		TextMuted: "#727169",
	}
}

func RosePine() Theme {
	return Theme{
		Name: "Rose Pine",

		Background: "#191724",
		Foreground: "#E0DEF4",
		Surface:    "#1F1D2E",

		Primary:   "#EB6F92",
		Secondary: "#C4A7E7",
		Accent:    "#F6C177",

		Success: "#9CCFD8",
		Warning: "#F6C177",
		Error:   "#EB6F92",
		Info:    "#3E8FB0",

		Border:    "#26233A",
		Highlight: "#EB6F92",
		Selection: "#2A273F",
		StatusBar: "#1F1D2E",
		Sidebar:   "#1F1D2E",

		Text:      "#E0DEF4",
		TextDim:   "#C4C0DB",
		TextMuted: "#6E6A86",
	}
}
