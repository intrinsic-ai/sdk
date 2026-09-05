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

// Package skillmanifest has utilities used by inbuild to work with skill manifests.
package skillmanifest

import (
	"context"
	"fmt"

	"intrinsic/skills/skillmanifest"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"

	smpb "intrinsic/skills/proto/skill_manifest_go_proto"

	dpb "google.golang.org/protobuf/types/descriptorpb"
)

// LoadManifestAndFileDescriptorSets loads a skill manifest and consolidates multiple file descriptor sets into one.
// If the file descriptor sets have source code info, then it is stripped for all types not used by
// the skill manifest.
func LoadManifestAndFileDescriptorSets(ctx context.Context, manifestPath string, fdsPaths []string, incompatibleDisallowManifestDependencies bool) (*smpb.SkillManifest, *dpb.FileDescriptorSet, error) {
	fds, err := registryutil.LoadFileDescriptorSets(fdsPaths)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to build FileDescriptorSet: %v", err)
	}
	types, err := registryutil.NewTypesFromFileDescriptorSet(fds)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to populate the types registry: %v", err)
	}
	m := new(smpb.SkillManifest)
	// Try loading as binary first.
	if err := protoio.ReadBinaryProto(manifestPath, m); err != nil {
		// If binary fails, try loading as text.
		if err := protoio.ReadTextProto(manifestPath, m, protoio.WithResolver(types)); err != nil {
			return nil, nil, fmt.Errorf("failed to read manifest as binary or text: %v", err)
		}
	}
	if incompatibleDisallowManifestDependencies && len(m.GetDependencies().GetRequiredEquipment()) > 0 {
		return nil, nil, fmt.Errorf("failed to validate manifest: dependencies declared in the manifest's dependencies field but --incompatible_disallow_manifest_dependencies is true")
	}
	if err := skillmanifest.PruneSourceCodeInfo(m, fds); err != nil {
		return nil, nil, fmt.Errorf("failed to prune source code info: %v", err)
	}

	return m, fds, nil
}
