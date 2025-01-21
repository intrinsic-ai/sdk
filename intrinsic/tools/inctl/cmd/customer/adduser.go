// Copyright 2023 Intrinsic Innovation LLC

package customer

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	pb "intrinsic/kubernetes/accounts/service/api/accesscontrol/v1/accesscontrol_go_grpc_proto"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"
)

func init() {
	addUserInit(customerCmd)
}

var (
	flagEmail        string
	flagOrganization string
	flagRoleCSV      string
)

func parseCSV(s string) []string {
	parts := strings.Split(s, ",")
	for idx := range parts {
		parts[idx] = strings.TrimSpace(parts[idx])
	}
	return slices.DeleteFunc(parts, func(p string) bool { return p == "" })
}

func addUserInit(root *cobra.Command) {
	addUser.Flags().StringVar(&flagEmail, "email", "", "The email address of the user to invite.")
	addUser.Flags().StringVar(&flagOrganization, "organization", "", "The organization to invite the user to.")
	addUser.Flags().StringVar(&flagRoleCSV, "roles", "", "Optional comma-separated list of roles to assign to the user when they accept the invitation.")
	addUser.MarkFlagRequired("email")
	addUser.MarkFlagRequired("organization")
	root.AddCommand(addUser)
	listUsers.Flags().StringVar(&flagOrganization, "organization", "", "The organization to list memberships for.")
	listUsers.MarkFlagRequired("organization")
	root.AddCommand(listUsers)
}

var addUserHelp = `
Invite a user to an organization by email address.

Use the --roles flag (comma-separated list) to assign roles to the user after they accept the invitation.

Example:

		inctl customer add-user --email=user@example.com --organization=myorg --roles=owner
`

var addUser = &cobra.Command{
	Use:   "add-user",
	Short: "Invite a user to an organization by email address.",
	Long:  addUserHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := withOrgID(cmd.Context())
		cl, err := newAccessControlV1Client(ctx)
		if err != nil {
			return err
		}
		req := pb.CreateOrganizationInvitationRequest{
			Parent: addPrefix(flagOrganization, "organizations/"),
			Invitation: &pb.OrganizationInvitation{
				Organization: flagOrganization,
				Email:        flagEmail,
				Roles:        addPrefixes(parseCSV(flagRoleCSV), "roles/"),
			},
		}
		if flagDebugRequests {
			protoPrint(&req)
		}
		op, err := cl.CreateOrganizationInvitation(ctx, &req)
		if err != nil {
			return fmt.Errorf("failed to create organization: %w", err)
		}
		if flagDebugRequests {
			protoPrint(op)
		}
		return nil
	},
}

type users struct {
	// pending organization invitations.
	is []*pb.OrganizationInvitation
	// memberships of the organization.
	ms []*pb.OrganizationMembership
	// role bindings of the organization.
	rs []*pb.RoleBinding
}

func (us *users) String() string {
	b := new(bytes.Buffer)
	w := tabwriter.NewWriter(b,
		/*minwidth=*/ 1 /*tabwidth=*/, 1 /*padding=*/, 1 /*padchar=*/, ' ' /*flags=*/, 0)
	fmt.Fprintf(w, "%s\t%s\t%s\n", "Email", "Roles", "Status")
	// iterate memberships
	slices.SortFunc(us.ms, func(a, b *pb.OrganizationMembership) int {
		return strings.Compare(a.GetEmail(), b.GetEmail())
	})
	urs := userRoles(us.rs)
	for _, m := range us.ms {
		roles := urs[m.GetEmail()]
		fmt.Fprintf(w, "%s\t%s\t%s\n", m.GetEmail(), formatRoles(roles), "active")
	}
	// iterate invitations
	slices.SortFunc(us.is, func(a, b *pb.OrganizationInvitation) int {
		return strings.Compare(a.GetEmail(), b.GetEmail())
	})
	for _, o := range us.is {
		fmt.Fprintf(w, "%s\t%s\t%s\n", o.GetEmail(), formatRoles(o.GetRoles()), "pending")
	}
	w.Flush()
	// Remove the trailing newline as the pretty-printer wrapper will add one.
	return strings.TrimSuffix(b.String(), "\n")
}

func userRoles(rs []*pb.RoleBinding) map[string][]string {
	var roles = make(map[string][]string)
	for _, r := range rs {
		subject := strings.TrimPrefix(r.GetSubject(), "users/")
		roles[subject] = append(roles[subject], r.GetRole())
	}
	return roles
}

func formatRoles(rs []string) string {
	roles := []string{}
	for _, r := range rs {
		roles = append(roles, strings.TrimPrefix(r, "roles/"))
	}
	slices.Sort(roles)
	return strings.Join(roles, ", ")
}

var listUsersHelp = `
List all memberships and invitations of an organization.

Example:

		inctl customer list-users --organization=myorg
`

var listUsers = &cobra.Command{
	Use:   "list-users",
	Short: "List all memberships of an organization.",
	Long:  listUsersHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := withOrgID(cmd.Context())
		cl, err := newAccessControlV1Client(ctx)
		if err != nil {
			return err
		}
		ms, err := listMemberships(ctx, cl)
		if err != nil {
			return err
		}
		ois, err := listInvitations(ctx, cl)
		if err != nil {
			return err
		}
		rs, err := listRolesBindings(ctx, cl)
		if err != nil {
			return err
		}
		// format and print the results
		prtr, err := printer.NewPrinter(root.FlagOutput)
		if err != nil {
			return err
		}
		prtr.Print(&users{ms: ms, is: ois, rs: rs})
		return nil
	},
}

func listRolesBindings(ctx context.Context, cl accessControlV1Client) ([]*pb.RoleBinding, error) {
	req := pb.ListOrganizationRoleBindingsRequest{
		Parent: addPrefix(flagOrganization, "organizations/"),
	}
	if flagDebugRequests {
		protoPrint(&req)
	}
	op, err := cl.ListOrganizationRoleBindings(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to list role bindings: %w", err)
	}
	if flagDebugRequests {
		protoPrint(op)
	}
	return op.GetRoleBindings(), nil
}

func listMemberships(ctx context.Context, cl accessControlV1Client) ([]*pb.OrganizationMembership, error) {
	req := pb.ListOrganizationMembershipsRequest{
		Parent: addPrefix(flagOrganization, "organizations/"),
	}
	if flagDebugRequests {
		protoPrint(&req)
	}
	op, err := cl.ListOrganizationMemberships(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to list memberships: %w", err)
	}
	if flagDebugRequests {
		protoPrint(op)
	}
	return op.GetMemberships(), nil
}

func listInvitations(ctx context.Context, cl accessControlV1Client) ([]*pb.OrganizationInvitation, error) {
	req := pb.ListOrganizationInvitationsRequest{
		Parent: addPrefix(flagOrganization, "organizations/"),
	}
	if flagDebugRequests {
		protoPrint(&req)
	}
	op, err := cl.ListOrganizationInvitations(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to list organization invitations: %w", err)
	}
	if flagDebugRequests {
		protoPrint(op)
	}
	return op.GetInvitations(), nil
}
