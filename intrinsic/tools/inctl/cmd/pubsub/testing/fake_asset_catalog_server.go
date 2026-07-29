// Copyright 2023 Intrinsic Innovation LLC

package pubsubtesting

import (
	"context"

	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
)

type FakeAssetCatalogServer struct {
	acgrpcpb.UnimplementedAssetCatalogServer
	GetAssetFn func(ctx context.Context, req *acgrpcpb.GetAssetRequest) (*acgrpcpb.Asset, error)
}

func NewFakeAssetCatalogServer() *FakeAssetCatalogServer {
	return &FakeAssetCatalogServer{}
}

func (s *FakeAssetCatalogServer) GetAsset(ctx context.Context, req *acgrpcpb.GetAssetRequest) (*acgrpcpb.Asset, error) {
	if s.GetAssetFn != nil {
		return s.GetAssetFn(ctx, req)
	}
	return nil, nil
}
