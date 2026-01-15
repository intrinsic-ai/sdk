// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "intrinsic/logging/proto/bag_packager_service_go_grpc_proto"
)

type mockBagPackagerClient struct {
	pb.BagPackagerClient
	GetBagFunc func(ctx context.Context, in *pb.GetBagRequest, opts ...grpc.CallOption) (*pb.GetBagResponse, error)
}

func (m *mockBagPackagerClient) GetBag(ctx context.Context, in *pb.GetBagRequest, opts ...grpc.CallOption) (*pb.GetBagResponse, error) {
	if m.GetBagFunc != nil {
		return m.GetBagFunc(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "mock GetBag should not be called directly; set GetBagFunc in test case")
}

func TestGetRecordingE(t *testing.T) {
	const (
		testBagID   = "test-bag-id"
		testURL     = "https://fake.url/for/testing"
		testOrg     = "test-org"
		testProject = "test-project"
	)

	tests := []struct {
		name       string
		args       []string
		getBagFunc func(ctx context.Context, in *pb.GetBagRequest, opts ...grpc.CallOption) (*pb.GetBagResponse, error)
		wantErr    string
		wantOut    string
	}{
		{
			name: "Successful get with URL",
			args: []string{"--recording_id", testBagID, "--with_url", "--org", testOrg},
			getBagFunc: func(ctx context.Context, in *pb.GetBagRequest, opts ...grpc.CallOption) (*pb.GetBagResponse, error) {
				url := testURL
				return &pb.GetBagResponse{
					Url: &url,
				}, nil
			},
			wantOut: testURL,
		},
		{
			name:    "Missing recording ID",
			args:    []string{"--org", testOrg},
			wantErr: "required flag(s) \"recording_id\" not set",
		},
		{
			name: "Get bag error",
			args: []string{"--recording_id", testBagID, "--org", testOrg},
			getBagFunc: func(ctx context.Context, in *pb.GetBagRequest, opts ...grpc.CallOption) (*pb.GetBagResponse, error) {
				return nil, errors.New("failed to get bag record")
			},
			wantErr: "recording with id \"test-bag-id\" does not exist",
		},
		{
			name: "URL requested but file not generated",
			args: []string{"--recording_id", testBagID, "--with_url", "--org", testOrg},
			getBagFunc: func(ctx context.Context, in *pb.GetBagRequest, opts ...grpc.CallOption) (*pb.GetBagResponse, error) {
				return nil, errors.New("file does not exist")
			},
			wantErr: "download URL requested for recording with id \"test-bag-id\", but file is not generated yet, please generate it first with `inctl recordings generate`",
		},
	}

	// We disable the org check in tests because the test environment does not have
	// the necessary home directory configuration.
	originalCheckOrgExists := checkOrgExists
	checkOrgExists = false
	t.Cleanup(func() { checkOrgExists = originalCheckOrgExists })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			originalFlagBagID := flagBagID
			t.Cleanup(func() { flagBagID = originalFlagBagID })
			flagBagID = ""

			originalFlagURL := flagURL
			t.Cleanup(func() { flagURL = originalFlagURL })
			flagURL = false

			mockClient := &mockBagPackagerClient{
				GetBagFunc: tc.getBagFunc,
			}
			runner := &GetCmdRunner{
				NewClient: func(cmd *cobra.Command) (pb.BagPackagerClient, error) {
					return mockClient, nil
				},
			}

			var out bytes.Buffer
			rootCmd := &cobra.Command{Use: "inctl"}
			recordingsCmd := &cobra.Command{Use: "recordings"}
			getCmd := NewGetCmd(runner)
			recordingsCmd.AddCommand(getCmd)

			rootCmd.AddCommand(recordingsCmd)
			rootCmd.SetOut(&out)
			rootCmd.SetErr(&out)
			rootCmd.SetArgs(append([]string{"recordings", "get"}, tc.args...))

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
