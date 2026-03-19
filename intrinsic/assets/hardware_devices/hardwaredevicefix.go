// Copyright 2023 Intrinsic Innovation LLC

// Package hardwaredevicefix contains utils that adapt a given manifest to meet the requirements of
// the latest platform version.
package hardwaredevicefix

import (
	hdmpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"

	"intrinsic/assets/services/servicefix"
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

// Manifest updates a HardwareDeviceManifest to meet the requirements of the latest
// platform version.
func Manifest(manifest *hdmpb.HardwareDeviceManifest, options ...FixOption) error {
	// There is currently nothing that can be fixed here. The unprocessed hardware device manifest
	// does not directly store any information about its constituent assets. For example, the
	// hardware device's service asset is either represented as an asset IDVersion or a path to a
	// service bundle file.
	return nil
}

// ProcessedManifest updates a ProcessedHardwareDeviceManifest to meet the requirements of the
// latest platform version.
func ProcessedManifest(manifest *hdmpb.ProcessedHardwareDeviceManifest, options ...FixOption) error {
	opts := &fixOpts{}
	for _, opt := range options {
		opt(opts)
	}
	if manifest == nil {
		return nil
	}
	for _, pa := range manifest.GetAssets() {
		// Note that non-inlined Services that are stored in the catalog will need to be "fixed"
		// downstream when the asset is installed.
		if m := pa.GetService(); m != nil {
			servicefix.ProcessedManifest(m, servicefix.WithPopulateOldFields(opts.populateOldFields), servicefix.WithClearObsoleteFields(opts.clearObsoleteFields))
		}
	}
	return nil
}
