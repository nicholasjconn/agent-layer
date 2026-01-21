package main

import (
	"fmt"

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
			warnings, err := sync.Run(root)
			if err != nil {
				return err
			}
			// Print warnings to stderr
			for _, w := range warnings {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Warning: %s\n", w.Message)
			}
			return nil
		},
	}

	return cmd
}
