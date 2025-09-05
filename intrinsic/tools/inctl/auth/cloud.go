// Copyright 2023 Intrinsic Innovation LLC

package auth

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"intrinsic/config/environments"
	"intrinsic/kubernetes/acl/identity"
)

// ConnectionOpts contains the options for creating a new gRPC connection to a cloud service.
type ConnectionOpts struct {
	project       string
	targetProject string
	org           string
	opts          []grpc.DialOption
	apiKey        string
	cluster       string

	// callbacks
	onIdentity func(u *identity.User)
}

// WithProject sets the cloud-project to use for the connection.
// If set together with WithOrg, the organization's API key will be used but the address will be
// resolved using the project provided. You can use this if you want to target a different project
// than the one associated with the organization. This can be necessary for global services (e.g.
// accounts, assets, portal).
func WithProject(p string) ConnectionOptsFunc {
	return func(c *ConnectionOpts) {
		c.project = p
	}
}

// WithTargetProject sets the cloud-project to use for the connection.
// If set together with WithOrg / WithProject, the org and/or project's API key will be used but the
// address will be resolved using the project provided. You can use this if you want to target a
// different project than the one associated with the API key. This can be necessary for global
// services (e.g. accounts, assets, portal).
func WithTargetProject(p string) ConnectionOptsFunc {
	return func(c *ConnectionOpts) {
		c.targetProject = p
	}
}

// WithOrg sets the organization to use for the connection.
func WithOrg(o string) ConnectionOptsFunc {
	return func(c *ConnectionOpts) {
		c.org = o
	}
}

// ConnectionOptsFunc is a function that can be used to configure the connection.
type ConnectionOptsFunc func(*ConnectionOpts)

// WithDialOptions sets the dial options to use for the connection.
func WithDialOptions(opts ...grpc.DialOption) ConnectionOptsFunc {
	return func(c *ConnectionOpts) {
		c.opts = append(c.opts, opts...)
	}
}

// WithFlagValues sets the project and organization to use for the connection from the current inctl
// CLI flags such as --project and --org.
// Must not be used together with WithProject or WithOrg.
func WithFlagValues(v *viper.Viper) ConnectionOptsFunc {
	return func(c *ConnectionOpts) {
		c.project = v.GetString(KeyProject)
		c.org = v.GetString(KeyOrganization)
	}
}

// WithAPIKey sets the API key to use for the connection.
// Skips loading the API key from the configuration store.
func WithAPIKey(k string) ConnectionOptsFunc {
	return func(c *ConnectionOpts) {
		c.apiKey = k
	}
}

// WithOnIdentityCallback sets a callback to be called with the identity of the authenticated user.
// Guaranteed to be called exactly once before the connection is returned.
func WithOnIdentityCallback(f func(u *identity.User)) ConnectionOptsFunc {
	return func(c *ConnectionOpts) {
		c.onIdentity = f
	}
}

// WithCluster sets the target cluster to connect to via the cloud relay.
func WithCluster(cluster string) ConnectionOptsFunc {
	return func(c *ConnectionOpts) {
		c.cluster = cluster
	}
}

// ErrorDetails contains the error details for a failed connection.
// This is used to provide more details about the error to the user in PrintErrorDetails.
type ErrorDetails struct {
	Opts    *ConnectionOpts
	Env     string
	Message string
	Help    string
}

func (e *ErrorDetails) Error() string {
	var values []string
	for _, kv := range [][]string{
		{"project", e.Opts.project},
		{"targetProject", e.Opts.targetProject},
		{"org", e.Opts.org},
		{"env", e.Env},
		{"cluster", e.Opts.cluster},
	} {
		if kv[1] != "" {
			values = append(values, fmt.Sprintf("%s: %q", kv[0], kv[1]))
		}
	}
	msg := fmt.Sprintf("%s (%s)", e.Message, strings.Join(values, ", "))
	if e.Help != "" {
		msg = fmt.Sprintf("%s\n%s", msg, e.Help)
	}
	return msg
}

var (
	// ErrUnableToRetrieveToken is returned if the token cannot be retrieved.
	ErrUnableToRetrieveToken = errors.New("unable to retrieve token")
)

// NewCloudConnection creates a new gRPC connection to a cloud project.
//
// This should be used for all connections to cloud services from inctl. This ensures that the
// connection uses the correct authentication and adds the necessary metadata to the requests. It
// makes sure that the API key is valid before the connection is established.
//
// Use either with:
//   - WithFlagValues (preferred): The organization and project to use for the connection will be read from
//     the current configuration in the inctl config or CLI flags. The organization ID will be added to
//     the request metadata if --org was specified.
//   - WithOrg: The organization to use for the connection. This will read the API key
//     from the organization and use it for the connection. Additionally, the organization ID will be
//     added to the request metadata.
//   - WithProject: The project to use for the connection. This will read the API key from the configuration
//     store for the given project and use it for the connection. No organization ID will be added to the request metadata.
//   - WithOrg & WithProject: Both can be set, in this case the API key of the organization will be used but
//     the address will be resolved using the project provided. You can use this if you want to target
//     a different project than the one associated with the organization. This can be  necessary for
//     global services (e.g. accounts, assets, portal).
func NewCloudConnection(ctx context.Context, optsFuncs ...ConnectionOptsFunc) (*grpc.ClientConn, error) {
	opts, tkSource, addMd, err := newOrLoadTokenSource(ctx, optsFuncs...)
	if err != nil {
		return nil, err
	}
	return newConnection(ctx, opts, tkSource, addMd)
}

