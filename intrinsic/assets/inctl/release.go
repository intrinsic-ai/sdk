// Copyright 2023 Intrinsic Innovation LLC

// Package release defines the command that releases an asset to the catalog.
package release

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"intrinsic/assets/bundleio"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	"intrinsic/assets/imagetransfer"
	"intrinsic/assets/imageutils"
	"intrinsic/assets/services/bundleimages"
	"intrinsic/skills/tools/skill/cmd/directupload/directupload"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"

	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	rmpb "intrinsic/assets/catalog/proto/v1/release_metadata_go_proto"
)

const (
)

// GetCommand returns command to release an asset.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()

	cmd := &cobra.Command{
		Use:   "release bundle.tar",
		Short: "Release an Asset to the catalog.",
		Example: `
	Release an Asset to the catalog
	$ inctl asset release abc/bundle.tar --version=0.0.1
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printer, err := printer.NewPrinter(root.FlagOutput)
			if err != nil {
				return err
			}

			ctx, conn, err := clientutils.DialCatalogFromInctl(cmd, flags)
			if err != nil {
				return fmt.Errorf("failed to create client connection: %v", err)
			}
			defer conn.Close()

			client := acgrpcpb.NewAssetCatalogClient(conn)

			req, err := makeCreateAssetRequest(ctx, makeCreateAssetRequestOptions{
				acClient:       client,
				conn:           conn,
				flags:          flags,
				progressWriter: cmd.OutOrStdout(),
				target:         args[0],
			})
			if err != nil {
				return err
			}

			idVersion, err := idutils.IDVersionFromProto(req.GetAsset().GetMetadata().GetIdVersion())
			if err != nil {
				return err
			}
			printer.PrintSf("Releasing %q to the asset catalog", idVersion)

			if flags.GetFlagDryRun() {
				printer.PrintS("Skipping release: dry-run")
				return nil
			}

			return release(cmd.Context(), client, req, flags.GetFlagIgnoreExisting(), printer)
		},
	}
	flags.SetCommand(cmd)
	flags.AddFlagDefault("asset")
	flags.AddFlagDryRun()
	flags.AddFlagIgnoreExisting("asset")
	flags.AddFlagOrganizationOptional()
	flags.AddFlagOrgPrivate()
	flags.AddFlagReleaseNotes("asset")
	flags.AddFlagVersion("asset")

	return cmd
}

func release(ctx context.Context, client acgrpcpb.AssetCatalogClient, req *acpb.CreateAssetRequest, ignoreExisting bool, printer printer.Printer) error {
	if _, err := client.CreateAsset(ctx, req); err != nil {
		if s, ok := status.FromError(err); ok && ignoreExisting && s.Code() == codes.AlreadyExists {
			printer.PrintS("Skipping release: asset already exists in the catalog")
			return nil
		}
		return fmt.Errorf("could not release the HardwareDevice: %w", err)
	}
	printer.PrintS("Finished releasing the HardwareDevice")
	return nil
}

type makeCreateAssetRequestOptions struct {
	acClient       acgrpcpb.AssetCatalogClient
	conn           *grpc.ClientConn
	flags          *cmdutils.CmdFlags
	progressWriter io.Writer
	target         string
}

func makeCreateAssetRequest(ctx context.Context, opts makeCreateAssetRequestOptions) (*acpb.CreateAssetRequest, error) {
	var transferer imagetransfer.Transferer
	if true {
		transferer = directupload.NewTransferer(ctx,
			directupload.WithDiscovery(directupload.NewCatalogTarget(opts.conn)),
			directupload.WithOutput(opts.progressWriter),
			directupload.WithFailOver(transferer),
		)
	}
	processor := bundleio.BundleProcessor{
		ImageProcessor: bundleimages.CreateImageProcessor(bundleimages.RegistryOptions{
			Transferer: transferer,
			URI:        imageutils.GetRegistry(clientutils.ResolveCatalogProjectFromInctl(opts.flags)),
		}),
		ProcessReferencedData:   bundleio.ToCatalogReferencedData(ctx, bundleio.WithACClient(opts.acClient)),
	}

	processedBundle, err := processor.Process(opts.target)
	if err != nil {
		return nil, fmt.Errorf("unable to process: %w", err)
	}

	asset := processedBundle.Release(bundleio.VersionDetails{
		Version:      opts.flags.GetFlagVersion(),
		ReleaseNotes: opts.flags.GetFlagReleaseNotes(),
		ReleaseMetadata: &rmpb.ReleaseMetadata{
			Default:    opts.flags.GetFlagDefault(),
			OrgPrivate: opts.flags.GetFlagOrgPrivate(),
		},
	})

	// This is a very asset-specific piece of validation, but it didn't seem to
	// make sense to bury it in Release, and force every other implementation
	// to add an error case.
	if len(asset.GetDeploymentData().GetHardwareDeviceSpecificDeploymentData().GetManifest().GetMetadata().GetAssetTags()) > 1 {
		return nil, fmt.Errorf("HardwareDevice %q specifies more than one asset tag, but at most one is allowed", idutils.IDFromProtoUnchecked(asset.GetMetadata().GetIdVersion().GetId()))
	}

	return &acpb.CreateAssetRequest{
		Asset: asset,
	}, nil
}
