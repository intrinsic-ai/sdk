// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	leasegrpcpb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_proto"
	bagmetadatapb "intrinsic/logging/proto/bag_metadata_go_proto"
	bagpackagerpb "intrinsic/logging/proto/bag_packager_service_go_proto"
	replaygrpcpb "intrinsic/logging/proto/replay_service_go_proto"
	"intrinsic/tools/inctl/cmd/root"
)

type mockLeaseClient struct {
	leasegrpcpb.VMPoolLeaseServiceClient
	LeaseFunc func(ctx context.Context, in *leasegrpcpb.LeaseRequest, opts ...grpc.CallOption) (*leasegrpcpb.LeaseResponse, error)
}

func (m *mockLeaseClient) Lease(ctx context.Context, in *leasegrpcpb.LeaseRequest, opts ...grpc.CallOption) (*leasegrpcpb.LeaseResponse, error) {
	if m.LeaseFunc != nil {
		return m.LeaseFunc(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "LeaseFunc not set")
}

type mockReplayClient struct {
	replaygrpcpb.ReplayClient
	VisualizeRecordingFunc func(ctx context.Context, in *replaygrpcpb.VisualizeRecordingRequest, opts ...grpc.CallOption) (*replaygrpcpb.VisualizeRecordingResponse, error)
}

func (m *mockReplayClient) VisualizeRecording(ctx context.Context, in *replaygrpcpb.VisualizeRecordingRequest, opts ...grpc.CallOption) (*replaygrpcpb.VisualizeRecordingResponse, error) {
	if m.VisualizeRecordingFunc != nil {
		return m.VisualizeRecordingFunc(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "VisualizeRecordingFunc not set")
}

type visualizeMockBagPackagerClient struct {
	bagpackagerpb.BagPackagerClient
	GetBagFunc func(ctx context.Context, in *bagpackagerpb.GetBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GetBagResponse, error)
}

func (m *visualizeMockBagPackagerClient) GetBag(ctx context.Context, in *bagpackagerpb.GetBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GetBagResponse, error) {
	if m.GetBagFunc != nil {
		return m.GetBagFunc(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "GetBagFunc not set")
}

func TestVisualizeRecordingE(t *testing.T) {
	const (
		testRecordingID = "test-recording-id"
		testDuration    = "1h"
		testURL         = "https://fake.url/for/testing"
		testOrg         = "test-org"
	)

	tests := []struct {
		name                   string
		args                   []string
		leaseFunc              func(ctx context.Context, in *leasegrpcpb.LeaseRequest, opts ...grpc.CallOption) (*leasegrpcpb.LeaseResponse, error)
		visualizeRecordingFunc func(ctx context.Context, in *replaygrpcpb.VisualizeRecordingRequest, opts ...grpc.CallOption) (*replaygrpcpb.VisualizeRecordingResponse, error)
		getBagFunc             func(ctx context.Context, in *bagpackagerpb.GetBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GetBagResponse, error)
		wantErr                string
		wantOut                string
	}{
		{
			name: "creates visualization host and generates URL on success",
			args: []string{"--recording_id", testRecordingID, "--duration", testDuration, "--org", testOrg},
			leaseFunc: func(ctx context.Context, in *leasegrpcpb.LeaseRequest, opts ...grpc.CallOption) (*leasegrpcpb.LeaseResponse, error) {
				return &leasegrpcpb.LeaseResponse{
					Lease: &leasegrpcpb.Lease{
						Instance: "test-instance",
						Expires:  timestamppb.New(time.Now().Add(1 * time.Hour)),
					},
				}, nil
			},
			visualizeRecordingFunc: func(ctx context.Context, in *replaygrpcpb.VisualizeRecordingRequest, opts ...grpc.CallOption) (*replaygrpcpb.VisualizeRecordingResponse, error) {
				return &replaygrpcpb.VisualizeRecordingResponse{
					Url: testURL,
				}, nil
			},
			getBagFunc: func(ctx context.Context, in *bagpackagerpb.GetBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GetBagResponse, error) {
				return &bagpackagerpb.GetBagResponse{
					Bag: &bagpackagerpb.BagRecord{
						BagMetadata: &bagmetadatapb.BagMetadata{
							EventSources: []*bagmetadatapb.EventSourceMetadata{
								{EventSourceWithTypeHints: &bagmetadatapb.EventSourceWithTypeHints{EventSource: "executive.operation_state"}},
							},
						},
					},
				}, nil
			},
			wantOut: testURL,
		},
		{
			name:    "returns error when recording_id flag is missing",
			args:    []string{"--duration", testDuration, "--org", testOrg},
			wantErr: `required flag(s) "recording_id" not set`,
		},
		{
			name:    "returns error when duration flag is missing",
			args:    []string{"--recording_id", testRecordingID, "--org", testOrg},
			wantErr: `required flag(s) "duration" not set`,
		},
		{
			name: "returns error when lease creation fails",
			args: []string{"--recording_id", testRecordingID, "--duration", testDuration, "--org", testOrg},
			leaseFunc: func(ctx context.Context, in *leasegrpcpb.LeaseRequest, opts ...grpc.CallOption) (*leasegrpcpb.LeaseResponse, error) {
				return nil, status.Error(codes.PermissionDenied, "permission denied")
			},
			wantErr: "visualization host create request failed",
		},
		{
			name: "returns error when visualization creation fails",
			args: []string{"--recording_id", testRecordingID, "--duration", testDuration, "--org", testOrg},
			leaseFunc: func(ctx context.Context, in *leasegrpcpb.LeaseRequest, opts ...grpc.CallOption) (*leasegrpcpb.LeaseResponse, error) {
				return &leasegrpcpb.LeaseResponse{
					Lease: &leasegrpcpb.Lease{
						Instance: "test-instance",
						Expires:  timestamppb.New(time.Now().Add(1 * time.Hour)),
					},
				}, nil
			},
			visualizeRecordingFunc: func(ctx context.Context, in *replaygrpcpb.VisualizeRecordingRequest, opts ...grpc.CallOption) (*replaygrpcpb.VisualizeRecordingResponse, error) {
				return nil, status.Error(codes.AlreadyExists, "already exists")
			},
			getBagFunc: func(ctx context.Context, in *bagpackagerpb.GetBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GetBagResponse, error) {
				return &bagpackagerpb.GetBagResponse{
					Bag: &bagpackagerpb.BagRecord{
						BagMetadata: &bagmetadatapb.BagMetadata{},
					},
				}, nil
			},
			wantErr: "already exists",
		},
		{
			name: "constructs precise denylist that excludes debug data but retains skill data overlap",
			args: []string{"--recording_id", testRecordingID, "--duration", testDuration, "--org", testOrg, "--exclude_debug_data"},
			leaseFunc: func(ctx context.Context, in *leasegrpcpb.LeaseRequest, opts ...grpc.CallOption) (*leasegrpcpb.LeaseResponse, error) {
				return &leasegrpcpb.LeaseResponse{
					Lease: &leasegrpcpb.Lease{
						Instance: "test-instance",
						Expires:  timestamppb.New(time.Now().Add(1 * time.Hour)),
					},
				}, nil
			},
			visualizeRecordingFunc: func(ctx context.Context, in *replaygrpcpb.VisualizeRecordingRequest, opts ...grpc.CallOption) (*replaygrpcpb.VisualizeRecordingResponse, error) {
				// Verify denylist includes operation_state but NOT extended_status.
				denylist := in.GetVisualizationOptions().GetDefaultVisualizerFilters().GetEventSources().GetDenylistRegexes()
				assert.Contains(t, denylist, "^executive\\.operation_state$")
				assert.NotContains(t, denylist, "^executive\\.extended_status$")
				return &replaygrpcpb.VisualizeRecordingResponse{
					Url: testURL,
				}, nil
			},
			getBagFunc: func(ctx context.Context, in *bagpackagerpb.GetBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GetBagResponse, error) {
				return &bagpackagerpb.GetBagResponse{
					Bag: &bagpackagerpb.BagRecord{
						BagMetadata: &bagmetadatapb.BagMetadata{
							EventSources: []*bagmetadatapb.EventSourceMetadata{
								{EventSourceWithTypeHints: &bagmetadatapb.EventSourceWithTypeHints{EventSource: "executive.operation_state"}},
								{EventSourceWithTypeHints: &bagmetadatapb.EventSourceWithTypeHints{EventSource: "executive.extended_status"}},
							},
						},
					},
				}, nil
			},
			wantOut: testURL,
		},
		{
			name: "JSON output returns correctly structured data without prompts",
			args: []string{"--recording_id", testRecordingID, "--duration", testDuration, "--org", testOrg},
			leaseFunc: func(ctx context.Context, in *leasegrpcpb.LeaseRequest, opts ...grpc.CallOption) (*leasegrpcpb.LeaseResponse, error) {
				return &leasegrpcpb.LeaseResponse{
					Lease: &leasegrpcpb.Lease{
						Instance: "test-instance",
						Expires:  timestamppb.New(time.Now().Add(1 * time.Hour)),
					},
				}, nil
			},
			visualizeRecordingFunc: func(ctx context.Context, in *replaygrpcpb.VisualizeRecordingRequest, opts ...grpc.CallOption) (*replaygrpcpb.VisualizeRecordingResponse, error) {
				return &replaygrpcpb.VisualizeRecordingResponse{
					Url: testURL,
				}, nil
			},
			getBagFunc: func(ctx context.Context, in *bagpackagerpb.GetBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GetBagResponse, error) {
				return &bagpackagerpb.GetBagResponse{
					Bag: &bagpackagerpb.BagRecord{
						BagMetadata: &bagmetadatapb.BagMetadata{},
					},
				}, nil
			},
			wantOut: `"status": "success"`,
		},
	}

	// We disable the org check in tests because the test environment does not have
	// the necessary home directory configuration.
	originalCheckOrgExists := checkOrgExists
	checkOrgExists = false
	t.Cleanup(func() { checkOrgExists = originalCheckOrgExists })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			originalFlagRecordingID := flagRecordingID
			t.Cleanup(func() { flagRecordingID = originalFlagRecordingID })
			flagRecordingID = ""

			originalFlagDuration := flagDuration
			t.Cleanup(func() { flagDuration = originalFlagDuration })
			flagDuration = ""

			originalFlagOutput := root.FlagOutput
			t.Cleanup(func() { root.FlagOutput = originalFlagOutput })
			if strings.Contains(tc.name, "JSON") {
				root.FlagOutput = "json"
			} else {
				root.FlagOutput = ""
			}

			// Provide a default getBagFunc if the test case didn't define one
			// to avoid panics on tests that don't care about the bag fetch logic.
			getBagFunc := tc.getBagFunc
			if getBagFunc == nil {
				getBagFunc = func(ctx context.Context, in *bagpackagerpb.GetBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GetBagResponse, error) {
					return &bagpackagerpb.GetBagResponse{
						Bag: &bagpackagerpb.BagRecord{
							BagMetadata: &bagmetadatapb.BagMetadata{},
						},
					}, nil
				}
			}

			runner := &VisualizeCmdRunner{
				NewLeaseClient: func(cmd *cobra.Command) (leasegrpcpb.VMPoolLeaseServiceClient, error) {
					return &mockLeaseClient{LeaseFunc: tc.leaseFunc}, nil
				},
				ReplayClientFactory: func(ctx context.Context, v *viper.Viper, clusterName string) (replaygrpcpb.ReplayClient, io.Closer, error) {
					return &mockReplayClient{VisualizeRecordingFunc: tc.visualizeRecordingFunc}, nil, nil
				},
				NewBagPackagerClient: func(cmd *cobra.Command) (bagpackagerpb.BagPackagerClient, error) {
					return &visualizeMockBagPackagerClient{GetBagFunc: getBagFunc}, nil
				},
			}

			var out bytes.Buffer
			rootCmd := &cobra.Command{Use: "inctl"}
			recordingsCmd := &cobra.Command{Use: "recordings"}
			visualizeCmd := &cobra.Command{Use: "visualize"}
			createCmd := NewVisualizeCmd(runner)

			visualizeCmd.AddCommand(createCmd)
			recordingsCmd.AddCommand(visualizeCmd)
			rootCmd.AddCommand(recordingsCmd)

			rootCmd.SetOut(&out)
			rootCmd.SetErr(&out)
			rootCmd.SetArgs(append([]string{"recordings", "visualize", "create"}, tc.args...))

			err := rootCmd.Execute()

			if tc.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			} else {
				assert.NoError(t, err)

				if strings.Contains(tc.name, "JSON") {
					var parsed map[string]interface{}
					assert.NoError(t, json.Unmarshal([]byte(out.String()), &parsed), "Output should be valid JSON")
					assert.Equal(t, "success", parsed["status"])
					data, ok := parsed["data"].(map[string]interface{})
					assert.True(t, ok, "Data should be a map")
					assert.Equal(t, testURL, data["url"])
				} else {
					assert.Contains(t, out.String(), tc.wantOut)
				}
			}
		})
	}
}
