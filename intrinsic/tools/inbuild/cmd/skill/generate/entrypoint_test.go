// Copyright 2023 Intrinsic Innovation LLC

package entrypoint

import (
	"os"
	"strings"
	"testing"

	"intrinsic/util/testing/testio"
)

const (
	examplePyManifestPath = "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_py.manifest.pbtxt"
	exampleCCManifestPath = "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_cc.manifest.pbtxt"
	defaultPythonFile     = "main.py"
	defaultCppFile        = "main.cc"
)

func TestEntryPoint(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantFile string
	}{
		{
			name: "cpp skill entry point",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, exampleCCManifestPath),
				"--language", "cpp",
				"--cc_header", "foobar.h",
			},
			wantFile: defaultCppFile,
		},
		{
			name: "python skill entry point",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, examplePyManifestPath),
				"--language", "python",
			},
			wantFile: defaultPythonFile,
		},
		{
			name: "C++ custom output",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, exampleCCManifestPath),
				"--language", "cpp",
				"--cc_header", "foobar.h",
				"--output", "foobar.cc",
			},
			wantFile: "foobar.cc",
		},
		{
			name: "python custom output",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, examplePyManifestPath),
				"--language", "python",
				"--output", "foobar.py",
			},
			wantFile: "foobar.py",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make sure we're in a writable directory.
			os.Chdir(t.TempDir())
			// Prevent state leaking between tests.
			resetEntryPointCommand()

			EntryPointCmd.SetArgs(tt.args)

			if err := EntryPointCmd.Execute(); err != nil {
				t.Fatalf("EntryPointCmd.Execute() failed: %v", err)
			}

			gotFile, err := os.ReadFile(tt.wantFile)
			if err != nil {
				t.Fatalf("os.ReadFile(%q) failed: %v", tt.wantFile, err)
			}
			if len(gotFile) == 0 {
				t.Errorf("gotFile is empty")
			}
		})
	}
}

func TestEntryPointErrors(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		wantErrorContains string
	}{
		{
			name: "invalid languagee",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, exampleCCManifestPath),
				"--language", "ada",
			},
			wantErrorContains: "--language must be one of",
		},
		{
			name: "no manifest",
			args: []string{
				"--language", "python",
			},
			wantErrorContains: "--manifest is required",
		},
		{
			name: "cpp but no cc_header",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, exampleCCManifestPath),
				"--language", "cpp",
			},
			wantErrorContains: "--cc_header is required",
		},
		{
			name: "python with cc_header",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, examplePyManifestPath),
				"--language", "python",
				"--cc_header", "foobar.h",
			},
			wantErrorContains: "--cc_header is only supported when",
		},
		{
			name: "No manifest",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, "does_not_exist.pbtxt"),
				"--language", "python",
			},
			wantErrorContains: "unable to read manifest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make sure we're in a writable directory.
			os.Chdir(t.TempDir())
			// Prevent state leaking between tests.
			resetEntryPointCommand()

			EntryPointCmd.SetArgs(tt.args)
			err := EntryPointCmd.Execute()
			if err == nil {
				t.Fatalf("EntryPointCmd.Execute() did not fail: %v", tt.args)
			}
			if !strings.Contains(err.Error(), tt.wantErrorContains) {
				t.Errorf("EntryPointCmd.Execute() returned error: %v, want error containing: %v", err, tt.wantErrorContains)
			}
		})
	}
}
