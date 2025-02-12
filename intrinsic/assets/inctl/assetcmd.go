// Copyright 2023 Intrinsic Innovation LLC

// Package assetcmd contains the root command for the inctl asset command.
package assetcmd

import (
	"github.com/spf13/cobra"
	"intrinsic/assets/inctl/list"
	"intrinsic/assets/inctl/listreleased"
	"intrinsic/assets/inctl/listreleasedversions"
	"intrinsic/assets/inctl/uninstall"
	"intrinsic/tools/inctl/cmd/root"
)

var assetCmd = &cobra.Command{
	Use:   root.AssetCmdName,
	Short: "Manages assets",
	Long:  "Manages assets",
}

func init() {
	assetCmd.AddCommand(list.GetCommand())
	assetCmd.AddCommand(listreleased.GetCommand())
	assetCmd.AddCommand(listreleasedversions.GetCommand())
	assetCmd.AddCommand(uninstall.GetCommand())

	root.RootCmd.AddCommand(assetCmd)
}
