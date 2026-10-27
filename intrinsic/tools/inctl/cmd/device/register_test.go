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

package device

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/anypb"

	clustermanagergrpcpb "intrinsic/frontend/cloud/api/v1/clustermanager_api_go_proto"
	clustermanagerpb "intrinsic/frontend/cloud/api/v1/clustermanager_api_go_proto"
)

type fakeClustersServiceClient struct {
	clustermanagergrpcpb.ClustersServiceClient
	gotReq         *clustermanagerpb.CreateClusterRequest
	retClusterName string
	retErr         error
}

func (f *fakeClustersServiceClient) CreateCluster(_ context.Context, req *clustermanagerpb.CreateClusterRequest, _ ...grpc.CallOption) (*lropb.Operation, error) {
	f.gotReq = req
	if f.retErr != nil {
		return nil, f.retErr
	}
	clusterName := f.retClusterName
	if clusterName == "" {
		clusterName = req.GetClusterInfo().GetClusterId()
	}
	resp, err := anypb.New(&clustermanagerpb.Cluster{ClusterName: clusterName})
	if err != nil {
		return nil, err
	}
	return &lropb.Operation{
		Name:   "createCluster:" + clusterName + ":" + req.GetClusterInfo().GetPending().GetDeviceId(),
		Result: &lropb.Operation_Response{Response: resp},
	}, nil
}

func TestValidHostname(t *testing.T) {
	testCases := []struct {
		name     string
		hostname string
		valid    bool
		index    int
	}{
		{
			name:     "empty",
			hostname: "",
			valid:    false,
			index:    0,
		},
		{
			name:     "single-char",
			hostname: "h",
			valid:    true,
		},
		{
			name:     "simple",
			hostname: "host",
			valid:    true,
		},
		{
			name:     "complex",
			hostname: "host-13245-123413240-998",
			valid:    true,
		},
		{
			name:     "generated",
			hostname: "node-7da4bec3-3aa2-4859-97dd-b07dcce37c19",
			valid:    true,
		},
		{
			name:     "with-alpha",
			hostname: "host123",
			valid:    true,
		},
		{
			name:     "with-dash",
			hostname: "host-123",
			valid:    true,
		},
		{
			name:     "no-starting-dash",
			hostname: "-host-123",
			valid:    false,
			index:    0,
		},
		{
			name:     "no-ending-dash",
			hostname: "host-123-",
			valid:    false,
			index:    9,
		},
		{
			name:     "no-capital",
			hostname: "HOST-123",
			valid:    false,
		},
		{
			name:     "no-underscore",
			hostname: "host_123",
			valid:    false,
			index:    5,
		},
		{
			name:     "no-overlength",
			hostname: "host-123456790-123456790-123456790-123456790-123456790-123456790",
			valid:    false,
			index:    43,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			indexGot, validGot := validHostname(tc.hostname)
			if got, want := validGot, tc.valid; got != want {
				t.Errorf("validHostname(%q) = %t, want %t", tc.hostname, got, want)
			}

			if got, want := indexGot, tc.index; got != want {
				t.Errorf("validHostname(%q).index = %v, want %v", tc.hostname, got, want)
			}
		})
	}
}

