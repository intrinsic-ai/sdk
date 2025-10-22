// Copyright 2023 Intrinsic Innovation LLC

// Package root contains the main entry point for the inbuild CLI.
package root

import (
	"os"

	"flag"
	"github.com/spf13/cobra"
	"intrinsic/production/intrinsic"
	"intrinsic/tools/inbuild/cmd/service/service"
	"intrinsic/tools/inbuild/cmd/skill/skill"
)

// RootCmd is the top level command of inbuild.
var RootCmd = &cobra.Command{
	Use:   "inbuild",
	Short: "inbuild builds assets",
	Long:  "inbuild builds assets for Flowstate.",
}

// Inbuild launches the main inbuild CLI.
func Inbuild() {
	intrinsic.Init()

	RootCmd.SetArgs(flag.Args())

	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

// The init function adds subcommands to `inbuild`.
func init() {
	RootCmd.AddCommand(service.ServiceCmd)
	RootCmd.AddCommand(skill.SkillCmd)
}
