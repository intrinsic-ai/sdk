// Copyright 2023 Intrinsic Innovation LLC

// Package release defines the service release command which releases a service to the catalog.
package release

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"intrinsic/assets/bundleio"
	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	rmpb "intrinsic/assets/catalog/proto/v1/release_metadata_go_proto"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	"intrinsic/assets/imagetransfer"
	"intrinsic/assets/imageutils"
	atpb "intrinsic/assets/proto/asset_type_go_proto"
	mpb "intrinsic/assets/proto/metadata_go_proto"
	"intrinsic/assets/services/bundleimages"
	"intrinsic/skills/tools/skill/cmd/directupload/directupload"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"
)

type imageTransfererOpts struct {
	cmd             *cobra.Command
	conn            *grpc.ClientConn
	useDirectUpload bool
}

func imageTransferer(opts imageTransfererOpts) imagetransfer.Transferer {
	var transferer imagetransfer.Transferer
	if opts.useDirectUpload {
		dopts := []directupload.Option{
			directupload.WithDiscovery(directupload.NewCatalogTarget(opts.conn)),
			directupload.WithOutput(opts.cmd.OutOrStdout()),
		}
		transferer = directupload.NewTransferer(opts.cmd.Context(), dopts...)
	}
	return transferer
}

func remoteOpt() remote.Option {
	return remote.WithAuthFromKeychain(google.Keychain)
}

func release(ctx context.Context, client acgrpcpb.AssetCatalogClient, req *acpb.CreateAssetRequest, ignoreExisting bool, printer printer.Printer) error {
	if _, err := client.CreateAsset(ctx, req); err != nil {
		if s, ok := status.FromError(err); ok && ignoreExisting && s.Code() == codes.AlreadyExists {
			printer.PrintS("Skipping release: service already exists in the catalog")
			return nil
		}
		return errors.Wrap(err, "could not release the service")
	}
	printer.PrintS("Finished releasing the service")
	return nil
}

func processAsset(target string, transferer imagetransfer.Transferer, flags *cmdutils.CmdFlags) (*acpb.Asset, error) {
	if flags.GetFlagDryRun() {
		manifest, err := bundleio.ReadServiceManifest(target)
		if err != nil {
			return nil, fmt.Errorf("failed to read service manifest from bundle: %v", err)
		}
		id := manifest.GetMetadata().GetId()
		idVersion, err := idutils.IDVersionProtoFrom(id.GetPackage(), id.GetName(), flags.GetFlagVersion())
		if err != nil {
			return nil, err
		}
		return &acpb.Asset{
			Metadata: &mpb.Metadata{
				IdVersion:     idVersion,
				AssetType:     atpb.AssetType_ASSET_TYPE_SERVICE,
				AssetTag:      manifest.GetMetadata().GetAssetTag(),
				DisplayName:   manifest.GetMetadata().GetDisplayName(),
				Documentation: manifest.GetMetadata().GetDocumentation(),
				Vendor:        manifest.GetMetadata().GetVendor(),
				ReleaseNotes:  flags.GetFlagReleaseNotes(),
			},
		}, nil
	}

	opts := bundleio.ProcessServiceOpts{
		ImageProcessor: bundleimages.CreateImageProcessor(bundleimages.RegistryOptions{
			Transferer: transferer,
			URI:        imageutils.GetRegistry(clientutils.ResolveCatalogProjectFromInctl(flags)),
		}),
	}
	psm, err := bundleio.ProcessService(target, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to process service bundle: %v", err)
	}
	version := flags.GetFlagVersion()
	metadata := psm.GetMetadata()
	idVersion, err := idutils.IDVersionProtoFrom(metadata.GetId().GetPackage(), metadata.GetId().GetName(), version)
	if err != nil {
		return nil, err
	}
	return &acpb.Asset{
		Metadata: &mpb.Metadata{
			IdVersion:     idVersion,
			AssetType:     atpb.AssetType_ASSET_TYPE_SERVICE,
			AssetTag:      metadata.GetAssetTag(),
			DisplayName:   metadata.GetDisplayName(),
			Documentation: metadata.GetDocumentation(),
			Vendor:        metadata.GetVendor(),
			ReleaseNotes:  flags.GetFlagReleaseNotes(),
		},
		ReleaseMetadata: &rmpb.ReleaseMetadata{
			Default:    flags.GetFlagDefault(),
			OrgPrivate: flags.GetFlagOrgPrivate(),
		},
		DeploymentData: &acpb.Asset_AssetDeploymentData{
			AssetSpecificDeploymentData: &acpb.Asset_AssetDeploymentData_ServiceSpecificDeploymentData{
				ServiceSpecificDeploymentData: &acpb.Asset_ServiceDeploymentData{
					Manifest: psm,
				},
			},
		},
	}, nil
}

// GetCommand returns command to release services.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()

	cmd := &cobra.Command{
		Use:   "release bundle.tar",
		Short: "Release a service to the catalog",
		Example: `
	Release a service to the catalog
	$ service release abc/bundle.tar --version=0.0.1
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			dryRun := flags.GetFlagDryRun()

			useDirectUpload := true

			var conn *grpc.ClientConn
			var transferer imagetransfer.Transferer
			ctx := cmd.Context()
			if !dryRun {
				var err error
				ctx, conn, err = clientutils.DialCatalogFromInctl(cmd, flags)
				if err != nil {
					return fmt.Errorf("failed to create client connection: %v", err)
				}
				defer conn.Close()
				transferer = imageTransferer(imageTransfererOpts{
					cmd:             cmd,
					conn:            conn,
					useDirectUpload: useDirectUpload,
				})
			}

			asset, err := processAsset(target, transferer, flags)
			if err != nil {
				return err
			}
			printer, err := printer.NewPrinter(root.FlagOutput)
			if err != nil {
				return err
			}
			idVersion := idutils.IDVersionFromProtoUnchecked(asset.GetMetadata().GetIdVersion())
			printer.PrintSf("Releasing service %q to the asset catalog", idVersion)
			if dryRun {
				printer.PrintS("Skipping call to asset catalog (dry-run)")
				return nil
			}
			client := acgrpcpb.NewAssetCatalogClient(conn)
			req := &acpb.CreateAssetRequest{
				Asset: asset,
			}
			return release(ctx, client, req, flags.GetFlagIgnoreExisting(), printer)
		},
	}
	flags.SetCommand(cmd)
	flags.AddFlagDefault("service")
	flags.AddFlagDryRun()
	flags.AddFlagIgnoreExisting("service")
	flags.AddFlagOrganizationOptional()
	flags.AddFlagOrgPrivate()
	flags.AddFlagReleaseNotes("service")
	flags.AddFlagVersion("service")

	return cmd
}
