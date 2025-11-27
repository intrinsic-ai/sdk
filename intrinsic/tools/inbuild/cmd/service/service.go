// Copyright 2023 Intrinsic Innovation LLC

// Package service defines the `inbuild service` command.
package service

import (
	"intrinsic/tools/inbuild/cmd/service/bundle"

	"github.com/spf13/cobra"
)

// ServiceCmd organizes commands for building services.
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Commands for building services",
	Long:  "Commands for building services for Flowstate.",
}

// The init function adds subcommands to `inbuild service`.
func init() {
	ServiceCmd.AddCommand(bundle.BundleCmd)
}
