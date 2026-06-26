// Copyright 2023 Intrinsic Innovation LLC

package pubsub

import (
	"bufio"
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
	idpb "intrinsic/assets/proto/id_go_proto"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	aigrpcpb "intrinsic/assets/proto/v1/asset_instances_go_proto"
	endpointpb "intrinsic/platform/pubsub/connect/cloud/proto/v1alpha1/endpoint_spec_go_proto"
	pubsubtesting "intrinsic/tools/inctl/cmd/pubsub/testing"
)

func TestHubServiceCreateRunE(t *testing.T) {
	tests := []struct {
		name                string
		spokeWorkcells      []string
		setupFakeInstServer func(s *pubsubtesting.FakeAssetInstancesServer)
		setupFakeDepServer  func(s *pubsubtesting.FakeAssetDeploymentServer)
		setupFakeOpServer   func(s *pubsubtesting.FakeOperationsServer)
		setupFakeIAServer   func(s *pubsubtesting.FakeInstalledAssetsServer)
		expectedOutput      []string
		expectErr           bool
		expectErrContains   string
	}{
		{
			name:           "Successful Install",
			spokeWorkcells: []string{"spoke1@local"},
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
			setupFakeOpServer: func(s *pubsubtesting.FakeOperationsServer) {
				s.GetOperationFn = func(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.GetInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
					return nil, grpcstatus.Error(codes.NotFound, "not found")
				}
				s.CreateInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.CreateInstalledAssetRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op2"}, nil
				}
			},
			expectedOutput: []string{"Successfully installed line orchestration relay"},
			expectErr:      false,
		},
		{
			name:           "Successful Update",
			spokeWorkcells: []string{"spoke1@local"},
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
				s.CreateResourceFromCatalogFn = func(ctx context.Context, in *adgrpcpb.CreateResourceFromCatalogRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
				s.DeleteResourceFn = func(ctx context.Context, in *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op1"}, nil
				}
			},
			setupFakeOpServer: func(s *pubsubtesting.FakeOperationsServer) {
				s.GetOperationFn = func(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op2"}, nil
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.GetInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
					return nil, grpcstatus.Error(codes.NotFound, "not found") // Return not found so it triggers install
				}
				s.CreateInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.CreateInstalledAssetRequest) (*lropb.Operation, error) {
					return &lropb.Operation{Done: true, Name: "op3"}, nil
				}
			},
			expectedOutput: []string{
				fmt.Sprintf("Deleting an instance of the relay service named %q", hubServiceName),
				fmt.Sprintf("Deleting an instance of the relay service named %q", anotherHubServiceName),
				"Successfully deleted an instance of the relay service",
				"Successfully installed line orchestration relay",
			},
			expectErr: false,
		},
		{
			name:           "GetAsset Error",
			spokeWorkcells: []string{"spoke1@local"},
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return nil, errors.New("backend down")
				}
			},
			expectErr:         true,
			expectErrContains: "backend down",
		},
		{
			name:              "No Spoke Workcells",
			spokeWorkcells:    []string{},
			expectErr:         true,
			expectErrContains: "at least one spoke endpoint must be specified using --spoke-endpoint",
		},
		{
			name:           "CreateResource Error",
			spokeWorkcells: []string{"spoke1@local"},
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
					// Simulate already installed with the correct version so it bypasses install and goes to create resource directly
					return &iagrpcpb.InstalledAsset{
						Metadata: &metadatapb.Metadata{
							IdVersion: &idpb.IdVersion{Version: defaultHubServiceVersion},
						},
					}, nil
				}
			},
			expectErr:         true,
			expectErrContains: "could not create resource",
		},
		{
			name:           "DeleteResource Error",
			spokeWorkcells: []string{"spoke1@local"},
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
			expectErr:         true,
			expectErrContains: "could not delete instance of the relay service",
		},
		{
			name:           "GetInstalledAsset Error",
			spokeWorkcells: []string{"spoke1@local"},
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{}, nil
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.GetInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
					return nil, errors.New("backend down")
				}
			},
			expectErr:         true,
			expectErrContains: "failed to determine version of the relay service asset",
		},
		{
			name:           "CreateInstalledAsset Error",
			spokeWorkcells: []string{"spoke1@local"},
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{}, nil
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.GetInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
					return nil, grpcstatus.Error(codes.NotFound, "not found")
				}
				s.CreateInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.CreateInstalledAssetRequest) (*lropb.Operation, error) {
					return nil, errors.New("install error")
				}
			},
			expectErr:         true,
			expectErrContains: "could not install relay service asset",
		},
		{
			name:           "CreateResource Wait Operation Error",
			spokeWorkcells: []string{"spoke1@local"},
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
					return &iagrpcpb.InstalledAsset{
						Metadata: &metadatapb.Metadata{
							IdVersion: &idpb.IdVersion{Version: defaultHubServiceVersion},
						},
					}, nil
				}
			},
			expectErr:         true,
			expectErrContains: "operation failed",
		},
		{
			name:           "Invalid Endpoint Spec",
			spokeWorkcells: []string{"invalid"},
			setupFakeInstServer: func(s *pubsubtesting.FakeAssetInstancesServer) {
				s.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
					return &aigrpcpb.ListAssetInstancesResponse{}, nil
				}
			},
			setupFakeIAServer: func(s *pubsubtesting.FakeInstalledAssetsServer) {
				s.GetInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
					return &iagrpcpb.InstalledAsset{
						Metadata: &metadatapb.Metadata{
							IdVersion: &idpb.IdVersion{Version: defaultHubServiceVersion},
						},
					}, nil
				}
			},
			expectErr:         true,
			expectErrContains: "failed to create service config from command line flags",
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
			runner := &HubServiceCreateCmdRunner{
				HubServiceCmdRunnerBase: *newHubServiceCmdRunnerBase(conn, &buf, "testcluster"),
				requestedVersion:        defaultHubServiceVersion,
				spokeEndpoints:          tt.spokeWorkcells,
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

func TestMakeConfig(t *testing.T) {
	tests := []struct {
		name              string
		spokeEndpoints    []string
		expectedSpokes    []string
		expectErr         bool
		expectErrContains string
		validateSpoke     func(t *testing.T, spec *endpointpb.EndpointSpec)
	}{
		{
			name:           "Local Endpoint",
			spokeEndpoints: []string{"spoke1@local"},
			expectedSpokes: []string{"spoke1"},
			expectErr:      false,
			validateSpoke: func(t *testing.T, spec *endpointpb.EndpointSpec) {
				if spec.GetLocal() == nil {
					t.Errorf("expected ConnectionSpec to be LocalConnectionSpec, got %T", spec.ConnectionSpec)
				}
			},
		},
		{
			name:           "URL Endpoint",
			spokeEndpoints: []string{"spoke1@custom-router.app-my-namespace.svc.cluster.local:7447"},
			expectedSpokes: []string{"spoke1"},
			expectErr:      false,
			validateSpoke: func(t *testing.T, spec *endpointpb.EndpointSpec) {
				if spec.GetUrl() != "custom-router.app-my-namespace.svc.cluster.local:7447" {
					t.Errorf("expected URL to be custom-router.app-my-namespace.svc.cluster.local:7447, got %v", spec.GetUrl())
				}
			},
		},
		{
			name:              "Remote Endpoint Error",
			spokeEndpoints:    []string{"spoke1@remote"},
			expectErr:         true,
			expectErrContains: "remote endpoints are not supported",
		},
		{
			name:              "Parse Endpoint Spec Error",
			spokeEndpoints:    []string{"invalid"},
			expectErr:         true,
			expectErrContains: "Each endpoint spec should consist of two parts separated by @",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			outputWriter := bufio.NewWriter(&buf)
			runner := &HubServiceCreateCmdRunner{
				HubServiceCmdRunnerBase: *newHubServiceCmdRunnerBase(nil, outputWriter, "test-hub-workcell"),
				requestedVersion:        defaultHubServiceVersion,
				spokeEndpoints:          tt.spokeEndpoints,
			}

			config, err := runner.makeConfig()

			if tt.expectErr {
				if err == nil {
					t.Errorf("makeConfig succeeded when it was supposed to fail")
				}
				if !strings.Contains(err.Error(), tt.expectErrContains) {
					t.Errorf("expected error to contain %q, got %v", tt.expectErrContains, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			} else {
				if config.HubWorkcellName != "test-hub-workcell" {
					t.Errorf("expected HubWorkcellName to be 'test-hub-workcell', got %v", config.HubWorkcellName)
				}

				if len(config.SpokeEndpointSpecs) != len(tt.expectedSpokes) {
					t.Errorf("expected %d spoke endpoints, got %d", len(tt.expectedSpokes), len(config.SpokeEndpointSpecs))
				}

				for i, spec := range config.SpokeEndpointSpecs {
					if spec.WorkcellName != tt.expectedSpokes[i] {
						t.Errorf("expected workcell name %v, got %v", tt.expectedSpokes[i], spec.WorkcellName)
					}
					if tt.validateSpoke != nil {
						tt.validateSpoke(t, spec)
					}
				}
			}
		})
	}
}
