// Copyright 2023 Intrinsic Innovation LLC

package pubsubtesting

import (
	"context"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"

	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
)

type FakeInstalledAssetsServer struct {
	iagrpcpb.UnimplementedInstalledAssetsServer
	GetInstalledAssetFn    func(ctx context.Context, req *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error)
	CreateInstalledAssetFn func(ctx context.Context, req *iagrpcpb.CreateInstalledAssetRequest) (*lropb.Operation, error)
	DeleteInstalledAssetFn func(ctx context.Context, req *iagrpcpb.DeleteInstalledAssetRequest) (*lropb.Operation, error)
}

func NewFakeInstalledAssetsServer() *FakeInstalledAssetsServer {
	return &FakeInstalledAssetsServer{}
}

func (s *FakeInstalledAssetsServer) GetInstalledAsset(ctx context.Context, req *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
	if s.GetInstalledAssetFn != nil {
		return s.GetInstalledAssetFn(ctx, req)
	}
	return nil, nil
}

func (s *FakeInstalledAssetsServer) CreateInstalledAsset(ctx context.Context, req *iagrpcpb.CreateInstalledAssetRequest) (*lropb.Operation, error) {
	if s.CreateInstalledAssetFn != nil {
		return s.CreateInstalledAssetFn(ctx, req)
	}
	return nil, nil
}

func (s *FakeInstalledAssetsServer) DeleteInstalledAsset(ctx context.Context, req *iagrpcpb.DeleteInstalledAssetRequest) (*lropb.Operation, error) {
	if s.DeleteInstalledAssetFn != nil {
		return s.DeleteInstalledAssetFn(ctx, req)
	}
	return nil, nil
}
