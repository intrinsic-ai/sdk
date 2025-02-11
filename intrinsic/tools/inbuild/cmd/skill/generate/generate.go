// Copyright 2023 Intrinsic Innovation LLC

// Package generate defines the `inbuild skill generate` command.
package generate

import (
	"github.com/spf13/cobra"
	"intrinsic/tools/inbuild/cmd/skill/generate/config"
	"intrinsic/tools/inbuild/cmd/skill/generate/entrypoint"
)

// GenerateCmd organizes commands for generating code for skills.
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Commands for generating code for skills",
	Long:  "Commands for generating code for skills for Flowstate.",
}

// The init function adds subcommands to `inbuild skill generate`.
func init() {
	GenerateCmd.AddCommand(entrypoint.EntryPointCmd)
	GenerateCmd.AddCommand(config.ConfigCmd)
}
