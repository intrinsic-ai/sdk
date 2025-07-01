// Copyright 2023 Intrinsic Innovation LLC

// Package state contains commands for introspecting and modifying the state of a running
// service asset in a solution.
package state

import (
	"github.com/spf13/cobra"
	"intrinsic/assets/services/inctl/state/get"
	"intrinsic/assets/services/inctl/state/list"
)

// ServiceStateCmd is the super-command for commands to introspect and modify the state of a
// running service asset in a solution.
var ServiceStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Display or modify the state of a running Service instance in a solution.",
	Long:  `Display or modify the state of a running Service instance in a solution.`,
}

func init() {
	ServiceStateCmd.AddCommand(get.Command())
	ServiceStateCmd.AddCommand(list.Command())
}
