// Copyright 2023 Intrinsic Innovation LLC

// Package service defines the `inbuild service` command.
package service

import (
	"github.com/spf13/cobra"
	"intrinsic/tools/inbuild/cmd/service/bundle"
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
