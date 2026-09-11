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
	"slices"
	"strings"

	provisionerpb "intrinsic/platform/pubsub/cloud_router_provisioner/v1/provisioner_go_proto"
	endpointpb "intrinsic/platform/pubsub/connect/onprem/relay_router_service/endpoint_spec_go_proto"

	"google.golang.org/grpc"
)

// HubServiceCreateRunner implements high-level logic of the hub-service-create command.
//
// Its primary responsibility is to decide what kind of line orchestration network to create:
// a local-only network that consists entirely of workcells, or a mixed network that may
// include VMs.
type HubServiceCreateRunner struct {
	project        string
	org            string
	hubEndpoint    string
	spokeEndpoints []string

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

// parseEndpointSpec creates an EndpointSpec based on a spoke-endpoint or hub-endpoint command line flag.
func parseEndpointSpec(flagValue string) (*endpointpb.EndpointSpec, error) {
	parts := strings.Split(flagValue, endpointSpecSeparator)
	if len(parts) != 2 {
		return nil, fmt.Errorf(
			"Failed to parse %q. Each endpoint spec should consist of two parts separated by %v",
			flagValue, endpointSpecSeparator)
	}

	result := &endpointpb.EndpointSpec{
		WorkcellName: parts[0],
	}

	switch parts[1] {
	case localEndpointDesignation:
		result.ConnectionSpec = &endpointpb.EndpointSpec_Local{
			Local: &endpointpb.LocalConnectionSpec{},
		}
	case remoteEndpointDesignation:
		result.ConnectionSpec = &endpointpb.EndpointSpec_Remote{
			Remote: &endpointpb.RemoteConnectionSpec{},
		}
	default:
		result.ConnectionSpec = &endpointpb.EndpointSpec_Url{Url: parts[1]}
	}

	return result, nil
}

// parseEndpointSpecs parses all spoke endpoint specs.
func parseEndpointSpecs(endpointSpecStrings []string) ([]*endpointpb.EndpointSpec, error) {
	result := make([]*endpointpb.EndpointSpec, 0, len(endpointSpecStrings))
	for _, spec := range endpointSpecStrings {
		endpointSpec, err := parseEndpointSpec(spec)
		if err != nil {
			return nil, fmt.Errorf("failed to parse endpoint spec %q: %w", spec, err)
		}
		result = append(result, endpointSpec)
	}

	return result, nil
}

// isRemoteEndpoint checks whether the endpoint spec points to a remote endpoint.
func isRemoteEndpoint(endpointSpec *endpointpb.EndpointSpec) bool {
	return (endpointSpec.GetRemote() != nil)
}

// remoteEndpointsExist checks whether any endpoint specified in the command line
// points to a remote endpoint.
func remoteEndpointsExist(spokeEndpointSpecs []*endpointpb.EndpointSpec, hubEndpointSpec *endpointpb.EndpointSpec) bool {
	return isRemoteEndpoint(hubEndpointSpec) || slices.ContainsFunc(spokeEndpointSpecs, isRemoteEndpoint)
}

// makeLineId creates a line identifier based on the organization id and the hub cluster id.
//
// The resulting identifier consists of the organization id and the cluster id, separated by a hyphen.
// All underscores are replaced with hyphens.
//
// Examples:
//
//	makeLineId("my_org", "vmp-123") = "my-org-vmp-123"
func makeLineId(org, cluster string) string {
	return fmt.Sprintf(
		"%s-%s",
		strings.ReplaceAll(org, "_", "-"),
		strings.ReplaceAll(cluster, "_", "-"))
}

// provisionLineLevelRouter provisions a line-level router for the given line.
func (r *HubServiceCreateRunner) provisionLineLevelRouter(ctx context.Context, lineId string) error {
	conn, err := r.dialCloudCluster(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to the cloud: %w", err)
	}
	defer conn.Close()

	client := provisionerpb.NewLineRouterProvisionerClient(conn)
	_, err = client.Provision(ctx, &provisionerpb.ProvisionLineRouterRequest{
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
		ServiceInstallingCmdRunner: ServiceInstallingCmdRunner{
			CmdRunnerBase: *newCmdRunnerBase(
				conn,
				out,
				cluster,
				onpremToLineRouterRelayServicePackage,
				onpremToLineRouterRelayServiceName),
			requestedVersion: versionToInstall,
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
	versionToInstall string,
	hubClusterId string,
	spokeEndpointSpecs []*endpointpb.EndpointSpec) error {

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
		ServiceInstallingCmdRunner: ServiceInstallingCmdRunner{
			CmdRunnerBase: *newCmdRunnerBase(
				conn,
				out,
				hubClusterId,
				hubServicePackage,
				hubServiceName),
			requestedVersion: versionToInstall,
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
	out io.Writer,
	versionToInstall string,
	hubClusterId string,
	spokeEndpointSpecs []*endpointpb.EndpointSpec) error {

	lineId := makeLineId(r.org, hubClusterId)
	fmt.Fprintf(out, "Creating a mixed line orchestration network. Line id: %q\n", lineId)
	fmt.Fprintf(out, "Provisioning a line-level router...\n")
	if err := r.provisionLineLevelRouter(ctx, lineId); err != nil {
		return fmt.Errorf("failed to provision line-level router: %w", err)
	}
	hubAndSpokeClusterIds := make([]string, 0, len(spokeEndpointSpecs)+1)
	hubAndSpokeClusterIds = append(hubAndSpokeClusterIds, hubClusterId) // Hub
	for _, spec := range spokeEndpointSpecs {
		hubAndSpokeClusterIds = append(hubAndSpokeClusterIds, spec.WorkcellName) // Spoke
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
		"All relay routers has been installed, the line orchestration network is ready.\n")
	return nil
}

// run implements high-level logic of the hub-service-create command.
func (r *HubServiceCreateRunner) run(ctx context.Context, out io.Writer) error {
	spokeEndpointSpecs, err := parseEndpointSpecs(r.spokeEndpoints)
	if err != nil {
		return err
	}
	if len(spokeEndpointSpecs) == 0 {
		return fmt.Errorf("at least one spoke endpoint must be specified using --spoke-endpoint")
	}

	hubEndpointSpec, err := parseEndpointSpec(r.hubEndpoint)
	if err != nil {
		return err
	}

	if remoteEndpointsExist(spokeEndpointSpecs, hubEndpointSpec) {
		versionToInstall, err := r.getServiceVersionToInstall(
			onpremToLineRouterRelayServicePackage,
			onpremToLineRouterRelayServiceName)
		if err != nil {
			return err
		}
		return r.createMixedLineOrchestrationNetwork(
			ctx,
			out,
			versionToInstall,
			hubEndpointSpec.WorkcellName,
			spokeEndpointSpecs)
	}

	versionToInstall, err := r.getServiceVersionToInstall(hubServicePackage, hubServiceName)
	if err != nil {
		return err
	}
	return r.createLocalOnlyLineOrchestrationNetwork(
		ctx,
		out,
		versionToInstall,
		hubEndpointSpec.WorkcellName,
		spokeEndpointSpecs)
}
