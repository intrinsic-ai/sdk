// Copyright 2023 Intrinsic Innovation LLC

package solution

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/config/operationmode"
	"intrinsic/tools/inctl/util/printer"

	opmodepb "intrinsic/config/proto/operation_mode_go_proto"
	deploygrpcpb "intrinsic/kubernetes/workcell_spec/proto/deploy_go_grpc_proto"
	deploypb "intrinsic/kubernetes/workcell_spec/proto/deploy_go_grpc_proto"
)

func startSolution(ctx context.Context, conn *grpc.ClientConn, solutionID string, operationMode string) error {
	client := deploygrpcpb.NewDeployServiceClient(conn)
	opMode := operationmode.FromString(operationMode)
	if opMode == opmodepb.OperationMode_OPERATION_MODE_UNSPECIFIED {
		return fmt.Errorf("invalid operation mode: %q", operationMode)
	}
	req := &deploypb.StartSolutionRequest{SolutionId: solutionID, OperationMode: opMode}
	_, err := client.StartSolution(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to start solution: %w", err)
	}

	return nil
}

// StartCmd returns the solution start command.
func StartCmd() *cobra.Command {
	var flagOperationMode string

	viperLocal := viper.New()
	flags := cmdutils.NewCmdFlagsWithViper(viperLocal)

	solutionStartCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a solution",
		Long:  "Start a solution by name on a given cluster",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			solutionName := args[0]
			printer, err := printer.NewPrinter(cmd.Flags().Lookup("output").Value.String())
			if err != nil {
				return err
			}

			ctx := cmd.Context()

			_, clusterFlag, _, err := flags.GetFlagsAddressClusterSolution()
			if err != nil {
				return err
			}
			printer.PrintSf("Starting solution '%s' on cluster '%s'\n", solutionName, clusterFlag)

			ctx, conn, _, err := clientutils.DialClusterFromInctl(ctx, flags)
			if err != nil {
				return err
			}
			defer conn.Close()

			if err := startSolution(ctx, conn, solutionName, flagOperationMode); err != nil {
				// intrinsic:*:begin_strip
				if status.Code(err) == codes.NotFound && strings.HasSuffix(solutionName, "_APPLIC") {
					printer.PrintSf(
						"Failed to start solution '%s' which ends with '_APPLIC', which normally will not work with 'inctl solution start'."+
							"Did you mean to use 'inctl app start-from-catalog' instead?\n",
						solutionName)
				}
				// intrinsic:*:end_strip
				return err
			}

			return nil
		},
	}

	flags.SetCommand(solutionStartCmd)
	flags.AddFlagsProjectOrg()
	flags.AddFlagsAddressClusterSolution()

	solutionStartCmd.Flags().StringVar(&flagOperationMode, "operation-mode", "sim",
		"The operation mode to start the solution in, one of 'sim' (default) or 'real'.")

	return solutionStartCmd
}
