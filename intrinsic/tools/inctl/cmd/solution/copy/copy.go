// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package copy provides the inctl solution copy command.
package copy

import (
	"bytes"
	"cmp"
	"context"
	"fmt"
	"io"
	"text/tabwriter"

	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/color"
	"intrinsic/tools/inctl/util/orgutil"
	"intrinsic/tools/inctl/util/promptutil"
	"intrinsic/util/status/extstatus"

	longrunningpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"

	solutioncopyservicepb "intrinsic/solution_copy/proto/v1/solution_copy_service_go_proto"
	snapshotpb "intrinsic/solution_versions/proto/v1/snapshot_go_proto"
)

var (
	flagBranchID   string
	flagSnapshotID string
	flagTargetOrg  string

	flagTargetBranchID          string
	flagTargetBranchDisplayName string
	flagTitle                   string
	flagDescription             string
	flagYes                     bool
)

// newCloudConnection opens the gRPC connection to the SolutionCopyService. Tests replace it with a
// connection to a local fake server, because the default needs cloud credentials.
var newCloudConnection = auth.NewCloudConnection

type previewParams struct {
	branchID   string
	snapshotID string
	targetOrg  string
}

// formatPermissionState returns a visual status symbol and human-readable explanation for a permission state.
func formatPermissionState(state solutioncopyservicepb.PermissionState) (symbol, description string) {
	switch state {
	case solutioncopyservicepb.PermissionState_PERMISSION_STATE_UNSPECIFIED:
		return "?", "Not specified"
	case solutioncopyservicepb.PermissionState_PERMISSION_STATE_GRANTED:
		return "✓", "-"
	case solutioncopyservicepb.PermissionState_PERMISSION_STATE_MISSING:
		return "✗", "Missing access to asset in the target organization. Request access from the asset owner."
	case solutioncopyservicepb.PermissionState_PERMISSION_STATE_ORG_PRIVATE:
		return "✗", "Source organization marked this asset as private. Contact source organization admin to request access."
	case solutioncopyservicepb.PermissionState_PERMISSION_STATE_NOT_FOUND:
		return "✗", "Asset not found. Verify that the asset exists and has not been deleted."
	default:
		return "?", fmt.Sprintf("Unknown state: %d", state)
	}
}

func runPreview(ctx context.Context, scsc solutioncopyservicepb.SolutionCopyServiceClient, params previewParams) (*solutioncopyservicepb.PreviewCopyResponse, error) {
	req := solutioncopyservicepb.PreviewCopyRequest_builder{
		BranchId:   proto.String(params.branchID),
		SnapshotId: proto.String(params.snapshotID),
		TargetOrg:  proto.String(params.targetOrg),
	}.Build()

	resp, err := scsc.PreviewCopy(ctx, req)
	if err != nil {
		if ext, ok := extstatus.FromGRPCError(err); ok {
			return nil, fmt.Errorf("preview copy failed: %w\n%v", err, ext)
		}
		return nil, fmt.Errorf("preview copy failed: %w", err)
	}

	return resp, nil
}

