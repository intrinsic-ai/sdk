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

	bmpb "intrinsic/logging/proto/bag_metadata_go_proto"
	pb "intrinsic/logging/proto/bag_packager_service_go_grpc_proto"
)

type mockBagPackagerClientForGenerate struct {
	pb.BagPackagerClient
	GetBagFunc      func(ctx context.Context, in *pb.GetBagRequest, opts ...grpc.CallOption) (*pb.GetBagResponse, error)
	GenerateBagFunc func(ctx context.Context, in *pb.GenerateBagRequest, opts ...grpc.CallOption) (*pb.GenerateBagResponse, error)
}

func (m *mockBagPackagerClientForGenerate) GetBag(ctx context.Context, in *pb.GetBagRequest, opts ...grpc.CallOption) (*pb.GetBagResponse, error) {
	if m.GetBagFunc != nil {
		return m.GetBagFunc(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "mock GetBag should not be called directly; set GetBagFunc in test case")
}

func (m *mockBagPackagerClientForGenerate) GenerateBag(ctx context.Context, in *pb.GenerateBagRequest, opts ...grpc.CallOption) (*pb.GenerateBagResponse, error) {
	if m.GenerateBagFunc != nil {
		return m.GenerateBagFunc(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "mock GenerateBag should not be called directly; set GenerateBagFunc in test case")
}

func TestGenerateRecordingE(t *testing.T) {
	const (
		testBagID = "test-bag-id"
		testOrg   = "test-org"
	)

	tests := []struct {
		name            string
		args            []string
		getBagFunc      func(ctx context.Context, in *pb.GetBagRequest, opts ...grpc.CallOption) (*pb.GetBagResponse, error)
		generateBagFunc func(ctx context.Context, in *pb.GenerateBagRequest, opts ...grpc.CallOption) (*pb.GenerateBagResponse, error)
		wantErr         string
		wantOut         string
	}{
		{
			name: "Successful generation",
			args: []string{"--recording_id", testBagID, "--org", testOrg},
			getBagFunc: func(ctx context.Context, in *pb.GetBagRequest, opts ...grpc.CallOption) (*pb.GetBagResponse, error) {
				return &pb.GetBagResponse{
					Bag: &pb.BagRecord{
						BagMetadata: &bmpb.BagMetadata{},
					},
				}, nil
			},
			generateBagFunc: func(ctx context.Context, in *pb.GenerateBagRequest, opts ...grpc.CallOption) (*pb.GenerateBagResponse, error) {
				return &pb.GenerateBagResponse{}, nil
			},
			wantOut: "Generated recording file",
		},
		{
			name: "Already generated",
			args: []string{"--recording_id", testBagID, "--org", testOrg},
			getBagFunc: func(ctx context.Context, in *pb.GetBagRequest, opts ...grpc.CallOption) (*pb.GetBagResponse, error) {
				return &pb.GetBagResponse{
					Bag: &pb.BagRecord{
						BagFile: &bmpb.BagFileReference{},
					},
				}, nil
			},
			wantErr: "is already generated",
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
				return nil, errors.New("does not exist")
			},
			wantErr: "does not exist",
		},
		{
			name: "Generate bag error",
			args: []string{"--recording_id", testBagID, "--org", testOrg},
			getBagFunc: func(ctx context.Context, in *pb.GetBagRequest, opts ...grpc.CallOption) (*pb.GetBagResponse, error) {
				return &pb.GetBagResponse{
					Bag: &pb.BagRecord{
						BagMetadata: &bmpb.BagMetadata{},
					},
				}, nil
			},
			generateBagFunc: func(ctx context.Context, in *pb.GenerateBagRequest, opts ...grpc.CallOption) (*pb.GenerateBagResponse, error) {
				return nil, errors.New("gRPC error")
			},
			wantErr: "gRPC error",
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

			mockClient := &mockBagPackagerClientForGenerate{
				GetBagFunc:      tc.getBagFunc,
				GenerateBagFunc: tc.generateBagFunc,
			}
			runner := &GenerateCmdRunner{
				NewClient: func(cmd *cobra.Command) (pb.BagPackagerClient, error) {
					return mockClient, nil
				},
			}

			var out bytes.Buffer
			rootCmd := &cobra.Command{Use: "inctl"}
			recordingsCmd := &cobra.Command{Use: "recordings"}
			generateCmd := NewGenerateCmd(runner)
			recordingsCmd.AddCommand(generateCmd)

			rootCmd.AddCommand(recordingsCmd)
			rootCmd.SetOut(&out)
			rootCmd.SetErr(&out)
			rootCmd.SetArgs(append([]string{"recordings", "generate"}, tc.args...))

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
