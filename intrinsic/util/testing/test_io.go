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

// Package testio provides helper functions for handling files in tests.
package testio

import (
	"os"
	"path/filepath"
	"testing"

	"intrinsic/util/proto/protoio"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// MustCreateParentDirectory creates the full file path to a specified file
// name.
func MustCreateParentDirectory(t *testing.T, path string) {
	t.Helper()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("Unable to create %q: %v", dir, err)
	}
}

// MustCreateFile creates a file at the given path with the specified content.
func MustCreateFile(t *testing.T, content []byte, path string) {
	t.Helper()
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("Write %q failed: %v", path, err)
	}
}

// MustCreateBinaryProto creates a serialized binary proto file at the given
// path.
func MustCreateBinaryProto(t *testing.T, p proto.Message, path string) {
	t.Helper()
	b, err := proto.Marshal(p)
	if err != nil {
		t.Fatalf("Failed to marshal proto: %v", err)
	}
	MustCreateFile(t, b, path)
}

// MustCreateTextProto creates a textproto file from a given proto message.
// The output is not stable, so should not be relied upon for matching exactly.
// This should only be used to create artifacts that are immediately parsed by
// tests.
func MustCreateTextProto(t *testing.T, p proto.Message, path string) {
	t.Helper()
	b, err := prototext.Marshal(p)
	if err != nil {
		t.Fatalf("Failed to marshal proto: %v", err)
	}
	MustCreateFile(t, b, path)
}

// MustReadTextProto reads a textproto file into the given proto message.
func MustReadTextProto(t *testing.T, path string, p proto.Message, opts ...protoio.TextReadOption) {
	t.Helper()
	if err := protoio.ReadTextProto(path, p, opts...); err != nil {
		t.Fatalf("Failed to read text proto: %v", err)
	}
}

// MustResolveRlocation resolves the given rlocation path, failing the test if
// the runfile cannot be resolved.
func MustResolveRlocation(t *testing.T, path string) string {
	t.Helper()
	// Use CallerRepository so apparent repository names in `path` are resolved
	// relative to the caller's Bazel module rather than testio's module.
	resolved, err := runfiles.RlocationFrom(path, runfiles.CallerRepository())
	if err != nil {
		t.Fatalf("runfiles.RlocationFrom(%q) failed: %v", path, err)
	}
	return resolved
}