// printPreview renders the preview summary and the asset tables to out. snapshotID is the snapshot
// that the preview resolved, which is the branch tip if params has no snapshot ID.
func printPreview(out io.Writer, params previewParams, snapshotID string, res *solutioncopyservicepb.PreviewCopyResponse) {
	snapshot := cmp.Or(snapshotID, "latest (tip of branch)")
	if snapshotID != "" && params.snapshotID == "" {
		snapshot += " (latest)"
	}
	color.C.Yellow().Fprintf(out, "\nPreview Copy summary:\n")
	color.C.Magenta().Fprintf(out, "* Solution:\t%s\n", params.branchID)
	color.C.Magenta().Fprintf(out, "* Snapshot:\t%s\n", snapshot)
	color.C.Magenta().Fprintf(out, "* Target organization:\t%s\n\n", params.targetOrg)

	if res.GetAssetCount() > 0 {
		color.C.Yellow().Fprintf(out, "The source solution contains %d catalog assets:\n", res.GetAssetCount())
		buf := new(bytes.Buffer)
		w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
		color.C.White().Fprintf(w, "Package name\tName\tVersion\tStatus\tDetails")
		fmt.Fprintln(w)

		for _, aps := range res.GetAssetPermissionStates() {
			symbol, desc := formatPermissionState(aps.GetState())
			format := "%s\t%s\t%s\t%s\t%s"
			pkgName := aps.GetIdVersion().GetId().GetPackage()
			name := aps.GetIdVersion().GetId().GetName()
			version := aps.GetIdVersion().GetVersion()

			switch aps.GetState() {
			case solutioncopyservicepb.PermissionState_PERMISSION_STATE_GRANTED:
				color.C.Green().Fprintf(w, format, pkgName, name, version, symbol, desc)
			default:
				color.C.Red().Fprintf(w, format, pkgName, name, version, symbol, desc)
			}
			fmt.Fprintln(w)
		}
		_ = w.Flush()
		fmt.Fprintln(out, buf.String())
	} else {
		color.C.Yellow().Fprintf(out, "* No catalog assets will be used\n")
	}

	if len(res.GetSideloadedAssets()) > 0 {
		color.C.Yellow().Fprintf(out, "The source solution contains %d local assets\n", len(res.GetSideloadedAssets()))
		buf := new(bytes.Buffer)
		w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
		color.C.White().Fprintf(w, "Package name\tName")
		fmt.Fprintln(w)
		for _, sa := range res.GetSideloadedAssets() {
			color.C.Magenta().Fprintf(w, "%s\t%s", sa.GetId().GetPackage(), sa.GetId().GetName())
			fmt.Fprintln(w)
		}
		_ = w.Flush()
		fmt.Fprintln(out, buf.String())
	} else {
		color.C.Yellow().Fprintf(out, "* The source solution contains no local assets\n")
	}
}

// checkPermissions returns an error if any catalog asset in res lacks copy permission in the
// target organization.
func checkPermissions(res *solutioncopyservicepb.PreviewCopyResponse, targetOrg string) error {
	var assetsWithIssues int
	for _, aps := range res.GetAssetPermissionStates() {
		if aps.GetState() != solutioncopyservicepb.PermissionState_PERMISSION_STATE_GRANTED {
			assetsWithIssues++
		}
	}
	if assetsWithIssues > 0 {
		return fmt.Errorf("copy aborted: %d asset(s) do not have copy permissions granted in target organization %q; see the Details column above to resolve each one", assetsWithIssues, targetOrg)
	}
	return nil
}

// confirmCopy asks the user to confirm the copy of snapshotID. isTip reports whether snapshotID is
// the branch tip. It returns false on any answer other than yes, including when stdin is not
// interactive.
func confirmCopy(cmd *cobra.Command, snapshotID string, isTip bool) (bool, error) {
	target := fmt.Sprintf("as new solution %q", flagTargetBranchDisplayName)
	if flagTargetBranchID != "" {
		target = fmt.Sprintf("as a new version of existing solution %q", flagTargetBranchID)
	}
	snapshot := "the latest snapshot"
	if snapshotID != "" {
		snapshot = fmt.Sprintf("snapshot %q", snapshotID)
		if isTip {
			snapshot += " (latest)"
		}
	}
	prompt := fmt.Sprintf("\nCopy %s of solution %q to organization %q %s?", snapshot, flagBranchID, flagTargetOrg, target)
	return promptutil.PromptYesNo(cmd, prompt, promptutil.DefaultNo, promptutil.ReemitPromptOnInvalidInput)
}

