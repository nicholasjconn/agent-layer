package dispatch

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicholasjconn/agent-layer/internal/version"
)

// readPinnedVersion reads and normalizes the pinned version from .agent-layer/al.version.
func readPinnedVersion(rootDir string) (string, bool, error) {
	path := filepath.Join(rootDir, ".agent-layer", "al.version")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("read %s: %w", path, err)
	}
	raw := strings.TrimSpace(string(data))
	if raw == "" {
		return "", false, fmt.Errorf("pin file %s is empty", path)
	}
	normalized, err := version.Normalize(raw)
	if err != nil {
		return "", false, fmt.Errorf("invalid pinned version in %s: %w", path, err)
	}
	return normalized, true, nil
}
