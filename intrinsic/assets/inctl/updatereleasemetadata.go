// Copyright 2023 Intrinsic Innovation LLC

// Package updatereleasemetadata defines the command to update the metadata of a released asset.
package updatereleasemetadata

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"

	fmpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	rmpb "intrinsic/assets/catalog/proto/v1/release_metadata_go_proto"
)

// GetCommand returns the command to get asset deployment data.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()

	cmd := &cobra.Command{
		Use: "update_release_metadata id_version",
		Short: strings.Join([]string{
			"Update the release metadata of the specified asset id_version.",
		}, "\n"),
		Example: strings.Join([]string{
			"$ inctl asset update_release_metadata some.package.my_skill.0.0.1 --org_private=false",
		}, "\n"),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ivp, err := idutils.NewIDVersionParts(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse id_version: %v", err)
			}

			ctx, conn, err := clientutils.DialCatalogFromInctl(cmd, flags)
			if err != nil {
				return fmt.Errorf("failed to create client connection: %v", err)
			}
			defer conn.Close()

			rm := &rmpb.ReleaseMetadata{}
			var updateMask []string
			if flags.GetFlagOrgPrivateIsSet() {
				op := flags.GetFlagOrgPrivate()
				if !op && !flags.GetFlagSkipPrompts() {
					fmt.Println("WARNING: org_private cannot be set to true after being set to false.")
					fmt.Println("Do you want to continue? [y/N] ")
					reader := bufio.NewReader(os.Stdin)
					input, _, err := reader.ReadRune()
					if err != nil {
						return fmt.Errorf("could not read response: %v", err)
					}
					if unicode.ToLower(input) != 'y' {
						return fmt.Errorf("aborted")
					}
				}

				rm.OrgPrivate = op
				updateMask = append(updateMask, "org_private")
			}
			if flags.GetFlagDefaultIsSet() {
				rm.Default = flags.GetFlagDefault()
				updateMask = append(updateMask, "default")
			}

			client := acgrpcpb.NewAssetCatalogClient(conn)
			newRM, err := client.UpdateReleaseMetadata(ctx, &acpb.UpdateReleaseMetadataRequest{
				IdVersion:       ivp.IDVersionProto(),
				ReleaseMetadata: rm,
				UpdateMask:      &fmpb.FieldMask{Paths: updateMask},
			})
			if err != nil {
				return fmt.Errorf("failed to update release metadata: %v", err)
			}

			prtr, err := printer.NewPrinter(root.FlagOutput)
			if err != nil {
				return err
			}
			prtr.Print(newRM)

			return nil
		},
	}

	flags.SetCommand(cmd)
	flags.AddFlagDefault("asset")
	flags.AddFlagOrganizationOptional()
	flags.AddFlagOrgPrivate()
	flags.AddFlagSkipPrompts()

	return cmd
}
