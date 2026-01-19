package config

import "path/filepath"

// Paths holds resolved paths for config files and directories.
type Paths struct {
	Root             string
	ConfigPath       string
	EnvPath          string
	InstructionsDir  string
	SlashCommandsDir string
	CommandsAllow    string
}

// DefaultPaths returns the default config paths for a repo root.
func DefaultPaths(root string) Paths {
	return Paths{
		Root:             root,
		ConfigPath:       filepath.Join(root, ".agent-layer", "config.toml"),
		EnvPath:          filepath.Join(root, ".agent-layer", ".env"),
		InstructionsDir:  filepath.Join(root, ".agent-layer", "instructions"),
		SlashCommandsDir: filepath.Join(root, ".agent-layer", "slash-commands"),
		CommandsAllow:    filepath.Join(root, ".agent-layer", "commands.allow"),
	}
}
