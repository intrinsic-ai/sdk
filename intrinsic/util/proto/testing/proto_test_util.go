// Copyright 2023 Intrinsic Innovation LLC

// Package prototestutil contains helpers for handling protos in tests.
package prototestutil

import (
	"fmt"
	"strings"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"

	descpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// MustWrapInAny is a test helper that wraps a proto message in an Any proto.
func MustWrapInAny(t *testing.T, m proto.Message) *anypb.Any {
	t.Helper()
	p, err := anypb.New(m)
	if err != nil {
		t.Fatalf("Unable to wrap proto message: %v ", err)
	}
	return p
}

// MustMarshalJSON marshals the given proto to JSON (with the default resolver)
// and fails the test if marshalling fails.
func MustMarshalJSON(t *testing.T, m proto.Message) string {
	t.Helper()
	b, err := protojson.Marshal(m)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	return string(b)
}

// MustMarshalText marshals the given proto to textproto (with the default resolver)
// and fails the test if marshalling fails.
func MustMarshalText(t *testing.T, m proto.Message) string {
	t.Helper()
	b, err := prototext.Marshal(m)
	if err != nil {
		t.Fatalf("Failed to marshal textproto: %v", err)
	}
	return string(b)
}

// FileDescriptorSet creates a FileDescriptorSet proto from the given messages.
func FileDescriptorSet(messages ...proto.Message) *descpb.FileDescriptorSet {
	fds := &descpb.FileDescriptorSet{}
	for _, m := range messages {
		f := m.ProtoReflect().Descriptor().ParentFile()
		fds.File = append(fds.File, protodesc.ToFileDescriptorProto(f))
	}
	return fds
}

// descriptorProtoBuildPackage is the package name of the `descriptor.proto`
// file that was used to build this code.
var descriptorProtoBuildPackage = string((&descpb.FileDescriptorSet{}).ProtoReflect().Descriptor().ParentFile().Package())

// WithDescriptorProtoPackage updates the package name of the `descriptor.proto`
// file in the file descriptor set to the given package name and updates all
// references in other files in the file descriptor set to match. The file
// descriptor set should contain the `descriptor.proto` file and may contain any
// other files (exceptions below).
//
// This function uses a simple text-based replacement. The passed descriptor set
// must NOT contain files that use `google.protobuf` as the package name
// internally and externally (`Timestamp`, `Duration`, etc.).
func WithDescriptorProtoPackage(t *testing.T, fds *descpb.FileDescriptorSet, pkg string) *descpb.FileDescriptorSet {
	txtBytes, err := prototext.Marshal(fds)
	if err != nil {
		t.Fatalf("Failed to marshal file descriptor set: %v", err)
	}
	txt := string(txtBytes)
	txt = strings.ReplaceAll(txt, fmt.Sprintf(".%s.", descriptorProtoBuildPackage), fmt.Sprintf(".%s.", pkg))
	txt = strings.ReplaceAll(txt, fmt.Sprintf("package:%q", descriptorProtoBuildPackage), fmt.Sprintf("package:%q", pkg))
	fds = &descpb.FileDescriptorSet{}
	if err := prototext.Unmarshal([]byte(txt), fds); err != nil {
		t.Fatalf("Failed to unmarshal file descriptor set: %v", err)
	}
	return fds
}
