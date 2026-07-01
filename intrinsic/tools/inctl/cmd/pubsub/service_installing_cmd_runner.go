// Copyright 2023 Intrinsic Innovation LLC

package pubsub

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"

	"intrinsic/assets/idutils"
	adpb "intrinsic/assets/proto/asset_deployment_go_proto"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
)

// ServiceInstallingCmdRunner is the base class for command runners
// that can install a service asset.
type ServiceInstallingCmdRunner struct {
	CmdRunnerBase

	requestedVersion string
}

// installServiceAsset installs a service asset to the current solution.
func (r *ServiceInstallingCmdRunner) installServiceAsset(ctx context.Context) error {
	idVersion, err := idutils.IDVersionProtoFrom(r.packageName, r.serviceName, r.requestedVersion)
	if err != nil {
		return err
	}

	op, err := r.installedAssetsClient.CreateInstalledAsset(ctx, &iagrpcpb.CreateInstalledAssetRequest{
		Asset: &iagrpcpb.CreateInstalledAssetRequest_Asset{
			Variant: &iagrpcpb.CreateInstalledAssetRequest_Asset_Catalog{
				Catalog: idVersion,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("could not install %v service asset: %w", r.serviceName, err)
	}

	if _, err := waitForOperation(ctx, r.operationsClient, op, r.outputWriter); err != nil {
		return err
	}

	fmt.Fprintf(r.outputWriter, "Successfully installed %v service asset.\n", r.serviceName)
	return nil
}

// addServiceInstance adds an instance of a service to the current solution.
func (r *ServiceInstallingCmdRunner) addServiceInstance(ctx context.Context, wrappedConfig *anypb.Any) error {
	typeIDVersion, err := idutils.IDVersionFrom(r.packageName, r.serviceName, r.requestedVersion)
	if err != nil {
		return fmt.Errorf("failed to create type id version: %w", err)
	}

	op, err := r.deploymentClient.CreateResourceFromCatalog(ctx, &adpb.CreateResourceFromCatalogRequest{
		TypeIdVersion: typeIDVersion,
		Configuration: &adpb.ResourceInstanceConfiguration{
			Name:          r.serviceName,
			Configuration: wrappedConfig,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create resource: %w", err)
	}

	if _, err := waitForOperation(ctx, r.operationsClient, op, r.outputWriter); err != nil {
		return err
	}

	fmt.Fprintf(r.outputWriter, "Successfully added an instance of the %v service.\n", r.serviceName)
	return nil
}

// updateInstalledServiceInstances implements the core logic of the command:
//   - Deletes existing instances of the service.
//   - Installs the requested version of the service asset.
//   - Creates a new instance of the service.
func (r *ServiceInstallingCmdRunner) updateInstalledServiceInstances(ctx context.Context, config proto.Message) error {
	fmt.Fprintf(
		r.outputWriter,
		"Deleting existing instances of the %v service in the current solution.\n", r.serviceName)
	numInstances, err := r.deleteExistingServiceInstances(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete existing instances of the %v service: %w", r.serviceName, err)
	}
	fmt.Fprintf(r.outputWriter, "%v instances have been deleted.\n", numInstances)

	fmt.Fprintf(
		r.outputWriter,
		"Checking version of the %v service asset installed in the current solution.\n",
		r.serviceName)
	currentVersion, err := r.getInstalledServiceAssetVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to determine version of the %v service asset: %w", r.serviceName, err)
	}

	var shouldInstallAsset bool
	if len(currentVersion) == 0 {
		fmt.Fprintf(
			r.outputWriter,
			"The %v service asset is currently not installed. Will install.\n",
			r.serviceName)
		shouldInstallAsset = true
	} else if currentVersion != r.requestedVersion {
		fmt.Fprintf(
			r.outputWriter,
			"Current version of the %v service asset is %v, requested version is %v. Will install the requested version.\n",
			r.serviceName, currentVersion, r.requestedVersion)
		shouldInstallAsset = true
	} else {
		fmt.Fprintf(
			r.outputWriter,
			"Current version of the %v service asset is the same as the requested version (%v). Will use currently installed asset.\n",
			r.serviceName,
			currentVersion)
		shouldInstallAsset = false
	}

	if shouldInstallAsset {
		fmt.Fprintf(
			r.outputWriter,
			"Installing the %v service asset, version %v.\n",
			r.serviceName,
			r.requestedVersion)
		if err = r.installServiceAsset(ctx); err != nil {
			return fmt.Errorf("failed to install %v service asset: %w", r.serviceName, err)
		}
	}

	fmt.Fprintf(r.outputWriter, "Adding instance of the %v service.\n", r.serviceName)
	wrappedConfig, err := anypb.New(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	return r.addServiceInstance(ctx, wrappedConfig)
}
