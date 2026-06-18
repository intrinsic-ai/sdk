// Copyright 2023 Intrinsic Innovation LLC

package any

import (
	"errors"
	"strings"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const testDescriptorSetPath = "ai_intrinsic_sdks/intrinsic/httpjson/any/test_messages_descriptors_transitive_set_sci.proto.bin"

func TestFileDescriptorSetResolver_EmptyPaths(t *testing.T) {
	_, err := NewFileDescriptorSetResolver([]string{})
	if err == nil {
		t.Fatal("NewFileDescriptorSetResolver() with empty paths succeeded, expected error")
	}
	expectedSubstr := "no file descriptor set loaded"
	if !strings.Contains(err.Error(), expectedSubstr) {
		t.Errorf("NewFileDescriptorSetResolver() error %q does not contain %q", err.Error(), expectedSubstr)
	}
}

func TestFileDescriptorSetResolver_FindMessageByName(t *testing.T) {
	resolver, err := NewFileDescriptorSetResolver([]string{testDescriptorSetPath})
	if err != nil {
		t.Fatalf("NewFileDescriptorSetResolver() failed: %v", err)
	}

	tests := []struct {
		name        string
		messageName protoreflect.FullName
		wantErr     error
	}{
		{
			name:        "Success",
			messageName: "intrinsic.httpjson.any.test.TestMessage",
			wantErr:     nil,
		},
		{
			name:        "NotFound",
			messageName: "intrinsic.httpjson.any.test.NonExistent",
			wantErr:     protoregistry.NotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mt, err := resolver.FindMessageByName(tc.messageName)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("FindMessageByName(%q) got error %v, want %v", tc.messageName, err, tc.wantErr)
			}
			if tc.wantErr == nil && (mt == nil || mt.Descriptor().FullName() != tc.messageName) {
				t.Errorf("FindMessageByName(%q) got message type %v, want %q", tc.messageName, mt, tc.messageName)
			}
		})
	}
}

func TestFileDescriptorSetResolver_FindMessageByURL(t *testing.T) {
	resolver, err := NewFileDescriptorSetResolver([]string{testDescriptorSetPath})
	if err != nil {
		t.Fatalf("NewFileDescriptorSetResolver() failed: %v", err)
	}

	tests := []struct {
		name    string
		url     string
		wantErr error
	}{
		{
			name:    "Success",
			url:     "type.googleapis.com/intrinsic.httpjson.any.test.TestMessage",
			wantErr: nil,
		},
		{
			name:    "NotFound",
			url:     "type.googleapis.com/intrinsic.httpjson.any.test.NonExistent",
			wantErr: protoregistry.NotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mt, err := resolver.FindMessageByURL(tc.url)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("FindMessageByURL(%q) got error %v, want %v", tc.url, err, tc.wantErr)
			}
			if tc.wantErr == nil && mt == nil {
				t.Errorf("FindMessageByURL(%q) returned nil, want non-nil type", tc.url)
			}
		})
	}
}

func TestFileDescriptorSetResolver_FindExtensionByName(t *testing.T) {
	resolver, err := NewFileDescriptorSetResolver([]string{testDescriptorSetPath})
	if err != nil {
		t.Fatalf("NewFileDescriptorSetResolver() failed: %v", err)
	}

	tests := []struct {
		name      string
		fieldName protoreflect.FullName
		wantErr   error
	}{
		{
			name:      "Success",
			fieldName: "intrinsic.httpjson.any.test.test_extension",
			wantErr:   nil,
		},
		{
			name:      "NotFound",
			fieldName: "some.extension",
			wantErr:   protoregistry.NotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotExt, err := resolver.FindExtensionByName(tc.fieldName)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("FindExtensionByName(%q) got error %v, want %v", tc.fieldName, err, tc.wantErr)
			}
			if tc.wantErr == nil && (gotExt == nil || gotExt.TypeDescriptor().FullName() != tc.fieldName) {
				t.Errorf("FindExtensionByName(%q) got extension %v, want %q", tc.fieldName, gotExt, tc.fieldName)
			}
		})
	}
}

func TestFileDescriptorSetResolver_FindExtensionByNumber(t *testing.T) {
	resolver, err := NewFileDescriptorSetResolver([]string{testDescriptorSetPath})
	if err != nil {
		t.Fatalf("NewFileDescriptorSetResolver() failed: %v", err)
	}

	tests := []struct {
		name        string
		messageName protoreflect.FullName
		fieldNumber protoreflect.FieldNumber
		wantErr     error
	}{
		{
			name:        "Success",
			messageName: "intrinsic.httpjson.any.test.TestMessage",
			fieldNumber: 101,
			wantErr:     nil,
		},
		{
			name:        "NotFound",
			messageName: "intrinsic.httpjson.any.test.TestMessage",
			fieldNumber: 102,
			wantErr:     protoregistry.NotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotExt, err := resolver.FindExtensionByNumber(tc.messageName, tc.fieldNumber)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("FindExtensionByNumber(%q, %d) got error %v, want %v", tc.messageName, tc.fieldNumber, err, tc.wantErr)
			}
			if tc.wantErr == nil && (gotExt == nil || gotExt.TypeDescriptor().Number() != tc.fieldNumber) {
				t.Errorf("FindExtensionByNumber(%q, %d) got extension field number %v, want %d", tc.messageName, tc.fieldNumber, gotExt, tc.fieldNumber)
			}
		})
	}
}
