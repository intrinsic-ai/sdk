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
