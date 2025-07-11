// Copyright 2023 Intrinsic Innovation LLC

// Package get contains commands for introspecting the state of a running service asset in a
// solution.
package get

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/services/inctl/state/stateutils"
	"intrinsic/tools/inctl/util/printer"

	systemservicestategrpcpb "intrinsic/assets/services/proto/v1/system_service_state_go_grpc_proto"
	systemservicestatepb "intrinsic/assets/services/proto/v1/system_service_state_go_grpc_proto"
)

// Command returns a command to get the state of a running service asset in a solution.
func Command() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Get the state of a running Service instance in a solution.",
		Long:  `Get the state of a running Service instance in a solution. The default output is a condensed view of the state. Use --output=json to get a more verbose output.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer stop()

			ctx, conn, address, err := clientutils.DialClusterFromInctl(ctx, flags)
			if err != nil {
				return fmt.Errorf("could not create connection to cluster: %w", err)
			}
			defer conn.Close()

			client := systemservicestategrpcpb.NewSystemServiceStateClient(conn)
			authCtx := clientutils.AuthInsecureConn(ctx, address, flags.GetFlagProject())

			res, err := client.GetInstanceState(authCtx, &systemservicestatepb.GetInstanceStateRequest{
				Name: args[0],
			})
			if err != nil {
				return fmt.Errorf("could not get state: %w", err)
			}

			fmt.Println(&stateutils.StatePrinter{
				Proto:      res,
				OutputType: printer.GetFlagOutputType(cmd),
			})
			return nil
		},
	}

	flags.SetCommand(cmd)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagsProjectOrg()
	return cmd
}
