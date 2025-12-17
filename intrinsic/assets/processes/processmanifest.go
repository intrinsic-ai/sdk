// Copyright 2023 Intrinsic Innovation LLC

// Package processmanifest provides utils for working with Process manifests.
package processmanifest

import (
	"fmt"

	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	papb "intrinsic/assets/processes/proto/process_asset_go_proto"
	pmpb "intrinsic/assets/processes/proto/process_manifest_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	docpb "intrinsic/assets/proto/documentation_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	mpb "intrinsic/assets/proto/metadata_go_proto"
	btpb "intrinsic/executive/proto/behavior_tree_go_proto"
	skpb "intrinsic/skills/proto/skills_go_proto"
)

var (
	// ErrMissingBehaviorTree is returned when a Process is missing a behavior tree.
	ErrMissingBehaviorTree = status.Errorf(codes.InvalidArgument, "'behavior_tree' must be specified")

	// ErrBehaviorTreeNameInconsistent is returned when a Process has an inconsistent behavior tree
	// name.
	ErrBehaviorTreeNameInconsistent = status.Errorf(
		codes.InvalidArgument,
		"'behavior_tree.name' must match 'metadata.display_name'")

	// ErrSkillProtoMissing is returned when the behavior tree of a Process does not have a Skill
	// proto.
	ErrSkillProtoMissing = status.Errorf(
		codes.InvalidArgument,
		"'behavior_tree.description' must be set")

	// ErrBehaviorTreeDescriptionMissing is returned when the behavior tree of a Process has a Skill
	// proto but the behavior tree description is missing.
	ErrBehaviorTreeDescriptionMissing = status.Errorf(
		codes.InvalidArgument,
		"'behavior_tree.description.behavior_tree_description' must be set")

	// ErrSkillNameMissing is returned when the behavior tree of a Process has a Skill proto but the
	// Skill name is missing.
	ErrSkillNameMissing = status.Errorf(
		codes.InvalidArgument,
		"'behavior_tree.description.skill_name' is missing but must be equal to 'metadata.id_version.id.name'")

	// ErrSkillNameInconsistent is returned when the behavior tree of a Process has a Skill proto but
	// the Skill name does not match the Asset name.
	ErrSkillNameInconsistent = status.Errorf(
		codes.InvalidArgument,
		"'behavior_tree.description.skill_name' must be equal to 'metadata.id_version.id.name'")

	// ErrSkillPackageNameMissing is returned when the behavior tree of a Process has a Skill proto
	// but the Skill package name is missing.
	ErrSkillPackageNameMissing = status.Errorf(
		codes.InvalidArgument,
		"'behavior_tree.description.package_name' is missing but must be equal to 'metadata.id_version.id.package'")

	// ErrSkillPackageNameInconsistent is returned when the behavior tree of a Process has a Skill
	// proto but the Skill package name does not match the Asset package name.
	ErrSkillPackageNameInconsistent = status.Errorf(
		codes.InvalidArgument,
		"'behavior_tree.description.package_name' must be equal to 'metadata.id_version.id.package'")

	// ErrSkillIDMissing is returned when the behavior tree of a Process has a Skill proto but the
	// Skill ID is missing.
	ErrSkillIDMissing = status.Errorf(
		codes.InvalidArgument,
		"'behavior_tree.description.id' is missing but must be equal to 'metadata.id_version.id'")

	// ErrSkillIDInconsistent is returned when the behavior tree of a Process has a Skill proto but
	// the Skill ID does not match the Asset ID.
	ErrSkillIDInconsistent = status.Errorf(
		codes.InvalidArgument,
		"'behavior_tree.description.id' must be equal to 'metadata.id_version.id'")

	// ErrSkillIDVersionMissing is returned when the behavior tree of a Process has a Skill proto but
	// the Skill ID version is missing.
	ErrSkillIDVersionMissing = status.Errorf(
		codes.InvalidArgument,
		"'behavior_tree.description.id_version' is missing but must be equal to 'metadata.id_version'")

	// ErrSkillIDVersionInconsistent is returned when the behavior tree of a Process has a Skill proto
	// but the skill ID version does not match the Asset ID version.
	ErrSkillIDVersionInconsistent = status.Errorf(
		codes.InvalidArgument,
		"'behavior_tree.description.id_version' must be equal to 'metadata.id_version'")

	// ErrSkillDescriptionMissing is returned when the behavior tree of a Process has a Skill proto
	// but the Skill description is missing.
	ErrSkillDescriptionMissing = status.Errorf(
		codes.InvalidArgument,
		"'behavior_tree.description.description' is missing but must be equal to 'metadata.documentation.description'")

	// ErrSkillDescriptionInconsistent is returned when the behavior tree of a Process has a Skill
	// proto but the Skill description does not match the Asset documentation.
	ErrSkillDescriptionInconsistent = status.Errorf(
		codes.InvalidArgument,
		"'behavior_tree.description.description' must be equal to 'metadata.documentation.description'")
)

