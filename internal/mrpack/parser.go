package mrpack

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
)

func Parse(path string) (*PackIndex, *PackMetadata, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open mrpack: %w", err)
	}
	defer r.Close()

	var index *PackIndex
	for _, f := range r.File {
		if f.Name == "modrinth.index.json" {
			rc, err := f.Open()
			if err != nil {
				return nil, nil, fmt.Errorf("open index: %w", err)
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return nil, nil, fmt.Errorf("read index: %w", err)
			}

			index = &PackIndex{}
			if err := json.Unmarshal(data, index); err != nil {
				return nil, nil, fmt.Errorf("parse index: %w", err)
			}
			break
		}
	}

	if index == nil {
		return nil, nil, fmt.Errorf("modrinth.index.json not found in mrpack")
	}
	if index.FormatVersion != 1 {
		return nil, nil, fmt.Errorf("unsupported format version: %d", index.FormatVersion)
	}

	meta := buildMetadata(index, r.File)
	return index, meta, nil
}

func ParseIndexFromReader(rc io.ReadCloser) (*PackIndex, error) {
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("read index: %w", err)
	}
	index := &PackIndex{}
	if err := json.Unmarshal(data, index); err != nil {
		return nil, fmt.Errorf("parse index: %w", err)
	}
	return index, nil
}

func buildMetadata(index *PackIndex, files []*zip.File) *PackMetadata {
	meta := &PackMetadata{
		Name:      index.Name,
		VersionID: index.VersionID,
		Summary:   index.Summary,
		FileCount: len(index.Files),
	}

	if gv, ok := index.Dependencies["minecraft"]; ok {
		meta.GameVersion = gv
	}
	for _, dep := range []string{"fabric-loader", "quilt-loader", "forge", "neoforge"} {
		if _, ok := index.Dependencies[dep]; ok {
			meta.Loader = dep
			break
		}
	}

	for _, f := range index.Files {
		meta.TotalSize += f.FileSize
	}

	for _, zf := range files {
		name := filepath.ToSlash(zf.Name)
		switch {
		case hasPrefix(name, "overrides/"):
			meta.OverrideCount++
		case hasPrefix(name, "client-overrides/"):
			meta.ClientOverrideCount++
		case hasPrefix(name, "server-overrides/"):
			meta.ServerOverrideCount++
		}
	}

	return meta
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
