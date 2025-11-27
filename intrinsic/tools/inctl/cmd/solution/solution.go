// Copyright 2023 Intrinsic Innovation LLC

// Package solution contains all commands for solution handling.
package solution

import (
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/cmd/solution/delete/delete"
	"intrinsic/tools/inctl/cmd/solution/get/get"
	"intrinsic/tools/inctl/cmd/solution/list/list"
	"intrinsic/tools/inctl/cmd/solution/start/start"
	"intrinsic/tools/inctl/cmd/solution/stop/stop"
	"intrinsic/tools/inctl/util/orgutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	viperLocal := viper.New()

	SolutionCmd := orgutil.WrapCmd(&cobra.Command{
		Use:                root.SolutionCmdName,
		Aliases:            []string{root.SolutionsCmdName},
		Short:              "Solution interacts with solutions",
		DisableFlagParsing: true,
	}, viperLocal)

	SolutionCmd.AddCommand(get.NewCommand())
	SolutionCmd.AddCommand(list.NewCommand())
	SolutionCmd.AddCommand(start.NewCommand())
	SolutionCmd.AddCommand(stop.NewCommand())
	SolutionCmd.AddCommand(delete.NewCommand())

	root.RootCmd.AddCommand(SolutionCmd)
}
