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
	"errors"
	"fmt"
	"io"

	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"
	lineconfigutils "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common/line_config_utils"
	servicedeletionutils "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common/service_deletion_utils"

	provisionerpb "intrinsic/platform/pubsub/cloud_router_provisioner/v1/provisioner_go_proto"
	endpointpb "intrinsic/platform/pubsub/connect/common/proto/line_configuration/v1/endpoint_spec_go_proto"

	"google.golang.org/grpc"
)

// HubServiceCreateRunner implements high-level logic of the hub-service-create command.
//
// Its primary responsibility is to decide what kind of line orchestration network to create:
// a local-only network that consists entirely of workcells, or a mixed network that may
// include VMs.
type HubServiceCreateRunner struct {
	project                         string
	org                             string
	hubEndpoint                     string
	spokeEndpoints                  []string
	forceLocalOnly                  bool
	shouldRetainServiceAsset        bool
	ignoreOnpremErrorsDuringCleanup bool

	// dialOnpremCluster creates a connection to the given spoke cluster.
	// The default implementation connects to the real workcell or VM.
	// Unit tests can provide an implementation that connects to a fake server.
	dialOnpremCluster func(context.Context, string, string, string) (context.Context, *grpc.ClientConn, string, error)

	// dialCloudCluster connects to the cloud-robotics cluster.
	// The default implementation connects to the real cloud-robotics cluster.
	// Unit tests can provide an implementation that connects to a fake server.
	dialCloudCluster func(context.Context) (*grpc.ClientConn, error)

	// getServiceVersionToInstall determines the version of the service asset that should be installed
	// on workcells / VMs.
	// The default implementation may call remote APIs.
	// Unit tests can provide an implementation that returns a predefined version.
	getServiceVersionToInstall func(string, string) (string, error)
}

// provisionLineLevelRouter provisions a line-level router for the given line.
func (r *HubServiceCreateRunner) provisionLineLevelRouter(ctx context.Context, conn *grpc.ClientConn, lineId string) error {
	client := provisionerpb.NewLineRouterProvisionerClient(conn)
	_, err := client.Provision(ctx, &provisionerpb.ProvisionLineRouterRequest{
		LineId: lineId,
	})
	if err != nil {
		return fmt.Errorf("failed to invoke cloud API: %w", err)
	}

	return nil
}

// installOnpremToCloudLineRouterRelay installs the relay service that connects a workcell / VM
// to the line-level router in the cloud.
func (r *HubServiceCreateRunner) installOnpremToCloudLineRouterRelay(ctx context.Context, out io.Writer, cluster, lineId, versionToInstall string) error {
	ctx, conn, _, err := r.dialOnpremCluster(ctx, r.project, r.org, cluster)
	if err != nil {
		return err
	}
	defer conn.Close()

	runner := &OnpremToLineRouterHubServiceCreateCmdRunner{
		ServiceInstallingCmdRunner: common.ServiceInstallingCmdRunner{
			CmdRunnerBase: *common.NewCmdRunnerBase(
				conn,
				out,
				cluster,
				common.OnpremToLineRouterRelayServicePackage,
				common.OnpremToLineRouterRelayServiceName),
			RequestedVersion: versionToInstall,
		},
		orgId:  r.org,
		lineId: lineId,
	}
	return runner.run(ctx)
}

// createLocalOnlyLineOrchestrationNetwork creates a local-only line orchestration
// network that consists entirely of workcells.
//
// Traffic in the local-only network flows through a single relay service installed
// on the hub workcell.
func (r *HubServiceCreateRunner) createLocalOnlyLineOrchestrationNetwork(
	ctx context.Context,
	out io.Writer,
	hubClusterId string,
	spokeEndpointSpecs []*endpointpb.EndpointSpec) error {

	versionToInstall, err := r.getServiceVersionToInstall(common.HubServicePackage, common.HubServiceName)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		out,
		"Creating a local-only line orchestration network with the hub at %q\n",
		hubClusterId)
	ctx, conn, _, err := r.dialOnpremCluster(ctx, r.project, r.org, hubClusterId)
	if err != nil {
		return err
	}
	defer conn.Close()

	runner := &LocalOnlyHubServiceCreateCmdRunner{
		ServiceInstallingCmdRunner: common.ServiceInstallingCmdRunner{
			CmdRunnerBase: *common.NewCmdRunnerBase(
				conn,
				out,
				hubClusterId,
				common.HubServicePackage,
				common.HubServiceName),
			RequestedVersion: versionToInstall,
		},
		spokeEndpointSpecs: spokeEndpointSpecs,
	}

	return runner.run(ctx)
}

