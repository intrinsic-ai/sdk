// Copyright 2023 Intrinsic Innovation LLC

// Package clientutils provides utils for supporting catalog clients and authentication.
package clientutils

import (
	"context"
	"fmt"
	"io/fs"
	"regexp"
	"strings"


	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"intrinsic/assets/baseclientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/config/environments"
	clusterdiscoverypb "intrinsic/frontend/cloud/api/v1/clusterdiscovery_api_go_grpc_proto"
	solutiondiscoverygrpcpb "intrinsic/frontend/cloud/api/v1/solutiondiscovery_api_go_grpc_proto"
	solutiondiscoverypb "intrinsic/frontend/cloud/api/v1/solutiondiscovery_api_go_grpc_proto"
	"intrinsic/kubernetes/acl/identity"
	"intrinsic/tools/inctl/auth/auth"
)

const (
	defaultCatalogProject = "intrinsic-assets-prod"
)

var (
	catalogAssetAddressRegex = regexp.MustCompile(`(^|/)assets[-]?([^\.]*)\.intrinsic\.ai`)

	assetAddressToProjectSuffix = map[string]string{
		"":    "prod",
		"dev": "dev",
		"qa":  "staging",
	}

	errNoInfoToAuthenticate = errors.New("GCP project or api-key is required for authentication, but neither could be determined (try providing `--org=<org>@<project>` to specify the project)")
)

// DialClusterFromInctl creates a connection to a cluster from an inctl command.
func DialClusterFromInctl(ctx context.Context, flags *cmdutils.CmdFlags) (context.Context, *grpc.ClientConn, string, error) {
	project := flags.GetFlagProject()
	org := flags.GetFlagOrganization()
	address, cluster, solution, err := flags.GetFlagsAddressClusterSolution()
	if err != nil {
		return ctx, nil, "", err
	}

	// Only lookup cluster using the solution if cluster is not set.  We
	// probably shouldn't be in this scenario right now, but this isn't the
	// spot to enforce that.
	if solution != "" && cluster == "" {
		ctx, conn, _, err := dialConnectionCtx(ctx, dialInfoParams{
			Address: address,
			CredOrg: org,
			Project: project,
		})
		if err != nil {
			return ctx, nil, "", fmt.Errorf("could not create connection options for cluster: %v", err)
		}
		defer conn.Close()

		cluster, err = getClusterNameFromSolution(ctx, conn, solution)
		if err != nil {
			return ctx, nil, "", fmt.Errorf("could not get cluster name from solution: %v", err)
		}
	}

	ctx, conn, address, err := dialConnectionCtx(ctx, dialInfoParams{
		Address: address,
		Cluster: cluster,
		CredOrg: org,
		Project: project,
	})
	if err != nil {
		return ctx, nil, "", fmt.Errorf("could not create connection options for the installer: %v", err)
	}

	return ctx, conn, address, nil
}

// DialCatalogFromInctl creates a connection to an asset catalog service from an inctl command.
func DialCatalogFromInctl(cmd *cobra.Command, flags *cmdutils.CmdFlags) (context.Context, *grpc.ClientConn, error) {

	return DialCatalog(
		cmd.Context(), DialCatalogOptions{
			Address: "",
			APIKey: "",
			Org:          flags.GetFlagOrganization(),
			Project:      flags.GetFlagProject(),
		},
	)
}

// DialCatalogOptions specifies the options for DialCatalog.
type DialCatalogOptions struct {
	Address      string
	APIKey       string
	Org          string
	Project      string // Defaults to the global assets project.
}

