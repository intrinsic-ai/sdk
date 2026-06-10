// Copyright 2023 Intrinsic Innovation LLC

// Package solution contains all commands for solution handling.
package solution

import (
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/cmd/solution/create/create"
	"intrinsic/tools/inctl/cmd/solution/delete/delete"
	"intrinsic/tools/inctl/cmd/solution/get/get"
	"intrinsic/tools/inctl/cmd/solution/list/list"
	"intrinsic/tools/inctl/cmd/solution/share/share"
	"intrinsic/tools/inctl/cmd/solution/start/start"
	"intrinsic/tools/inctl/cmd/solution/stop/stop"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:                root.SolutionCmdName,
		Aliases:            []string{root.SolutionsCmdName},
		Short:              "Solution interacts with solutions",
		DisableFlagParsing: true,
	}

	cmd.AddCommand(get.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(share.NewCommand())
	cmd.AddCommand(start.NewCommand())
	cmd.AddCommand(stop.NewCommand())
	cmd.AddCommand(delete.NewCommand())

	root.RootCmd.AddCommand(cmd)
}
