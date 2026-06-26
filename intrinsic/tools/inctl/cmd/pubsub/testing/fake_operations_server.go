// Copyright 2023 Intrinsic Innovation LLC

package pubsubtesting

import (
	"context"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
)

type FakeOperationsServer struct {
	lropb.UnimplementedOperationsServer
	GetOperationFn func(ctx context.Context, req *lropb.GetOperationRequest) (*lropb.Operation, error)
}

func NewFakeOperationsServer() *FakeOperationsServer {
	return &FakeOperationsServer{}
}

func (s *FakeOperationsServer) GetOperation(ctx context.Context, req *lropb.GetOperationRequest) (*lropb.Operation, error) {
	if s.GetOperationFn != nil {
		return s.GetOperationFn(ctx, req)
	}
	return nil, nil
}
