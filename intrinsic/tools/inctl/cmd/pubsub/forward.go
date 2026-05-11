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
	flagForwardParams             = viper.New()
	flagForwardWorkcell           string
	flagForwardAllowedTopics      []string
	flagForwardAllowedKVStoreKeys []string
)

// ForwardCmdRunner handles execution of the forward (configure spoke) subcommand.
type ForwardCmdRunner struct {
	NewClient func(ctx context.Context) (pubsubpb.PubSubConnectServiceClient, error)
}

// RunE implements the core logic for setting forwarding topics.
func (r *ForwardCmdRunner) RunE(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	client, err := r.NewClient(ctx)
	if err != nil {
		return err
	}

	req := &pubsubpb.ConfigureSpokeRequest{
		SpokeWorkcellName:  flagForwardWorkcell,
		AllowedTopics:      flagForwardAllowedTopics,
		AllowedKvstoreKeys: flagForwardAllowedKVStoreKeys,
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Configuring traffic forwarding for spoke %q...\n", flagForwardWorkcell)
	ctx = metadata.AppendToOutgoingContext(ctx, "organization-id", flagForwardParams.GetString(orgutil.KeyOrganization))
	op, err := client.ConfigureSpoke(ctx, req)
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

	if !status.GetSucceeded() {
		return fmt.Errorf("backend failed to delete hub: %v", status.GetErrorMessage())
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Successfully configured forwarding allowed topics for spoke %q.\n", flagForwardWorkcell)
	return nil
}

// NewForwardCmd returns the initialized cobra command for forward.
func NewForwardCmd(runner *ForwardCmdRunner) *cobra.Command {
	if runner == nil {
		runner = &ForwardCmdRunner{
			NewClient: func(ctx context.Context) (pubsubpb.PubSubConnectServiceClient, error) {
				return newPubSubClient(ctx, flagForwardParams)
			},
		}
	}

	cmd := &cobra.Command{
		Use:   "forward",
		Short: "Configures the topics/key-store-keys allowed to forward",
		Args:  cobra.NoArgs,
		RunE:  runner.RunE,
	}

	flags := cmd.Flags()
	flags.StringVar(&flagForwardWorkcell, "workcell", "", "The name of the workcell to setup traffic forwarding for.")
	flags.StringSliceVar(&flagForwardAllowedTopics, "topic", []string{}, "Topics to allow for forwarding (can be specified multiple times).")
	flags.StringSliceVar(&flagForwardAllowedKVStoreKeys, "kvstore-key", []string{}, "KeyStoreKeys to allow for forwarding (can be specified multiple times).")
	cmd.MarkFlagRequired("workcell")

	return orgutil.WrapCmd(cmd, flagForwardParams)
}

func init() {
	PubsubCmd.AddCommand(NewForwardCmd(nil))
}
