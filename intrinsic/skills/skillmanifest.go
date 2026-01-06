// Copyright 2023 Intrinsic Innovation LLC

// Package skillmanifest provides utils for SkillManifests.
package skillmanifest

import (
	"intrinsic/util/proto/sourcecodeinfoview"

	smpb "intrinsic/skills/proto/skill_manifest_go_proto"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// PruneSourceCodeInfo removes source code info from the FileDescriptorSet for all message types
// except those that are referenced by the SkillManifest.
func PruneSourceCodeInfo(m *smpb.SkillManifest, fds *dpb.FileDescriptorSet) {
	var fullNames []string
	if name := m.GetParameter().GetMessageFullName(); name != "" {
		fullNames = append(fullNames, name)
	}
	if name := m.GetReturnType().GetMessageFullName(); name != "" {
		fullNames = append(fullNames, name)
	}
	sourcecodeinfoview.PruneSourceCodeInfo(fullNames, fds)
}
