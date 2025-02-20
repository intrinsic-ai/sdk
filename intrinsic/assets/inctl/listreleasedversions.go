// Copyright 2023 Intrinsic Innovation LLC

// Package listreleasedversions defines the list_released_versions command that lists versions of an asset in the catalog.
package listreleasedversions

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"intrinsic/assets/catalog/assetdescriptions"
	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/listutils"
	viewpb "intrinsic/assets/proto/view_go_proto"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"
)

const pageSize int64 = 50

// GetCommand returns a command to list versions of a released asset in the catalog.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{Use: "list_released_versions id",
		Short: "List versions of a released asset in the catalog",
		Args:  cobra.ExactArgs(1), // id
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := clientutils.DialCatalogFromInctl(cmd, flags)
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
			assets, err := listutils.ListAllAssets(cmd.Context(), client, pageSize, viewpb.AssetViewType_ASSET_VIEW_TYPE_VERSIONS, filter)
			if err != nil {
				return errors.Wrap(err, "could not list asset versions")
			}
			ad, err := assetdescriptions.FromCatalogAssets(assets)
			if err != nil {
				return err
			}
			prtr.Print(assetdescriptions.IDVersionsStringView{Descriptions: ad})

			return nil
		},
	}
	flags.SetCommand(cmd)
	flags.AddFlagAssetTypes("")

	return cmd
}
