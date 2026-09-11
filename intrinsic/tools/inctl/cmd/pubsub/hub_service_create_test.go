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

package pubsub

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"

	adgrpcpb "intrinsic/assets/proto/asset_deployment_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	aigrpcpb "intrinsic/assets/proto/v1/asset_instances_go_proto"
	provisionerpb "intrinsic/platform/pubsub/cloud_router_provisioner/v1/provisioner_go_proto"
	endpointpb "intrinsic/platform/pubsub/connect/onprem/relay_router_service/endpoint_spec_go_proto"
	pubsubtesting "intrinsic/tools/inctl/cmd/pubsub/testing"
)

const (
	testHubServiceVersion = "0.0.1"
)

func verifyExpectedOutputAndError(t *testing.T, receivedOutput *bytes.Buffer, receivedError error, expectError bool, expectErrorContains string, expectedOutput []string) {
	t.Helper()
	if expectError {
		if receivedError == nil {
			t.Fatalf("expected error containing %q, got nil", expectErrorContains)
		}
		if !strings.Contains(receivedError.Error(), expectErrorContains) {
			t.Errorf("expected error to contain %q, got %v", expectErrorContains, receivedError)
		}
	} else {
		if receivedError != nil {
			t.Fatalf("expected no error, got %v", receivedError)
		}
		for _, expectedOutputFragment := range expectedOutput {
			if !strings.Contains(receivedOutput.String(), expectedOutputFragment) {
				t.Errorf("expected output to contain %q, got %q", expectedOutputFragment, receivedOutput.String())
			}
		}
	}
}

func configureTestServerForSuccessfulOnpremInstallation(res *testServerResources) {
	res.instServer.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
		return &aigrpcpb.ListAssetInstancesResponse{}, nil
	}
	res.depServer.CreateResourceFromCatalogFn = func(ctx context.Context, in *adgrpcpb.CreateResourceFromCatalogRequest) (*lropb.Operation, error) {
		return &lropb.Operation{Done: true, Name: "op1"}, nil
	}
	res.opServer.GetOperationFn = func(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
		return &lropb.Operation{Done: true, Name: "op1"}, nil
	}
	res.iaServer.GetInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
		return nil, grpcstatus.Error(codes.NotFound, "not found")
	}
	res.iaServer.CreateInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.CreateInstalledAssetRequest) (*lropb.Operation, error) {
		return &lropb.Operation{Done: true, Name: "op2"}, nil
	}
}

func TestInstallationAndUpdateOfOnpremService(t *testing.T) {
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
			expectedOutput: []string{
				fmt.Sprintf("Successfully installed %v", hubServiceName),
				fmt.Sprintf("Successfully added an instance of the %v service", hubServiceName),
			},
			expectErr: false,
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
				fmt.Sprintf("Deleting an instance of the %v service named %q", hubServiceName, hubServiceName),
				fmt.Sprintf("Deleting an instance of the %v service named %q", hubServiceName, anotherHubServiceName),
				fmt.Sprintf("Successfully deleted an instance of the %v service", hubServiceName),
				fmt.Sprintf("Successfully installed %v", hubServiceName),
				fmt.Sprintf("Successfully added an instance of the %v service", hubServiceName),
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
							IdVersion: &idpb.IdVersion{Version: testHubServiceVersion},
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
			expectErrContains: fmt.Sprintf("could not delete instance of the %v service", hubServiceName),
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
			expectErrContains: fmt.Sprintf("failed to determine version of the %v service asset", hubServiceName),
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
			expectErrContains: fmt.Sprintf("could not install %v service asset", hubServiceName),
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
							IdVersion: &idpb.IdVersion{Version: testHubServiceVersion},
						},
					}, nil
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
			var buf bytes.Buffer
			runner := &HubServiceCreateRunner{
				hubEndpoint:    "testcluster@local",
				spokeEndpoints: tt.spokeWorkcells,
				dialOnpremCluster: func(ctx context.Context, project string, org string, cluster string) (context.Context, *grpc.ClientConn, string, error) {
					conn, err := dialTestServer(ctx, res.listener)
					if err != nil {
						t.Fatalf("Failed to connect to the test server: %v", err)
					}
					// Not closing the connection here because it will be closed by the calling code.
					return ctx, conn, "", err
				},
				getServiceVersionToInstall: func(packageName string, serviceName string) (string, error) {
					return testHubServiceVersion, nil
				},
			}

			err := runner.run(ctx, &buf)
			verifyExpectedOutputAndError(t, &buf, err, tt.expectErr, tt.expectErrContains, tt.expectedOutput)
		})
	}
}

