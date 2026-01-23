package sync

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var lookPath = exec.LookPath

// resolvePromptServerCommand returns the command and args used to run the internal MCP prompt server.
// It prefers the globally installed "al mcp-prompts" and falls back to "go run <root>/cmd/al mcp-prompts" for dev usage.
// It returns an error when it cannot resolve a runnable command.
func resolvePromptServerCommand(root string) (string, []string, error) {
	if _, err := lookPath("al"); err == nil {
		return "al", []string{"mcp-prompts"}, nil
	}

	if root == "" {
		return "", nil, fmt.Errorf("al not found on PATH and no repo root available for go run")
	}

	sourcePath := filepath.Join(root, "cmd", "al")
	info, err := os.Stat(sourcePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil, fmt.Errorf("missing prompt server source at %s", sourcePath)
		}
		return "", nil, fmt.Errorf("check %s: %w", sourcePath, err)
	}
	if !info.IsDir() {
		return "", nil, fmt.Errorf("prompt server source path %s is not a directory", sourcePath)
	}

	if _, err := lookPath("go"); err != nil {
		return "", nil, fmt.Errorf("missing go on PATH for prompt server: %w", err)
	}

	return "go", []string{"run", sourcePath, "mcp-prompts"}, nil
}
