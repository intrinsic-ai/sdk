// Copyright 2023 Intrinsic Innovation LLC

// Package fakeserver has a fake server for unit testing anyresolver
package fakeserver

import (
	"context"
	"fmt"
	"net"
	"testing"

	prpb "intrinsic/proto_tools/proto/proto_registry_go_proto"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	fds *descriptorpb.FileDescriptorSet
}

func (s *FakeServer) ListInstalledAssets(ctx context.Context, req *iagrpcpb.ListInstalledAssetsRequest) (*iagrpcpb.ListInstalledAssetsResponse, error) {
	return &iagrpcpb.ListInstalledAssetsResponse{
		InstalledAssets: []*iagrpcpb.InstalledAsset{
			{
				Metadata: &metadatapb.Metadata{
					IdVersion: &id_go_proto.IdVersion{
						Id: &id_go_proto.Id{
							Package: "com.example",
							Name:    "Foobar",
						},
						Version: "1.0.0",
					},
				},
			},
		},
	}, nil
}

func (s *FakeServer) BatchGetInstalledAssets(ctx context.Context, req *iagrpcpb.BatchGetInstalledAssetsRequest) (*iagrpcpb.BatchGetInstalledAssetsResponse, error) {
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
						Version: "1.0.0",
					},
					FileDescriptorSet: fds,
				},
			},
		},
	}, nil
}

func MustMakeFakeServer(t *testing.T) string {
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

	return lis.Addr().String()
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
