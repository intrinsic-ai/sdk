// Copyright 2023 Intrinsic Innovation LLC

// Package main is the entrypoint to HardwareDevice asset bundle creation.
package main

import (
	"flag"
	log "github.com/golang/glog"
	"intrinsic/assets/hardware_devices/build_defs/hardwaredevicegen"
	"intrinsic/production/intrinsic"
	intrinsicflag "intrinsic/util/flag"
)

var (
	manifest      = flag.String("manifest", "", "Path to the HardwareDeviceManifest .textproto file.")
	localAssets   = intrinsicflag.MultiString("local_asset", nil, "Path to serialized AssetLocalInfo proto for an asset to add to the manifest. Can be repeated.")
	catalogAssets = intrinsicflag.MultiString("catalog_asset", nil, "Path to serialized AssetCatalogRefInfo proto for an asset to add to the manifest. Can be repeated.")
	outputBundle  = flag.String("output_bundle", "", "Path to the output .tar bundle file.")
)

func main() {
	intrinsic.Init()

	if err := hardwaredevicegen.CreateHardwareDeviceBundle(hardwaredevicegen.CreateHardwareDeviceBundleOptions{
		AssetCatalogRefInfoPaths: *catalogAssets,
		AssetLocalInfoPaths:      *localAssets,
		ManifestPath:             *manifest,
		OutputBundlePath:         *outputBundle,
	}); err != nil {
		log.Exitf("Could not create HardwareDevice asset bundle: %v", err)
	}
}
