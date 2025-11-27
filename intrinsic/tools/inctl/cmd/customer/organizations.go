// Copyright 2023 Intrinsic Innovation LLC

package customer

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"intrinsic/tools/inctl/util/accounts/accounts"
	"intrinsic/tools/inctl/util/orgutil"
	"intrinsic/tools/inctl/util/printer"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	accresourcemanager1pb "intrinsic/kubernetes/accounts/service/api/resourcemanager/v1/resourcemanager_go_grpc_proto"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
)

func init() {
	organizationsInit(customerCmd)
}

var (
	flagOrgDisplayName string
	flagYes            bool
)

func organizationsInit(root *cobra.Command) {
	createCmd.Flags().StringVar(&flagCustomer, "customer", "", "The human-friendly identifier of the organization to create (format: ^[a-z][a-z0-9_-]{0,63}$). Keep empty to auto-generate a random identifier.")
	createCmd.Flags().StringVar(&flagOrgDisplayName, "display-name", "", "The display name of the organization to create.")
	createCmd.MarkFlagRequired("display-name")
	root.AddCommand(createCmd)
	deleteCmd.Flags().StringVar(&flagCustomer, "customer", "", "The human-friendly identifier of the organization to delete.")
	deleteCmd.MarkFlagRequired("customer")
	deleteCmd.Flags().BoolVar(&flagYes, "yes", false, "Skip the confirmation prompt and directly delete the organization.")
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
		ctx := cmd.Context()
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
		if flagCustomer == "" {
			fmt.Printf("Creating organization with random identifier (display name: %q).\n", flagOrgDisplayName)
		} else {
			fmt.Printf("Creating organization %q (%q).\n", flagCustomer, flagOrgDisplayName)
		}

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
		operr := op.GetError()
		if operr != nil {
			return fmt.Errorf("failed to create organization: %v", operr)
		}
		if flagCustomer == "" && op.GetResult() != nil {
			res := op.GetResult().(*lropb.Operation_Response).Response
			org := &accresourcemanager1pb.Organization{}
			if err := res.UnmarshalTo(org); err != nil {
				return fmt.Errorf("failed to unmarshal organization: %w", err)
			}
			fmt.Printf("Created organization with identifier %q.\n", org.GetName())
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
		ctx := cmd.Context()
		cl, err := accounts.NewResourceManagerV1Client(ctx, vipr)
		if err != nil {
			return err
		}
		// if the user did not specify --yes, ask for confirmation
		// also skips client-side input validations (like org state)
		if !flagYes {
			if err := confirmDelete(ctx, cl); err != nil {
				return err
			}
		}
		req := accresourcemanager1pb.DeleteOrganizationRequest{
			Name: addPrefix(flagCustomer, "organizations/"),
		}
		if flagDebugRequests {
			protoPrint(&req)
		}
		fmt.Printf("Deleting organization %q ...\n", flagCustomer)
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

// confirmDelete first fetches the organization to check if it exists and to display its display name to the user.
// Also performs some client-side input validations (like org state).
func confirmDelete(ctx context.Context, cl accounts.ResourceManagerV1Client) error {
	o, err := cl.GetOrganization(ctx, &accresourcemanager1pb.GetOrganizationRequest{Name: addPrefix(flagCustomer, "organizations/")})
	if err != nil {
		return fmt.Errorf("failed to find organization %q: %w", flagCustomer, err)
	}
	if o.GetState() == accresourcemanager1pb.Organization_DELETED {
		return fmt.Errorf("organization %q is already deleted", flagCustomer)
	}
	// if the user did not specify --yes, ask for confirmation
	if !flagYes {
		fmt.Printf("You are about to delete organization %q (%q). Please type the organization identifier to confirm: ",
			flagCustomer, o.DisplayName)
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		if input != flagCustomer {
			return fmt.Errorf("input does not match organization identifier %q", flagCustomer)
		}
	}
	return nil
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
		ctx := cmd.Context()
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
