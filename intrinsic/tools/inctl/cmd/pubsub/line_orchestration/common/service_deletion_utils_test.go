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

package servicedeletionutils

import (
	"bytes"
	"context"

	"errors"
	"fmt"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"

	adgrpcpb "intrinsic/assets/proto/asset_deployment_go_proto"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	aigrpcpb "intrinsic/assets/proto/v1/asset_instances_go_proto"
	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"
	pubsubtesting "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/testing"
)

func TestServiceDeletionCommand(t *testing.T) {
	tests := []struct {
		name                     string
		setupFakeInstServer      func(s *pubsubtesting.FakeAssetInstancesServer)
		setupFakeDepServer       func(s *pubsubtesting.FakeAssetDeploymentServer)
		setupFakeOpServer        func(s *pubsubtesting.FakeOperationsServer)
		setupFakeIAServer        func(s *pubsubtesting.FakeInstalledAssetsServer)
		expectedOutput           []string
		shouldRetainServiceAsset bool
		expectErr                bool
		expectErrContains        string
		packageName              string
		serviceName              string
	}{
		{
			name: "successful_uninstall",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{
						AssetInstances: []*aigrpcpb.AssetInstance{
							{Name: common.HubServiceName},
							{Name: pubsubtesting.AnotherHubServiceName},
						},
					}, nil
				}
			},
			setupFakeDepServer: func(s *pubsubtesting.FakeAssetDeploymentServer) {
				s.DeleteResourceFn = func(ctx context.Context, in *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			setupFakeOpServer: func(s *pubsubtesting.FakeOperationsServer) {
				s.GetOperationFn = func(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.DeleteInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.DeleteInstalledAssetRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			expectedOutput: []string{
				fmt.Sprintf("Deleting an instance of the %v service named %q", common.HubServiceName, common.HubServiceName),
				fmt.Sprintf("Deleting an instance of the %v service named %q", common.HubServiceName, pubsubtesting.AnotherHubServiceName),
				fmt.Sprintf("Successfully uninstalled the %v service asset", common.HubServiceName),
			},
			shouldRetainServiceAsset: false,
			expectErr:                false,
			packageName:              common.HubServicePackage,
			serviceName:              common.HubServiceName,
		},
		{
			name: "successful_uninstall_different_name",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{
						AssetInstances: []*aigrpcpb.AssetInstance{
							{Name: common.ForwardingServiceName},
						},
					}, nil
				}
			},
			setupFakeDepServer: func(s *pubsubtesting.FakeAssetDeploymentServer) {
				s.DeleteResourceFn = func(ctx context.Context, in *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			setupFakeOpServer: func(s *pubsubtesting.FakeOperationsServer) {
				s.GetOperationFn = func(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.DeleteInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.DeleteInstalledAssetRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			expectedOutput: []string{
				fmt.Sprintf("Deleting an instance of the %v service named %q", common.ForwardingServiceName, common.ForwardingServiceName),
				fmt.Sprintf("Successfully uninstalled the %v service asset", common.ForwardingServiceName),
			},
			shouldRetainServiceAsset: false,
			expectErr:                false,
			packageName:              common.ForwardingServicePackage,
			serviceName:              common.ForwardingServiceName,
		},
		{
			name: "remove_instances_dont_uninstall_asset",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{
						AssetInstances: []*aigrpcpb.AssetInstance{
							{Name: common.HubServiceName},
							{Name: pubsubtesting.AnotherHubServiceName},
						},
					}, nil
				}
			},
			setupFakeDepServer: func(s *pubsubtesting.FakeAssetDeploymentServer) {
				s.DeleteResourceFn = func(ctx context.Context, in *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			setupFakeOpServer: func(s *pubsubtesting.FakeOperationsServer) {
				s.GetOperationFn = func(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			expectedOutput: []string{
				fmt.Sprintf(
					"Deleting an instance of the %v service named %q",
					common.HubServiceName, common.HubServiceName),
				fmt.Sprintf(
					"Deleting an instance of the %v service named %q",
					common.HubServiceName, pubsubtesting.AnotherHubServiceName),
				fmt.Sprintf(
					"The %v option is enabled, won't try to uninstall the %v service asset.",
					common.KeyRetainServiceAsset, common.HubServiceName),
			},
			shouldRetainServiceAsset: true,
			expectErr:                false,
			packageName:              common.HubServicePackage,
			serviceName:              common.HubServiceName,
		},
		{
			name: "service_not_installed",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{}, nil
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.GetInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
					return nil, grpcstatus.Error(codes.NotFound, "not found")
				}
			},
			expectedOutput: []string{
				fmt.Sprintf("%v service asset is not installed, nothing else to do.", common.HubServiceName),
			},
			shouldRetainServiceAsset: false,
			expectErr:                false,
			packageName:              common.HubServicePackage,
			serviceName:              common.HubServiceName,
		},
		{
			name: "get_asset_error",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return nil, errors.New("backend down")
				}
			},
			expectErr:                true,
			shouldRetainServiceAsset: false,
			expectErrContains:        "backend down",
			packageName:              common.HubServicePackage,
			serviceName:              common.HubServiceName,
		},
		{
			name: "delete_resource_error",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{AssetInstances: []*aigrpcpb.AssetInstance{{Name: common.HubServiceName}}}, nil
				}
			},
			setupFakeDepServer: func(s *pubsubtesting.FakeAssetDeploymentServer) {
				s.DeleteResourceFn = func(ctx context.Context, in *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
					return nil, errors.New("delete error")
				}
			},
			expectErr:                true,
			shouldRetainServiceAsset: false,
			expectErrContains:        fmt.Sprintf("could not delete instance of the %v service", common.HubServiceName),
			packageName:              common.HubServicePackage,
			serviceName:              common.HubServiceName,
		},
		{
			name: "get_installed_asset_error",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{}, nil
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.GetInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
					return nil, errors.New("installed asset backend down")
				}
			},
			expectErr:                true,
			shouldRetainServiceAsset: false,
			expectErrContains:        fmt.Sprintf("failed to determine whether the %v service is installed", common.HubServiceName),
			packageName:              common.HubServicePackage,
			serviceName:              common.HubServiceName,
		},
		{
			name: "delete_resource_wait_operation_error",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{AssetInstances: []*aigrpcpb.AssetInstance{{Name: common.HubServiceName}}}, nil
				}
			},
			setupFakeDepServer: func(s *pubsubtesting.FakeAssetDeploymentServer) {
				s.DeleteResourceFn = func(ctx context.Context, in *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: false, Name: "op1"}, nil
				}
			},
			setupFakeOpServer: func(s *pubsubtesting.FakeOperationsServer) {
				s.GetOperationFn = func(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
					return nil, errors.New("operation failed")
				}
			},
			expectErr:                true,
			shouldRetainServiceAsset: false,
			expectErrContains:        "operation failed",
			packageName:              common.HubServicePackage,
			serviceName:              common.HubServiceName,
		},
		{
			name: "uninstall_service_asset_error",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{AssetInstances: []*aigrpcpb.AssetInstance{{Name: common.HubServiceName}}}, nil
				}
			},
			setupFakeDepServer: func(s *pubsubtesting.FakeAssetDeploymentServer) {
				s.DeleteResourceFn = func(ctx context.Context, in *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			setupFakeOpServer: func(s *pubsubtesting.FakeOperationsServer) {
				s.GetOperationFn = func(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.DeleteInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.DeleteInstalledAssetRequest) (*lropb.Operation, error) {
					return nil, errors.New("uninstall error")
				}
			},
			expectErr:                true,
			shouldRetainServiceAsset: false,
			expectErrContains:        fmt.Sprintf("failed to uninstall %v service asset", common.HubServiceName),
			packageName:              common.HubServicePackage,
			serviceName:              common.HubServiceName,
		},
		{
			name: "uninstall_service_asset_operation_error",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{AssetInstances: []*aigrpcpb.AssetInstance{{Name: common.HubServiceName}}}, nil
				}
			},
			setupFakeDepServer: func(s *pubsubtesting.FakeAssetDeploymentServer) {
				s.DeleteResourceFn = func(ctx context.Context, in *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			setupFakeOpServer: func(s *pubsubtesting.FakeOperationsServer) {
				s.GetOperationFn = func(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
					if in.Name == "op1" {
						return &lropb.Operation{Done: true, Name: "op1"}, nil
					}
					return nil, errors.New("uninstall op failed")
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.DeleteInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.DeleteInstalledAssetRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: false, Name: "op2"}, nil
				}
			},
			expectErr:                true,
			shouldRetainServiceAsset: false,
			expectErrContains:        "uninstall op failed",
			packageName:              common.HubServicePackage,
			serviceName:              common.HubServiceName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := pubsubtesting.SetupTestServer(t)

			if tt.setupFakeInstServer != nil {
				tt.setupFakeInstServer(res.InstServer)
			}
			if tt.setupFakeDepServer != nil {
				tt.setupFakeDepServer(res.DepServer)
			}
			if tt.setupFakeOpServer != nil {
				tt.setupFakeOpServer(res.OpServer)
			}
			if tt.setupFakeIAServer != nil {
				tt.setupFakeIAServer(res.IaServer)
			}

			ctx := context.Background()
			conn, err := pubsubtesting.DialTestServer(ctx, res.Listener)
			if err != nil {
				t.Fatalf("Failed to connect to the test server: %v", err)
			}
			defer conn.Close()

			var buf bytes.Buffer
			runner := &ServiceDeleteCmdRunner{
				CmdRunnerBase: *common.NewCmdRunnerBase(
					conn,
					&buf,
					"testcluster",
					tt.packageName,
					tt.serviceName),
				ShouldRetainServiceAsset: tt.shouldRetainServiceAsset,
			}

			err = runner.Run(ctx)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.expectErrContains)
				}
				if !strings.Contains(err.Error(), tt.expectErrContains) {
					t.Errorf("expected error to contain %q, got %v", tt.expectErrContains, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				for _, expectedOutputFragment := range tt.expectedOutput {
					if !strings.Contains(buf.String(), expectedOutputFragment) {
						t.Errorf("expected output to contain %q, got %q", expectedOutputFragment, buf.String())
					}
				}
			}
		})
	}
}
