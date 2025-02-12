// Copyright 2023 Intrinsic Innovation LLC

// Package list defines the command to list installed assets.
package list

import (
	"fmt"

	"github.com/spf13/cobra"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
	spb "intrinsic/assets/proto/v1/search_go_proto"
)

type outputType string

const (
	keyOutputType = "output_type"

	outputTypeID        = outputType("id")
	outputTypeIDVersion = outputType("id_version")
)

var allOutputTypes = []outputType{outputTypeID, outputTypeIDVersion}

// GetCommand returns the command to list installed assets in a cluster.
func GetCommand() *cobra.Command {
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

			var filter *iapb.ListInstalledAssetsRequest_Filter
			if assetTypes, err := flags.GetFlagAssetTypes(); err != nil {
				return err
			} else if len(assetTypes) > 0 {
				filter = &iapb.ListInstalledAssetsRequest_Filter{
					AssetTypes: assetTypes,
				}
			}

			var outputFrom func(*iapb.InstalledAsset) string
			switch outputType := outputType(flags.GetString(keyOutputType)); outputType {
			case outputTypeID:
				outputFrom = func(asset *iapb.InstalledAsset) string {
					return idutils.IDFromProtoUnchecked(asset.GetMetadata().GetIdVersion().GetId())
				}
			case outputTypeIDVersion:
				outputFrom = func(asset *iapb.InstalledAsset) string {
					return idutils.IDVersionFromProtoUnchecked(asset.GetMetadata().GetIdVersion())
				}
			default:
				return fmt.Errorf("invalid output type: %q. must be one of: %v", outputType, allOutputTypes)
			}

			ctx, conn, _, err := clientutils.DialClusterFromInctl(ctx, flags)
			if err != nil {
				return err
			}
			defer conn.Close()

			var pageToken string
			for {
				client := iagrpcpb.NewInstalledAssetsClient(conn)
				resp, err := client.ListInstalledAssets(ctx, &iapb.ListInstalledAssetsRequest{
					StrictFilter: filter,
					OrderBy:      spb.OrderBy_ORDER_BY_ID,
					PageToken:    pageToken,
				})
				if err != nil {
					return fmt.Errorf("could not list assets: %v", err)
				}
				for _, asset := range resp.GetInstalledAssets() {
					fmt.Println(outputFrom(asset))
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
	flags.AddFlagAssetTypes()
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagsProjectOrg()
	flags.OptionalString(keyOutputType, "id", fmt.Sprintf("The output type of the list command. One of: %v.", allOutputTypes))

	return cmd
}
