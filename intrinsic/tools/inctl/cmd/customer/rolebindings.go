// Copyright 2023 Intrinsic Innovation LLC

package customer

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	pb "intrinsic/kubernetes/accounts/service/api/accesscontrol/v1/accesscontrol_go_grpc_proto"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
	"intrinsic/tools/inctl/util/printer"
)

var rolebindingsCmd = cobrautil.ParentOfNestedSubcommands("role-bindings", "List the role bindings on a given resource.")

func init() {
	customerCmd.AddCommand(rolebindingsCmd)
	rolebindingsInit(rolebindingsCmd)
}

var (
	flagResource string
	flagRole     string
	flagSubject  string
	flagName     string
)

func rolebindingsInit(root *cobra.Command) {
	listRoleBindingsCmd.Flags().StringVar(&flagResource, "resource", "", "The resource to list role-bindings for.")
	listRoleBindingsCmd.MarkFlagRequired("resource")
	root.AddCommand(listRoleBindingsCmd)
	grantRoleBindingCmd.Flags().StringVar(&flagResource, "resource", "", "The resource to attach the role-binding to.")
	grantRoleBindingCmd.Flags().StringVar(&flagSubject, "subject", "", "The subject grant the role.")
	grantRoleBindingCmd.Flags().StringVar(&flagRole, "role", "", "The role to grant.")
	grantRoleBindingCmd.MarkFlagRequired("resource")
	grantRoleBindingCmd.MarkFlagRequired("subject")
	grantRoleBindingCmd.MarkFlagRequired("role")
	root.AddCommand(grantRoleBindingCmd)
	revokeRoleBindingCmd.Flags().StringVar(&flagName, "name", "", "The name of the role-binding to revoke taken from the output of the list command.")
	revokeRoleBindingCmd.MarkFlagRequired("name")
	root.AddCommand(revokeRoleBindingCmd)
}

var grantRoleBindingCmdHelp = `
Grant a user a role on a given resource and all its descendants.

		inctl customer role-bindings grant --resource=organizations/exampleorg --subject=users/user@example.com --role=owner
`

var grantRoleBindingCmd = &cobra.Command{
	Use:   "grant",
	Short: "Grant a user a role on a given resource.",
	Long:  grantRoleBindingCmdHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cl, err := newAccessControlV1Client(ctx)
		if err != nil {
			return err
		}
		req := &pb.CreateRoleBindingRequest{
			RoleBinding: &pb.RoleBinding{
				Resource: flagResource,
				Role:     addPrefix(flagRole, "roles/"),
				Subject:  flagSubject,
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
		lrop, err = waitForOperation(ctx, cl.GetOperation, lrop, 10*time.Minute)
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
		cl, err := newAccessControlV1Client(ctx)
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
		lrop, err = waitForOperation(ctx, cl.GetOperation, lrop, 10*time.Minute)
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

		inctl customer role-bindings list --resource=organizations/exampleorg
`

var listRoleBindingsCmd = &cobra.Command{
	Use:   "list",
	Short: "List the role bindings on a given resource.",
	Long:  listRoleBindingsCmdHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if !strings.HasPrefix(flagResource, "organizations/") {
			return fmt.Errorf("only organizations are supported at the moment")
		}
		cl, err := newAccessControlV1Client(ctx)
		if err != nil {
			return err
		}
		req := &pb.ListOrganizationRoleBindingsRequest{
			Parent: flagResource,
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
