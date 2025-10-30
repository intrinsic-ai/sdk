// Copyright 2023 Intrinsic Innovation LLC

// Package utils provides utility functions for Asset dependencies.
package utils

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	anypb "google.golang.org/protobuf/types/known/anypb"
	dasgrpcpb "intrinsic/assets/data/proto/v1/data_assets_go_grpc_proto"
	daspb "intrinsic/assets/data/proto/v1/data_assets_go_grpc_proto"
	rdpb "intrinsic/assets/proto/v1/resolved_dependency_go_proto"
)

const ingressAddress = "istio-ingressgateway.app-ingress.svc.cluster.local:80"

var (
	errMissingInterface = errors.New("interface not found in resolved dependency")
	errNotGRPC          = errors.New("interface is not gRPC or no connection information is available")
	errNotData          = errors.New("interface is not data or no data dependency information is available")
)

// Connect creates a gRPC connection for communicating with the provider of the specified interface.
//
// It also returns a new context that includes any needed metadata for communicating with the
// provider.
func Connect(ctx context.Context, dep *rdpb.ResolvedDependency, iface string) (*grpc.ClientConn, context.Context, error) {
	ifaceProto, err := findInterface(dep, iface)
	if err != nil {
		return nil, nil, err
	}
	connection := ifaceProto.GetGrpcConnection()
	if connection == nil {
		return nil, nil, fmt.Errorf("%w: %q", errNotGRPC, iface)
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

type getDataOptions struct {
	dataAssetsClient dasgrpcpb.DataAssetsClient
}

// GetDataOption is an option for configuring the GetData function.
type GetDataOption func(*getDataOptions)

// WithDataAssetsClient sets the DataAssets client to use.
func WithDataAssetsClient(client dasgrpcpb.DataAssetsClient) GetDataOption {
	return func(opts *getDataOptions) {
		opts.dataAssetsClient = client
	}
}

// GetDataPayload returns the DataAsset payload for the specified interface.
//
// If no DataAssets client is provided, an insecure connection to the DataAssets service via the
// ingress gateway will be created. This connection is valid for services running in the same
// cluster as the DataAssets service.
func GetDataPayload(ctx context.Context, dep *rdpb.ResolvedDependency, iface string, options ...GetDataOption) (*anypb.Any, error) {
	opts := &getDataOptions{}
	for _, opt := range options {
		opt(opts)
	}

	ifaceProto, err := findInterface(dep, iface)
	if err != nil {
		return nil, err
	}
	dataDependency := ifaceProto.GetData()
	if dataDependency == nil {
		return nil, fmt.Errorf("%w: %q", errNotData, iface)
	}

	if opts.dataAssetsClient == nil {
		client, conn, err := makeDefaultDataAssetsClient()
		if err != nil {
			return nil, fmt.Errorf("failed to create default DataAssets client: %w", err)
		}
		defer conn.Close()
		opts.dataAssetsClient = client
	}

	// Get the DataAsset proto from the DataAssets service.
	da, err := opts.dataAssetsClient.GetDataAsset(ctx, &daspb.GetDataAssetRequest{
		Id: dataDependency.GetId(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get DataAsset proto for %q: %w", dataDependency.GetId(), err)
	}

	return da.GetData(), nil
}

func findInterface(dep *rdpb.ResolvedDependency, iface string) (*rdpb.ResolvedDependency_Interface, error) {
	ifaceProto, ok := dep.GetInterfaces()[iface]
	if !ok {
		var explanation string
		if len(dep.GetInterfaces()) == 0 {
			explanation = "no interfaces provided"
		} else {
			keys := slices.Collect(maps.Keys(dep.GetInterfaces()))
			explanation = fmt.Sprintf("got interfaces: %v", strings.Join(keys, ", "))
		}
		return nil, fmt.Errorf("%w: (want %q, %s)", errMissingInterface, iface, explanation)
	}
	return ifaceProto, nil
}

func makeDefaultDataAssetsClient() (dasgrpcpb.DataAssetsClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(ingressAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create gRPC client for DataAssets service: %w", err)
	}
	return dasgrpcpb.NewDataAssetsClient(conn), conn, nil
}
