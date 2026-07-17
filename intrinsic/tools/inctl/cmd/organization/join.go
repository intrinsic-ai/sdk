// Copyright 2023 Intrinsic Innovation LLC

package organization

import (
	"fmt"
	"time"

	"intrinsic/tools/inctl/util/accounts/accounts"

	"github.com/spf13/cobra"

	pb "intrinsic/kubernetes/accounts/service/api/invitations/v1/invitations_go_proto"
)

func init() {
	joinCmd.Flags().StringVar(&flagInvitationToken, "token", "", "The token of the invitation to accept.")
	organizationCmd.AddCommand(joinCmd)
}

var joinCmdHelp = `
Accept an invitation to join an organization using an invitation token.

Example:

		inctl organization join 24d7ab1f-8c55-4903-9352-4ce421bef264
		inctl organization join --token=24d7ab1f-8c55-4903-9352-4ce421bef264
`

var joinCmd = &cobra.Command{
	Use:   "join [token]",
	Short: "Accept an invitation to join an organization.",
	Long:  joinCmdHelp,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := resolveInvitationTokenArgOrFlag(args)
		if err != nil {
			return err
		}
		ctx := cmd.Context()
		cl, lroCl, err := newInvitationsV1Client(ctx)
		if err != nil {
			return err
		}
		req := &pb.ApplyInvitationRequest{
			Token: token,
		}
		if flagDebugRequests {
			protoPrint(req)
		}
		fmt.Printf("Joining organization via invitation token %q...\n", token)
		lrop, err := cl.ApplyInvitation(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to apply invitation: %w", err)
		}
		if flagDebugRequests {
			protoPrint(lrop)
		}
		lrop, err = accounts.WaitForOperation(ctx, lroCl.GetOperation, lrop, 2*time.Minute)
		if err != nil {
			return fmt.Errorf("failed to wait for operation: %w", err)
		}
		if flagDebugRequests {
			protoPrint(lrop)
		}
		if lrop.GetError() != nil {
			return fmt.Errorf("failed to join organization: %v", lrop.GetError())
		}
		fmt.Printf("Successfully joined organization.\n")
		return nil
	},
}
