// Copyright 2023 Intrinsic Innovation LLC

package pubsub

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
)

const (
	hubServicePackage        = "ai.intrinsic"
	hubServiceName           = "line_orchestration_relay"
	defaultHubServiceVersion = "0.0.1"

	forwardingServicePackage        = "ai.intrinsic"
	forwardingServiceName           = "line_orchestration_forwarder"
	defaultForwardingServiceVersion = "0.0.1"

	listAssetInstancesPageSize = 200
)

// CmdRunnerBase is the base class for command runners related to
// line orchestration.
//
// It provides functions for connecting to and calling gRPC services.
// Those functions can be mocked in tests.
type CmdRunnerBase struct {
	outputWriter io.Writer

	clusterId   string
	packageName string
	serviceName string

	installedAssetsClient iagrpcpb.InstalledAssetsClient
	assetInstancesClient  aigrpcpb.AssetInstancesClient
	deploymentClient      adgrpcpb.AssetDeploymentServiceClient
	operationsClient      lrogrpcpb.OperationsClient
}

func newCmdRunnerBase(conn *grpc.ClientConn, outputWriter io.Writer, clusterId, packageName, serviceName string) *CmdRunnerBase {
	return &CmdRunnerBase{
		outputWriter:          outputWriter,
		clusterId:             clusterId,
		packageName:           packageName,
		serviceName:           serviceName,
		installedAssetsClient: iagrpcpb.NewInstalledAssetsClient(conn),
		assetInstancesClient:  aigrpcpb.NewAssetInstancesClient(conn),
		deploymentClient:      adgrpcpb.NewAssetDeploymentServiceClient(conn),
		operationsClient:      lrogrpcpb.NewOperationsClient(conn),
	}
}

// getInstalledServiceAssetVersion returns the version of the service asset
// installed in the current solution.
//
// If the asset is not installed, it returns an empty string and no error.
func (r *CmdRunnerBase) getInstalledServiceAssetVersion(ctx context.Context) (string, error) {
	idProto, err := idutils.IDProtoFrom(r.packageName, r.serviceName)
	if err != nil {
		return "", err
	}
	resp, err := r.installedAssetsClient.GetInstalledAsset(ctx, &iagrpcpb.GetInstalledAssetRequest{
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

func (r *CmdRunnerBase) deleteExistingServiceInstances(ctx context.Context) (int, error) {
	numDeletedInstances := 0
	idProto, err := idutils.IDProtoFrom(r.packageName, r.serviceName)
	if err != nil {
		return 0, err
	}

	for {
		// Requesting one page of asset instances and deleting them.
		// Repeating this process until no more instances are found.
		resp, err := r.assetInstancesClient.ListAssetInstances(ctx, &aigrpcpb.ListAssetInstancesRequest{
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
	fmt.Fprintf(r.outputWriter, "Deleting an instance of the %v service named %q.\n", r.serviceName, instanceName)
	op, err := r.deploymentClient.DeleteResource(ctx, &adpb.DeleteResourceRequest{
		Name: instanceName,
	})
	if err != nil {
		return fmt.Errorf("could not delete instance of the %v service: %w", r.serviceName, err)
	}

	if _, err := waitForOperation(ctx, r.operationsClient, op, r.outputWriter); err != nil {
		return err
	}

	fmt.Fprintf(r.outputWriter, "Successfully deleted an instance of the %v service\n", r.serviceName)
	return nil
}

// uninstallServiceAsset uninstalls a service asset from the current solution.
func (r *CmdRunnerBase) uninstallServiceAsset(ctx context.Context) error {
	idProto, err := idutils.IDProtoFrom(r.packageName, r.serviceName)
	if err != nil {
		return err
	}
	op, err := r.installedAssetsClient.DeleteInstalledAsset(ctx, &iagrpcpb.DeleteInstalledAssetRequest{
		Asset:  idProto,
		Policy: iagrpcpb.DeletePolicy_POLICY_REJECT_USED,
	})
	if err != nil {
		return err
	}

	if _, err := waitForOperation(ctx, r.operationsClient, op, r.outputWriter); err != nil {
		return err
	}

	fmt.Fprintf(r.outputWriter, "Successfully uninstalled the %v service asset\n", r.serviceName)
	return nil
}
