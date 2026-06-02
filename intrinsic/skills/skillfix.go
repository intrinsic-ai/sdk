// Copyright 2023 Intrinsic Innovation LLC

// Package skillfix contains utils that adapt a given manifest to meet the requirements of
// the latest platform version.
package skillfix

import (
	psmpb "intrinsic/skills/proto/processed_skill_manifest_go_proto"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
)

// fixOpts contains options for fixing a manifest.
type fixOpts struct {
	populateOldFields   bool
	clearObsoleteFields bool
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
	backfillSkillOptions(manifest.GetOptions(), opts)
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
	backfillSkillOptions(manifest.GetDetails().GetOptions(), opts)
	return nil
}

func backfillSkillOptions(options *smpb.Options, opts *fixOpts) {
	if options == nil {
		return
	}

	// intrinsic:assets_platform_provided_dependencies:strip_begin
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
	// intrinsic:assets_platform_provided_dependencies:strip_end
}
