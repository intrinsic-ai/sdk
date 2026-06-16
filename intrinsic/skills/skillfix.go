// Copyright 2023 Intrinsic Innovation LLC

// Package skillfix contains utils that adapt a given manifest to meet the requirements of
// the latest platform version.
package skillfix

import (
	_ "embed"
	"fmt"
	"sync"

	"intrinsic/util/proto/descriptor"
	"intrinsic/util/proto/sourcecodeinfoview"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"

	psmpb "intrinsic/skills/proto/processed_skill_manifest_go_proto"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
	skillservicepb "intrinsic/skills/proto/skill_service_go_proto"

	dpb "google.golang.org/protobuf/types/descriptorpb"
)

//go:embed generator/skill_services_provided_to_platform_transitive_set_sci.proto.bin
var providedToPlatformFDSBytes []byte

var (
	cachedProvidedToPlatformFDS *dpb.FileDescriptorSet
	providedToPlatformFDSOnce   sync.Once
	providedToPlatformFDSErr    error
	// platformFDSOverride is only used to override the platform FDS in tests.
	platformFDSOverride *dpb.FileDescriptorSet
)

func providedToPlatformFDS() (*dpb.FileDescriptorSet, error) {
	providedToPlatformFDSOnce.Do(func() {
		fds := &dpb.FileDescriptorSet{}
		if err := proto.Unmarshal(providedToPlatformFDSBytes, fds); err != nil {
			providedToPlatformFDSErr = fmt.Errorf("failed to unmarshal platform descriptors: %w", err)
			return
		}
		if err := sourcecodeinfoview.PruneSourceCodeInfo(fds); err != nil {
			providedToPlatformFDSErr = fmt.Errorf("failed to prune source code info: %w", err)
			return
		}
		cachedProvidedToPlatformFDS = fds
	})
	return cachedProvidedToPlatformFDS, providedToPlatformFDSErr
}

// fixOpts contains options for fixing a manifest.
type fixOpts struct {
	populateOldFields                         bool
	clearObsoleteFields                       bool
	mergeMissingProvidedToPlatformDescriptors bool
}

// FixOption is an option for fixing a manifest.
type FixOption func(*fixOpts)

// WithPopulateOldFields specifies whether to backfill old deprecated fields if empty.
func WithPopulateOldFields(populate bool) FixOption {
	return func(opts *fixOpts) {
		opts.populateOldFields = populate
	}
}

// WithClearObsoleteFields specifies whether to clear obsolete manifest fields. A field is
// considered obsolete if the platform no longer uses it.
func WithClearObsoleteFields(clear bool) FixOption {
	return func(opts *fixOpts) {
		opts.clearObsoleteFields = clear
	}
}

// WithMergeMissingProvidedToPlatformDescriptors specifies whether to merge missing platform
// provided descriptors into the processed skill manifest.
func WithMergeMissingProvidedToPlatformDescriptors(merge bool) FixOption {
	return func(opts *fixOpts) {
		opts.mergeMissingProvidedToPlatformDescriptors = merge
	}
}

// Manifest updates a SkillManifest to meet the requirements of the latest platform version.
func Manifest(manifest *smpb.SkillManifest, options ...FixOption) error {
	opts := &fixOpts{}
	for _, opt := range options {
		opt(opts)
	}
	if manifest == nil {
		return nil
	}

	if manifest.GetOptions() == nil {
		manifest.Options = &smpb.Options{}
	}
	backfillSkillOptions(manifest.GetOptions())
	return nil
}

// ProcessedManifest updates a ProcessedSkillManifest to meet the requirements of the latest
// platform version.
func ProcessedManifest(manifest *psmpb.ProcessedSkillManifest, options ...FixOption) error {
	opts := &fixOpts{}
	for _, opt := range options {
		opt(opts)
	}
	if manifest == nil {
		return nil
	}

	if manifest.GetDetails() == nil {
		manifest.Details = &psmpb.SkillDetails{}
	}
	if manifest.GetDetails().GetOptions() == nil {
		manifest.Details.Options = &smpb.Options{}
	}
	backfillSkillOptions(manifest.GetDetails().GetOptions())

	if opts.mergeMissingProvidedToPlatformDescriptors {
		if err := mergeMissingProvidedToPlatformDescriptors(manifest); err != nil {
			return err
		}
	}

	return nil
}