type validateBehaviorTreeOptions struct {
	assetIDVersion             *idpb.IdVersion
	assetDisplayName           string
	assetDocumentation         *docpb.Documentation
	requireFilledSkillMetadata bool
}

// validateBehaviorTree validates the given behavior tree for a Process.
//
// In particular, checks the consistency of the behavior tree's name and skill proto (if present)
// with the Asset metadata.
func validateBehaviorTree(bt *btpb.BehaviorTree, options validateBehaviorTreeOptions) error {
	if bt == nil {
		return ErrMissingBehaviorTree
	}

	if bt.GetName() != options.assetDisplayName {
		return ErrBehaviorTreeNameInconsistent
	}

	skill := bt.GetDescription()

	if skill == nil {
		if options.requireFilledSkillMetadata {
			return ErrSkillProtoMissing
		}
		return nil
	}

	if options.requireFilledSkillMetadata && skill.BehaviorTreeDescription == nil {
		return ErrBehaviorTreeDescriptionMissing
	}

	// These metadata fields are redundant with the Asset's metadata. They should match the Asset's
	// metadata or can be empty (if `requireFilledSkillMetadata` is false).
	if skill.SkillName != "" {
		if skill.SkillName != options.assetIDVersion.GetId().GetName() {
			return ErrSkillNameInconsistent
		}
	} else if options.requireFilledSkillMetadata {
		return ErrSkillNameMissing
	}

	if skill.PackageName != "" {
		if skill.PackageName != options.assetIDVersion.GetId().GetPackage() {
			return ErrSkillPackageNameInconsistent
		}
	} else if options.requireFilledSkillMetadata {
		return ErrSkillPackageNameMissing
	}

	if skill.Id != "" {
		if skill.Id != idutils.IDFromProtoUnchecked(options.assetIDVersion.GetId()) {
			return ErrSkillIDInconsistent
		}
	} else if options.requireFilledSkillMetadata {
		return ErrSkillIDMissing
	}

	if skill.IdVersion != "" {
		if skill.IdVersion != idutils.IDVersionFromProtoUnchecked(options.assetIDVersion) {
			return ErrSkillIDVersionInconsistent
		}
	} else if options.requireFilledSkillMetadata && options.assetIDVersion.GetVersion() != "" {
		return ErrSkillIDVersionMissing
	}

	if skill.Description != "" {
		if skill.Description != options.assetDocumentation.GetDescription() {
			return ErrSkillDescriptionInconsistent
		}
	} else if options.requireFilledSkillMetadata && options.assetDocumentation.GetDescription() != "" {
		return ErrSkillDescriptionMissing
	}

	return nil
}

// ValidateProcessManifest validates a ProcessManifest.
func ValidateProcessManifest(manifest *pmpb.ProcessManifest) error {
	if manifest == nil {
		return fmt.Errorf("ProcessManifest must not be nil")
	}

	metadata := manifest.GetMetadata()

	if err := metadatautils.ValidateManifestMetadata(metadata); err != nil {
		return fmt.Errorf("invalid ProcessManifest metadata: %w", err)
	}

	return validateBehaviorTree(manifest.GetBehaviorTree(), validateBehaviorTreeOptions{
		assetIDVersion: &idpb.IdVersion{
			Id: metadata.GetId(),
			// Version is not specified in manifest.
		},
		assetDisplayName:   metadata.GetDisplayName(),
		assetDocumentation: metadata.GetDocumentation(),
		// In the manifest the Skill proto in the BehaviorTree is allowed to be missing or can be filled
		// partially. However, metadata fields which are filled in the Skill proto must be consistent
		// with the Asset metadata.
		requireFilledSkillMetadata: false,
	})
}

// ValidateProcess validates a Process.
//
// By default, the metadata of the given Process must specify the type ASSET_TYPE_PROCESS.
// Additional metadata validation options can be passed via `metadataOpts`.
func ValidateProcess(processAsset *papb.ProcessAsset, metadataOpts ...metadatautils.ValidateMetadataOption) error {
	metadata := processAsset.GetMetadata()

	validateMetadataOpts := append(
		[]metadatautils.ValidateMetadataOption{
			metadatautils.WithRequiredAssetType(atypepb.AssetType_ASSET_TYPE_PROCESS),
		},
		metadataOpts...,
	)
	err := metadatautils.ValidateMetadata(metadata, validateMetadataOpts...)
	if err != nil {
		return err
	}

	return validateBehaviorTree(processAsset.GetBehaviorTree(), validateBehaviorTreeOptions{
		assetIDVersion:     metadata.GetIdVersion(),
		assetDisplayName:   metadata.GetDisplayName(),
		assetDocumentation: metadata.GetDocumentation(),
		// In the processed asset the Skill proto in the BehaviorTree must be set and must be filled
		// consistently with the asset metadata.
		requireFilledSkillMetadata: true,
	})
}

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
