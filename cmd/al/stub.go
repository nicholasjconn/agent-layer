package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newStubCmd(name string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("%s (not implemented yet)", name),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("%s is not implemented in this phase", name)
		},
	}
}
