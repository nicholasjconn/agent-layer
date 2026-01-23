package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolvePromptServerCommand_MissingSource(t *testing.T) {
	root := t.TempDir()
	originalLookPath := lookPath
	t.Cleanup(func() { lookPath = originalLookPath })
	lookPath = func(file string) (string, error) {
		return "", os.ErrNotExist
	}

	// Neither al on PATH nor ./cmd/al exists
	_, _, err := resolvePromptServerCommand(root)
	if err == nil {
		t.Fatalf("expected error for missing source")
	}
	if !strings.Contains(err.Error(), "missing prompt server source") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolvePromptServerCommand_SourceNotDir(t *testing.T) {
	root := t.TempDir()
	originalLookPath := lookPath
	t.Cleanup(func() { lookPath = originalLookPath })
	lookPath = func(file string) (string, error) {
		return "", os.ErrNotExist
	}

	// ./cmd/al is a file
	cmdDir := filepath.Join(root, "cmd")
	if err := os.MkdirAll(cmdDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cmdDir, "al"), []byte("file"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, _, err := resolvePromptServerCommand(root)
	if err == nil {
		t.Fatalf("expected error for file source")
	}
	if !strings.Contains(err.Error(), "is not a directory") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolvePromptServerCommand_SourceStatError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping stat error test as root")
	}
	root := t.TempDir()

	originalLookPath := lookPath
	t.Cleanup(func() { lookPath = originalLookPath })
	lookPath = func(file string) (string, error) {
		return "", os.ErrNotExist
	}

	// We want os.Stat(root/cmd/al) to fail with permission denied.
	// root must be searchable; cmd must be unsearchable.

	cmdDir := filepath.Join(root, "cmd")
	if err := os.Mkdir(cmdDir, 0o000); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	defer func() { _ = os.Chmod(cmdDir, 0o755) }()

	_, _, err := resolvePromptServerCommand(root)
	if err == nil {
		t.Fatalf("expected error for source stat failure")
	}
	// It should fail at os.Stat(sourcePath)
	if !strings.Contains(err.Error(), "check") {
		t.Fatalf("unexpected error: %v", err)
	}
}
