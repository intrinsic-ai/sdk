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
