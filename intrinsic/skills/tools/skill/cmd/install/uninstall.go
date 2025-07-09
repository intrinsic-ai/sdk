// Copyright 2023 Intrinsic Innovation LLC

// Package uninstall defines the skill command which uninstalls a skill.
package uninstall

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	"intrinsic/assets/imageutils"
	"intrinsic/skills/tools/skill/cmd/cmd"
	"intrinsic/skills/tools/skill/cmd/skillio"

	lrogrpcpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
)

var cmdFlags = cmdutils.NewCmdFlags()

var uninstallCmd = &cobra.Command{
	Use:   "uninstall --type=TYPE TARGET",
	Short: "Remove a skill",
	Example: `Stop a running skill using its build target
$ inctl skill uninstall --type=build //abc:skill_bundle --context=minikube

Stop a running skill by specifying its id
$ inctl skill uninstall --type=id com.foo.skill
`,
	Args: cobra.ExactArgs(1),
	Aliases: []string{
		"stop",
		"unload",
	},
	RunE: func(command *cobra.Command, args []string) error {
		ctx := command.Context()
		target := args[0]

		targetType := imageutils.TargetType(cmdFlags.GetFlagSideloadStopType())
		if targetType != imageutils.Build && targetType != imageutils.ID {
			return fmt.Errorf("type must be one of (%s, %s)", imageutils.Build, imageutils.ID)
		}

		var skillID string
		var err error
		switch targetType {
		case imageutils.Build:
			skillID, err = skillio.SkillIDFromBuildTarget(target)
			if err != nil {
				return err
			}
		case imageutils.ID:
			skillID = target
		default:
			return fmt.Errorf("unimplemented target type: %v", targetType)
		}
		asset, err := idutils.NewIDProto(skillID)
		if err != nil {
			return fmt.Errorf("invalid id: %w", err)
		}

		ctx, conn, _, err := clientutils.DialClusterFromInctl(ctx, cmdFlags)
		if err != nil {
			return err
		}
		defer conn.Close()

		log.Printf("Uninstalling skill %q", skillID)
		client := iagrpcpb.NewInstalledAssetsClient(conn)
		op, err := client.DeleteInstalledAsset(ctx, &iapb.DeleteInstalledAssetRequest{
			Asset: asset,
		})
		if err != nil {
			return fmt.Errorf("could not install the skill: %v", err)
		}
		log.Printf("Awaiting completion of the uninstall")
		lroClient := lrogrpcpb.NewOperationsClient(conn)
		for !op.GetDone() {
			op, err = lroClient.WaitOperation(ctx, &lropb.WaitOperationRequest{
				Name: op.GetName(),
			})
			if err != nil {
				return fmt.Errorf("unable to check status of uninstall: %v", err)
			}
		}

		if err := status.ErrorProto(op.GetError()); err != nil {
			return fmt.Errorf("uninstall failed: %w", err)
		}
		log.Print("Finished uninstalling the skill")

		return nil
	},
}

func init() {
	cmd.SkillCmd.AddCommand(uninstallCmd)
	cmdFlags.SetCommand(uninstallCmd)

	cmdFlags.AddFlagsAddressClusterSolution()
	cmdFlags.AddFlagsProjectOrg()
	cmdFlags.AddFlagSideloadStopType("skill")
}
