// Copyright 2023 Intrinsic Innovation LLC

// Package processmanifest provides utils for ProcessManifests.
package processmanifest

import (
	"intrinsic/assets/idutils"

	mpb "intrinsic/assets/proto/metadata_go_proto"
	btpb "intrinsic/executive/proto/behavior_tree_go_proto"
	skpb "intrinsic/skills/proto/skills_go_proto"
)

// FillInSkillIDVersionFromAssetMetadata overwrites in the given Skill proto the `id_version`
// according to the given asset metadata.
//
// Use this method to make the top-level Skill description in a BehaviorTree proto consistent with
// the Asset metadata after changing/setting the Asset's version.
func FillInSkillIDVersionFromAssetMetadata(skill *skpb.Skill, metadata *mpb.Metadata) {
	if skill == nil {
		return
	}

	if metadata.GetIdVersion().GetVersion() != "" {
		skill.IdVersion = idutils.IDVersionFromProtoUnchecked(metadata.GetIdVersion())
	} else {
		skill.IdVersion = ""
	}
}

// FillInSkillMetadataFromAssetMetadata overwrites in the Skill proto of the given BehaviorTree the
// Asset related metadata such as `id_version` and `description` with the values from the given
// Asset metadata. The Skill's "interface description" (including `parameter_description`,
// `return_value_description` and `resource_selectors`) remains unchanged. The Skill proto is
// created if it is currently missing.
//
// Use this method to make the top-level Skill description in a BehaviorTree proto consistent with
// the Asset metadata.
func FillInSkillMetadataFromAssetMetadata(behaviorTree *btpb.BehaviorTree, metadata *mpb.Metadata) {
	if behaviorTree.Description == nil {
		behaviorTree.Description = &skpb.Skill{}
	}

	skill := behaviorTree.Description

	skill.SkillName = metadata.GetIdVersion().GetId().GetName()
	skill.PackageName = metadata.GetIdVersion().GetId().GetPackage()
	skill.Id = idutils.IDFromProtoUnchecked(metadata.GetIdVersion().GetId())
	skill.Description = metadata.GetDocumentation().GetDescription()
	skill.BehaviorTreeDescription = &skpb.BehaviorTreeDescription{}

	FillInSkillIDVersionFromAssetMetadata(skill, metadata)
}
