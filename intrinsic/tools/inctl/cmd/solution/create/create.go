// Copyright 2023 Intrinsic Innovation LLC

// Package create provides a command to create a solution.
package create

import (
	"context"
	"fmt"

	"intrinsic/assets/cmdutils"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/color"
	"intrinsic/tools/inctl/util/orgutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	branchpb "intrinsic/solution_versions/proto/v1/branch_go_proto"
	snapshotpb "intrinsic/solution_versions/proto/v1/snapshot_go_proto"
	solutionversionservicepb "intrinsic/solution_versions/proto/v1/solution_version_service_go_proto"
)

type createParams struct {
	displayName       string
	snapshotID        string
	sourceBranchID    string
	templateID        string
	commitTitle       string
	commitDescription string
}

func makeCreateVersionBranchRequest(p createParams) *solutionversionservicepb.CreateBranchRequest {
	var commitMessage *snapshotpb.CommitMessage
	if p.commitTitle != "" || p.commitDescription != "" {
		commitMessage = &snapshotpb.CommitMessage{
			Title:       p.commitTitle,
			Description: p.commitDescription,
		}
	}

	req := &solutionversionservicepb.CreateBranchRequest{
		Branch: &branchpb.Branch{
			DisplayName: p.displayName,
			BranchType: branchpb.Branch_BRANCH_TYPE_MAIN,
		},
		CommitMessage: commitMessage,
	}

	if p.snapshotID != "" {
		req.From = &solutionversionservicepb.CreateBranchRequest_SnapshotSource_{
			SnapshotSource: &solutionversionservicepb.CreateBranchRequest_SnapshotSource{
				SnapshotId: p.snapshotID,
				BranchId:   p.sourceBranchID,
			},
		}
	} else {
		templateSource := &solutionversionservicepb.CreateBranchRequest_TemplateSource{}
		if p.templateID != "" {
			templateSource.TemplateId = p.templateID
		}

		req.From = &solutionversionservicepb.CreateBranchRequest_TemplateSource_{
			TemplateSource: templateSource,
		}
	}

	return req
}

func createSolution(ctx context.Context, svsC solutionversionservicepb.SolutionVersionServiceClient, reqVersion *solutionversionservicepb.CreateBranchRequest) error {
	// 1. Create Version Branch
	respVersion, err := svsC.CreateBranch(ctx, reqVersion)
	if err != nil {
		return fmt.Errorf("failed to create version branch: %v. Try again or contact support the issue persists.", err)
	}

	color.C.Blue().Printf("Successfully created version branch %q with ID %q\n", respVersion.GetDisplayName(), respVersion.GetId())

	// 2. Create Deployment Branch
	reqDeployment := &solutionversionservicepb.CreateBranchRequest{
		Branch: &branchpb.Branch{
			DisplayName: reqVersion.GetBranch().GetDisplayName(),
			BranchType:       branchpb.Branch_BRANCH_TYPE_UNSPECIFIED,
			UpstreamBranchId: respVersion.GetId(),
		},
	}

	respDeployment, err := svsC.CreateBranch(ctx, reqDeployment)
	if err != nil {
		return fmt.Errorf("failed to create deployment branch: %v. Try again or contact support if the issue persists.", err)
	}

	color.C.Blue().Printf("Successfully created deployment branch %q with ID %q\n", respDeployment.GetDisplayName(), respDeployment.GetId())
	return nil
}

func NewCommand() *cobra.Command {
	viperLocal := viper.New()
	flags := cmdutils.NewCmdFlagsWithViper(viperLocal)

	var (
		flagDisplayName       string
		flagSnapshotID        string
		flagSourceBranchID    string
		flagTemplateID        string
		flagCommitTitle       string
		flagCommitDescription string
	)

	solutionCreateCmd := orgutil.WrapCmd(&cobra.Command{
		Use:   "create",
		Short: "Create a new solution (version and deployment branches).",
		Long:  `Create a new solution by creating a version branch and a tracking deployment branch. You can create an empty solution, a solution from an existing snapshot or from a template.`,
		Example: `  # Create an empty solution
  inctl solution create --display-name "my-solution"

  # Create a solution from an existing snapshot
  inctl solution create --display-name "my-solution" \
    --snapshot-id <snapshot_id> \
    --source-branch-id <branch_id>

  # Create a solution from a template
  inctl solution create --display-name "my-solution" --template-id <template_id>`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			conn, err := auth.NewCloudConnection(ctx, auth.WithFlagValues(viperLocal))
			if err != nil {
				return err
			}
			defer conn.Close()

			svsC := solutionversionservicepb.NewSolutionVersionServiceClient(conn)
			reqVersion := makeCreateVersionBranchRequest(createParams{
				displayName:       flagDisplayName,
				snapshotID:        flagSnapshotID,
				sourceBranchID:    flagSourceBranchID,
				templateID:        flagTemplateID,
				commitTitle:       flagCommitTitle,
				commitDescription: flagCommitDescription,
			})
			return createSolution(ctx, svsC, reqVersion)
		},
	}, viperLocal)

	flags.SetCommand(solutionCreateCmd)

	solutionCreateCmd.Flags().StringVar(&flagDisplayName, "display-name", "", "The display name of the solution.")
	solutionCreateCmd.Flags().StringVar(&flagSnapshotID, "snapshot-id", "", "The ID of the snapshot to create the version branch from.")
	solutionCreateCmd.Flags().StringVar(&flagSourceBranchID, "source-branch-id", "", "The branch ID that contains the snapshot in its history.")
	solutionCreateCmd.Flags().StringVar(&flagTemplateID, "template-id", "", "The ID of the template to create a version branch from.")
	solutionCreateCmd.Flags().StringVar(&flagCommitTitle, "commit-title", "", "Optional commit title for the initial version.")
	solutionCreateCmd.Flags().StringVar(&flagCommitDescription, "commit-description", "", "Optional commit description for the initial version.")

	solutionCreateCmd.MarkFlagRequired("display-name")
	solutionCreateCmd.MarkFlagsMutuallyExclusive("template-id", "snapshot-id")
	solutionCreateCmd.MarkFlagsMutuallyExclusive("template-id", "source-branch-id")
	solutionCreateCmd.MarkFlagsRequiredTogether("snapshot-id", "source-branch-id")

	return solutionCreateCmd
}
