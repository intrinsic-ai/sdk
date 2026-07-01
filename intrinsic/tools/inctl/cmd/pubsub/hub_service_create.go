// Copyright 2023 Intrinsic Innovation LLC

package pubsub

import (
	"context"
	"fmt"
	"strings"

	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	endpointpb "intrinsic/platform/pubsub/connect/onprem/relay_router_service/endpoint_spec_go_proto"
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
	ServiceInstallingCmdRunner

	spokeEndpoints []string
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

// run creates a configuration proto for the relay service,
// and triggers installation of that service.
func (r *HubServiceCreateCmdRunner) run(ctx context.Context) error {
	if len(r.spokeEndpoints) == 0 {
		return fmt.Errorf("at least one spoke endpoint must be specified using --spoke-endpoint; cannot install or update the relay")
	}

	config, err := r.makeConfig()
	if err != nil {
		return fmt.Errorf("failed to create service config from command line flags: %w", err)
	}
	return r.updateInstalledServiceInstances(ctx, config)
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
		ServiceInstallingCmdRunner: ServiceInstallingCmdRunner{
			CmdRunnerBase: *newCmdRunnerBase(
				conn,
				cmd.OutOrStdout(),
				e.cmdFlags.GetString(cmdutils.KeyCluster),
				hubServicePackage,
				hubServiceName),
			requestedVersion: e.cmdFlags.GetString(keyHubServiceVersion),
		},
		spokeEndpoints: e.cmdFlags.GetStringSlice(keySpokeEndpoints),
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
