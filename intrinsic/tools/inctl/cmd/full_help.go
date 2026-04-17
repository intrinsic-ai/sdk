// Copyright 2023 Intrinsic Innovation LLC

package root

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// PrintFullTree prints the entire command tree of the given command.
func PrintFullTree(cmd *cobra.Command, indent string) {
	if cmd.Hidden {
		return
	}

	out := cmd.OutOrStdout()

	fmt.Fprintf(out, "%s%s: %s\n", indent, cmd.Name(), cmd.Short)

	// Print local flags
	if usages := cmd.LocalFlags().FlagUsages(); usages != "" {
		for _, line := range strings.Split(usages, "\n") {
			if line != "" {
				fmt.Fprintf(out, "%s  %s\n", indent, line)
			}
		}
	}

	// Recurse into subcommands
	for _, subCmd := range cmd.Commands() {
		PrintFullTree(subCmd, indent+"    ")
	}
}
