// Copyright 2023 Intrinsic Innovation LLC

// Package install defines the command to install a Data asset.
package install

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
	"intrinsic/assets/bundleio"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"

	lrogrpcpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
)

// GetCommand returns a command to install a Data asset.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "install <bundle>",
		Short: "Install a Data asset.",
		Example: `
	Install a Data asset to the specified solution:
	$ inctl data install abc/data_bundle.tar \
			--org my_org \
			--solution my_solution_id

	To find a running solution's id, run:
	$ inctl solution list --project my-project --filter "running_on_hw,running_in_sim" --output json

	The Data asset can also be installed by specifying the cluster on which the solution is running:
	$ inctl data install abc/data_bundle.tar \
			--org my_org \
			--cluster my_cluster
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			target := args[0]

			policy, err := flags.GetFlagPolicy()
			if err != nil {
				return err
			}

			ctx, conn, address, err := clientutils.DialClusterFromInctl(ctx, flags)
			if err != nil {
				return err
			}
			defer conn.Close()

			client := iagrpcpb.NewInstalledAssetsClient(conn)

			da, err := bundleio.ReadDataAsset(target, bundleio.WithProcessReferencedData(bundleio.ToPortableReferencedData))
			if err != nil {
				return errors.Wrapf(err, "could not read bundle file: %s", target)
			}

			id, err := idutils.IDFromProto(da.GetMetadata().GetIdVersion().GetId())
			if err != nil {
				return fmt.Errorf("invalid id: %v", err)
			}
			log.Printf("Installing Data asset %q", id)

			// This needs an authorized context to pull from the catalog if not available.
			authCtx := clientutils.AuthInsecureConn(ctx, address, flags.GetFlagProject())
			op, err := client.CreateInstalledAsset(authCtx, &iapb.CreateInstalledAssetRequest{
				Policy: policy,
				Asset: &iapb.CreateInstalledAssetRequest_Asset{
					Variant: &iapb.CreateInstalledAssetRequest_Asset_Data{
						Data: da,
					},
				},
			})
			if err != nil {
				return fmt.Errorf("could not install the Data asset: %v", err)
			}

			log.Printf("Awaiting completion of the installation")
			if op, err := lrogrpcpb.NewOperationsClient(conn).WaitOperation(ctx, &lropb.WaitOperationRequest{
				Name: op.GetName(),
			}); err != nil {
				return fmt.Errorf("waiting for installation failed: %v", err)
			} else if err := status.ErrorProto(op.GetError()); err != nil {
				return fmt.Errorf("installation failed: %w", err)
			}
			log.Printf("Finished installing %q", id)

			return nil
		},
	}

	flags.SetCommand(cmd)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagPolicy("data")
	flags.AddFlagsProjectOrg()

	return cmd
}
