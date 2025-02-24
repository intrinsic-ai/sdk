// Copyright 2023 Intrinsic Innovation LLC

// Package cobrautil provides common cobra utility functions.
package cobrautil

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ParentOfNestedSubcommandsPreRun returns the parent command used for nested subcommands.
// persistentPreRunE can be passed in to pre-process shared arguments.
func ParentOfNestedSubcommandsPreRun(use string, short string, persistentPreRunE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		// While this only changes the output by a single line, cobra defaults to returning 0
		// when it cannot find a subcommand.
		// This ensures that there's a proper error code for the shell to handle.
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("%s requires a valid subcommand.\n%s", cmd.Name(), cmd.UsageString())
		},
		PersistentPreRunE: persistentPreRunE,
	}
}

// ParentOfNestedSubcommands returns the parent command used for nested subcommands.
func ParentOfNestedSubcommands(use string, short string) *cobra.Command {
	return ParentOfNestedSubcommandsPreRun(use, short, nil)
}
