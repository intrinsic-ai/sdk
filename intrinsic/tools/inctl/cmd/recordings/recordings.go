// Copyright 2023 Intrinsic Innovation LLC

// Package recordings provides an implementation of the recordings command.
package recordings

import (
	"intrinsic/tools/inctl/cmd/root"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Exposed for testing
var (
	checkOrgExists = true
)

// Shared flags across subcommands.
var (
	flagBagID string
)

const (
	keyProjectShort = "p"
)

var RecordingsCmd = &cobra.Command{
	Use:   "recordings",
	Short: "Provides access to recordings for a given workcell.",
	Long:  "Provides access to recordings for a given workcell.",
	// Catching common typos and potential alternatives
	SuggestFor: []string{"recording", "record", "bag"},
}

func init() {
	root.RootCmd.AddCommand(RecordingsCmd)
}

// resolveCluster returns the cluster name, preferring the "cluster" parameter
// and falling back to the deprecated "workcell" parameter.
func resolveCluster(params *viper.Viper) string {
	cluster := params.GetString("cluster")
	if cluster == "" {
		cluster = params.GetString("workcell")
	}
	return cluster
}
