// Copyright 2023 Intrinsic Innovation LLC

// Package uninstall defines the command to uninstall a Service.
package uninstall

import (
	"fmt"
	"log"
	"time"

	lrogrpcpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	idpb "intrinsic/assets/proto/id_go_proto"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
)

// GetCommand returns a command to uninstall a Service.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "uninstall ID",
		Short: "Uninstall a Service type (Note: This will fail if there are instances of it in the solution)",
		Example: `
		$ inctl service uninstall ai.intrinsic.realtime_control_service \
				--project my_project \
				--solution my_solution_id

				To find a Service's id_version, run:
				$ inctl service list --org my_organization --solution my_solution_id

				To find a running solution's id, run:
				$ inctl solution list --project my-project --filter "running_on_hw,running_in_sim" --output json
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			id := args[0]
			idv, err := idutils.IDOrIDVersionProtoFrom(id)
			if err != nil {
				return fmt.Errorf("invalid identifier: %v", err)
			}
			if v := idv.GetVersion(); v != "" {
				log.Print("Warning: specifying the version of an asset is deprecated, and soon will cause an error")
			}

			ctx, conn, _, err := clientutils.DialClusterFromInctl(ctx, flags)
			if err != nil {
				return fmt.Errorf("could not connect to cluster: %w", err)
			}
			defer conn.Close()

			client := iagrpcpb.NewInstalledAssetsClient(conn)
			op, err := client.DeleteInstalledAssets(ctx, &iapb.DeleteInstalledAssetsRequest{
				Assets: []*idpb.Id{
					idv.GetId(),
				},
			})
			if err != nil {
				return fmt.Errorf("could not uninstall the Service: %w", err)
			}

			log.Printf("Awaiting completion of the uninstallation")
			lroClient := lrogrpcpb.NewOperationsClient(conn)
			for !op.GetDone() {
				time.Sleep(15 * time.Millisecond)
				op, err = lroClient.GetOperation(ctx, &lropb.GetOperationRequest{
					Name: op.GetName(),
				})
				if err != nil {
					return fmt.Errorf("unable to check status of uninstallation: %v", err)
				}
			}

			if err := status.ErrorProto(op.GetError()); err != nil {
				return fmt.Errorf("uninstalling failed: %w", err)
			}
			log.Printf("Finished uninstalling %q", id)

			return nil
		},
	}

	flags.SetCommand(cmd)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagsProjectOrg()

	return cmd
}
