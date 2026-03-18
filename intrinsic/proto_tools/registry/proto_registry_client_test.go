// Copyright 2023 Intrinsic Innovation LLC

package protoregistryclient

import (
	"context"
	"fmt"
	"testing"

	"intrinsic/testing/grpctest"
	"intrinsic/util/proto/registryutil"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	protoregistrypb "intrinsic/proto_tools/proto/proto_registry_go_proto"
)

type mockProtoRegistry struct {
	protoregistrypb.UnimplementedProtoRegistryServer

	errToReturn                 error
	fileDescriptorSetsByTypeURL map[string]*protoregistrypb.NamedFileDescriptorSet
}

func NewMockProtoRegistry() *mockProtoRegistry {
	return &mockProtoRegistry{
		fileDescriptorSetsByTypeURL: make(map[string]*protoregistrypb.NamedFileDescriptorSet),
	}
}

func (r *mockProtoRegistry) addFileDescriptorSet(typeURL string, fileDescriptorSet *descriptorpb.FileDescriptorSet) {
	r.fileDescriptorSetsByTypeURL[typeURL] = &protoregistrypb.NamedFileDescriptorSet{
		// In practice the name would need to be unique, but we don't need that in this test.
		Name:              "irrelevant_name",
		FileDescriptorSet: fileDescriptorSet,
	}
}

func (r *mockProtoRegistry) setReturnError(err error) {
	r.errToReturn = err
}

func (r *mockProtoRegistry) clearFileDescriptorSets() {
	r.fileDescriptorSetsByTypeURL = make(map[string]*protoregistrypb.NamedFileDescriptorSet)
}

func (r *mockProtoRegistry) GetNamedFileDescriptorSet(_ context.Context, request *protoregistrypb.GetNamedFileDescriptorSetRequest) (*protoregistrypb.NamedFileDescriptorSet, error) {
	if r.errToReturn != nil {
		return nil, r.errToReturn
	}

	if request.GetTypeUrl() == "" {
		return nil, status.Error(codes.InvalidArgument, "Request does not have 'type_url' set")
	}

	namedFds, ok := r.fileDescriptorSetsByTypeURL[request.GetTypeUrl()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "File descriptor set for type URL %q not found", request.GetTypeUrl())
	}

	return namedFds, nil
}

type protoRegistryResolverFixture struct {
	protoRegistry *mockProtoRegistry
	resolver      *ProtoRegistryResolver
}

func mustCreateProtoRegistryResolver(t *testing.T, defaultResolvers []Resolver) protoRegistryResolverFixture {
	t.Helper()

	protoRegistry := NewMockProtoRegistry()

	server := grpc.NewServer()
	protoregistrypb.RegisterProtoRegistryServer(server, protoRegistry)

	conn, err := grpc.NewClient(
		grpctest.StartServerT(t, server),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Unable to connect to test server: %v", err)
	}

	resolver := NewProtoRegistryResolver(
		context.Background(),
		protoregistrypb.NewProtoRegistryClient(conn),
		defaultResolvers,
	)

	return protoRegistryResolverFixture{
		protoRegistry: protoRegistry,
		resolver:      resolver,
	}
}

func messageTypeWithNFields(name string, numFields int) *descriptorpb.DescriptorProto {
	result := &descriptorpb.DescriptorProto{
		Name: &name,
	}
	for i := range numFields {
		result.Field = append(result.Field, &descriptorpb.FieldDescriptorProto{
			Name:   proto.String(fmt.Sprintf("my_field_%d", i)),
			Number: proto.Int32(int32(i + 1)),
			Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
		})
	}
	return result
}

func fileDescriptorSetWithMessages(protoPackage string, messageTypes ...*descriptorpb.DescriptorProto) *descriptorpb.FileDescriptorSet {
	return &descriptorpb.FileDescriptorSet{
		File: []*descriptorpb.FileDescriptorProto{
			{
				Name:        proto.String("my_messages.proto"),
				Package:     &protoPackage,
				MessageType: messageTypes,
			},
		},
	}
}

func mustCreateTypesFromFileDescriptorSet(t *testing.T, fileDescriptorSet *descriptorpb.FileDescriptorSet) *protoregistry.Types {
	t.Helper()

	types, err := registryutil.NewTypesFromFileDescriptorSet(fileDescriptorSet)
	if err != nil {
		t.Fatalf("failed creating types from file descriptor set: %v", err)
	}
	return types
}

