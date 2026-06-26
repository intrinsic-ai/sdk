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

func TestHubServiceDeleteRunE(t *testing.T) {
	tests := []struct {
		name                        string
		setupFakeInstServer         func(s *pubsubtesting.FakeAssetInstancesServer)
		setupFakeDepServer          func(s *pubsubtesting.FakeAssetDeploymentServer)
		setupFakeOpServer           func(s *pubsubtesting.FakeOperationsServer)
		setupFakeIAServer           func(s *pubsubtesting.FakeInstalledAssetsServer)
		expectedOutput              []string
		shouldUninstallServiceAsset bool
		expectErr                   bool
		expectErrContains           string
	}{
		{
			name: "Successful Uninstall",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{
						AssetInstances: []*aigrpcpb.AssetInstance{
							{Name: hubServiceName},
							{Name: anotherHubServiceName},
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
				fmt.Sprintf("Deleting an instance of the relay service named %q", hubServiceName),
				fmt.Sprintf("Deleting an instance of the relay service named %q", anotherHubServiceName),
				"Successfully uninstalled the relay service asset",
			},
			shouldUninstallServiceAsset: true,
			expectErr:                   false,
		},
		{
			name: "Remove instances but do not uninstall service asset",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{
						AssetInstances: []*aigrpcpb.AssetInstance{
							{Name: hubServiceName},
							{Name: anotherHubServiceName},
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
				fmt.Sprintf("Deleting an instance of the relay service named %q", hubServiceName),
				fmt.Sprintf("Deleting an instance of the relay service named %q", anotherHubServiceName),
				fmt.Sprintf("The %v option is disabled, won't try to uninstall the relay service asset.", keyUninstallServiceAsset),
			},
			shouldUninstallServiceAsset: false,
			expectErr:                   false,
		},
		{
			name: "Service Not Installed",
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
			expectedOutput:              []string{"Relay service asset is not installed, nothing else to do."},
			shouldUninstallServiceAsset: true,
			expectErr:                   false,
		},
		{
			name: "GetAsset Error",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return nil, errors.New("backend down")
				}
			},
			expectErr:                   true,
			shouldUninstallServiceAsset: true,
			expectErrContains:           "backend down",
		},
		{
			name: "DeleteResource Error",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{AssetInstances: []*aigrpcpb.AssetInstance{{Name: hubServiceName}}}, nil
				}
			},
			setupFakeDepServer: func(s *pubsubtesting.FakeAssetDeploymentServer) {
				s.DeleteResourceFn = func(ctx context.Context, in *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
					return nil, errors.New("delete error")
				}
			},
			expectErr:                   true,
			shouldUninstallServiceAsset: true,
			expectErrContains:           "could not delete instance of the relay service",
		},
		{
			name: "GetInstalledAsset Error",
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
			expectErr:                   true,
			shouldUninstallServiceAsset: true,
			expectErrContains:           "failed to determine whether the relay service is installed",
		},
		{
			name: "DeleteResource Wait Operation Error",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{AssetInstances: []*aigrpcpb.AssetInstance{{Name: hubServiceName}}}, nil
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
			expectErr:                   true,
			shouldUninstallServiceAsset: true,
			expectErrContains:           "operation failed",
		},
		{
			name: "UninstallRelayServiceAsset Error",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{AssetInstances: []*aigrpcpb.AssetInstance{{Name: hubServiceName}}}, nil
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
			expectErr:                   true,
			shouldUninstallServiceAsset: true,
			expectErrContains:           "failed to uninstall relay service asset",
		},
		{
			name: "UninstallRelayServiceAsset Operation Error",
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{AssetInstances: []*aigrpcpb.AssetInstance{{Name: hubServiceName}}}, nil
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
			expectErr:                   true,
			shouldUninstallServiceAsset: true,
			expectErrContains:           "uninstall op failed",
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

			ctx := context.Background()
			conn, err := dialTestServer(ctx, res.listener)
			if err != nil {
				t.Fatalf("Failed to connect to the test server: %v", err)
			}
			defer conn.Close()

			var buf bytes.Buffer
			runner := &HubServiceDeleteCmdRunner{
				HubServiceCmdRunnerBase:     *newHubServiceCmdRunnerBase(conn, &buf, "testcluster"),
				shouldUninstallServiceAsset: tt.shouldUninstallServiceAsset,
			}

			err = runner.run(ctx)

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
