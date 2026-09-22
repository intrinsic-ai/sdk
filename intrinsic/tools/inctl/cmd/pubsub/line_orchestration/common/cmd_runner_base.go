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

package common

import (
	"context"
	"fmt"
	"io"

	lrogrpcpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	"intrinsic/assets/idutils"
	adgrpcpb "intrinsic/assets/proto/asset_deployment_go_proto"
	adpb "intrinsic/assets/proto/asset_deployment_go_proto"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	aigrpcpb "intrinsic/assets/proto/v1/asset_instances_go_proto"
	viewpb "intrinsic/assets/proto/view_go_proto"
	lroutils "intrinsic/tools/inctl/cmd/pubsub/long_running_operation_utils"
)

const (
	HubServicePackage = "ai.intrinsic"
	HubServiceName    = "line_orchestration_relay"

	OnpremToLineRouterRelayServicePackage = "ai.intrinsic"
	OnpremToLineRouterRelayServiceName    = "onprem_to_line_router_relay"

	ForwardingServicePackage = "ai.intrinsic"
	ForwardingServiceName    = "line_orchestration_forwarder"

	listAssetInstancesPageSize = 200
)

// CmdRunnerBase is the base class for command runners related to
// line orchestration.
//
// It provides functions for connecting to and calling gRPC services.
// Those functions can be mocked in tests.
type CmdRunnerBase struct {
	OutputWriter io.Writer

	ClusterId   string
	PackageName string
	ServiceName string

	InstalledAssetsClient iagrpcpb.InstalledAssetsClient
	AssetInstancesClient  aigrpcpb.AssetInstancesClient
	DeploymentClient      adgrpcpb.AssetDeploymentServiceClient
	OperationsClient      lrogrpcpb.OperationsClient
}

func NewCmdRunnerBase(conn *grpc.ClientConn, outputWriter io.Writer, clusterId, packageName, serviceName string) *CmdRunnerBase {
	return &CmdRunnerBase{
		OutputWriter:          outputWriter,
		ClusterId:             clusterId,
		PackageName:           packageName,
		ServiceName:           serviceName,
		InstalledAssetsClient: iagrpcpb.NewInstalledAssetsClient(conn),
		AssetInstancesClient:  aigrpcpb.NewAssetInstancesClient(conn),
		DeploymentClient:      adgrpcpb.NewAssetDeploymentServiceClient(conn),
		OperationsClient:      lrogrpcpb.NewOperationsClient(conn),
	}
}

// GetInstalledServiceAssetVersion returns the version of the service asset
// installed in the current solution.
//
// If the asset is not installed, it returns an empty string and no error.
func (r *CmdRunnerBase) GetInstalledServiceAssetVersion(ctx context.Context) (string, error) {
	idProto, err := idutils.IDProtoFrom(r.PackageName, r.ServiceName)
	if err != nil {
		return "", err
	}
	resp, err := r.InstalledAssetsClient.GetInstalledAsset(ctx, &iagrpcpb.GetInstalledAssetRequest{
		Id:   idProto,
		View: viewpb.AssetViewType_ASSET_VIEW_TYPE_BASIC,
	})
	if err != nil {
		if grpcstatus.Code(err) == codes.NotFound {
			return "", nil
		} else {
			return "", fmt.Errorf("failed to get installed asset: %w", err)
		}
	}

	return resp.Metadata.IdVersion.Version, nil
}

func (r *CmdRunnerBase) DeleteExistingServiceInstances(ctx context.Context) (int, error) {
	numDeletedInstances := 0
	idProto, err := idutils.IDProtoFrom(r.PackageName, r.ServiceName)
	if err != nil {
		return 0, err
	}

	for {
		// Requesting one page of asset instances and deleting them.
		// Repeating this process until no more instances are found.
		resp, err := r.AssetInstancesClient.ListAssetInstances(ctx, &aigrpcpb.ListAssetInstancesRequest{
			PageSize: listAssetInstancesPageSize,
			StrictFilters: []*aigrpcpb.ListAssetInstancesRequest_Filter{
				{
					Id: idProto,
				},
			},
		})
		if err != nil {
			return numDeletedInstances, fmt.Errorf("failed to list asset instances: %w", err)
		}

		for _, instance := range resp.AssetInstances {
			if err := r.deleteServiceInstance(ctx, instance.Name); err != nil {
				return numDeletedInstances, err
			}
			numDeletedInstances++
		}

		if len(resp.AssetInstances) == 0 || resp.NextPageToken == "" {
			break
		}
	}

	return numDeletedInstances, nil
}

// deleteServiceInstance deletes an instance of a service.
func (r *CmdRunnerBase) deleteServiceInstance(ctx context.Context, instanceName string) error {
	fmt.Fprintf(r.OutputWriter, "Deleting an instance of the %v service named %q.\n", r.ServiceName, instanceName)
	op, err := r.DeploymentClient.DeleteResource(ctx, &adpb.DeleteResourceRequest{
		Name: instanceName,
	})
	if err != nil {
		return fmt.Errorf("could not delete instance of the %v service: %w", r.ServiceName, err)
	}

	if _, err := lroutils.WaitForOperation(ctx, r.OperationsClient, op, r.OutputWriter); err != nil {
		return err
	}

	fmt.Fprintf(r.OutputWriter, "Successfully deleted an instance of the %v service\n", r.ServiceName)
	return nil
}

// UninstallServiceAsset uninstalls a service asset from the current solution.
func (r *CmdRunnerBase) UninstallServiceAsset(ctx context.Context) error {
	idProto, err := idutils.IDProtoFrom(r.PackageName, r.ServiceName)
	if err != nil {
		return err
	}
	op, err := r.InstalledAssetsClient.DeleteInstalledAsset(ctx, &iagrpcpb.DeleteInstalledAssetRequest{
		Asset:  idProto,
		Policy: iagrpcpb.DeletePolicy_POLICY_REJECT_USED,
	})
	if err != nil {
		return err
	}

	if _, err := lroutils.WaitForOperation(ctx, r.OperationsClient, op, r.OutputWriter); err != nil {
		return err
	}

	fmt.Fprintf(r.OutputWriter, "Successfully uninstalled the %v service asset\n", r.ServiceName)
	return nil
}
