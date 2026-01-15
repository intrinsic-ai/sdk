// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	leasegrpcpb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_grpc_proto"
	replaygrpcpb "intrinsic/logging/proto/replay_service_go_grpc_proto"
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
		wantErr                string
		wantOut                string
	}{
		{
			name: "Successful visualization",
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
			wantOut: testURL,
		},
		{
			name:    "Missing recording ID",
			args:    []string{"--duration", testDuration, "--org", testOrg},
			wantErr: `required flag(s) "recording_id" not set`,
		},
		{
			name:    "Missing duration",
			args:    []string{"--recording_id", testRecordingID, "--org", testOrg},
			wantErr: `required flag(s) "duration" not set`,
		},
		{
			name: "Lease error",
			args: []string{"--recording_id", testRecordingID, "--duration", testDuration, "--org", testOrg},
			leaseFunc: func(ctx context.Context, in *leasegrpcpb.LeaseRequest, opts ...grpc.CallOption) (*leasegrpcpb.LeaseResponse, error) {
				return nil, status.Error(codes.PermissionDenied, "permission denied")
			},
			wantErr: "visualization host create request failed",
		},
		{
			name: "Visualize error",
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
			wantErr: "already exists",
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

			runner := &VisualizeCmdRunner{
				NewLeaseClient: func(cmd *cobra.Command) (leasegrpcpb.VMPoolLeaseServiceClient, error) {
					return &mockLeaseClient{LeaseFunc: tc.leaseFunc}, nil
				},
				ReplayClientFactory: func(ctx context.Context, v *viper.Viper, clusterName string) (replaygrpcpb.ReplayClient, io.Closer, error) {
					return &mockReplayClient{VisualizeRecordingFunc: tc.visualizeRecordingFunc}, nil, nil
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
				assert.Contains(t, out.String(), tc.wantOut)
			}
		})
	}
}
