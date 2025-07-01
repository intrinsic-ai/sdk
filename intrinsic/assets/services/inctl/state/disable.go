// Copyright 2023 Intrinsic Innovation LLC

// Package disable provides a command to disable a running service asset in a solution.
package disable

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

// Command returns a command to disable a running service asset in a solution.
func Command() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "disable <name>",
		Short: "Disable a running Service instance in a solution.",
		Long:  `Disable a running Service instance in a solution.`,
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
			if _, err := client.DisableService(authCtx, &systemservicestatepb.DisableServiceRequest{
				Name: name,
			}); err != nil {
				return fmt.Errorf("could not disable service: %w", err)
			}

			fmt.Printf("Service %q has been disabled.\n", name)

			return nil
		},
	}
	flags.SetCommand(cmd)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagsProjectOrg()
	return cmd
}
