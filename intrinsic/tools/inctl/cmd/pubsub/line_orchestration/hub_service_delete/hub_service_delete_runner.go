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
	"errors"
	"fmt"
	"io"

	"intrinsic/platform/pubsub/connect/common/lineidutils"
	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"
	servicedeletionutils "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/service_deletion_utils"

	provisionerpb "intrinsic/platform/pubsub/cloud_router_provisioner/v1/provisioner_go_proto"
	lineconfigstoragepb "intrinsic/platform/pubsub/connect/cloud/proto/line_configuration_storage/v1/line_configuration_storage_go_proto"
	lineconfigpb "intrinsic/platform/pubsub/connect/common/proto/line_configuration/v1/line_configuration_go_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HubServiceDeleteRunner implements high-level logic of the hub-service-delete command.
//
// Its primary responsibility is to determine the topology of the network (local-only or mixed),
// and delete all constituent parts of that network.
type HubServiceDeleteRunner struct {
	project            string
	org                string
	hubEndpoint        string
	forceLocalOnly     bool
	ignoreOnpremErrors bool

	shouldUninstallServiceAsset bool

	// dialOnpremCluster creates a connection to the given spoke cluster.
	// The default implementation connects to the real workcell or VM.
	// Unit tests can provide an implementation that connects to a fake server.
	dialOnpremCluster func(context.Context, string, string, string) (context.Context, *grpc.ClientConn, string, error)

	// dialCloudCluster connects to the cloud-robotics cluster.
	// The default implementation connects to the real cloud-robotics cluster.
	// Unit tests can provide an implementation that connects to a fake server.
	dialCloudCluster func(context.Context) (*grpc.ClientConn, error)
}

func (r *HubServiceDeleteRunner) fetchLineConfiguration(
	ctx context.Context,
	conn *grpc.ClientConn,
	lineId string) (*lineconfigpb.LineConfiguration, error) {
	client := lineconfigstoragepb.NewLineConfigurationStorageServiceClient(conn)
	req := &lineconfigstoragepb.GetLineConfigurationRequest{
		Name: lineidutils.ConvertIDToName(lineId),
	}
	lineConfig, err := client.GetLineConfiguration(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, common.ErrLineConfigNotFound
			case codes.Unimplemented:
				return nil, common.ErrConfigStorageNotImplemented
			}
		}

		return nil, err
	}

	return lineConfig, nil
}

