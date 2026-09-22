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

package hubservicecreate

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

	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"

	adgrpcpb "intrinsic/assets/proto/asset_deployment_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	aigrpcpb "intrinsic/assets/proto/v1/asset_instances_go_proto"
	provisionerpb "intrinsic/platform/pubsub/cloud_router_provisioner/v1/provisioner_go_proto"
	lineconfigstoragepb "intrinsic/platform/pubsub/connect/cloud/proto/line_configuration_storage/v1/line_configuration_storage_go_proto"
	endpointpb "intrinsic/platform/pubsub/connect/common/proto/line_configuration/v1/endpoint_spec_go_proto"
	lineconfigpb "intrinsic/platform/pubsub/connect/common/proto/line_configuration/v1/line_configuration_go_proto"
	pubsubtesting "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/testing"
)

const (
	testHubServiceVersion = "0.0.1"
)

func configureTestServerForSuccessfulOnpremInstallation(res *pubsubtesting.TestServerResources) {
	res.InstServer.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
		return &aigrpcpb.ListAssetInstancesResponse{}, nil
	}
	res.DepServer.CreateResourceFromCatalogFn = func(ctx context.Context, in *adgrpcpb.CreateResourceFromCatalogRequest) (*lropb.Operation, error) {
		return &lropb.Operation{Done: true, Name: "op1"}, nil
	}
	res.OpServer.GetOperationFn = func(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
		return &lropb.Operation{Done: true, Name: "op1"}, nil
	}
	res.IaServer.GetInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
		return nil, grpcstatus.Error(codes.NotFound, "not found")
	}
	res.IaServer.CreateInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.CreateInstalledAssetRequest) (*lropb.Operation, error) {
		return &lropb.Operation{Done: true, Name: "op2"}, nil
	}
}

