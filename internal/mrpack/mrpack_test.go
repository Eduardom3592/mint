package mrpack

import (
	"archive/zip"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestMrpack(t *testing.T, dir string, index *PackIndex, files map[string]string) string {
	t.Helper()
	path := filepath.Join(dir, "test.mrpack")
	z, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer z.Close()

	w := zip.NewWriter(z)
	defer w.Close()

	if index != nil {
		rc, err := w.Create("modrinth.index.json")
		if err != nil {
			t.Fatal(err)
		}
		if err := json.NewEncoder(rc).Encode(index); err != nil {
			t.Fatal(err)
		}
	}

	for name, content := range files {
		rc, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		_, err = rc.Write([]byte(content))
		if err != nil {
			t.Fatal(err)
		}
	}

	return path
}

func defaultIndex() *PackIndex {
	return &PackIndex{
		FormatVersion: 1,
		Game:          "minecraft",
		VersionID:     "1.0.0",
		Name:          "Test Pack",
		Summary:       "A test modpack",
		Files: []PackFile{
			{
				Path:      "mods/test.jar",
				Hashes:    map[string]string{"sha1": "abc123"},
				Downloads: []string{"https://example.com/test.jar"},
				FileSize:  1024,
			},
		},
		Dependencies: map[string]string{
			"minecraft":     "1.20.1",
			"fabric-loader": "0.15.0",
		},
	}
}

func TestParse(t *testing.T) {
	dir := t.TempDir()
	path := writeTestMrpack(t, dir, defaultIndex(), map[string]string{
		"overrides/config/foo.txt": "hello",
	})

	index, meta, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if index == nil {
		t.Fatal("Parse() returned nil index")
	}
	if index.Name != "Test Pack" {
		t.Errorf("expected name 'Test Pack', got %q", index.Name)
	}
	if meta == nil {
		t.Fatal("Parse() returned nil metadata")
	}
	if meta.FileCount != 1 {
		t.Errorf("expected 1 file, got %d", meta.FileCount)
	}
}

func TestParseInvalidZip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.mrpack")
	if err := os.WriteFile(path, []byte("not a zip"), 0644); err != nil {
		t.Fatal(err)
	}
	_, _, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid zip")
	}
}

func TestParseNoIndex(t *testing.T) {
	dir := t.TempDir()
	path := writeTestMrpack(t, dir, nil, nil)
	_, _, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for missing modrinth.index.json")
	}
}

func TestValidateIndex(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		err := ValidateIndex(defaultIndex())
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("nil index", func(t *testing.T) {
		err := ValidateIndex(nil)
		if err == nil {
			t.Fatal("expected error for nil index")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		idx := defaultIndex()
		idx.Name = ""
		err := ValidateIndex(idx)
		if err == nil {
			t.Fatal("expected error for empty name")
		}
	})

	t.Run("empty version id", func(t *testing.T) {
		idx := defaultIndex()
		idx.VersionID = ""
		err := ValidateIndex(idx)
		if err == nil {
			t.Fatal("expected error for empty version id")
		}
	})

	t.Run("no dependencies", func(t *testing.T) {
		idx := defaultIndex()
		idx.Dependencies = nil
		err := ValidateIndex(idx)
		if err == nil {
			t.Fatal("expected error for nil dependencies")
		}
	})

	t.Run("no minecraft dep", func(t *testing.T) {
		idx := defaultIndex()
		idx.Dependencies = map[string]string{"fabric-loader": "0.15.0"}
		err := ValidateIndex(idx)
		if err == nil {
			t.Fatal("expected error for missing minecraft dep")
		}
	})
}

func TestValidateZipPath(t *testing.T) {
	tests := []struct {
		path    string
		invalid bool
	}{
		{"mods/test.jar", false},
		{"config/foo.txt", false},
		{"/etc/passwd", true},
		{"../escape", true},
		{"foo/../../bar", true},
		{"overrides/config.txt", false},
		{"C:\\windows\\system32", true},
	}
	for _, tt := range tests {
		err := ValidateZipPath(tt.path)
		if tt.invalid && err == nil {
			t.Errorf("expected error for path %q", tt.path)
		}
		if !tt.invalid && err != nil {
			t.Errorf("unexpected error for path %q: %v", tt.path, err)
		}
	}
}

func TestHashFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	content := "hello world"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	sha1Hash, err := HashFileSHA1(path)
	if err != nil {
		t.Fatalf("HashFileSHA1() error = %v", err)
	}
	h := sha1.Sum([]byte(content))
	expected := hex.EncodeToString(h[:])
	if sha1Hash != expected {
		t.Errorf("SHA1 mismatch: got %q, want %q", sha1Hash, expected)
	}

	sha512Hash, err := HashFileSHA512(path)
	if err != nil {
		t.Fatalf("HashFileSHA512() error = %v", err)
	}
	h512 := sha512.Sum512([]byte(content))
	expected512 := hex.EncodeToString(h512[:])
	if sha512Hash != expected512 {
		t.Errorf("SHA512 mismatch: got %q, want %q", sha512Hash, expected512)
	}
}

func TestVerifyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("hello world"), 0644); err != nil {
		t.Fatal(err)
	}

	h := sha512.Sum512([]byte("hello world"))
	goodHash := hex.EncodeToString(h[:])

	t.Run("valid sha512", func(t *testing.T) {
		ok, err := VerifyFile(path, map[string]string{"sha512": goodHash})
		if err != nil {
			t.Fatalf("VerifyFile() error = %v", err)
		}
		if !ok {
			t.Fatal("expected verification to pass")
		}
	})

	t.Run("invalid sha512", func(t *testing.T) {
		ok, err := VerifyFile(path, map[string]string{"sha512": "badhash"})
		if err != nil {
			t.Fatalf("VerifyFile() error = %v", err)
		}
		if ok {
			t.Fatal("expected verification to fail")
		}
	})

	t.Run("empty hashes returns true", func(t *testing.T) {
		ok, err := VerifyFile(path, nil)
		if err != nil {
			t.Fatalf("VerifyFile() error = %v", err)
		}
		if !ok {
			t.Fatal("expected empty hashes to return true")
		}
	})
}

func TestExtract(t *testing.T) {
	dir := t.TempDir()
	outDir := filepath.Join(dir, "out")
	path := writeTestMrpack(t, dir, defaultIndex(), map[string]string{
		"overrides/config/foo.txt": "hello",
		"mods/test.jar":            "fake jar content",
	})

	if err := Extract(path, outDir, nil); err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(outDir, "overrides/config/foo.txt")); err != nil {
		t.Errorf("extracted file missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "mods/test.jar")); err != nil {
		t.Errorf("extracted file missing: %v", err)
	}
}

func TestExtractPathTraversal(t *testing.T) {
	dir := t.TempDir()
	outDir := filepath.Join(dir, "out")

	path := filepath.Join(dir, "bad.mrpack")
	z, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	w := zip.NewWriter(z)
	rc, err := w.Create("../outside.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = rc.Write([]byte("bad"))
	w.Close()
	z.Close()

	err = Extract(path, outDir, nil)
	if err == nil {
		t.Fatal("expected error for path traversal")
	}
	if !strings.Contains(err.Error(), "unsafe") && !strings.Contains(err.Error(), "invalid") {
		t.Logf("got error: %v", err)
	}
}

func TestExtractOverrides(t *testing.T) {
	dir := t.TempDir()
	outDir := filepath.Join(dir, "out")
	overrideDir := "overrides"

	path := writeTestMrpack(t, dir, defaultIndex(), map[string]string{
		"overrides/config/foo.txt": "override content",
	})

	if err := ExtractOverrides(path, outDir, overrideDir); err != nil {
		t.Fatalf("ExtractOverrides() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "config/foo.txt"))
	if err != nil {
		t.Fatalf("read extracted override: %v", err)
	}
	if string(data) != "override content" {
		t.Errorf("expected 'override content', got %q", string(data))
	}
}

func TestInstall(t *testing.T) {
	dir := t.TempDir()
	installDir := filepath.Join(dir, "install")
	zipPath := writeTestMrpack(t, dir, defaultIndex(), map[string]string{
		"overrides/config/foo.txt": "hello",
	})

	record, err := Install(zipPath, installDir)
	if err != nil {
		t.Fatalf("Install() error = %v", err)
	}
	if record == nil {
		t.Fatal("Install() returned nil record")
	}
	if record.Name != "Test Pack" {
		t.Errorf("expected name 'Test Pack', got %q", record.Name)
	}
	if !record.Verified {
		t.Error("expected verified to be true")
	}

	if _, err := os.Stat(filepath.Join(installDir, "Test Pack", "overrides/config/foo.txt")); err != nil {
		t.Errorf("installed file missing: %v", err)
	}
}

func TestInstallInvalid(t *testing.T) {
	dir := t.TempDir()
	_, err := Install(filepath.Join(dir, "nonexistent.mrpack"), dir)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestApplyOverrides(t *testing.T) {
	dir := t.TempDir()
	dstDir := filepath.Join(dir, "dest")

	path := writeTestMrpack(t, dir, defaultIndex(), map[string]string{
		"overrides/config/foo.txt": "config override",
		"overrides/mods/bar.jar":   "mod override",
	})

	if err := ApplyOverrides(path, dstDir); err != nil {
		t.Fatalf("ApplyOverrides() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(dstDir, "config/foo.txt")); err != nil {
		t.Errorf("override file missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dstDir, "mods/bar.jar")); err != nil {
		t.Errorf("override file missing: %v", err)
	}
}

func TestApplyAllOverrides(t *testing.T) {
	dir := t.TempDir()
	dstDir := filepath.Join(dir, "dest")

	path := writeTestMrpack(t, dir, defaultIndex(), map[string]string{
		"overrides/config/base.txt":          "base",
		"client-overrides/config/client.txt": "client",
		"server-overrides/config/server.txt": "server",
	})

	if err := ApplyAllOverrides(path, dstDir); err != nil {
		t.Fatalf("ApplyAllOverrides() error = %v", err)
	}

	tests := []struct {
		rel   string
		exist bool
	}{
		{"config/base.txt", true},
		{"config/client.txt", true},
		{"config/server.txt", true},
	}
	for _, tt := range tests {
		_, err := os.Stat(filepath.Join(dstDir, tt.rel))
		if tt.exist && err != nil {
			t.Errorf("expected %s to exist: %v", tt.rel, err)
		}
		if !tt.exist && err == nil {
			t.Errorf("expected %s to not exist", tt.rel)
		}
	}
}

func TestOverrideFileCount(t *testing.T) {
	dir := t.TempDir()
	path := writeTestMrpack(t, dir, defaultIndex(), map[string]string{
		"overrides/config/foo.txt": "hello",
		"overrides/mods/bar.jar":   "world",
	})

	count, err := OverrideFileCount(path, "overrides")
	if err != nil {
		t.Fatalf("OverrideFileCount() error = %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 override files, got %d", count)
	}
}
