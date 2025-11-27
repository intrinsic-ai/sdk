// Copyright 2023 Intrinsic Innovation LLC

// Package install defines the command to install an asset.
package install

import (
	"fmt"
	"log"
	"os"

	"intrinsic/assets/bundleio"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	"intrinsic/assets/imagetransfer"
	"intrinsic/assets/services/bundleimages"
	"intrinsic/kubernetes/acl/identity"
	"intrinsic/skills/tools/skill/cmd/directupload/directupload"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"

	iagrpcpb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_grpc_proto"

	lrogrpcpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
)

const (
)

// GetCommand returns a command to install an asset.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "install <asset_id_version>",
		Short: "Install an asset",
		Example: `
	Install a specific asset ID version from a catalog to the specified solution:
	$ inctl asset install ai.intrinsic.calculator_service.0.20250320.0-RC01+insrc \
			--org my_org \
			--solution my_solution_id

	Install a local bundle to the specified solution:
	$ inctl asset install abc/bundle.tar \
			--org my_org \
			--solution my_solution_id

	To find a running solution's id, run:
	$ inctl solution list --org my_org --filter "running_on_hw,running_in_sim" --output json

	The asset can also be installed by specifying the cluster on which the solution is running:
	$ inctl asset install ai.intrinsic.calculator_service.0.20250320.0-RC01+insrc \
			--org my_org \
			--cluster my_cluster

	Install a local bundle to into a solution running on the specified cluster:
	$ inctl asset install abc/bundle.tar \
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

			ctx, err = identity.OrgToContext(ctx, flags.GetFlagOrganization())
			if err != nil {
				return fmt.Errorf("failed to add org information to context: %w", err)
			}
			ctx, conn, address, err := clientutils.DialClusterFromInctl(ctx, flags)
			if err != nil {
				return err
			}
			defer conn.Close()

			// Determine the image transferer to use. Default to direct injection into the cluster.
			registry := flags.GetFlagRegistry()
			remoteOpt, err := clientutils.RemoteOpt(flags)
			if err != nil {
				return err
			}
			transfer := imagetransfer.RemoteTransferer(remoteOpt)
			if !flags.GetFlagSkipDirectUpload() {
				opts := []directupload.Option{
					directupload.WithDiscovery(directupload.NewFromConnection(conn)),
					directupload.WithOutput(cmd.OutOrStdout()),
				}
				if registry != "" {
					// User set external registry, so we can use it as failover.
					opts = append(opts, directupload.WithFailOver(transfer))
				} else {
					// Fake name that ends in .local in order to indicate that this is local, directly
					// uploaded image.
					registry = "direct.upload.local"
				}
				transfer = directupload.NewTransferer(opts...)
			}
			client := iagrpcpb.NewInstalledAssetsClient(conn)
			authCtx := clientutils.AuthInsecureConn(ctx, address, flags.GetFlagProject())

			processor := bundleio.BundleProcessor{
				ImageProcessor:          bundleimages.CreateImageProcessor(flags.CreateRegistryOptsWithTransferer(ctx, transfer, registry)),
				ProcessReferencedData:   bundleio.ToPortableReferencedData,
			}

			var fileExists bool
			if _, err := os.Stat(target); err == nil {
				fileExists = true
			}
			var asset *iapb.CreateInstalledAssetRequest_Asset
			if idvParts, err := idutils.NewIDVersionParts(target); err != nil {
				if !fileExists {
					return fmt.Errorf("%q is neither a file nor a valid id_version (package.name.version); check that either the file path is correct or that the id_version is formatted correctly", target)
				}
				processedBundle, err := processor.Process(ctx, target)
				if err != nil {
					return fmt.Errorf("unable to process: %w", err)
				}
				asset = processedBundle.Install()
			} else {
				if fileExists {
					return fmt.Errorf("input is ambiguous; %q is both a file and an id_version", target)
				}
				asset = &iapb.CreateInstalledAssetRequest_Asset{
					Variant: &iapb.CreateInstalledAssetRequest_Asset_Catalog{
						Catalog: idvParts.IDVersionProto(),
					},
				}
			}

			op, err := client.CreateInstalledAsset(authCtx, &iapb.CreateInstalledAssetRequest{
				Policy: policy,
				Asset:  asset,
			})
			if err != nil {
				return fmt.Errorf("could not install the asset: %v", err)
			}

			log.Printf("Awaiting completion of the installation")
			lroClient := lrogrpcpb.NewOperationsClient(conn)
			for !op.GetDone() {
				op, err = lroClient.WaitOperation(ctx, &lropb.WaitOperationRequest{
					Name: op.GetName(),
				})
				if err != nil {
					return fmt.Errorf("unable to check status of installation: %v", err)
				}
			}

			if err := status.ErrorProto(op.GetError()); err != nil {
				return fmt.Errorf("installation failed: %w", err)
			}
			installed := &iapb.InstalledAsset{}
			if err := op.GetResponse().UnmarshalTo(installed); err != nil {
				return fmt.Errorf("unable to parse result from successful installation: %w", err)
			}
			log.Printf("Finished installing %q", idutils.IDVersionFromProtoUnchecked(installed.GetMetadata().GetIdVersion()))
			return nil
		},
	}

	flags.SetCommand(cmd)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagPolicy("asset")
	flags.AddFlagsProjectOrg()
	flags.AddFlagRegistry()
	flags.AddFlagsRegistryAuthUserPassword()
	flags.AddFlagSkipDirectUpload("asset")

	return cmd
}
