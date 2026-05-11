// Copyright 2023 Intrinsic Innovation LLC

package pubsub

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	commandpb "intrinsic/platform/pubsub/common/command_go_proto"
	pubsubpb "intrinsic/platform/pubsub/connect/cloud/proto/v1alpha1/pubsub_connect_go_proto"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestForwardCmd(t *testing.T) {
	origPoll := operationPollInterval
	operationPollInterval = 1 * time.Millisecond
	defer func() { operationPollInterval = origPoll }()

	tests := []struct {
		name           string
		setupMock      func() *mockPubSubClient
		expectedOutput string
		expectErr      bool
	}{
		{
			name: "Success immediate",
			setupMock: func() *mockPubSubClient {
				return &mockPubSubClient{
					configureSpokeFn: func(ctx context.Context, req *pubsubpb.ConfigureSpokeRequest, opts ...grpc.CallOption) (*lropb.Operation, error) {
						resp := &commandpb.CommandExecutionStatus{Succeeded: true}
						anyResp, _ := anypb.New(resp)
						return &lropb.Operation{Done: true, Result: &lropb.Operation_Response{Response: anyResp}, Name: "op1"}, nil
					},
				}
			},
			expectedOutput: "Successfully configured forwarding allowed topics for spoke",
			expectErr:      false,
		},
		{
			name: "Backend error",
			setupMock: func() *mockPubSubClient {
				return &mockPubSubClient{
					configureSpokeFn: func(ctx context.Context, req *pubsubpb.ConfigureSpokeRequest, opts ...grpc.CallOption) (*lropb.Operation, error) {
						resp := &commandpb.CommandExecutionStatus{Succeeded: false, ErrorMessage: "failed configure"}
						anyResp, _ := anypb.New(resp)
						return &lropb.Operation{Done: true, Result: &lropb.Operation_Response{Response: anyResp}, Name: "op1"}, nil
					},
				}
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.setupMock()
			runner := &ForwardCmdRunner{
				NewClient: func(ctx context.Context) (pubsubpb.PubSubConnectServiceClient, error) {
					return mockClient, nil
				},
			}
			cmd := NewForwardCmd(runner)
			cmd.SetContext(context.Background())

			flagForwardWorkcell = "test-workcell"
			flagForwardAllowedTopics = nil
			flagForwardAllowedKVStoreKeys = nil
			flagForwardParams = viper.New()
			buf := &bytes.Buffer{}
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.RunE(cmd, []string{})
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected err state %v, got %v", tt.expectErr, err)
			}

			if tt.expectedOutput != "" && !strings.Contains(buf.String(), tt.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tt.expectedOutput, buf.String())
			}
		})
	}
}
