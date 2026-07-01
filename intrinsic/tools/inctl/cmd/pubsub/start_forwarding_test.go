// Copyright 2023 Intrinsic Innovation LLC

package pubsub

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
	pubsubtesting "intrinsic/tools/inctl/cmd/pubsub/testing"
)

func TestStartForwarding(t *testing.T) {
	tests := []struct {
		name                string
		topics              []string
		kvStorePaths        []string
		setupFakeInstServer func(s *pubsubtesting.FakeAssetInstancesServer)
		setupFakeDepServer  func(s *pubsubtesting.FakeAssetDeploymentServer)
		setupFakeOpServer   func(s *pubsubtesting.FakeOperationsServer)
		setupFakeIAServer   func(s *pubsubtesting.FakeInstalledAssetsServer)
		expectedOutput      []string
		expectErr           bool
		expectErrContains   string
	}{
		{
			name:         "successful_install",
			topics:       []string{"topic1"},
			kvStorePaths: []string{"path1"},
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{}, nil
				}
			},
			setupFakeDepServer: func(s *pubsubtesting.FakeAssetDeploymentServer) {
				s.CreateResourceFromCatalogFn = func(ctx context.Context, in *adgrpcpb.CreateResourceFromCatalogRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.GetInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
					return nil, grpcstatus.Error(codes.NotFound, "not found") // Trigger install
				}
				s.CreateInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.CreateInstalledAssetRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op2"}, nil
				}
			},
			expectedOutput: []string{fmt.Sprintf("Successfully installed %v", forwardingServiceName)},
			expectErr:      false,
		},
		{
			name:         "successful_update",
			topics:       []string{"topic1"},
			kvStorePaths: []string{"path1"},
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{
						AssetInstances: []*aigrpcpb.AssetInstance{
							{Name: forwardingServiceName},
						},
					}, nil
				}
			},
			setupFakeDepServer: func(s *pubsubtesting.FakeAssetDeploymentServer) {
				s.CreateResourceFromCatalogFn = func(ctx context.Context, in *adgrpcpb.CreateResourceFromCatalogRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
				s.DeleteResourceFn = func(ctx context.Context, in *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.GetInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
					return nil, grpcstatus.Error(codes.NotFound, "not found")
				}
				s.CreateInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.CreateInstalledAssetRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op3"}, nil
				}
			},
			expectedOutput: []string{
				fmt.Sprintf("Deleting an instance of the %v service named %q", forwardingServiceName, forwardingServiceName),
				fmt.Sprintf("Successfully deleted an instance of the %v service", forwardingServiceName),
				fmt.Sprintf("Successfully installed %v", forwardingServiceName),
			},
			expectErr: false,
		},
		{
			name:              "nothing_to_forward",
			topics:            []string{},
			kvStorePaths:      []string{},
			expectErr:         true,
			expectErrContains: "no topics or KV store keys specified to forward",
		},
		{
			name:   "get_asset_error",
			topics: []string{"topic1"},
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return nil, errors.New("backend down")
				}
			},
			expectErr:         true,
			expectErrContains: "backend down",
		},
		{
			name:   "create_resource_error",
			topics: []string{"topic1"},
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{}, nil
				}
			},
			setupFakeDepServer: func(s *pubsubtesting.FakeAssetDeploymentServer) {
				s.CreateResourceFromCatalogFn = func(ctx context.Context, in *adgrpcpb.CreateResourceFromCatalogRequest) (*lropb.Operation, error) {
					return nil, errors.New("create error")
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.GetInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
					return nil, grpcstatus.Error(codes.NotFound, "not found") // Trigger install
				}
				s.CreateInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.CreateInstalledAssetRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op2"}, nil
				}
			},
			expectErr:         true,
			expectErrContains: "could not create resource",
		},
		{
			name:   "create_resource_wait_operation_error",
			topics: []string{"topic1"},
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{}, nil
				}
			},
			setupFakeDepServer: func(s *pubsubtesting.FakeAssetDeploymentServer) {
				s.CreateResourceFromCatalogFn = func(ctx context.Context, in *adgrpcpb.CreateResourceFromCatalogRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: false, Name: "op1"}, nil
				}
			},
			setupFakeOpServer: func(s *pubsubtesting.FakeOperationsServer) {
				s.GetOperationFn = func(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
					return nil, errors.New("operation failed")
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.GetInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
					return nil, grpcstatus.Error(codes.NotFound, "not found") // Trigger install
				}
				s.CreateInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.CreateInstalledAssetRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op2"}, nil
				}
			},
			expectErr:         true,
			expectErrContains: "operation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := setupTestServer(t)

			if tt.setupFakeInstServer != nil {
				tt.setupFakeInstServer(res.instServer)
			}
			if tt.setupFakeDepServer != nil {
				tt.setupFakeDepServer(res.depServer)
			}
			if tt.setupFakeOpServer != nil {
				tt.setupFakeOpServer(res.opServer)
			}
			if tt.setupFakeIAServer != nil {
				tt.setupFakeIAServer(res.iaServer)
			}

			ctx := t.Context()
			conn, err := dialTestServer(ctx, res.listener)
			if err != nil {
				t.Fatalf("Failed to connect to the test server: %v", err)
			}
			defer conn.Close()

			var buf bytes.Buffer
			runner := &StartForwardingCmdRunner{
				ServiceInstallingCmdRunner: ServiceInstallingCmdRunner{
					CmdRunnerBase: *newCmdRunnerBase(
						conn,
						&buf,
						"testcluster",
						forwardingServicePackage,
						forwardingServiceName),
					requestedVersion: defaultForwardingServiceVersion,
				},
				topics:       tt.topics,
				kvStorePaths: tt.kvStorePaths,
			}

			err = runner.run(ctx)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("runner.run() returned nil, want error containing %q", tt.expectErrContains)
				}
				if !strings.Contains(err.Error(), tt.expectErrContains) {
					t.Errorf("runner.run() returned %v, want error containing %q", err, tt.expectErrContains)
				}
			} else {
				if err != nil {
					t.Fatalf("runner.run() returned %v, want nil", err)
				}
				for _, expectedOutputFragment := range tt.expectedOutput {
					if !strings.Contains(buf.String(), expectedOutputFragment) {
						t.Errorf("buf.String() = %q, want string containing %q", buf.String(), expectedOutputFragment)
					}
				}
			}
		})
	}
}
