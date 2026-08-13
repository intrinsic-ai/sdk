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

package openapi

import (
	"os"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

func TestOpenAPIFile(t *testing.T) {
	r, err := runfiles.New()
	if err != nil {
		t.Fatalf("Failed to initialize runfiles: %v", err)
	}

	openapiYaml, err := r.Rlocation("ai_intrinsic_sdks/intrinsic/httpjson/test/_inventory_service_openapi/openapi.yaml")
	if err != nil {
		t.Fatalf("failed to find openapi.yaml: %v", err)
	}
	content, err := os.ReadFile(openapiYaml)
	if err != nil {
		t.Fatalf("failed to read openapi.yaml: %v", err)
	}
	if !strings.Contains(string(content), "/v1/skus") {
		t.Errorf("openapi.yaml does not contain '/v1/skus'")
	}
}
