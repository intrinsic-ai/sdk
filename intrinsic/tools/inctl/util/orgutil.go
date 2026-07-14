// Copyright 2023 Intrinsic Innovation LLC

// Package orgutil provides common utility to handle projects/organizations in inctl.
package orgutil

import (
	"fmt"
	env "intrinsic/config/environments"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/viperutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
const (
	// KeyProject is used as central flag name for passing a project name to inctl.
	KeyProject = auth.KeyProject
	// KeyOrganization is used as central flag name for passing an organization name to inctl.
	KeyOrganization = auth.KeyOrganization
	// KeyEnvironment is used as central flag name for passing an environment name to inctl.
	//
	// The environment can be one of prod, staging or dev.
	KeyEnvironment = auth.KeyEnvironment
)

var (
	// Exposed for testing
	authStore        = auth.NewStore()
	errNoOrg         = fmt.Errorf("expected --%s=<org>", KeyOrganization)
	errOrgAndProject = fmt.Errorf("do not set --%s, use --%s=<org> or <org@project> instead", KeyProject, KeyOrganization)
)

// ErrOrgNotFound indicates that the lookup for a given credential
// name failed.
type ErrOrgNotFound struct {
	err           error
	CandidateOrgs []string
	OrgName       string
}

func (e *ErrOrgNotFound) Error() string {
	return fmt.Sprintf("credentials not found: %q", e.OrgName)
}

func (e *ErrOrgNotFound) Unwrap() error {
	return e.err
}

func editDistance(left, right string) int {
	length := len([]rune(right))
	if length == 0 {
		return len([]rune(left))
	}

	dist1 := make([]int, length+1)
	dist2 := make([]int, length+1)

	// initialize dist1 (the previous row of distances)
	// this row is A[0][i]: edit distance from an empty left to right;
	// that distance is the number of characters to append to  left to make right.
	for i := 0; i < length+1; i++ {
		dist1[i] = i
		dist2[i] = 0
	}

	for i, vLeft := range []rune(left) {
		dist2[0] = i + 1

		for j, vRight := range []rune(right) {
			deletionCost := dist1[j+1] + 1
			insertionCost := dist2[j] + 1
			var substitutionCost int
			if vLeft == vRight {
				substitutionCost = dist1[j]
			} else {
				substitutionCost = dist1[j] + 1
			}

			// dist2[j + 1] = min(insertionCost, deletionCost, substitutionCost)
			if deletionCost <= insertionCost && deletionCost <= substitutionCost {
				dist2[j+1] = deletionCost
			} else if insertionCost <= deletionCost && insertionCost <= substitutionCost {
				dist2[j+1] = insertionCost
			} else {
				dist2[j+1] = substitutionCost
			}
		}

		copy(dist1, dist2)
	}

	return dist1[length]
}

func makeOrgNotFound(inner error, org string) error {
	candidates := []string{}
	orgs, err := auth.NewStore().ListOrgs()
	// We can only do this, if there's NO error!
	if err == nil {
		hasAt := strings.Contains(org, "@")
		for _, candidate := range orgs {
			target := candidate
			if !hasAt {
				parts := strings.Split(candidate, "@")
				target = parts[0]
			}
			if editDistance(org, target) < 3 {
				candidates = append(candidates, candidate)
			}
		}
	}

	return &ErrOrgNotFound{err: inner, CandidateOrgs: candidates, OrgName: org}
}

// ValidateEnvironment validates the environment value in a cobra command.
func ValidateEnvironment(vipr *viper.Viper) error {
	e := vipr.GetString(KeyEnvironment)

	switch e {
	case env.Prod, env.Staging, env.Dev, "":
		// Valid environments
	default:
		return fmt.Errorf("invalid --%s value %q. It must be one of %v", KeyEnvironment, e, strings.Join(env.All, ", "))
	}

	// If a project is explicitly set, check if it's a known central project
	// and reject contradicting explicit environment values.
	project := vipr.GetString(KeyProject)
	if e != "" && project != "" {
		if inferred := env.FromAnyProject(project); inferred != e && env.IsKnownProject(project) {
			return fmt.Errorf("environment mismatch: explicitly provided environment %q contradicts environment %q inferred from project name %q", e, inferred, project)
		}
	}

	return nil
}

// ResolveOrg searches stored organization credentials to match the given short name.
//
// If a user inputs a short organization name (e.g. `--org=my-org`), this function will
// look up all credentials in the local auth store. If it finds a unique matching
// fully-qualified organization (e.g. `my-org@my-project`), it will return the OrgInfo.
//
// Returns:
//   - auth.OrgInfo: The unique matching organization info, if found.
//   - error: ErrOrgNotFound if no matching organization exists; an error listing the
//     matching fully-qualified candidates if the short name is ambiguous (matches
//     multiple stored credentials); or any filesystem read error.
func ResolveOrg(requestedOrg string) (auth.OrgInfo, error) {
	// First, try a direct lookup by the given org name.
	if info, err := authStore.ReadOrgInfo(requestedOrg); err == nil {
		return info, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return auth.OrgInfo{}, err
	}

	// Try to resolve it as a short name by searching stored credentials.
	orgs, err := authStore.ListOrgs()
	if err != nil {
		return auth.OrgInfo{}, fmt.Errorf("list orgs: %w", err)
	}

	reqShort := requestedOrg
	var reqProject string
	if strings.Contains(requestedOrg, "@") {
		parts := strings.Split(requestedOrg, "@")
		reqShort = parts[0]
		reqProject = parts[1]
	}

	var matches []auth.OrgInfo
	for _, candidate := range orgs {
		parts := strings.Split(candidate, "@")
		if parts[0] == reqShort {
			info, err := authStore.ReadOrgInfo(candidate)
			if err == nil {
				if reqProject == "" || info.Project == reqProject {
					matches = append(matches, info)
				}
			}
		}
	}

	if len(matches) == 0 {
		return auth.OrgInfo{}, makeOrgNotFound(os.ErrNotExist, requestedOrg)
	}
	if len(matches) > 1 {
		var options []string
		for _, m := range matches {
			options = append(options, m.Organization)
		}
		return auth.OrgInfo{}, fmt.Errorf("organization %q is ambiguous. Please specify one of the following fully-qualified organizations: %s", requestedOrg, strings.Join(options, ", "))
	}
	return matches[0], nil
}

// preRunOrganizationOptional provides the organization/project flag handling as PersistentPreRunE
// of a cobra command. This is done automatically with the WrapCmdOptional() function.
//
// However, it lets the user skip setting --org in case they prefer --context with a local context /
// alias.
func preRunOrganizationOptional(cmd *cobra.Command, vipr *viper.Viper, enableOrgExistsCheck func() bool) error {
	if err := validateProjectOrg(cmd, vipr, enableOrgExistsCheck); err != nil {
		return fmt.Errorf("error validating project/org: %w", err)
	}

	if err := ValidateEnvironment(vipr); err != nil {
		return fmt.Errorf("error validating environment: %w", err)
	}

	return nil
}

func validateProjectOrg(cmd *cobra.Command, vipr *viper.Viper, enableOrgExistsCheck func() bool) error {
	projectFlag := cmd.PersistentFlags().Lookup(KeyProject)
	orgFlag := cmd.PersistentFlags().Lookup(KeyOrganization)

	org := vipr.GetString(KeyOrganization)
	project := vipr.GetString(KeyProject)

	if project != "" && org != "" {
		return errOrgAndProject
	}

	if org == "" {
		// When using --project, the org is unknown and no further logic is required.
		return nil
	}

	// Cleanup the org parameter, it could be org@project.
	// The full name is only required to lookup the correct project. So we can clean it up here
	parts := strings.Split(org, "@")

	orgFlag.Value.Set(parts[0])
	vipr.Set(KeyOrganization, parts[0])
	if len(parts) > 1 {
		projectFlag.Value.Set(parts[1])
		vipr.Set(KeyProject, parts[1])
	}

	if !enableOrgExistsCheck() {
		return nil
	}

	// Resolve the organization (either direct lookup or fallback to short name matching).
	resolved, err := ResolveOrg(org)
	if err != nil {
		return fmt.Errorf("error resolving org: %w", err)
	}

	projectFlag.Value.Set(resolved.Project)
	vipr.Set(KeyProject, resolved.Project)

	return nil
}

// preRunOrganization checks organization/project flags as PersistentPreRunE of a cobra command.
// This is done automatically with the WrapCmd() function. preRunOrganization() doesn't call
// preRunOrganizationOptional() itself.
//
// It enforces that exactly one of --project or --org is set.
func preRunOrganization(cmd *cobra.Command, vipr *viper.Viper) (noOrg bool, err error) {
	org := vipr.GetString(KeyOrganization)
	project := vipr.GetString(KeyProject)

	if project == "" && org == "" {
		return false, errNoOrg
	}
	if org == "" {
		fmt.Fprintf(os.Stderr, "\ninctl was called without an organization. This is deprecated and will soon be an error. Please use --org intrinsic@%v.\n", project)
		return true, nil
	}

	return false, nil
}

// WrapCmdOption is an interface for options that configure the behavior of
// WrapCmd and WrapCmdOptional.
type WrapCmdOption interface {
	set(*wrapCmdOptions)
}

type wrapCmdOptions struct {
	enableOrgExistsCheck func() bool
}

type wrapCmdOption func(*wrapCmdOptions)

func (o wrapCmdOption) set(opts *wrapCmdOptions) {
	o(opts)
}

// WithOrgExistsCheck provides a way to selectively enable or disable the org
// credentials check in a command.  It's provided as a closure to allow it to
// be evaluated after flag parsing has occurred.
func WithOrgExistsCheck(enable func() bool) WrapCmdOption {
	return wrapCmdOption(func(opts *wrapCmdOptions) {
		opts.enableOrgExistsCheck = enable
	})
}

// WrapCmdOptional injects KeyProject and KeyOrganization as PersistentFlags into the command and
// sets up shared handling for them.
//
// However, it lets the user skip setting --org in case they prefer --context with a local context /
// alias.
func WrapCmdOptional(cmd *cobra.Command, vipr *viper.Viper, options ...WrapCmdOption) *cobra.Command {
	opts := &wrapCmdOptions{
		enableOrgExistsCheck: func() bool { return true },
	}
	for _, option := range options {
		option.set(opts)
	}
	cmd.PersistentFlags().StringP(KeyProject, "p", "",
		`The Google Cloud Project (GCP) project to use. You can set the environment variable
	INTRINSIC_PROJECT=project_name to set a default project name.`)
	cmd.PersistentFlags().StringP(KeyOrganization, "", "",
		`The Intrinsic organization to use. You can set the environment variable
	INTRINSIC_ORG=organization to set a default organization.`)
	envUsage := fmt.Sprintf("Auth environment to use. This should be one of %v. The environment is automatically inferred from the project by default, or explicitly defaults to %q if the project is unknown. "+
		"Each cloud project is associated with exactly one environment.", strings.Join(env.All, ", "), env.Prod)
	cmd.PersistentFlags().String(KeyEnvironment, "", envUsage)
	_ = cmd.PersistentFlags().MarkHidden(KeyEnvironment)
	oldPreRunE := cmd.PersistentPreRunE
	cmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		// This is required to cooperate with cobrautil.
		// Cobrautil's way to force an error instead of 0 return code when there's no subcommand
		// causes cobra to run the PersistentPreRunE either way. So we need to short-circuit
		// the flag handling here.
		if !c.DisableFlagParsing {
			if err := preRunOrganizationOptional(cmd, vipr, opts.enableOrgExistsCheck); err != nil {
				return err
			}
		}

		if oldPreRunE != nil {
			return oldPreRunE(c, args)
		}
		return nil
	}

	viperutil.BindFlags(vipr, cmd.PersistentFlags(), viperutil.BindToListEnv(KeyOrganization, KeyEnvironment))

	return cmd
}

