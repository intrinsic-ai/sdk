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

package hubservicedelete

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"
	servicedeletionutils "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common/service_deletion_utils"
	pubsubcmd "intrinsic/tools/inctl/cmd/pubsub/pubsub_cmd"
	"intrinsic/tools/inctl/util/agents"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	hubServiceDeleteViper = viper.New()
)

// serviceDeleteCmdEnvironment is the execution environment for commands
// that delete services used for line orchestration.
type serviceDeleteCmdEnvironment struct {
	cmdFlags *cmdutils.CmdFlags
}

// RunE sets up the execution environment and invokes HubServiceDeleteRunner.run.
func (e *serviceDeleteCmdEnvironment) RunE(cmd *cobra.Command, _ []string) error {
	if err := agents.Check(cmd); err != nil {
		return err
	}
	ctx := cmd.Context()
	project := e.cmdFlags.GetFlagProject()
	org := e.cmdFlags.GetFlagOrganization()
	hubEndpoint, err := common.GetHubEndpoint(cmd.OutOrStdout(), e.cmdFlags.GetString(common.KeyClusterDeprecated), e.cmdFlags.GetString(common.KeyHubEndpoint))
	if err != nil {
		return fmt.Errorf("could not get hub endpoint: %w", err)
	}

	runner := &HubServiceDeleteRunner{
		ServiceDeleter: servicedeletionutils.ServiceDeleter{
			ProjectID:                project,
			OrgID:                    org,
			IgnoreOnpremErrors:       e.cmdFlags.GetBool(common.KeyIgnoreOnpremErrors),
			ShouldRetainServiceAsset: e.cmdFlags.GetBool(common.KeyRetainServiceAsset),
			DialOnpremCluster: func(ctx context.Context, project, org, cluster string) (context.Context, *grpc.ClientConn, string, error) {
				return clientutils.DialCluster(ctx, project, org, "" /* address */, cluster, "" /* solution */)
			},
		},

		HubEndpoint:    hubEndpoint,
		ForceLocalOnly: e.cmdFlags.GetBool(common.KeyForceLocalOnly),
		DialCloudCluster: func(ctx context.Context) (*grpc.ClientConn, error) {
			return auth.NewCloudConnection(ctx, auth.WithFlagValues(hubServiceDeleteViper))
		},
	}

	return runner.Run(ctx, cmd.OutOrStdout())
}

// NewHubServiceDeleteCmd returns the initialized cobra command for
// deletion of a factory line network.
func NewHubServiceDeleteCmd(use, short string) *cobra.Command {
	flags := cmdutils.NewCmdFlagsWithViper(hubServiceDeleteViper)
	env := serviceDeleteCmdEnvironment{
		cmdFlags: flags,
	}

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.NoArgs,
		RunE:  env.RunE,
	}

	flags.SetCommand(cmd)
	flags.AddFlagsProjectOrg()

	flags.OptionalString(common.KeyHubEndpoint, "", "Hub endpoint specification (<workcell_name>@{local|remote|url})")
	flags.OptionalString(common.KeyClusterDeprecated, "", "Hub cluster (DEPRECATED, use hub-endpoint instead)")
	flags.OptionalBool(common.KeyForceLocalOnly, false, "Whether to delete a local-only network without trying to read its configuration")
	flags.OptionalBool(common.KeyIgnoreOnpremErrors, false, "Whether to proceed with deletion when uninstallation of onprem assets fails")
	cmd.MarkFlagsMutuallyExclusive(common.KeyHubEndpoint, common.KeyClusterDeprecated)
	cmd.MarkFlagsOneRequired(common.KeyHubEndpoint, common.KeyClusterDeprecated)

	// Can be useful during development, when a service built from source is sideloaded
	// into a solution. In this case, it may be better to keep it installed instead of
	// rebuilding it from source every time (it takes about 20 minutes).
	flags.OptionalBool(common.KeyRetainServiceAsset, false, "Whether to retain the service asset")

	return cmd
}

func init() {
	pubsubcmd.PubsubCmd.AddCommand(
		NewHubServiceDeleteCmd(
			"hub-service-delete",
			"Deletes the line orchestration network."))
}
