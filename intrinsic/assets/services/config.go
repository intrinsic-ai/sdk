// Copyright 2023 Intrinsic Innovation LLC

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