// WrapCmd injects KeyProject, KeyOrganization and KeyEnvironment as PersistentFlags into the command and sets up
// shared handling for them.
//
// It enforces that exactly one of --project or --org is set.
func WrapCmd(cmd *cobra.Command, vipr *viper.Viper, options ...WrapCmdOption) *cobra.Command {
	cmd = WrapCmdOptional(cmd, vipr, options...)

	var noOrg bool
	oldPreRunE := cmd.PersistentPreRunE
	cmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		// This is required to cooperate with cobrautil.
		// Cobrautil's way to force an error instead of 0 return code when there's no subcommand
		// causes cobra to run the PersistentPreRunE either way. So we need to short-circuit
		// the flag handling here.
		if !c.DisableFlagParsing {
			var err error
			noOrg, err = preRunOrganization(cmd, vipr)
			if err != nil {
				return err
			}
		}

		if oldPreRunE != nil {
			return oldPreRunE(c, args)
		}
		return nil
	}
	oldPostRunE := cmd.PersistentPostRunE
	cmd.PersistentPostRunE = func(c *cobra.Command, args []string) error {
		if noOrg {
			fmt.Fprintf(os.Stderr, "\ninctl was called without an organization. This is deprecated and will soon be an error. Please use --org.\n")
		}

		if oldPostRunE != nil {
			return oldPostRunE(c, args)
		}
		return nil
	}

	return cmd
}

// QualifiedOrg returns a "unique" org name, adding an @project suffix for orgs that are present in
// multiple projects. This undoes the "cleaning" applied by preRunOrganization when using WrapCmd().
func QualifiedOrg(projectName, orgName string) string {
	if orgName == "" { // fallback, not sure if this is really required
		return fmt.Sprintf("intrinsic@%s", projectName)
	}
	// for most customer organizations there is no project required
	if projectName == "" {
		return orgName
	}
	// for organizations with multiple projects
	return fmt.Sprintf("%s@%s", orgName, projectName)
}
