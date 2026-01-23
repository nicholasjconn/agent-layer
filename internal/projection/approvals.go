package projection

import "github.com/conn-castle/agent-layer/internal/config"

// Approvals captures the resolved approvals policy and allowlist.
type Approvals struct {
	AllowCommands bool
	AllowMCP      bool
	Commands      []string
}

// BuildApprovals resolves approvals.mode into per-feature flags.
func BuildApprovals(cfg config.Config, commands []string) Approvals {
	allowCommands := cfg.Approvals.Mode == "all" || cfg.Approvals.Mode == "commands"
	allowMCP := cfg.Approvals.Mode == "all" || cfg.Approvals.Mode == "mcp"

	return Approvals{
		AllowCommands: allowCommands,
		AllowMCP:      allowMCP,
		Commands:      commands,
	}
}
