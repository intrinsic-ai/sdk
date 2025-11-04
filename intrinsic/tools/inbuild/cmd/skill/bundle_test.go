// Copyright 2023 Intrinsic Innovation LLC

package bundle

import (
	"os"
	"testing"

	"intrinsic/assets/bundleio"
	"intrinsic/util/testing/testio"

	_ "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_go_proto" // Needed to resolve the proto message
)

const (
	exampleManifestPath          = "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_py.manifest.pbtxt"
	exampleFileDescriptorSetPath = "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_proto-descriptor-set.proto.bin"
	exampleImagePath             = "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_py_image.tar"
)

type bundleCheck func(t *testing.T)

func checkManifestHasID(t *testing.T, bundlePath string, wantPackage string, wantName string) bundleCheck {
	return func(t *testing.T) {
		t.Helper()
		manifest, err := bundleio.ReadSkillManifest(t.Context(), bundlePath)
		if err != nil {
			t.Fatalf("bundleio.ReadSkillManifest(%q) failed: %v", bundlePath, err)
		}
		if got := manifest.GetId().GetPackage(); got != wantPackage {
			t.Errorf("manifest.GetId().GetPackage() = %q, want %q", got, wantPackage)
		}
		if got := manifest.GetId().GetName(); got != wantName {
			t.Errorf("manifest.GetId().GetName() = %q, want %q", got, wantName)
		}
	}
}

func checkHasFile(t *testing.T, bundlePath string, wantFile string) bundleCheck {
	return func(t *testing.T) {
		t.Helper()
		_, gotContents, err := bundleio.ReadSkill(t.Context(), bundlePath)
		if err != nil {
			t.Fatalf("bundleio.ReadSkill(%q) failed: %v", bundlePath, err)
		}
		_, ok := gotContents[wantFile]
		if !ok {
			t.Errorf("gotContents[%q] = nil, want non-nil", wantFile)
		}
	}
}

func TestBundleCreate(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		checks []bundleCheck
	}{
		{
			name: "create skill",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, exampleManifestPath),
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
				"--oci_image", testio.MustCreateRunfilePath(t, exampleImagePath),
			},
			checks: []bundleCheck{
				checkManifestHasID(t, "skill.bundle.tar", "com.example", "example_skill"),
				checkHasFile(t, "skill.bundle.tar", "example_skill_py_image.tar"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make sure we're in a writable directory.
			os.Chdir(t.TempDir())
			// Prevent state leaking between tests.
			resetBundleCommand()

			BundleCmd.SetArgs(tt.args)

			if err := BundleCmd.Execute(); err != nil {
				t.Fatalf("BundleCmd.Execute() failed: %v", err)
			}

			for _, check := range tt.checks {
				check(t)
			}
		})
	}
}
