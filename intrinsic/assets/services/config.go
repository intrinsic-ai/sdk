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

// Package config provides utilities for working with service configurations.
package config

import (
	"fmt"

	"intrinsic/util/proto/protoio"

	rcpb "intrinsic/resources/proto/runtime_context_go_proto"
)

// runtimeContextPath is a var so we can override it in tests.
var runtimeContextPath = "/etc/intrinsic/runtime_config.pb"

// LoadRuntimeContext loads the Service's runtime context.
func LoadRuntimeContext() (*rcpb.RuntimeContext, error) {
	rc := &rcpb.RuntimeContext{}
	if err := protoio.ReadBinaryProto(runtimeContextPath, rc); err != nil {
		return nil, fmt.Errorf("failed to load runtime context: %w", err)
	}
	return rc, nil
}
