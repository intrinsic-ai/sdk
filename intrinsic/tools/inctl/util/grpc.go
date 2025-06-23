// Copyright 2023 Intrinsic Innovation LLC

// Package grpc contains grpc client helpers.
package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"intrinsic/assets/baseclientutils"
	"intrinsic/kubernetes/acl/identity"
	"intrinsic/tools/inctl/auth/auth"
)

// NewIPCGRPCClient creates a new gRPC client for contacting an IPC through the API relay.
func NewIPCGRPCClient(ctx context.Context, projectName, orgName, clusterName string) (context.Context, *grpc.ClientConn, error) {
	address := fmt.Sprintf("dns:///www.endpoints.%s.cloud.goog:443", projectName)
	configuration, err := auth.NewStore().GetConfiguration(projectName)
	if err != nil {
		return ctx, nil, fmt.Errorf("credentials not found: %w", err)
	}
	creds, err := configuration.GetDefaultCredentials()
	if err != nil {
		return ctx, nil, fmt.Errorf("get default credentials: %w", err)
	}
	tcOption, err := baseclientutils.GetTransportCredentialsDialOption()
	if err != nil {
		return ctx, nil, fmt.Errorf("cannot retrieve transport credentials: %w", err)
	}
	dialerOpts := append(baseclientutils.BaseDialOptions(),
		grpc.WithPerRPCCredentials(creds),
		tcOption,
	)
	conn, err := grpc.NewClient(address, dialerOpts...)
	if err != nil {
		return ctx, nil, fmt.Errorf("dialing context: %w", err)
	}
	ctx, err = identity.OrgToContext(ctx, orgName)
	if err != nil {
		return ctx, nil, fmt.Errorf("unable to setup the context: %w", err)
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "x-server-name", clusterName)
	return ctx, conn, nil
}
