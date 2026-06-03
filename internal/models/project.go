package models

import "time"

type ProjectType string

const (
	ProjectTypeMod          ProjectType = "mod"
	ProjectTypeModpack      ProjectType = "modpack"
	ProjectTypeResourcePack ProjectType = "resourcepack"
	ProjectTypeShader       ProjectType = "shader"
	ProjectTypeDatapack     ProjectType = "datapack"
	ProjectTypePlugin       ProjectType = "plugin"
)

type License struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type DonationURL struct {
	ID       string `json:"id"`
	Platform string `json:"platform"`
	URL      string `json:"url"`
}

type GalleryEntry struct {
	URL         string `json:"url"`
	Featured    bool   `json:"featured"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"ordering"`
}

type Project struct {
	ID                   string         `json:"id"`
	ProjectType          string         `json:"project_type"`
	Slug                 string         `json:"slug"`
	Title                string         `json:"title"`
	Description          string         `json:"description"`
	Body                 string         `json:"body"`
	BodyURL              string         `json:"body_url"`
	Published            time.Time      `json:"published"`
	Updated              time.Time      `json:"updated"`
	Approved             time.Time      `json:"approved"`
	Queued               *time.Time     `json:"queued"`
	Status               string         `json:"status"`
	RequestedStatus      *string        `json:"requested_status"`
	ModerationMessage    *string        `json:"moderation_message"`
	License              License        `json:"license"`
	ClientSide           string         `json:"client_side"`
	ServerSide           string         `json:"server_side"`
	Downloads            int            `json:"downloads"`
	Followers            int            `json:"followers"`
	Categories           []string       `json:"categories"`
	AdditionalCategories []string       `json:"additional_categories"`
	Loaders              []string       `json:"loaders"`
	GameVersions         []string       `json:"game_versions"`
	Versions             []string       `json:"versions"`
	IconURL              string         `json:"icon_url"`
	IssuesURL            *string        `json:"issues_url"`
	SourceURL            *string        `json:"source_url"`
	WikiURL              *string        `json:"wiki_url"`
	DiscordURL           *string        `json:"discord_url"`
	DonationURLs         []DonationURL  `json:"donation_urls"`
	Gallery              []GalleryEntry `json:"gallery"`
	Color                *int           `json:"color"`
	ThreadID             *string        `json:"thread_id"`
	MonetizationStatus   string         `json:"monetization_status"`
}

type ProjectResult struct {
	Project Project `json:"project"`
	Score   float64 `json:"score"`
}

type SearchResponse struct {
	Hits      []SearchHit `json:"hits"`
	Offset    int         `json:"offset"`
	Limit     int         `json:"limit"`
	TotalHits int         `json:"total_hits"`
}

type SearchHit struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	ProjectID     string   `json:"project_id"`
	ProjectType   string   `json:"project_type"`
	Slug          string   `json:"slug"`
	Author        string   `json:"author"`
	IconURL       string   `json:"icon_url"`
	Downloads     int      `json:"downloads"`
	Follows       int      `json:"follows"`
	Color         *int     `json:"color"`
	Categories    []string `json:"categories"`
	Loaders       []string `json:"loaders"`
	GameVersions  []string `json:"game_versions"`
	Versions      []string `json:"versions"`
	DateCreated   string   `json:"date_created"`
	DateModified  string   `json:"date_modified"`
	LatestVersion string   `json:"latest_version"`
	License       string   `json:"license"`
	ClientSide    string   `json:"client_side"`
	ServerSide    string   `json:"server_side"`
	OpenSource    *bool    `json:"open_source"`
}
