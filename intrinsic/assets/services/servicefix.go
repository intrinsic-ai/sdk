// Copyright 2023 Intrinsic Innovation LLC

// Package servicefix contains utils that adapt a given manifest to meet the requirements of
// the latest platform version.
package servicefix

import (

	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
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

}
