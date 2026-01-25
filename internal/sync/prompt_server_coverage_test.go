package sync

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolvePromptServerCommand_MissingSource(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	sys := &MockSystem{
		LookPathFunc: func(string) (string, error) {
			return "", os.ErrNotExist
		},
		StatFunc: func(string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		},
	}

	// Neither al on PATH nor ./cmd/al exists
	_, _, err := resolvePromptServerCommand(sys, root)
	if err == nil {
		t.Fatalf("expected error for missing source")
	}
	if !strings.Contains(err.Error(), "missing prompt server source") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolvePromptServerCommand_SourceNotDir(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	sys := &MockSystem{
		LookPathFunc: func(string) (string, error) {
			return "", os.ErrNotExist
		},
		StatFunc: func(name string) (os.FileInfo, error) {
			if name == filepath.Join(root, "cmd", "al") {
				return &mockFileInfo{isDir: false}, nil
			}
			return nil, os.ErrNotExist
		},
	}

	_, _, err := resolvePromptServerCommand(sys, root)
	if err == nil {
		t.Fatalf("expected error for file source")
	}
	if !strings.Contains(err.Error(), "is not a directory") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolvePromptServerCommand_SourceStatError(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	sys := &MockSystem{
		LookPathFunc: func(string) (string, error) {
			return "", os.ErrNotExist
		},
		StatFunc: func(string) (os.FileInfo, error) {
			return nil, errors.New("stat failed")
		},
	}

	_, _, err := resolvePromptServerCommand(sys, root)
	if err == nil {
		t.Fatalf("expected error for source stat failure")
	}
	// It should fail at os.Stat(sourcePath)
	if !strings.Contains(err.Error(), "check") {
		t.Fatalf("unexpected error: %v", err)
	}
}
