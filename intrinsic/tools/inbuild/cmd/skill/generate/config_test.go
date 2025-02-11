// Copyright 2023 Intrinsic Innovation LLC

package config

import (
	"os"
	"strings"
	"testing"

	"intrinsic/util/testing/testio"
)

const (
	examplePyManifestPath        = "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_py.manifest.pbtxt"
	exampleCCManifestPath        = "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_cc.manifest.pbtxt"
	exampleFileDescriptorSetPath = "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_proto-descriptor-set.proto.bin"
	defaultOutput                = "config.pbbin"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantFile string
	}{
		{
			name: "cpp skill service",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, exampleCCManifestPath),
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
			},
			wantFile: defaultOutput,
		},
		{
			name: "python skill service",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, examplePyManifestPath),
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
			},
			wantFile: defaultOutput,
		},
		{
			name: "cpp skill service custom output",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, exampleCCManifestPath),
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
				"--output", "foobar_cpp.pbbin",
			},
			wantFile: "foobar_cpp.pbbin",
		},
		{
			name: "python skill service custom output",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, examplePyManifestPath),
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
				"--output", "foobar_py.pbbin",
			},
			wantFile: "foobar_py.pbbin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make sure we're in a writable directory.
			os.Chdir(t.TempDir())
			// Prevent state leaking between tests.
			resetConfigCommand()

			ConfigCmd.SetArgs(tt.args)

			if err := ConfigCmd.Execute(); err != nil {
				t.Fatalf("ConfigCmd.Execute() failed: %v", err)
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

func TestConfigErrors(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		wantErrorContains string
	}{
		{
			name: "No manifest",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, "does_not_exist.pbtxt"),
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
			},
			wantErrorContains: "failed to read manifest",
		},
		{
			name: "No file descriptor set",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, exampleCCManifestPath),
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, "does_not_exist.pbbin"),
			},
			wantErrorContains: "failed to read file descriptor set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make sure we're in a writable directory.
			os.Chdir(t.TempDir())
			// Prevent state leaking between tests.
			resetConfigCommand()

			ConfigCmd.SetArgs(tt.args)
			err := ConfigCmd.Execute()
			if err == nil {
				t.Fatalf("ConfigCmd.Execute() did not fail: %v", tt.args)
			}
			if !strings.Contains(err.Error(), tt.wantErrorContains) {
				t.Errorf("ConfigCmd.Execute() returned error: %v, want error containing: %v", err, tt.wantErrorContains)
			}
		})
	}
}
