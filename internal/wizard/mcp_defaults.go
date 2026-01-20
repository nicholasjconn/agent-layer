package wizard

import (
	"fmt"
	"strings"

	"github.com/nicholasjconn/agent-layer/internal/config"
	"github.com/nicholasjconn/agent-layer/internal/templates"
)

// DefaultMCPServer describes a default MCP server and its required env vars.
type DefaultMCPServer struct {
	ID          string
	RequiredEnv []string
}

// loadDefaultMCPServers returns default MCP servers derived from the template config.
func loadDefaultMCPServers() ([]DefaultMCPServer, error) {
	cfg, err := config.LoadTemplateConfig()
	if err != nil {
		return nil, err
	}
	defaults := make([]DefaultMCPServer, 0, len(cfg.MCP.Servers))
	for _, server := range cfg.MCP.Servers {
		required := config.RequiredEnvVarsForMCPServer(server)
		defaults = append(defaults, DefaultMCPServer{
			ID:          server.ID,
			RequiredEnv: required,
		})
	}
	if len(defaults) == 0 {
		return nil, fmt.Errorf("template config contains no MCP servers")
	}
	return defaults, nil
}

// missingDefaultMCPServers returns default MCP server IDs absent from the current config.
// defaults is the list of default servers; servers is the current config server list.
func missingDefaultMCPServers(defaults []DefaultMCPServer, servers []config.MCPServer) []string {
	existing := make(map[string]bool, len(servers))
	for _, srv := range servers {
		if srv.ID == "" {
			continue
		}
		existing[srv.ID] = true
	}

	var missing []string
	for _, def := range defaults {
		if !existing[def.ID] {
			missing = append(missing, def.ID)
		}
	}
	return missing
}

// appendMissingDefaultMCPServers appends template MCP server blocks for the provided IDs.
// content is the current config; missing lists server IDs to append; returns updated content or error.
func appendMissingDefaultMCPServers(content string, missing []string) (string, error) {
	if len(missing) == 0 {
		return content, nil
	}

	blocks, err := defaultMCPServerBlocks()
	if err != nil {
		return "", err
	}

	toAppend := make([]string, 0, len(missing))
	for _, id := range missing {
		block, ok := blocks[id]
		if !ok {
			return "", fmt.Errorf("missing default MCP server template for %q", id)
		}
		toAppend = append(toAppend, block)
	}

	trimmed := strings.TrimRight(content, "\n")
	if trimmed != "" {
		trimmed += "\n\n"
	}
	trimmed += strings.Join(toAppend, "\n\n")
	trimmed += "\n"
	return trimmed, nil
}

// defaultMCPServerBlocks loads MCP server blocks from the embedded config template.
// It returns a map of server ID to the block text as written in the template.
func defaultMCPServerBlocks() (map[string]string, error) {
	data, err := templates.Read("config.toml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config template: %w", err)
	}
	return parseMCPServerBlocks(string(data))
}

// parseMCPServerBlocks extracts [[mcp.servers]] blocks keyed by ID.
// content should be TOML text that includes MCP server blocks.
func parseMCPServerBlocks(content string) (map[string]string, error) {
	lines := strings.Split(content, "\n")
	blocks := make(map[string]string)

	for i := 0; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != "[[mcp.servers]]" {
			continue
		}

		start := i
		end := i + 1
		for end < len(lines) && strings.TrimSpace(lines[end]) != "[[mcp.servers]]" {
			end++
		}

		blockLines := lines[start:end]
		id, err := findMCPServerID(blockLines)
		if err != nil {
			return nil, err
		}
		if _, exists := blocks[id]; exists {
			return nil, fmt.Errorf("duplicate MCP server id %q in template", id)
		}
		blocks[id] = strings.TrimRight(strings.Join(blockLines, "\n"), "\n")
		i = end - 1
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("no MCP server blocks found in template")
	}

	return blocks, nil
}

// findMCPServerID finds the id field value within a server block.
// blockLines must include the [[mcp.servers]] header; returns an error if id is missing.
func findMCPServerID(blockLines []string) (string, error) {
	for _, line := range blockLines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "id") {
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "id"))
			if strings.HasPrefix(rest, "=") {
				val := strings.TrimSpace(strings.TrimPrefix(rest, "="))
				val = strings.Trim(val, "\"")
				if val == "" {
					return "", fmt.Errorf("empty MCP server id in template")
				}
				return val, nil
			}
		}
	}
	return "", fmt.Errorf("missing MCP server id in template block")
}
