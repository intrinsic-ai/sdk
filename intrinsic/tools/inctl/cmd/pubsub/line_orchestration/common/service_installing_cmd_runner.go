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

	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"

	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	adpb "intrinsic/assets/proto/asset_deployment_go_proto"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	viewpb "intrinsic/assets/proto/view_go_proto"
	lroutils "intrinsic/tools/inctl/cmd/pubsub/long_running_operation_utils"
)

// ServiceInstallingCmdRunner is the base class for command runners
// that can install a service asset.
type ServiceInstallingCmdRunner struct {
	CmdRunnerBase

	RequestedVersion string
}

// installServiceAsset installs a service asset to the current solution.
func (r *ServiceInstallingCmdRunner) installServiceAsset(ctx context.Context) error {
	idVersion, err := idutils.IDVersionProtoFrom(r.PackageName, r.ServiceName, r.RequestedVersion)
	if err != nil {
		return err
	}

	op, err := r.InstalledAssetsClient.CreateInstalledAsset(ctx, &iagrpcpb.CreateInstalledAssetRequest{
		Asset: &iagrpcpb.CreateInstalledAssetRequest_Asset{
			Variant: &iagrpcpb.CreateInstalledAssetRequest_Asset_Catalog{
				Catalog: idVersion,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("could not install %v service asset: %w", r.ServiceName, err)
	}

	if _, err := lroutils.WaitForOperation(ctx, r.OperationsClient, op, r.OutputWriter); err != nil {
		return err
	}

	fmt.Fprintf(r.OutputWriter, "Successfully installed %v service asset.\n", r.ServiceName)
	return nil
}

// addServiceInstance adds an instance of a service to the current solution.
func (r *ServiceInstallingCmdRunner) addServiceInstance(ctx context.Context, wrappedConfig *anypb.Any) error {
	typeIDVersion, err := idutils.IDVersionFrom(r.PackageName, r.ServiceName, r.RequestedVersion)
	if err != nil {
		return fmt.Errorf("failed to create type id version: %w", err)
	}

	op, err := r.DeploymentClient.CreateResourceFromCatalog(ctx, &adpb.CreateResourceFromCatalogRequest{
		TypeIdVersion: typeIDVersion,
		Configuration: &adpb.ResourceInstanceConfiguration{
			Name:          r.ServiceName,
			Configuration: wrappedConfig,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create resource: %w", err)
	}

	if _, err := lroutils.WaitForOperation(ctx, r.OperationsClient, op, r.OutputWriter); err != nil {
		return err
	}

	fmt.Fprintf(r.OutputWriter, "Successfully added an instance of the %v service.\n", r.ServiceName)
	return nil
}

// UpdateInstalledServiceInstances implements the core logic of the command:
//   - Deletes existing instances of the service.
//   - Installs the requested version of the service asset.
//   - Creates a new instance of the service.
func (r *ServiceInstallingCmdRunner) UpdateInstalledServiceInstances(ctx context.Context, config proto.Message) error {
	fmt.Fprintf(
		r.OutputWriter,
		"Deleting existing instances of the %v service in the current solution.\n", r.ServiceName)
	numInstances, err := r.DeleteExistingServiceInstances(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete existing instances of the %v service: %w", r.ServiceName, err)
	}
	fmt.Fprintf(r.OutputWriter, "%v instances have been deleted.\n", numInstances)

	fmt.Fprintf(
		r.OutputWriter,
		"Checking version of the %v service asset installed in the current solution.\n",
		r.ServiceName)
	currentVersion, err := r.GetInstalledServiceAssetVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to determine version of the %v service asset: %w", r.ServiceName, err)
	}

	var shouldInstallAsset bool
	if len(currentVersion) == 0 {
		fmt.Fprintf(
			r.OutputWriter,
			"The %v service asset is currently not installed. Will install.\n",
			r.ServiceName)
		shouldInstallAsset = true
	} else if currentVersion != r.RequestedVersion {
		fmt.Fprintf(
			r.OutputWriter,
			"Current version of the %v service asset is %v, requested version is %v. Will install the requested version.\n",
			r.ServiceName, currentVersion, r.RequestedVersion)
		shouldInstallAsset = true
	} else {
		fmt.Fprintf(
			r.OutputWriter,
			"Current version of the %v service asset is the same as the requested version (%v). Will use currently installed asset.\n",
			r.ServiceName,
			currentVersion)
		shouldInstallAsset = false
	}

	if shouldInstallAsset {
		fmt.Fprintf(
			r.OutputWriter,
			"Installing the %v service asset, version %v.\n",
			r.ServiceName,
			r.RequestedVersion)
		if err = r.installServiceAsset(ctx); err != nil {
			return fmt.Errorf("failed to install %v service asset: %w", r.ServiceName, err)
		}
	}

	fmt.Fprintf(r.OutputWriter, "Adding instance of the %v service.\n", r.ServiceName)
	wrappedConfig, err := anypb.New(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	return r.addServiceInstance(ctx, wrappedConfig)
}

// GetDefaultVersion fetches the default version of the given service from the asset catalog.
func GetDefaultVersion(ctx context.Context, cmdFlags *cmdutils.CmdFlags, outputWriter io.Writer, packageName, serviceName string) (string, error) {
	assetCatalogProject, err := clientutils.GetAssetCatalogProject(cmdFlags.GetFlagProject())
	if err != nil {
		return "", fmt.Errorf(
			"failed to determine asset catalog project for %v: %w",
			cmdFlags.GetFlagProject(), err)
	}
	fmt.Fprintf(outputWriter, "Connecting to the catalog.\n")
	ctx, conn, err := clientutils.DialCatalog(ctx, clientutils.DialCatalogOptions{
		Address: "",
		APIKey: "",
		Org:     cmdFlags.GetFlagOrganization(),
		Project: assetCatalogProject,
	})
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return GetDefaultVersionFromCatalog(ctx, conn, outputWriter, packageName, serviceName)
}

// GetDefaultVersion fetches the default version of the given service from the asset catalog.
func GetDefaultVersionFromCatalog(ctx context.Context, catalogConn *grpc.ClientConn, outputWriter io.Writer, packageName, serviceName string) (string, error) {
	fmt.Fprintf(outputWriter, "Checking which version of %v is the default.\n", serviceName)
	assetIDProto, err := idutils.IDProtoFrom(packageName, serviceName)
	if err != nil {
		return "", fmt.Errorf("failed to generate asset id proto: %w", err)
	}

	client := acgrpcpb.NewAssetCatalogClient(catalogConn)
	resp, err := client.GetAsset(ctx, &acgrpcpb.GetAssetRequest{
		AssetId: &acpb.GetAssetRequest_Id{
			Id: assetIDProto,
		},
		View: viewpb.AssetViewType_ASSET_VIEW_TYPE_VERSIONS,
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch %v from catalog: %w", serviceName, err)
	}

	result := resp.Metadata.IdVersion.Version
	fmt.Fprintf(outputWriter, "The default version is %v.\n", result)
	return result, nil
}
