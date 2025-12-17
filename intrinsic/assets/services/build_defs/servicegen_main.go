// Copyright 2023 Intrinsic Innovation LLC

// package main is the entrypoint for creating Service Asset bundles.
package main

import (
	"flag"

	"intrinsic/assets/services/build_defs/servicegen"
	"intrinsic/production/intrinsic"
	intrinsicflag "intrinsic/util/flag"

	log "github.com/golang/glog"
)

var (
	manifestPath           = flag.String("manifest", "", "Path to the ServiceManifest textproto file.")
	defaultConfigPath      = flag.String("default_config", "", "Optional path to default config proto.")
	fileDescriptorSetPaths = intrinsicflag.MultiString("file_descriptor_set", nil, "Path to binary file descriptor set proto associated with the manifest. Can be repeated.")
	imageTarPaths          = intrinsicflag.MultiString("image_tar", nil, "Full path to .tar archive for an image. Can be repeated.")
	outputBundlePath       = flag.String("output_bundle", "", "Output path for the .tar bundle.")
)

func main() {
	intrinsic.Init()

	if err := servicegen.CreateServiceBundle(&servicegen.CreateServiceBundleOptions{
		ManifestPath:           *manifestPath,
		DefaultConfigPath:      *defaultConfigPath,
		FileDescriptorSetPaths: *fileDescriptorSetPaths,
		ImageTarPaths:          *imageTarPaths,
		OutputBundlePath:       *outputBundlePath,
	}); err != nil {
		log.Exitf("failed to create Service Asset bundle: %v", err)
	}
}
