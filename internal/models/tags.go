package models

type Category struct {
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	ProjectType string `json:"project_type"`
	Header      string `json:"header"`
}

type Loader struct {
	Icon      string `json:"icon"`
	Name      string `json:"name"`
	Supported bool   `json:"supported"`
}

type GameVersion struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	VersionType string `json:"version_type"`
	Date        string `json:"date"`
	Major       bool   `json:"major"`
}

type LicenseTag struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	URL  string `json:"url"`
}

type ReportType struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type TagBundle struct {
	Categories   []Category    `json:"categories"`
	Loaders      []Loader      `json:"loaders"`
	GameVersions []GameVersion `json:"game_versions"`
	Licenses     []LicenseTag  `json:"licenses"`
	ReportTypes  []ReportType  `json:"report_types"`
}
