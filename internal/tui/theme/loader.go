package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const themesDir = ".config/mint/themes"

func themesDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	return filepath.Join(home, themesDir), nil
}

func SaveTheme(t Theme) error {
	dir, err := themesDirPath()
	if err != nil {
		return err
	}
	if err = os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create themes dir: %w", err)
	}

	filename := filepath.Join(dir, sanitizeName(t.Name)+".json")
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal theme: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("write theme: %w", err)
	}

	return nil
}

func LoadCustomThemes() ([]Theme, error) {
	dir, err := themesDirPath()
	if err != nil {
		return nil, err
	}

	if _, statErr := os.Stat(dir); os.IsNotExist(statErr) {
		return nil, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read themes dir: %w", err)
	}

	var themes []Theme
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var t Theme
		if err := json.Unmarshal(data, &t); err != nil {
			continue
		}

		themes = append(themes, t)
	}

	return themes, nil
}

func sanitizeName(name string) string {
	result := make([]byte, 0, len(name))
	for _, c := range []byte(name) {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			result = append(result, c)
		} else if c == ' ' {
			result = append(result, '_')
		}
	}
	return string(result)
}
