// Copyright 2023 Intrinsic Innovation LLC

package pubsub

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"

	commandpb "intrinsic/platform/pubsub/common/command_go_proto"
	pubsubpb "intrinsic/platform/pubsub/connect/cloud/proto/v1alpha1/pubsub_connect_go_proto"
)

func TestHubCreateCmd(t *testing.T) {
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
					createHubFn: func(ctx context.Context, req *pubsubpb.CreateHubRequest, opts ...grpc.CallOption) (*lropb.Operation, error) {
						status := &commandpb.CommandExecutionStatus{
							Succeeded: true,
						}
						anyResp, _ := anypb.New(status)
						return &lropb.Operation{Done: true, Result: &lropb.Operation_Response{Response: anyResp}, Name: "op1"}, nil
					},
				}
			},
			expectedOutput: "Successfully created PubSub Hub",
			expectErr:      false,
		},
		{
			name: "Failure",
			setupMock: func() *mockPubSubClient {
				return &mockPubSubClient{
					createHubFn: func(ctx context.Context, req *pubsubpb.CreateHubRequest, opts ...grpc.CallOption) (*lropb.Operation, error) {
						status := &commandpb.CommandExecutionStatus{
							Succeeded:    false,
							ErrorMessage: "failed",
						}
						anyResp, _ := anypb.New(status)
						return &lropb.Operation{Done: true, Result: &lropb.Operation_Response{Response: anyResp}, Name: "op1"}, nil
					},
				}
			},
			expectedOutput: "failed",
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.setupMock()
			runner := &HubCreateCmdRunner{
				NewClient: func(ctx context.Context) (pubsubpb.PubSubConnectServiceClient, error) {
					return mockClient, nil
				},
			}
			cmd := NewHubCreateCmd(runner)
			cmd.SetContext(context.Background())

			flagHubCreateHubWorkcell = "hub1"
			flagHubCreateStaticEndpoints = []string{"e1", "e2"}
			flagHubCreateParams = viper.New()
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
