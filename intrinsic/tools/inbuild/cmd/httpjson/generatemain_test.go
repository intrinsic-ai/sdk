// Copyright 2023 Intrinsic Innovation LLC

package generatemain

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"intrinsic/util/testing/testio"
)

const (
	openapiPath = "intrinsic/tools/inbuild/cmd/httpjson/test_data/openapi.yaml"
)

func TestGenerate(t *testing.T) {
	// Setup temporary environment
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "asdf.go")

	// Reset global variables used for CLI flags
	resetGenerateCommand()

	args := []string{
		"--service_go_importpath", "foo/bar/service_go_proto",
		"--grpc_service", "foo.bar.Baz",
		"--openapi_path", testio.MustCreateRunfilePath(t, openapiPath),
		"--output", outputPath,
	}

	GenerateMainCmd.SetArgs(args)

	if err := GenerateMainCmd.Execute(); err != nil {
		t.Fatalf("GenerateMainCmd.Execute() failed: %v", err)
	}

	// Check main.go content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) failed: %v", outputPath, err)
	}
	wantMainLine := `	pb "foo/bar/service_go_proto"`
	if !strings.Contains(string(content), wantMainLine) {
		t.Errorf("main.go missing expected line: %s. Have lines: %v", wantMainLine, string(content))
	}
}

func TestGenerateErrors(t *testing.T) {
	// Helper to create string pointers easily
	ptr := func(s string) *string { return &s }
	tmpDir := t.TempDir()

	type flags struct {
		serviceGoImportpath *string
		grpcService         *string
		openapiPath         *string
		outputPath          *string
	}

	// The "Golden" state: all flags are valid
	knownGood := flags{
		serviceGoImportpath: ptr("foo/bar/service_go_proto"),
		grpcService:         ptr("foo.bar.Baz"),
		openapiPath:         ptr(testio.MustCreateRunfilePath(t, openapiPath)),
		outputPath:          ptr(filepath.Join(tmpDir, "asdf.go")),
	}

	tests := []struct {
		name            string
		modify          func(f *flags)
		wantErrContains string
	}{
		{
			name: "Error when service_go_importpath is missing",
			modify: func(f *flags) {
				f.serviceGoImportpath = nil
			},
			wantErrContains: "--service_go_importpath is required",
		},
		{
			name: "Error when grpc_service is missing",
			modify: func(f *flags) {
				f.grpcService = nil
			},
			wantErrContains: "--grpc_service is required",
		},
		{
			name: "Error when openapi_path is missing",
			modify: func(f *flags) {
				f.openapiPath = nil
			},
			wantErrContains: "--openapi_path is required",
		},
		{
			name: "Error when openapi_path does not exist",
			modify: func(f *flags) {
				f.openapiPath = ptr("path/to/nowhere.yaml")
			},
			wantErrContains: "path/to/nowhere.yaml",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetGenerateCommand()

			// Start with a fresh copy of known good flags
			f := knownGood
			// Apply the modification for this specific test case
			tc.modify(&f)

			// Construct the args slice based on non-nil pointers
			var args []string
			if f.serviceGoImportpath != nil {
				args = append(args, "--service_go_importpath", *f.serviceGoImportpath)
			}
			if f.grpcService != nil {
				args = append(args, "--grpc_service", *f.grpcService)
			}
			if f.openapiPath != nil {
				args = append(args, "--openapi_path", *f.openapiPath)
			}
			if f.outputPath != nil {
				args = append(args, "--output", *f.outputPath)
			}

			GenerateMainCmd.SetArgs(args)
			err := GenerateMainCmd.Execute()

			if err == nil {
				t.Errorf("expected error containing %q, but got nil", tc.wantErrContains)
				return
			}

			if !strings.Contains(err.Error(), tc.wantErrContains) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.wantErrContains)
			}
		})
	}
}
