// Copyright 2023 Intrinsic Innovation LLC

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
