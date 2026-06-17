// Copyright 2023 Intrinsic Innovation LLC

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
