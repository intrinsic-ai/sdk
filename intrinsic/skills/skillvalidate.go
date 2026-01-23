// Copyright 2023 Intrinsic Innovation LLC

// Package skillvalidate provides utils for validating Skills.
package skillvalidate

import (
	"fmt"

	"intrinsic/assets/dependencies/utils"
	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	psmpb "intrinsic/skills/proto/processed_skill_manifest_go_proto"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
)

var errMixOfDependencyModels = fmt.Errorf("cannot declare dependencies in both the manifest's dependencies field (required equipment) and in the skill's parameter proto (annotation-based dependencies)")

type skillManifestOptions struct {
	files                                    *protoregistry.Files
	incompatibleDisallowManifestDependencies bool
}

// SkillManifestOption is an option for validating a SkillManifest.
type SkillManifestOption func(*skillManifestOptions)

// WithFiles provides a Files for validating proto messages.
func WithFiles(files *protoregistry.Files) SkillManifestOption {
	return func(opts *skillManifestOptions) {
		opts.files = files
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
	if opts.files == nil {
		return fmt.Errorf("files option must be specified")
	}

	if m == nil {
		return fmt.Errorf("SkillManifest must not be nil")
	}

	if err := metadatautils.ValidateManifestMetadata(m); err != nil {
		return fmt.Errorf("invalid SkillManifest metadata: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetId())

	sd := &psmpb.SkillDetails{
		Dependencies:  m.GetDependencies(),
		Parameter:     m.GetParameter(),
		ExecuteResult: m.GetReturnType(),
	}
	if err := validateSkillDetails(sd, &validateSkillDetailsOptions{
		files:                                    opts.files,
		incompatibleDisallowManifestDependencies: opts.incompatibleDisallowManifestDependencies,
	}); err != nil {
		return fmt.Errorf("invalid Skill details for %q: %w", id, err)
	}

	return nil
}

// ProcessedSkillManifest validates a ProcessedSkillManifest.
func ProcessedSkillManifest(m *psmpb.ProcessedSkillManifest) error {
	if m == nil {
		return fmt.Errorf("ProcessedSkillManifest must not be nil")
	}

	if err := metadatautils.ValidateManifestMetadata(m.GetMetadata()); err != nil {
		return fmt.Errorf("invalid ProcessedSkillManifest metadata: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetMetadata().GetId())

	if m.GetAssets() == nil || m.GetAssets().GetFileDescriptorSet() == nil {
		return fmt.Errorf("ProcessedSkillManifest file descriptor set must not be nil")
	}
	files, err := protodesc.NewFiles(m.GetAssets().GetFileDescriptorSet())
	if err != nil {
		return fmt.Errorf("failed to populate the registry: %w", err)
	}

	if err := validateSkillDetails(m.GetDetails(), &validateSkillDetailsOptions{
		files: files,
	}); err != nil {
		return fmt.Errorf("invalid Skill details for %q: %w", id, err)
	}

	return nil
}

type validateSkillDetailsOptions struct {
	files                                    *protoregistry.Files
	incompatibleDisallowManifestDependencies bool
}

func validateSkillDetails(sd *psmpb.SkillDetails, opts *validateSkillDetailsOptions) error {
	if opts.incompatibleDisallowManifestDependencies && len(sd.GetDependencies().GetRequiredEquipment()) > 0 {
		return fmt.Errorf("dependencies declared in the manifest's dependencies field but --incompatible_disallow_manifest_dependencies is true")
	}

	if name := sd.GetParameter().GetMessageFullName(); name != "" {
		d, err := opts.files.FindDescriptorByName(protoreflect.FullName(name))
		if err != nil {
			return fmt.Errorf("cannot find parameter message %q: %w", name, err)
		}
		if md, ok := d.(protoreflect.MessageDescriptor); !ok {
			return fmt.Errorf("message %q is not a message", name)
		} else if parameterHasResolvedDependencies := utils.HasResolvedDependency(md); parameterHasResolvedDependencies && len(sd.GetDependencies().GetRequiredEquipment()) != 0 {
			return errMixOfDependencyModels
		}
	}
	if name := sd.GetExecuteResult().GetMessageFullName(); name != "" {
		if _, err := opts.files.FindDescriptorByName(protoreflect.FullName(name)); err != nil {
			return fmt.Errorf("cannot find return type message %q: %w", name, err)
		}
	}
	return nil
}
