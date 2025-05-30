// Copyright 2023 Intrinsic Innovation LLC

// Package assetcmd contains the root command for the inctl asset command.
package assetcmd

import (
	"github.com/spf13/cobra"
	"intrinsic/assets/inctl/getreleased"
	"intrinsic/assets/inctl/list"
	"intrinsic/assets/inctl/listreleased"
	"intrinsic/assets/inctl/listreleasedversions"
	"intrinsic/assets/inctl/uninstall"
	"intrinsic/assets/inctl/updatereleasemetadata"
	"intrinsic/tools/inctl/cmd/root"
)

var assetCmd = &cobra.Command{
	Use:   root.AssetCmdName,
	Short: "Manages assets",
	Long:  "Manages assets",
}

func init() {
	assetCmd.AddCommand(getreleased.GetCommand())
	assetCmd.AddCommand(list.GetCommand(""))
	assetCmd.AddCommand(listreleased.GetCommand())
	assetCmd.AddCommand(listreleasedversions.GetCommand())
	assetCmd.AddCommand(uninstall.GetCommand())
	assetCmd.AddCommand(updatereleasemetadata.GetCommand())

	root.RootCmd.AddCommand(assetCmd)
}
