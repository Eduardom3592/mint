package models

import "testing"

func TestProjectTypeConstants(t *testing.T) {
	tests := []struct {
		pt   ProjectType
		want string
	}{
		{ProjectTypeMod, "mod"},
		{ProjectTypeModpack, "modpack"},
		{ProjectTypeResourcePack, "resourcepack"},
		{ProjectTypeShader, "shader"},
		{ProjectTypeDatapack, "datapack"},
		{ProjectTypePlugin, "plugin"},
	}
	for _, tt := range tests {
		if string(tt.pt) != tt.want {
			t.Errorf("expected %s, got %s", tt.want, string(tt.pt))
		}
	}
}

func TestVersionTypeConstants(t *testing.T) {
	tests := []struct {
		vt   VersionType
		want string
	}{
		{VersionTypeRelease, "release"},
		{VersionTypeBeta, "beta"},
		{VersionTypeAlpha, "alpha"},
	}
	for _, tt := range tests {
		if string(tt.vt) != tt.want {
			t.Errorf("expected %s, got %s", tt.want, string(tt.vt))
		}
	}
}
