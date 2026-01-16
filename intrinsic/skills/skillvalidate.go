// Copyright 2023 Intrinsic Innovation LLC

// Package skillvalidate provides utils for validating Skills.
package skillvalidate

import (
	"fmt"

	"intrinsic/assets/dependencies/utils"
	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"

	"google.golang.org/protobuf/reflect/protoregistry"

	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
)

var errMixOfDependencyModels = fmt.Errorf("cannot declare dependencies in both the manifest's dependencies field (required equipment) and in the skill's parameter proto (annotation-based dependencies)")

type skillManifestOptions struct {
	types                                    *protoregistry.Types
	incompatibleDisallowManifestDependencies bool
}

// SkillManifestOption is an option for validating a SkillManifest.
type SkillManifestOption func(*skillManifestOptions)

// WithTypes provides a Types for validating proto messages.
func WithTypes(types *protoregistry.Types) SkillManifestOption {
	return func(opts *skillManifestOptions) {
		opts.types = types
	}
}

// WithIncompatibleDisallowManifestDependencies specifies whether to prevent the SkillManifest from
// declaring dependencies in the manifest.
func WithIncompatibleDisallowManifestDependencies(incompatible bool) SkillManifestOption {
	return func(opts *skillManifestOptions) {
		opts.incompatibleDisallowManifestDependencies = incompatible
	}
}

// SkillManifest validates a SkillManifest.
func SkillManifest(m *smpb.SkillManifest, options ...SkillManifestOption) error {
	opts := &skillManifestOptions{}
	for _, opt := range options {
		opt(opts)
	}
	if opts.types == nil {
		return fmt.Errorf("types option must be specified")
	}

	if m == nil {
		return fmt.Errorf("SkillManifest must not be nil")
	}

	if err := metadatautils.ValidateManifestMetadata(m); err != nil {
		return fmt.Errorf("invalid SkillManifest metadata: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetId())

	if opts.incompatibleDisallowManifestDependencies && len(m.GetDependencies().GetRequiredEquipment()) > 0 {
		return fmt.Errorf("Skill %q declares dependencies in the manifest's dependencies field but --incompatible_disallow_manifest_dependencies is true", id)
	}

	if name := m.GetParameter().GetMessageFullName(); name != "" {
		mt, err := opts.types.FindMessageByURL(name)
		if err != nil {
			return fmt.Errorf("cannot find parameter message %q for Skill %q: %w", name, id, err)
		}
		parameterHasResolvedDependencies := utils.HasResolvedDependency(mt.Descriptor())
		if parameterHasResolvedDependencies && len(m.GetDependencies().GetRequiredEquipment()) != 0 {
			return errMixOfDependencyModels
		}
	}
	if name := m.GetReturnType().GetMessageFullName(); name != "" {
		if _, err := opts.types.FindMessageByURL(name); err != nil {
			return fmt.Errorf("cannot find return type message %q for Skill %q: %w", name, id, err)
		}
	}
	return nil
}
