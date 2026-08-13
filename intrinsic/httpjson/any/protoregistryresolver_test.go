// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package any

import (
	"errors"
	"testing"

	"google.golang.org/protobuf/reflect/protoregistry"
)

func TestProtoRegistryResolver_FindMessageByURL(t *testing.T) {
	fakeServerURL := MustMakeFakeServer(t)

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

func TestProtoRegistryResolver_FindExtension(t *testing.T) {
	resolver, err := NewProtoRegistryResolver("localhost:0")
	if err != nil {
		t.Fatalf("Failed to create resolver: %v", err)
	}

	if ext, err := resolver.FindExtensionByName("some.extension"); ext != nil || !errors.Is(err, protoregistry.NotFound) {
		t.Errorf("FindExtensionByName: expected nil, NotFound; got %v, %v", ext, err)
	}

	if ext, err := resolver.FindExtensionByNumber("some.message", 42); ext != nil || !errors.Is(err, protoregistry.NotFound) {
		t.Errorf("FindExtensionByNumber: expected nil, NotFound; got %v, %v", ext, err)
	}
}
