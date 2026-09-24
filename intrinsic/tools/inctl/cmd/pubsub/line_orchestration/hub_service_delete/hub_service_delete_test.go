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

package hubservicedelete

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"google.golang.org/protobuf/types/known/emptypb"

	servicedeletionutils "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common/service_deletion_utils"

	adgrpcpb "intrinsic/assets/proto/asset_deployment_go_proto"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	aigrpcpb "intrinsic/assets/proto/v1/asset_instances_go_proto"
	provisionerpb "intrinsic/platform/pubsub/cloud_router_provisioner/v1/provisioner_go_proto"
	lineconfigstoragepb "intrinsic/platform/pubsub/connect/cloud/proto/line_configuration_storage/v1/line_configuration_storage_go_proto"
	endpointpb "intrinsic/platform/pubsub/connect/common/proto/line_configuration/v1/endpoint_spec_go_proto"
	lineconfigpb "intrinsic/platform/pubsub/connect/common/proto/line_configuration/v1/line_configuration_go_proto"
	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"
	pubsubtesting "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/testing"
)

func TestDeletionOfEntireNetwork(t *testing.T) {
	tests := []struct {
		name                      string
		onPremServiceName         string
		hubEndpoint               string
		lineConfig                *lineconfigpb.LineConfiguration
		forceLocalOnly            bool
		getLineConfigError        error
		deleteInstalledAssetError error
		ignoreOnpremErrors        bool
		expectedOutput            []string
		unexpectedOutput          []string
		expectFinalError          bool
		expectFinalErrorContains  string
	}{
		{
			name:              "Delete mixed network",
			onPremServiceName: common.OnpremToLineRouterRelayServiceName,
			hubEndpoint:       pubsubtesting.LocalHubEndpoint,
			lineConfig: &lineconfigpb.LineConfiguration{
				HubEndpoint: &endpointpb.EndpointSpec{
					WorkcellName: "node-hub",
					ConnectionSpec: &endpointpb.EndpointSpec_Local{
						Local: &endpointpb.LocalConnectionSpec{},
					},
				},
				SpokeEndpoints: []*endpointpb.EndpointSpec{
					{
						WorkcellName: "vmp-spoke",
						ConnectionSpec: &endpointpb.EndpointSpec_Remote{
							Remote: &endpointpb.RemoteConnectionSpec{},
						},
					},
				},
			},
			expectedOutput: []string{
				"Connecting to the cloud API endpoint",
				"Connecting to node-hub",
				"Uninstalling ai.intrinsic.onprem_to_line_router_relay from node-hub",
				"Connecting to vmp-spoke",
				"Uninstalling ai.intrinsic.onprem_to_line_router_relay from vmp-spoke",
				"Uninstalling cloud router",
				"Deleting line configuration stored in the cloud",
			},
		},
		{
			name:              "Delete local-only network",
			onPremServiceName: common.HubServiceName,
			hubEndpoint:       pubsubtesting.LocalHubEndpoint,
			lineConfig: &lineconfigpb.LineConfiguration{
				HubEndpoint: &endpointpb.EndpointSpec{
					WorkcellName: "node-hub",
					ConnectionSpec: &endpointpb.EndpointSpec_Local{
						Local: &endpointpb.LocalConnectionSpec{},
					},
				},
				SpokeEndpoints: []*endpointpb.EndpointSpec{
					{
						WorkcellName: "node-spoke",
						ConnectionSpec: &endpointpb.EndpointSpec_Local{
							Local: &endpointpb.LocalConnectionSpec{},
						},
					},
				},
			},

			// In a local-only network, the only thing that should be deleted
			// is the relay service running on the hub.
			expectedOutput: []string{
				"Connecting to the cloud API endpoint",
				"Connecting to node-hub",
				"Uninstalling ai.intrinsic.line_orchestration_relay from node-hub",
				"Deleting line configuration stored in the cloud",
			},
		},
		{
			name:                     "Missing configuration prompts user to use force-local-only flag",
			onPremServiceName:        common.HubServiceName,
			hubEndpoint:              pubsubtesting.LocalHubEndpoint,
			getLineConfigError:       status.Errorf(codes.NotFound, "test error"),
			expectFinalError:         true,
			expectFinalErrorContains: "line configuration not found",
			expectedOutput: []string{
				"Connecting to the cloud API endpoint",
				"Configuration for line \"test-org-node-hub\" not found",
				"To delete it, try running this command with the '--force-local-only' flag",
			},
		},
		{
			name: "Failure to read config aborts deletion",

			// Any error other than "not found" is a failure.
			getLineConfigError: fmt.Errorf("test error"),
			hubEndpoint:        pubsubtesting.LocalHubEndpoint,
			expectedOutput: []string{
				"Connecting to the cloud API endpoint",
			},
			expectFinalError:         true,
			expectFinalErrorContains: "test error",
		},
		{
			name:                     "Unimplemented error prompts user to delete local-only network",
			hubEndpoint:              pubsubtesting.LocalHubEndpoint,
			getLineConfigError:       status.Errorf(codes.Unimplemented, ""),
			expectFinalError:         true,
			expectFinalErrorContains: "line configuration storage is not implemented",
			expectedOutput: []string{
				fmt.Sprintf("Try running the command with the '--%v' flag", common.KeyForceLocalOnly),
			},
		},
		{
			name:              "force-local-only flag disables interactions with cloud services",
			onPremServiceName: common.OnpremToLineRouterRelayServiceName,
			hubEndpoint:       pubsubtesting.LocalHubEndpoint,
			forceLocalOnly:    true,
			expectedOutput: []string{
				"Deleting the line orchestration network assuming that it is local-only",
				"Connecting to node-hub",
				"Uninstalling ai.intrinsic.line_orchestration_relay from node-hub",
			},
			unexpectedOutput: []string{
				"Connecting to the cloud API endpoint",
				"Connecting to vmp-spoke",
			},
		},
		{
			name:                     "force-local-only flag is not compatible with remote endpoints",
			onPremServiceName:        common.OnpremToLineRouterRelayServiceName,
			hubEndpoint:              pubsubtesting.RemoteHubEndpoint,
			forceLocalOnly:           true,
			expectFinalError:         true,
			expectFinalErrorContains: fmt.Sprintf("The '--%v' flag is present, but the hub endpoint is remote.", common.KeyForceLocalOnly),
		},
		{
			name:                     "Invalid endpoint spec is rejected",
			hubEndpoint:              "node-###@local",
			expectFinalError:         true,
			expectFinalErrorContains: "contains invalid characters",
			expectedOutput: []string{
				"Cannot create a valid line id",
			},
		},
		{
			name:              "Failure to delete onprem asset aborts deletion",
			onPremServiceName: common.OnpremToLineRouterRelayServiceName,
			hubEndpoint:       pubsubtesting.LocalHubEndpoint,
			lineConfig: &lineconfigpb.LineConfiguration{
				HubEndpoint: &endpointpb.EndpointSpec{
					WorkcellName: "node-hub",
					ConnectionSpec: &endpointpb.EndpointSpec_Local{
						Local: &endpointpb.LocalConnectionSpec{},
					},
				},
				SpokeEndpoints: []*endpointpb.EndpointSpec{
					{
						WorkcellName: "vmp-spoke",
						ConnectionSpec: &endpointpb.EndpointSpec_Remote{
							Remote: &endpointpb.RemoteConnectionSpec{},
						},
					},
				},
			},
			deleteInstalledAssetError: fmt.Errorf("test onprem error"),
			expectFinalError:          true,
			expectFinalErrorContains:  "test onprem error",
		},
		{
			name:              "Failure to delete onprem asset ignored",
			onPremServiceName: common.OnpremToLineRouterRelayServiceName,
			hubEndpoint:       pubsubtesting.LocalHubEndpoint,
			lineConfig: &lineconfigpb.LineConfiguration{
				HubEndpoint: &endpointpb.EndpointSpec{
					WorkcellName: "node-hub",
					ConnectionSpec: &endpointpb.EndpointSpec_Local{
						Local: &endpointpb.LocalConnectionSpec{},
					},
				},
				SpokeEndpoints: []*endpointpb.EndpointSpec{
					{
						WorkcellName: "vmp-spoke",
						ConnectionSpec: &endpointpb.EndpointSpec_Remote{
							Remote: &endpointpb.RemoteConnectionSpec{},
						},
					},
				},
			},
			deleteInstalledAssetError: fmt.Errorf("test onprem error"),
			ignoreOnpremErrors:        true,
			expectFinalError:          false,
			expectedOutput: []string{
				"Failed to delete \"onprem_to_line_router_relay\" from \"node-hub\": failed to uninstall onprem_to_line_router_relay service asset: rpc error: code = Unknown desc = test onprem error, but the '--ignore-onprem-errors' flag is present. Moving on",
				"Failed to delete \"onprem_to_line_router_relay\" from \"vmp-spoke\": failed to uninstall onprem_to_line_router_relay service asset: rpc error: code = Unknown desc = test onprem error, but the '--ignore-onprem-errors' flag is present. Moving on",
				"Uninstalling cloud router",
				"Deleting line configuration stored in the cloud",
				"Line orchestration network with the hub at \"node-hub\" (line id \"test-org-node-hub\") has been deleted",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := pubsubtesting.SetupTestServer(t)
			res.InstServer.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
				return &aigrpcpb.ListAssetInstancesResponse{
					AssetInstances: []*aigrpcpb.AssetInstance{
						{Name: tt.onPremServiceName},
					},
				}, nil
			}

			res.DepServer.DeleteResourceFn = func(ctx context.Context, in *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
				return &lropb.Operation{Done: true, Name: "op1"}, nil
			}

			res.OpServer.GetOperationFn = func(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
				return &lropb.Operation{Done: true, Name: "op1"}, nil
			}

			res.IaServer.DeleteInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.DeleteInstalledAssetRequest) (*lropb.Operation, error) {
				if tt.deleteInstalledAssetError != nil {
					return nil, tt.deleteInstalledAssetError
				}
				return &lropb.Operation{Done: true, Name: "op1"}, nil
			}

			res.LineConfigStorageServer.GetFn = func(ctx context.Context, in *lineconfigstoragepb.GetLineConfigurationRequest) (*lineconfigpb.LineConfiguration, error) {
				if tt.getLineConfigError != nil {
					return nil, tt.getLineConfigError
				}

				return tt.lineConfig, nil
			}

			res.LineConfigStorageServer.DeleteFn = func(ctx context.Context, in *lineconfigstoragepb.DeleteLineConfigurationRequest) (*emptypb.Empty, error) {
				return &emptypb.Empty{}, nil
			}

			res.ProvisionerServer.DeleteFn = func(ctx context.Context, req *provisionerpb.DeleteLineRouterRequest) (*provisionerpb.LineRouterDeleteResponse, error) {
				return &provisionerpb.LineRouterDeleteResponse{}, nil
			}

			ctx := t.Context()
			var buf bytes.Buffer
			runner := &HubServiceDeleteRunner{
				ServiceDeleter: servicedeletionutils.ServiceDeleter{
					ProjectID:                "test-project",
					OrgID:                    "test-org",
					IgnoreOnpremErrors:       tt.ignoreOnpremErrors,
					ShouldRetainServiceAsset: false,
					DialOnpremCluster: func(ctx context.Context, project string, org string, cluster string) (context.Context, *grpc.ClientConn, string, error) {
						conn, err := pubsubtesting.DialTestServer(ctx, res.Listener)
						if err != nil {
							t.Fatalf("Failed to connect to the test server: %v", err)
						}
						// Not closing the connection here because it will be closed by the calling code.
						return ctx, conn, "", err
					},
				},
				HubEndpoint:    tt.hubEndpoint,
				ForceLocalOnly: tt.forceLocalOnly,

				DialCloudCluster: func(ctx context.Context) (*grpc.ClientConn, error) {
					conn, err := pubsubtesting.DialTestServer(ctx, res.Listener)
					if err != nil {
						t.Fatalf("Failed to connect to the test server: %v", err)
					}
					// Not closing the connection here because it will be closed by the calling code.
					return conn, err
				},
			}

			err := runner.Run(ctx, &buf)
			pubsubtesting.VerifyExpectedOutputAndError(t, &buf, err, tt.expectFinalError, tt.expectFinalErrorContains, tt.expectedOutput)
			pubsubtesting.EnsureNoUnexpectedOutput(t, &buf, tt.unexpectedOutput)
		})
	}
}
