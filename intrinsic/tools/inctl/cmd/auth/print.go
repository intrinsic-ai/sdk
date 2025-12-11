// Copyright 2023 Intrinsic Innovation LLC

package auth

import (
	"fmt"
	"net/http"

	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/orgutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flagFlowstateAddr string

func init() {
	authCmd.AddCommand(printAPIKeyCmd)
	printAPIKeyCmd.Flags().MarkHidden(orgutil.KeyProject)

	authCmd.AddCommand(printAccessTokenCmd)
	printAccessTokenCmd.Flags().StringVar(&flagFlowstateAddr, "flowstate", "flowstate.intrinsic.ai", "Flowstate address.")
	printAccessTokenCmd.Flags().MarkHidden(orgutil.KeyProject)
}

var printAPIKeyParams = viper.New()

var printAPIKeyCmd = orgutil.WrapCmd(&cobra.Command{
	Use:   "print-api-key",
	Short: "Prints the API key for a project.",
	Long:  "Prints the API key for a project.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		project := printAPIKeyParams.GetString(orgutil.KeyProject)
		store, err := authStore.GetConfiguration(project)
		if err != nil {
			return fmt.Errorf("failed to get configuration for project %q: %v", project, err)
		}
		key, err := store.GetDefaultCredentials()
		if err != nil {
			return fmt.Errorf("failed to get default API key for project %q: %v", project, err)
		}
		fmt.Print(key.APIKey)
		return nil
	},
}, printAPIKeyParams, orgutil.WithOrgExistsCheck(func() bool { return checkOrgExists }))

var makeHTTPClient = func() *http.Client { // for unit-tests
	return &http.Client{}
}

var printAccessTokenHelp = `
Print an access token.

Can be used to authenticate with the Flowstate API.

Example:

		inctl auth print-access-token --org=myorganization

Example (curl):

		curl -s -X GET -H "Authorization: Bearer $(inctl auth print-access-token --org=myorganization)" https://flowstate.intrinsic.ai/api/v1/cloud-projects-orgs -H 'Content-Type: application/json'
`

var printAccessTokenParams = viper.New()

var printAccessTokenCmd = orgutil.WrapCmd(&cobra.Command{
	Use:   "print-access-token",
	Short: "Print an access token.",
	Long:  printAccessTokenHelp,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		project := printAccessTokenParams.GetString(orgutil.KeyProject)
		store, err := authStore.GetConfiguration(project)
		if err != nil {
			return fmt.Errorf("failed to get configuration for project %q: %v", project, err)
		}
		key, err := store.GetDefaultCredentials()
		if err != nil {
			return fmt.Errorf("failed to get default API key for project %q: %v", project, err)
		}
		ctx := cmd.Context()
		resp, err := auth.GetIDToken(ctx, makeHTTPClient(), flagFlowstateAddr, &auth.GetIDTokenRequest{
			APIKey:   key.APIKey,
			DoFanOut: true,
		})
		if err != nil {
			return fmt.Errorf("failed to get ID token : %v", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s", resp.IDToken)
		return nil
	},
}, printAccessTokenParams, orgutil.WithOrgExistsCheck(func() bool { return checkOrgExists }))
