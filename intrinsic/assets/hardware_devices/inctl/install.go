// Copyright 2023 Intrinsic Innovation LLC

// Package install defines the command to install a HardwareDevice.
package install

import (
	"fmt"
	"log"
	"os"

	lrogrpcpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"intrinsic/assets/bundleio"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	"intrinsic/assets/imagetransfer"
	"intrinsic/assets/imageutils"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
	"intrinsic/skills/tools/resource/cmd/bundleimages"
	"intrinsic/skills/tools/skill/cmd/directupload/directupload"

	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

const (
)

// GetCommand returns a command to install a HardwareDevice.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "install <bundle>",
		Short: "Install a HardwareDevice.",
		Example: `
	Install a HardwareDevice to the specified solution:
	$ inctl hardware_device install abc/bundle.tar \
			--org my_org \
			--solution my_solution_id

	To find a running solution's id, run:
	$ inctl solution list --project my-project --filter "running_on_hw,running_in_sim" --output json

	The HardwareDevice can also be installed by specifying the cluster on which the solution is running:
	$ inctl hardware_device install abc/bundle.tar \
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

			useDirectUpload := true
			var transferer imagetransfer.Transferer
			if useDirectUpload {
				dopts := []directupload.Option{
					directupload.WithDiscovery(directupload.NewCatalogTarget(conn)),
					directupload.WithOutput(cmd.OutOrStdout()),
				}
				transferer = directupload.NewTransferer(ctx, dopts...)
			}
			assetInliner := bundleio.NewLocalAssetInliner(bundleio.LocalAssetInlinerOptions{
				ImageProcessor: bundleimages.CreateImageProcessor(bundleimages.RegistryOptions{
					Transferer: transferer,
					URI:        imageutils.GetRegistry(clientutils.ResolveCatalogProjectFromInctl(flags)),
				}),
			})

			localAssetsDir, err := os.MkdirTemp("", "local-assets")
			if err != nil {
				return fmt.Errorf("could not create temporary directory for local assets: %w", err)
			}
			defer os.RemoveAll(localAssetsDir)

			hwd, err := bundleio.ProcessHardwareDevice(target,
				bundleio.WithProcessAsset(assetInliner.Process),
				bundleio.WithReadOptions(
					bundleio.WithExtractLocalAssetsDir(localAssetsDir),
				),
			)
			if err != nil {
				return fmt.Errorf("could not process HardwareDevice bundle: %w", err)
			}

			id, err := idutils.IDFromProto(hwd.GetMetadata().GetId())
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}
			log.Printf("Installing HardwareDevice %q", id)

			// This needs an authorized context to pull from the catalog if not available.
			authCtx := clientutils.AuthInsecureConn(ctx, address, flags.GetFlagProject())
			op, err := client.CreateInstalledAsset(authCtx, &iapb.CreateInstalledAssetRequest{
				Policy: policy,
				Asset: &iapb.CreateInstalledAssetRequest_Asset{
					Variant: &iapb.CreateInstalledAssetRequest_Asset_HardwareDevice{
						HardwareDevice: hwd,
					},
				},
			})
			if err != nil {
				return fmt.Errorf("could not install the HardwareDevice: %w", err)
			}

			log.Printf("Awaiting completion of the installation")
			lroClient := lrogrpcpb.NewOperationsClient(conn)
			for !op.GetDone() {
				op, err = lroClient.WaitOperation(ctx, &lropb.WaitOperationRequest{
					Name: op.GetName(),
				})
				if err != nil {
					return fmt.Errorf("waiting for installation failed: %w", err)
				}
			}

			if err := status.ErrorProto(op.GetError()); err != nil {
				return fmt.Errorf("installation failed: %w", err)
			}
			log.Printf("Finished installing %q", id)

			return nil
		},
	}

	flags.SetCommand(cmd)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagPolicy("hardware device")
	flags.AddFlagsProjectOrg()
	flags.AddFlagRegistry()
	flags.AddFlagsRegistryAuthUserPassword()
	flags.AddFlagSkipDirectUpload("hardware device")

	return cmd
}

func remoteOpt() remote.Option {
	return remote.WithAuthFromKeychain(google.Keychain)
}
