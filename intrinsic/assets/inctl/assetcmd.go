// Copyright 2023 Intrinsic Innovation LLC

// Package assetcmd contains the root command for the inctl asset command.
package assetcmd

import (
	"intrinsic/assets/inctl/getreleased"
	"intrinsic/assets/inctl/list"
	"intrinsic/assets/inctl/listreleased"
	"intrinsic/assets/inctl/listreleasedversions"
	"intrinsic/assets/inctl/uninstall"
	"intrinsic/assets/inctl/updatereleasemetadata"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
)

func init() {
	cmd := cobrautil.ParentOfNestedSubcommands(root.AssetCmdName, "Manage assets.")
	cmd.AddCommand(getreleased.GetCommand())
	cmd.AddCommand(list.GetCommand(""))
	cmd.AddCommand(listreleased.GetCommand())
	cmd.AddCommand(listreleasedversions.GetCommand())
	cmd.AddCommand(uninstall.GetCommand())
	cmd.AddCommand(updatereleasemetadata.GetCommand())
	root.RootCmd.AddCommand(cmd)
}
