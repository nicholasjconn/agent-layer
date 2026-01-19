package main

import (
	"github.com/spf13/cobra"

	"github.com/nicholasjconn/agent-layer/internal/sync"
)

func newSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Regenerate client outputs from .agent-layer",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := getwd()
			if err != nil {
				return err
			}
			return sync.Run(root)
		},
	}

	return cmd
}
