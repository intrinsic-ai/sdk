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
	"bytes"
	"context"
	"net"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"

	adgrpcpb "intrinsic/assets/proto/asset_deployment_go_proto"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	aigrpcpb "intrinsic/assets/proto/v1/asset_instances_go_proto"
	provisionerpb "intrinsic/platform/pubsub/cloud_router_provisioner/v1/provisioner_go_proto"
	lineconfigstoragepb "intrinsic/platform/pubsub/connect/cloud/proto/line_configuration_storage/v1/line_configuration_storage_go_proto"
)

const (
	// Used in tests with multiple instances of the relay service.
	AnotherHubServiceName = "another-hub-service"

	LocalHubEndpoint    = "node-hub@local"
	RemoteHubEndpoint   = "vmp-hub@remote"
	RemoteSpokeEndpoint = "vmp-spoke@remote"
)

func VerifyExpectedOutputAndError(t *testing.T, receivedOutput *bytes.Buffer, receivedError error, expectError bool, expectErrorContains string, expectedOutput []string) {
	t.Helper()
	for _, expectedOutputFragment := range expectedOutput {
		if !strings.Contains(receivedOutput.String(), expectedOutputFragment) {
			t.Errorf("expected output to contain %q, got %q", expectedOutputFragment, receivedOutput.String())
		}
	}

	if expectError {
		if receivedError == nil {
			t.Fatalf("expected error containing %q, got nil", expectErrorContains)
		}
		if !strings.Contains(receivedError.Error(), expectErrorContains) {
			t.Errorf("expected error to contain %q, got %v", expectErrorContains, receivedError)
		}
	} else if receivedError != nil {
		t.Fatalf("expected no error, got %v", receivedError)
	}
}

func EnsureNoUnexpectedOutput(t *testing.T, receivedOutput *bytes.Buffer, unexpectedOutput []string) {
	t.Helper()
	for _, outputFragment := range unexpectedOutput {
		if strings.Contains(receivedOutput.String(), outputFragment) {
			t.Errorf("output should not contain %q, got %q", outputFragment, receivedOutput.String())
		}
	}
}

type TestServerResources struct {
	grpcServer              *grpc.Server
	Listener                *bufconn.Listener
	InstServer              *FakeAssetInstancesServer
	DepServer               *FakeAssetDeploymentServer
	OpServer                *FakeOperationsServer
	IaServer                *FakeInstalledAssetsServer
	ProvisionerServer       *FakeLineRouterProvisionerServer
	LineConfigStorageServer *FakeLineConfigurationStorageServer
}

func SetupTestServer(t *testing.T) *TestServerResources {
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()

	instServer := NewFakeAssetInstancesServer()
	depServer := NewFakeAssetDeploymentServer()
	opServer := NewFakeOperationsServer()
	iaServer := NewFakeInstalledAssetsServer()
	provisionerServer := NewFakeLineRouterProvisionerServer()
	lineConfigStorageServer := NewFakeLineConfigurationStorageServer()

	aigrpcpb.RegisterAssetInstancesServer(s, instServer)
	adgrpcpb.RegisterAssetDeploymentServiceServer(s, depServer)
	lropb.RegisterOperationsServer(s, opServer)
	iagrpcpb.RegisterInstalledAssetsServer(s, iaServer)
	provisionerpb.RegisterLineRouterProvisionerServer(s, provisionerServer)
	lineconfigstoragepb.RegisterLineConfigurationStorageServiceServer(s, lineConfigStorageServer)

	go func() {
		if err := s.Serve(lis); err != nil {
			// Ignore error on close
		}
	}()

	t.Cleanup(func() {
		s.Stop()
		lis.Close()
	})

	return &TestServerResources{
		grpcServer:              s,
		Listener:                lis,
		InstServer:              instServer,
		DepServer:               depServer,
		OpServer:                opServer,
		IaServer:                iaServer,
		ProvisionerServer:       provisionerServer,
		LineConfigStorageServer: lineConfigStorageServer,
	}
}

func DialTestServer(ctx context.Context, lis *bufconn.Listener) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithInsecure())
}
