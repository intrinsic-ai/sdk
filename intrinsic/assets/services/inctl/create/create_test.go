// Copyright 2023 Intrinsic Innovation LLC

package create

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func mustMakeBazelWorkspace(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "bazel_ws_*")
	if err != nil {
		t.Fatal()
	}
	// Create MODULE.bazel
	if err := os.WriteFile(filepath.Join(dir, "MODULE.bazel"), []byte{}, 0644); err != nil {
		t.Fatal()
	}
	return dir
}

func TestCreateCommandErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name: "Invalid language",
			args: []string{
				"--language", "not_a_supported_language",
				"my_org.subpackage.my_service",
			},
			wantErr: "unknown language not_a_supported_language",
		},
		{
			name: "Missing language",
			args: []string{
				"my_org.subpackage.my_service",
			},
			wantErr: "required flag(s) \"language\" not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceRoot := mustMakeBazelWorkspace(t)

			args := append([]string{"--output_path", workspaceRoot}, tt.args...)

			cmd := Command()
			cmd.SetArgs(args)

			err := cmd.Execute()
			if err == nil {
				t.Errorf("No error but expected %s", tt.wantErr)
			} else if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Got error %v but expected %s", err, tt.wantErr)
			}
		})
	}
}

func TestCreatesFiles(t *testing.T) {
	tests := []struct {
		name       string
		language   string
		outputPath string // Relative path to be joined with workspaceRoot
		assetID    string
		wantFiles  []string
	}{
		{
			name:       "Python in workspace root",
			language:   "python",
			outputPath: "",
			assetID:    "com.foo.py_service",
			wantFiles: []string{
				"BUILD",
				"py_service.proto",
				"py_service_defaults.textproto",
				"py_service.py",
				"py_service_manifest.textproto",
			},
		},
		{
			name:       "C++ in nested subdirectory",
			language:   "cpp",
			outputPath: "gen/cpp",
			assetID:    "com.foo.cpp_service",
			wantFiles: []string{
				"gen/cpp/BUILD",
				"gen/cpp/cpp_service.proto",
				"gen/cpp/cpp_service_defaults.textproto",
				"gen/cpp/cpp_service.cc",
				"gen/cpp/cpp_service_manifest.textproto",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceRoot := mustMakeBazelWorkspace(t)

			// Build the command arguments dynamically
			var args []string
			if tt.language != "" {
				args = append(args, "--language", tt.language)
			}

			args = append(args, "--output_path", filepath.Join(workspaceRoot, tt.outputPath))
			args = append(args, tt.assetID)

			cmd := Command()
			cmd.SetArgs(args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute() failed: %v", err)
			}

			for _, f := range tt.wantFiles {
				path := filepath.Join(workspaceRoot, f)

				info, err := os.Stat(path)
				if os.IsNotExist(err) {
					t.Errorf("Expected file %q was not created", f)
					continue
				} else if err != nil {
					t.Errorf("Stat error for %q: %v", f, err)
					continue
				}

				// Assert the file is not empty
				if info.Size() == 0 {
					t.Errorf("File %q is empty, expected content", f)
				}
			}
		})
	}
}
