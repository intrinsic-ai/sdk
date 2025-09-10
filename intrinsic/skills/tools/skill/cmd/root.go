// Copyright 2023 Intrinsic Innovation LLC

// Package cmd contains the root command for the skill installer tool.
package cmd

import (
	"github.com/spf13/cobra"
	"intrinsic/skills/tools/skill/cmd/create/create"
	"intrinsic/skills/tools/skill/cmd/install/install"
	"intrinsic/skills/tools/skill/cmd/list/list"
	"intrinsic/skills/tools/skill/cmd/logs/logs"
	"intrinsic/skills/tools/skill/cmd/release/release"
	"intrinsic/tools/inctl/cmd/root"
)

// SkillCmd is the super-command for everything skill management.
var SkillCmd = &cobra.Command{
	Use:   root.SkillCmdName,
	Short: "Manages skills",
	Long:  "Manages skills in a workcell, in a local repository or in the asset catalog",
}

func init() {
	SkillCmd.AddCommand(create.Command())
	SkillCmd.AddCommand(install.Command())
	SkillCmd.AddCommand(list.Command())
	SkillCmd.AddCommand(logs.Command())
	SkillCmd.AddCommand(release.Command())
	root.RootCmd.AddCommand(SkillCmd)
}
