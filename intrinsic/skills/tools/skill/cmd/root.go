// Copyright 2023 Intrinsic Innovation LLC

// Package cmd contains the root command for the skill installer tool.
package cmd

import (
	"intrinsic/skills/tools/skill/cmd/create/create"
	"intrinsic/skills/tools/skill/cmd/install/install"
	"intrinsic/skills/tools/skill/cmd/list/list"
	"intrinsic/skills/tools/skill/cmd/logs/logs"
	"intrinsic/skills/tools/skill/cmd/release/release"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
)

func init() {
	cmd := cobrautil.ParentOfNestedSubcommands(root.SkillCmdName, "Manage Skill assets.")
	cmd.AddCommand(create.Command())
	cmd.AddCommand(install.Command())
	cmd.AddCommand(list.Command())
	cmd.AddCommand(logs.Command())
	cmd.AddCommand(release.Command())
	root.RootCmd.AddCommand(cmd)
}
