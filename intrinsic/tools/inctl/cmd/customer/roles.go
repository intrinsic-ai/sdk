// Copyright 2023 Intrinsic Innovation LLC

package customer

import (
	"fmt"

	"intrinsic/tools/inctl/util/accounts/accounts"
	"intrinsic/tools/inctl/util/cobrautil"

	"github.com/spf13/cobra"

	accaccesscontrolv1pb "intrinsic/kubernetes/accounts/service/api/accesscontrol/v1/accesscontrol_go_proto"
)

var rolesCmd = cobrautil.ParentOfNestedSubcommands("roles", "List available roles.")

func init() {
	customerCmd.AddCommand(rolesCmd)
	rolesInit(rolesCmd)
}

func rolesInit(root *cobra.Command) {
	root.AddCommand(listRolesCmd)
}

var listRolesCmdHelp = `
List available roles.

Example:
  inctl customer roles list
`

var listRolesCmd = &cobra.Command{
	Use:   "list",
	Short: "List available roles.",
	Long:  listRolesCmdHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cl, err := accounts.NewAccessControlV1Client(ctx, vipr)
		if err != nil {
			return err
		}
		rs, err := cl.ListRoles(ctx, &accaccesscontrolv1pb.ListRolesRequest{})
		if err != nil {
			return err
		}
		for _, r := range rs.GetRoles() {
			fmt.Printf("%s - %s - %s\n", r.GetName(), r.GetDisplayName(), r.GetDescription())
		}
		return nil
	},
}
