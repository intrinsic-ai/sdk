// Copyright 2023 Intrinsic Innovation LLC

// Package list defines the skill list command which lists skills in a registry.
package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	skillCmd "intrinsic/skills/tools/skill/cmd/cmd"
	"intrinsic/skills/tools/skill/cmd/listutil"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"

	skillregistrygrpcpb "intrinsic/skills/proto/skill_registry_go_grpc_proto"
	skillregistrypb "intrinsic/skills/proto/skill_registry_go_grpc_proto"
	spb "intrinsic/skills/proto/skills_go_proto"
)

const (
	keyFilter = "filter"

	sideloadedFilter = "sideloaded"
	releasedFilter   = "released"
)

var (
	cmdFlags = cmdutils.NewCmdFlags()

	filterOptions = []string{sideloadedFilter, releasedFilter}
)

type listSkillsParams struct {
	filter   string
	printer  printer.Printer
	pageSize int32 // This can be set in tests to verify pagination behavior.
}

func listSkills(ctx context.Context, client skillregistrygrpcpb.SkillRegistryClient, params *listSkillsParams) error {
	filter := ""
	if params.filter == sideloadedFilter {
		filter = sideloadedFilter
	} else if params.filter == releasedFilter {
		filter = fmt.Sprintf("-%s", sideloadedFilter)
	}

	var (
		skills        []*spb.Skill
		nextPageToken string
	)
	for {
		resp, err := client.ListSkills(ctx, &skillregistrypb.ListSkillsRequest{
			Filter:    filter,
			PageSize:  params.pageSize,
			PageToken: nextPageToken,
		})
		if err != nil {
			return errors.Wrap(err, "could not list skills")
		}
		skills = append(skills, resp.GetSkills()...)
		nextPageToken = resp.GetNextPageToken()
		if nextPageToken == "" {
			break
		}
	}

	params.printer.Print(listutil.SkillDescriptionsFromSkills(skills))

	return nil
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List skills that are loaded into a solution.",
	Example: `List skills of a running solution (solution id, not display name)
$ inctl skill list --project my-project --solution my-solution-id

	To find a running solution's id, run:
	$ inctl solution list --project my-project --filter "running_on_hw,running_in_sim" --output json

Set the cluster on which the solution is running
$ inctl skill list --project my-project --cluster my-cluster
`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		ctx, conn, _, err := clientutils.DialClusterFromInctl(ctx, cmdFlags)
		if err != nil {
			return err
		}
		defer conn.Close()

		prtr, err := printer.NewPrinter(root.FlagOutput)
		if err != nil {
			return err
		}

		client := skillregistrygrpcpb.NewSkillRegistryClient(conn)
		err = listSkills(ctx, client, &listSkillsParams{
			filter:  cmdFlags.GetString(keyFilter),
			printer: prtr,
		})
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	skillCmd.SkillCmd.AddCommand(listCmd)
	cmdFlags.SetCommand(listCmd)

	cmdFlags.AddFlagsAddressClusterSolution()
	cmdFlags.AddFlagsProjectOrg()

	cmdFlags.OptionalString(keyFilter, "", fmt.Sprintf("Filter skills by the way they where loaded into the solution. One of: %s.", strings.Join(filterOptions, ", ")))
}
