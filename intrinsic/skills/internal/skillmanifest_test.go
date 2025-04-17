// Copyright 2023 Intrinsic Innovation LLC

package skillmanifest

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"
	"intrinsic/util/testing/testio"
)

const (
	manifestFilename   = "intrinsic/skills/build_defs/tests/no_op_skill_cc_manifest.pbbin"
	descriptorFilename = "intrinsic/skills/build_defs/tests/no_op_skill_cc_manifest_filedescriptor.pbbin"
)

func mustLoadManifest(t *testing.T, path string) *smpb.SkillManifest {
	t.Helper()
	realPath := testio.MustCreateRunfilePath(t, path)
	m := new(smpb.SkillManifest)
	if err := protoio.ReadBinaryProto(realPath, m); err != nil {
		t.Fatalf("failed to read manifest: %v", err)
	}
	return m
}

func TestValidateManifest(t *testing.T) {
	set, err := registryutil.LoadFileDescriptorSets([]string{
		testio.MustCreateRunfilePath(t, descriptorFilename),
	})
	if err != nil {
		t.Fatalf("unable to build FileDescriptorSet: %v", err)
	}
	types, err := registryutil.NewTypesFromFileDescriptorSet(set)
	if err != nil {
		t.Fatalf("failed to populate the registry: %v", err)
	}

	tests := []struct {
		name      string
		manifest  *smpb.SkillManifest
		opts      []ValidateSkillManifestOption
		wantError bool
	}{
		{
			name:      "C++ no op",
			manifest:  mustLoadManifest(t, manifestFilename),
			wantError: false,
		},
		{
			name: "C++ no op with types",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, manifestFilename)).(*smpb.SkillManifest)
				m.Parameter.MessageFullName = "intrinsic_proto.skills.NoOpSkillParams"
				return m
			}(),
			opts:      []ValidateSkillManifestOption{WithTypes(types)},
			wantError: false,
		},
		{
			name: "C++ no op invalid name",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, manifestFilename)).(*smpb.SkillManifest)
				m.Id.Name = ""
				return m
			}(),
			wantError: true,
		},
		{
			name: "C++ no op name too long",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, manifestFilename)).(*smpb.SkillManifest)
				m.Id.Name = strings.Repeat("a", 1024)
				return m
			}(),
			wantError: true,
		},
		{
			name: "C++ no op invalid display name",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, manifestFilename)).(*smpb.SkillManifest)
				m.DisplayName = ""
				return m
			}(),
			wantError: true,
		},
		{
			name: "C++ no op display name too long",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, manifestFilename)).(*smpb.SkillManifest)
				m.DisplayName = strings.Repeat("a", 1024)
				return m
			}(),
			wantError: true,
		},
		{
			name: "C++ no op description too long",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, manifestFilename)).(*smpb.SkillManifest)
				m.Documentation.Description = strings.Repeat("a", 4096)
				return m
			}(),
			wantError: true,
		},
		{
			name: "C++ no op invalid parameter type",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, manifestFilename)).(*smpb.SkillManifest)
				m.Parameter.MessageFullName = "invalid.type"
				return m
			}(),
			opts:      []ValidateSkillManifestOption{WithTypes(types)},
			wantError: true,
		},
	}

	for _, tc := range tests {
		err := ValidateSkillManifest(tc.manifest, tc.opts...)
		if gotError := (err != nil); tc.wantError != gotError {
			t.Fatalf("wantErr: %v gotError: %v err: %v", tc.wantError, gotError, err)
		}
	}
}

func TestPruneSourceCodeInfo(t *testing.T) {
	fds, err := registryutil.LoadFileDescriptorSets([]string{
		testio.MustCreateRunfilePath(t, descriptorFilename),
	})
	if err != nil {
		t.Fatalf("unable to build FileDescriptorSet: %v", err)
	}
	m := mustLoadManifest(t, manifestFilename)

	PruneSourceCodeInfo(m, fds)
	for _, file := range fds.GetFile() {
		if strings.HasSuffix(file.GetName(), "no_op_skill.proto") {
			if file.GetSourceCodeInfo() == nil {
				t.Fatalf("%v has no source code info, but it should not have been pruned", file.GetName())
			}
		} else if file.GetSourceCodeInfo() != nil {
			t.Fatalf("%v has source code info, but it should have been pruned", file.GetName())
		}
	}
}
