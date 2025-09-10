// Copyright 2023 Intrinsic Innovation LLC

// Package data contains all commands for handling Data assets.
package data

import (
	"github.com/spf13/cobra"
	"intrinsic/assets/data/inctl/install"
	"intrinsic/assets/data/inctl/release"
	"intrinsic/tools/inctl/cmd/root"
)

// dataCmd is the super-command for everything to manage Data assets.
var dataCmd = &cobra.Command{
	Use:   root.DataCmdName,
	Short: "Manages Data assets.",
	Long:  "Manages Data assets.",
}

func init() {
	dataCmd.AddCommand(install.GetCommand())
	dataCmd.AddCommand(release.GetCommand())

	root.RootCmd.AddCommand(dataCmd)
}
