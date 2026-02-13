// Copyright 2023 Intrinsic Innovation LLC

package customer

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/accounts/accounts"
	"intrinsic/tools/inctl/util/cobrautil"
	"intrinsic/tools/inctl/util/printer"

	"github.com/spf13/cobra"

	pb "intrinsic/kubernetes/accounts/service/api/accesscontrol/v1/accesscontrol_go_proto"
)

var rolebindingsCmd = cobrautil.ParentOfNestedSubcommands("role-bindings", "List the role bindings on a given resource.")

func init() {
	customerCmd.AddCommand(rolebindingsCmd)
	rolebindingsInit(rolebindingsCmd)
}

func rolebindingsInit(root *cobra.Command) {
	listRoleBindingsCmd.Flags().StringVar(&flagCustomer, "customer", "", "The customer organization to list role-bindings for.")
	listRoleBindingsCmd.MarkFlagRequired("customer")
	root.AddCommand(listRoleBindingsCmd)
	grantRoleBindingCmd.Flags().StringVar(&flagCustomer, "customer", "", "The customer organization to attach the role-binding to.")
	grantRoleBindingCmd.Flags().StringVar(&flagEmail, "email", "", "The email address of the user to grant the role to.")
	grantRoleBindingCmd.Flags().StringVar(&flagRole, "role", "", "The role to grant.")
	grantRoleBindingCmd.MarkFlagRequired("customer")
	grantRoleBindingCmd.MarkFlagRequired("email")
	grantRoleBindingCmd.MarkFlagRequired("role")
	root.AddCommand(grantRoleBindingCmd)
	revokeRoleBindingCmd.Flags().StringVar(&flagName, "name", "", "The name of the role-binding to revoke taken from the output of the list command.")
	revokeRoleBindingCmd.MarkFlagRequired("name")
	root.AddCommand(revokeRoleBindingCmd)
}

var grantRoleBindingCmdHelp = `
Grant a user a role on a given resource and all its descendants.

		inctl customer role-bindings grant --customer=exampleorg --email=user@example.com --role=owner
`

var grantRoleBindingCmd = &cobra.Command{
	Use:   "grant",
	Short: "Grant a user a role on a given resource.",
	Long:  grantRoleBindingCmdHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cl, err := accounts.NewAccessControlV1Client(ctx, vipr)
		if err != nil {
			return err
		}
		req := &pb.CreateRoleBindingRequest{
			RoleBinding: &pb.RoleBinding{
				Resource: addPrefix(flagCustomer, "organizations/"),
				Role:     addPrefix(flagRole, "roles/"),
				Subject:  addPrefix(flagEmail, "users/"),
			},
		}
		if flagDebugRequests {
			protoPrint(req)
		}
		lrop, err := cl.CreateRoleBinding(ctx, req)
		if err != nil {
			return err
		}
		if flagDebugRequests {
			protoPrint(lrop)
		}
		lrop, err = accounts.WaitForOperation(ctx, cl.GetOperation, lrop, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("failed to wait for operation: %w", err)
		}
		if flagDebugRequests {
			protoPrint(lrop)
		}
		return nil
	},
}

var revokeRoleBindingCmdHelp = `
Revoke a given role binding.

		inctl customer role-bindings revoke --name=rolebindings/7iawfQMYZAMkx6XdmQdtqJfW+gCZeoT83PcYw0daIrg=
`

var revokeRoleBindingCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke a given role binding.",
	Long:  revokeRoleBindingCmdHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cl, err := accounts.NewAccessControlV1Client(ctx, vipr)
		if err != nil {
			return err
		}
		req := &pb.DeleteRoleBindingRequest{
			Name: addPrefix(flagName, "rolebindings/"),
		}
		if flagDebugRequests {
			protoPrint(req)
		}
		lrop, err := cl.DeleteRoleBinding(ctx, req)
		if err != nil {
			return err
		}
		if flagDebugRequests {
			protoPrint(lrop)
		}
		lrop, err = accounts.WaitForOperation(ctx, cl.GetOperation, lrop, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("failed to wait for operation: %w", err)
		}
		if flagDebugRequests {
			protoPrint(lrop)
		}
		return nil
	},
}

type printableRoleBindings []*pb.RoleBinding

func (r printableRoleBindings) String() string {
	b := new(bytes.Buffer)
	w := tabwriter.NewWriter(b,
		/*minwidth=*/ 1 /*tabwidth=*/, 1 /*padding=*/, 1 /*padchar=*/, ' ' /*flags=*/, 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "Name", "Resource", "Role", "Subject")
	for _, rb := range r {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", rb.GetName(), rb.GetResource(), rb.GetRole(), rb.GetSubject())
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n")
}

var listRoleBindingsCmdHelp = `
List the role bindings on a given resource.

		inctl customer role-bindings list --customer=exampleorg
`

var listRoleBindingsCmd = &cobra.Command{
	Use:   "list",
	Short: "List the role bindings on a given resource.",
	Long:  listRoleBindingsCmdHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cOrg := addPrefix(flagCustomer, "organizations/")
		cl, err := accounts.NewAccessControlV1Client(ctx, vipr)
		if err != nil {
			return err
		}
		req := &pb.ListOrganizationRoleBindingsRequest{
			Parent: cOrg,
		}
		if flagDebugRequests {
			protoPrint(req)
		}
		ret, err := cl.ListOrganizationRoleBindings(ctx, req)
		if err != nil {
			return err
		}
		if flagDebugRequests {
			protoPrint(ret)
		}
		prtr, err := printer.NewPrinter(root.FlagOutput)
		if err != nil {
			return err
		}
		prtr.Print(printableRoleBindings(ret.GetRoleBindings()))
		return nil
	},
}
