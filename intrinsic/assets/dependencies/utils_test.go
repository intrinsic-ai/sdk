// Copyright 2023 Intrinsic Innovation LLC

package utils

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"intrinsic/assets/data/fakedataassets"
	"intrinsic/testing/grpctest"
	"intrinsic/util/proto/descriptor"
	"intrinsic/util/proto/testing/prototestutil"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"

	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	tcpb "intrinsic/assets/dependencies/testing/test_configs_go_proto"
	tsgrpcpb "intrinsic/assets/dependencies/testing/test_service_go_grpc_proto"
	tspb "intrinsic/assets/dependencies/testing/test_service_go_grpc_proto"
	atpb "intrinsic/assets/proto/asset_type_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	mpb "intrinsic/assets/proto/metadata_go_proto"
	gcpb "intrinsic/assets/proto/v1/grpc_connection_go_proto"
	rdpb "intrinsic/assets/proto/v1/resolved_dependency_go_proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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
					"grpc://intrinsic_proto.assets.dependencies.testing.TestService": {
						Protocol: &rdpb.ResolvedDependency_Interface_Grpc_{
							Grpc: &rdpb.ResolvedDependency_Interface_Grpc{
								Connection: &gcpb.GrpcConnection{
									Address: serverAddr,
									Metadata: []*gcpb.GrpcConnection_Metadata{
										{
											Key:   "test_key",
											Value: "test_value1",
										},
										{
											Key:   "test_key",
											Value: "test_value2",
										},
									},
								},
							},
						},
					},
				},
			},
			iface: "grpc://intrinsic_proto.assets.dependencies.testing.TestService",
			wantMetadata: map[string][]string{
				"test_key": {"test_value1", "test_value2"},
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
					"data://google.protobuf.Empty": {
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
					"data://google.protobuf.Empty": {
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
					"data://google.protobuf.Empty": {
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
					"grpc://intrinsic_proto.assets.dependencies.testing.TestService": {
						Protocol: &rdpb.ResolvedDependency_Interface_Grpc_{
							Grpc: &rdpb.ResolvedDependency_Interface_Grpc{
								Connection: &gcpb.GrpcConnection{
									Address: "localhost:12345",
								},
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
					"grpc://intrinsic_proto.assets.dependencies.testing.TestService": {
						Protocol: &rdpb.ResolvedDependency_Interface_Grpc_{
							Grpc: &rdpb.ResolvedDependency_Interface_Grpc{
								Connection: &gcpb.GrpcConnection{
									Address: "localhost:12345",
								},
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

func TestHasDependencies(t *testing.T) {
	tests := []struct {
		desc       string
		descriptor protoreflect.MessageDescriptor
		options    []ResolvedDepsIntrospectionOption
		want       bool
	}{
		{
			desc:       "simple_grpc_dependency_config",
			descriptor: (&tcpb.SimpleGrpcDependencyConfig{}).ProtoReflect().Descriptor(),
			want:       true,
		},
		{
			desc:       "simple_grpc_dependency_config_with_dependency_annotation_check",
			descriptor: (&tcpb.SimpleGrpcDependencyConfig{}).ProtoReflect().Descriptor(),
			options:    []ResolvedDepsIntrospectionOption{WithDependencyAnnotation()},
			want:       true,
		},
		{
			desc:       "simple_grpc_dependency_config_with_dependency_annotation_and_skill_annotations_check",
			descriptor: (&tcpb.SimpleGrpcDependencyConfig{}).ProtoReflect().Descriptor(),
			options:    []ResolvedDepsIntrospectionOption{WithDependencyAnnotation(), WithSkillAnnotations()},
			want:       false,
		},
		{
			desc:       "simple_grpc_dependency_config_with_skill_annotations_check",
			descriptor: (&tcpb.SimpleGrpcDependencyConfig{}).ProtoReflect().Descriptor(),
			options:    []ResolvedDepsIntrospectionOption{WithSkillAnnotations()},
			want:       false,
		},
		{
			desc:       "simple_data_dependency_config",
			descriptor: (&tcpb.SimpleDataDependencyConfig{}).ProtoReflect().Descriptor(),
			want:       true,
		},
		{
			desc:       "simple_object_dependency_config",
			descriptor: (&tcpb.SimpleObjectDependencyConfig{}).ProtoReflect().Descriptor(),
			want:       true,
		},
		{
			desc:       "config_with_repeated_dependency",
			descriptor: (&tcpb.WithRepeatedDependency{}).ProtoReflect().Descriptor(),
			want:       true,
		},
		{
			desc:       "config_with_nested_fields",
			descriptor: (&tcpb.WithNestedFields{}).ProtoReflect().Descriptor(),
			want:       true,
		},
		{
			desc:       "config_with_multiple_dependencies",
			descriptor: (&tcpb.WithMultipleDependencies{}).ProtoReflect().Descriptor(),
			want:       true,
		},
		{
			desc:       "config_with_unannotated_dependency",
			descriptor: (&tcpb.WithUnannotatedDependency{}).ProtoReflect().Descriptor(),
			want:       true,
		},
		{
			desc:       "config_with_incorrect_dependency_type",
			descriptor: (&tcpb.WithIncorrectMessageType{}).ProtoReflect().Descriptor(),
			want:       false,
		},
		{
			desc:       "config_with_multiple_interfaces",
			descriptor: (&tcpb.WithMultipleInterfaces{}).ProtoReflect().Descriptor(),
			want:       true,
		},
		{
			desc:       "config_with_map_dependency",
			descriptor: (&tcpb.WithMapDependency{}).ProtoReflect().Descriptor(),
			want:       true,
		},
		{
			desc:       "config_with_optional_dependency",
			descriptor: (&tcpb.WithOptionalDependency{}).ProtoReflect().Descriptor(),
			want:       true,
		},
		{
			desc:       "config_with_required_dependency",
			descriptor: (&tcpb.WithRequiredDependency{}).ProtoReflect().Descriptor(),
			want:       true,
		},
		{
			desc:       "config_with_any_proto",
			descriptor: (&tcpb.WithAnyProto{}).ProtoReflect().Descriptor(),
			want:       false,
		},
		{
			desc:       "skill_with_annotations_dependency_check",
			descriptor: (&tcpb.WithAlwaysProvideConnectionInfo{}).ProtoReflect().Descriptor(),
			want:       true,
		},
		{
			desc:       "skill_with_annotations_only_dependency_annotation_check",
			descriptor: (&tcpb.WithAlwaysProvideConnectionInfo{}).ProtoReflect().Descriptor(),
			options:    []ResolvedDepsIntrospectionOption{WithDependencyAnnotation()},
			want:       true,
		},
		{
			desc:       "skill_with_annotations_skill_annotation_check",
			descriptor: (&tcpb.WithAlwaysProvideConnectionInfo{}).ProtoReflect().Descriptor(),
			options:    []ResolvedDepsIntrospectionOption{WithDependencyAnnotation(), WithSkillAnnotations()},
			want:       true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := HasResolvedDependency(tc.descriptor, (tc.options)...)
			if got != tc.want {
				t.Errorf("HasResolvedDependency(%v) = %v, want: %v", tc.descriptor, got, tc.want)
			}
		})
	}
}
