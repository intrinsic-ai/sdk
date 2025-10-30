// Copyright 2023 Intrinsic Innovation LLC

package utils

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/testing/protocmp"
	"intrinsic/assets/data/fakedataassets"
	"intrinsic/testing/grpctest"
	"intrinsic/util/proto/descriptor"
	"intrinsic/util/proto/testing/prototestutil"

	anypb "google.golang.org/protobuf/types/known/anypb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	tsgrpcpb "intrinsic/assets/dependencies/testing/test_service_go_grpc_proto"
	tspb "intrinsic/assets/dependencies/testing/test_service_go_grpc_proto"
	atpb "intrinsic/assets/proto/asset_type_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	mpb "intrinsic/assets/proto/metadata_go_proto"
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
		desc          string
		dep           *rdpb.ResolvedDependency
		iface         string
		wantMetadata  map[string][]string
		wantErrorType error
		wantError     string
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
			desc:          "no interfaces",
			dep:           &rdpb.ResolvedDependency{},
			iface:         "grpc://intrinsic_proto.assets.dependencies.testing.TestService",
			wantErrorType: errMissingInterface,
			wantError:     "no interfaces provided",
		},
		{
			desc: "wrong interface type",
			dep: &rdpb.ResolvedDependency{
				Interfaces: map[string]*rdpb.ResolvedDependency_Interface{
					"data://google.protobuf.Empty": &rdpb.ResolvedDependency_Interface{
						Protocol: &rdpb.ResolvedDependency_Interface_Data_{
							Data: &rdpb.ResolvedDependency_Interface_Data{
								Id: &idpb.Id{Package: "ai.intrinsic", Name: "data_asset"},
							},
						},
					},
				},
			},
			iface:         "grpc://intrinsic_proto.assets.dependencies.testing.TestService",
			wantErrorType: errMissingInterface,
			wantError:     "got interfaces: data://google.protobuf.Empty",
		},
		{
			desc: "not a gRPC interface",
			dep: &rdpb.ResolvedDependency{
				Interfaces: map[string]*rdpb.ResolvedDependency_Interface{
					"data://google.protobuf.Empty": &rdpb.ResolvedDependency_Interface{
						Protocol: &rdpb.ResolvedDependency_Interface_Data_{
							Data: &rdpb.ResolvedDependency_Interface_Data{
								Id: &idpb.Id{Package: "ai.intrinsic", Name: "data_asset"},
							},
						},
					},
				},
			},
			iface:         "data://google.protobuf.Empty",
			wantErrorType: errNotGRPC,
			wantError:     "is not gRPC",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			conn, ctx, err := Connect(context.Background(), tc.dep, tc.iface)
			if tc.wantErrorType != nil || tc.wantError != "" {
				if err == nil {
					t.Errorf("Connect() returned no error, want error (%q, %q)", tc.wantErrorType, tc.wantError)
				} else if tc.wantErrorType != nil && !errors.Is(err, tc.wantErrorType) {
					t.Errorf("Connect() returned error %q, want error type %q", err, tc.wantErrorType)
				} else if !strings.Contains(err.Error(), tc.wantError) {
					t.Errorf("Connect() returned error %q, want error %q", err, tc.wantError)
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

func TestGetDataPayload(t *testing.T) {
	payload := &emptypb.Empty{}
	payloadAny := prototestutil.MustWrapInAny(t, payload)
	da := &dapb.DataAsset{
		Data:              payloadAny,
		FileDescriptorSet: descriptor.FileDescriptorSetFrom(payload),
		Metadata: &mpb.Metadata{
			AssetType: atpb.AssetType_ASSET_TYPE_DATA,
			IdVersion: &idpb.IdVersion{
				Id: &idpb.Id{
					Package: "ai.intrinsic",
					Name:    "data_asset",
				},
				Version: "0.0.1",
			},
		},
	}

	tests := []struct {
		desc          string
		dep           *rdpb.ResolvedDependency
		iface         string
		wantPayload   *anypb.Any
		wantErrorType error
		wantError     string
	}{
		{
			desc: "success",
			dep: &rdpb.ResolvedDependency{
				Interfaces: map[string]*rdpb.ResolvedDependency_Interface{
					"data://google.protobuf.Empty": &rdpb.ResolvedDependency_Interface{
						Protocol: &rdpb.ResolvedDependency_Interface_Data_{
							Data: &rdpb.ResolvedDependency_Interface_Data{
								Id: &idpb.Id{Package: "ai.intrinsic", Name: "data_asset"},
							},
						},
					},
				},
			},
			iface:       "data://google.protobuf.Empty",
			wantPayload: payloadAny,
		},
		{
			desc:          "no interfaces",
			dep:           &rdpb.ResolvedDependency{},
			iface:         "data://google.protobuf.Empty",
			wantErrorType: errMissingInterface,
			wantError:     "no interfaces provided",
		},
		{
			desc: "wrong interface type",
			dep: &rdpb.ResolvedDependency{
				Interfaces: map[string]*rdpb.ResolvedDependency_Interface{
					"grpc://intrinsic_proto.assets.dependencies.testing.TestService": &rdpb.ResolvedDependency_Interface{
						Protocol: &rdpb.ResolvedDependency_Interface_GrpcConnection_{
							GrpcConnection: &rdpb.ResolvedDependency_Interface_GrpcConnection{
								Address: "localhost:12345",
							},
						},
					},
				},
			},
			iface:         "data://google.protobuf.Empty",
			wantErrorType: errMissingInterface,
			wantError:     "got interfaces: grpc://intrinsic_proto.assets.dependencies.testing.TestService",
		},
		{
			desc: "not a data interface",
			dep: &rdpb.ResolvedDependency{
				Interfaces: map[string]*rdpb.ResolvedDependency_Interface{
					"grpc://intrinsic_proto.assets.dependencies.testing.TestService": &rdpb.ResolvedDependency_Interface{
						Protocol: &rdpb.ResolvedDependency_Interface_GrpcConnection_{
							GrpcConnection: &rdpb.ResolvedDependency_Interface_GrpcConnection{
								Address: "localhost:12345",
							},
						},
					},
				},
			},
			iface:         "grpc://intrinsic_proto.assets.dependencies.testing.TestService",
			wantErrorType: errNotData,
			wantError:     "is not data",
		},
	}

	ctx := context.Background()
	fake := fakedataassets.StartServer(ctx, t, fakedataassets.WithDataAssets([]*dapb.DataAsset{da}))

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			gotPayload, err := GetDataPayload(ctx, tc.dep, tc.iface, WithDataAssetsClient(fake.Client))
			if tc.wantErrorType != nil || tc.wantError != "" {
				if err == nil {
					t.Errorf("GetDataPayload() returned no error, want error (%q, %q)", tc.wantErrorType, tc.wantError)
				} else if tc.wantErrorType != nil && !errors.Is(err, tc.wantErrorType) {
					t.Errorf("GetDataPayload() returned error %q, want error type %q", err, tc.wantErrorType)
				} else if !strings.Contains(err.Error(), tc.wantError) {
					t.Errorf("GetDataPayload() returned error %q, want error %q", err, tc.wantError)
				}
			} else if err != nil {
				t.Errorf("GetDataPayload() returned an unexpected error: %v", err)
			} else if diff := cmp.Diff(tc.wantPayload, gotPayload, protocmp.Transform()); diff != "" {
				t.Errorf("GetDataPayload() returned an unexpected diff (-want +got): %v", diff)
			}
		})
	}
}
