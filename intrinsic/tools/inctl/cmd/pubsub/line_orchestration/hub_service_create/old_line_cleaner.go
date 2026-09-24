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
	"io"
	"slices"

	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"
	servicedeletionutils "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common/service_deletion_utils"
	"intrinsic/util/go/set"

	endpointpb "intrinsic/platform/pubsub/connect/common/proto/line_configuration/v1/endpoint_spec_go_proto"
	lineconfigpb "intrinsic/platform/pubsub/connect/common/proto/line_configuration/v1/line_configuration_go_proto"

	"google.golang.org/grpc"
)

// lineTypeStr returns a string describing topology of a line orchestration
// network in one word.
func lineTypeStr(remoteEndpointsPresent bool) string {
	if remoteEndpointsPresent {
		return "mixed"
	}
	return "local-only"
}

// makeSetOfClusterIDs creates a set of cluster ids that the line orchestration
// network consists of.
func makeSetOfClusterIDs(hubEndpointSpec *endpointpb.EndpointSpec, spokeEndpointSpecs []*endpointpb.EndpointSpec) set.Set[string] {
	result := set.New[string]()
	result.Add(hubEndpointSpec.GetWorkcellName())
	for _, spokeSpec := range spokeEndpointSpecs {
		result.Add(spokeSpec.GetWorkcellName())
	}
	return result
}

// oldLineCleaner cleans up resources of a line orchestration network.
type oldLineCleaner struct {
	servicedeletionutils.ServiceDeleter

	lineID             string
	hubEndpointSpec    *endpointpb.EndpointSpec
	spokeEndpointSpecs []*endpointpb.EndpointSpec

	oldLineConfig                   *lineconfigpb.LineConfiguration
	remoteEndpointsPresentInOldLine bool
	remoteEndpointsPresentInNewLine bool
}

// cleanupOnTopologyChange cleans up resources of a line orchestration network
// when its topology changes:
//   - If the topology changes from local-only to mixed, then it uninstalls the
//     relay router from the hub cluster (because local-only and mixed networks
//     use different relay routers).
//   - If the topology changes from mixed to local-only, then it tears down the
//     entire mixed network. The line-level router is uninstalled, and the relay
//     service is deleted from all clusters.
func (c *oldLineCleaner) cleanupOnTopologyChange(ctx context.Context, out io.Writer, conn *grpc.ClientConn) error {
	hubClusterID := c.hubEndpointSpec.GetWorkcellName()
	packageName := common.HubServicePackage
	serviceName := common.HubServiceName
	fmt.Fprintf(out,
		"Topology of the line orchestration network is changing from %v to %v.\n",
		lineTypeStr(c.remoteEndpointsPresentInOldLine),
		lineTypeStr(c.remoteEndpointsPresentInNewLine))

	if c.remoteEndpointsPresentInOldLine {
		fmt.Fprintf(out, "The existing network needs to be torn down before the new one can be created.\n")
		fmt.Fprintf(out, "Tearing down the existing network now.\n")
		err := c.DeleteMixedLineOrchestrationNetwork(ctx, out, conn, c.lineID, c.oldLineConfig)
		if err != nil {
			return fmt.Errorf("failed to tear down existing line orchestration network: %w", err)
		}
		fmt.Fprintf(out, "The existing line orchestration network has been torn down.\n")
		return nil
	}

	// If we got here, it means that topology is changing from local-only to mixed.
	// In this case, we need to uninstall the relay service from the hub, because
	// local-only and mixed networks use different relay services.
	fmt.Fprintf(out,
		"Uninstalling %v from %v\n",
		serviceName,
		hubClusterID)

	err := c.DeleteOnpremServiceAndMaybeIgnoreErrors(
		ctx,
		out,
		hubClusterID,
		packageName,
		serviceName)
	if err != nil {
		fmt.Fprintf(out, `
Failed to delete %q from the hub cluster %q.

Cleanup of the existing line orchestration network will be aborted,
and the network will not be recreated.

If the hub is no longer available, try deleting the network by running
the following command:

inctl pubsub hub-service-delete --%v=%v --org=%v --%v
`,
			serviceName,
			hubClusterID,
			common.KeyHubEndpoint,
			common.EndpointSpecString(c.hubEndpointSpec),
			c.OrgID,
			common.KeyIgnoreOnpremErrors)
		return fmt.Errorf("failed to delete %q from the hub cluster %q: %w", serviceName, hubClusterID, err)
	}

	fmt.Fprintf(out, "Cleanup of the existing network complete.\n")
	return nil
}

// cleanupMixedNetwork cleans up resources of a mixed line orchestration network.
//
// More specifically, it uninstalls relay routers from clusters that are removed
// from the network.
func (c *oldLineCleaner) cleanupMixedNetwork(ctx context.Context, out io.Writer) error {
	oldClusters := makeSetOfClusterIDs(c.oldLineConfig.GetHubEndpoint(), c.oldLineConfig.GetSpokeEndpoints())
	newClusters := makeSetOfClusterIDs(c.hubEndpointSpec, c.spokeEndpointSpecs)
	serviceName := common.OnpremToLineRouterRelayServiceName

	for _, removedClusterID := range slices.Sorted(oldClusters.Difference(newClusters)) {
		fmt.Fprintf(
			out,
			"Cluster %q is being removed from the network. Uninstalling %v from it.\n",
			removedClusterID,
			serviceName)
		// If something goes wrong on the spoke that is being removed from the network,
		// it's not necessarily a fatal error. It's possible that the spoke VM is
		// no longer alive (its lease has expired and it was returned to the pool).
		// In this case, we should allow the user to ignore errors on spokes.
		err := c.DeleteOnpremServiceAndMaybeIgnoreErrors(
			ctx,
			out,
			removedClusterID,
			common.OnpremToLineRouterRelayServicePackage,
			serviceName)
		if err != nil {
			return fmt.Errorf("failed to delete %q from %q: %w", serviceName, removedClusterID, err)
		}
	}

	fmt.Fprintf(out, "Cleanup of the existing network complete.\n")
	return nil
}

// cleanup decides which cleanup actions to perform, and initiates cleanup.
func (c *oldLineCleaner) cleanup(ctx context.Context, out io.Writer, conn *grpc.ClientConn) error {
	if c.remoteEndpointsPresentInOldLine != c.remoteEndpointsPresentInNewLine {
		return c.cleanupOnTopologyChange(ctx, out, conn)
	}

	if c.remoteEndpointsPresentInNewLine {
		// Network topology stays mixed, so we only need to uninstall the relay service
		// from clusters that are being removed from the network.
		return c.cleanupMixedNetwork(ctx, out)
	}

	fmt.Fprintf(out, "Network topology remains local-only, nothing to clean up.\n")
	return nil
}
