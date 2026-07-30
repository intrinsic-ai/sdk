// Copyright 2023 Intrinsic Innovation LLC

// Package platformlevelswitch reexports enterprise-only CLI, identity, and authentication utilities for enterprise builds, while defaulting to no-op implementations for IOC core builds.
package platformlevelswitch

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppendOrgToOutgoingContext appends the org option into outgoing context metadata.
var AppendOrgToOutgoingContext = func(ctx context.Context, org string) (context.Context, error) {
	return ctx, nil
}

// WrapCmd injects project/org flags and sets up required organization verification for cobra commands.
var WrapCmd = func(cmd *cobra.Command, vipr *viper.Viper) *cobra.Command {
	return cmd
}

// WrapCmdOptional injects optional project/org flags and handling into the command.
var WrapCmdOptional = func(cmd *cobra.Command, vipr *viper.Viper) *cobra.Command {
	return cmd
}
