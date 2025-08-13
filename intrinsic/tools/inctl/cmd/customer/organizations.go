// Copyright 2023 Intrinsic Innovation LLC

package customer

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"intrinsic/tools/inctl/util/accounts/accounts"

	accresourcemanager1pb "intrinsic/kubernetes/accounts/service/api/resourcemanager/v1/resourcemanager_go_grpc_proto"
)

func init() {
	organizationsInit(customerCmd)
}

var (
	flagOrgDisplayName  string
	flagSkipPaymentPlan bool
)

func organizationsInit(root *cobra.Command) {
	createCmd.Flags().StringVar(&flagCustomer, "customer", "", "The human-friendly identifier of the organization to create.")
	createCmd.Flags().StringVar(&flagOrgDisplayName, "display-name", "", "The display name of the organization to create.")
	createCmd.Flags().BoolVar(&flagSkipPaymentPlan, "skip-payment-plan", false, "Skip creating a payment plan for the organization.")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("display-name")
	createCmd.MarkFlagRequired("customer")
	root.AddCommand(createCmd)
	deleteCmd.Flags().StringVar(&flagCustomer, "customer", "", "The human-friendly identifier of the organization to delete.")
	deleteCmd.MarkFlagRequired("customer")
	root.AddCommand(deleteCmd)
}

var createCmdHelp = `
Create a new empty organization.

You must have permissions to create new organization on your current organization.

Example:

		inctl customer create --customer=exampleorg --display-name="My Organization"
`

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new organization.",
	Long:  createCmdHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := accounts.WithOrgID(cmd.Context(), vipr)
		cl, err := accounts.NewResourceManagerV1Client(ctx, vipr)
		if err != nil {
			return err
		}
		req := accresourcemanager1pb.CreateOrganizationRequest{
			OrganizationId: flagCustomer,
			Organization: &accresourcemanager1pb.Organization{
				DisplayName: flagOrgDisplayName,
			},
		}
		if flagDebugRequests {
			protoPrint(&req)
		}
		fmt.Printf("Creating organization %q.\n", flagCustomer)
		op, err := cl.CreateOrganization(ctx, &req)
		if err != nil {
			return fmt.Errorf("failed to create organization: %w", err)
		}
		if flagDebugRequests {
			protoPrint(op)
		}
		op, err = accounts.WaitForOperation(ctx, cl.GetOperation, op, 10*time.Minute)
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
			Parent: addPrefix(flagCustomer, "organizations/"),
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
		op, err = accounts.WaitForOperation(ctx, cl.GetOperation, op, 10*time.Minute)
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

The delete command marks the organization as soft-deleted. A soft-deleted organization can be
recovered for at least 30 days by contacting support.

Example:

		inctl customer delete --customer=exampleorg --org=myorg
`

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an organization.",
	Long:  deleteCmdHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := accounts.WithOrgID(cmd.Context(), vipr)
		cl, err := accounts.NewResourceManagerV1Client(ctx, vipr)
		if err != nil {
			return err
		}
		req := accresourcemanager1pb.DeleteOrganizationRequest{
			Name: addPrefix(flagCustomer, "organizations/"),
		}
		if flagDebugRequests {
			protoPrint(&req)
		}
		fmt.Printf("Deleting organization %q.\n", flagCustomer)
		op, err := cl.DeleteOrganization(ctx, &req)
		if err != nil {
			return fmt.Errorf("failed to delete organization: %w", err)
		}
		if flagDebugRequests {
			protoPrint(op)
		}
		op, err = accounts.WaitForOperation(ctx, cl.GetOperation, op, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("failed to wait for operation: %w", err)
		}
		if flagDebugRequests {
			protoPrint(op)
		}
		return nil
	},
}
