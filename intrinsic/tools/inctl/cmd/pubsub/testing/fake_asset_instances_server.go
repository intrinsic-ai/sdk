// Copyright 2023 Intrinsic Innovation LLC

package pubsubtesting

import (
	"context"

	aigrpcpb "intrinsic/assets/proto/v1/asset_instances_go_proto"
)

type FakeAssetInstancesServer struct {
	aigrpcpb.UnimplementedAssetInstancesServer
	GetAssetInstanceFn   func(ctx context.Context, req *aigrpcpb.GetAssetInstanceRequest) (*aigrpcpb.AssetInstance, error)
	ListAssetInstancesFn func(ctx context.Context, req *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error)
}

func NewFakeAssetInstancesServer() *FakeAssetInstancesServer {
	return &FakeAssetInstancesServer{}
}

func (s *FakeAssetInstancesServer) GetAssetInstance(ctx context.Context, req *aigrpcpb.GetAssetInstanceRequest) (*aigrpcpb.AssetInstance, error) {
	if s.GetAssetInstanceFn != nil {
		return s.GetAssetInstanceFn(ctx, req)
	}
	return nil, nil
}

func (s *FakeAssetInstancesServer) ListAssetInstances(ctx context.Context, req *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
	if s.ListAssetInstancesFn != nil {
		return s.ListAssetInstancesFn(ctx, req)
	}
	return nil, nil
}
