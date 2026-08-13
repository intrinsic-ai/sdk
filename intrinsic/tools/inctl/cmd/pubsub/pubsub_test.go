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

package pubsub

import (
	"bytes"
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
)

type mockPubSubClient struct {
	lropb.OperationsClient
	getOperationFn func(context.Context, *lropb.GetOperationRequest, ...grpc.CallOption) (*lropb.Operation, error)
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
