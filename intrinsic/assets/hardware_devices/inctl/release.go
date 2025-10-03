// Copyright 2023 Intrinsic Innovation LLC

// Package release defines the command that releases a HardwareDevice to the catalog.
package release

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/remote"
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
	atagpb "intrinsic/assets/proto/asset_tag_go_proto"
	atpb "intrinsic/assets/proto/asset_type_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	mpb "intrinsic/assets/proto/metadata_go_proto"
	"intrinsic/assets/services/bundleimages"
	"intrinsic/skills/tools/skill/cmd/directupload/directupload"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"
)

const (
)

// GetCommand returns command to release HardwareDevices.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()

	cmd := &cobra.Command{
		Use:   "release bundle.tar",
		Short: "Release a HardwareDevice to the catalog.",
		Example: `
	Release a HardwareDevice to the catalog
	$ hardware_device release abc/bundle.tar --version=0.0.1
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
			printer.PrintSf("Releasing HardwareDevice %q to the asset catalog", idVersion)

			if flags.GetFlagDryRun() {
				printer.PrintS("Skipping release: dry-run")
				return nil
			}

			return release(cmd.Context(), client, req, flags.GetFlagIgnoreExisting(), printer)
		},
	}
	flags.SetCommand(cmd)
	flags.AddFlagDefault("hardware device")
	flags.AddFlagDryRun()
	flags.AddFlagIgnoreExisting("hardware device")
	flags.AddFlagOrganizationOptional()
	flags.AddFlagOrgPrivate()
	flags.AddFlagReleaseNotes("hardware device")
	flags.AddFlagVersion("hardware device")

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
	useDirectUpload := true
	var transferer imagetransfer.Transferer
	if useDirectUpload {
		dopts := []directupload.Option{
			directupload.WithDiscovery(directupload.NewCatalogTarget(opts.conn)),
			directupload.WithOutput(opts.progressWriter),
		}
		transferer = directupload.NewTransferer(ctx, dopts...)
	}
	assetInliner := bundleio.NewLocalAssetInliner(bundleio.LocalAssetInlinerOptions{
		ImageProcessor: bundleimages.CreateImageProcessor(bundleimages.RegistryOptions{
			Transferer: transferer,
			URI:        imageutils.GetRegistry(clientutils.ResolveCatalogProjectFromInctl(opts.flags)),
		}),
		ProcessReferencedData:   bundleio.ToCatalogReferencedData(ctx, bundleio.WithACClient(opts.acClient)),
	})

	localAssetsDir, err := os.MkdirTemp("", "local-assets")
	if err != nil {
		return nil, fmt.Errorf("could not create temporary directory for local assets: %w", err)
	}
	defer os.RemoveAll(localAssetsDir)

	hwd, err := bundleio.ProcessHardwareDevice(opts.target,
		bundleio.WithProcessAsset(assetInliner.Process),
		bundleio.WithReadOptions(
			bundleio.WithExtractLocalAssetsDir(localAssetsDir),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("could not process HardwareDevice bundle: %w", err)
	}
	var tag atagpb.AssetTag
	if len(hwd.GetMetadata().GetAssetTags()) > 1 {
		return nil, fmt.Errorf("HardwareDevice %q specifies more than one asset tag, but at most one is allowed", idutils.IDFromProtoUnchecked(hwd.GetMetadata().GetId()))
	}
	if len(hwd.GetMetadata().GetAssetTags()) == 1 {
		tag = hwd.GetMetadata().GetAssetTags()[0]
	}

	m := &mpb.Metadata{
		IdVersion: &idpb.IdVersion{
			Id:      hwd.GetMetadata().GetId(),
			Version: opts.flags.GetFlagVersion(),
		},
		AssetType:     atpb.AssetType_ASSET_TYPE_HARDWARE_DEVICE,
		AssetTag:      tag,
		DisplayName:   hwd.GetMetadata().GetDisplayName(),
		Documentation: hwd.GetMetadata().GetDocumentation(),
		Vendor:        hwd.GetMetadata().GetVendor(),
		ReleaseNotes:  opts.flags.GetFlagReleaseNotes(),
	}

	return &acpb.CreateAssetRequest{
		Asset: &acpb.Asset{
			Metadata: m,
			DeploymentData: &acpb.Asset_AssetDeploymentData{
				AssetSpecificDeploymentData: &acpb.Asset_AssetDeploymentData_HardwareDeviceSpecificDeploymentData{
					HardwareDeviceSpecificDeploymentData: &acpb.Asset_HardwareDeviceDeploymentData{
						Manifest: hwd,
					},
				},
			},
			ReleaseMetadata: &rmpb.ReleaseMetadata{
				Default:    opts.flags.GetFlagDefault(),
				OrgPrivate: opts.flags.GetFlagOrgPrivate(),
			},
		},
	}, nil
}

func remoteOpt() remote.Option {
	return remote.WithAuthFromKeychain(google.Keychain)
}