func runCopy(ctx context.Context, scsc solutioncopyservicepb.SolutionCopyServiceClient, snapshotID string) (*longrunningpb.Operation, error) {
	b := solutioncopyservicepb.CopyRequest_builder{
		BranchId:   proto.String(flagBranchID),
		SnapshotId: proto.String(snapshotID),
		TargetOrg:  proto.String(flagTargetOrg),
	}
	if flagTargetBranchID != "" {
		b.TargetBranchId = proto.String(flagTargetBranchID)
	} else {
		b.TargetBranchDisplayName = proto.String(flagTargetBranchDisplayName)
	}
	// Leave the commit message unset when neither flag is given, so that the service generates one.
	if flagTitle != "" || flagDescription != "" {
		b.CommitMessage = &snapshotpb.CommitMessage{Title: flagTitle, Description: flagDescription}
	}

	op, err := scsc.Copy(ctx, b.Build())
	if err != nil {
		if ext, ok := extstatus.FromGRPCError(err); ok {
			return nil, fmt.Errorf("copy failed: %w\n%v", err, ext)
		}
		return nil, fmt.Errorf("copy failed: %w", err)
	}
	return op, nil
}

// NewCommand returns the copy command.
func NewCommand() *cobra.Command {
	viperLocal := viper.New()

	solutionCopyCmd := orgutil.WrapCmd(&cobra.Command{
		Use:   "copy",
		Short: "Copy a solution to a different organization",
		Long: `Copy a solution to a different organization.

This feature enables users to copy solutions from one organization to another.

Before triggering the copy operation of a solution, the user is presented with a preview of this operation,
detailing which assets are sideloaded, which catalog assets are referenced in the solution,
and to which assets the target organization has access.

The branch must be a version branch, such as the version branch ID printed by "inctl solution create".
Solution IDs listed by "inctl solution list" are deployment branches and cannot be copied directly.

If --snapshot-id is omitted, the tip snapshot of the branch is used. The tip is resolved once, during
the preview, and that snapshot is the one copied.

The copy either creates a new solution in the target organization (--target-branch-display-name),
or adds a new version to an existing solution of the target organization (--target-branch-id).
--target-branch-id must be the solution's version branch, ending in _MAIN, such as the version
branch ID printed by "inctl solution create". Solution IDs listed by "inctl solution list" cannot be used.

The copy only starts after you confirm it. Pass --yes to skip the confirmation, for example in scripts.

During the copy, the service ensures that members of the target organization have access to all assets
of the solution and copies all required data to the target organization.`,
		Example: `  # Preview and copy the latest snapshot (branch tip) of a solution
  inctl solution copy \
    --org "source_org_id" \
    --branch-id "f47ac10b-58cc-4372-a567-0e02b2c3d479_MAIN" \
    --target-org "destination_org_id" \
    --target-branch-display-name "My copied solution"

  # Preview a solution snapshot copy, confirm it, and start the copy
  inctl solution copy \
    --org "source_org_id" \
    --branch-id "f47ac10b-58cc-4372-a567-0e02b2c3d479_MAIN" \
    --snapshot-id "c56a4180-65aa-42ec-a945-5fd21dec0538_SNAP" \
    --target-org "destination_org_id" \
    --target-branch-display-name "My copied solution"

  # Start the copy without confirmation, with a custom commit message
  inctl solution copy \
    --org "source_org_id" \
    --branch-id "f47ac10b-58cc-4372-a567-0e02b2c3d479_MAIN" \
    --snapshot-id "c56a4180-65aa-42ec-a945-5fd21dec0538_SNAP" \
    --target-org "destination_org_id" \
    --target-branch-display-name "My copied solution" \
    --title "Copy of v1.0" \
    --description "Initial copy for the customer organization" \
    --yes

  # Copy a snapshot into an existing solution of the target organization, as a new version
  inctl solution copy \
    --org "source_org_id" \
    --branch-id "f47ac10b-58cc-4372-a567-0e02b2c3d479_MAIN" \
    --snapshot-id "c56a4180-65aa-42ec-a945-5fd21dec0538_SNAP" \
    --target-org "destination_org_id" \
    --target-branch-id "0b5e4a6c-1f2d-4e3a-9c8b-7d6e5f4a3b2c_MAIN"`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// The service replaces the whole commit message when the title is empty (b/507928991), so a
			// description without a title would be silently lost.
			if flagDescription != "" && flagTitle == "" {
				return descriptionWithoutTitleError()
			}
			if flagTargetBranchID == "" && flagTargetBranchDisplayName == "" {
				return emptyTargetError()
			}

			ctx := cmd.Context()
			conn, err := newCloudConnection(ctx, auth.WithFlagValues(viperLocal))
			if err != nil {
				return err
			}
			defer conn.Close()

			params := previewParams{
				branchID:   flagBranchID,
				snapshotID: flagSnapshotID,
				targetOrg:  flagTargetOrg,
			}

			scsc := solutioncopyservicepb.NewSolutionCopyServiceClient(conn)
			res, err := runPreview(ctx, scsc, params)
			if err != nil {
				return err
			}

			// Copy the snapshot that the preview resolved, so that a tip that moves after the preview
			// does not change what is copied. Services that do not return the snapshot ID leave it
			// empty; fall back to --snapshot-id then.
			snapshotID := cmp.Or(res.GetSnapshotId(), params.snapshotID)
			printPreview(cmd.OutOrStdout(), params, snapshotID, res)
			if err := checkPermissions(res, params.targetOrg); err != nil {
				return err
			}

			if !flagYes {
				ok, err := confirmCopy(cmd, snapshotID, params.snapshotID == "")
				if err != nil {
					return err
				}
				if !ok {
					return copyNotConfirmedError()
				}
			}

			op, err := runCopy(ctx, scsc, snapshotID)
			if err != nil {
				return err
			}
			color.C.Green().Fprintf(cmd.OutOrStdout(), "\nCopy started. Operation: %s\n", op.GetName())
			return nil
		},
	}, viperLocal)

	solutionCopyCmd.Flags().StringVar(&flagBranchID, "branch-id", "", "ID of the version branch to copy (not the solution ID listed by 'inctl solution list').")
	solutionCopyCmd.Flags().StringVar(&flagSnapshotID, "snapshot-id", "", "ID of the snapshot to copy. If not set, the tip snapshot of the branch, as resolved by the preview, is used.")
	solutionCopyCmd.Flags().StringVar(&flagTargetOrg, "target-org", "", "ID of the organization to copy the solution to.")
	solutionCopyCmd.Flags().StringVar(&flagTargetBranchDisplayName, "target-branch-display-name", "", "Display name of the new solution created in the target organization (at most 80 characters). Mutually exclusive with --target-branch-id.")
	solutionCopyCmd.Flags().StringVar(&flagTargetBranchID, "target-branch-id", "", "ID of the existing solution in the target organization to add the copy to as a new version: its version branch, ending in _MAIN (not the solution ID listed by 'inctl solution list'). Mutually exclusive with --target-branch-display-name.")
	solutionCopyCmd.Flags().StringVar(&flagTitle, "title", "", "Title of the commit created in the target organization (at most 80 characters). Generated by the service if empty.")
	solutionCopyCmd.Flags().StringVar(&flagDescription, "description", "", "Description of the commit created in the target organization. Requires --title.")
	solutionCopyCmd.Flags().BoolVar(&flagYes, "yes", false, "Start the copy without asking for confirmation.")

	solutionCopyCmd.MarkFlagRequired("branch-id")
	solutionCopyCmd.MarkFlagRequired("target-org")
	solutionCopyCmd.MarkFlagsMutuallyExclusive("target-branch-display-name", "target-branch-id")
	solutionCopyCmd.MarkFlagsOneRequired("target-branch-display-name", "target-branch-id")

	return solutionCopyCmd
}
