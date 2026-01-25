package sync

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/conn-castle/agent-layer/internal/messages"
)

// resolvePromptServerCommand returns the command and args used to run the internal MCP prompt server.
// It prefers the globally installed "al mcp-prompts" and falls back to "go run <root>/cmd/al mcp-prompts" for dev usage.
// It returns an error when it cannot resolve a runnable command.
func resolvePromptServerCommand(sys System, root string) (string, []string, error) {
	if _, err := sys.LookPath("al"); err == nil {
		return "al", []string{"mcp-prompts"}, nil
	}

	if root == "" {
		return "", nil, fmt.Errorf(messages.SyncMissingPromptServerNoRoot)
	}

	sourcePath := filepath.Join(root, "cmd", "al")
	info, err := sys.Stat(sourcePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil, fmt.Errorf(messages.SyncMissingPromptServerSourceFmt, sourcePath)
		}
		return "", nil, fmt.Errorf(messages.SyncCheckPathFmt, sourcePath, err)
	}
	if !info.IsDir() {
		return "", nil, fmt.Errorf(messages.SyncPromptServerNotDirFmt, sourcePath)
	}

	if _, err := sys.LookPath("go"); err != nil {
		return "", nil, fmt.Errorf(messages.SyncMissingGoForPromptServerFmt, err)
	}

	return "go", []string{"run", sourcePath, "mcp-prompts"}, nil
}
