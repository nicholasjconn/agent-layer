package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/conn-castle/agent-layer/internal/config"
	"github.com/conn-castle/agent-layer/internal/mcp"
)

var runPromptServer = mcp.RunPromptServer

func newMcpPromptsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp-prompts",
		Short: "Run the internal MCP prompt server over stdio",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveRepoRoot()
			if err != nil {
				return err
			}
			project, err := config.LoadProjectConfig(root)
			if err != nil {
				return err
			}
			return runPromptServer(context.Background(), Version, project.SlashCommands)
		},
	}

	return cmd
}
