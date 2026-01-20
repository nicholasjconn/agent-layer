package wizard

import (
	"fmt"
	"os"

	"github.com/nicholasjconn/agent-layer/internal/envfile"
	"github.com/nicholasjconn/agent-layer/internal/fsutil"
	"github.com/nicholasjconn/agent-layer/internal/sync"
)

// applyChanges writes config/env updates and runs sync.
// root/configPath/envPath identify files; c holds wizard selections; returns an error on failure.
func applyChanges(root, configPath, envPath string, c *Choices) error {
	// Config
	rawConfig, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	// Backup
	if err := os.WriteFile(configPath+".bak", rawConfig, 0644); err != nil {
		return fmt.Errorf("failed to backup config: %w", err)
	}
	// Patch
	newConfig, err := PatchConfig(string(rawConfig), c)
	if err != nil {
		return fmt.Errorf("failed to patch config: %w", err)
	}
	if err := fsutil.WriteFileAtomic(configPath, []byte(newConfig), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Env
	// Backup if exists
	rawEnv, err := os.ReadFile(envPath)
	if err == nil {
		if err := os.WriteFile(envPath+".bak", rawEnv, 0600); err != nil {
			return fmt.Errorf("failed to backup .env: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	// Patch
	newEnv := envfile.Patch(string(rawEnv), c.Secrets)
	if err := fsutil.WriteFileAtomic(envPath, []byte(newEnv), 0600); err != nil {
		return fmt.Errorf("failed to write .env: %w", err)
	}

	// Sync
	fmt.Println("Running sync...")
	return sync.Run(root)
}
