// Copyright 2023 Intrinsic Innovation LLC

package utils

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"intrinsic/testing/grpctest"

	tsgrpcpb "intrinsic/assets/dependencies/testing/test_service_go_grpc_proto"
	tspb "intrinsic/assets/dependencies/testing/test_service_go_grpc_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	rdpb "intrinsic/assets/proto/v1/resolved_dependency_go_proto"
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
		dep          *rdpb.ResolvedDependency
		iface        string
		wantMetadata map[string][]string
		wantErr      string
	}{
		{
			desc: "success",
			dep: &rdpb.ResolvedDependency{
				Interfaces: map[string]*rdpb.ResolvedDependency_Interface{
					"grpc://intrinsic_proto.assets.dependencies.testing.TestService": &rdpb.ResolvedDependency_Interface{
						Protocol: &rdpb.ResolvedDependency_Interface_GrpcConnection_{
							GrpcConnection: &rdpb.ResolvedDependency_Interface_GrpcConnection{
								Address: serverAddr,
								Metadata: []*rdpb.ResolvedDependency_Interface_GrpcConnection_Metadata{
									&rdpb.ResolvedDependency_Interface_GrpcConnection_Metadata{
										Key:   "test_key",
										Value: "test_value1",
									},
									&rdpb.ResolvedDependency_Interface_GrpcConnection_Metadata{
										Key:   "test_key",
										Value: "test_value2",
									},
								},
							},
						},
					},
				},
			},
			iface: "grpc://intrinsic_proto.assets.dependencies.testing.TestService",
			wantMetadata: map[string][]string{
				"test_key": []string{"test_value1", "test_value2"},
			},
		},
		{
			desc:    "missing interface",
			dep:     &rdpb.ResolvedDependency{},
			iface:   "grpc://intrinsic_proto.assets.dependencies.testing.TestService",
			wantErr: "not found in resolved dependency",
		},
		{
			desc: "not a gRPC interface",
			dep: &rdpb.ResolvedDependency{
				Interfaces: map[string]*rdpb.ResolvedDependency_Interface{
					"data://foo": &rdpb.ResolvedDependency_Interface{
						Protocol: &rdpb.ResolvedDependency_Interface_Data_{
							Data: &rdpb.ResolvedDependency_Interface_Data{
								Id: &idpb.Id{Package: "foo", Name: "bar"},
							},
						},
					},
				},
			},
			iface:   "data://foo",
			wantErr: "is not gRPC",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			conn, ctx, err := Connect(context.Background(), tc.dep, tc.iface)
			if tc.wantErr != "" {
				if err == nil {
					t.Errorf("Connect() returned no error, want error %q", tc.wantErr)
				} else if !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("Connect() returned error %q, want error %q", err, tc.wantErr)
				}
			} else if err != nil {
				t.Errorf("Connect() returned an unexpected error: %v", err)
			} else {
				defer conn.Close()

				client := tsgrpcpb.NewTestServiceClient(conn)
				response, err := client.Test(ctx, &tspb.TestRequest{})
				if err != nil {
					t.Fatalf("Test() returned an unexpected error: %v", err)
				}
				for k, vs := range tc.wantMetadata {
					gotVs := response.GetContextMetadata()[k].GetValues()
					if diff := cmp.Diff(vs, gotVs); diff != "" {
						t.Errorf("Test() returned an unexpected diff for metadata key %q (-want +got): %v", k, diff)
					}
				}
			}
		})
	}
}
