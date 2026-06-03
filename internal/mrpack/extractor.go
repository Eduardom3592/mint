package mrpack

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Extract(path, destDir string, skipPaths []string) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("open mrpack: %w", err)
	}
	defer r.Close()

	skip := make(map[string]bool)
	for _, s := range skipPaths {
		skip[filepath.ToSlash(s)] = true
	}

	for _, f := range r.File {
		name := filepath.ToSlash(f.Name)
		if skip[name] {
			continue
		}

		if err := ValidateZipPath(name); err != nil {
			return err
		}

		target := filepath.Join(destDir, name)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return fmt.Errorf("create dir %s: %w", filepath.Dir(target), err)
		}

		if f.FileInfo().IsDir() {
			continue
		}

		if err := extractFile(f, target); err != nil {
			return fmt.Errorf("extract %s: %w", name, err)
		}
	}

	return nil
}

func ExtractOverrides(mrpackPath, destDir, overrideDir string) error {
	srcPath := filepath.ToSlash(overrideDir) + "/"
	r, err := zip.OpenReader(mrpackPath)
	if err != nil {
		return fmt.Errorf("open mrpack: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		name := filepath.ToSlash(f.Name)
		if !strings.HasPrefix(name, srcPath) {
			continue
		}

		rel := strings.TrimPrefix(name, srcPath)
		if rel == "" {
			continue
		}

		if err := ValidateZipPath(name); err != nil {
			return err
		}

		target := filepath.Join(destDir, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return fmt.Errorf("create dir: %w", err)
		}

		if f.FileInfo().IsDir() {
			continue
		}

		if err := extractFile(f, target); err != nil {
			return fmt.Errorf("extract override %s: %w", name, err)
		}
	}

	return nil
}

func extractFile(f *zip.File, target string) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("open zip entry: %w", err)
	}
	defer rc.Close()

	out, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return err
}
