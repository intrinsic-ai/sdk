// Copyright 2023 Intrinsic Innovation LLC

package installedassetsresolver

import (
	"testing"

	"google.golang.org/protobuf/reflect/protoregistry"

	"intrinsic/httpjson/any/fakeserver"
)

func TestFindMessageByURL(t *testing.T) {
	fakeServerURL := fakeserver.MustMakeFakeServer(t)

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
