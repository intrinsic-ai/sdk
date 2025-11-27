// Copyright 2023 Intrinsic Innovation LLC

package accounts

import (
	"context"

	"intrinsic/config/environments"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/orgutil"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	accaccesscontrolgrpcpb "intrinsic/kubernetes/accounts/service/api/accesscontrol/v1/accesscontrol_go_grpc_proto"
	accresourcemanagergrpcpb "intrinsic/kubernetes/accounts/service/api/resourcemanager/v1/resourcemanager_go_grpc_proto"
)

// NewResourceManagerV1Client creates a new secure ResourceManagerServiceClient using API keys.
// Aliases for convenience
var NewResourceManagerV1Client = func(ctx context.Context, vipr *viper.Viper) (ResourceManagerV1Client, error) {
	conn, err := newConn(ctx, vipr)
	if err != nil {
		return nil, err
	}
	return accresourcemanagergrpcpb.NewResourceManagerServiceClient(conn), nil
}

// NewAccessControlV1Client creates a new secure AccessControlServiceClient using API keys.
// Aliases for convenience
var NewAccessControlV1Client = func(ctx context.Context, vipr *viper.Viper) (AccessControlV1Client, error) {
	conn, err := newConn(ctx, vipr)
	if err != nil {
		return nil, err
	}
	return accaccesscontrolgrpcpb.NewAccessControlServiceClient(conn), nil
}

// newConn creates a new secure connection to the accounts service using auth.NewCloudConnection.
// The accounts service is a global service and the project used for the connection is determined
// from the --project flag. If the flag is not set, prod is used as default.
// The connection is authenticated using the API key of the organization or project.
func newConn(ctx context.Context, vipr *viper.Viper) (*grpc.ClientConn, error) {
	project := vipr.GetString(orgutil.KeyProject)
	accountsProject := environments.AccountsProjectFromProject(project)
	if accountsProject == "" { // always fallback to prod
		accountsProject = environments.AccountsProjectProd
	}
	return auth.NewCloudConnection(ctx, auth.WithFlagValues(vipr), auth.WithTargetProject(accountsProject))
}

// ResourceManagerV1Client is the client for the ResourceManagerService.
// Aliases for convenience
type ResourceManagerV1Client = accresourcemanagergrpcpb.ResourceManagerServiceClient

// AccessControlV1Client is the client for the AccessControlService.
// Aliases for convenience
type AccessControlV1Client = accaccesscontrolgrpcpb.AccessControlServiceClient
