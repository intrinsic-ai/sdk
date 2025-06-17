// Copyright 2023 Intrinsic Innovation LLC

// Package delete defines the command which deletes a service instance from the
// solution.
package delete

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	lrogrpcpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/spf13/cobra"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	adgrpcpb "intrinsic/assets/proto/asset_deployment_go_grpc_proto"
	adpb "intrinsic/assets/proto/asset_deployment_go_grpc_proto"
)

// GetCommand returns a command to delete a service instance from a solution.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "delete name",
		Short: "Delete a service instance from a solution",
		Example: `
Delete a service instance with the specified name
$ inctl service delete --project=my_project --cluster=some_cluster my_instance
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Generally try to cancel calls if the user hits ctrl-c
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer stop()
			name := args[0]

			ctx, conn, _, err := clientutils.DialClusterFromInctl(ctx, flags)
			if err != nil {
				return fmt.Errorf("could not create connection to cluster: %w", err)
			}
			defer conn.Close()

			log.Printf("Requesting deletion of %q", name)
			client := adgrpcpb.NewAssetDeploymentServiceClient(conn)
			op, err := client.DeleteResource(ctx, &adpb.DeleteResourceRequest{
				Name:             name,
				DeletionStrategy: adpb.DeleteResourceRequest_DELETE_INSTANCE_ONLY,
			})
			if err != nil {
				return fmt.Errorf("could not delete service %q: %v", name, err)
			}

			lroClient := lrogrpcpb.NewOperationsClient(conn)
			defer func() {
				if !op.GetDone() {
					log.Printf("Cancelling unfinished operation")
					// Assume ctx has been cancelled if we're here.
					ctx, cancel := context.WithTimeout(cmd.Context(), 1*time.Second)
					defer cancel()
					if _, err = lroClient.CancelOperation(ctx, &lropb.CancelOperationRequest{
						Name: op.GetName(),
					}); err != nil {
						log.Printf("Cancelling failed: %v", err)
					}
				}
			}()
			log.Printf("Awaiting completion of the delete operation")
			for !op.GetDone() {
				op, err = lroClient.WaitOperation(ctx, &lropb.WaitOperationRequest{
					Name: op.GetName(),
				})
				if err != nil {
					return fmt.Errorf("unable to check status of delete operation for %q: %v", name, err)
				}
			}

			if err := op.GetError(); err != nil {
				return fmt.Errorf("failed to delete %q: %v", name, err)
			}

			log.Printf("Deleted service %q", name)
			return nil
		},
	}

	flags.SetCommand(cmd)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagsProjectOrg()

	return cmd
}
