package main

import (
	"github.com/spf13/cobra"

	"github.com/nicholasjconn/agent-layer/internal/clients"
	"github.com/nicholasjconn/agent-layer/internal/clients/vscode"
	"github.com/nicholasjconn/agent-layer/internal/config"
)

func newVSCodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vscode",
		Short: "Sync and launch VS Code with CODEX_HOME configured",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveRepoRoot()
			if err != nil {
				return err
			}
			return clients.Run(root, "vscode", func(cfg *config.Config) *bool {
				return cfg.Agents.VSCode.Enabled
			}, vscode.Launch)
		},
	}

	return cmd
}
