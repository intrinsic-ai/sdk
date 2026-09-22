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

package stopforwarding

import (
	"bytes"
	"context"
	"testing"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"

	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"
	servicedeletionutils "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/service_deletion_utils"
	pubsubtesting "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/testing"

	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	adgrpcpb "intrinsic/assets/proto/asset_deployment_go_proto"
	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	aigrpcpb "intrinsic/assets/proto/v1/asset_instances_go_proto"
)

func TestStopForwarding(t *testing.T) {
	tests := []struct {
		name               string
		forwarderInstalled bool
		shouldUninstall    bool
		expectedOutput     []string
		unexpectedOutput   []string
	}{
		{
			name:               "Successful stop and uninstallation",
			forwarderInstalled: true,
			shouldUninstall:    true,
			expectedOutput: []string{
				"Successfully deleted an instance of the line_orchestration_forwarder service",
				"Successfully uninstalled the line_orchestration_forwarder service asset",
			},
		},
		{
			name:               "Successful stop without uninstallation",
			forwarderInstalled: true,
			shouldUninstall:    false,
			expectedOutput: []string{
				"Successfully deleted an instance of the line_orchestration_forwarder service",
			},
			unexpectedOutput: []string{
				"Successfully uninstalled the line_orchestration_forwarder service asset",
			},
		},
		{
			name:               "No-op when no forwarding",
			forwarderInstalled: false,
			shouldUninstall:    true,
			expectedOutput: []string{
				"Deleting existing instances of the line_orchestration_forwarder service in the current solution",
				"0 instances have been deleted",
				"line_orchestration_forwarder service asset is not installed, nothing else to do",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := pubsubtesting.SetupTestServer(t)
			res.InstServer.ListAssetInstancesFn = func(ctx context.Context, in *aigrpcpb.ListAssetInstancesRequest) (*aigrpcpb.ListAssetInstancesResponse, error) {
				if !tt.forwarderInstalled {
					return &aigrpcpb.ListAssetInstancesResponse{}, nil
				}
				return &aigrpcpb.ListAssetInstancesResponse{
					AssetInstances: []*aigrpcpb.AssetInstance{
						{Name: common.ForwardingServiceName},
					},
				}, nil
			}

			if !tt.forwarderInstalled {
				res.IaServer.GetInstalledAssetFn = func(ctx context.Context, req *iagrpcpb.GetInstalledAssetRequest) (*iagrpcpb.InstalledAsset, error) {
					return nil, grpcstatus.Errorf(codes.NotFound, "not found")
				}
			}

			res.DepServer.DeleteResourceFn = func(ctx context.Context, in *adgrpcpb.DeleteResourceRequest) (*lropb.Operation, error) {
				return &lropb.Operation{Done: true, Name: "op1"}, nil
			}

			res.OpServer.GetOperationFn = func(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
				return &lropb.Operation{Done: true, Name: "op1"}, nil
			}

			res.IaServer.DeleteInstalledAssetFn = func(ctx context.Context, in *iagrpcpb.DeleteInstalledAssetRequest) (*lropb.Operation, error) {
				return &lropb.Operation{Done: true, Name: "op1"}, nil
			}

			ctx := t.Context()
			conn, err := pubsubtesting.DialTestServer(ctx, res.Listener)
			if err != nil {
				t.Fatalf("Failed to connect to the test server: %v", err)
			}
			var buf bytes.Buffer
			runner := &servicedeletionutils.ServiceDeleteCmdRunner{
				CmdRunnerBase: *common.NewCmdRunnerBase(
					conn,
					&buf,
					"test-cluster",
					common.ForwardingServicePackage,
					common.ForwardingServiceName),
				ShouldUninstallServiceAsset: tt.shouldUninstall,
			}

			err = runner.Run(ctx)
			if err != nil {
				t.Fatalf("%v", err)
			}
			pubsubtesting.VerifyExpectedOutputAndError(
				t,
				&buf,
				err,
				false, /* expectFinalError */
				"",    /* expectFinalErrorContains */
				tt.expectedOutput)
			pubsubtesting.EnsureNoUnexpectedOutput(t, &buf, tt.unexpectedOutput)
		})
	}
}
