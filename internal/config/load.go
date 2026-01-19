package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// LoadProjectConfig reads and validates the full Agent Layer config from disk.
func LoadProjectConfig(root string) (*ProjectConfig, error) {
	paths := DefaultPaths(root)

	cfg, err := LoadConfig(paths.ConfigPath)
	if err != nil {
		return nil, err
	}

	env, err := LoadEnv(paths.EnvPath)
	if err != nil {
		return nil, err
	}

	instructions, err := LoadInstructions(paths.InstructionsDir)
	if err != nil {
		return nil, err
	}

	slashCommands, err := LoadSlashCommands(paths.SlashCommandsDir)
	if err != nil {
		return nil, err
	}

	commandsAllow, err := LoadCommandsAllow(paths.CommandsAllow)
	if err != nil {
		return nil, err
	}

	return &ProjectConfig{
		Config:        *cfg,
		Env:           env,
		Instructions:  instructions,
		SlashCommands: slashCommands,
		CommandsAllow: commandsAllow,
		Root:          root,
	}, nil
}

// LoadConfig reads .agent-layer/config.toml and validates it.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("missing config file %s: %w", path, err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config file %s: %w", path, err)
	}

	if err := cfg.Validate(path); err != nil {
		return nil, err
	}

	return &cfg, nil
}
