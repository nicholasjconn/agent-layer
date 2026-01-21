package sync

import (
	"fmt"
	"sort"

	"github.com/nicholasjconn/agent-layer/internal/config"
)

// EstimatedCharsPerToken is the standard approximation for token estimation.
// This is a well-established heuristic: ~4 characters per token for English text.
const EstimatedCharsPerToken = 4

// Warning represents a sync-time warning message.
type Warning struct {
	Message string
}

// CheckInstructionSize checks if the generated instructions exceed the token threshold.
// Returns nil if the warning is disabled (threshold is nil).
func CheckInstructionSize(instructions []config.InstructionFile, cfg config.WarningsConfig) *Warning {
	if cfg.InstructionTokenThreshold == nil {
		return nil
	}

	threshold := *cfg.InstructionTokenThreshold

	var totalChars int
	for _, inst := range instructions {
		totalChars += len(inst.Content)
	}

	estimatedTokens := totalChars / EstimatedCharsPerToken
	if estimatedTokens > threshold {
		return &Warning{
			Message: fmt.Sprintf(
				"Generated instructions exceed token threshold: ~%d tokens (threshold: %d). Consider reducing instruction size.",
				estimatedTokens, threshold,
			),
		}
	}

	return nil
}

// CheckMCPServerCount checks if any enabled client has too many MCP servers.
// Returns warnings for each client that exceeds the threshold.
// Returns empty slice if the warning is disabled (threshold is nil).
func CheckMCPServerCount(servers []config.MCPServer, enabledClients []string, cfg config.WarningsConfig) []Warning {
	if cfg.MCPServerThreshold == nil {
		return nil
	}

	threshold := *cfg.MCPServerThreshold
	var warnings []Warning

	for _, client := range enabledClients {
		count := countEnabledServersForClient(servers, client)
		if count > threshold {
			warnings = append(warnings, Warning{
				Message: fmt.Sprintf(
					"Client %q has %d MCP servers enabled (threshold: %d). Consider reducing server count.",
					client, count, threshold,
				),
			})
		}
	}

	return warnings
}

// countEnabledServersForClient counts enabled MCP servers that apply to the given client.
func countEnabledServersForClient(servers []config.MCPServer, client string) int {
	count := 0
	for _, server := range servers {
		if server.Enabled == nil || !*server.Enabled {
			continue
		}
		if !server.AppliesToClient(client) {
			continue
		}
		count++
	}
	return count
}

// EnabledClientNames returns the names of enabled agents from the config.
func EnabledClientNames(agents config.AgentsConfig) []string {
	var names []string
	if agents.Gemini.Enabled != nil && *agents.Gemini.Enabled {
		names = append(names, "gemini")
	}
	if agents.Claude.Enabled != nil && *agents.Claude.Enabled {
		names = append(names, "claude")
	}
	if agents.Codex.Enabled != nil && *agents.Codex.Enabled {
		names = append(names, "codex")
	}
	if agents.VSCode.Enabled != nil && *agents.VSCode.Enabled {
		names = append(names, "vscode")
	}
	if agents.Antigravity.Enabled != nil && *agents.Antigravity.Enabled {
		names = append(names, "antigravity")
	}
	sort.Strings(names)
	return names
}
