package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := `
# comment
export API_KEY="abc123"
NAME=plain
QUOTED='value with spaces'
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write env: %v", err)
	}

	env, err := LoadEnv(path)
	if err != nil {
		t.Fatalf("LoadEnv error: %v", err)
	}

	if env["API_KEY"] != "abc123" {
		t.Fatalf("expected API_KEY to be abc123, got %q", env["API_KEY"])
	}
	if env["NAME"] != "plain" {
		t.Fatalf("expected NAME to be plain, got %q", env["NAME"])
	}
	if env["QUOTED"] != "value with spaces" {
		t.Fatalf("unexpected QUOTED value: %q", env["QUOTED"])
	}
}

func TestLoadEnvMissing(t *testing.T) {
	_, err := LoadEnv(filepath.Join(t.TempDir(), ".env"))
	if err == nil {
		t.Fatalf("expected missing env error")
	}
	if !strings.Contains(err.Error(), "missing env file") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadEnvInvalidLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("not-an-env-line"), 0o644); err != nil {
		t.Fatalf("write env: %v", err)
	}

	_, err := LoadEnv(path)
	if err == nil {
		t.Fatalf("expected invalid env error")
	}
	if !strings.Contains(err.Error(), "expected KEY=VALUE") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadEnvMissingKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("=value"), 0o644); err != nil {
		t.Fatalf("write env: %v", err)
	}

	_, err := LoadEnv(path)
	if err == nil {
		t.Fatalf("expected invalid env error")
	}
	if !strings.Contains(err.Error(), "expected KEY=VALUE") {
		t.Fatalf("unexpected error: %v", err)
	}
}