func TestProtoRegistryResolverResolvesIntrinsicTypeURL(t *testing.T) {
	fixture := mustCreateProtoRegistryResolver(t, []Resolver{})

	typeURL := "type.intrinsic.ai/foo/bar/my_package.MyMessage"
	msgName := "my_package.MyMessage"
	fixture.protoRegistry.addFileDescriptorSet(
		typeURL, fileDescriptorSetWithMessages("my_package", messageTypeWithNFields("MyMessage", 0)),
	)

	msgType, err := fixture.resolver.FindMessageByURL(typeURL)
	if err != nil {
		t.Fatalf("Unexpected error from resolver.FindMessageByURL(%q) = %v, want nil", typeURL, err)
	}

	if msgType.Descriptor().FullName() != protoreflect.FullName(msgName) {
		t.Errorf(
			"Unexpected message type returned from resolver.FindMessageByURL(%q), got message name %v, want %q",
			typeURL, msgType.Descriptor().FullName(), msgName,
		)
	}

	// Clear proto registry and call resolver.FindMessageByURL() again. This
	// should hit the local cache of ProtoRegistryResolver and succeed instead
	// of forwarding a not-found error from the proto registry.
	fixture.protoRegistry.clearFileDescriptorSets()

	msgType, err = fixture.resolver.FindMessageByURL(typeURL)
	if err != nil {
		t.Fatalf("Unexpected error from resolver.FindMessageByURL(%q) = %v, want nil", typeURL, err)
	}

	if msgType.Descriptor().FullName() != protoreflect.FullName(msgName) {
		t.Errorf(
			"Unexpected message type returned from resolver.FindMessageByURL(%q), got message name %v, want %q",
			typeURL, msgType.Descriptor().FullName(), msgName,
		)
	}
}

func TestProtoRegistryResolverForwardsErrorsFromProtoRegistry(t *testing.T) {
	fixture := mustCreateProtoRegistryResolver(t, []Resolver{})

	typeURL := "type.intrinsic.ai/foo/bar/my_package.MyMessage"
	fixture.protoRegistry.setReturnError(
		status.Error(codes.Internal, "Some error"),
	)

	_, err := fixture.resolver.FindMessageByURL(typeURL)
	if status.Code(err) != codes.Internal {
		t.Fatalf(
			"Unexpected error code from resolver.FindMessageByURL(%q) = %v, want Internal",
			typeURL, status.Code(err),
		)
	}
}

func TestProtoRegistryResolverResolvesNonIntrinsicTypeURL(t *testing.T) {
	typesOne := mustCreateTypesFromFileDescriptorSet(t,
		fileDescriptorSetWithMessages(
			"my_package",
			messageTypeWithNFields("MessageInBoth", 1),
			messageTypeWithNFields("MessageInFirst", 0),
		),
	)
	typesTwo := mustCreateTypesFromFileDescriptorSet(t,
		fileDescriptorSetWithMessages(
			"my_package",
			messageTypeWithNFields("MessageInBoth", 3),
			messageTypeWithNFields("MessageInSecond", 0),
		),
	)

	fixture := mustCreateProtoRegistryResolver(t, []Resolver{typesOne, typesTwo})

	// MessageInBoth should be returned from the first Types and not from the
	// second Types and thus have 1 field.
	typeURL := "type.googleapis.com/my_package.MessageInBoth"
	msgType, err := fixture.resolver.FindMessageByURL(typeURL)
	if err != nil {
		t.Fatalf("Unexpected error from resolver.FindMessageByURL(%q) = %v, want nil", typeURL, err)
	}

	if msgType.Descriptor().Fields().Len() != 1 {
		t.Errorf(
			"Unexpected message type returned from resolver.FindMessageByURL(%q), got message with %d fields, want 1",
			typeURL,
			msgType.Descriptor().Fields().Len(),
		)
	}

	// MessageInFirst should be returned successfully from the first Types.
	typeURL = "some.custom.type.domain/my_package.MessageInFirst"
	msgName := "my_package.MessageInFirst"
	msgType, err = fixture.resolver.FindMessageByURL(typeURL)
	if err != nil {
		t.Fatalf("Unexpected error from resolver.FindMessageByURL(%q) = %v, want nil", typeURL, err)
	}

	if msgType.Descriptor().FullName() != protoreflect.FullName(msgName) {
		t.Errorf(
			"Unexpected message type returned from resolver.FindMessageByURL(%q), got message name %v, want %q",
			typeURL, msgType.Descriptor().FullName(), msgName,
		)
	}

	// MessageInSecond should be returned successfully from the second Types.
	typeURL = "another.custom.type/domain/my_package.MessageInSecond"
	msgName = "my_package.MessageInSecond"
	msgType, err = fixture.resolver.FindMessageByURL(typeURL)
	if err != nil {
		t.Fatalf("Unexpected error from resolver.FindMessageByURL(%q) = %v, want nil", typeURL, err)
	}

	if msgType.Descriptor().FullName() != protoreflect.FullName(msgName) {
		t.Errorf(
			"Unexpected message type returned from resolver.FindMessageByURL(%q), got message name %v, want %q",
			typeURL, msgType.Descriptor().FullName(), msgName,
		)
	}

	// A message not defined in either of the two Types should lead to an error.
	typeURL = "type.googleapis.com/my_package.NonExistingMessage"
	msgType, err = fixture.resolver.FindMessageByURL(typeURL)
	wantErr := protoregistry.NotFound
	if diff := cmp.Diff(wantErr, err, cmpopts.EquateErrors()); diff != "" {
		t.Fatalf("resolver.FindMessageByURL(%q) returned unexpected error diff (-want +got):\n%s", typeURL, diff)
	}
}
