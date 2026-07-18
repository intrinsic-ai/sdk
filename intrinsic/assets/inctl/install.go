// Copyright 2023 Intrinsic Innovation LLC

// Package install defines the command to install an asset.
package install

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"intrinsic/assets/bundle"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	"intrinsic/assets/imagetransfer"
	"intrinsic/assets/referenceddata"
	"intrinsic/assets/services/bundleimages"
	"intrinsic/kubernetes/acl/identity"
	"intrinsic/skills/tools/skill/cmd/directupload/directupload"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"

	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_proto"
	assetartifactspb "intrinsic/assets/proto/v1/asset_artifacts_go_proto"
	rpb "intrinsic/assets/proto/v1/reference_go_proto"

	lrogrpcpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
)

// solutionAssetRegex matches strings of the form <branch_id>/<asset_id> that indicate an Asset in
// a specific Solution branch.
//
// Further validation is needed on both parts to ensure validity.
var solutionAssetRegex = regexp.MustCompile(`^(?P<branch>[A-Za-z0-9_\-]+)/(?P<asset>[a-z0-9_\.]+)$`)

const (
)

// GetCommand returns a command to install an asset.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "install <asset_id_version>",
		Short: "Install an Asset",
		Example: `
  Install a local Asset bundle into the specified Solution:
  $ inctl asset install abc/bundle.tar \
      --org $my_org \
      --solution $my_solution_id

  Install an Asset from the catalog into the specified Solution:
  $ inctl asset install ai.intrinsic.calculator_service.0.20260126.0-RC00 \
      --org $my_org \
      --solution $my_solution_id

  Install an Asset from another Solution into the specified Solution:
  $ inctl asset install $source_solution_id/ai.intrinsic.calculator_service \
      --org $my_org \
      --solution $my_solution_id

  To find a running Solution's id, run:
  $ inctl solution list --org $my_org --filter "running_on_hw,running_in_sim" --output json

  The Asset can also be installed by specifying the cluster on which the Solution is running:
  $ inctl asset install $my_asset --org $my_org --cluster $my_cluster
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			target := args[0]

			policy, err := flags.GetFlagPolicy()
			if err != nil {
				return err
			}

			ctx, err = identity.AppendToOutgoingContext(ctx, identity.WithOrg(flags.GetFlagOrganization()))
			if err != nil {
				return fmt.Errorf("failed to add org information to context: %w", err)
			}
			ctx, conn, address, err := clientutils.DialClusterFromInctl(ctx, flags)
			if err != nil {
				return err
			}
			defer conn.Close()

			// Determine the image transferer to use. Default to direct injection into the cluster.
			var transfer imagetransfer.Transferer
			if registry := flags.GetFlagRegistry(); registry != "" {
				user, pwd := flags.GetFlagsRegistryAuthUserPassword()
				transfer = imagetransfer.RemoteTransferer(registry, user, pwd)
			}
			if !flags.GetFlagSkipDirectUpload() {
				transfer = directupload.NewTransferer(
					directupload.WithDiscovery(directupload.NewFromConnection(conn)),
					directupload.WithOutput(cmd.OutOrStdout()),
					directupload.WithFailOver(transfer),
				)
			}
			if transfer == nil {
				return fmt.Errorf("--registry must be specified if --skip-direct-upload is used")
			}
			client := iagrpcpb.NewInstalledAssetsClient(conn)
			authCtx := clientutils.AuthInsecureConn(ctx, address, flags.GetFlagProject())

			processor := &bundle.Processor{
				ImageProcessor: bundleimages.CreateImageProcessor(transfer),
				ReferencedDataProcessor: referenceddata.NewProcessor(
					assetartifactspb.NewAssetArtifactsClient(conn),
					lropb.NewOperationsClient(conn),
					referenceddata.WithProgressWriter(cmd.OutOrStdout()),
				),
			}

			asset, err := assetFromTarget(authCtx, target, processor.ProcessFile)
			if err != nil {
				return err
			}

			op, err := client.CreateInstalledAsset(authCtx, &iapb.CreateInstalledAssetRequest{
				Policy: policy,
				Asset:  asset,
			})
			if err != nil {
				return fmt.Errorf("could not install the asset: %v", err)
			}

			log.Printf("Awaiting completion of the installation")
			lroClient := lrogrpcpb.NewOperationsClient(conn)
			for !op.GetDone() {
				op, err = lroClient.WaitOperation(ctx, &lropb.WaitOperationRequest{
					Name: op.GetName(),
				})
				if err != nil {
					return fmt.Errorf("unable to check status of installation: %v", err)
				}
			}

			if err := status.ErrorProto(op.GetError()); err != nil {
				return fmt.Errorf("installation failed: %w", err)
			}
			installed := &iapb.InstalledAsset{}
			if err := op.GetResponse().UnmarshalTo(installed); err != nil {
				return fmt.Errorf("unable to parse result from successful installation: %w", err)
			}
			log.Printf("Finished installing %q", idutils.IDVersionFromProtoUnchecked(installed.GetMetadata().GetIdVersion()))
			return nil
		},
	}

	flags.SetCommand(cmd)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagPolicy("asset")
	flags.AddFlagsProjectOrg()
	flags.AddFlagRegistry()
	flags.AddFlagsRegistryAuthUserPassword()
	flags.AddFlagSkipDirectUpload("asset")

	return cmd
}

type processBundle func(ctx context.Context, path string) (bundle.ProcessedBundle, error)

func assetFromTarget(ctx context.Context, target string, process processBundle) (*iapb.CreateInstalledAssetRequest_Asset, error) {
	fileExists := false
	if _, err := os.Stat(target); err == nil {
		fileExists = true
	}

	idvParts, err := idutils.NewIDVersionParts(target)
	isIDVersion := err == nil

	solutionAssetMatches := solutionAssetRegex.FindStringSubmatch(target)
	isSolutionAsset := false
	if solutionAssetMatches != nil {
		isSolutionAsset = idutils.IsID(solutionAssetMatches[solutionAssetRegex.SubexpIndex("asset")])
	}

	if (isIDVersion || isSolutionAsset) && fileExists {
		return nil, fmt.Errorf("input is ambiguous; %q is both a file and an id_version or Solution Asset", target)
	}

	if isIDVersion {
		return &iapb.CreateInstalledAssetRequest_Asset{
			Variant: &iapb.CreateInstalledAssetRequest_Asset_Catalog{
				Catalog: idvParts.IDVersionProto(),
			},
		}, nil
	}

	if isSolutionAsset {
		return &iapb.CreateInstalledAssetRequest_Asset{
			Variant: &iapb.CreateInstalledAssetRequest_Asset_SolutionAsset{
				SolutionAsset: &rpb.SolutionAsset{
					BranchId: solutionAssetMatches[solutionAssetRegex.SubexpIndex("branch")],
					Name:     solutionAssetMatches[solutionAssetRegex.SubexpIndex("asset")],
				},
			},
		}, nil
	}

	if fileExists {
		processedBundle, err := process(ctx, target)
		if err != nil {
			return nil, fmt.Errorf("unable to process file: %w", err)
		}
		return processedBundle.Install(), nil
	}

	return nil, fmt.Errorf("%q is not a file, id_version, or Solution Asset; check that the input is formatted correctly", target)
}
