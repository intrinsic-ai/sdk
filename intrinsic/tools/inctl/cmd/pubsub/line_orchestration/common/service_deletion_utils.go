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

package servicedeletionutils

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"

	provisionerpb "intrinsic/platform/pubsub/cloud_router_provisioner/v1/provisioner_go_proto"
	lineconfigpb "intrinsic/platform/pubsub/connect/common/proto/line_configuration/v1/line_configuration_go_proto"
)

// ServiceDeleteCmdRunner handles execution of commands that delete service assets
// used for line orchestration.
type ServiceDeleteCmdRunner struct {
	common.CmdRunnerBase

	// ShouldRetainServiceAsset specifies whether a service asset should be
	// kept in a solution.
	//
	// When it's false, service instances are removed, and the service asset is
	// uninstalled. When it's true, service instances are removed, but the
	// service asset remains installed.
	ShouldRetainServiceAsset bool
}

// Run implements the core logic of the command:
//   - Deletes all instances of the service.
//   - Optionally uninstalls the service asset.
func (r *ServiceDeleteCmdRunner) Run(ctx context.Context) error {
	fmt.Fprintf(
		r.OutputWriter,
		"Deleting existing instances of the %v service in the current solution.\n",
		r.ServiceName)
	numInstances, err := r.DeleteExistingServiceInstances(ctx)
	if err != nil {
		return fmt.Errorf(
			"failed to delete existing instances of the %v service: %w",
			r.ServiceName, err)
	}
	fmt.Fprintf(r.OutputWriter, "%v instances have been deleted.\n", numInstances)

	if r.ShouldRetainServiceAsset {
		fmt.Fprintf(
			r.OutputWriter,
			"The %v option is enabled, won't try to uninstall the %v service asset.\n",
			common.KeyRetainServiceAsset,
			r.ServiceName)
		return nil
	}

	if numInstances == 0 {
		// There were no instances of the service, but it is still possible that
		// the service asset is installed. Checking it here.
		fmt.Fprintf(
			r.OutputWriter,
			"Checking if the %v service asset is installed.\n",
			r.ServiceName)
		currentVersion, err := r.GetInstalledServiceAssetVersion(ctx)
		if err != nil {
			return fmt.Errorf(
				"failed to determine whether the %v service is installed: %w",
				r.ServiceName,
				err)
		}
		if len(currentVersion) == 0 {
			fmt.Fprintf(
				r.OutputWriter,
				"%v service asset is not installed, nothing else to do.\n",
				r.ServiceName)
			return nil
		}
		fmt.Fprintf(
			r.OutputWriter,
			"%v service, version %v, is installed. Will uninstall.\n", r.ServiceName, currentVersion)
	}

	fmt.Fprintf(r.OutputWriter, "Uninstalling %v service asset.\n", r.ServiceName)
	if err = r.UninstallServiceAsset(ctx); err != nil {
		return fmt.Errorf("failed to uninstall %v service asset: %w", r.ServiceName, err)
	}
	fmt.Fprintf(
		r.OutputWriter,
		"The %v service asset has been uninstalled.\n",
		r.ServiceName)

	return nil
}

// ServiceDeleter deletes a service installed in a solution.
type ServiceDeleter struct {
	// Cloud project that the cluster belongs to.
	ProjectID string

	// Organization that the cluster belongs to.
	OrgID string

	// Whether the service asset should be retained.
	//
	// When false, service instances are deleted, and the service
	// asset is uninstalled.
	// When true, service instances are deleted, but the service
	// asset remains installed.
	ShouldRetainServiceAsset bool

	// Whether failures should be ignored.
	//
	// It may be necessary to ignore them if the cluster is
	// no longer available (e.g. VM lease expired), and a
	// failure to delete a service from that cluster should
	// not block a bigger workflow (e.g. deletion of a line
	// orchestration network).
	IgnoreOnpremErrors bool

	// DialOnpremCluster creates a connection to the given spoke cluster.
	// The default implementation connects to the real workcell or VM.
	// Unit tests can provide an implementation that connects to a fake server.
	DialOnpremCluster func(ctx context.Context, projectID string, orgID string, clusterID string) (context.Context, *grpc.ClientConn, string, error)
}

