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

package connect

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"intrinsic/testing/grpctest"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	tsgrpcpb "intrinsic/assets/instances/connect/testing/test_service_go_proto"
	tspb "intrinsic/assets/instances/connect/testing/test_service_go_proto"
	gcpb "intrinsic/assets/proto/v1/grpc_connection_go_proto"
)

type testService struct{}

func (s *testService) Test(ctx context.Context, req *tspb.TestRequest) (*tspb.TestResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata found in context")
	}

	response := &tspb.TestResponse{
		ContextMetadata: make(map[string]*tspb.TestResponse_MetadataValues),
	}
	for k, vs := range md {
		if strings.HasPrefix(k, ":") || k == "content-type" || k == "user-agent" {
			continue
		}
		response.GetContextMetadata()[k] = &tspb.TestResponse_MetadataValues{
			Values: vs,
		}
	}

	return response, nil
}

func newTestServer(t *testing.T) string {
	t.Helper()

	var s *testService
	grpcServer := grpc.NewServer()
	tsgrpcpb.RegisterTestServiceServer(grpcServer, s)

	return grpctest.StartServerT(t, grpcServer)
}

func TestConnect(t *testing.T) {
	serverAddr := newTestServer(t)

	tests := []struct {
		desc         string
		conn         *gcpb.GrpcConnection
		wantMetadata map[string][]string
		wantError    error
	}{
		{
			desc: "success with metadata",
			conn: &gcpb.GrpcConnection{
				Address: serverAddr,
				Metadata: []*gcpb.GrpcConnection_Metadata{
					{Key: "test_key", Value: "test_value1"},
					{Key: "test_key", Value: "test_value2"},
				},
			},
			wantMetadata: map[string][]string{
				"test_key": {"test_value1", "test_value2"},
			},
		},
		{
			desc: "success without metadata",
			conn: &gcpb.GrpcConnection{
				Address: serverAddr,
			},
		},
		{
			desc:      "nil connection",
			conn:      nil,
			wantError: errNilConnection,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			conn, ctx, err := Connect(context.Background(), tc.conn)
			if diff := cmp.Diff(tc.wantError, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("Connect() error mismatch (-want +got):\n%s", diff)
			}
			if tc.wantError != nil {
				return
			}
			defer conn.Close()

			client := tsgrpcpb.NewTestServiceClient(conn)
			resp, err := client.Test(ctx, &tspb.TestRequest{})
			if err != nil {
				t.Fatalf("Test() error = %v", err)
			}

			gotMetadata := make(map[string][]string)
			for k, v := range resp.GetContextMetadata() {
				gotMetadata[k] = v.GetValues()
			}
			wantMetadata := tc.wantMetadata
			if wantMetadata == nil {
				wantMetadata = make(map[string][]string)
			}
			if diff := cmp.Diff(wantMetadata, gotMetadata, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Context metadata mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
