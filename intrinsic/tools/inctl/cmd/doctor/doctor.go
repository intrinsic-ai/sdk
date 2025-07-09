// Copyright 2023 Intrinsic Innovation LLC

// Package doctor contains the root command for the inctl doctor command.
package doctor

import (
	"github.com/spf13/cobra"
	"intrinsic/tools/inctl/cmd/doctor/api/api"
	"intrinsic/tools/inctl/cmd/doctor/cmds/check/check"
	"intrinsic/tools/inctl/cmd/doctor/cmds/list_checks/listchecks"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"
)

func setPrinterFromOutputFlag(command *cobra.Command, args []string) (err error) {
	if out, err := printer.NewPrinter(root.FlagOutput); err == nil {
		command.SetOut(out)
	}
	return
}

var doctorCmd = &cobra.Command{
	Use:               root.DoctorCmdName,
	Short:             check.CheckCmd.Short,
	Long:              "Check your environment for common problems and return a non-zero exit code if any problems are found.",
	PersistentPreRunE: setPrinterFromOutputFlag,
}

func init() {
	api.CmdFlags.SetCommand(doctorCmd)
	api.CmdFlags.AddFlagsAddressClusterSolution()
	api.CmdFlags.AddFlagsProjectOrgOptional()

	doctorCmd.AddCommand(check.CheckCmd)
	doctorCmd.AddCommand(listchecks.ListChecksCmd)

	root.RootCmd.AddCommand(doctorCmd)
}
