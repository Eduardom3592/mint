package models

import "time"

type VersionType string

const (
	VersionTypeRelease VersionType = "release"
	VersionTypeBeta    VersionType = "beta"
	VersionTypeAlpha   VersionType = "alpha"
)

type VersionFile struct {
	Hashes   map[string]string `json:"hashes"`
	URL      string            `json:"url"`
	Filename string            `json:"filename"`
	Primary  bool              `json:"primary"`
	Size     int64             `json:"size"`
	FileType string            `json:"file_type"`
}

type Dependency struct {
	VersionID      *string `json:"version_id"`
	ProjectID      *string `json:"project_id"`
	FileName       *string `json:"file_name"`
	DependencyType string  `json:"dependency_type"`
}

type Version struct {
	ID              string        `json:"id"`
	ProjectID       string        `json:"project_id"`
	AuthorID        string        `json:"author_id"`
	Name            string        `json:"name"`
	VersionNumber   string        `json:"version_number"`
	Changelog       string        `json:"changelog"`
	ChangelogURL    *string       `json:"changelog_url"`
	DatePublished   time.Time     `json:"date_published"`
	Downloads       int           `json:"downloads"`
	VersionType     VersionType   `json:"version_type"`
	Status          string        `json:"status"`
	RequestedStatus *string       `json:"requested_status"`
	Files           []VersionFile `json:"files"`
	Dependencies    []Dependency  `json:"dependencies"`
	GameVersions    []string      `json:"game_versions"`
	Loaders         []string      `json:"loaders"`
	Featured        bool          `json:"featured"`
}
