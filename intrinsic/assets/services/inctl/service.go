// Copyright 2023 Intrinsic Innovation LLC

// Package service contains all commands for handling service assets.
package service

import (
	"intrinsic/assets/inctl/install"
	"intrinsic/assets/services/inctl/add"
	deletecmd "intrinsic/assets/services/inctl/delete"
	"intrinsic/assets/services/inctl/release"
	servicestate "intrinsic/assets/services/inctl/state/state"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
)

func init() {
	cmd := cobrautil.ParentOfNestedSubcommands(root.ServiceCmdName, "Manage Service assets.")
	cmd.AddCommand(add.GetCommand())
	cmd.AddCommand(deletecmd.GetCommand())
	cmd.AddCommand(install.GetCommand())
	cmd.AddCommand(release.GetCommand())
	cmd.AddCommand(servicestate.ServiceStateCmd)
	root.RootCmd.AddCommand(cmd)
}
