// Copyright 2023 Intrinsic Innovation LLC

// Package skillmanifest contains tools for working with SkillManifest.
package skillmanifest

import (
	"fmt"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/reflect/protoregistry"
	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatafieldlimits"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
	"intrinsic/util/proto/sourcecodeinfoview"
)

// ValidateManifest checks that a SkillManifest is consistent and valid.
func ValidateManifest(m *smpb.SkillManifest, types *protoregistry.Types) error {
	id, err := idutils.IDFromProto(m.GetId())
	if err != nil {
		return fmt.Errorf("invalid name or package: %v", err)
	}
	if m.GetDisplayName() == "" {
		return fmt.Errorf("missing display name for skill %q", id)
	}
	if m.GetVendor().GetDisplayName() == "" {
		return fmt.Errorf("missing vendor display name")
	}
	if name := m.GetParameter().GetMessageFullName(); name != "" {
		if _, err := types.FindMessageByURL(name); err != nil {
			return fmt.Errorf("problem with parameter message name %q: %w", name, err)
		}
	}
	if name := m.GetReturnType().GetMessageFullName(); name != "" {
		if _, err := types.FindMessageByURL(name); err != nil {
			return fmt.Errorf("problem with return message name %q: %w", name, err)
		}
	}
	if err := metadatafieldlimits.ValidateNameLength(m.GetId().GetName(), atypepb.AssetType_ASSET_TYPE_SKILL); err != nil {
		return fmt.Errorf("invalid name for skill: %v", err)
	}
	if err := metadatafieldlimits.ValidateDescriptionLength(m.GetDocumentation().GetDescription()); err != nil {
		return fmt.Errorf("invalid description for skill: %v", err)
	}
	if err := metadatafieldlimits.ValidateDisplayNameLength(m.GetDisplayName()); err != nil {
		return fmt.Errorf("invalid display name for skill: %v", err)
	}
	return nil
}

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
