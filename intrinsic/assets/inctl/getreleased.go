// Copyright 2023 Intrinsic Innovation LLC

// Package getreleased defines the command to get information about a released asset.
package getreleased

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"
)

// GetCommand returns the command to get asset deployment data.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()

	cmd := &cobra.Command{
		Use: "get_released id_version",
		Short: strings.Join([]string{
			"Get information about the specified asset id_version.",
		}, "\n"),
		Example: strings.Join([]string{
			"$ inctl asset get_released some.package.my_skill.0.0.1",
		}, "\n"),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ivp, err := idutils.NewIDVersionParts(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse id_version: %v", err)
			}

			view, err := flags.GetFlagView()
			if err != nil {
				return fmt.Errorf("failed to parse view: %v", err)
			}

			ctx, conn, err := clientutils.DialCatalogFromInctl(cmd, flags)
			if err != nil {
				return fmt.Errorf("failed to create client connection: %v", err)
			}
			defer conn.Close()

			client := acgrpcpb.NewAssetCatalogClient(conn)
			asset, err := client.GetAsset(ctx, &acpb.GetAssetRequest{
				AssetId: &acpb.GetAssetRequest_IdVersion{
					IdVersion: ivp.IDVersionProto(),
				},
				View: view,
			})
			if err != nil {
				return fmt.Errorf("failed to get asset: %v", err)
			}

			prtr, err := printer.NewPrinter(root.FlagOutput)
			if err != nil {
				return err
			}
			prtr.Print(asset)

			return nil
		},
	}

	flags.SetCommand(cmd)
	flags.AddFlagOrganizationOptional()
	flags.AddFlagView()

	return cmd
}
