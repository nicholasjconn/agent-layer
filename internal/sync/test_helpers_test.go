package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func writePromptServerBinary(t *testing.T, root string) {
	t.Helper()
	path := filepath.Join(root, "al")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write prompt server binary: %v", err)
	}
	t.Setenv("PATH", root+string(os.PathListSeparator)+os.Getenv("PATH"))
}