// DialCatalog creates a connection to a asset catalog service.
func DialCatalog(ctx context.Context, opts DialCatalogOptions) (context.Context, *grpc.ClientConn, error) {
	catalogProject := ResolveCatalogProject(opts.Project)
	// Get the catalog address.
	address, err := resolveCatalogAddress(ctx, resolveCatalogAddressOptions{
		Address:      opts.Address,
		Project:      catalogProject,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("cannot resolve catalog address: %w", err)
	}

	optsOpts := getDialContextOptionsOptions{
		Address:                     address,
		APIKey:                      opts.APIKey,
		ClusterProject:              catalogProject,
		CredOrg:                     opts.Org,
		SkipCredsForInsecureAddress: false,
	}
	dialCtx, dialOpts, err := getDialContextOptions(ctx, optsOpts)
	if err != nil {
		// If we don't have enough info to authenticate, try one more time using the catalog project
		// as the credential project. (This supports mainly legacy usage where the project is specified
		// but not the org, and the user can authenticate in the catalog project.)
		if errors.Is(err, errNoInfoToAuthenticate) {
			optsOpts.CredProject = catalogProject
			var catalogProjectErr error
			if dialCtx, dialOpts, catalogProjectErr = getDialContextOptions(ctx, optsOpts); catalogProjectErr == nil {
				err = nil
			}
		}
		if err != nil {
			return nil, nil, fmt.Errorf("cannot get dial context options for catalog: %w", err)
		}
	}

	conn, err := grpc.NewClient(address, dialOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot dial catalog: %w", err)
	}
	return dialCtx, conn, nil
}

// ResolveCatalogProjectFromInctl returns the project to use for communicating with a catalog.
func ResolveCatalogProjectFromInctl(flags *cmdutils.CmdFlags) string {
	return ResolveCatalogProject(flags.GetFlagProject())
}

// ResolveCatalogProject returns the project to use for communicating with a catalog.
func ResolveCatalogProject(project string) string {
	if project == "" {
		return defaultCatalogProject
	}
	return project
}

// RemoteOpt returns the remote option to use for the given flags.
func RemoteOpt(flags *cmdutils.CmdFlags) (remote.Option, error) {
	authUser, authPwd := flags.GetFlagsRegistryAuthUserPassword()
	if len(authUser) != 0 && len(authPwd) != 0 {
		return remote.WithAuth(authn.FromConfig(authn.AuthConfig{
			Username: authUser,
			Password: authPwd,
		})), nil
	}
	return remote.WithAuthFromKeychain(google.Keychain), nil
}

type resolveCatalogAddressOptions struct {
	Address      string
	Project      string
}

func resolveCatalogAddress(ctx context.Context, opts resolveCatalogAddressOptions) (string, error) {
	// Check for user-provided address.
	if opts.Address != "" {
		return opts.Address, nil
	}

	// Derive the address from the project.
	if opts.Project == "" {
		return "", fmt.Errorf("cannot determine catalog address when project is empty")
	}
	address, err := getCatalogAddressForProject(ctx, opts)
	if err != nil {
		return "", err
	}

	addDNS := true
	if addDNS && !strings.HasPrefix(address, "dns:///") {
		address = fmt.Sprintf("dns:///%s", address)
	}

	return address, nil
}

func resolveTokenExchangeAddress(ctx context.Context, project string) (string, error) {
	var environment string
	var err error
	if environment, err = environments.FromProject(project); err != nil {
		environment = environments.FromComputeProject(project)
	}

	return environments.PortalDomain(environment), nil
}

func defaultGetCatalogAddressForProject(ctx context.Context, opts resolveCatalogAddressOptions) (address string, err error) {
	if opts.Project != "intrinsic-assets-prod" {
		return "", fmt.Errorf("unsupported project %s", opts.Project)
	}
	address = fmt.Sprintf("assets.intrinsic.ai:443")

	return address, nil
}

var (
	getCatalogAddressForProject = defaultGetCatalogAddressForProject
)

type getDialContextOptionsOptions struct {
	Address string
	APIKey  string
	// ClusterProject is the project in which the cluster being contacted is running. Used for
	// resolving the address that the cluster will use to contact the token exchange service.
	ClusterProject              string
	CredOrg                     string
	CredProject                 string
	SkipCredsForInsecureAddress bool
}

// getDialContextOptions returns metadata for dialing a gRPC connection to a cloud/on-prem cluster.
//
// It uses the provided ctx to manage the lifecycle of connection created. ctx may be modified, so
// the caller should use the returned context instead.
func getDialContextOptions(ctx context.Context, opts getDialContextOptionsOptions) (context.Context, []grpc.DialOption, error) {
	credProject := opts.CredProject
	if opts.CredOrg != "" {
		info, err := getOrgInfo(opts.CredOrg)
		if err != nil {
			return nil, nil, err
		}
		// If the credential project has not been specified, use the project from the org.
		if credProject == "" {
			credProject = info.Project
		}
		ctx, err = identity.OrgToContext(ctx, info.Organization)
		if err != nil {
			return nil, nil, err
		}
	}

	dialOpts := baseclientutils.BaseDialOptions()

	creds, err := getPerRPCCredentials(ctx, getPerRPCCredentialsOptions{
		Address:                     opts.Address,
		APIKey:                      opts.APIKey,
		ClusterProject:              opts.ClusterProject,
		CredProject:                 credProject,
		SkipCredsForInsecureAddress: opts.SkipCredsForInsecureAddress,
	})
	if err != nil {
		return nil, nil, err
	}
	if creds != nil {
		dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(creds))
	}

	var tcOption grpc.DialOption
	if baseclientutils.UseInsecureCredentials(opts.Address) {
		tcOption = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		var err error
		tcOption, err = baseclientutils.GetTransportCredentialsDialOption()
		if err != nil {
			return nil, nil, fmt.Errorf("cannot get transport credentials: %w", err)
		}
	}
	dialOpts = append(dialOpts, tcOption)

	return ctx, dialOpts, nil
}

func getOrgInfo(org string) (auth.OrgInfo, error) {
	// Try first to parse the org and project directly from the org string, so we don't read the
	// org info from disk unless we need to.
	orgAndProject := strings.Split(org, "@")
	if len(orgAndProject) == 2 {
		return auth.OrgInfo{
			Organization: orgAndProject[0],
			Project:      orgAndProject[1],
		}, nil
	}

	info, err := auth.NewStore().ReadOrgInfo(org)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return auth.OrgInfo{Organization: org}, nil
		}
		return auth.OrgInfo{}, err
	}
	return info, nil
}

