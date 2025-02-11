// Copyright 2023 Intrinsic Innovation LLC

// Package skillmanifest has utilities used by inbuild to work with skill manifests
package skillmanifest

import (
	"fmt"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"intrinsic/skills/internal/skillmanifest"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"
)

// LoadManifestAndFileDescriptorSets loads a skill manifest and consolidates multiple file descriptor sets into one.
// If the file descriptor sets have source code info, then it is stripped for all types not used by
// the skill manifest.
func LoadManifestAndFileDescriptorSets(manifestPath string, fdsPaths []string) (*smpb.SkillManifest, *dpb.FileDescriptorSet, error) {
	fds, err := registryutil.LoadFileDescriptorSets(fdsPaths)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to build FileDescriptorSet: %v", err)
	}
	types, err := registryutil.NewTypesFromFileDescriptorSet(fds)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to populate the registry: %v", err)
	}
	m := new(smpb.SkillManifest)
	if err := protoio.ReadTextProto(manifestPath, m, protoio.WithResolver(types)); err != nil {
		return nil, nil, fmt.Errorf("failed to read manifest: %v", err)
	}
	if err := skillmanifest.ValidateManifest(m, types); err != nil {
		return nil, nil, fmt.Errorf("failed to validate manifest: %v", err)
	}
	skillmanifest.PruneSourceCodeInfo(m, fds)

	return m, fds, nil
}
