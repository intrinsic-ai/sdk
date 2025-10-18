// Copyright 2023 Intrinsic Innovation LLC

// Package list provides a command to list solutions.
package list

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/orgutil"
	"intrinsic/tools/inctl/util/printer"

	clusterdiscoverypb "intrinsic/frontend/cloud/api/v1/clusterdiscovery_api_go_grpc_proto"
	solutiondiscoverygrpcpb "intrinsic/frontend/cloud/api/v1/solutiondiscovery_api_go_grpc_proto"
	solutiondiscoverypb "intrinsic/frontend/cloud/api/v1/solutiondiscovery_api_go_grpc_proto"
)

const (
	keyFilter = "filter"
)

var (
	flagFilter     []string
	allowedFilters = []string{"not_running", "running_in_sim", "running_on_hw"}
)

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
	}
	solutions := make([]solution, len(res.m.GetSolutions()))
	for i, c := range res.m.GetSolutions() {
		solutions[i] = solution{
			Name:        c.GetName(),
			State:       c.GetState().String(),
			DisplayName: c.GetDisplayName(),
			ClusterName: c.GetClusterName(),
		}
	}
	return json.Marshal(struct {
		// solution intentionally not omitted when empty
		Solutions []solution `json:"solutions"`
	}{Solutions: solutions})
}

type listSolutionsParams struct {
	filter     []string
	printer    printer.CommandPrinter
	outputType printer.OutputType
}

type solutionRow struct {
	Name        string `json:"name,omitempty"`
	State       string `json:"state,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	ClusterName string `json:"clusterName,omitempty"`
}

func asSolutionRow(p *solutiondiscoverypb.SolutionDescription) *solutionRow {
	return &solutionRow{
		Name:        p.GetName(),
		State:       p.GetState().String(),
		DisplayName: p.GetDisplayName(),
		ClusterName: p.GetClusterName(),
	}
}

func getSolutionRowCommandPrinter(cmd *cobra.Command) printer.CommandPrinter {
	ot := printer.GetFlagOutputType(cmd)
	if ot == printer.OutputTypeText {
		ot = printer.OutputTypeTAB
	}
	cp, err := printer.NewPrinterOfType(
		ot,
		cmd,
		printer.WithDefaultsFromValue(&solutionRow{}, func(columns []string) []string {
			return []string{"name", "state", "displayName"}
		}),
	)
	if err != nil {
		cmd.PrintErrf("Error setting up output: %v\n", err)
		cp = printer.GetDefaultPrinter(cmd)
	}
	return cp
}

func (s *solutionRow) Tabulated(columns []string) []string {
	name := s.DisplayName
	if name == "" {
		name = s.Name
	}

	statusStr := strings.TrimPrefix(s.State, "SOLUTION_STATE_")
	if s.ClusterName != "" {
		statusStr = fmt.Sprintf("%s on %s", statusStr, s.ClusterName)
	}

	return []string{name, statusStr, s.Name}
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

// listSolutions lists solutions matching the given filters. If the output type is JSON, the
// solutions are printed in a custom JSON format. Otherwise, the default solutionRow printer is
// used.
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

	// Keep the JSON output format as is, in case there are consumers relying on it.
	if params.outputType == printer.OutputTypeJSON {
		params.printer.Print(&ListSolutionDescriptionsResponse{m: resp})
	} else {
		var view printer.View = nil // this is to reuse reflectors in default views
		for _, p := range resp.GetSolutions() {
			view = printer.NextView(asSolutionRow(p), view)
			params.printer.Println(view)
		}
	}
	return printer.Flush(params.printer)
}

// NewCommand returns the list command.
func NewCommand() *cobra.Command {
	viperLocal := viper.New()

	solutionListCmd := orgutil.WrapCmd(&cobra.Command{
		Use:   "list",
		Short: "List solutions in a project",
		Long:  "List solutions on the given project.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			prtr := getSolutionRowCommandPrinter(cmd)
			ctx := cmd.Context()
			conn, err := auth.NewCloudConnection(ctx, auth.WithFlagValues(viperLocal))
			if err != nil {
				return err
			}
			defer conn.Close()

			err = listSolutions(ctx, conn, &listSolutionsParams{
				filter:     flagFilter,
				printer:    prtr,
				outputType: printer.GetFlagOutputType(cmd),
			})
			if err != nil {
				return err
			}
			return nil
		},
	}, viperLocal)

	solutionListCmd.PersistentFlags().StringSliceVarP(&flagFilter, keyFilter, "", []string{},
		fmt.Sprintf("Filter solutions by state. Available filters: %s."+
			" Separate multiple filters with a comma (without whitespaces in between).",
			strings.Join(allowedFilters, ",")))

	return solutionListCmd
}
