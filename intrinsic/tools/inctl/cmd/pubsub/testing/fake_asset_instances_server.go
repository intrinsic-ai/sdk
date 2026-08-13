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
