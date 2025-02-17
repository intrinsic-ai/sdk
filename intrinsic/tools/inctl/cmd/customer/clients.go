// Copyright 2023 Intrinsic Innovation LLC

package customer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	grpccredentials "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"intrinsic/config/environments"
	accaccesscontrolgrpcpb "intrinsic/kubernetes/accounts/service/api/accesscontrol/v1/accesscontrol_go_grpc_proto"
	accresourcemanagergrpcpb "intrinsic/kubernetes/accounts/service/api/resourcemanager/v1/resourcemanager_go_grpc_proto"
	"intrinsic/kubernetes/acl/cookies"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/orgutil"
)

func authFromVipr() (string, string) {
	authOrg := vipr.GetString(orgutil.KeyOrganization)
	authProject := vipr.GetString(orgutil.KeyProject)
	authEnv := vipr.GetString(orgutil.KeyEnvironment)
	org := authOrg
	if authProject != "" {
		org = authOrg + "@" + authProject
	}
	return authEnv, org
}

// Aliases for convenience
var newresourcemanagerClient = func(ctx context.Context) (resourcemanagerClient, error) {
	env, org := authFromVipr()
	return newSecureAccountsResourceManagerAPIKeyClient(ctx, env, org)
}

var newAccessControlV1Client = func(ctx context.Context) (accessControlV1Client, error) {
	env, org := authFromVipr()
	return newSecureAccountsAccessControlAPIKeyClient(ctx, env, org)
}

// Aliases for convenience
type resourcemanagerClient = accresourcemanagergrpcpb.ResourceManagerServiceClient
type accessControlV1Client = accaccesscontrolgrpcpb.AccessControlServiceClient

func newSecureAccountsAccessControlAPIKeyClient(ctx context.Context, env, org string) (accaccesscontrolgrpcpb.AccessControlServiceClient, error) {
	conn, err := newConnAuthStore(ctx, environments.AccountsDomain(env), org)
	if err != nil {
		return nil, err
	}
	return accaccesscontrolgrpcpb.NewAccessControlServiceClient(conn), nil
}

// newSecureAccountsTokensServiceAPIKeyClient creates a new secure ResourceManagerClient using API keys.
// Suitable for calling the ressourcemanager via HTTPS from any environment.
func newSecureAccountsResourceManagerAPIKeyClient(ctx context.Context, env, org string) (accresourcemanagergrpcpb.ResourceManagerServiceClient, error) {
	conn, err := newConnAuthStore(ctx, environments.AccountsDomain(env), org)
	if err != nil {
		return nil, err
	}
	return accresourcemanagergrpcpb.NewResourceManagerServiceClient(conn), nil
}

// Can be overwridden/injected in tests.
var authStore = auth.NewStore()

func newConnAuthStore(ctx context.Context, addr, org string) (*grpc.ClientConn, error) {
	// determine address
	orgInfo, err := authStore.ReadOrgInfo(org)
	if err != nil {
		return nil, fmt.Errorf("failed to read org info for %q: %v", org, err)
	}
	// fetch API key
	project := orgInfo.Project
	cfg, err := auth.NewStore().GetConfiguration(project)
	if err != nil {
		return nil, fmt.Errorf("failed to get project configuration for project %q: %v", project, err)
	}
	creds, err := cfg.GetDefaultCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to get default credentials for project %q: %v", project, err)
	}
	return newConn(ctx, addr, grpc.WithPerRPCCredentials(creds))
}

func newConn(ctx context.Context, addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// create connection
	var grpcOpts = []grpc.DialOption{
		grpc.WithStatsHandler(new(ocgrpc.ClientHandler)),
		grpc.WithTransportCredentials(grpccredentials.NewTLS(&tls.Config{})),
	}
	grpcOpts = append(grpcOpts, opts...)
	conn, err := grpc.NewClient(addr+":443", grpcOpts...)
	if err != nil {
		return nil, errors.Wrapf(err, "grpc.Dial(%q)", addr)
	}
	return conn, nil
}

const (
	// orgIDCookie is the cookie key for the organization identifier.
	orgIDCookie = "org-id"
)

// withOrgID adds the org ID to the outgoing RCP context.
func withOrgID(ctx context.Context) context.Context {
	o := vipr.GetString(orgutil.KeyOrganization)
	md := cookies.ToMDString(&http.Cookie{Name: orgIDCookie, Value: o})
	return metadata.AppendToOutgoingContext(ctx, md...)
}
