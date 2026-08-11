// Copyright 2023 Intrinsic Innovation LLC

package world

import (
	"intrinsic/assets/cmdutils"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
)

var (
	flags *cmdutils.CmdFlags
)

var WorldCmd = cobrautil.ParentOfNestedSubcommands("world", "Manage and introspect worlds.")

func init() {
	flags = cmdutils.NewCmdFlags()
	flags.SetCommand(WorldCmd)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagsProjectOrg()

	root.RootCmd.AddCommand(WorldCmd)
}
