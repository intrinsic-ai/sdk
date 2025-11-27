// Copyright 2023 Intrinsic Innovation LLC

// Package list provides a command to list the states of all running service assets in a
// solution.
package list

import (
	"fmt"
	"os"
	"os/signal"
	"sort"

	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/services/inctl/state/stateutils"
	"intrinsic/tools/inctl/util/printer"

	"github.com/spf13/cobra"

	systemservicestategrpcpb "intrinsic/assets/services/proto/v1/system_service_state_go_grpc_proto"
	systemservicestatepb "intrinsic/assets/services/proto/v1/system_service_state_go_grpc_proto"
)

// Command returns a command to list the states of all running service assets in a solution.
func Command() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the states of all running Service instances in a solution.",
		Long:  `List the states of all running Service instances in a solution. The default output is a condensed view of the state. Use --output=json to get a more verbose output.`,
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

			res, err := client.ListInstanceStates(authCtx, &systemservicestatepb.ListInstanceStatesRequest{})
			if err != nil {
				return fmt.Errorf("could not get state: %w", err)
			}

			states := res.GetStates()
			sort.Slice(states, func(i, j int) bool {
				return states[i].GetName() < states[j].GetName()
			})
			fmt.Print(&stateutils.ListStatesPrinter{
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
