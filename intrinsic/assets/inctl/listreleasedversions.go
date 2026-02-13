// Copyright 2023 Intrinsic Innovation LLC

// Package listreleasedversions defines the list_released_versions command that lists versions of an asset in the catalog.
package listreleasedversions

import (
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/inctl/assetviews"
	"intrinsic/assets/listutils"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"

	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	viewpb "intrinsic/assets/proto/view_go_proto"
)

const pageSize int64 = 50

// GetCommand returns a command to list versions of a released asset in the catalog.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "list_released_versions id",
		Short: "List versions of a released asset in the catalog",
		Args:  cobra.ExactArgs(1), // id
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, conn, err := clientutils.DialCatalogFromInctl(cmd, flags)
			if err != nil {
				return errors.Wrap(err, "failed to create client connection")
			}
			defer conn.Close()

			client := acgrpcpb.NewAssetCatalogClient(conn)
			prtr, err := printer.NewPrinter(root.FlagOutput)
			if err != nil {
				return err
			}

			assetTypes, err := flags.GetFlagAssetTypes()
			if err != nil {
				return err
			}

			filter := &acpb.ListAssetsRequest_AssetFilter{
				Id:         proto.String(args[0]),
				AssetTypes: assetTypes,
			}
			assets, err := listutils.ListAllAssets(ctx, client, pageSize, viewpb.AssetViewType_ASSET_VIEW_TYPE_VERSIONS, filter)
			if err != nil {
				return errors.Wrap(err, "could not list asset versions")
			}
			for _, asset := range assets {
				prtr.Print(assetviews.FromAsset(asset))
			}

			return nil
		},
	}
	flags.SetCommand(cmd)
	flags.AddFlagAssetTypes("")
	flags.AddFlagOrganizationOptional()

	return cmd
}
