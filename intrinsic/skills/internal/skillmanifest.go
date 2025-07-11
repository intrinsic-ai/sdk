// Copyright 2023 Intrinsic Innovation LLC

// Package skillmanifest contains tools for working with SkillManifest.
package skillmanifest

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoregistry"
	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"
	"intrinsic/util/proto/sourcecodeinfoview"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
)

// ValidateSkillManifestOptions contains options for validating a SkillManifest.
type ValidateSkillManifestOptions struct {
	types *protoregistry.Types
}

// ValidateSkillManifestOption is an option for validating a SkillManifest.
type ValidateSkillManifestOption func(*ValidateSkillManifestOptions)

// WithTypes adds the proto types to the validation options.
func WithTypes(types *protoregistry.Types) ValidateSkillManifestOption {
	return func(opts *ValidateSkillManifestOptions) {
		opts.types = types
	}
}

// ValidateSkillManifest checks that a SkillManifest is consistent and valid.
func ValidateSkillManifest(m *smpb.SkillManifest, options ...ValidateSkillManifestOption) error {
	opts := &ValidateSkillManifestOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if err := metadatautils.ValidateManifestMetadata(m); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetId())

	if opts.types != nil {
		if name := m.GetParameter().GetMessageFullName(); name != "" {
			if _, err := opts.types.FindMessageByURL(name); err != nil {
				return fmt.Errorf("cannot find parameter message %q for Skill %q: %w", name, id, err)
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
