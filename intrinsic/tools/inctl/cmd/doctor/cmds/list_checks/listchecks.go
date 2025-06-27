// Copyright 2023 Intrinsic Innovation LLC

// Package listchecks contains the command for the inctl doctor list_checks command.
package listchecks

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"intrinsic/tools/inctl/cmd/doctor/checks/checks"
	"intrinsic/tools/inctl/util/printer"
)

// ListChecksCmd is the entry point for the inctl doctor list_checks command.
var ListChecksCmd = &cobra.Command{
	Use:   "list_checks",
	Short: "List all available checks",
	RunE:  listChecksCommandE,
}

func listChecksCommandE(cmd *cobra.Command, args []string) error {
	out, ok := printer.AsPrinter(cmd.OutOrStdout(), printer.TextOutputFormat)
	if !ok {
		return fmt.Errorf("invalid output configuration")
	}

	return runListChecksCmd(out)
}

func runListChecksCmd(prtr printer.Printer) error {
	// Map keys are not guaranteed to be sorted, so we sort them arbitrarily here.
	sortedKeys := make([]string, 0, len(checks.Checks))
	for k := range checks.Checks {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, checkName := range sortedKeys {
		check, ok := checks.Checks[checkName]
		if !ok {
			return fmt.Errorf("check '%s' unexpectedly not in the checks map", checkName)
		}
		prtr.Print(check.Name + ": " + check.Description)
	}
	return nil
}