func TestInstallationAndUpdateOfOnpremService(t *testing.T) {
	tests := []struct {
		name                   string
		spokeWorkcells         []string
		setupFakeInstServer    func(s *pubsubtesting.FakeAssetInstancesServer)
		setupFakeDepServer     func(s *pubsubtesting.FakeAssetDeploymentServer)
		setupFakeOpServer      func(s *pubsubtesting.FakeOperationsServer)
		setupFakeIAServer      func(s *pubsubtesting.FakeInstalledAssetsServer)
		setupLineConfigStorage func(s *pubsubtesting.FakeLineConfigurationStorageServer)
		expectedOutput         []string
		expectErr              bool
		expectErrContains      string
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
				fmt.Sprintf("Successfully installed %v", common.HubServiceName),
				fmt.Sprintf("Successfully added an instance of the %v service", common.HubServiceName),
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
							{Name: common.HubServiceName},
							{Name: pubsubtesting.AnotherHubServiceName},
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
				fmt.Sprintf("Deleting an instance of the %v service named %q", common.HubServiceName, common.HubServiceName),
				fmt.Sprintf("Deleting an instance of the %v service named %q", common.HubServiceName, pubsubtesting.AnotherHubServiceName),
				fmt.Sprintf("Successfully deleted an instance of the %v service", common.HubServiceName),
				fmt.Sprintf("Successfully installed %v", common.HubServiceName),
				fmt.Sprintf("Successfully added an instance of the %v service", common.HubServiceName),
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
					return &aigrpcpb.ListAssetInstancesResponse{AssetInstances: []*aigrpcpb.AssetInstance{{Name: common.HubServiceName}}}, nil
				}
			},
			setupFakeDepServer: func(s *pubsubtesting.FakeAssetDeploymentServer) {
				s.DeleteResourceFn = func(ctx context.Context, in *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
					return nil, errors.New("delete error")
				}
			},
			expectErr:         true,
			expectErrContains: fmt.Sprintf("could not delete instance of the %v service", common.HubServiceName),
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
			expectErrContains: fmt.Sprintf("failed to determine version of the %v service asset", common.HubServiceName),
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
			expectErrContains: fmt.Sprintf("could not install %v service asset", common.HubServiceName),
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
		{
			name:           "Failure to save line config aborts creation",
			spokeWorkcells: []string{"spoke1@local"},
			setupLineConfigStorage: func(s *pubsubtesting.FakeLineConfigurationStorageServer) {
				s.SetFn = func(context.Context, *lineconfigstoragepb.UpdateLineConfigurationRequest) (*lineconfigpb.LineConfiguration, error) {
					return nil, fmt.Errorf("test error")
				}
			},
			expectErr:         true,
			expectErrContains: "test error",
		},
		{
			name:           "Unimplemented error prompts user to create local-only network",
			spokeWorkcells: []string{"spoke1@local"},
			setupLineConfigStorage: func(s *pubsubtesting.FakeLineConfigurationStorageServer) {
				s.SetFn = func(context.Context, *lineconfigstoragepb.UpdateLineConfigurationRequest) (*lineconfigpb.LineConfiguration, error) {
					return nil, grpcstatus.Errorf(codes.Unimplemented, "")
				}
			},
			expectErr:         true,
			expectErrContains: "line configuration storage is not implemented",
			expectedOutput: []string{
				fmt.Sprintf("Try running the command with the '--%v' flag", common.KeyForceLocalOnly),
			},
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
			if tt.setupLineConfigStorage != nil {
				tt.setupLineConfigStorage(res.LineConfigStorageServer)
			}

			ctx := t.Context()
			var buf bytes.Buffer
			runner := &HubServiceCreateRunner{
				org:            "test-org",
				hubEndpoint:    pubsubtesting.LocalHubEndpoint,
				spokeEndpoints: tt.spokeWorkcells,
				dialCloudCluster: func(ctx context.Context) (*grpc.ClientConn, error) {
					conn, err := pubsubtesting.DialTestServer(ctx, res.Listener)
					if err != nil {
						t.Fatalf("Failed to connect to the test server: %v", err)
					}
					// Not closing the connection here because it will be closed by the calling code.
					return conn, err
				},
				dialOnpremCluster: func(ctx context.Context, project string, org string, cluster string) (context.Context, *grpc.ClientConn, string, error) {
					conn, err := pubsubtesting.DialTestServer(ctx, res.Listener)
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
			pubsubtesting.VerifyExpectedOutputAndError(t, &buf, err, tt.expectErr, tt.expectErrContains, tt.expectedOutput)
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
			endpointSpecs, err := common.ParseEndpointSpecs(tt.spokeEndpoints)

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
				fmt.Sprintf("Successfully installed %v", common.HubServiceName),
				fmt.Sprintf("Successfully added an instance of the %v service", common.HubServiceName),
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
			res := pubsubtesting.SetupTestServer(t)
			configureTestServerForSuccessfulOnpremInstallation(res)

			ctx := t.Context()
			var buf bytes.Buffer
			runner := &HubServiceCreateRunner{
				org:            "test-org",
				hubEndpoint:    "test-hub@local",
				spokeEndpoints: tt.spokeEndpoints,
				dialCloudCluster: func(ctx context.Context) (*grpc.ClientConn, error) {
					conn, err := pubsubtesting.DialTestServer(ctx, res.Listener)
					if err != nil {
						t.Fatalf("Failed to connect to the test server: %v", err)
					}
					// Not closing the connection here because it will be closed by the calling code.
					return conn, err
				},
				dialOnpremCluster: func(ctx context.Context, project string, org string, cluster string) (context.Context, *grpc.ClientConn, string, error) {
					conn, err := pubsubtesting.DialTestServer(ctx, res.Listener)
					if err != nil {
						t.Fatalf("Failed to connect to the test server: %v", err)
					}
					// Not closing the connection here because it will be closed by the calling code.
					return ctx, conn, "", err
				},
				getServiceVersionToInstall: func(packageName string, serviceName string) (string, error) {
					if packageName == common.HubServicePackage && serviceName == common.HubServiceName {
						return testHubServiceVersion, nil
					}
					return "", fmt.Errorf(
						"unexpected parameters of getServiceVersionToInstall: got(%q, %q), want (%q, %q)",
						packageName, serviceName,
						common.HubServicePackage, common.HubServiceName,
					)
				},
			}

			err := runner.run(ctx, &buf)
			pubsubtesting.VerifyExpectedOutputAndError(t, &buf, err, tt.expectError, tt.expectErrorContains, tt.expectedOutput)
		})
	}
}

