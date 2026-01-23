package root

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	agentLayerDir = ".agent-layer"
	gitDir        = ".git"
)

// FindAgentLayerRoot walks upward from start until it finds a directory containing .agent-layer/.
// It returns the root path, whether it was found, and any error encountered.
func FindAgentLayerRoot(start string) (string, bool, error) {
	if start == "" {
		return "", false, fmt.Errorf("start path is required")
	}
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", false, fmt.Errorf("resolve path %s: %w", start, err)
	}

	dir := abs
	for {
		candidate := filepath.Join(dir, agentLayerDir)
		info, err := os.Stat(candidate)
		if err == nil {
			if !info.IsDir() {
				return "", false, fmt.Errorf("%s exists but is not a directory", candidate)
			}
			return dir, true, nil
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", false, fmt.Errorf("check %s: %w", candidate, err)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false, nil
		}
		dir = parent
	}
}

// FindRepoRoot returns the repo root for initialization.
// It prefers an existing .agent-layer directory, then a .git directory or file, and falls back to start.
func FindRepoRoot(start string) (string, error) {
	if start == "" {
		return "", fmt.Errorf("start path is required")
	}
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", fmt.Errorf("resolve path %s: %w", start, err)
	}

	if root, found, err := FindAgentLayerRoot(abs); err != nil {
		return "", err
	} else if found {
		return root, nil
	}

	dir := abs
	for {
		candidate := filepath.Join(dir, gitDir)
		info, err := os.Stat(candidate)
		if err == nil {
			if info.IsDir() || info.Mode().IsRegular() {
				return dir, nil
			}
			return "", fmt.Errorf("%s exists but is not a directory or file", candidate)
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("check %s: %w", candidate, err)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return abs, nil
		}
		dir = parent
	}
}
