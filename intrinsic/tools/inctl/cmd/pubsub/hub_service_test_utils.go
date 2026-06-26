// Copyright 2023 Intrinsic Innovation LLC

package pubsub

import (
	"context"
	"net"
	"testing"

	pubsubtesting "intrinsic/tools/inctl/cmd/pubsub/testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"

	adgrpcpb "intrinsic/assets/proto/asset_deployment_go_proto"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	aigrpcpb "intrinsic/assets/proto/v1/asset_instances_go_proto"
)

const (
	// Used in tests with multiple instances of the relay service.
	anotherHubServiceName = "another-hub-service"
)

type testServerResources struct {
	grpcServer *grpc.Server
	listener   *bufconn.Listener
	instServer *pubsubtesting.FakeAssetInstancesServer
	depServer  *pubsubtesting.FakeAssetDeploymentServer
	opServer   *pubsubtesting.FakeOperationsServer
	iaServer   *pubsubtesting.FakeInstalledAssetsServer
}

func setupTestServer(t *testing.T) *testServerResources {
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()

	instServer := pubsubtesting.NewFakeAssetInstancesServer()
	depServer := pubsubtesting.NewFakeAssetDeploymentServer()
	opServer := pubsubtesting.NewFakeOperationsServer()
	iaServer := pubsubtesting.NewFakeInstalledAssetsServer()

	aigrpcpb.RegisterAssetInstancesServer(s, instServer)
	adgrpcpb.RegisterAssetDeploymentServiceServer(s, depServer)
	lropb.RegisterOperationsServer(s, opServer)
	iagrpcpb.RegisterInstalledAssetsServer(s, iaServer)

	go func() {
		if err := s.Serve(lis); err != nil {
			// Ignore error on close
		}
	}()

	t.Cleanup(func() {
		s.Stop()
		lis.Close()
	})

	return &testServerResources{
		grpcServer: s,
		listener:   lis,
		instServer: instServer,
		depServer:  depServer,
		opServer:   opServer,
		iaServer:   iaServer,
	}
}

func dialTestServer(ctx context.Context, lis *bufconn.Listener) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithInsecure())
}
