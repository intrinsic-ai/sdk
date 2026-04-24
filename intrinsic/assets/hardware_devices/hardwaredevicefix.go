// Copyright 2023 Intrinsic Innovation LLC

// Package hardwaredevicefix contains utils that adapt a given manifest to meet the requirements of
// the latest platform version.
package hardwaredevicefix

import (
	"fmt"

	hdmpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"

	"intrinsic/assets/data/datafix"
	"intrinsic/assets/scene_objects/sceneobjectfix"
	"intrinsic/assets/services/servicefix"
)

// fixOpts contains options for fixing a manifest.
type fixOpts struct {
	serviceOptions []servicefix.FixOption
}

// FixOption is an option for fixing a manifest.
type FixOption func(*fixOpts)

// WithClearObsoleteFields specifies whether to clear obsolete manifest fields. A field is
// considered obsolete if the platform no longer uses it.
func WithClearObsoleteFields(clear bool) FixOption {
	return WithServiceOptions(servicefix.WithClearObsoleteFields(clear))
}

// WithPopulateOldFields specifies whether to backfill old deprecated fields if empty.
func WithPopulateOldFields(populate bool) FixOption {
	return WithServiceOptions(servicefix.WithPopulateOldFields(populate))
}

// WithServiceOptions appends options to use for fixing Service Assets.
func WithServiceOptions(options ...servicefix.FixOption) FixOption {
	return func(opts *fixOpts) {
		opts.serviceOptions = append(opts.serviceOptions, options...)
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
	for k, pa := range manifest.GetAssets() {
		// Note that non-inlined Assets that are stored in the catalog will need to be "fixed"
		// downstream when the asset is installed.
		switch pa.GetVariant().(type) {
		case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Data:
			if err := datafix.DataAsset(pa.GetData()); err != nil {
				return fmt.Errorf("failed to fix Data Asset %q: %w", k, err)
			}
		case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_SceneObject:
			if err := sceneobjectfix.ProcessedManifest(pa.GetSceneObject()); err != nil {
				return fmt.Errorf("failed to fix SceneObject Asset %q: %w", k, err)
			}
		case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Service:
			if err := servicefix.ProcessedManifest(pa.GetService(), opts.serviceOptions...); err != nil {
				return fmt.Errorf("failed to fix Service Asset %q: %w", k, err)
			}
		}
	}
	return nil
}
