package mrpack

type PackIndex struct {
	FormatVersion int               `json:"formatVersion"`
	Game          string            `json:"game"`
	VersionID     string            `json:"versionId"`
	Name          string            `json:"name"`
	Summary       string            `json:"summary,omitempty"`
	Files         []PackFile        `json:"files"`
	Dependencies  map[string]string `json:"dependencies"`
}

type PackFile struct {
	Path      string            `json:"path"`
	Hashes    map[string]string `json:"hashes"`
	Env       map[string]string `json:"env,omitempty"`
	Downloads []string          `json:"downloads"`
	FileSize  int64             `json:"fileSize,omitempty"`
}

type EnvSide string

const (
	EnvRequired    EnvSide = "required"
	EnvOptional    EnvSide = "optional"
	EnvUnsupported EnvSide = "unsupported"
)

type DependencyType string

const (
	DepMinecraft    DependencyType = "minecraft"
	DepFabricLoader DependencyType = "fabric-loader"
	DepQuiltLoader  DependencyType = "quilt-loader"
	DepForge        DependencyType = "forge"
	DepNeoForge     DependencyType = "neoforge"
)

type PackMetadata struct {
	Path                string
	Name                string
	VersionID           string
	Summary             string
	GameVersion         string
	Loader              string
	FileCount           int
	OverrideCount       int
	ClientOverrideCount int
	ServerOverrideCount int
	TotalSize           int64
}

type InstallRecord struct {
	ID          int64
	Name        string
	VersionID   string
	GameVersion string
	Loader      string
	Path        string
	InstalledAt string
	FileCount   int
	Verified    bool
}
