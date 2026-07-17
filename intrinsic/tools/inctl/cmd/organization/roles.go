// Copyright 2023 Intrinsic Innovation LLC

package organization

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
	"intrinsic/tools/inctl/util/printer"

	"github.com/spf13/cobra"

	accaccesscontrolv1pb "intrinsic/kubernetes/accounts/service/api/accesscontrol/v1/accesscontrol_go_proto"
)

var rolesCmd = cobrautil.ParentOfNestedSubcommands("roles", "List available roles.")

func init() {
	organizationCmd.AddCommand(rolesCmd)
	rolesInit(rolesCmd)
}

func rolesInit(root *cobra.Command) {
	root.AddCommand(listRolesCmd)
}

var listRolesCmdHelp = `
List available roles.

Example:
  inctl organization roles list
`

type printableRoles []*accaccesscontrolv1pb.Role

func (r printableRoles) String() string {
	b := new(bytes.Buffer)
	w := tabwriter.NewWriter(b,
		/*minwidth=*/ 1 /*tabwidth=*/, 1 /*padding=*/, 1 /*padchar=*/, ' ' /*flags=*/, 0)
	fmt.Fprintf(w, "%s\t%s\t%s\n", "Name", "Display Name", "Description")
	for _, role := range r {
		fmt.Fprintf(w, "%s\t%s\t%s\n", role.GetName(), role.GetDisplayName(), role.GetDescription())
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n")
}

func (r printableRoles) MarshalJSON() ([]byte, error) {
	return marshalProtoSlice(r)
}

var listRolesCmd = &cobra.Command{
	Use:   "list",
	Short: "List available roles.",
	Long:  listRolesCmdHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkOrgNotIntrinsic(); err != nil {
			return err
		}
		ctx := cmd.Context()
		cl, err := newAccessControlV1Client(ctx)
		if err != nil {
			return err
		}
		rs, err := cl.ListRoles(ctx, &accaccesscontrolv1pb.ListRolesRequest{})
		if err != nil {
			return err
		}
		prtr, err := printer.NewPrinter(root.FlagOutput)
		if err != nil {
			return err
		}
		prtr.Print(printableRoles(rs.GetRoles()))
		return nil
	},
}
