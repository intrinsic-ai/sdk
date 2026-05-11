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

func TestHubDeleteCmd(t *testing.T) {
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
					deleteHubFn: func(ctx context.Context, req *pubsubpb.DeleteHubRequest, opts ...grpc.CallOption) (*lropb.Operation, error) {
						status := &commandpb.CommandExecutionStatus{
							Succeeded: true,
						}
						anyResp, _ := anypb.New(status)
						return &lropb.Operation{Done: true, Result: &lropb.Operation_Response{Response: anyResp}, Name: "op1"}, nil
					},
				}
			},
			expectedOutput: "Successfully deleted PubSub Hub",
			expectErr:      false,
		},
		{
			name: "Backend error",
			setupMock: func() *mockPubSubClient {
				return &mockPubSubClient{
					deleteHubFn: func(ctx context.Context, req *pubsubpb.DeleteHubRequest, opts ...grpc.CallOption) (*lropb.Operation, error) {
						status := &commandpb.CommandExecutionStatus{Succeeded: false, ErrorMessage: "delete problem"}
						anyResp, _ := anypb.New(status)
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
			runner := &HubDeleteCmdRunner{
				NewClient: func(ctx context.Context) (pubsubpb.PubSubConnectServiceClient, error) {
					return mockClient, nil
				},
			}
			cmd := NewHubDeleteCmd(runner)
			cmd.SetContext(context.Background())

			flagHubDeleteHubWorkcell = "hub1"
			flagHubDeleteParams = viper.New()
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
