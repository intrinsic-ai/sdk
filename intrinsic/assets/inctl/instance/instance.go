// Copyright 2023 Intrinsic Innovation LLC

// Package instance contains the inctl asset instance command.
package instance

import (
	"intrinsic/assets/inctl/instance/get"
	"intrinsic/assets/inctl/instance/list"

	"github.com/spf13/cobra"
)

// Command returns the parent command for asset instances and registers subcommands.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "instance",
		Aliases: []string{"instances"},
		Short:   "Manage Asset instances.",
	}
	cmd.AddCommand(get.Command())
	cmd.AddCommand(list.Command())
	return cmd
}
