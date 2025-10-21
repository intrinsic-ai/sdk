// Copyright 2023 Intrinsic Innovation LLC

// Package utils provides utility functions for Asset dependencies.
package utils

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	rdpb "intrinsic/assets/proto/v1/resolved_dependency_go_proto"
)

// Connect creates a gRPC connection for communicating with the provider of the specified interface.
//
// It also returns a new context that includes any needed metadata for communicating with the
// provider.
func Connect(ctx context.Context, dep *rdpb.ResolvedDependency, iface string) (*grpc.ClientConn, context.Context, error) {
	// Create a connection to the provider.
	ifaceProto, ok := dep.GetInterfaces()[iface]
	if !ok {
		return nil, nil, fmt.Errorf("interface %q not found in resolved dependency", iface)
	}

	connection := ifaceProto.GetGrpcConnection()
	if connection == nil {
		return nil, nil, fmt.Errorf("interface %q is not gRPC or no connection information is available", iface)
	}

	conn, err := grpc.NewClient(connection.GetAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create gRPC client for interface %q: %w", iface, err)
	}

	// Add any needed metadata to the context.
	for _, m := range connection.GetMetadata() {
		ctx = metadata.AppendToOutgoingContext(ctx, m.GetKey(), m.GetValue())
	}

	return conn, ctx, nil
}