// deleteOnpremService deletes the given service from the given cluster.
//
// Returns an error if deletion fails.
func (d *ServiceDeleter) deleteOnpremService(
	ctx context.Context,
	out io.Writer,
	clusterID, packageName, serviceName string) error {
	fmt.Fprintf(out, "Connecting to %v\n", clusterID)
	ctx, conn, _, err := d.DialOnpremCluster(ctx, d.ProjectID, d.OrgID, clusterID)
	if err != nil {
		return fmt.Errorf("failed to connect to cluster %q: %w", clusterID, err)
	}
	defer conn.Close()

	fmt.Fprintf(out, "Uninstalling %v.%v from %v\n", packageName, serviceName, clusterID)
	runner := &ServiceDeleteCmdRunner{
		CmdRunnerBase: *common.NewCmdRunnerBase(
			conn,
			out,
			clusterID,
			packageName,
			serviceName),
		ShouldRetainServiceAsset: d.ShouldRetainServiceAsset,
	}

	return runner.Run(ctx)
}

// DeleteOnpremServiceAndMaybeIgnoreErrors deletes the given service from the given cluster.
//
// Returns an error if deletion fails AND IgnoreOnpremErrors is false.
// When IgnoreOnpremErrors is true, always returns nil.
func (d *ServiceDeleter) DeleteOnpremServiceAndMaybeIgnoreErrors(
	ctx context.Context,
	out io.Writer,
	clusterID, packageName, serviceName string) error {
	err := d.deleteOnpremService(ctx, out, clusterID, packageName, serviceName)
	if err != nil && d.IgnoreOnpremErrors {
		fmt.Fprintf(
			out,
			"Failed to delete %q from %q: %v, but the '--%v' flag is present. Moving on.\n",
			serviceName,
			clusterID,
			err,
			common.KeyIgnoreOnpremErrors)
		return nil
	}

	return err
}

// DeleteMixedLineOrchestrationNetwork tears down a mixed line orchestration network:
//   - Uninstall relay routers from all onprem clusters.
//   - Uninstalls the cloud router.
func (d *ServiceDeleter) DeleteMixedLineOrchestrationNetwork(
	ctx context.Context,
	out io.Writer,
	cloudConn *grpc.ClientConn,
	lineId string,
	lineConfig *lineconfigpb.LineConfiguration) error {

	clusterIDs := make([]string, 0, len(lineConfig.GetSpokeEndpoints())+1)
	clusterIDs = append(clusterIDs, lineConfig.GetHubEndpoint().GetWorkcellName())
	for _, spokeEndpoint := range lineConfig.GetSpokeEndpoints() {
		clusterIDs = append(clusterIDs, spokeEndpoint.GetWorkcellName())
	}
	for _, clusterID := range clusterIDs {
		err := d.DeleteOnpremServiceAndMaybeIgnoreErrors(
			ctx,
			out,
			clusterID,
			common.OnpremToLineRouterRelayServicePackage,
			common.OnpremToLineRouterRelayServiceName)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(
		out,
		"Line orchestration service assets have been uninstalled from all clusters in line %q.\n",
		lineId)

	fmt.Fprintf(out, "Uninstalling cloud router.\n")
	if err := d.uninstallLineLevelCloudRouter(ctx, cloudConn, lineId); err != nil {
		return err
	}

	return nil
}

func (d *ServiceDeleter) uninstallLineLevelCloudRouter(
	ctx context.Context,
	conn *grpc.ClientConn,
	lineId string) error {
	client := provisionerpb.NewLineRouterProvisionerClient(conn)
	_, err := client.Delete(ctx, &provisionerpb.DeleteLineRouterRequest{
		LineId: lineId,
	})
	if err != nil {
		return fmt.Errorf("failed to uninstall cloud router for %q: %w", lineId, err)
	}

	return nil
}
