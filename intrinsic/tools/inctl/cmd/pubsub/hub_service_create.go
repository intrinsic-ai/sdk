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

package pubsub

import (
	"context"
	"fmt"
	"io"
	"strings"

	"google.golang.org/grpc"

	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/tools/inctl/auth/auth"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyHubEndpoint       = "hub-endpoint"
	keySpokeEndpoints    = "spoke-endpoint"
	keyHubServiceVersion = "hub-service-version"

	// The `cluster` flag in `hub-service-create / delete` commands is deprecated
	// because it can only contain a cluster id. `hub-service-create` and `delete`
	// commands need additional information (whether that cluster is on a local
	// network or in the cloud). This information can be provided in the `hub-endpoint`
	// flag.
	keyClusterDeprecated = "cluster"

	endpointSpecSeparator     = "@"
	localEndpointDesignation  = "local"
	remoteEndpointDesignation = "remote"
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
	versionToInstall := e.cmdFlags.GetString(keyHubServiceVersion)
	var err error
	if len(versionToInstall) == 0 {
		versionToInstall, err = getDefaultVersion(
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

// getHubEndpoint computes the hub endpoint spec.
//
// It is intended to preserve backwards compatibility with previous versions of
// the hub-service-create command, where the `--cluster` flag was used to specify
// the hub. Newer versions use the `--hub-endpoint` flag. Unlike `--cluster`, that
// flag allows `@local` and `@remote` suffixes. This, in turn, enables users to
// unambiguously specify whether the hub is on a local network or a remote server.
func getHubEndpoint(out io.Writer, cluster, hubEndpoint string) (string, error) {
	if len(hubEndpoint) != 0 {
		return hubEndpoint, nil
	}

	if len(cluster) != 0 {
		fmt.Fprintf(out, "WARNING: the --cluster flag is deprecated, use --%s instead.\n", keyHubEndpoint)
		if strings.HasPrefix(cluster, "vmp-") {
			fmt.Fprintf(out, "Assuming that %q is a remote server.\n", cluster)
			return fmt.Sprintf("%s%s%s", cluster, endpointSpecSeparator, remoteEndpointDesignation), nil
		}

		fmt.Fprintf(out, "Assuming that %q is on the local network.\n", cluster)
		return fmt.Sprintf("%s%s%s", cluster, endpointSpecSeparator, localEndpointDesignation), nil
	}

	return "", fmt.Errorf("neither cluster nor hub endpoint are specified")
}

// RunE sets up the execution environment and invokes HubServiceCreateRunner.run.
func (e *hubServiceCreateCmdEnvironment) RunE(cmd *cobra.Command, _ []string) error {
	project := e.cmdFlags.GetFlagProject()
	org := e.cmdFlags.GetFlagOrganization()

	hubEndpoint, err := getHubEndpoint(cmd.OutOrStdout(), e.cmdFlags.GetString(keyClusterDeprecated), e.cmdFlags.GetString(keyHubEndpoint))
	if err != nil {
		return fmt.Errorf("could not get hub endpoint: %w", err)
	}

	runner := &HubServiceCreateRunner{
		project:        project,
		org:            org,
		hubEndpoint:    hubEndpoint,
		spokeEndpoints: e.cmdFlags.GetStringSlice(keySpokeEndpoints),
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
		keySpokeEndpoints,
		[]string{},
		"Spoke endpoint specifications (<workcell_name>@{local|remote|url})")
	flags.OptionalString(
		keyHubServiceVersion,
		"",
		"Version of the service asset to install. If not specified, the current default version will be installed.")
	flags.OptionalString(keyHubEndpoint, "", "Hub endpoint specification (<workcell_name>@{local|remote|url})")
	flags.OptionalString(keyClusterDeprecated, "", "Hub cluster (DEPRECATED, use hub-endpoint instead)")
	cmd.MarkFlagsMutuallyExclusive(keyHubEndpoint, keyClusterDeprecated)
	cmd.MarkFlagsOneRequired(keyHubEndpoint, keyClusterDeprecated)

	return cmd
}

func init() {
	PubsubCmd.AddCommand(NewHubServiceCreateCmd())
}
