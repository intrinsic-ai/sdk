// Copyright 2023 Intrinsic Innovation LLC

package pubsub

import (
	"context"
	"fmt"
	"strings"

	anypb "google.golang.org/protobuf/types/known/anypb"

	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	adpb "intrinsic/assets/proto/asset_deployment_go_proto"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	endpointpb "intrinsic/platform/pubsub/connect/cloud/proto/v1alpha1/endpoint_spec_go_proto"
	pb "intrinsic/platform/pubsub/connect/onprem/relay_router_service/relay_router_service_go_proto"

	"github.com/spf13/cobra"
)

const (
	keySpokeEndpoints    = "spoke-endpoint"
	keyHubServiceVersion = "hub-service-version"

	endpointSpecSeparator     = "@"
	localEndpointDesignation  = "local"
	remoteEndpointDesignation = "remote"
)

// HubServiceCreateCmdRunner handles execution of the hub-service-create command.
// That subcommand installs or updates the relay service used for line orchestration.
type HubServiceCreateCmdRunner struct {
	HubServiceCmdRunnerBase

	spokeEndpoints   []string
	requestedVersion string
}

// parseEndpointSpec creates an EndpointSpec based on a spoke-endpoint command line flag.
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

// makeConfig generates configuration of the relay service from command line flags.
func (r *HubServiceCreateCmdRunner) makeConfig() (*pb.RelayRouterServiceConfig, error) {
	spokeWorkcells := r.spokeEndpoints
	config := &pb.RelayRouterServiceConfig{
		HubWorkcellName: r.clusterId,
		SpokeEndpointSpecs: make(
			[]*endpointpb.EndpointSpec,
			len(spokeWorkcells)),
	}

	for i, spokeWorkcell := range spokeWorkcells {
		endpointSpec, err := parseEndpointSpec(spokeWorkcell)
		if err != nil {
			return nil, err
		}

		if endpointSpec.GetRemote() != nil {
			return nil, fmt.Errorf("remote endpoints are not supported")
		}

		config.SpokeEndpointSpecs[i] = endpointSpec
	}

	fmt.Fprintf(r.outputWriter, "--- Config ---\n%v\n--- End of config ---\n", config)

	return config, nil
}

// installRelayServiceAsset installs the relay service asset to the current solution.
func (r *HubServiceCreateCmdRunner) installRelayServiceAsset(ctx context.Context) error {
	idVersion, err := idutils.IDVersionProtoFrom(hubServicePackage, hubServiceName, r.requestedVersion)
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
		return fmt.Errorf("could not install relay service asset: %w", err)
	}

	if _, err := waitForOperation(ctx, r.operationsClient, op, r.outputWriter); err != nil {
		return err
	}

	fmt.Fprintf(r.outputWriter, "Successfully installed line orchestration relay service asset.\n")
	return nil
}