func TestParseEndpointSpecs(t *testing.T) {
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
			name:              "Parse Endpoint Spec Error",
			spokeEndpoints:    []string{"invalid"},
			expectErr:         true,
			expectErrContains: "Each endpoint spec should consist of two parts separated by @",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpointSpecs, err := parseEndpointSpecs(tt.spokeEndpoints)

			if tt.expectErr {
				if err == nil {
					t.Errorf("parseEndpointSpecs succeeded when it was supposed to fail")
				}
				if !strings.Contains(err.Error(), tt.expectErrContains) {
					t.Errorf("expected error to contain %q, got %v", tt.expectErrContains, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			} else {
				if len(endpointSpecs) != len(tt.expectedSpokes) {
					t.Errorf("expected %d spoke endpoints, got %d", len(tt.expectedSpokes), len(endpointSpecs))
				}

				for i, spec := range endpointSpecs {
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

func TestCreatingLocalOnlyNetwork(t *testing.T) {
	tests := []struct {
		name                string
		spokeEndpoints      []string
		expectedOutput      []string
		expectError         bool
		expectErrorContains string
	}{
		{
			name:           "Successful creation",
			spokeEndpoints: []string{"spoke1@local"},
			expectedOutput: []string{
				"Creating a local-only line orchestration network with the hub at \"test-hub\"",
				fmt.Sprintf("Successfully installed %v", hubServiceName),
				fmt.Sprintf("Successfully added an instance of the %v service", hubServiceName),
			},
			expectError: false,
		},
		{
			name:                "Invalid endpoint spec",
			spokeEndpoints:      []string{"zzz"},
			expectError:         true,
			expectErrorContains: "Failed to parse \"zzz\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := setupTestServer(t)
			configureTestServerForSuccessfulOnpremInstallation(res)

			ctx := t.Context()
			var buf bytes.Buffer
			runner := &HubServiceCreateRunner{
				hubEndpoint:    "test-hub@local",
				spokeEndpoints: tt.spokeEndpoints,
				dialOnpremCluster: func(ctx context.Context, project string, org string, cluster string) (context.Context, *grpc.ClientConn, string, error) {
					conn, err := dialTestServer(ctx, res.listener)
					if err != nil {
						t.Fatalf("Failed to connect to the test server: %v", err)
					}
					// Not closing the connection here because it will be closed by the calling code.
					return ctx, conn, "", err
				},
				getServiceVersionToInstall: func(packageName string, serviceName string) (string, error) {
					if packageName == hubServicePackage && serviceName == hubServiceName {
						return testHubServiceVersion, nil
					}
					return "", fmt.Errorf(
						"unexpected parameters of getServiceVersionToInstall: got(%q, %q), want (%q, %q)",
						packageName, serviceName,
						hubServicePackage, hubServiceName,
					)
				},
			}

			err := runner.run(ctx, &buf)
			verifyExpectedOutputAndError(t, &buf, err, tt.expectError, tt.expectErrorContains, tt.expectedOutput)
		})
	}
}

func TestCreatingMixedNetwork(t *testing.T) {
	tests := []struct {
		name                      string
		hubEndpoint               string
		spokeEndpoints            []string
		provisioningShouldSucceed bool
		expectedOutput            []string
		expectError               bool
		expectErrorContains       string
	}{
		{
			name:                      "All VMs",
			hubEndpoint:               "vmp-123@remote",
			spokeEndpoints:            []string{"vmp-234@remote"},
			provisioningShouldSucceed: true,
			expectedOutput: []string{
				"Creating a mixed line orchestration network. Line id: \"intrinsic-vmp-123\"",
				"Line-level router has been provisioned",
				"Installing the relay router on \"vmp-123\"",
				"Installing the relay router on \"vmp-234\"",
				"All relay routers has been installed",
			},
			expectError: false,
		},
		{
			name:                      "Provisioning failure",
			hubEndpoint:               "vmp-123@remote",
			spokeEndpoints:            []string{"vmp-234@remote"},
			provisioningShouldSucceed: false,
			expectError:               true,
			expectErrorContains:       "simulated provisioning failure",
		},
		{
			name:                      "Only hub on VM",
			hubEndpoint:               "vmp-123@remote",
			spokeEndpoints:            []string{"node-234@local"},
			provisioningShouldSucceed: true,
			expectedOutput: []string{
				"Creating a mixed line orchestration network. Line id: \"intrinsic-vmp-123\"",
				"Line-level router has been provisioned",
				"Installing the relay router on \"vmp-123\"",
				"Installing the relay router on \"node-234\"",
				"All relay routers has been installed",
			},
			expectError: false,
		},
		{
			name:                      "Hub workcell and spoke VM",
			hubEndpoint:               "node-234@local",
			spokeEndpoints:            []string{"vmp-123@remote"},
			provisioningShouldSucceed: true,
			expectedOutput: []string{
				"Creating a mixed line orchestration network. Line id: \"intrinsic-node-234\"",
				"Line-level router has been provisioned",
				"Installing the relay router on \"vmp-123\"",
				"Installing the relay router on \"node-234\"",
				"All relay routers has been installed",
			},
			expectError: false,
		},
		{
			name:                      "Invalid endpoint spec",
			hubEndpoint:               "vmp-123@remote",
			spokeEndpoints:            []string{"zzz"},
			provisioningShouldSucceed: true,
			expectError:               true,
			expectErrorContains:       "Failed to parse \"zzz\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := setupTestServer(t)
			configureTestServerForSuccessfulOnpremInstallation(res)
			res.provisionerServer.ProvisionFn = func(ctx context.Context, in *provisionerpb.ProvisionLineRouterRequest) (*provisionerpb.LineRouterProvisionResponse, error) {
				if tt.provisioningShouldSucceed {
					return &provisionerpb.LineRouterProvisionResponse{}, nil
				}

				return nil, fmt.Errorf("simulated provisioning failure")
			}

			ctx := t.Context()

			var buf bytes.Buffer
			runner := &HubServiceCreateRunner{
				org:            "intrinsic",
				hubEndpoint:    tt.hubEndpoint,
				spokeEndpoints: tt.spokeEndpoints,
				dialOnpremCluster: func(ctx context.Context, project string, org string, cluster string) (context.Context, *grpc.ClientConn, string, error) {
					conn, err := dialTestServer(ctx, res.listener)
					if err != nil {
						t.Fatalf("Failed to connect to the test server : %v", err)
					}

					// Not closing the connection here because it will be closed by the calling code.
					return ctx, conn, "", nil
				},
				dialCloudCluster: func(ctx context.Context) (*grpc.ClientConn, error) {
					conn, err := dialTestServer(ctx, res.listener)
					if err != nil {
						t.Fatalf("Failed to connect to the test server: %v", err)
					}

					// Not closing the connection here because it will be closed by the calling code.
					return conn, nil
				},
				getServiceVersionToInstall: func(packageName string, serviceName string) (string, error) {
					if packageName == onpremToLineRouterRelayServicePackage && serviceName == onpremToLineRouterRelayServiceName {
						return testHubServiceVersion, nil
					}
					return "", fmt.Errorf(
						"unexpected parameters of getServiceVersionToInstall: got(%q, %q), want (%q, %q)",
						packageName, serviceName,
						onpremToLineRouterRelayServicePackage, onpremToLineRouterRelayServiceName,
					)
				},
			}

			err := runner.run(ctx, &buf)
			verifyExpectedOutputAndError(t, &buf, err, tt.expectError, tt.expectErrorContains, tt.expectedOutput)
		})
	}
}

func TestGetHubEndpoint(t *testing.T) {
	const deprecationWarning = "WARNING: the --cluster flag is deprecated, use --hub-endpoint instead."
	tests := []struct {
		name                string
		cluster             string
		hubEndpoint         string
		expectedResult      string
		expectError         bool
		expectErrorContains string
		expectedOutput      []string
	}{
		{
			name:           "VM cluster",
			cluster:        "vmp-123",
			hubEndpoint:    "",
			expectedResult: "vmp-123@remote",
			expectError:    false,
			expectedOutput: []string{
				deprecationWarning,
				"Assuming that \"vmp-123\" is a remote server",
			},
		},
		{
			name:           "Workcell cluster",
			cluster:        "node-123-456",
			hubEndpoint:    "",
			expectedResult: "node-123-456@local",
			expectError:    false,
			expectedOutput: []string{
				deprecationWarning,
				"Assuming that \"node-123-456\" is on the local network",
			},
		},
		{
			name:           "Hub endpoint has higher priority",
			cluster:        "vmp-123",
			hubEndpoint:    "vmp-234@remote",
			expectedResult: "vmp-234@remote",
			expectError:    false,
			expectedOutput: []string{}, // No output when hub endpoint is specified
		},
		{
			name:           "Hub endpoint and no cluster",
			cluster:        "",
			hubEndpoint:    "vmp-345@remote",
			expectedResult: "vmp-345@remote",
			expectError:    false,
			expectedOutput: []string{}, // No output when hub endpoint is specified
		},
		{
			name:                "None are specified",
			cluster:             "",
			hubEndpoint:         "",
			expectError:         true,
			expectErrorContains: "neither cluster nor hub endpoint are specified",
			expectedOutput:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			got, err := getHubEndpoint(&buf, tt.cluster, tt.hubEndpoint)
			if got != tt.expectedResult {
				t.Errorf("Unexpected hub endpoint: got %q, want %q", got, tt.expectedResult)
			}
			verifyExpectedOutputAndError(t, &buf, err, tt.expectError, tt.expectErrorContains, tt.expectedOutput)
		})
	}
}
