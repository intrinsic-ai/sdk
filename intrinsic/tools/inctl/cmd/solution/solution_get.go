// Copyright 2023 Intrinsic Innovation LLC

package solution

import (
	"context"
	"encoding/json"
	"fmt"

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

// GetSolutionParams contains all parameters to retrieve a solution via GetSolution
type GetSolutionParams struct {
	// The name of the solution to get
	SolutionName string
	// The server to connect to
	ServerAddr string
}

type getSolutionDescriptionWrapper struct {
	solutionName string
	solution     *solutiondiscoverypb.SolutionDescription
}

// MarshalJSON converts a getSolutionDescriptionWrapper to a byte slice.
func (wrapper *getSolutionDescriptionWrapper) MarshalJSON() ([]byte, error) {
	var stateStr string
	if wrapper.solution.GetState() == clusterdiscoverypb.SolutionState_SOLUTION_STATE_UNSPECIFIED {
		stateStr = ""
	} else {
		stateStr = wrapper.solution.GetState().String()
	}

	return json.Marshal(struct {
		Name        string `json:"name,omitempty"`
		State       string `json:"state,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		ClusterName string `json:"clusterName,omitempty"`
	}{
		Name:        wrapper.solution.GetName(),
		State:       stateStr,
		DisplayName: wrapper.solution.GetDisplayName(),
		ClusterName: wrapper.solution.GetClusterName(),
	},
	)
}

// String converts a getSolutionDescriptionWrapper into a human-readable string
func (wrapper *getSolutionDescriptionWrapper) String() string {
	switch wrapper.solution.GetState() {
	case clusterdiscoverypb.SolutionState_SOLUTION_STATE_NOT_RUNNING:
		return fmt.Sprintf("Solution %q is currently not running.",
			wrapper.solution.GetName())
	case clusterdiscoverypb.SolutionState_SOLUTION_STATE_RUNNING_ON_HW:
		return fmt.Sprintf("Solution %q (%s) is running on hardware cluster %q.",
			wrapper.solution.GetDisplayName(), wrapper.solution.GetName(),
			wrapper.solution.GetClusterName())
	case clusterdiscoverypb.SolutionState_SOLUTION_STATE_RUNNING_IN_SIM:
		return fmt.Sprintf("Solution %q (%s) is running in simulation on cluster %q.",
			wrapper.solution.GetDisplayName(), wrapper.solution.GetName(),
			wrapper.solution.GetClusterName())
	}
	return fmt.Sprintf("Cannot determine the current state of solution %q (%s)",
		wrapper.solution.GetName(), wrapper.solution.GetState().String())
}

// GetSolution gets solution data by name
func GetSolution(ctx context.Context, conn *grpc.ClientConn, solutionName string) (*solutiondiscoverypb.SolutionDescription, error) {
	client := solutiondiscoverygrpcpb.NewSolutionDiscoveryServiceClient(conn)
	req := &solutiondiscoverypb.GetSolutionDescriptionRequest{Name: solutionName}
	resp, err := client.GetSolutionDescription(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get solution description: %w", err)
	}

	return resp.GetSolution(), nil
}

func getAndPrintSolution(ctx context.Context, conn *grpc.ClientConn, solutionName string, printer printer.Printer) error {
	resp, err := GetSolution(ctx, conn, solutionName)
	if err != nil {
		return fmt.Errorf("request to get solution '%s' failed: %w", solutionName, err)
	}

	printer.Print(&getSolutionDescriptionWrapper{solutionName: solutionName, solution: resp})
	return nil
}

var solutionGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get solution in a project",
	Long:  "Get a solution by name (the unique identifier - not the display name) on the given project.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		solutionName := args[0]
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

		if err = getAndPrintSolution(ctx, conn, solutionName, prtr); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	SolutionCmd.AddCommand(solutionGetCmd)
}
