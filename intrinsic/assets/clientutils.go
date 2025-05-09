// Copyright 2023 Intrinsic Innovation LLC

// Package clientutils provides utils for supporting catalog clients and authentication.
package clientutils

import (
	"context"
	"fmt"
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
	clusterdiscoverypb "intrinsic/frontend/cloud/api/v1/clusterdiscovery_api_go_grpc_proto"
	solutiondiscoverygrpcpb "intrinsic/frontend/cloud/api/v1/solutiondiscovery_api_go_grpc_proto"
	solutiondiscoverypb "intrinsic/frontend/cloud/api/v1/solutiondiscovery_api_go_grpc_proto"
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
)

// DialCatalogOptions specifies the options for DialCatalog.
type DialCatalogOptions struct {
	Address      string
	APIKey       string
	Project      string // Defaults to the global assets project.
}

// DialClusterFromInctl creates a connection to a cluster from an inctl command.
func DialClusterFromInctl(ctx context.Context, flags *cmdutils.CmdFlags) (context.Context, *grpc.ClientConn, string, error) {
	project := flags.GetFlagProject()
	org := flags.GetFlagOrganization()
	address, cluster, solution, err := flags.GetFlagsAddressClusterSolution()
	if err != nil {
		return ctx, nil, "", err
	}

	if solution != "" {
		ctx, conn, _, err := dialConnectionCtx(ctx, dialInfoParams{
			Address:  address,
			CredName: project,
			CredOrg:  org,
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
		Address:  address,
		Cluster:  cluster,
		CredName: project,
		CredOrg:  org,
	})
	if err != nil {
		return ctx, nil, "", fmt.Errorf("could not create connection options for the installer: %v", err)
	}

	return ctx, conn, address, nil
}

// DialCatalogFromInctl creates a connection to an asset catalog service from an inctl command.
func DialCatalogFromInctl(cmd *cobra.Command, flags *cmdutils.CmdFlags) (*grpc.ClientConn, error) {

	return DialCatalog(
		cmd.Context(), DialCatalogOptions{
			Address:      "",
			APIKey: "",
			Project:      ResolveCatalogProjectFromInctl(flags),
		},
	)
}

// DialCatalog creates a connection to a asset catalog service.
func DialCatalog(ctx context.Context, opts DialCatalogOptions) (*grpc.ClientConn, error) {
	opts.Project = ResolveCatalogProject(opts.Project)

	// Get the catalog address.
	address, err := resolveCatalogAddress(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "cannot resolve address")
	}

	options := baseclientutils.BaseDialOptions()

	if baseclientutils.IsLocalAddress(opts.Address) { // Use insecure creds.
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else { // Use api-key creds.
		rpcCreds, err := getAPIKeyPerRPCCredentials(opts.APIKey, opts.Project)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get api-key credentials")
		}
		tcOption, err := baseclientutils.GetTransportCredentialsDialOption()
		if err != nil {
			return nil, errors.Wrap(err, "cannot get transport credentials")
		}
		options = append(options, grpc.WithPerRPCCredentials(rpcCreds), tcOption)
	}

	return grpc.DialContext(ctx, address, options...)
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

func resolveCatalogAddress(ctx context.Context, opts DialCatalogOptions) (string, error) {
	// Check for user-provided address.
	if opts.Address != "" {
		return opts.Address, nil
	}

	// Derive the address from the project.
	if opts.Project == "" {
		return "", fmt.Errorf("project is empty")
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

func defaultGetCatalogAddressForProject(ctx context.Context, opts DialCatalogOptions) (address string, err error) {
	if opts.Project != "intrinsic-assets-prod" {
		return "", fmt.Errorf("unsupported project %s", opts.Project)
	}
	address = fmt.Sprintf("assets.intrinsic.ai:443")

	return address, nil
}

var (
	getCatalogAddressForProject = defaultGetCatalogAddressForProject
)

// getAPIKeyPerRPCCredentials returns api-key PerRPCCredentials.
func getAPIKeyPerRPCCredentials(apiKey string, project string) (credentials.PerRPCCredentials, error) {
	var token *auth.ProjectToken

	if apiKey != "" {
		// User-provided api-key.
		token = &auth.ProjectToken{APIKey: apiKey}
	} else {
		// Load api-key from the auth store.
		configuration, err := auth.NewStore().GetConfiguration(project)
		if err != nil {
			return nil, err
		}

		token, err = configuration.GetDefaultCredentials()
		if err != nil {
			return nil, err
		}
	}

	return token, nil
}

type dialInfoParams struct {
	Address   string // The address of a cloud/on-prem cluster
	Cluster   string // The name of the server to install to
	CredName  string // The name of the credentials to load from auth.Store
	CredAlias string // Optional alias for key to load
	CredOrg   string // Optional the org-id header to set
	CredToken string // Optional the credential value itself. This bypasses the store
}

func dialConnectionCtx(ctx context.Context, params dialInfoParams) (context.Context, *grpc.ClientConn, string, error) {

	ctx, dialerOpts, address, err := dialInfoCtx(ctx, params)
	if err != nil {
		return nil, nil, "", fmt.Errorf("dial info: %w", err)
	}

	conn, err := grpc.DialContext(ctx, address, *dialerOpts...)
	if err != nil {
		return nil, nil, "", fmt.Errorf("dialing context: %w", err)
	}

	return ctx, conn, address, nil
}

// dialInfoCtx returns the metadata for dialing a gRPC connection to a cloud/on-prem cluster.
//
// Function uses provided ctx to manage lifecycle of connection created. Ctx may be
// modified on return, caller is encouraged to switch to returned context if appropriate.
//
// DialInfoParams.Cluster optionally has to be set to the name of the target cluster if
// DialInfoParams.Address is the address of a cloud cluster and the connection will be used to send
// a request to an on-prem service via the relay running in the cloud cluster.
//
// Returns insecure connection data if the address is a local network address (such as
// `localhost:17080`), otherwise retrieves cert from system cert pool, and sets up the metadata for
// a TLS cert with per-RPC basic auth credentials.
func dialInfoCtx(ctx context.Context, params dialInfoParams) (context.Context, *[]grpc.DialOption, string, error) {
	address, err := resolveClusterAddress(params.Address, params.CredName)
	if err != nil {
		return ctx, nil, "", err
	}
	params.Address = address

	if params.CredOrg != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, auth.OrgIDHeader, strings.Split(params.CredOrg, "@")[0])
	}

	if baseclientutils.UseInsecureCredentials(params.Address) {
		finalOpts := append(baseclientutils.BaseDialOptions(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		return ctx, &finalOpts, params.Address, nil
	}

	if params.Cluster != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-server-name", params.Cluster)
	}

	rpcCredentials, err := createCredentials(params)
	if err != nil {
		return nil, nil, "", fmt.Errorf("cannot retrieve connection credentials: %w", err)
	}
	tcOption, err := baseclientutils.GetTransportCredentialsDialOption()
	if err != nil {
		return nil, nil, "", fmt.Errorf("cannot retrieve transport credentials: %w", err)
	}

	finalOpts := append(baseclientutils.BaseDialOptions(),
		grpc.WithPerRPCCredentials(rpcCredentials),
		tcOption,
	)

	return ctx, &finalOpts, params.Address, nil
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

func createCredentials(params dialInfoParams) (credentials.PerRPCCredentials, error) {
	if params.CredToken != "" {
		return &auth.ProjectToken{APIKey: params.CredToken}, nil
	}

	if params.CredName != "" {
		configuration, err := auth.NewStore().GetConfiguration(params.CredName)
		if err != nil {
			return nil, fmt.Errorf("credentials not found: %w", err)
		}

		if params.CredAlias == "" {
			return configuration.GetDefaultCredentials()
		}
		return configuration.GetCredentials(params.CredAlias)
	}

	if baseclientutils.IsLocalAddress(params.Address) {
		// local calls do not require any authentication
		return nil, nil
	}
	// credential name is required for non-local calls to resolve
	// the corresponding API key.
	return nil, fmt.Errorf("credential name is required")
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
