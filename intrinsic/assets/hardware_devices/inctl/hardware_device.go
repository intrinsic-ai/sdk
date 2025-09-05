// Copyright 2023 Intrinsic Innovation LLC

// Package hardwaredevice contains all commands for handling HardwareDevices.
package hardwaredevice

import (
	"github.com/spf13/cobra"
	"intrinsic/assets/hardware_devices/inctl/install"
	"intrinsic/assets/hardware_devices/inctl/release"
	"intrinsic/tools/inctl/cmd/root"
)

// hardwareDeviceCmd is the super-command for everything to manage HardwareDevices.
var hardwareDeviceCmd = &cobra.Command{
	Use:   root.HardwareDeviceCmdName,
	Short: "Manages HardwareDevices.",
	Long:  "Manages HardwareDevices.",
}

func init() {
	hardwareDeviceCmd.AddCommand(install.GetCommand())
	hardwareDeviceCmd.AddCommand(release.GetCommand())
	root.RootCmd.AddCommand(hardwareDeviceCmd)
}
