// Copyright 2023 Intrinsic Innovation LLC

// Package restart provides a command to restart a running service asset in a solution.
package restart

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/cobra"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	systemservicestategrpcpb "intrinsic/assets/services/proto/v1/system_service_state_go_grpc_proto"
	systemservicestatepb "intrinsic/assets/services/proto/v1/system_service_state_go_grpc_proto"
)

var (
	skipConfirmation = false
)

// Command returns a command to restart a running service asset in a solution.
func Command() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "restart <instance_name>",
		Short: "Restart a running Service instance in a solution.",
		Long:  `Restart a running Service instance in a solution.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer stop()

			if !skipConfirmation {
				consoleIO := bufio.NewReadWriter(
					bufio.NewReader(cmd.InOrStdin()),
					bufio.NewWriter(cmd.OutOrStdout()),
				)

				fmt.Fprintf(consoleIO,
					"The service will be temporarily unavailable during a restart. Any ongoing work may be interrupted or lost.\n\nAre you sure you want to restart %q? [Y/n] ", args[0])
				consoleIO.Flush()
				response, err := consoleIO.ReadString('\n')
				if err != nil {
					return fmt.Errorf("read response: %w", err)
				}

				response = strings.ToLower(strings.TrimSpace(response))
				if response != "y" {
					return nil
				}
			}

			ctx, conn, address, err := clientutils.DialClusterFromInctl(ctx, flags)
			if err != nil {
				return fmt.Errorf("could not create connection to cluster: %w", err)
			}
			defer conn.Close()

			client := systemservicestategrpcpb.NewSystemServiceStateClient(conn)
			authCtx := clientutils.AuthInsecureConn(ctx, address, flags.GetFlagProject())

			name := args[0]
			if _, err := client.RestartService(authCtx, &systemservicestatepb.RestartServiceRequest{
				Name: name,
			}); err != nil {
				return fmt.Errorf("could not restart service: %w", err)
			}

			fmt.Printf("Service %q is restarting. This may take a few seconds.\n", name)
			return nil
		},
	}
	flags.SetCommand(cmd)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagsProjectOrg()
	cmd.Flags().BoolVarP(&skipConfirmation, "skip_confirmation", "y", false, "Skip confirmation prompt")
	return cmd
}
