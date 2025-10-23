// Copyright 2023 Intrinsic Innovation LLC

// Package hardwaredevice contains all commands for handling HardwareDevices.
package hardwaredevice

import (
	"intrinsic/assets/inctl/install"
	"intrinsic/assets/inctl/release"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
)

func init() {
	cmd := cobrautil.ParentOfNestedSubcommands(root.HardwareDeviceCmdName, "Manage HardwareDevices.")
	cmd.AddCommand(install.GetCommand())
	cmd.AddCommand(release.GetCommand())
	root.RootCmd.AddCommand(cmd)
}
