package main

import (
	"github.com/spf13/cobra"

	"github.com/nicholasjconn/agent-layer/internal/clients"
	"github.com/nicholasjconn/agent-layer/internal/clients/antigravity"
	"github.com/nicholasjconn/agent-layer/internal/config"
)

func newAntigravityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "antigravity",
		Short: "Sync and launch Antigravity",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveRepoRoot()
			if err != nil {
				return err
			}
			return clients.Run(root, "antigravity", func(cfg *config.Config) *bool {
				return cfg.Agents.Antigravity.Enabled
			}, antigravity.Launch)
		},
	}

	return cmd
}
