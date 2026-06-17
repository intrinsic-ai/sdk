// Copyright 2023 Intrinsic Innovation LLC

package any

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

func (d *DummyResolver) FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error) {
	return nil, protoregistry.NotFound
}

func (d *DummyResolver) FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error) {
	return nil, protoregistry.NotFound
}

func TestGreedyResolver(t *testing.T) {
	noTypes := new(protoregistry.Types)
	dummyTypes := new(DummyResolver)

	noTypesResolver := NewGreedyResolver([]Resolver{noTypes})
	_, err := noTypesResolver.FindMessageByName("google.protobuf.Any")
	if err != protoregistry.NotFound {
		t.Errorf("expected noTypesResolver.FindMessageByName to return NotFound, got %v", err)
	}
	_, err = noTypesResolver.FindMessageByURL("type.googleapis.com/google.protobuf.Any")
	if err != protoregistry.NotFound {
		t.Errorf("expected noTypesResolver.FindMessageByURL to return NotFound, got %v", err)
	}

	noTypesThenDummyTypesResolver := NewGreedyResolver([]Resolver{noTypes, dummyTypes})
	_, err = noTypesThenDummyTypesResolver.FindMessageByName("google.protobuf.Any")
	if err != nil {
		t.Errorf("expected noTypesThenDummyTypesResolver.FindMessageByName to succeed, got %v", err)
	}
	_, err = noTypesThenDummyTypesResolver.FindMessageByURL("type.googleapis.com/google.protobuf.Any")
	if err != nil {
		t.Errorf("expected noTypesThenDummyTypesResolver.FindMessageByURL to succeed, got %v", err)
	}
}

func mustFindExtension(t *testing.T) protoreflect.ExtensionType {
	t.Helper()
	var foundExt protoreflect.ExtensionType
	protoregistry.GlobalTypes.RangeExtensionsByMessage("google.protobuf.FieldOptions", func(ext protoreflect.ExtensionType) bool {
		foundExt = ext
		return false // stop iteration
	})

	if foundExt == nil {
		t.Fatal("No extensions found in protoregistry.GlobalTypes")
	}
	return foundExt
}

func TestGreedyResolver_FindExtensionByName(t *testing.T) {
	noTypes := new(protoregistry.Types)
	r := NewGreedyResolver([]Resolver{noTypes, protoregistry.GlobalTypes})

	foundExt := mustFindExtension(t)
	extDesc := foundExt.TypeDescriptor()

	tests := []struct {
		name      string
		fieldName protoreflect.FullName
		wantExt   protoreflect.ExtensionType
		wantErr   error
	}{
		{
			name:      "Found",
			fieldName: extDesc.FullName(),
			wantExt:   foundExt,
			wantErr:   nil,
		},
		{
			name:      "NotFound",
			fieldName: "unknown.extension",
			wantExt:   nil,
			wantErr:   protoregistry.NotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotExt, err := r.FindExtensionByName(tc.fieldName)
			if err != tc.wantErr {
				t.Fatalf("FindExtensionByName(%q) got error %v, want %v", tc.fieldName, err, tc.wantErr)
			}
			if gotExt != tc.wantExt {
				t.Errorf("FindExtensionByName(%q) got extension %v, want %v", tc.fieldName, gotExt, tc.wantExt)
			}
		})
	}
}

func TestGreedyResolver_FindExtensionByNumber(t *testing.T) {
	noTypes := new(protoregistry.Types)
	r := NewGreedyResolver([]Resolver{noTypes, protoregistry.GlobalTypes})

	foundExt := mustFindExtension(t)
	extDesc := foundExt.TypeDescriptor()

	tests := []struct {
		name        string
		messageName protoreflect.FullName
		fieldNumber protoreflect.FieldNumber
		wantExt     protoreflect.ExtensionType
		wantErr     error
	}{
		{
			name:        "Found",
			messageName: extDesc.ContainingMessage().FullName(),
			fieldNumber: extDesc.Number(),
			wantExt:     foundExt,
			wantErr:     nil,
		},
		{
			name:        "NotFound",
			messageName: extDesc.ContainingMessage().FullName(),
			fieldNumber: extDesc.Number() + 1,
			wantExt:     nil,
			wantErr:     protoregistry.NotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotExt, err := r.FindExtensionByNumber(tc.messageName, tc.fieldNumber)
			if err != tc.wantErr {
				t.Fatalf("FindExtensionByNumber(%q, %d) got error %v, want %v", tc.messageName, tc.fieldNumber, err, tc.wantErr)
			}
			if gotExt != tc.wantExt {
				t.Errorf("FindExtensionByNumber(%q, %d) got extension %v, want %v", tc.messageName, tc.fieldNumber, gotExt, tc.wantExt)
			}
		})
	}
}