// NewCloudClient creates a new http.Client that is authenticated for the cloud project.
//
// This should be used for all HTTP connections to cloud services from inctl. This ensures that the
// connection uses the correct authentication and adds the necessary metadata to the requests. It
// makes sure that the API key is valid before the connection is established.
//
// See NewCloudConnection for more details on how to configure the connection.
func NewCloudClient(ctx context.Context, optsFuncs ...ConnectionOptsFunc) (*http.Client, error) {
	_, tkSource, addMd, err := newOrLoadTokenSource(ctx, optsFuncs...)
	if err != nil {
		return nil, err
	}
	hc := &http.Client{
		Transport: &authenticatedTransport{
			base: http.DefaultTransport,
			ts:   tkSource,
			md:   addMd,
		},
	}
	return hc, nil
}

type authenticatedTransport struct {
	base http.RoundTripper
	ts   *cachedTokenSource
	md   *AddMetadata
}

func (t *authenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	tk, err := t.ts.Token(req.Context())
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "auth-proxy", Value: tk})

	for k, v := range t.md.metadata {
		if k == "" || v == "" {
			continue
		}
		req.Header.Set(k, v)
	}
	for k, v := range t.md.cookies {
		if k == "" || v == "" {
			continue
		}
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	return t.base.RoundTrip(req)
}

func newOrLoadTokenSource(ctx context.Context, optsFuncs ...ConnectionOptsFunc) (*ConnectionOpts, *cachedTokenSource, *AddMetadata, error) {
	opts := ConnectionOpts{}
	for _, f := range optsFuncs {
		f(&opts)
	}
	if opts.project == "" && opts.org == "" {
		return nil, nil, nil, fmt.Errorf("either project or org must be set")
	}

	errDetails := &ErrorDetails{
		Opts: &opts,
	}

	ak, err := loadAPIKey(&opts)
	if err != nil {
		return nil, nil, nil, err
	}

	// Determine the environment from the project.
	// The environment is important because the ID tokens issues by the
	// accounts service are tied to a specific environment.
	env, err := environments.FromProject(opts.project)
	if err != nil {
		env = environments.FromComputeProject(opts.project)
		// default to prod if we cannot determine the environment
	}
	errDetails.Env = env

	tkSource, err := newTokenSource(env, ak)
	if err != nil {
		errDetails.Message = "unable to create API key token source"
		return nil, nil, nil, errors.Join(err, errDetails)
	}

	tk, err := tkSource.Token(ctx)
	if err != nil {
		errDetails.Message = "unable to retrieve token"
		errDetails.Help = "This often indicates that your API key is expired or got invalidated. Please run `inctl auth login` and follow the instructions."
		return nil, nil, nil, errors.Join(err, errDetails)
	}
	// if requested, return the identity of the authenticated user
	if opts.onIdentity != nil {
		u, err := identity.UserFromJWT(tk)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to get identity from context: %w", err)
		}
		opts.onIdentity(u)
	}
	md := &AddMetadata{
		cookies:  map[string]string{},
		metadata: map[string]string{},
	}
	if opts.org != "" {
		md.cookies["org-id"] = opts.org
	}
	if opts.cluster != "" {
		md.metadata["x-server-name"] = opts.cluster
	}
	return &opts, tkSource, md, nil
}

func newConnection(ctx context.Context, opts *ConnectionOpts, tkSource *cachedTokenSource, md *AddMetadata) (*grpc.ClientConn, error) {
	grpcOpts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(&perRPCCreds{ts: tkSource, md: md}),
		grpc.WithStatsHandler(new(ocgrpc.ClientHandler)),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	}

	project := opts.project
	if opts.targetProject != "" {
		project = opts.targetProject
	}
	addr := fmt.Sprintf("dns:///%s:443", environments.Domain(project))

	return grpc.NewClient(addr, grpcOpts...)
}

// loadAPIKey loads the API key to use for the connection.
// If WithAPIKey was set, it will be used. Otherwise, the API key will be loaded from the
// configuration store for the given project or organization.
func loadAPIKey(opts *ConnectionOpts) (string, error) {
	if opts.apiKey != "" {
		return opts.apiKey, nil
	}
	if opts.org != "" {
		orgInfo, err := NewStore().ReadOrgInfo(opts.org)
		if err == nil { // if no error
			// take the project from the org if WithProject was not set.
			if opts.project == "" {
				opts.project = orgInfo.Project
			}
		}
	}
	cfg, err := NewStore().GetConfiguration(opts.project)
	if err != nil {
		return "", err
	}
	creds, err := cfg.GetDefaultCredentials()
	if err != nil {
		return "", err
	}
	return creds.APIKey, nil
}

// newTokenSource creates a new API key token source for the given API key and environment.
// The token source will add the given metadata to the request.
func newTokenSource(env, key string) (*cachedTokenSource, error) {
	// This is portal and not accounts because we are using the grpc-http gateway for token requests.
	fsAddr := environments.PortalDomain(env)
	if fsAddr == "" { // default to prod if we cannot determine the environment
		fsAddr = environments.PortalDomain(environments.Prod)
	}
	factory := getSharedTokenSourceFactory()
	tsc, err := factory.LoadOrNew(fsAddr, key)
	if err != nil {
		return nil, fmt.Errorf("cannot create token exchange: %w", err)
	}
	return tsc, nil
}
