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
	flagHubDeleteParams      = viper.New()
	flagHubDeleteHubWorkcell string
)

// HubDeleteCmdRunner handles execution of the hub-delete subcommand.
type HubDeleteCmdRunner struct {
	NewClient func(ctx context.Context) (pubsubpb.PubSubConnectServiceClient, error)
}

// RunE implements the core logic of the delete command.
func (r *HubDeleteCmdRunner) RunE(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	client, err := r.NewClient(ctx)
	if err != nil {
		return err
	}

	req := &pubsubpb.DeleteHubRequest{
		HubWorkcellName: flagHubDeleteHubWorkcell,
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Initiating deletion of PubSub Hub for workcell %q...\n", flagHubDeleteHubWorkcell)
	ctx = metadata.AppendToOutgoingContext(ctx, "organization-id", flagHubDeleteParams.GetString(orgutil.KeyOrganization))
	op, err := client.DeleteHub(ctx, req)
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

	fmt.Fprintf(cmd.OutOrStdout(), "Successfully deleted PubSub Hub for workcell %q.\n", flagHubDeleteHubWorkcell)
	return nil
}

// NewHubDeleteCmd returns the initialized cobra command for hub-delete.
func NewHubDeleteCmd(runner *HubDeleteCmdRunner) *cobra.Command {
	if runner == nil {
		runner = &HubDeleteCmdRunner{
			NewClient: func(ctx context.Context) (pubsubpb.PubSubConnectServiceClient, error) {
				return newPubSubClient(ctx, flagHubDeleteParams)
			},
		}
	}

	cmd := &cobra.Command{
		Use:   "hub-delete",
		Short: "Deletes an existing PubSub Hub connection.",
		Args:  cobra.NoArgs,
		RunE:  runner.RunE,
	}

	flags := cmd.Flags()
	flags.StringVar(&flagHubDeleteHubWorkcell, "hub-workcell", "", "The name of the hub workcell.")
	cmd.MarkFlagRequired("hub-workcell")

	return orgutil.WrapCmd(cmd, flagHubDeleteParams)
}

func init() {
	PubsubCmd.AddCommand(NewHubDeleteCmd(nil))
}
