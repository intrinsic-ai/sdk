// Copyright 2023 Intrinsic Innovation LLC

// Package skill defines the `inbuild skill` command.
package skill

import (
	"github.com/spf13/cobra"
	"intrinsic/tools/inbuild/cmd/skill/bundle"
	"intrinsic/tools/inbuild/cmd/skill/generate/generate"
)

// SkillCmd organizes commands for building skills.
var SkillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Commands for building skills",
	Long:  "Commands for building skills for Flowstate.",
}

// The init function adds subcommands to `inbuild skill`.
func init() {
	SkillCmd.AddCommand(bundle.BundleCmd)
	SkillCmd.AddCommand(generate.GenerateCmd)
}
