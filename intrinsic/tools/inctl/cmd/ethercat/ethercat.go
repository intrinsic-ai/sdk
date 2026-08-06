// Copyright 2023 Intrinsic Innovation LLC

// Package ethercat contains the externally available commands for ethercat handling.
package ethercat

import (
	"intrinsic/assets/clientutils"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
)

var (
	// EtherCATCmd is the `inctl ethercat` command.
	EtherCATCmd = cobrautil.ParentOfNestedSubcommands(
		"ethercat", "Workcell EtherCAT handling")

	// Defined here, so that it can be replaced in tests.
	clientutilsDialClusterFromInctl = clientutils.DialClusterFromInctl
)

func init() {
	root.RootCmd.AddCommand(EtherCATCmd)
}
