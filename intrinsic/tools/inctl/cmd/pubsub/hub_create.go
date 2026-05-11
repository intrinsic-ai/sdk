// Copyright 2023 Intrinsic Innovation LLC

package pubsub

import (
	"context"
	"fmt"

	"intrinsic/tools/inctl/util/orgutil"

	commandpb "intrinsic/platform/pubsub/common/command_go_proto"
	pubsubpb "intrinsic/platform/pubsub/connect/cloud/proto/v1alpha1/pubsub_connect_go_proto"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"
)

var (
	flagHubCreateParams          = viper.New()
	flagHubCreateHubWorkcell     string
	flagHubCreateStaticEndpoints []string
)

// HubCreateCmdRunner handles execution of the hub-create subcommand.
type HubCreateCmdRunner struct {
	NewClient func(ctx context.Context) (pubsubpb.PubSubConnectServiceClient, error)
}

// RunE implements the core logic of the command.
func (r *HubCreateCmdRunner) RunE(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	client, err := r.NewClient(ctx)
	if err != nil {
		return err
	}

	var spokeEndpoints []*pubsubpb.Endpoint
	for _, s := range flagHubCreateStaticEndpoints {
		spokeEndpoints = append(spokeEndpoints, &pubsubpb.Endpoint{
			Endpoint: &pubsubpb.Endpoint_StaticEndpoint{
				StaticEndpoint: s,
			},
		})
	}

	req := &pubsubpb.CreateHubRequest{
		HubWorkcellName: flagHubCreateHubWorkcell,
		SpokeEndpoints:  spokeEndpoints,
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Initiating creation of PubSub Hub for workcell %q...\n", flagHubCreateHubWorkcell)
	ctx = metadata.AppendToOutgoingContext(ctx, "organization-id", flagHubCreateParams.GetString(orgutil.KeyOrganization))
	op, err := client.CreateHub(ctx, req)
	if err != nil {
		return err
	}

	op, err = waitForOperation(ctx, client, op, cmd.OutOrStdout())
	if err != nil {
		return err
	}

	status := &commandpb.CommandExecutionStatus{}
	if err := op.GetResponse().UnmarshalTo(status); err != nil {
		return fmt.Errorf("failed to parse response status: %v", err)
	}

	if status.GetErrorMessage() != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Hub connection failed: %v\n", status.GetErrorMessage())
		return fmt.Errorf(status.GetErrorMessage())
	}

	if !status.GetSucceeded() {
		return fmt.Errorf("Hub connection failed")
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\nSuccessfully created PubSub Hub for workcell %q.\n", flagHubCreateHubWorkcell)
	return nil
}

// NewHubCreateCmd returns the initialized cobra command for hub-create.
func NewHubCreateCmd(runner *HubCreateCmdRunner) *cobra.Command {
	if runner == nil {
		runner = &HubCreateCmdRunner{
			NewClient: func(ctx context.Context) (pubsubpb.PubSubConnectServiceClient, error) {
				return newPubSubClient(ctx, flagHubCreateParams)
			},
		}
	}

	cmd := &cobra.Command{
		Use:   "hub-create",
		Short: "Creates a PubSub Hub connection.",
		Args:  cobra.NoArgs,
		RunE:  runner.RunE,
	}

	flags := cmd.Flags()
	flags.StringVar(&flagHubCreateHubWorkcell, "hub-workcell", "", "The name of the hub workcell.")
	flags.StringSliceVar(&flagHubCreateStaticEndpoints, "spoke-static-endpoint", []string{}, "List of static endpoints to connect to the hub.")
	cmd.MarkFlagRequired("hub-workcell")
	cmd.MarkFlagRequired("spoke-static-endpoint")

	return orgutil.WrapCmd(cmd, flagHubCreateParams)
}

func init() {
	PubsubCmd.AddCommand(NewHubCreateCmd(nil))
}
