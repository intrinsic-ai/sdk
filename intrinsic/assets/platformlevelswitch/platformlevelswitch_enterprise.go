// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
