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
