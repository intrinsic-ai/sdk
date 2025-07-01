// Copyright 2023 Intrinsic Innovation LLC

// Package enable provides a command to enable a running service asset in a solution.
package enable

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	systemservicestategrpcpb "intrinsic/assets/services/proto/v1/system_service_state_go_grpc_proto"
	systemservicestatepb "intrinsic/assets/services/proto/v1/system_service_state_go_grpc_proto"
)

// Command returns a command to enable a running service asset in a solution.
func Command() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "enable <name>",
		Short: "Enable a running Service instance in a solution.",
		Long:  `Enable a running Service instance in a solution.`,
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

			name := args[0]
			if _, err := client.EnableService(authCtx, &systemservicestatepb.EnableServiceRequest{
				Name: name,
			}); err != nil {
				return fmt.Errorf("could not enable service: %w", err)
			}

			fmt.Printf("Service %q has been enabled.\n", name)

			return nil
		},
	}
	flags.SetCommand(cmd)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagsProjectOrg()
	return cmd
}
