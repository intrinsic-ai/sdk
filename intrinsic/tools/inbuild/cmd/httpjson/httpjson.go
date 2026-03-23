// Copyright 2023 Intrinsic Innovation LLC

// Package httpjson defines the `inbuild httpjson` command.
package httpjson

import (
	"intrinsic/tools/inbuild/cmd/httpjson/generatemain"

	"github.com/spf13/cobra"
)

// HttpJsonCmd organizes commands for building HTTP / JSON bridges.
var HttpJsonCmd = &cobra.Command{
	Use:   "httpjson",
	Short: "Commands for building HTTP / JSON bridges to gRPC services.",
	Long:  "Commands for building HTTP / JSON bridges to gRPC services.",
}

// The init function adds subcommands to `inbuild httpjson`.
func init() {
	HttpJsonCmd.AddCommand(generatemain.GenerateMainCmd)
}
