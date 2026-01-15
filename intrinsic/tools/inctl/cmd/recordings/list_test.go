// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	bmpb "intrinsic/logging/proto/bag_metadata_go_proto"
	pb "intrinsic/logging/proto/bag_packager_service_go_grpc_proto"
)

type mockBagPackagerClientForList struct {
	pb.BagPackagerClient
	ListBagsFunc func(ctx context.Context, in *pb.ListBagsRequest, opts ...grpc.CallOption) (*pb.ListBagsResponse, error)
}

func (m *mockBagPackagerClientForList) ListBags(ctx context.Context, in *pb.ListBagsRequest, opts ...grpc.CallOption) (*pb.ListBagsResponse, error) {
	if m.ListBagsFunc != nil {
		return m.ListBagsFunc(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "mock ListBags should not be called directly; set ListBagsFunc in test case")
}

func TestListRecordingsE(t *testing.T) {
	const (
		testOrg        = "test-org"
		testWorkcell   = "test-workcell"
		testRecording1 = "recording-1"
		testRecording2 = "recording-2"
	)

	tests := []struct {
		name         string
		args         []string
		listBagsFunc func(ctx context.Context, in *pb.ListBagsRequest, opts ...grpc.CallOption) (*pb.ListBagsResponse, error)
		promptResp   string
		wantErr      string
		wantOut      string
	}{
		{
			name: "Successful list",
			args: []string{"--workcell", testWorkcell, "--org", testOrg},
			listBagsFunc: func(ctx context.Context, in *pb.ListBagsRequest, opts ...grpc.CallOption) (*pb.ListBagsResponse, error) {
				return &pb.ListBagsResponse{
					Bags: []*pb.BagRecord{
						{BagMetadata: &bmpb.BagMetadata{BagId: testRecording1, StartTime: timestamppb.Now(), EndTime: timestamppb.Now()}},
						{BagMetadata: &bmpb.BagMetadata{BagId: testRecording2, StartTime: timestamppb.Now(), EndTime: timestamppb.Now()}},
					},
				}, nil
			},
			wantOut: "recording-1",
		},
		{
			name: "No recordings found",
			args: []string{"--workcell", testWorkcell, "--org", testOrg},
			listBagsFunc: func(ctx context.Context, in *pb.ListBagsRequest, opts ...grpc.CallOption) (*pb.ListBagsResponse, error) {
				return &pb.ListBagsResponse{Bags: []*pb.BagRecord{}}, nil
			},
			wantOut: "No recordings found",
		},
		{
			name: "List bags error",
			args: []string{"--workcell", testWorkcell, "--org", testOrg},
			listBagsFunc: func(ctx context.Context, in *pb.ListBagsRequest, opts ...grpc.CallOption) (*pb.ListBagsResponse, error) {
				return nil, errors.New("gRPC error")
			},
			wantErr: "gRPC error",
		},
		{
			name: "Pagination with continue",
			args: []string{"--workcell", testWorkcell, "--org", testOrg},
			listBagsFunc: func() func(context.Context, *pb.ListBagsRequest, ...grpc.CallOption) (*pb.ListBagsResponse, error) {
				callCount := 0
				return func(ctx context.Context, in *pb.ListBagsRequest, opts ...grpc.CallOption) (*pb.ListBagsResponse, error) {
					callCount++
					if callCount == 1 {
						return &pb.ListBagsResponse{
							Bags:           []*pb.BagRecord{{BagMetadata: &bmpb.BagMetadata{BagId: testRecording1, StartTime: timestamppb.Now(), EndTime: timestamppb.Now()}}},
							NextPageCursor: []byte("next-page"),
						}, nil
					}
					return &pb.ListBagsResponse{
						Bags: []*pb.BagRecord{{BagMetadata: &bmpb.BagMetadata{BagId: testRecording2, StartTime: timestamppb.Now(), EndTime: timestamppb.Now()}}},
					}, nil
				}
			}(),
			promptResp: "y\n",
			wantOut:    "Num recordings: 2",
		},
	}

	// We disable the org check in tests because the test environment does not have
	// the necessary home directory configuration.
	originalCheckOrgExists := checkOrgExists
	checkOrgExists = false
	t.Cleanup(func() { checkOrgExists = originalCheckOrgExists })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			originalFlagWorkcellName := flagWorkcellName
			t.Cleanup(func() { flagWorkcellName = originalFlagWorkcellName })
			flagWorkcellName = ""

			originalFlagStartTimestamp := flagStartTimestamp
			t.Cleanup(func() { flagStartTimestamp = originalFlagStartTimestamp })
			flagStartTimestamp = ""

			originalFlagEndTimestamp := flagEndTimestamp
			t.Cleanup(func() { flagEndTimestamp = originalFlagEndTimestamp })
			flagEndTimestamp = ""

			originalFlagMaxNumResults := flagMaxNumResults
			t.Cleanup(func() { flagMaxNumResults = originalFlagMaxNumResults })
			flagMaxNumResults = 10

			mockClient := &mockBagPackagerClientForList{
				ListBagsFunc: tc.listBagsFunc,
			}
			runner := &ListCmdRunner{
				NewClient: func(cmd *cobra.Command) (pb.BagPackagerClient, error) {
					return mockClient, nil
				},
				PromptContinue: func(cmd *cobra.Command) (bool, error) {
					var input string
					if _, err := fmt.Fscanln(cmd.InOrStdin(), &input); err != nil && err != io.EOF {
						return false, err
					}
					input = strings.ToLower(input)
					return input == "y" || input == "", nil
				},
			}

			var out bytes.Buffer
			rootCmd := &cobra.Command{Use: "inctl"}
			recordingsCmd := &cobra.Command{Use: "recordings"}
			listCmd := NewListCmd(runner)
			recordingsCmd.AddCommand(listCmd)
			rootCmd.AddCommand(recordingsCmd)

			in := bytes.NewBufferString(tc.promptResp)
			rootCmd.SetOut(&out)
			rootCmd.SetErr(&out)
			rootCmd.SetIn(in)
			rootCmd.SetArgs(append([]string{"recordings", "list"}, tc.args...))

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