type getPerRPCCredentialsOptions struct {
	Address string
	APIKey  string
	// ClusterProject is the project in which the cluster being contacted is running. Used for
	// resolving the address that the cluster will use to contact the token exchange service.
	ClusterProject string
	CredProject    string
	// SkipCredsForInsecureAddress is a flag to skip adding credentials to the connection if the
	// address is a candidate for an insecure connection.
	SkipCredsForInsecureAddress bool
}

func getPerRPCCredentials(ctx context.Context, opts getPerRPCCredentialsOptions) (credentials.PerRPCCredentials, error) {
	addressIsInsecure := baseclientutils.UseInsecureCredentials(opts.Address) || baseclientutils.IsLocalAddress(opts.Address)
	if addressIsInsecure && opts.SkipCredsForInsecureAddress {
		return nil, nil
	}

	if opts.APIKey != "" {
		// User-provided api-key.
		return &auth.ProjectToken{APIKey: opts.APIKey}, nil
	}

	if opts.CredProject != "" {
		configuration, err := auth.NewStore().GetConfiguration(opts.CredProject)
		if err != nil {
			return nil, err
		}

		projectToken, err := configuration.GetDefaultCredentials()
		if err != nil {
			return nil, err
		}

		// Convert the API token to a Google ID token. With this conversion, the JWT is included in the
		// credentials provided by the client (rather than requiring the auth-proxy to do it on
		// ingress.)
		tokenExchangeServiceAddress, err := resolveTokenExchangeAddress(ctx, opts.ClusterProject)
		if err != nil {
			return nil, fmt.Errorf("cannot resolve token exchange service address: %w", err)
		}
		opts := []auth.AsIDTokenCredentialsOption{
			auth.WithTokenExchangeServiceAddress(tokenExchangeServiceAddress),
		}
		if addressIsInsecure {
			opts = append(opts, auth.WithAPIKeyTokenSourceOptions(auth.WithAllowInsecure()))
		}

		return projectToken.AsIDTokenCredentials(opts...)
	}

	if addressIsInsecure {
		return nil, nil
	}

	return nil, errNoInfoToAuthenticate
}

type dialInfoParams struct {
	Address   string // The address of a cloud/on-prem cluster
	Cluster   string // The name of the server to install to
	CredOrg   string // Optional the org-id header to set
	CredToken string // Optional the credential value itself. This bypasses the store
	Project   string // The current GCP project
}

func dialConnectionCtx(ctx context.Context, params dialInfoParams) (context.Context, *grpc.ClientConn, string, error) {
	address, err := resolveClusterAddress(params.Address, params.Project)
	if err != nil {
		return ctx, nil, "", err
	}

	ctx, dialOpts, err := getDialContextOptions(ctx, getDialContextOptionsOptions{
		Address:                     address,
		APIKey:                      params.CredToken,
		ClusterProject:              params.Project,
		CredOrg:                     params.CredOrg,
		CredProject:                 params.Project,
		SkipCredsForInsecureAddress: true,
	})
	if err != nil {
		return ctx, nil, "", fmt.Errorf("cannot get dial context options for cluster: %w", err)
	}

	if params.Cluster != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-server-name", params.Cluster)
	}

	conn, err := grpc.NewClient(address, dialOpts...)
	if err != nil {
		return ctx, nil, "", fmt.Errorf("could not dial context: %w", err)
	}

	return ctx, conn, address, nil
}

// AuthInsecureConn returns a context with authentication information if the address is insecure.
func AuthInsecureConn(ctx context.Context, address string, project string) context.Context {
	authCtx := ctx

	return authCtx
}

func resolveClusterAddress(address string, project string) (string, error) {
	if address != "" {
		return address, nil
	}

	if project == "" {
		return "", fmt.Errorf("project is required if no address is specified")
	}

	return fmt.Sprintf("dns:///www.endpoints.%s.cloud.goog:443", project), nil
}

// getClusterNameFromSolution returns the cluster in which a solution currently runs.
func getClusterNameFromSolution(ctx context.Context, conn *grpc.ClientConn, solutionName string) (string, error) {
	solution, err := getSolution(ctx, conn, solutionName)
	if err != nil {
		return "", fmt.Errorf("failed to get solution: %w", err)
	}
	if solution.GetState() == clusterdiscoverypb.SolutionState_SOLUTION_STATE_NOT_RUNNING {
		return "", fmt.Errorf("solution is not running")
	}
	if solution.GetClusterName() == "" {
		return "", fmt.Errorf("unknown error: solution is running but cluster is empty")
	}
	return solution.GetClusterName(), nil
}

// getSolution gets solution data by name
func getSolution(ctx context.Context, conn *grpc.ClientConn, solutionName string) (*solutiondiscoverypb.SolutionDescription, error) {
	client := solutiondiscoverygrpcpb.NewSolutionDiscoveryServiceClient(conn)
	req := &solutiondiscoverypb.GetSolutionDescriptionRequest{Name: solutionName}
	resp, err := client.GetSolutionDescription(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get solution description: %w", err)
	}

	return resp.GetSolution(), nil
}
