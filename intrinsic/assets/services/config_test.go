// Copyright 2023 Intrinsic Innovation LLC

package config

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"intrinsic/util/proto/protoio"

	rcpb "intrinsic/resources/proto/runtime_context_go_proto"
)

func TestLoadRuntimeContext(t *testing.T) {
	want := &rcpb.RuntimeContext{
		Name: "test-service",
	}

	runtimeContextPath = fmt.Sprintf("%s/runtime_config.pb", t.TempDir())
	if err := protoio.WriteBinaryProto(runtimeContextPath, want); err != nil {
		t.Fatalf("Failed to write runtime context: %v", err)
	}

	got, err := LoadRuntimeContext()
	if err != nil {
		t.Fatalf("Failed to load runtime context: %v", err)
	}

	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("LoadRuntimeContext() returned diff (-want +got):\n%s", diff)
	}
}
