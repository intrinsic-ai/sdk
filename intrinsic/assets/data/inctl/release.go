// Copyright 2023 Intrinsic Innovation LLC

// Package release defines the command that releases a Data asset to the catalog.
package release

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"intrinsic/assets/bundleio"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"

	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	rmpb "intrinsic/assets/catalog/proto/v1/release_metadata_go_proto"
)

// GetCommand returns command to release Data assets.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()

	cmd := &cobra.Command{
		Use:   "release bundle.tar",
		Short: "Release a Data asset to the catalog.",
		Long: `Release a Data asset to the catalog.

A Data asset can have any proto message as its deployment data.
The bundle.tar file can be created with the intrinsic_data() BUILD rule.`,
		Example: `
	Release a Data asset to the catalog
	$ bazel build --config=intrinsic //path/to:intrinsic_data_target
	$ inctl data release bazel-bin/path/to/intrinsic_data_target.bundle.tar --version=0.0.1
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]

			printer, err := printer.NewPrinter(root.FlagOutput)
			if err != nil {
				return err
			}

			ctx, conn, err := clientutils.DialCatalogFromInctl(cmd, flags)
			if err != nil {
				return fmt.Errorf("could not dial catalog: %w", err)
			}
			defer conn.Close()

			client := acgrpcpb.NewAssetCatalogClient(conn)

			req, err := makeCreateAssetRequest(cmd.Context(), target, client, flags)
			if err != nil {
				return fmt.Errorf("could not create request for Data asset target %q: %w", target, err)
			}

			idVersion := idutils.IDVersionFromProtoUnchecked(req.GetAsset().GetMetadata().GetIdVersion())
			printer.PrintSf("Releasing Data asset %q to the asset catalog", idVersion)

			if flags.GetFlagDryRun() {
				printer.PrintS("Skipping release: dry-run")
				return nil
			}

			return release(ctx, client, req, flags.GetFlagIgnoreExisting(), printer)
		},
	}
	flags.SetCommand(cmd)
	flags.AddFlagDefault("data asset")
	flags.AddFlagDryRun()
	flags.AddFlagIgnoreExisting("data asset")
	flags.AddFlagOrganizationOptional()
	flags.AddFlagOrgPrivate()
	flags.AddFlagReleaseNotes("data asset")
	flags.AddFlagVersion("data asset")

	return cmd
}

func release(ctx context.Context, client acgrpcpb.AssetCatalogClient, req *acpb.CreateAssetRequest, ignoreExisting bool, printer printer.Printer) error {
	if _, err := client.CreateAsset(ctx, req); err != nil {
		if s, ok := status.FromError(err); ok && ignoreExisting && s.Code() == codes.AlreadyExists {
			printer.PrintS("Skipping release: asset already exists in the catalog")
			return nil
		}
		return fmt.Errorf("could not release the data asset: %w", err)
	}
	printer.PrintS("Finished releasing the data asset")
	return nil
}

func makeCreateAssetRequest(ctx context.Context, target string, client acgrpcpb.AssetCatalogClient, flags *cmdutils.CmdFlags) (*acpb.CreateAssetRequest, error) {
	referencedDataProcessor := bundleio.NoOpReferencedData()
	if !flags.GetFlagDryRun() {
		referencedDataProcessor = bundleio.ToCatalogReferencedData(ctx, bundleio.WithACClient(client))
	}
	da, err := bundleio.ReadDataAsset(target, bundleio.WithProcessReferencedData(referencedDataProcessor))
	if err != nil {
		return nil, fmt.Errorf("could not read Data asset: %w", err)
	}
	da.Metadata.IdVersion.Version = flags.GetFlagVersion()
	da.Metadata.ReleaseNotes = flags.GetFlagReleaseNotes()

	return &acpb.CreateAssetRequest{
		Asset: &acpb.Asset{
			Metadata: da.GetMetadata(),
			ReleaseMetadata: &rmpb.ReleaseMetadata{
				Default:    flags.GetFlagDefault(),
				OrgPrivate: flags.GetFlagOrgPrivate(),
			},
			DeploymentData: &acpb.Asset_AssetDeploymentData{
				AssetSpecificDeploymentData: &acpb.Asset_AssetDeploymentData_DataSpecificDeploymentData{
					DataSpecificDeploymentData: &acpb.Asset_DataDeploymentData{
						Data: da,
					},
				},
			},
		},
	}, nil
}
