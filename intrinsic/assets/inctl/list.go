// Copyright 2023 Intrinsic Innovation LLC

// Package list defines the command to list installed assets.
package list

import (
	"fmt"

	"github.com/spf13/cobra"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/inctl/assetviews"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"

	iagrpcpb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
	spb "intrinsic/assets/proto/v1/search_go_proto"
)

const (
	keyOutputType = "output_type"
)

// GetCommand returns the command to list installed assets in a cluster.
func GetCommand(defaultTypes string) *cobra.Command {
	flags := cmdutils.NewCmdFlags()

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed assets",
		Example: `
		List the assets installed in a solution:
		$ inctl asset list --org my_organization --solution my_solution_id

		To find a running solution's id, run:
		$ inctl solution list --project my_project --filter "running_on_hw,running_in_sim" --output json

		Can also use:
		$ inctl asset list --project my_project --address my_address
		or
		$ inctl asset list --project my_project --cluster my_cluster

		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			vt, err := assetviews.AssetTextViewTypeFromString(flags.GetString(keyOutputType))
			if err != nil {
				return err
			}

			var filter *iapb.ListInstalledAssetsRequest_Filter
			if assetTypes, err := flags.GetFlagAssetTypes(); err != nil {
				return err
			} else if len(assetTypes) > 0 {
				filter = &iapb.ListInstalledAssetsRequest_Filter{
					AssetTypes: assetTypes,
				}
			}

			ctx, conn, _, err := clientutils.DialClusterFromInctl(ctx, flags)
			if err != nil {
				return err
			}
			defer conn.Close()

			client := iagrpcpb.NewInstalledAssetsClient(conn)
			prtr, err := printer.NewPrinter(root.FlagOutput)
			if err != nil {
				return err
			}

			var pageToken string
			for {
				resp, err := client.ListInstalledAssets(ctx, &iapb.ListInstalledAssetsRequest{
					StrictFilter: filter,
					OrderBy:      spb.OrderBy_ORDER_BY_ID,
					PageToken:    pageToken,
				})
				if err != nil {
					return fmt.Errorf("could not list assets: %v", err)
				}
				for _, asset := range resp.GetInstalledAssets() {
					prtr.Print(assetviews.FromAsset(asset, assetviews.WithTextViewType(vt)))
				}
				pageToken = resp.GetNextPageToken()
				if pageToken == "" {
					break
				}
			}

			return nil
		},
	}

	flags.SetCommand(cmd)
	flags.AddFlagAssetTypes(defaultTypes)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagsProjectOrg()
	flags.OptionalString(keyOutputType, string(assetviews.AssetTextViewTypeID), fmt.Sprintf("The output type of the list command. One of: %v.", assetviews.AllAssetTextViewTypes))

	return cmd
}
