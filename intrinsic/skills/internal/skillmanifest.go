// Copyright 2023 Intrinsic Innovation LLC

// Package skillmanifest provides utils for working with Skill manifests.
package skillmanifest

import (
	"fmt"

	"intrinsic/assets/dependencies/utils"
	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"
	"intrinsic/util/proto/sourcecodeinfoview"

	"google.golang.org/protobuf/reflect/protoregistry"

	smpb "intrinsic/skills/proto/skill_manifest_go_proto"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

var errMixOfDependencyModels = fmt.Errorf("cannot declare dependencies in both the manifest's dependencies field (required equipment) and in the skill's parameter proto (annotation-based dependencies)")

type validateSkillManifestOptions struct {
	types                                    *protoregistry.Types
	incompatibleDisallowManifestDependencies bool
}

// ValidateSkillManifestOption is an option for validating a SkillManifest.
type ValidateSkillManifestOption func(*validateSkillManifestOptions)

// WithTypes provides a Types for validating proto messages.
func WithTypes(types *protoregistry.Types) ValidateSkillManifestOption {
	return func(opts *validateSkillManifestOptions) {
		opts.types = types
	}
}

// WithIncompatibleDisallowManifestDependencies specifies whether to prevent the SkillManifest from
// declaring dependencies in the manifest.
func WithIncompatibleDisallowManifestDependencies(incompatible bool) ValidateSkillManifestOption {
	return func(opts *validateSkillManifestOptions) {
		opts.incompatibleDisallowManifestDependencies = incompatible
	}
}

// ValidateSkillManifest validates a SkillManifest.
func ValidateSkillManifest(m *smpb.SkillManifest, options ...ValidateSkillManifestOption) error {
	opts := &validateSkillManifestOptions{}
	for _, opt := range options {
		opt(opts)
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

	if opts.types != nil {
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
