// Copyright 2023 Intrinsic Innovation LLC

package protoregistryresolver

import (
	"testing"

	"intrinsic/httpjson/any/fakeserver"

	"google.golang.org/protobuf/reflect/protoregistry"
)

const anyTypeUrl = "type.intrinsic.ai/google.protobuf.Any"

func TestFindMessageByURL(t *testing.T) {
	fakeServerURL := fakeserver.MustMakeFakeServer(t)

	resolver, err := NewProtoRegistryResolver(fakeServerURL)
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
			testName:    "DoesNotExist",
			typeUrl:     "type.intrinsic.ai/does.not.exist",
			wantMessage: false,
			wantErr:     protoregistry.NotFound,
		},
		{
			testName:    "UnsupportedTypeURL",
			typeUrl:     "does.not.exist/foo.bar",
			wantMessage: false,
			wantErr:     protoregistry.NotFound,
		},
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
