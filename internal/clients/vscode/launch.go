package vscode

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nicholasjconn/agent-layer/internal/clients"
	"github.com/nicholasjconn/agent-layer/internal/config"
	"github.com/nicholasjconn/agent-layer/internal/run"
)

// Launch starts VS Code with CODEX_HOME set for the repo.
func Launch(cfg *config.ProjectConfig, runInfo *run.Info, env []string) error {
	codexHome := filepath.Join(cfg.Root, ".codex")
	env = clients.SetEnv(env, "CODEX_HOME", codexHome)

	cmd := exec.Command("code", ".")
	cmd.Dir = cfg.Root
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("vscode exited with error: %w", err)
	}

	return nil
}
