// Copyright 2023 Intrinsic Innovation LLC

// Package solutionversion contains all commands for solution version handling.
package solutionversion

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/orgutil"
)

var (
	viperLocal = viper.New()
)

// SolutionVersionCmd is the `inctl solution_version` command.
var SolutionVersionCmd = orgutil.WrapCmd(&cobra.Command{
	Use:                root.SolutionVersionCmdName,
	Aliases:            []string{"sv"},
	Short:              "SolutionVersion interacts with solution versions",
	DisableFlagParsing: true,
}, viperLocal)

func init() {
	root.RootCmd.AddCommand(SolutionVersionCmd)
}