func TestRegisterControlPlane(t *testing.T) {
	testCases := []struct {
		name           string
		deviceID       string
		hostname       string
		region         string
		private        bool
		replace        bool
		noUpdate       bool
		createdByTest  string
		retClusterName string
		wantReq        *clustermanagerpb.CreateClusterRequest
		wantDeviceID   string
		wantHostname   string
	}{
		{
			name:     "defaults",
			deviceID: "node-1234",
			hostname: "node-1234",
			region:   "unspecified",
			wantReq: &clustermanagerpb.CreateClusterRequest{
				ClusterInfo: &clustermanagerpb.CreateClusterInfo{
					Pending: &clustermanagerpb.PendingDeviceDescription{
						DeviceId: "node-1234",
					},
					ClusterId:   "node-1234",
					DisplayName: "node-1234",
					Region:      "unspecified",
					UpdateMode:  clustermanagerpb.PlatformUpdateMode_PLATFORM_UPDATE_MODE_AUTOMATIC,
				},
			},
			wantDeviceID: "node-1234",
			wantHostname: "node-1234",
		},
		{
			name:           "registration-code-resolves-full-id",
			deviceID:       "0123456789ab",
			hostname:       "0123456789ab",
			region:         "unspecified",
			retClusterName: "node-00000000-0000-0000-0000-0123456789ab",
			wantReq: &clustermanagerpb.CreateClusterRequest{
				ClusterInfo: &clustermanagerpb.CreateClusterInfo{
					Pending: &clustermanagerpb.PendingDeviceDescription{
						DeviceId: "0123456789ab",
					},
					ClusterId:   "0123456789ab",
					DisplayName: "0123456789ab",
					Region:      "unspecified",
					UpdateMode:  clustermanagerpb.PlatformUpdateMode_PLATFORM_UPDATE_MODE_AUTOMATIC,
				},
			},
			wantDeviceID: "node-00000000-0000-0000-0000-0123456789ab",
			wantHostname: "node-00000000-0000-0000-0000-0123456789ab",
		},
		{
			name:          "custom-options",
			deviceID:      "0123456789ab",
			hostname:      "my-cluster",
			region:        "us-west1",
			private:       true,
			replace:       true,
			noUpdate:      true,
			createdByTest: "e2e-test-id",
			wantReq: &clustermanagerpb.CreateClusterRequest{
				ClusterInfo: &clustermanagerpb.CreateClusterInfo{
					Pending: &clustermanagerpb.PendingDeviceDescription{
						DeviceId: "0123456789ab",
					},
					ClusterId:     "my-cluster",
					DisplayName:   "my-cluster",
					Region:        "us-west1",
					PrivateDevice: true,
					Replace:       true,
					UpdateMode:    clustermanagerpb.PlatformUpdateMode_PLATFORM_UPDATE_MODE_ON,
				},
				CreatedByTest: "e2e-test-id",
			},
			wantDeviceID: "0123456789ab",
			wantHostname: "my-cluster",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deviceID = tc.deviceID
			hostname = tc.hostname
			deviceRegion = tc.region
			privateDevice = tc.private
			replaceDevice = tc.replace
			noUpdate = tc.noUpdate
			t.Setenv("INCTL_CREATED_BY_TEST", tc.createdByTest)

			fakeClient := &fakeClustersServiceClient{retClusterName: tc.retClusterName}
			client := authedClient{grpcClient: fakeClient}

			if err := registerControlPlane(t.Context(), client, "test-project", "test-org"); err != nil {
				t.Fatalf("registerControlPlane() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantReq, fakeClient.gotReq, protocmp.Transform()); diff != "" {
				t.Errorf("CreateClusterRequest diff (-want +got):\n%s", diff)
			}
			if deviceID != tc.wantDeviceID {
				t.Errorf("deviceID = %q, want %q", deviceID, tc.wantDeviceID)
			}
			if hostname != tc.wantHostname {
				t.Errorf("hostname = %q, want %q", hostname, tc.wantHostname)
			}
		})
	}
}

func TestRegisterControlPlaneErrors(t *testing.T) {
	deviceID = "node-1234"
	hostname = "my-cluster"

	testCases := []struct {
		name        string
		grpcErr     error
		wantErrText string
	}{
		{
			name:        "already-exists",
			grpcErr:     status.Error(codes.AlreadyExists, "cluster already exists"),
			wantErrText: `cluster "my-cluster" already exists`,
		},
		{
			name:        "not-found",
			grpcErr:     status.Error(codes.NotFound, "device not found"),
			wantErrText: `device "node-1234" does not exist`,
		},
		{
			name:        "unauthenticated",
			grpcErr:     status.Error(codes.Unauthenticated, "invalid token"),
			wantErrText: "your login key has expired or been replaced",
		},
		{
			name:        "permission-denied",
			grpcErr:     status.Error(codes.PermissionDenied, "missing cluster_add"),
			wantErrText: `you do not have the necessary permissions to add a cluster on organization "test-org@test-project"`,
		},
		{
			name:        "internal-error",
			grpcErr:     status.Error(codes.Internal, "internal failure"),
			wantErrText: "request failed:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeClient := &fakeClustersServiceClient{retErr: tc.grpcErr}
			client := authedClient{grpcClient: fakeClient}

			err := registerControlPlane(t.Context(), client, "test-project", "test-org")
			if err == nil {
				t.Fatalf("registerControlPlane() succeeded, want error containing %q", tc.wantErrText)
			}
			if !strings.Contains(err.Error(), tc.wantErrText) {
				t.Errorf("registerControlPlane() error = %q, want substring %q", err.Error(), tc.wantErrText)
			}
		})
	}
}

func TestRegisterCmdControlPlaneClusterNameValidation(t *testing.T) {
	deviceRole = "control-plane"
	deviceID = "node-1234"
	hostname = "my-cluster"
	clusterName = "other-cluster"
	t.Cleanup(func() {
		clusterName = ""
	})

	err := registerCmd.RunE(registerCmd, nil)
	if err == nil {
		t.Fatal("registerCmd.RunE() succeeded, want error when --cluster_name differs from --hostname on control-plane")
	}
	want := `--cluster_name must be empty or equal to --hostname for role "control-plane"`
	if !strings.Contains(err.Error(), want) {
		t.Errorf("registerCmd.RunE() error = %q, want substring %q", err.Error(), want)
	}
}

func TestWaitForStatusAvailableCanceledNoPanic(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	client := authedClient{
		client:  &http.Client{},
		baseURL: url.URL{Scheme: "https", Host: "example.com"},
	}
	err := waitForStatusAvailable(ctx, client, "my-cluster", "my-cluster")
	if !errors.Is(err, context.Canceled) {
		t.Errorf("waitForStatusAvailable() error = %v, want %v", err, context.Canceled)
	}
}
