// Copyright 2023 Intrinsic Innovation LLC

package customer

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"intrinsic/tools/inctl/util/accounts/accounts"
	"intrinsic/tools/inctl/util/orgutil"
	"intrinsic/tools/inctl/util/printer"

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
	root.AddCommand(listCmd)
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

var listCmdHelp = `
List customer organizations.

List all organizations that are children of the current organization. The membership status indicates
whether you are also a member of the organization or if you are granted access to the organization
through a role binding without being a member.

Example:

		inctl customer list --org=myorg
`

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List customer organizations.",
	Long:  listCmdHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := accounts.WithOrgID(cmd.Context(), vipr)
		cl, err := accounts.NewResourceManagerV1Client(ctx, vipr)
		if err != nil {
			return err
		}
		wg := errgroup.Group{}
		// list all the organizations the user has access to
		var os []*accresourcemanager1pb.Organization
		wg.Go(func() error {
			var err error
			os, err = listMyOrganizations(ctx, cl, "")
			if err != nil {
				return err
			}
			return nil
		})
		// list a second time with the filter to get only the organizations the user is a member of
		var msOs []*accresourcemanager1pb.Organization
		wg.Go(func() error {
			var err error
			msOs, err = listMyOrganizations(ctx, cl, "is_member()")
			if err != nil {
				return err
			}
			return nil
		})
		if err := wg.Wait(); err != nil {
			return err
		}
		msMap := make(map[string]bool)
		for _, o := range msOs {
			msMap[o.GetName()] = true
		}
		// convert to orgRows and sort them
		orgs := make([]*orgRow, 0, len(os))
		authOrg := vipr.GetString(orgutil.KeyOrganization)
		for _, o := range os {
			if o.GetParentOrganization() != addPrefix(authOrg, "organizations/") {
				continue
			}
			orgs = append(orgs, toOrgRow(o, msMap))
		}
		slices.SortFunc(orgs, func(a, b *orgRow) int { return strings.Compare(a.Name, b.Name) })
		if err := printOrganizations(cmd, orgs); err != nil {
			return err
		}
		return nil
	},
}

func getOrgRowCommandPrinter(cmd *cobra.Command) printer.CommandPrinter {
	ot := printer.GetFlagOutputType(cmd)
	if ot == printer.OutputTypeText {
		ot = printer.OutputTypeTAB
	}
	cp, err := printer.NewPrinterOfType(
		ot,
		cmd,
		printer.WithDefaultsFromValue(&accresourcemanager1pb.Organization{}, func(columns []string) []string {
			return orgColumns
		}),
		printer.WithHeaderOverride(func(s string) string {
			return orgColumnNames[s]
		}),
	)
	if err != nil {
		cmd.PrintErrf("Error setting up output: %v\n", err)
		cp = printer.GetDefaultPrinter(cmd)
	}
	return cp
}

var orgColumns = []string{"name", "display_name", "create_time", "membership_status"}

var orgColumnNames = map[string]string{
	"name":              "Name",
	"display_name":      "Display Name",
	"create_time":       "Create Time",
	"membership_status": "Membership Status",
}

type orgRow struct {
	Name             string `json:"name"`
	DisplayName      string `json:"display_name"`
	CreateTime       string `json:"create_time"`
	MembershipStatus bool   `json:"membership_status"`
}

func toOrgRow(o *accresourcemanager1pb.Organization, msMap map[string]bool) *orgRow {
	or := &orgRow{
		Name:             strings.TrimPrefix(o.GetName(), "organizations/"),
		DisplayName:      o.GetDisplayName(),
		CreateTime:       o.GetCreateTime().AsTime().Format(time.RFC3339),
		MembershipStatus: msMap[o.GetName()],
	}
	return or
}

func printOrganizations(cmd *cobra.Command, orgs []*orgRow) error {
	prtr := getOrgRowCommandPrinter(cmd)
	var view printer.View = nil // this is to reuse reflectors in default views
	for _, o := range orgs {
		view = printer.NextView(o, view)
		prtr.Println(view)
	}
	return printer.Flush(prtr)
}

func listMyOrganizations(ctx context.Context, cl accounts.ResourceManagerV1Client, filter string) ([]*accresourcemanager1pb.Organization, error) {
	req := &accresourcemanager1pb.ListOrganizationsRequest{
		Filter: filter,
	}
	resp, err := cl.ListOrganizations(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	return resp.GetOrganizations(), nil
}