func (r *HubServiceDeleteRunner) uninstallLineLevelCloudRouter(
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

func (r *HubServiceDeleteRunner) deleteLineConfig(
	ctx context.Context,
	conn *grpc.ClientConn,
	lineId string) error {
	client := lineconfigstoragepb.NewLineConfigurationStorageServiceClient(conn)
	_, err := client.DeleteLineConfiguration(ctx, &lineconfigstoragepb.DeleteLineConfigurationRequest{
		Name: lineidutils.ConvertIDToName(lineId),
	})
	if err != nil {
		return fmt.Errorf("failed to delete configuration for line %q: %w", lineId, err)
	}

	return nil
}

func (r *HubServiceDeleteRunner) deleteOnpremService(
	ctx context.Context,
	out io.Writer,
	clusterID, packageName, serviceName string) error {
	fmt.Fprintf(out, "Connecting to %v\n", clusterID)
	ctx, conn, _, err := r.dialOnpremCluster(ctx, r.project, r.org, clusterID)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Fprintf(out, "Uninstalling %v.%v from %v\n", packageName, serviceName, clusterID)
	runner := &servicedeletionutils.ServiceDeleteCmdRunner{
		CmdRunnerBase: *common.NewCmdRunnerBase(
			conn,
			out,
			clusterID,
			packageName,
			serviceName),
		ShouldUninstallServiceAsset: r.shouldUninstallServiceAsset,
	}

	return runner.Run(ctx)
}

func (r *HubServiceDeleteRunner) deleteOnpremServiceAndMaybeIgnoreErrors(
	ctx context.Context,
	out io.Writer,
	clusterID, packageName, serviceName string) error {
	err := r.deleteOnpremService(ctx, out, clusterID, packageName, serviceName)
	if err != nil && r.ignoreOnpremErrors {
		fmt.Fprintf(
			out,
			"Failed to delete %q from %q, but the '--%v' flag is present. Moving on.\n",
			serviceName,
			clusterID,
			servicedeletionutils.KeyIgnoreOnpremErrors)
		return nil
	}

	return err
}

func (r *HubServiceDeleteRunner) deleteMixedLineOrchestrationNetwork(
	ctx context.Context,
	cloudConn *grpc.ClientConn,
	out io.Writer,
	lineId string,
	lineConfig *lineconfigpb.LineConfiguration) error {

	clusterIDs := make([]string, 0, len(lineConfig.GetSpokeEndpoints())+1)
	clusterIDs = append(clusterIDs, lineConfig.GetHubEndpoint().GetWorkcellName())
	for _, spokeEndpoint := range lineConfig.GetSpokeEndpoints() {
		clusterIDs = append(clusterIDs, spokeEndpoint.GetWorkcellName())
	}
	for _, clusterID := range clusterIDs {
		err := r.deleteOnpremServiceAndMaybeIgnoreErrors(
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
	if err := r.uninstallLineLevelCloudRouter(ctx, cloudConn, lineId); err != nil {
		return err
	}

	return nil
}

func (r *HubServiceDeleteRunner) run(ctx context.Context, out io.Writer) error {
	hubEndpointSpec, err := common.ParseEndpointSpec(r.hubEndpoint)
	if err != nil {
		return err
	}

	if r.forceLocalOnly {
		if common.IsRemoteEndpoint(hubEndpointSpec) {
			return fmt.Errorf(
				"The '--%v' flag is present, but the hub endpoint is remote.",
				common.KeyForceLocalOnly)
		}
		fmt.Fprintf(
			out,
			`
The '--%v' flag is present.
Deleting the line orchestration network assuming that it is local-only.
`,
			common.KeyForceLocalOnly)
		return r.deleteOnpremServiceAndMaybeIgnoreErrors(
			ctx,
			out,
			hubEndpointSpec.WorkcellName,
			common.HubServicePackage,
			common.HubServiceName)
	}

	fmt.Fprintf(out, "Connecting to the cloud API endpoint\n")
	conn, err := r.dialCloudCluster(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to the cloud: %w", err)
	}
	defer conn.Close()

	lineID, err := common.MakeLineID(r.org, hubEndpointSpec.WorkcellName)
	if err != nil {
		fmt.Fprintf(out, "Cannot create a valid line id from %q and %q\n", r.org, hubEndpointSpec.WorkcellName)
		return err
	}

	lineConfig, err := r.fetchLineConfiguration(ctx, conn, lineID)
	if err != nil {
		if errors.Is(err, common.ErrLineConfigNotFound) {
			fmt.Fprintf(
				out,
				`
------------------------------------------------------------------
Configuration for line %q not found.

It is possible that a network for this line had been created
before we started saving line configurations in the cloud.

If that's the case, then that network is local-only.
To delete it, try running this command with the '--%v' flag.
------------------------------------------------------------------
`,
				lineID,
				common.KeyForceLocalOnly,
			)
			return err
		} else if errors.Is(err, common.ErrConfigStorageNotImplemented) {
			common.ExplainConfigStorageNotImplementedError(out, "hub-service-delete", common.KeyForceLocalOnly, "delete")
			return err
		} else {
			return fmt.Errorf("failed to fetch line configuration: %w", err)
		}
	}

	if common.RemoteEndpointsExist(lineConfig.GetSpokeEndpoints(), lineConfig.GetHubEndpoint()) {
		fmt.Fprintf(
			out,
			"Line %q has been identified as mixed (containing local and remote clusters, e.g. VMs).\n",
			lineID)
		if err = r.deleteMixedLineOrchestrationNetwork(ctx, conn, out, lineID, lineConfig); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(out, "Line %q has been identified as local-only.\n", lineID)
		err = r.deleteOnpremServiceAndMaybeIgnoreErrors(
			ctx,
			out,
			hubEndpointSpec.WorkcellName,
			common.HubServicePackage,
			common.HubServiceName)
		if err != nil {
			return fmt.Errorf(
				"failed to delete the relay router from %v: %v\n",
				hubEndpointSpec.WorkcellName,
				err)
		}
	}

	fmt.Fprintf(out, "Deleting line configuration stored in the cloud.\n")
	if err := r.deleteLineConfig(ctx, conn, lineID); err != nil {
		return err
	}

	fmt.Fprintf(
		out,
		"Line orchestration network with the hub at %q (line id %q) has been deleted.\n",
		hubEndpointSpec.WorkcellName,
		lineID)
	return nil
}
