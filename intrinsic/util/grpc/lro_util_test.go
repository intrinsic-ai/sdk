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

package lroutil

import (
	"context"
	"testing"

	"intrinsic/testing/grpctest"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/local"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
)

type fakeOperationsServer struct {
	lropb.UnimplementedOperationsServer

	resps []*lropb.Operation
	err   error
	reqs  []*lropb.WaitOperationRequest
}

func (s *fakeOperationsServer) WaitOperation(ctx context.Context, req *lropb.WaitOperationRequest) (*lropb.Operation, error) {
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	op := s.resps[0]
	s.resps = s.resps[1:]
	return op, nil
}

func TestWaitOperation(t *testing.T) {
	wantReq := &lropb.WaitOperationRequest{
		Name:    "operations/op-1",
		Timeout: durationpb.New(RelaySafeWaitTimeout),
	}

	tests := []struct {
		desc     string
		initial  *lropb.Operation
		resps    []*lropb.Operation
		rpcErr   error
		wantDone bool
		wantCode codes.Code
		wantReqs []*lropb.WaitOperationRequest
	}{
		{
			desc:     "already_done",
			initial:  &lropb.Operation{Name: "operations/op-1", Done: true},
			wantDone: true,
			wantCode: codes.OK,
		},
		{
			desc:    "polls_until_done",
			initial: &lropb.Operation{Name: "operations/op-1"},
			resps: []*lropb.Operation{
				{Name: "operations/op-1"},
				{Name: "operations/op-1", Done: true},
			},
			wantDone: true,
			wantCode: codes.OK,
			wantReqs: []*lropb.WaitOperationRequest{wantReq, wantReq},
		},
		{
			desc:    "operation_error",
			initial: &lropb.Operation{Name: "operations/op-1"},
			resps: []*lropb.Operation{
				{
					Name:   "operations/op-1",
					Done:   true,
					Result: &lropb.Operation_Error{Error: status.New(codes.FailedPrecondition, "install failed").Proto()},
				},
			},
			wantDone: true,
			wantCode: codes.FailedPrecondition,
			wantReqs: []*lropb.WaitOperationRequest{wantReq},
		},
		{
			desc:     "rpc_error",
			initial:  &lropb.Operation{Name: "operations/op-1"},
			rpcErr:   status.Error(codes.Unavailable, "relay unavailable"),
			wantCode: codes.Unavailable,
			wantReqs: []*lropb.WaitOperationRequest{wantReq},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			srv := &fakeOperationsServer{resps: tc.resps, err: tc.rpcErr}
			grpcSrv := grpc.NewServer()
			lropb.RegisterOperationsServer(grpcSrv, srv)
			addr := grpctest.StartServerT(t, grpcSrv)

			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(local.NewCredentials()))
			if err != nil {
				t.Fatalf("grpc.NewClient(%q) failed: %v", addr, err)
			}
			t.Cleanup(func() { conn.Close() })

			gotOp, err := WaitOperation(t.Context(), lropb.NewOperationsClient(conn), tc.initial)
			if got := status.Code(err); got != tc.wantCode {
				t.Fatalf("WaitOperation(%v) error code = %v, want %v (err: %v)", tc.initial, got, tc.wantCode, err)
			}
			if gotOp.GetName() != tc.initial.GetName() {
				t.Errorf("WaitOperation(%v) returned Name = %q, want %q", tc.initial, gotOp.GetName(), tc.initial.GetName())
			}
			if gotOp.GetDone() != tc.wantDone {
				t.Errorf("WaitOperation(%v) returned Done = %t, want %t", tc.initial, gotOp.GetDone(), tc.wantDone)
			}
			if diff := cmp.Diff(tc.wantReqs, srv.reqs, protocmp.Transform()); diff != "" {
				t.Errorf("WaitOperation(%v) requests mismatch (-want +got):\n%s", tc.initial, diff)
			}
		})
	}
}