// createMixedLineOrchestrationNetwork creates a mixed line orchestration
// network that may consist of workcells and VMs.
//
// Traffic in the mixed network flows through multiple Zenoh routers:
//   - Relay routers installed on _each_ participating workcell / VM.
//   - Line-level router in the cloud.
//
// This function installs all those routers.
func (r *HubServiceCreateRunner) createMixedLineOrchestrationNetwork(
	ctx context.Context,
	conn *grpc.ClientConn,
	out io.Writer,
	lineId string,
	hubClusterId string,
	spokeEndpointSpecs []*endpointpb.EndpointSpec) error {

	versionToInstall, err := r.getServiceVersionToInstall(
		common.OnpremToLineRouterRelayServicePackage,
		common.OnpremToLineRouterRelayServiceName)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "Creating a mixed line orchestration network. Line id: %q\n", lineId)
	fmt.Fprintf(out, "Provisioning a line-level router...\n")
	if err := r.provisionLineLevelRouter(ctx, conn, lineId); err != nil {
		return fmt.Errorf("failed to provision line-level router: %w", err)
	}
	hubAndSpokeClusterIds := make([]string, 0, len(spokeEndpointSpecs)+1)
	hubAndSpokeClusterIds = append(hubAndSpokeClusterIds, hubClusterId) // Hub
	for _, spec := range spokeEndpointSpecs {
		hubAndSpokeClusterIds = append(hubAndSpokeClusterIds, spec.GetWorkcellName()) // Spoke
	}
	fmt.Fprintf(
		out,
		"Line-level router has been provisioned. Will install relay routers on each of the %d clusters\n",
		len(hubAndSpokeClusterIds))
	for _, clusterID := range hubAndSpokeClusterIds {
		fmt.Fprintf(out, "Installing the relay router on %q\n", clusterID)
		err := r.installOnpremToCloudLineRouterRelay(ctx, out, clusterID, lineId, versionToInstall)
		if err != nil {
			return fmt.Errorf("failed to install a relay router on %q: %w", clusterID, err)
		}
	}
	fmt.Fprintf(
		out,
		"All relay routers have been installed, the line orchestration network is ready.\n")
	return nil
}

// cleanupLineIfExists checks if a line orchestration network already exists,
// and deletes resources that are no longer necessary.
func (r *HubServiceCreateRunner) cleanupLineIfExists(
	ctx context.Context,
	out io.Writer,
	conn *grpc.ClientConn,
	lineID string,
	hubEndpointSpec *endpointpb.EndpointSpec,
	spokeEndpointSpecs []*endpointpb.EndpointSpec) error {
	hubClusterID := hubEndpointSpec.GetWorkcellName()
	fmt.Fprintf(
		out,
		"Checking if the line orchestration network with the hub at %q already exists\n",
		hubClusterID)
	oldLineConfig, err := lineconfigutils.FetchLineConfiguration(ctx, conn, lineID)
	if err != nil {
		if errors.Is(err, common.ErrConfigStorageNotImplemented) {
			common.ExplainConfigStorageNotImplementedError(out, "hub-service-create", common.KeyForceLocalOnly, "create")
			return err
		}
		if errors.Is(err, common.ErrLineConfigNotFound) {
			fmt.Fprintf(
				out,
				"Line with the hub at %q not found, will create a new one.\n",
				hubClusterID)
			return nil
		}
		return fmt.Errorf("failed to check existing line configuration: %w", err)
	}

	fmt.Fprintf(out, "Line orchestration network with the hub at %q found.\n", hubClusterID)
	cleaner := oldLineCleaner{
		ServiceDeleter: servicedeletionutils.ServiceDeleter{
			ProjectID:                r.project,
			OrgID:                    r.org,
			ShouldRetainServiceAsset: r.shouldRetainServiceAsset,
			IgnoreOnpremErrors:       r.ignoreOnpremErrorsDuringCleanup,
			DialOnpremCluster:        r.dialOnpremCluster,
		},
		lineID:             lineID,
		hubEndpointSpec:    hubEndpointSpec,
		spokeEndpointSpecs: spokeEndpointSpecs,
		oldLineConfig:      oldLineConfig,
		remoteEndpointsPresentInOldLine: common.RemoteEndpointsExist(
			oldLineConfig.GetSpokeEndpoints(),
			oldLineConfig.GetHubEndpoint()),
		remoteEndpointsPresentInNewLine: common.RemoteEndpointsExist(spokeEndpointSpecs, hubEndpointSpec),
	}

	return cleaner.cleanup(ctx, out, conn)
}

