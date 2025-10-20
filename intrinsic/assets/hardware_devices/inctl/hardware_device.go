// Copyright 2023 Intrinsic Innovation LLC

// Package hardwaredevice contains all commands for handling HardwareDevices.
package hardwaredevice

import (
	"intrinsic/assets/hardware_devices/inctl/release"
	"intrinsic/assets/inctl/install"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
)

func init() {
	cmd := cobrautil.ParentOfNestedSubcommands(root.HardwareDeviceCmdName, "Manage HardwareDevices.")
	cmd.AddCommand(install.GetCommand())
	cmd.AddCommand(release.GetCommand())
	root.RootCmd.AddCommand(cmd)
}
