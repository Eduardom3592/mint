package mrpack

import (
	"fmt"
	"os"
	"path/filepath"
)

type InstallResult struct {
	Name         string
	VersionID    string
	GameVersion  string
	Loader       string
	Path         string
	FileCount    int
	OverridesDir string
	Verified     bool
}

func Install(mrpackPath, installDir string) (*InstallResult, error) {
	index, meta, err := Parse(mrpackPath)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	if err := ValidateIndex(index); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	packDir := filepath.Join(installDir, sanitizeDirName(index.Name))
	if err := os.MkdirAll(packDir, 0755); err != nil {
		return nil, fmt.Errorf("create install dir: %w", err)
	}

	if err := Extract(mrpackPath, packDir, []string{"modrinth.index.json"}); err != nil {
		return nil, fmt.Errorf("extract: %w", err)
	}

	if err := ApplyAllOverrides(mrpackPath, packDir); err != nil {
		return nil, fmt.Errorf("overrides: %w", err)
	}

	result := &InstallResult{
		Name:         index.Name,
		VersionID:    index.VersionID,
		GameVersion:  meta.GameVersion,
		Loader:       meta.Loader,
		Path:         packDir,
		FileCount:    len(index.Files),
		OverridesDir: filepath.Join(packDir, "overrides"),
		Verified:     true,
	}

	return result, nil
}

func sanitizeDirName(name string) string {
	result := make([]byte, 0, len(name))
	for i := 0; i < len(name); i++ {
		c := name[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' || c == ' ' {
			result = append(result, c)
		} else {
			result = append(result, '_')
		}
	}
	return string(result)
}
