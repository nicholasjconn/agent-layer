package sync

import (
	"errors"
	"os"
	"testing"

	"github.com/conn-castle/agent-layer/internal/config"
)

func TestBuildMCPConfig_PromptCommandError(t *testing.T) {
	t.Parallel()
	sys := &MockSystem{
		LookPathFunc: func(string) (string, error) {
			return "", errors.New("missing")
		},
		StatFunc: func(string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		},
	}

	// Create a root without al binary or source, so resolvePromptServerCommand fails
	root := t.TempDir()
	project := &config.ProjectConfig{Root: root}

	_, err := buildMCPConfig(sys, project)
	if err == nil {
		t.Fatalf("expected error from resolvePromptServerCommand")
	}
}
