// Copyright 2023 Intrinsic Innovation LLC

// Package platformlevelswitch overrides the default IOC no-op implementations with enterprise-only CLI, identity, and authentication utilities.
package platformlevelswitch

import (
	"context"

	"intrinsic/kubernetes/acl/identity"
	"intrinsic/tools/inctl/util/orgutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	AppendOrgToOutgoingContext = func(ctx context.Context, org string) (context.Context, error) {
		return identity.AppendToOutgoingContext(ctx, identity.WithOrg(org))
	}

	WrapCmd = func(cmd *cobra.Command, vipr *viper.Viper) *cobra.Command {
		return orgutil.WrapCmd(cmd, vipr)
	}

	WrapCmdOptional = func(cmd *cobra.Command, vipr *viper.Viper) *cobra.Command {
		return orgutil.WrapCmdOptional(cmd, vipr)
	}
}
