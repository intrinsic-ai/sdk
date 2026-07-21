// Copyright 2023 Intrinsic Innovation LLC

// Package data defines the `inbuild data` command.
package data

import (
	"intrinsic/tools/inbuild/cmd/data/bundle"

	"github.com/spf13/cobra"
)

// DataCmd organizes commands for building Data Assets.
var DataCmd = &cobra.Command{
	Use:   "data",
	Short: "Commands for building Data Assets",
	Long:  "Commands for building Data Assets for Flowstate.",
}

// The init function adds subcommands to `inbuild data`.
func init() {
	DataCmd.AddCommand(bundle.BundleCmd)
}
