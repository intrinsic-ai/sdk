// Copyright 2023 Intrinsic Innovation LLC

package anyresolver

import (
	"testing"

	"intrinsic/httpjson/any/fakeserver"

	"google.golang.org/protobuf/encoding/protojson"
)

const anyTypeUrl = "type.intrinsic.ai/google.protobuf.Any"

func TestImplementsResolver(t *testing.T) {
	mo := protojson.MarshalOptions{}
	mo.Resolver = (*AnyResolver)(nil)
}

func TestFindMessageByURL(t *testing.T) {
	fakeServerURL := fakeserver.MustMakeFakeServer(t)

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
