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

package any

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"

	prpb "intrinsic/proto_tools/proto/proto_registry_go_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"

	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"

	"intrinsic/assets/proto/id_go_proto"
	"intrinsic/util/proto/testing/prototestutil"
)

const anyTypeUrl = "type.intrinsic.ai/google.protobuf.Any"

type FakeServer struct {
	iagrpcpb.UnimplementedInstalledAssetsServer
	prpb.UnimplementedProtoRegistryServer
	mu            sync.Mutex
	version       string
	batchGetCount int
	fds           *descriptorpb.FileDescriptorSet
}

func (s *FakeServer) getVersionLocked() string {
	if s.version != "" {
		return s.version
	}
	return "1.0.0"
}

func (s *FakeServer) SetVersion(version string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.version = version
}

func (s *FakeServer) BatchGetCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.batchGetCount
}

func (s *FakeServer) ListInstalledAssets(ctx context.Context, req *iagrpcpb.ListInstalledAssetsRequest) (*iagrpcpb.ListInstalledAssetsResponse, error) {
	s.mu.Lock()
	version := s.getVersionLocked()
	s.mu.Unlock()

	return &iagrpcpb.ListInstalledAssetsResponse{
		InstalledAssets: []*iagrpcpb.InstalledAsset{
			{
				Metadata: &metadatapb.Metadata{
					IdVersion: &id_go_proto.IdVersion{
						Id: &id_go_proto.Id{
							Package: "com.example",
							Name:    "Foobar",
						},
						Version: version,
					},
				},
			},
		},
	}, nil
}

func (s *FakeServer) BatchGetInstalledAssets(ctx context.Context, req *iagrpcpb.BatchGetInstalledAssetsRequest) (*iagrpcpb.BatchGetInstalledAssetsResponse, error) {
	s.mu.Lock()
	s.batchGetCount++
	version := s.getVersionLocked()
	s.mu.Unlock()

	fds := prototestutil.FileDescriptorSet(&anypb.Any{})
	return &iagrpcpb.BatchGetInstalledAssetsResponse{
		InstalledAssets: []*iagrpcpb.InstalledAsset{
			{
				Metadata: &metadatapb.Metadata{
					IdVersion: &id_go_proto.IdVersion{
						Id: &id_go_proto.Id{
							Package: "com.example",
							Name:    "Foobar",
						},
						Version: version,
					},
					FileDescriptorSet: fds,
				},
			},
		},
	}, nil
}

func MustMakeFakeServerWithServer(t *testing.T) (*FakeServer, string) {
	t.Helper()
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	fakeServer := &FakeServer{}
	iagrpcpb.RegisterInstalledAssetsServer(s, fakeServer)
	fakeServer.fds = prototestutil.FileDescriptorSet(&anypb.Any{})
	prpb.RegisterProtoRegistryServer(s, fakeServer)

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server exited with error: %v", err)
		}
	}()

	t.Cleanup(func() {
		s.Stop()
	})

	return fakeServer, lis.Addr().String()
}

func MustMakeFakeServer(t *testing.T) string {
	_, addr := MustMakeFakeServerWithServer(t)
	return addr
}

func (s *FakeServer) GetNamedFileDescriptorSet(ctx context.Context, req *prpb.GetNamedFileDescriptorSetRequest) (*prpb.NamedFileDescriptorSet, error) {
	fmt.Printf("Got request for name %s\n", req.GetName())
	if req.GetTypeUrl() == anyTypeUrl {
		fmt.Printf("returning valid FDS for %s", req.GetName())
		return &prpb.NamedFileDescriptorSet{
			Name:              req.GetName(),
			FileDescriptorSet: s.fds,
		}, nil
	}
	return nil, status.Errorf(codes.NotFound, "Type URL %s not found", req.GetName())
}
