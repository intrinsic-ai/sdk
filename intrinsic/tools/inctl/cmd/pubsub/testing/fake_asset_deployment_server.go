// Copyright 2023 Intrinsic Innovation LLC

package pubsubtesting

import (
	"context"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"

	adgrpcpb "intrinsic/assets/proto/asset_deployment_go_proto"
)

type FakeAssetDeploymentServer struct {
	adgrpcpb.UnimplementedAssetDeploymentServiceServer
	CreateResourceFromCatalogFn func(ctx context.Context, req *adgrpcpb.CreateResourceFromCatalogRequest) (*lropb.Operation, error)
	DeleteResourceFn            func(ctx context.Context, req *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error)
}

func NewFakeAssetDeploymentServer() *FakeAssetDeploymentServer {
	return &FakeAssetDeploymentServer{}
}

func (s *FakeAssetDeploymentServer) CreateResourceFromCatalog(ctx context.Context, req *adgrpcpb.CreateResourceFromCatalogRequest) (*lropb.Operation, error) {
	if s.CreateResourceFromCatalogFn != nil {
		return s.CreateResourceFromCatalogFn(ctx, req)
	}
	return nil, nil
}

func (s *FakeAssetDeploymentServer) DeleteResource(ctx context.Context, req *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
	if s.DeleteResourceFn != nil {
		return s.DeleteResourceFn(ctx, req)
	}
	return nil, nil
}
