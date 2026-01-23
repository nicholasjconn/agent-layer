package main

import (
	"github.com/spf13/cobra"

	"github.com/conn-castle/agent-layer/internal/clients"
	"github.com/conn-castle/agent-layer/internal/clients/gemini"
	"github.com/conn-castle/agent-layer/internal/config"
)

func newGeminiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gemini",
		Short: "Sync and launch Gemini CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveRepoRoot()
			if err != nil {
				return err
			}
			return clients.Run(root, "gemini", func(cfg *config.Config) *bool {
				return cfg.Agents.Gemini.Enabled
			}, gemini.Launch)
		},
	}

	return cmd
}
