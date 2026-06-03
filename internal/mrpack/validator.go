package mrpack

import (
	"fmt"
	"path/filepath"
	"strings"
)

func ValidateIndex(index *PackIndex) error {
	if index == nil {
		return fmt.Errorf("nil index")
	}
	if index.FormatVersion != 1 {
		return fmt.Errorf("unsupported format version: %d", index.FormatVersion)
	}
	if index.Game == "" {
		return fmt.Errorf("missing game field")
	}
	if index.Name == "" {
		return fmt.Errorf("missing name")
	}
	if index.VersionID == "" {
		return fmt.Errorf("missing versionId")
	}
	if len(index.Files) == 0 {
		return fmt.Errorf("no files in index")
	}
	if index.Dependencies == nil {
		return fmt.Errorf("missing dependencies")
	}
	if _, ok := index.Dependencies["minecraft"]; !ok {
		return fmt.Errorf("missing minecraft dependency")
	}
	for _, f := range index.Files {
		if err := validateFilePath(f.Path); err != nil {
			return fmt.Errorf("invalid file path %q: %w", f.Path, err)
		}
		if len(f.Downloads) == 0 {
			return fmt.Errorf("file %q has no downloads", f.Path)
		}
	}
	return nil
}

func validateFilePath(p string) error {
	clean := filepath.Clean(p)
	if strings.HasPrefix(clean, "..") {
		return fmt.Errorf("path traversal detected")
	}
	if filepath.IsAbs(clean) {
		return fmt.Errorf("absolute path not allowed")
	}
	if strings.Contains(clean, "..") {
		return fmt.Errorf("path traversal detected")
	}
	return nil
}

func ValidateZipPath(p string) error {
	clean := filepath.ToSlash(filepath.Clean(p))
	if strings.HasPrefix(clean, "..") {
		return fmt.Errorf("path traversal: %s", p)
	}
	if filepath.IsAbs(clean) || strings.HasPrefix(clean, "/") || strings.HasPrefix(clean, `\`) {
		return fmt.Errorf("absolute path: %s", p)
	}
	if strings.Contains(clean, "..") {
		return fmt.Errorf("path traversal: %s", p)
	}
	if strings.Contains(clean, ":") {
		return fmt.Errorf("unsafe path: %s", p)
	}
	return nil
}
