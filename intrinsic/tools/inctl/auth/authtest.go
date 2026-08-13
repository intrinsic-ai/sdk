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

// Package authtest provides test helpers.
package authtest

import (
	"testing"

	"intrinsic/tools/inctl/auth/auth"
)

// NewStoreForTest creates a new auth.Store for use in tests.
func NewStoreForTest(t *testing.T) *auth.Store {
	configDir := t.TempDir()
	return &auth.Store{GetConfigDirFx: func() (string, error) { return configDir, nil }}
}
