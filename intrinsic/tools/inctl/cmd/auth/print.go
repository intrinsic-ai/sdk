// Copyright 2023 Intrinsic Innovation LLC

package auth

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/orgutil"
)

var (
	flagFlowstateAddr string
)

func init() {
	authCmd.AddCommand(printAPIKeyCmd)

	authCmd.AddCommand(printAccessTokenCmd)
	printAccessTokenCmd.Flags().StringVar(&flagFlowstateAddr, "flowstate", "flowstate.intrinsic.ai", "Flowstate address.")
}

func projectFromFlags(cmd *cobra.Command) (string, error) {
	org := viperLocal.GetString(orgutil.KeyOrganization)
	project := viperLocal.GetString(orgutil.KeyProject)
	if project != "" {
		return project, nil
	}
	if org != "" {
		orgInfo, err := authStore.ReadOrgInfo(org)
		if err != nil {
			return "", fmt.Errorf("failed to get organization info for %q: %v", org, err)
		}
		return orgInfo.Project, nil
	}
	return "", fmt.Errorf("either --%s or --%s is required", orgutil.KeyOrganization, orgutil.KeyProject)
}

var printAPIKeyCmd = &cobra.Command{
	Use:   "print-api-key",
	Short: "Prints the API key for a project.",
	Long:  "Prints the API key for a project.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, err := projectFromFlags(cmd)
		if err != nil {
			return fmt.Errorf("failed to determine project: %v", err)
		}
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
}

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

var printAccessTokenCmd = &cobra.Command{
	Use:   "print-access-token",
	Short: "Print an access token.",
	Long:  printAccessTokenHelp,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, err := projectFromFlags(cmd)
		if err != nil {
			return fmt.Errorf("failed to determine project: %v", err)
		}
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
		cmd.Printf("%s", resp.IDToken)
		return nil
	},
}
