// Copyright 2023 Intrinsic Innovation LLC

// Package main is the entrypoint for creating HardwareDevice Asset bundles.
package main

import (
	"flag"

	"intrinsic/assets/hardware_devices/build_defs/hardwaredevicegen"
	"intrinsic/production/intrinsic"
	intrinsicflag "intrinsic/util/flag"

	log "github.com/golang/glog"
)

var (
	manifestPath      = flag.String("manifest", "", "Path to the HardwareDeviceManifest textproto file.")
	localAssetPaths   = intrinsicflag.MultiString("local_asset", nil, "Path to serialized AssetLocalInfo proto for an asset to add to the manifest. Can be repeated.")
	catalogAssetPaths = intrinsicflag.MultiString("catalog_asset", nil, "Path to serialized AssetCatalogRefInfo proto for an asset to add to the manifest. Can be repeated.")
	outputBundlePath  = flag.String("output_bundle", "", "Output path for the .tar bundle.")
)

func main() {
	intrinsic.Init()

	if err := hardwaredevicegen.CreateHardwareDeviceBundle(&hardwaredevicegen.CreateHardwareDeviceBundleOptions{
		ManifestPath:             *manifestPath,
		AssetLocalInfoPaths:      *localAssetPaths,
		AssetCatalogRefInfoPaths: *catalogAssetPaths,
		OutputBundlePath:         *outputBundlePath,
	}); err != nil {
		log.Exitf("failed to create HardwareDevice Asset bundle: %v", err)
	}
}
