// Copyright 2023 Intrinsic Innovation LLC

package pubsub

import (
	"bytes"
	"context"
	"testing"
	"time"

	pubsubpb "intrinsic/platform/pubsub/connect/cloud/proto/v1alpha1/pubsub_connect_go_proto"

	"google.golang.org/grpc"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
)

type mockPubSubClient struct {
	pubsubpb.PubSubConnectServiceClient
	configureSpokeFn func(context.Context, *pubsubpb.ConfigureSpokeRequest, ...grpc.CallOption) (*lropb.Operation, error)
	getOperationFn   func(context.Context, *lropb.GetOperationRequest, ...grpc.CallOption) (*lropb.Operation, error)
}

func (m *mockPubSubClient) ConfigureSpoke(ctx context.Context, in *pubsubpb.ConfigureSpokeRequest, opts ...grpc.CallOption) (*lropb.Operation, error) {
	return m.configureSpokeFn(ctx, in, opts...)
}

func (m *mockPubSubClient) GetOperation(ctx context.Context, in *lropb.GetOperationRequest, opts ...grpc.CallOption) (*lropb.Operation, error) {
	return m.getOperationFn(ctx, in, opts...)
}

func TestWaitForOperation(t *testing.T) {
	origPoll := operationPollInterval
	operationPollInterval = 1 * time.Millisecond
	defer func() { operationPollInterval = origPoll }()

	ctx := context.Background()
	out := &bytes.Buffer{}

	t.Run("Immediate success", func(t *testing.T) {
		op := &lropb.Operation{Done: true, Name: "op1"}
		res, err := waitForOperation(ctx, &mockPubSubClient{}, op, out)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if res.GetName() != "op1" {
			t.Errorf("expected name op1, got %v", res.GetName())
		}
	})

	t.Run("Success after poll", func(t *testing.T) {
		callCount := 0
		client := &mockPubSubClient{
			getOperationFn: func(ctx context.Context, req *lropb.GetOperationRequest, opts ...grpc.CallOption) (*lropb.Operation, error) {
				callCount++
				if callCount == 1 {
					return &lropb.Operation{Done: false, Name: "op2"}, nil
				}
				return &lropb.Operation{Done: true, Name: "op2"}, nil
			},
		}
		op := &lropb.Operation{Done: false, Name: "op2"}
		res, err := waitForOperation(ctx, client, op, out)
		if err != nil {
			t.Fatalf("waitForOperation failed: %v", err)
		}
		if !res.GetDone() {
			t.Error("expected operation to be done")
		}
		if callCount != 2 {
			t.Errorf("expected 2 calls to GetOperation, got %d", callCount)
		}
	})
}
