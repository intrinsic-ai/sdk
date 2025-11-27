// Copyright 2023 Intrinsic Innovation LLC

// Package recordings provides an implementation of the recordings command.
package recordings

import (
	"intrinsic/assets/cmdutils"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Shared flags across subcommands.
var (
	flagBagID string
)

const (
	keyProjectShort = "p"
)

func setPrinterFromOutputFlag(command *cobra.Command, args []string) (err error) {
	if out, err := printer.NewPrinter(root.FlagOutput); err == nil {
		command.SetOut(out)
	}
	return
}

var (
	recordingsCmd = &cobra.Command{
		Use:   "recordings",
		Short: "Provides access to recordings for a given workcell.",
		Long:  "Provides access to recordings for a given workcell.",
		// Catching common typos and potential alternatives
		SuggestFor:        []string{"recording", "record", "bag"},
		PersistentPreRunE: setPrinterFromOutputFlag,
	}

	localViper = viper.New()
	cmdFlags   = cmdutils.NewCmdFlagsWithViper(localViper)
)

func init() {
	root.RootCmd.AddCommand(recordingsCmd)
	cmdFlags.SetCommand(recordingsCmd)
	cmdFlags.AddFlagsProjectOrg()
}
