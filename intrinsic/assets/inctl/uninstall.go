// Copyright 2023 Intrinsic Innovation LLC

// Package uninstall defines the command to uninstall an asset.
package uninstall

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"

	lrogrpcpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
)

// GetCommand returns a command to uninstall an asset.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "uninstall <id>",
		Short: "Uninstall an asset (Note: This will fail if there are instances of it in the solution.)",
		Example: `
		$ inctl asset uninstall ai.intrinsic.box \
				--project my_project \
				--solution my_solution_id

		To find a running solution's id, run:
		$ inctl solution list --project my-project --filter "running_on_hw,running_in_sim" --output json

		Can also use:
		$ inctl asset uninstall <id> --project my_project --address my_address
		or
		$ inctl asset uninstall <id> --project my_project --cluster my_cluster
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			idString := args[0]
			id, err := idutils.NewIDProto(idString)
			if err != nil {
				return err
			}

			ctx, conn, _, err := clientutils.DialClusterFromInctl(ctx, flags)
			if err != nil {
				return fmt.Errorf("could not connect to cluster: %w", err)
			}
			defer conn.Close()

			client := iagrpcpb.NewInstalledAssetsClient(conn)
			op, err := client.DeleteInstalledAsset(ctx, &iapb.DeleteInstalledAssetRequest{
				Asset: id,
			})
			if err != nil {
				return fmt.Errorf("could not uninstall the asset: %w", err)
			}

			log.Printf("Awaiting completion of the uninstallation")
			lroClient := lrogrpcpb.NewOperationsClient(conn)
			for !op.GetDone() {
				op, err = lroClient.WaitOperation(ctx, &lropb.WaitOperationRequest{
					Name: op.GetName(),
				})
				if err != nil {
					return fmt.Errorf("waiting for uninstallation failed: %w", err)
				}
			}

			if err := status.ErrorProto(op.GetError()); err != nil {
				return fmt.Errorf("uninstalling failed: %w", err)
			}
			log.Printf("Finished uninstalling %q", idString)

			return nil
		},
	}

	flags.SetCommand(cmd)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagsProjectOrg()

	return cmd
}