func TestCreatingMixedNetwork(t *testing.T) {
	tests := []struct {
		name                        string
		hubEndpoint                 string
		spokeEndpoints              []string
		provisioningShouldSucceed   bool
		saveLineConfigShouldSucceed bool
		forceLocalOnly              bool
		onpremPackageName           string
		onpremServiceName           string
		expectedOutput              []string
		unexpectedOutput            []string
		expectError                 bool
		expectErrorContains         string
	}{
		{
			name:                        "All VMs",
			hubEndpoint:                 "vmp-123@remote",
			spokeEndpoints:              []string{"vmp-234@remote"},
			provisioningShouldSucceed:   true,
			saveLineConfigShouldSucceed: true,
			onpremPackageName:           common.OnpremToLineRouterRelayServicePackage,
			onpremServiceName:           common.OnpremToLineRouterRelayServiceName,
			expectedOutput: []string{
				"Creating a mixed line orchestration network. Line id: \"intrinsic-vmp-123\"",
				"Line-level router has been provisioned",
				"Installing the relay router on \"vmp-123\"",
				"Installing the relay router on \"vmp-234\"",
				"All relay routers have been installed",
			},
			expectError: false,
		},
		{
			name:                        "Provisioning failure",
			hubEndpoint:                 "vmp-123@remote",
			spokeEndpoints:              []string{"vmp-234@remote"},
			provisioningShouldSucceed:   false,
			saveLineConfigShouldSucceed: true,
			onpremPackageName:           common.OnpremToLineRouterRelayServicePackage,
			onpremServiceName:           common.OnpremToLineRouterRelayServiceName,
			expectError:                 true,
			expectErrorContains:         "simulated provisioning failure",
		},
		{
			name:                        "Only hub on VM",
			hubEndpoint:                 "vmp-123@remote",
			spokeEndpoints:              []string{"node-234@local"},
			provisioningShouldSucceed:   true,
			saveLineConfigShouldSucceed: true,
			onpremPackageName:           common.OnpremToLineRouterRelayServicePackage,
			onpremServiceName:           common.OnpremToLineRouterRelayServiceName,
			expectedOutput: []string{
				"Creating a mixed line orchestration network. Line id: \"intrinsic-vmp-123\"",
				"Line-level router has been provisioned",
				"Installing the relay router on \"vmp-123\"",
				"Installing the relay router on \"node-234\"",
				"All relay routers have been installed",
			},
			expectError: false,
		},
		{
			name:                        "Hub workcell and spoke VM",
			hubEndpoint:                 "node-234@local",
			spokeEndpoints:              []string{"vmp-123@remote"},
			provisioningShouldSucceed:   true,
			saveLineConfigShouldSucceed: true,
			onpremPackageName:           common.OnpremToLineRouterRelayServicePackage,
			onpremServiceName:           common.OnpremToLineRouterRelayServiceName,
			expectedOutput: []string{
				"Creating a mixed line orchestration network. Line id: \"intrinsic-node-234\"",
				"Line-level router has been provisioned",
				"Installing the relay router on \"vmp-123\"",
				"Installing the relay router on \"node-234\"",
				"All relay routers have been installed",
			},
			expectError: false,
		},
		{
			name:                        "Invalid endpoint spec",
			hubEndpoint:                 "vmp-123@remote",
			spokeEndpoints:              []string{"zzz"},
			provisioningShouldSucceed:   true,
			saveLineConfigShouldSucceed: true,
			expectError:                 true,
			expectErrorContains:         "Failed to parse \"zzz\"",
		},
		{
			name:                        "force-local-only flag disables interactions with cloud services",
			hubEndpoint:                 "node-123@local",
			spokeEndpoints:              []string{"node-234@local"},
			provisioningShouldSucceed:   false,
			saveLineConfigShouldSucceed: false,
			forceLocalOnly:              true,
			onpremPackageName:           common.HubServicePackage,
			onpremServiceName:           common.HubServiceName,
			expectError:                 false,
			expectedOutput: []string{
				"Creating a local-only line orchestration network without saving its configuration in the cloud",
				"Creating a local-only line orchestration network with the hub at \"node-123\"",
			},
			unexpectedOutput: []string{
				"Connecting to the cloud API endpoint",
				"Saving line configuration",
				"Provisioning a line-level router",
			},
		},
		{
			name:                "force-local-only flag is not compatible with remote endpoints",
			hubEndpoint:         pubsubtesting.LocalHubEndpoint,
			spokeEndpoints:      []string{pubsubtesting.RemoteSpokeEndpoint},
			forceLocalOnly:      true,
			expectError:         true,
			expectErrorContains: fmt.Sprintf("The '--%v' flag is present, but at least one endpoint is remote", common.KeyForceLocalOnly),
		},
		{
			name:                "Invalid endpoint spec is rejected",
			hubEndpoint:         "node-###@local",
			spokeEndpoints:      []string{pubsubtesting.RemoteSpokeEndpoint},
			expectError:         true,
			expectErrorContains: "contains invalid characters",
			expectedOutput: []string{
				"Cannot create a valid line id",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := pubsubtesting.SetupTestServer(t)
			configureTestServerForSuccessfulOnpremInstallation(res)
			res.ProvisionerServer.ProvisionFn = func(ctx context.Context, in *provisionerpb.ProvisionLineRouterRequest) (*provisionerpb.ProvisionLineRouterResponse, error) {
				if tt.provisioningShouldSucceed {
					return &provisionerpb.ProvisionLineRouterResponse{}, nil
				}

				return nil, fmt.Errorf("simulated provisioning failure")
			}
			res.LineConfigStorageServer.SetFn = func(ctx context.Context, in *lineconfigstoragepb.UpdateLineConfigurationRequest) (*lineconfigpb.LineConfiguration, error) {
				if tt.saveLineConfigShouldSucceed {
					return &lineconfigpb.LineConfiguration{}, nil
				}

				return nil, fmt.Errorf("simulated write failure")
			}

			ctx := t.Context()

			var buf bytes.Buffer
			runner := &HubServiceCreateRunner{
				org:            "intrinsic",
				hubEndpoint:    tt.hubEndpoint,
				spokeEndpoints: tt.spokeEndpoints,
				forceLocalOnly: tt.forceLocalOnly,
				dialOnpremCluster: func(ctx context.Context, project string, org string, cluster string) (context.Context, *grpc.ClientConn, string, error) {
					conn, err := pubsubtesting.DialTestServer(ctx, res.Listener)
					if err != nil {
						t.Fatalf("Failed to connect to the test server : %v", err)
					}

					// Not closing the connection here because it will be closed by the calling code.
					return ctx, conn, "", nil
				},
				dialCloudCluster: func(ctx context.Context) (*grpc.ClientConn, error) {
					conn, err := pubsubtesting.DialTestServer(ctx, res.Listener)
					if err != nil {
						t.Fatalf("Failed to connect to the test server: %v", err)
					}

					// Not closing the connection here because it will be closed by the calling code.
					return conn, nil
				},
				getServiceVersionToInstall: func(packageName string, serviceName string) (string, error) {
					if packageName == tt.onpremPackageName && serviceName == tt.onpremServiceName {
						return testHubServiceVersion, nil
					}
					return "", fmt.Errorf(
						"unexpected parameters of getServiceVersionToInstall: got(%q, %q), want (%q, %q)",
						packageName, serviceName,
						tt.onpremPackageName, tt.onpremServiceName,
					)
				},
			}

			err := runner.run(ctx, &buf)
			pubsubtesting.VerifyExpectedOutputAndError(t, &buf, err, tt.expectError, tt.expectErrorContains, tt.expectedOutput)
			pubsubtesting.EnsureNoUnexpectedOutput(t, &buf, tt.unexpectedOutput)
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
			expectErrorContains: "neither --cluster nor --hub-endpoint is specified",
			expectedOutput:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			got, err := common.GetHubEndpoint(&buf, tt.cluster, tt.hubEndpoint)
			if got != tt.expectedResult {
				t.Errorf("Unexpected hub endpoint: got %q, want %q", got, tt.expectedResult)
			}
			pubsubtesting.VerifyExpectedOutputAndError(t, &buf, err, tt.expectError, tt.expectErrorContains, tt.expectedOutput)
		})
	}
}
