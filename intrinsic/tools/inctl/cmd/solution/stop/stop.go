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

	deploygrpcpb "intrinsic/kubernetes/workcell_spec/proto/deploy_go_grpc_proto"
	deploypb "intrinsic/kubernetes/workcell_spec/proto/deploy_go_grpc_proto"
)

func stopSolution(ctx context.Context, conn *grpc.ClientConn) error {
	client := deploygrpcpb.NewDeployServiceClient(conn)
	req := &deploypb.StopSolutionRequest{}
	_, err := client.StopSolution(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to stop solution: %w", err)
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