func backfillSkillOptions(options *smpb.Options) {
	if options == nil {
		return
	}

	// If SkillsServicesConfig is not present, we assume this skill provides the following skill
	// service gRPC interfaces to the platform.
	if options.GetSkillServicesConfig() == nil {
		options.SkillServicesConfig = &smpb.SkillServicesConfig{
			ServiceVersions: []smpb.SkillServicesConfig_ServiceVersion{
				smpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_PROJECTOR,
				smpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_EXECUTOR,
				smpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_SKILL_INFORMATION,
			},
		}
	}

}

func mergeMissingProvidedToPlatformDescriptors(manifest *psmpb.ProcessedSkillManifest) error {
	versions := manifest.GetDetails().GetOptions().GetSkillServicesConfig().GetServiceVersions()
	// Only merge platform descriptors if the skill exactly matches the initial platform service
	// versions. This is called after backfilling, ensuring that skills with missing service versions
	// (which are backfilled to initial) will still trigger this merge.
	if !hasOnlyInitialPlatformServiceVersions(versions) {
		return nil
	}

	if manifest.GetAssets() == nil {
		manifest.Assets = &psmpb.ProcessedSkillAssets{}
	}
	if manifest.GetAssets().GetFileDescriptorSet() == nil {
		manifest.Assets.FileDescriptorSet = &dpb.FileDescriptorSet{}
	}

	// Check if the platform service descriptors are already present.
	present, err := initialPlatformServicesPresent(manifest.GetAssets().GetFileDescriptorSet())
	if err != nil {
		return err
	}
	if present {
		return nil
	}

	// platformFDSOverride is only set in tests.
	platformFDS := platformFDSOverride
	if platformFDS == nil {
		var err error
		platformFDS, err = providedToPlatformFDS()
		if err != nil {
			return err
		}
	}

	// Loose comparison preprocessor that treats all descriptors as equal if they have the same name.
	// This ensures that the merge operation never fails due to conflicting definitions in duplicate
	// files.
	byNamePreprocessor := func(fds *dpb.FileDescriptorSet) (*dpb.FileDescriptorSet, error) {
		cmpFDS := &dpb.FileDescriptorSet{}
		for _, file := range fds.GetFile() {
			cmpFDS.File = append(cmpFDS.File, &dpb.FileDescriptorProto{
				Name: proto.String(file.GetName()),
			})
		}
		return cmpFDS, nil
	}

	mergedFDS, err := descriptor.MergeFileDescriptorSets(
		[]*dpb.FileDescriptorSet{
			manifest.GetAssets().GetFileDescriptorSet(),
			platformFDS,
		},
		descriptor.WithComparisonPreprocessor(byNamePreprocessor),
	)
	if err != nil {
		return fmt.Errorf("failed to merge platform descriptors: %w", err)
	}

	manifest.Assets.FileDescriptorSet = mergedFDS
	return nil
}

func hasOnlyInitialPlatformServiceVersions(versions []smpb.SkillServicesConfig_ServiceVersion) bool {
	if len(versions) != 3 {
		return false
	}
	hasProjector := false
	hasExecutor := false
	hasInfo := false
	for _, v := range versions {
		switch v {
		case smpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_PROJECTOR:
			hasProjector = true
		case smpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_EXECUTOR:
			hasExecutor = true
		case smpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_SKILL_INFORMATION:
			hasInfo = true
		}
	}
	return hasProjector && hasExecutor && hasInfo
}

func initialPlatformServicesPresent(fds *dpb.FileDescriptorSet) (bool, error) {
	if fds == nil || len(fds.GetFile()) == 0 {
		return false, nil
	}
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return false, fmt.Errorf("failed to parse file descriptor set: %w", err)
	}

	services := []string{
		skillservicepb.Projector_ServiceDesc.ServiceName,
		skillservicepb.Executor_ServiceDesc.ServiceName,
		skillservicepb.SkillInformation_ServiceDesc.ServiceName,
	}

	for _, svc := range services {
		_, err := files.FindDescriptorByName(protoreflect.FullName(svc))
		if err != nil {
			return false, nil
		}
	}
	return true, nil
}
