package main

import (
	"github.com/spf13/cobra"

	"github.com/nicholasjconn/agent-layer/internal/install"
)

func newInstallCmd() *cobra.Command {
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Initialize Agent Layer in this repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := getwd()
			if err != nil {
				return err
			}
			return install.Run(root, install.Options{Overwrite: overwrite})
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing template files")

	return cmd
}
