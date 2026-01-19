package config

// Config is the root configuration loaded from .agent-layer/config.toml.
type Config struct {
	Approvals ApprovalsConfig `toml:"approvals"`
	Agents    AgentsConfig    `toml:"agents"`
	MCP       MCPConfig       `toml:"mcp"`
}

// ApprovalsConfig controls auto-approval behavior per client.
type ApprovalsConfig struct {
	Mode string `toml:"mode"`
}

// AgentsConfig holds per-client enablement and model selection.
type AgentsConfig struct {
	Gemini      AgentConfig `toml:"gemini"`
	Claude      AgentConfig `toml:"claude"`
	Codex       CodexConfig `toml:"codex"`
	VSCode      AgentConfig `toml:"vscode"`
	Antigravity AgentConfig `toml:"antigravity"`
}

// AgentConfig is shared by agents that only need enablement and model selection.
type AgentConfig struct {
	Enabled *bool  `toml:"enabled"`
	Model   string `toml:"model"`
}

// CodexConfig extends AgentConfig with Codex-specific settings.
type CodexConfig struct {
	Enabled         *bool  `toml:"enabled"`
	Model           string `toml:"model"`
	ReasoningEffort string `toml:"reasoning_effort"`
}

// MCPConfig contains the external MCP servers configuration.
type MCPConfig struct {
	Servers []MCPServer `toml:"servers"`
}

// MCPServer defines a single MCP server entry.
type MCPServer struct {
	ID        string            `toml:"id"`
	Enabled   *bool             `toml:"enabled"`
	Clients   []string          `toml:"clients"`
	Transport string            `toml:"transport"`
	URL       string            `toml:"url"`
	Headers   map[string]string `toml:"headers"`
	Command   string            `toml:"command"`
	Args      []string          `toml:"args"`
	Env       map[string]string `toml:"env"`
}

// InstructionFile holds a single instruction fragment.
type InstructionFile struct {
	Name    string
	Content string
}

// SlashCommand represents a parsed slash command with metadata and body.
type SlashCommand struct {
	Name        string
	Description string
	Body        string
	SourcePath  string
}

// ProjectConfig is the fully loaded configuration state for sync and launch.
type ProjectConfig struct {
	Config        Config
	Env           map[string]string
	Instructions  []InstructionFile
	SlashCommands []SlashCommand
	CommandsAllow []string
	Root          string
}
