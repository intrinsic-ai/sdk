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

func TestInstalledAssetsResolver_FindMessageByURL(t *testing.T) {
	fakeServerURL := MustMakeFakeServer(t)

	resolver, err := NewInstalledAssetsResolver(fakeServerURL)
	if err != nil {
		t.Fatalf("Failed to create resolver: %v", err)
	}

	t.Run("BeforeRefresh", func(t *testing.T) {
		mt, err := resolver.FindMessageByURL("type.googleapis.com/google.protobuf.Any")
		if mt != nil {
			t.Errorf("Expected nil message type, got %v", mt)
		}
		if err != protoregistry.NotFound {
			t.Errorf("Expected protoregistry.NotFound, got %v", err)
		}
	})

	if err := resolver.RefreshInstalledAssets(); err != nil {
		t.Fatalf("RefreshInstalledAssets failed: %v", err)
	}

	t.Run("AfterRefresh", func(t *testing.T) {
		mt, err := resolver.FindMessageByURL("type.googleapis.com/google.protobuf.Any")
		if mt == nil {
			t.Fatal("Expected message type, got nil")
		}
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if mt.Descriptor().FullName() != "google.protobuf.Any" {
			t.Errorf("Expected google.protobuf.Any, got %v", mt.Descriptor().FullName())
		}
	})
}

func TestInstalledAssetsResolver_FindExtension(t *testing.T) {
	resolver, err := NewInstalledAssetsResolver("localhost:0")
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
