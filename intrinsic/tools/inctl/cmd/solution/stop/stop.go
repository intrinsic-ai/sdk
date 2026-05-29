// Copyright 2023 Intrinsic Innovation LLC

// Package stop provides a command to stop a solution.
package stop

import (
	"context"
	"fmt"

	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/tools/inctl/util/printer"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	solutiondeploymentpb "intrinsic/assets/proto/v1/solution_deployment_go_proto"
	deploygrpcpb "intrinsic/kubernetes/workcell_spec/proto/deploy_go_proto"
	deploypb "intrinsic/kubernetes/workcell_spec/proto/deploy_go_proto"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
)

func stopSolution(ctx context.Context, conn *grpc.ClientConn) error {
	client := solutiondeploymentpb.NewSolutionDeploymentServiceClient(conn)
	op, err := client.DeleteSolutionDeployment(ctx, &solutiondeploymentpb.DeleteSolutionDeploymentRequest{})
	if status.Code(err) == codes.Unimplemented {
		if _, err := deploygrpcpb.NewDeployServiceClient(conn).StopSolution(ctx, &deploypb.StopSolutionRequest{}); err != nil {
			return fmt.Errorf("failed to stop solution (fallback): %w", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to stop solution: %w", err)
	}

	name := op.GetName()
	lroClient := lropb.NewOperationsClient(conn)
	for !op.GetDone() {
		op, err = lroClient.WaitOperation(ctx, &lropb.WaitOperationRequest{
			Name: name,
		})
		if err != nil {
			return fmt.Errorf("unable to check status of solution stop operation %q: %w", name, err)
		}
	}

	if err := status.ErrorProto(op.GetError()); err != nil {
		return fmt.Errorf("solution stop operation %q failed: %w", name, err)
	}

	return nil
}

// NewCommand returns the solution stop command.
func NewCommand() *cobra.Command {
	viperLocal := viper.New()
	flags := cmdutils.NewCmdFlagsWithViper(viperLocal)

	solutionStopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the solution running on a cluster",
		Long:  "Stop the solution running on a given cluster",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			printer, err := printer.NewPrinter(cmd.Flags().Lookup("output").Value.String())
			if err != nil {
				return err
			}

			ctx := cmd.Context()

			_, clusterFlag, _, err := flags.GetFlagsAddressClusterSolution()
			if err != nil {
				return err
			}
			printer.PrintSf("Stopping solution on cluster '%s'\n", clusterFlag)

			ctx, conn, _, err := clientutils.DialClusterFromInctl(ctx, flags)
			if err != nil {
				return err
			}
			defer conn.Close()

			if err = stopSolution(ctx, conn); err != nil {
				return err
			}
			return nil
		},
	}

	flags.SetCommand(solutionStopCmd)
	flags.AddFlagsProjectOrg()
	flags.AddFlagsAddressClusterSolution()

	return solutionStopCmd
}
