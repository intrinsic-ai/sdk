// Copyright 2023 Intrinsic Innovation LLC

// Package servicefix contains utils that adapt a given manifest to meet the requirements of
// the latest platform version.
package servicefix

import (
	"slices" // intrinsic:assets_platform_provided_dependencies:strip

	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	drpb "intrinsic/assets/services/proto/v1/dynamic_reconfiguration_go_proto" // intrinsic:assets_platform_provided_dependencies:strip
	sspb "intrinsic/assets/services/proto/v1/service_state_go_proto"           // intrinsic:assets_platform_provided_dependencies:strip
)

// fixOpts contains options for fixing a manifest.
type fixOpts struct {
	populateOldFields bool
}

// FixOption is an option for fixing a manifest.
type FixOption func(*fixOpts)

// WithPopulateOldFields specifies whether to backfill old deprecated fields if empty.
func WithPopulateOldFields(populate bool) FixOption {
	return func(opts *fixOpts) {
		opts.populateOldFields = populate
	}
}

// Manifest updates a ServiceManifest to meet the requirements of the latest platform version.
func Manifest(manifest *smpb.ServiceManifest, options ...FixOption) error {
	opts := &fixOpts{}
	for _, opt := range options {
		opt(opts)
	}
	if manifest == nil {
		return nil
	}
	backfillServiceDef(manifest.GetServiceDef(), opts)
	return nil
}

// ProcessedManifest updates a ProcessedServiceManifest to meet the requirements of the latest
// platform version.
func ProcessedManifest(manifest *smpb.ProcessedServiceManifest, options ...FixOption) error {
	opts := &fixOpts{}
	for _, opt := range options {
		opt(opts)
	}
	if manifest == nil {
		return nil
	}
	backfillServiceDef(manifest.GetServiceDef(), opts)
	return nil
}

func backfillServiceDef(sd *smpb.ServiceDef, opts *fixOpts) {
	if sd == nil {
		return
	}

	// intrinsic:assets_platform_provided_dependencies:strip_begin
	// Populate the dynamic reconfiguration platform gRPC interface if only the deprecated boolean
	// setting is present and true.
	if conf := sd.GetDynamicReconfigurationConfig(); conf == nil && sd.GetSupportsDynamicReconfiguration() {
		sd.DynamicReconfigurationConfig = &drpb.DynamicReconfigurationConfig{
			ServiceVersions: []drpb.DynamicReconfigurationConfig_ServiceVersion{
				drpb.DynamicReconfigurationConfig_INTRINSIC_PROTO_SERVICES_V1_DYNAMIC_RECONFIGURATION,
			},
		}
	}

	// Populate the service state platform gRPC interface if only the deprecated boolean setting is
	// present and true.
	if conf := sd.GetServiceStateConfig(); conf == nil && sd.GetSupportsServiceState() {
		sd.ServiceStateConfig = &sspb.ServiceStateConfig{
			ServiceVersions: []sspb.ServiceStateConfig_ServiceVersion{
				sspb.ServiceStateConfig_INTRINSIC_PROTO_SERVICES_V1_SERVICE_STATE,
			},
		}
	}

	if opts.populateOldFields {
		// Backfill the deprecated SupportsDynamicReconfiguration field if the new config is present.
		if conf := sd.GetDynamicReconfigurationConfig(); conf != nil {
			if slices.Contains(conf.GetServiceVersions(), drpb.DynamicReconfigurationConfig_INTRINSIC_PROTO_SERVICES_V1_DYNAMIC_RECONFIGURATION) {
				sd.SupportsDynamicReconfiguration = true
			}
		}

		// Backfill the deprecated SupportsServiceState field if the new config is present.
		if conf := sd.GetServiceStateConfig(); conf != nil {
			if slices.Contains(conf.GetServiceVersions(), sspb.ServiceStateConfig_INTRINSIC_PROTO_SERVICES_V1_SERVICE_STATE) {
				sd.SupportsServiceState = true
			}
		}
	}
	// intrinsic:assets_platform_provided_dependencies:strip_end
}
