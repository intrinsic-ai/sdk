// Copyright 2023 Intrinsic Innovation LLC

package skillmanifest

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/proto"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"
	"intrinsic/util/testing/testio"

	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
)

const (
	noOpManifestFilename        = "intrinsic/skills/build_defs/tests/no_op_skill_cc_manifest.pbbin"
	noOpDescriptorFilename      = "intrinsic/skills/build_defs/tests/no_op_skill_cc_manifest_filedescriptor.pbbin"
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
		testio.MustCreateRunfilePath(t, noOpDescriptorFilename),
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
		wantError error
	}{
		{
			name:     "C++ no op",
			manifest: mustLoadManifest(t, noOpManifestFilename),
		},
		{
			name: "C++ no op with types",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, noOpManifestFilename)).(*smpb.SkillManifest)
				m.Parameter.MessageFullName = "intrinsic_proto.skills.NoOpSkillParams"
				return m
			}(),
			opts: []ValidateSkillManifestOption{WithTypes(types)},
		},
		{
			name: "C++ no op invalid name",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, noOpManifestFilename)).(*smpb.SkillManifest)
				m.Id.Name = ""
				return m
			}(),
			wantError: cmpopts.AnyError,
		},
		{
			name: "C++ no op name too long",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, noOpManifestFilename)).(*smpb.SkillManifest)
				m.Id.Name = strings.Repeat("a", 1024)
				return m
			}(),
			wantError: cmpopts.AnyError,
		},
		{
			name: "C++ no op invalid display name",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, noOpManifestFilename)).(*smpb.SkillManifest)
				m.DisplayName = ""
				return m
			}(),
			wantError: cmpopts.AnyError,
		},
		{
			name: "C++ no op display name too long",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, noOpManifestFilename)).(*smpb.SkillManifest)
				m.DisplayName = strings.Repeat("a", 1024)
				return m
			}(),
			wantError: cmpopts.AnyError,
		},
		{
			name: "C++ no op description too long",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, noOpManifestFilename)).(*smpb.SkillManifest)
				m.Documentation.Description = strings.Repeat("a", 4096)
				return m
			}(),
			wantError: cmpopts.AnyError,
		},
		{
			name: "C++ no op invalid parameter type",
			manifest: func() *smpb.SkillManifest {
				m := proto.Clone(mustLoadManifest(t, noOpManifestFilename)).(*smpb.SkillManifest)
				m.Parameter.MessageFullName = "invalid.type"
				return m
			}(),
			opts:      []ValidateSkillManifestOption{WithTypes(types)},
			wantError: cmpopts.AnyError,
		},
	}

	for _, tc := range tests {
		err := ValidateSkillManifest(tc.manifest, tc.opts...)
		if diff := cmp.Diff(tc.wantError, err, cmpopts.EquateErrors()); diff != "" {
			t.Errorf("ValidateSkillManifest(%v) returned unexpected error (-want +got):\n%s", tc.manifest, diff)
		}
	}
}

func TestPruneSourceCodeInfo(t *testing.T) {
	fds, err := registryutil.LoadFileDescriptorSets([]string{
		testio.MustCreateRunfilePath(t, noOpDescriptorFilename),
	})
	if err != nil {
		t.Fatalf("unable to build FileDescriptorSet: %v", err)
	}
	m := mustLoadManifest(t, noOpManifestFilename)

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
