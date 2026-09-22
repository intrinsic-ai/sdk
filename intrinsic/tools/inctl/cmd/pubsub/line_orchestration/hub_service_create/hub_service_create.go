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

package hubservicecreate

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"
	pubsubcmd "intrinsic/tools/inctl/cmd/pubsub/pubsub_cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	hubServiceCreateViper = viper.New()
)

// hubServiceCreateCmdEnvironment is the execution environment for the
// hub-service-create command. That environment contains command line
// flags.
type hubServiceCreateCmdEnvironment struct {
	cmdFlags *cmdutils.CmdFlags
}

// getServiceVersionToInstall determines the version of the service asset that should
// be installed. It can be the version specified in the command line, or the default
// version in the asset catalog.
func (e *hubServiceCreateCmdEnvironment) getServiceVersionToInstall(cmd *cobra.Command, packageName, serviceName string) (string, error) {
	versionToInstall := e.cmdFlags.GetString(common.KeyHubServiceVersion)
	var err error
	if len(versionToInstall) == 0 {
		versionToInstall, err = common.GetDefaultVersion(
			cmd.Context(),
			e.cmdFlags,
			cmd.OutOrStdout(),
			packageName,
			serviceName)
		if err != nil {
			return "", fmt.Errorf(
				"failed to determine which version of %v to install: %w",
				serviceName,
				err)
		}
	}
	return versionToInstall, nil
}

// RunE sets up the execution environment and invokes HubServiceCreateRunner.run.
func (e *hubServiceCreateCmdEnvironment) RunE(cmd *cobra.Command, _ []string) error {
	project := e.cmdFlags.GetFlagProject()
	org := e.cmdFlags.GetFlagOrganization()

	hubEndpoint, err := common.GetHubEndpoint(cmd.OutOrStdout(), e.cmdFlags.GetString(common.KeyClusterDeprecated), e.cmdFlags.GetString(common.KeyHubEndpoint))
	if err != nil {
		return fmt.Errorf("could not get hub endpoint: %w", err)
	}

	runner := &HubServiceCreateRunner{
		project:        project,
		org:            org,
		hubEndpoint:    hubEndpoint,
		spokeEndpoints: e.cmdFlags.GetStringSlice(common.KeySpokeEndpoints),
		forceLocalOnly: e.cmdFlags.GetBool(common.KeyForceLocalOnly),
		dialOnpremCluster: func(ctx context.Context, project, org, cluster string) (context.Context, *grpc.ClientConn, string, error) {
			return clientutils.DialCluster(ctx, project, org, "" /* address */, cluster, "" /* solution */)
		},
		dialCloudCluster: func(ctx context.Context) (*grpc.ClientConn, error) {
			return auth.NewCloudConnection(ctx, auth.WithFlagValues(hubServiceCreateViper))
		},
		getServiceVersionToInstall: func(packageName string, serviceName string) (string, error) {
			return e.getServiceVersionToInstall(cmd, packageName, serviceName)
		},
	}
	return runner.run(cmd.Context(), cmd.OutOrStdout())
}

// NewHubServiceCreateCmd returns the initialized cobra command for hub-service-create.
func NewHubServiceCreateCmd() *cobra.Command {
	flags := cmdutils.NewCmdFlagsWithViper(hubServiceCreateViper)
	commandWrapper := &hubServiceCreateCmdEnvironment{cmdFlags: flags}

	cmd := &cobra.Command{
		Use:   "hub-service-create",
		Short: "Creates or updates the PubSub Hub service in the currently running solution.",
		Args:  cobra.NoArgs,
		RunE:  commandWrapper.RunE,
	}

	flags.SetCommand(cmd)
	flags.AddFlagsProjectOrg()

	flags.StringSlice(
		common.KeySpokeEndpoints,
		[]string{},
		"Spoke endpoint specifications (<workcell_name>@{local|remote|url})")
	flags.OptionalString(
		common.KeyHubServiceVersion,
		"",
		"Version of the service asset to install. If not specified, the current default version will be installed.")
	flags.OptionalString(common.KeyHubEndpoint, "", "Hub endpoint specification (<workcell_name>@{local|remote|url})")
	flags.OptionalString(common.KeyClusterDeprecated, "", "Hub cluster (DEPRECATED, use hub-endpoint instead)")
	flags.OptionalBool(common.KeyForceLocalOnly, false, "Whether to create a local-only network without saving its configuration")
	cmd.MarkFlagsMutuallyExclusive(common.KeyHubEndpoint, common.KeyClusterDeprecated)
	cmd.MarkFlagsOneRequired(common.KeyHubEndpoint, common.KeyClusterDeprecated)

	return cmd
}

func init() {
	pubsubcmd.PubsubCmd.AddCommand(NewHubServiceCreateCmd())
}