// addRelayServiceInstance adds an instance of the relay service to the current solution.
func (r *HubServiceCreateCmdRunner) addRelayServiceInstance(ctx context.Context) error {
	config, err := r.makeConfig()
	if err != nil {
		return fmt.Errorf("failed to create service config from command line flags: %w", err)
	}
	wrappedConfig, err := anypb.New(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	typeIdVersion, err := idutils.IDVersionFrom(hubServicePackage, hubServiceName, r.requestedVersion)
	if err != nil {
		return fmt.Errorf("failed to create type id version: %w", err)
	}

	op, err := r.deploymentClient.CreateResourceFromCatalog(ctx, &adpb.CreateResourceFromCatalogRequest{
		TypeIdVersion: typeIdVersion,
		Configuration: &adpb.ResourceInstanceConfiguration{
			Name:          hubServiceName,
			Configuration: wrappedConfig,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create resource: %v", err)
	}

	if _, err := waitForOperation(ctx, r.operationsClient, op, r.outputWriter); err != nil {
		return err
	}

	fmt.Fprintf(r.outputWriter, "Successfully added an instance of the line orchestration relay service.\n")
	return nil
}

// run implements the core logic of the command:
//   - Deletes existing instances of the relay service.
//   - Installs the requested version of the service asset.
//   - Creates a new instance of the relay service.
func (r *HubServiceCreateCmdRunner) run(ctx context.Context) error {
	if len(r.spokeEndpoints) == 0 {
		return fmt.Errorf("at least one spoke endpoint must be specified using --spoke-endpoint; cannot install or update the relay.")
	}

	fmt.Fprintf(
		r.outputWriter,
		"Deleting existing instances of the relay service in the current solution.\n")
	numInstances, err := r.deleteExistingRelayServiceInstances(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete existing instances of the relay service: %w", err)
	}
	fmt.Fprintf(r.outputWriter, "%v instances have been deleted.\n", numInstances)

	fmt.Fprintf(r.outputWriter, "Checking version of the relay service asset installed in the current solution.\n")
	requestedVersion := r.requestedVersion
	currentVersion, err := r.getInstalledRelayServiceAssetVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to determine version of the relay service asset: %w", err)
	}
	var shouldInstallAsset bool
	if len(currentVersion) == 0 {
		fmt.Fprintf(r.outputWriter, "The relay service asset is currently not installed. Will install.\n")
		shouldInstallAsset = true
	} else if currentVersion != requestedVersion {
		fmt.Fprintf(
			r.outputWriter,
			"Current version of the relay service asset is %v, requested version is %v. Will install the requested version.\n",
			currentVersion, requestedVersion)
		shouldInstallAsset = true
	} else {
		fmt.Fprintf(
			r.outputWriter,
			"Current version of the relay service asset is the same as the requested version (%v). Will use currently installed asset.\n",
			currentVersion)
		shouldInstallAsset = false
	}

	if shouldInstallAsset {
		fmt.Fprintf(
			r.outputWriter,
			"Installing the relay service asset, version %v.\n",
			requestedVersion)
		if err = r.installRelayServiceAsset(ctx); err != nil {
			return fmt.Errorf("failed to install relay service asset: %w", err)
		}
	}

	fmt.Fprintf(r.outputWriter, "Adding instance of the relay service.\n")
	return r.addRelayServiceInstance(ctx)
}

// hubServiceCreateCmdEnvironment is the execution environment for the
// hub-service-create command. That environment contains command line
// flags and a connection to the gRPC service.
type hubServiceCreateCmdEnvironment struct {
	cmdFlags *cmdutils.CmdFlags
}

// RunE sets up the execution environment and invokes HubServiceCreateCmdRunner.run.
func (e *hubServiceCreateCmdEnvironment) RunE(cmd *cobra.Command, _ []string) error {
	ctx, conn, _, err := clientutils.DialClusterFromInctl(cmd.Context(), e.cmdFlags)
	if err != nil {
		return err
	}
	defer conn.Close()

	runner := &HubServiceCreateCmdRunner{
		HubServiceCmdRunnerBase: *newHubServiceCmdRunnerBase(conn, cmd.OutOrStdout(), e.cmdFlags.GetString(cmdutils.KeyCluster)),
		requestedVersion:        e.cmdFlags.GetString(keyHubServiceVersion),
		spokeEndpoints:          e.cmdFlags.GetStringSlice(keySpokeEndpoints),
	}
	return runner.run(ctx)
}

// NewHubServiceCreateCmd returns the initialized cobra command for hub-service-create.
func NewHubServiceCreateCmd() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	commandWrapper := &hubServiceCreateCmdEnvironment{cmdFlags: flags}

	cmd := &cobra.Command{
		Use:   "hub-service-create",
		Short: "Creates or updates the PubSub Hub service in the currently running solution.",
		Args:  cobra.NoArgs,
		RunE:  commandWrapper.RunE,
	}

	flags.SetCommand(cmd)

	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagsProjectOrg()
	flags.StringSlice(
		keySpokeEndpoints,
		[]string{},
		"Spoke endpoint specifications (<workcell_name>@{local|remote|url})")
	flags.OptionalString(
		keyHubServiceVersion,
		defaultHubServiceVersion,
		fmt.Sprintf(
			"Version of the service asset to install. The default value is %v",
			defaultHubServiceVersion))

	return cmd
}

func init() {
	PubsubCmd.AddCommand(NewHubServiceCreateCmd())
}
