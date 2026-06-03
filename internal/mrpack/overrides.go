package mrpack

import (
	"fmt"
	"os"
	"path/filepath"
)

type OverrideLayer int

const (
	OverrideBase   OverrideLayer = 0
	OverrideClient OverrideLayer = 1
	OverrideServer OverrideLayer = 2
)

func ApplyOverrides(mrpackPath, installDir string) error {
	if err := ExtractOverrides(mrpackPath, installDir, "overrides"); err != nil {
		return fmt.Errorf("base overrides: %w", err)
	}
	return nil
}

func ApplyClientOverrides(mrpackPath, installDir string) error {
	return ExtractOverrides(mrpackPath, installDir, "client-overrides")
}

func ApplyServerOverrides(mrpackPath, installDir string) error {
	return ExtractOverrides(mrpackPath, installDir, "server-overrides")
}

func ApplyAllOverrides(mrpackPath, installDir string) error {
	layers := []struct {
		name string
		fn   func(string, string) error
	}{
		{"overrides", ApplyOverrides},
		{"client-overrides", ApplyClientOverrides},
		{"server-overrides", ApplyServerOverrides},
	}

	for _, layer := range layers {
		if err := layer.fn(mrpackPath, installDir); err != nil {
			return fmt.Errorf("%s: %w", layer.name, err)
		}
	}

	return nil
}

func OverrideFileCount(mrpackPath, overrideDir string) (int, error) {
	count := 0
	tmp := filepath.Join(os.TempDir(), "mint-override-count")
	os.RemoveAll(tmp)
	defer os.RemoveAll(tmp)
	if err := ExtractOverrides(mrpackPath, tmp, overrideDir); err != nil {
		return 0, err
	}
	_ = filepath.Walk(tmp, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, nil
}
