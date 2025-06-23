// Copyright 2023 Intrinsic Innovation LLC

package customer

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	accresourcemanager1pb "intrinsic/kubernetes/accounts/service/api/resourcemanager/v1/resourcemanager_go_grpc_proto"
)

func init() {
	organizationsInit(customerCmd)
}

var (
	flagOrgIdentifier   string
	flagOrgDisplayName  string
	flagSkipPaymentPlan bool
)

func organizationsInit(root *cobra.Command) {
	createCmd.Flags().StringVar(&flagOrgIdentifier, "identifier", "", "The human-friendly identifier of the organization to create.")
	createCmd.Flags().StringVar(&flagOrgDisplayName, "display-name", "", "The display name of the organization to create.")
	createCmd.Flags().BoolVar(&flagSkipPaymentPlan, "skip-payment-plan", false, "Skip creating a payment plan for the organization.")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("display-name")
	root.AddCommand(createCmd)
	root.AddCommand(deleteCmd)
}

var createCmdHelp = `
Create a new empty organization.

You must have permissions to create new organization on your current organization.

Example:

		inctl customer create --identifier=exampleorg --display-name="My Organization"
`

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new organization.",
	Long:  createCmdHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := withOrgID(cmd.Context())
		cl, err := newresourcemanagerClient(ctx)
		if err != nil {
			return err
		}
		req := accresourcemanager1pb.CreateOrganizationRequest{
			OrganizationId: flagOrgIdentifier,
			Organization: &accresourcemanager1pb.Organization{
				DisplayName: flagOrgDisplayName,
			},
		}
		if flagDebugRequests {
			protoPrint(&req)
		}
		fmt.Printf("Creating organization %q.\n", flagOrgIdentifier)
		op, err := cl.CreateOrganization(ctx, &req)
		if err != nil {
			return fmt.Errorf("failed to create organization: %w", err)
		}
		if flagDebugRequests {
			protoPrint(op)
		}
		op, err = waitForOperation(ctx, cl.GetOperation, op, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("failed to wait for operation: %w", err)
		}
		if flagDebugRequests {
			protoPrint(op)
		}
		if flagSkipPaymentPlan {
			fmt.Println("Warning: skipping payment plan creation. The organization will have no quota assigned.")
			return nil
		}
		preq := &accresourcemanager1pb.CreateOrganizationPaymentPlanRequest{
			Parent: "organizations/" + flagOrgIdentifier,
		}
		if flagDebugRequests {
			protoPrint(preq)
		}
		fmt.Println("Creating a payment plan for the organization.")
		op, err = cl.CreateOrganizationPaymentPlan(ctx, preq)
		if err != nil {
			return fmt.Errorf("failed to create organization payment plan: %w", err)
		}
		if flagDebugRequests {
			protoPrint(op)
		}
		op, err = waitForOperation(ctx, cl.GetOperation, op, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("failed to wait for operation: %w", err)
		}
		if flagDebugRequests {
			protoPrint(op)
		}
		return nil
	},
}

var deleteCmdHelp = `
Delete an organization.

You must have permissions to delete existing organizations on your current organization and your
organization must have owner role for the organization you are trying to delete.

Example:

		inctl customer delete --identifier=exampleorg --org=myorg
`

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an organization.",
	Long:  createCmdHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := withOrgID(cmd.Context())
		cl, err := newresourcemanagerClient(ctx)
		if err != nil {
			return err
		}
		req := accresourcemanager1pb.DeleteOrganizationRequest{
			Name: "organizations/" + flagOrgIdentifier,
		}
		if flagDebugRequests {
			protoPrint(&req)
		}
		fmt.Printf("Deleting organization %q.\n", flagOrgIdentifier)
		op, err := cl.DeleteOrganization(ctx, &req)
		if err != nil {
			return fmt.Errorf("failed to create organization: %w", err)
		}
		if flagDebugRequests {
			protoPrint(op)
		}
		op, err = waitForOperation(ctx, cl.GetOperation, op, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("failed to wait for operation: %w", err)
		}
		if flagDebugRequests {
			protoPrint(op)
		}
		return nil
	},
}
