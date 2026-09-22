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

	"google.golang.org/protobuf/types/known/emptypb"

	lineconfigstoragepb "intrinsic/platform/pubsub/connect/cloud/proto/line_configuration_storage/v1/line_configuration_storage_go_proto"
	lineconfigpb "intrinsic/platform/pubsub/connect/common/proto/line_configuration/v1/line_configuration_go_proto"
)

// FakeLineConfigurationStorageServer is a fake implementation of LineConfigurationStorageServer.
//
// Unit tests can provide their own implementation of any of its functions.
type FakeLineConfigurationStorageServer struct {
	lineconfigstoragepb.UnimplementedLineConfigurationStorageServiceServer

	SetFn    func(context.Context, *lineconfigstoragepb.UpdateLineConfigurationRequest) (*lineconfigpb.LineConfiguration, error)
	GetFn    func(context.Context, *lineconfigstoragepb.GetLineConfigurationRequest) (*lineconfigpb.LineConfiguration, error)
	DeleteFn func(context.Context, *lineconfigstoragepb.DeleteLineConfigurationRequest) (*emptypb.Empty, error)
}

// NewFakeLineConfigurationStorageServer creates a new server.
func NewFakeLineConfigurationStorageServer() *FakeLineConfigurationStorageServer {
	return &FakeLineConfigurationStorageServer{}
}

// Get is the fake implementation of the API that reads factory line's metadata from Firestore.
//
// The default implementation returns an empty response.
func (s *FakeLineConfigurationStorageServer) GetLineConfiguration(ctx context.Context, req *lineconfigstoragepb.GetLineConfigurationRequest) (*lineconfigpb.LineConfiguration, error) {
	if s.GetFn != nil {
		return s.GetFn(ctx, req)
	}
	return &lineconfigpb.LineConfiguration{}, nil
}

// Set is the fake implementation of the API that stores factory line's metadata in Firestore.
//
// The default implementation returns an empty response that indicates successful write.
func (s *FakeLineConfigurationStorageServer) UpdateLineConfiguration(ctx context.Context, req *lineconfigstoragepb.UpdateLineConfigurationRequest) (*lineconfigpb.LineConfiguration, error) {
	if s.SetFn != nil {
		return s.SetFn(ctx, req)
	}
	return &lineconfigpb.LineConfiguration{}, nil
}

// Delete is the fake implementation of the API that deletes factory line's metadata from Firestore.
//
// The default implementation returns an empty response that indicates successful deletion.
func (s *FakeLineConfigurationStorageServer) DeleteLineConfiguration(ctx context.Context, req *lineconfigstoragepb.DeleteLineConfigurationRequest) (*emptypb.Empty, error) {
	if s.DeleteFn != nil {
		return s.DeleteFn(ctx, req)
	}
	return &emptypb.Empty{}, nil
}
