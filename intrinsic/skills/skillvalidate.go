// Copyright 2023 Intrinsic Innovation LLC

// Package skillvalidate provides utils for validating Skills.
package skillvalidate

import (
	"context"
	"fmt"

	"intrinsic/assets/dependencies/platform"
	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"
	validationerrors "intrinsic/assets/validation/errors"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	psmpb "intrinsic/skills/proto/processed_skill_manifest_go_proto"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
)


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
func SkillManifest(ctx context.Context, m *smpb.SkillManifest, options ...SkillManifestOption) error {
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

type processedSkillManifestOptions struct {
	report                            *validationerrors.Report
	requiredPlatformSkillDependencies []string
	requiredRegistry                  string
}

// ProcessedSkillManifestOption is an option for validating a ProcessedSkillManifest.
type ProcessedSkillManifestOption func(*processedSkillManifestOptions)

// WithRequiredRegistry specifies the registry that must have been used for all images.
func WithRequiredRegistry(registry string) ProcessedSkillManifestOption {
	return func(opts *processedSkillManifestOptions) {
		opts.requiredRegistry = registry
	}
}

// WithReport sets the shared validation Report to use for collecting warnings.
func WithReport(report *validationerrors.Report) ProcessedSkillManifestOption {
	return func(opts *processedSkillManifestOptions) {
		opts.report = report
	}
}

// WithRequiredProvidedToPlatformInterfaces specifies the protocol-prefixed interfaces a Skill must
// implement to be compatible with the current version of the platform. For example, if called with
// 'grpc://intrinsic_proto.skills.Executor', the Skill validator will generate an error if the Skill
// does not provide the Executor gRPC service to the platform.
func WithRequiredProvidedToPlatformInterfaces(required ...string) ProcessedSkillManifestOption {
	return func(opts *processedSkillManifestOptions) {
		opts.requiredPlatformSkillDependencies = append(opts.requiredPlatformSkillDependencies, required...)
	}
}

// ProcessedSkillManifest validates a ProcessedSkillManifest.
func ProcessedSkillManifest(ctx context.Context, m *psmpb.ProcessedSkillManifest, options ...ProcessedSkillManifestOption) error {
	opts := &processedSkillManifestOptions{}
	WithReport(validationerrors.NewReport())(opts)
	for _, opt := range options {
		opt(opts)
	}

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

	if opts.requiredRegistry != "" {
		switch d := m.GetAssets().GetDeploymentType().(type) {
		case *psmpb.ProcessedSkillAssets_Image:
			if d.Image.GetRegistry() != opts.requiredRegistry {
				return fmt.Errorf("unexpected registry specified (expected %q, got %q)", opts.requiredRegistry, d.Image.GetRegistry())
			}
		}
	}

	if len(opts.requiredPlatformSkillDependencies) > 0 {
		if err := validatePlatformSkillDependencies(m, opts.requiredPlatformSkillDependencies); err != nil {
			return fmt.Errorf("invalid platform skill dependencies: %w", err)
		}
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
		if _, ok := d.(protoreflect.MessageDescriptor); !ok {
			return fmt.Errorf("message %q is not a message", name)
		}

	}
	if name := sd.GetExecuteResult().GetMessageFullName(); name != "" {
		if _, err := opts.files.FindDescriptorByName(protoreflect.FullName(name)); err != nil {
			return fmt.Errorf("cannot find return type message %q: %w", name, err)
		}
	}
	return nil
}

func validatePlatformSkillDependencies(manifest *psmpb.ProcessedSkillManifest, required []string) error {
	interfaces := platform.ProvidedByProcessedSkillManifest(manifest)

	// provided is a map where keys are required interfaces and values indicate whether the Skill
	// provides them.
	provided := make(map[string]bool, len(required))
	for _, r := range required {
		provided[r] = false
	}
	for _, iface := range interfaces {
		ifaceURI := iface.GetUri()
		if _, ok := provided[ifaceURI]; ok {
			provided[ifaceURI] = true
		}
	}
	for uri, found := range provided {
		if !found {
			return fmt.Errorf("this platform version requires that each Skill provide %q", uri)
		}
	}
	return nil
}
