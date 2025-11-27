// Copyright 2023 Intrinsic Innovation LLC

package auth

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"slices"
	"strings"

	env "intrinsic/config/environments"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/orgutil"
	"intrinsic/tools/inctl/util/viperutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"

	accdiscoverv1grpcpb "intrinsic/kubernetes/accounts/service/api/v1/discoveryapi_go_grpc_proto"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
	keyNoBrowser = "no_browser"

	orgTokenURLFmt     = "https://%s/o/%s/generate-keys"
	projectTokenURLFmt = "https://%s/project/%s/generate-keys"
	// We are going to use system defaults to ensure we open web-url correctly.
	// For dev container running via VS Code the sensible-browser redirects
	// call into code client from server to ensure URL is opened in valid
	// client browser.
	sensibleBrowser = "/usr/bin/sensible-browser"
)

// Exposed for testing
var (
	queryProjects = queryProjectsForAPIKey
)

var (
	loginParams = viper.New()
	loginCmd    = orgutil.WrapCmd(
		&cobra.Command{
			Use:   "login",
			Short: "Logs in user into Flowstate",
			Long:  "Logs in user into Flowstate to allow interactions with solutions.",
			Args:  cobra.NoArgs,
			RunE:  loginCmdE,
			PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
				if err := orgutil.ValidateEnvironment(loginParams); err != nil {
					return err
				}

				return nil
			},
		},
		loginParams,
		orgutil.WithOrgExistsCheck(func() bool {
			// The login command only creates the org, so we must disable the flag check which happens before.
			return false
		}),
	)
)

func readAPIKeyFromPipe(reader *bufio.Reader) (string, error) {
	fi, _ := os.Stdin.Stat()
	// Check if input comes from pipe. Taken from
	// https://www.socketloop.com/tutorials/golang-check-if-os-stdin-input-data-is-piped-or-from-terminal
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			return "", err
		}

		return strings.TrimSpace(string(bytes)), nil
	}
	return "", nil
}

func queryForAPIKey(ctx context.Context, writer io.Writer, in *bufio.Reader, organization, project string) (string, error) {
	e := loginParams.GetString(orgutil.KeyEnvironment)
	if e == "" {
		e = env.FromComputeProject(project)
	}
	portal := env.PortalDomain(e)
	if portal == "" {
		return "", fmt.Errorf("unknown environment %q", e)
	}
	authorizationURL := fmt.Sprintf(projectTokenURLFmt, portal, project)
	if organization != "" {
		authorizationURL = fmt.Sprintf(orgTokenURLFmt, portal, url.PathEscape(organization))
	}
	fmt.Fprintf(writer, "Open URL in your browser to obtain authorization token: %s\n", authorizationURL)

	ignoreBrowser := loginParams.GetBool(keyNoBrowser)
	if !ignoreBrowser {
		_, _ = fmt.Fprintln(writer, "Attempting to open URL in your browser...")
		browser := exec.CommandContext(ctx, sensibleBrowser, authorizationURL)
		browser.Stdout = io.Discard
		browser.Stderr = io.Discard
		if err := browser.Start(); err != nil {
			fmt.Fprintf(writer, "Failed to open URL in your browser, please run command again with '--%s'.\n", keyNoBrowser)
			return "", fmt.Errorf("rerun with '--%s', got error %w", keyNoBrowser, err)
		}
	}
	fmt.Fprintf(writer, "\nPaste access token from website: ")

	apiKey, err := in.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("cannot read from input device: %w", err)
	}

	// Move the cursor back to the beginning of the line and clear the line
	fmt.Fprintf(writer, "\033[1A\033[2K")
	// Overwrite the line with a placeholder
	fmt.Fprintf(writer, "Paste access token from website: [redacted]\n")

	return strings.TrimSpace(apiKey), nil
}