// run implements high-level logic of the hub-service-create command.
func (r *HubServiceCreateRunner) run(ctx context.Context, out io.Writer) error {
	spokeEndpointSpecs, err := common.ParseEndpointSpecs(r.spokeEndpoints)
	if err != nil {
		return err
	}
	if len(spokeEndpointSpecs) == 0 {
		return fmt.Errorf("at least one spoke endpoint must be specified using --spoke-endpoint")
	}

	hubEndpointSpec, err := common.ParseEndpointSpec(r.hubEndpoint)
	if err != nil {
		return err
	}

	remoteEndpointsPresent := common.RemoteEndpointsExist(spokeEndpointSpecs, hubEndpointSpec)
	if r.forceLocalOnly {
		if remoteEndpointsPresent {
			return fmt.Errorf(
				"The '--%v' flag is present, but at least one endpoint is remote",
				common.KeyForceLocalOnly)
		}
		fmt.Fprintf(
			out,
			`
The '--%v' flag is present.
Creating a local-only line orchestration network without saving its configuration in the cloud.
`,
			common.KeyForceLocalOnly)
		return r.createLocalOnlyLineOrchestrationNetwork(
			ctx,
			out,
			hubEndpointSpec.GetWorkcellName(),
			spokeEndpointSpecs)
	}

	fmt.Fprintf(out, "Connecting to the cloud API endpoint\n")
	conn, err := r.dialCloudCluster(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to the cloud: %w", err)
	}
	defer conn.Close()

	lineID, err := common.MakeLineID(r.org, hubEndpointSpec.GetWorkcellName())
	if err != nil {
		fmt.Fprintf(out, "Cannot create a valid line id from %q and %q\n", r.org, hubEndpointSpec.GetWorkcellName())
		return err
	}

	err = r.cleanupLineIfExists(ctx, out, conn, lineID, hubEndpointSpec, spokeEndpointSpecs)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "Saving line configuration\n")
	err = lineconfigutils.SaveLineConfiguration(ctx, conn, lineID, hubEndpointSpec, spokeEndpointSpecs)
	if err != nil {
		return fmt.Errorf("failed to save line configuration: %w", err)
	}

	var result error = nil
	if remoteEndpointsPresent {
		result = r.createMixedLineOrchestrationNetwork(
			ctx,
			conn,
			out,
			lineID,
			hubEndpointSpec.GetWorkcellName(),
			spokeEndpointSpecs)
	} else {
		result = r.createLocalOnlyLineOrchestrationNetwork(
			ctx,
			out,
			hubEndpointSpec.GetWorkcellName(),
			spokeEndpointSpecs)
	}

	if result != nil {
		fmt.Fprintf(out, `
------------------------------------------------------------------
Failed to create a line orchestration network: %v

That error occurred after network configuration had been persisted
in the cloud. Please run the following command to delete that
configuration:

inctl pubsub hub-service-delete --hub-endpoint=%v --org=%v
------------------------------------------------------------------
`,
			result, r.hubEndpoint, r.org)
	}

	return result
}
