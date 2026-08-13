// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
