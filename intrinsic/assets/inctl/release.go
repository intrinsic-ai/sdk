// Copyright 2023 Intrinsic Innovation LLC

// Package release defines the command that releases an asset to the catalog.
package release

import (
	"fmt"

	"intrinsic/assets/catalog/releaseasset"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/imagetransfer"
	"intrinsic/assets/imageutils"
	"intrinsic/skills/tools/skill/cmd/directupload/directupload"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"

	"github.com/spf13/cobra"
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

			var transferer imagetransfer.Transferer
			if true {
				transferer = directupload.NewTransferer(
					directupload.WithDiscovery(directupload.NewCatalogTarget(conn)),
					directupload.WithOutput(cmd.OutOrStdout()),
					directupload.WithFailOver(transferer),
					directupload.WithCatalogOptions(flags.GetFlagImageUploadParallelism()), // this allows uploading images with max size of the single layer of 2GiB.
				)
			}

			return releaseasset.FromBundle(ctx, args[0],
				releaseasset.WithConnection(conn),
				releaseasset.WithDryRun(flags.GetFlagDryRun()),
				releaseasset.WithFlagDefault(flags.GetFlagDefault()),
				releaseasset.WithFlagOrgPrivate(flags.GetFlagOrgPrivate()),
				releaseasset.WithIgnoreExisting(flags.GetFlagIgnoreExisting()),
				releaseasset.WithImageTransferer(transferer),
				releaseasset.WithPrinter(printer.PrintSf),
				releaseasset.WithRegistry(imageutils.GetRegistry(clientutils.ResolveCatalogProjectFromInctl(flags))),
				releaseasset.WithReleaseNotes(flags.GetFlagReleaseNotes()),
				releaseasset.WithVersion(flags.GetFlagVersion()),
			)
		},
	}
	flags.SetCommand(cmd)
	flags.AddFlagDefault("asset")
	flags.AddFlagDryRun()
	flags.AddFlagIgnoreExisting("asset")
	flags.AddFlagImageUploadParallelism(1)
	flags.AddFlagOrganizationOptional()
	flags.AddFlagOrgPrivate()
	flags.AddFlagReleaseNotes("asset")
	flags.AddFlagVersion("asset")

	return cmd
}