// queryProjectsForAPIKey discovers the projects the given API key has access to.
// If optionalOrg is set, it will be used as a filter to only return projects the given organization
// is part of.
func queryProjectsForAPIKey(ctx context.Context, apiKey string, optionalOrg string) ([]string, error) {
	e := loginParams.GetString(orgutil.KeyEnvironment)
	if e == "" {
		e = env.Prod
	}
	accProject := env.AccountsProjectFromEnv(e)
	conn, err := auth.NewCloudConnection(ctx, auth.WithProject(accProject), auth.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := accdiscoverv1grpcpb.NewAccountsDiscoveryServiceClient(conn)
	resp, err := client.ListOrganizations(ctx, &emptypb.Empty{})
	if err != nil {
		fmt.Println("Could not find the project for this token. Please restart the login process and make sure to provide the exact key shown by the portal.")
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	// multiple organizations are fine, but they must be all on the same project
	orgs := resp.GetOrganizations()
	projects := map[string]struct{}{}
	for _, org := range orgs {
		// filter by org if specified
		if optionalOrg != "" && optionalOrg != org.GetName() {
			continue
		}
		projects[org.GetProject()] = struct{}{}
	}
	return maps.Keys(projects), nil
}

func loginCmdE(cmd *cobra.Command, _ []string) (err error) {
	writer := cmd.OutOrStdout()
	projectName := loginParams.GetString(orgutil.KeyProject)
	orgName := loginParams.GetString(orgutil.KeyOrganization)
	org := orgutil.QualifiedOrg(projectName, orgName)
	in := bufio.NewReader(cmd.InOrStdin())
	// In the future multiple aliases should be supported for one project.
	alias := auth.AliasDefaultToken
	isBatch := loginParams.GetBool(keyBatch)

	apiKey, err := readAPIKeyFromPipe(in)
	if err != nil {
		return err
	}

	if apiKey == "" && !isBatch {
		apiKey, err = queryForAPIKey(cmd.Context(), writer, in, org, projectName)
		if err != nil {
			return err
		}
	}

	if apiKey == "" {
		return fmt.Errorf("API key is empty. Please provide an API key")
	}

	// If we are passed a pure org, we don't know the project yet
	if projectName == "" {
		projects, err := queryProjects(cmd.Context(), apiKey, orgName)
		if err != nil {
			return fmt.Errorf("query project: %w", err)
		}
		if len(projects) == 0 {
			return fmt.Errorf("no project found for API key. Please double check the value of the -org flag for typos" +
				".")
		}
		if len(projects) > 1 {
			slices.Sort(projects)
			return fmt.Errorf("multiple projects found for API key (and org %q): %+v", org, projects)
		}
		// exactly one found
		projectName = projects[0]
	}
	if org != "" {
		if err := authStore.WriteOrgInfo(&auth.OrgInfo{Organization: org, Project: projectName}); err != nil {
			return fmt.Errorf("store org info: %w", err)
		}
	}

	var config *auth.ProjectConfiguration
	if authStore.HasConfiguration(projectName) {
		if config, err = authStore.GetConfiguration(projectName); err != nil {
			return fmt.Errorf("cannot load '%s' configuration: %w", projectName, err)
		}
	} else {
		config = auth.NewConfiguration(projectName)
	}

	config, err = config.SetCredentials(alias, apiKey)
	if err != nil {
		return fmt.Errorf("aborting, invalid credentials: %w", err)
	}

	_, err = authStore.WriteConfiguration(config)

	return err
}

func init() {
	authCmd.AddCommand(loginCmd)

	flags := loginCmd.Flags()
	// we will use viper to fetch data, we do not need local variables
	flags.Bool(keyNoBrowser, false, "Disables attempt to open login URL in browser automatically")
	flags.Bool(keyBatch, false, "Suppresses command prompts and assume Yes or default as an answer. Use with shell scripts.")
	flags.String(orgutil.KeyEnvironment, "", fmt.Sprintf("Auth environment to use. This should be one of %v. %q is used by default. See http://go/intrinsic-users#environments for the compatible environment corresponding to a cloud project.", strings.Join(env.All, ", "), env.Prod))
	flags.MarkHidden(orgutil.KeyEnvironment)
	flags.MarkHidden(orgutil.KeyProject)

	viperutil.BindFlags(loginParams, flags, viperutil.BindToListEnv(orgutil.KeyEnvironment))
}
