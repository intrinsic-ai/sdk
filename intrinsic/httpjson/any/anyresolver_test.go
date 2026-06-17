// Copyright 2023 Intrinsic Innovation LLC

package any

import (
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func TestImplementsResolver(t *testing.T) {
	mo := protojson.MarshalOptions{}
	mo.Resolver = (*AnyResolver)(nil)
}

func TestAnyResolver_FindMessageByURL(t *testing.T) {
	fakeServerURL := MustMakeFakeServer(t)

	resolver, err := NewAnyResolver(fakeServerURL, fakeServerURL)
	if err != nil {
		t.Fatalf("Failed to create resolver: %v", err)
	}

	tests := []struct {
		testName    string
		typeUrl     string
		wantMessage bool
		wantErr     error
	}{
		{
			testName:    "DoesExist",
			typeUrl:     anyTypeUrl,
			wantMessage: true,
			wantErr:     nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mt, err := resolver.FindMessageByURL(test.typeUrl)
			if (mt != nil) != test.wantMessage {
				t.Errorf("FindMessageByURL(%q) got message type %v, want message type: %v", test.typeUrl, mt != nil, test.wantMessage)
			}
			if err != test.wantErr {
				t.Errorf("FindMessageByURL(%q) got error %v, want error %v", test.typeUrl, err, test.wantErr)
			}
		})
	}
}

func TestAnyResolver_FindExtensionByName(t *testing.T) {
	resolver, err := NewAnyResolver("localhost:0", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create resolver: %v", err)
	}

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
			fieldName: "some.extension",
			wantExt:   nil,
			wantErr:   protoregistry.NotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotExt, err := resolver.FindExtensionByName(tc.fieldName)
			if err != tc.wantErr {
				t.Fatalf("FindExtensionByName(%q) got error %v, want %v", tc.fieldName, err, tc.wantErr)
			}
			if gotExt != tc.wantExt {
				t.Errorf("FindExtensionByName(%q) got extension %v, want %v", tc.fieldName, gotExt, tc.wantExt)
			}
		})
	}
}

func TestAnyResolver_FindExtensionByNumber(t *testing.T) {
	resolver, err := NewAnyResolver("localhost:0", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create resolver: %v", err)
	}

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
			gotExt, err := resolver.FindExtensionByNumber(tc.messageName, tc.fieldNumber)
			if err != tc.wantErr {
				t.Fatalf("FindExtensionByNumber(%q, %d) got error %v, want %v", tc.messageName, tc.fieldNumber, err, tc.wantErr)
			}
			if gotExt != tc.wantExt {
				t.Errorf("FindExtensionByNumber(%q, %d) got extension %v, want %v", tc.messageName, tc.fieldNumber, gotExt, tc.wantExt)
			}
		})
	}
}
