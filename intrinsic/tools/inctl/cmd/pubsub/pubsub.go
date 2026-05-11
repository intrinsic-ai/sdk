// Copyright 2023 Intrinsic Innovation LLC

// Package pubsub implements commands for managing pubsub network components.
package pubsub

import (
	"context"
	"fmt"
	"io"
	"time"

	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/cmd/root"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	pubsubpb "intrinsic/platform/pubsub/connect/cloud/proto/v1alpha1/pubsub_connect_go_proto"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
)

var operationPollInterval = 2 * time.Second

// PubsubCmd provides the parent command for all pubsub management commands.
var PubsubCmd = &cobra.Command{
	Use:        "pubsub",
	Short:      "Manages Intrinsic PubSub connected services.",
	Long:       "Manages Intrinsic PubSub connected services including Hub creation, deletion and traffic configuration.",
	SuggestFor: []string{"pub-sub", "mq", "bus"},
}

func init() {
	root.RootCmd.AddCommand(PubsubCmd)
}

// newPubSubClient establishes an authenticated connection to the PubSub Connect Service.
func newPubSubClient(ctx context.Context, v *viper.Viper) (pubsubpb.PubSubConnectServiceClient, error) {
	conn, err := auth.NewCloudConnection(ctx, auth.WithFlagValues(v))
	if err != nil {
		return nil, err
	}
	return pubsubpb.NewPubSubConnectServiceClient(conn), nil
}

// waitForOperation continuously polls the long running operation using client.GetOperation
// until it reaches the completed state, per requirement.
func waitForOperation(ctx context.Context, client pubsubpb.PubSubConnectServiceClient, op *lropb.Operation, out io.Writer) (*lropb.Operation, error) {
	if op == nil {
		return nil, fmt.Errorf("no operation to wait for")
	}
	if op.GetDone() {
		if op.GetError() != nil {
			return nil, fmt.Errorf("operation %q failed immediately: %v", op.GetName(), op.GetError())
		}
		return op, nil
	}

	ticker := time.NewTicker(operationPollInterval)
	defer ticker.Stop()

	req := &lropb.GetOperationRequest{Name: op.GetName()}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			updatedOp, err := client.GetOperation(ctx, req)
			if err != nil {
				return nil, err
			}
			if updatedOp.GetDone() {
				if updatedOp.GetError() != nil {
					return nil, fmt.Errorf("operation %q failed: %v", updatedOp.GetName(), updatedOp.GetError())
				}
				return updatedOp, nil
			}
			fmt.Fprintf(out, "Waiting for operation %q to complete...\n", updatedOp.GetName())
		}
	}
}
