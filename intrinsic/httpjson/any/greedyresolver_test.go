// Copyright 2023 Intrinsic Innovation LLC

package greedyresolver

import (
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

type DummyResolver struct{}

func (d *DummyResolver) FindMessageByName(message protoreflect.FullName) (protoreflect.MessageType, error) {
	return (&anypb.Any{}).ProtoReflect().Type(), nil
}

func (d *DummyResolver) FindMessageByURL(url string) (protoreflect.MessageType, error) {
	return (&anypb.Any{}).ProtoReflect().Type(), nil
}

func TestGreedyResolver(t *testing.T) {
	noTypes := new(protoregistry.Types)
	dummyTypes := new(DummyResolver)

	noTypesResolver := NewGreedyResolver([]protoregistry.MessageTypeResolver{noTypes})
	_, err := noTypesResolver.FindMessageByName("google.protobuf.Any")
	if err != protoregistry.NotFound {
		t.Errorf("expected noTypesResolver.FindMessageByName to return NotFound, got %v", err)
	}
	_, err = noTypesResolver.FindMessageByURL("type.googleapis.com/google.protobuf.Any")
	if err != protoregistry.NotFound {
		t.Errorf("expected noTypesResolver.FindMessageByURL to return NotFound, got %v", err)
	}

	noTypesThenDummyTypesResolver := NewGreedyResolver([]protoregistry.MessageTypeResolver{noTypes, dummyTypes})
	_, err = noTypesThenDummyTypesResolver.FindMessageByName("google.protobuf.Any")
	if err != nil {
		t.Errorf("expected noTypesThenDummyTypesResolver.FindMessageByName to succeed, got %v", err)
	}
	_, err = noTypesThenDummyTypesResolver.FindMessageByURL("type.googleapis.com/google.protobuf.Any")
	if err != nil {
		t.Errorf("expected noTypesThenDummyTypesResolver.FindMessageByURL to succeed, got %v", err)
	}
}
