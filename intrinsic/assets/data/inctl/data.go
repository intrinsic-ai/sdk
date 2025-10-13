// Copyright 2023 Intrinsic Innovation LLC

// Package data contains all commands for handling Data assets.
package data

import (
	"intrinsic/assets/data/inctl/release"
	"intrinsic/assets/inctl/install"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
)

func init() {
	cmd := cobrautil.ParentOfNestedSubcommands(root.DataCmdName, "Manage Data assets.")
	cmd.AddCommand(install.GetCommand())
	cmd.AddCommand(release.GetCommand())
	root.RootCmd.AddCommand(cmd)
}
