// Copyright 2023 Intrinsic Innovation LLC

package entrypoint

import (
	"os"
	slice "slices"
	"strings"
	"testing"

	"intrinsic/util/proto/protoio"
	"intrinsic/util/testing/testio"

	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
)

const (
	examplePyManifestPath        = "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_py.manifest.pbtxt"
	exampleCCManifestPath        = "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_cc.manifest.pbtxt"
	exampleFileDescriptorSetPath = "intrinsic/tools/inbuild/cmd/skill/test_data/example_skill_proto-descriptor-set.proto.bin"
	defaultPythonFile            = "main.py"
	defaultCppFile               = "main.cc"
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
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
				"--augmented_manifest_out", "manifest_out.pbbin",
				"--augmented_file_descriptor_set_out", "fds_out.pbbin",
			},
			wantFile: defaultCppFile,
		},
		{
			name: "python skill entry point",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, examplePyManifestPath),
				"--language", "python",
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
				"--augmented_manifest_out", "manifest_out.pbbin",
				"--augmented_file_descriptor_set_out", "fds_out.pbbin",
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
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
				"--augmented_manifest_out", "manifest_out.pbbin",
				"--augmented_file_descriptor_set_out", "fds_out.pbbin",
			},
			wantFile: "foobar.cc",
		},
		{
			name: "python custom output",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, examplePyManifestPath),
				"--language", "python",
				"--output", "foobar.py",
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
				"--augmented_manifest_out", "manifest_out.pbbin",
				"--augmented_file_descriptor_set_out", "fds_out.pbbin",
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

			// Verify generated augmented files exist and contain appropriate service versions.
			manifestOut := new(smpb.SkillManifest)
			if err := protoio.ReadBinaryProto("manifest_out.pbbin", manifestOut); err != nil {
				t.Fatalf("unable to read augmented manifest: %v", err)
			}
			gotVersions := manifestOut.GetOptions().GetSkillServicesConfig().GetServiceVersions()
			wantVersions := []smpb.SkillServicesConfig_ServiceVersion{
				smpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_PROJECTOR,
				smpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_EXECUTOR,
				smpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_SKILL_INFORMATION,
			}
			if !slice.Equal(gotVersions, wantVersions) {
				t.Errorf("manifestOut.GetOptions().GetSkillServicesConfig().GetServiceVersions() = %v, want %v", gotVersions, wantVersions)
			}

			if _, err := os.Stat("fds_out.pbbin"); err != nil {
				t.Errorf("augmented file descriptor set was not created: %v", err)
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
			name: "invalid language",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, exampleCCManifestPath),
				"--language", "ada",
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
				"--augmented_manifest_out", "manifest_out.pbbin",
				"--augmented_file_descriptor_set_out", "fds_out.pbbin",
			},
			wantErrorContains: "--language must be one of",
		},
		{
			name: "no manifest",
			args: []string{
				"--language", "python",
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
				"--augmented_manifest_out", "manifest_out.pbbin",
				"--augmented_file_descriptor_set_out", "fds_out.pbbin",
			},
			wantErrorContains: "--manifest is required",
		},
		{
			name: "no file descriptor set",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, exampleCCManifestPath),
				"--language", "python",
				"--augmented_manifest_out", "manifest_out.pbbin",
				"--augmented_file_descriptor_set_out", "fds_out.pbbin",
			},
			wantErrorContains: "--file_descriptor_set is required",
		},
		{
			name: "cpp but no cc_header",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, exampleCCManifestPath),
				"--language", "cpp",
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
				"--augmented_manifest_out", "manifest_out.pbbin",
				"--augmented_file_descriptor_set_out", "fds_out.pbbin",
			},
			wantErrorContains: "--cc_header is required",
		},
		{
			name: "python with cc_header",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, examplePyManifestPath),
				"--language", "python",
				"--cc_header", "foobar.h",
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
				"--augmented_manifest_out", "manifest_out.pbbin",
				"--augmented_file_descriptor_set_out", "fds_out.pbbin",
			},
			wantErrorContains: "--cc_header is only supported when",
		},
		{
			name: "No manifest",
			args: []string{
				"--manifest", testio.MustCreateRunfilePath(t, "does_not_exist.pbtxt"),
				"--language", "python",
				"--file_descriptor_set", testio.MustCreateRunfilePath(t, exampleFileDescriptorSetPath),
				"--augmented_manifest_out", "manifest_out.pbbin",
				"--augmented_file_descriptor_set_out", "fds_out.pbbin",
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
