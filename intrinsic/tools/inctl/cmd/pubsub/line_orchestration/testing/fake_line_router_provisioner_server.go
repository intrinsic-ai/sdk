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

	provisionerpb "intrinsic/platform/pubsub/cloud_router_provisioner/v1/provisioner_go_proto"
)

// FakeLineRouterProvisionerServer is a fake implementation of LineRouterProvisionerServer.
type FakeLineRouterProvisionerServer struct {
	provisionerpb.UnimplementedLineRouterProvisionerServer

	// ProvisionFn is called from the server's Provision method.
	ProvisionFn func(ctx context.Context, req *provisionerpb.ProvisionLineRouterRequest) (*provisionerpb.ProvisionLineRouterResponse, error)

	// DeleteFn is called from the server's Delete method.
	DeleteFn func(ctx context.Context, req *provisionerpb.DeleteLineRouterRequest) (*provisionerpb.LineRouterDeleteResponse, error)
}

// NewFakeLineRouterProvisionerServer creates a new FakeLineRouterProvisionerServer.
func NewFakeLineRouterProvisionerServer() *FakeLineRouterProvisionerServer {
	return &FakeLineRouterProvisionerServer{}
}

// Provision is the fake implementation of the API that provisions a line-level router.
func (s *FakeLineRouterProvisionerServer) Provision(ctx context.Context, req *provisionerpb.ProvisionLineRouterRequest) (*provisionerpb.ProvisionLineRouterResponse, error) {
	if s.ProvisionFn != nil {
		return s.ProvisionFn(ctx, req)
	}
	return nil, nil
}

// Delete is the fake implementation of the API that deletes a line-level router.
//
// The default implementation returns an empty response that indicates successful deletion.
func (s *FakeLineRouterProvisionerServer) Delete(ctx context.Context, req *provisionerpb.DeleteLineRouterRequest) (*provisionerpb.LineRouterDeleteResponse, error) {
	if s.DeleteFn != nil {
		return s.DeleteFn(ctx, req)
	}

	return &provisionerpb.LineRouterDeleteResponse{}, nil
}
