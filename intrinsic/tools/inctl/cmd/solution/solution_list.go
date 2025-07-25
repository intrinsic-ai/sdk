// Copyright 2023 Intrinsic Innovation LLC

package solution

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"intrinsic/skills/tools/skill/cmd/dialerutil"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/orgutil"
	"intrinsic/tools/inctl/util/printer"

	clusterdiscoverypb "intrinsic/frontend/cloud/api/v1/clusterdiscovery_api_go_grpc_proto"
	solutiondiscoverygrpcpb "intrinsic/frontend/cloud/api/v1/solutiondiscovery_api_go_grpc_proto"
	solutiondiscoverypb "intrinsic/frontend/cloud/api/v1/solutiondiscovery_api_go_grpc_proto"
)

var (
	flagFilter     []string
	allowedFilters = []string{"not_running", "running_in_sim", "running_on_hw"}
)

type listSolutionsParams struct {
	filter                    []string
	printer                   printer.Printer
	solutionVersioningEnabled bool
}

// ListSolutionDescriptionsResponse embeds solutiondiscoverypb.ListSolutionDescriptionsResponse.
type ListSolutionDescriptionsResponse struct {
	m *solutiondiscoverypb.ListSolutionDescriptionsResponse
}

// MarshalJSON converts a ListSolutionDescriptionsResponse to a byte slice.
func (res *ListSolutionDescriptionsResponse) MarshalJSON() ([]byte, error) {
	type solution struct {
		Name        string `json:"name,omitempty"`
		State       string `json:"state,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		ClusterName string `json:"clusterName,omitempty"`
		Version     string `json:"version,omitempty"`
	}
	solutions := make([]solution, len(res.m.GetSolutions()))
	for i, c := range res.m.GetSolutions() {
		solutions[i] = solution{
			Name:        c.GetName(),
			State:       c.GetState().String(),
			DisplayName: c.GetDisplayName(),
			ClusterName: c.GetClusterName(),
			Version:     c.GetVersion(),
		}
	}
	return json.Marshal(struct {
		// solution intentionally not omitted when empty
		Solutions []solution `json:"solutions"`
	}{Solutions: solutions})
}

// String converts a ListSolutionDescriptionsResponse to a string
func (res *ListSolutionDescriptionsResponse) String() string {
	const formatString = "%-50s %-15s %-50s %-50s"
	lines := []string{
		fmt.Sprintf(formatString, "Name", "State", "ID", "Version"),
	}
	for _, c := range res.m.GetSolutions() {
		name := c.GetDisplayName()
		if name == "" {
			name = c.GetName()
		}

		statusStr := strings.TrimPrefix(c.GetState().String(), "SOLUTION_STATE_")
		if c.GetClusterName() != "" {
			statusStr = fmt.Sprintf("%s on %s", statusStr, c.GetClusterName())
		}

		versionStr := "Legacy"
		if c.GetVersion() != "" {
			versionStr = c.GetVersion()
		}

		lines = append(
			lines,
			fmt.Sprintf(formatString, name, statusStr, c.GetName(), versionStr))
	}
	return strings.Join(lines, "\n")
}

func validateAndGetFilters(filterNames []string) ([]clusterdiscoverypb.SolutionState, error) {
	filters := []clusterdiscoverypb.SolutionState{}

	if len(filterNames) == 0 {
		return filters, nil
	}

	for _, filterName := range filterNames {
		filter, ok := clusterdiscoverypb.SolutionState_value["SOLUTION_STATE_"+strings.ToUpper(filterName)]
		if !ok {
			return filters,
				fmt.Errorf("Filter needs to be one of %s but is %s",
					strings.Join(allowedFilters, ", "), filterName)
		}
		filters = append(filters, clusterdiscoverypb.SolutionState(filter))
	}

	return filters, nil

}

func listSolutions(ctx context.Context, conn *grpc.ClientConn, params *listSolutionsParams) error {
	filters, err := validateAndGetFilters(params.filter)
	if err != nil {
		return err
	}

	listSolutionsRequest := &solutiondiscoverypb.ListSolutionDescriptionsRequest{Filters: filters}
	client := solutiondiscoverygrpcpb.NewSolutionDiscoveryServiceClient(conn)
	resp, err := client.ListSolutionDescriptions(
		ctx, listSolutionsRequest)

	if err != nil {
		return fmt.Errorf("request to list solutions failed: %w", err)
	}

	params.printer.Print(&ListSolutionDescriptionsResponse{m: resp})
	return nil
}

var solutionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List solutions in a project",
	Long:  "List solutions on the given project.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		prtr, err := printer.NewPrinter(root.FlagOutput)
		if err != nil {
			return err
		}

		projectName := viperLocal.GetString(orgutil.KeyProject)
		orgName := viperLocal.GetString(orgutil.KeyOrganization)
		ctx, conn, err := dialerutil.DialConnectionCtx(cmd.Context(), dialerutil.DialInfoParams{
			CredName: projectName,
			CredOrg:  orgName,
		})
		if err != nil {
			return fmt.Errorf("failed to create client connection: %w", err)
		}
		defer conn.Close()

		err = listSolutions(ctx, conn, &listSolutionsParams{
			filter:  flagFilter,
			printer: prtr,
		})
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	SolutionCmd.AddCommand(solutionListCmd)
	solutionListCmd.PersistentFlags().StringSliceVarP(&flagFilter, keyFilter, "", []string{},
		fmt.Sprintf("Filter solutions by state. Available filters: %s."+
			" Separate multiple filters with a comma (without whitespaces in between).",
			strings.Join(allowedFilters, ",")))
}
