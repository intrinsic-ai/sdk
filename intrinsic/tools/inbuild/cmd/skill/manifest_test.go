// Copyright 2023 Intrinsic Innovation LLC

package manifest

import (
	"os"
	"testing"

	"intrinsic/util/testing/testio"
)

const (
	exampleManifestPath          = "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_py.manifest.pbtxt"
	exampleFileDescriptorSetPath = "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_proto-descriptor-set.proto.bin"
)

func TestManifestCreate(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "create manifest",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, exampleManifestPath),
				"--file_descriptor_sets", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
				"--output", "skill.manifest.pb",
				"--file_descriptor_set_out", "fds.pb",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make sure we're in a writable directory.
			os.Chdir(t.TempDir())
			// Prevent state leaking between tests.
			resetManifestCommand()

			ManifestCmd.SetArgs(tt.args)

			if err := ManifestCmd.Execute(); err != nil {
				t.Fatalf("ManifestCmd.Execute() failed: %v", err)
			}

			if _, err := os.Stat("skill.manifest.pb"); os.IsNotExist(err) {
				t.Errorf("skill.manifest.pb was not created")
			}
			if _, err := os.Stat("fds.pb"); os.IsNotExist(err) {
				t.Errorf("fds.pb was not created")
			}
		})
	}
}
