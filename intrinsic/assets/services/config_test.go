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

package config

import (
	"fmt"
	"testing"

	"intrinsic/util/proto/protoio"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

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
