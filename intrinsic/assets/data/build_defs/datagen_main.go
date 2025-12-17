// Copyright 2023 Intrinsic Innovation LLC

// package main is the entrypoint for creating Data Asset bundles.
package main

import (
	"flag"

	"intrinsic/assets/data/build_defs/datagen"
	"intrinsic/production/intrinsic"
	intrinsicflag "intrinsic/util/flag"

	log "github.com/golang/glog"
)

var (
	manifestPath                = flag.String("manifest", "", "Path to the DataManifest textproto file.")
	excludeReferencedFilePaths  = intrinsicflag.MultiString("exclude_referenced_file", nil, "Path to a referenced data file that should not be copied into the bundle. Can be repeated.")
	expectedReferencedFilePaths = intrinsicflag.MultiString("expected_referenced_file", nil, "Path to a file that is expected to be referenced by the Data Asset. Can be repeated.")
	fileDescriptorSetPaths      = intrinsicflag.MultiString("file_descriptor_set", nil, "Path to a binary file descriptor set proto to be used to resolve the data payload. Can be repeated.")
	outputBundlePath            = flag.String("output_bundle", "", "Output path for the .tar bundle.")
)

func main() {
	intrinsic.Init()

	if err := datagen.CreateDataBundle(&datagen.CreateDataBundleOptions{
		ManifestPath:                *manifestPath,
		ExcludedReferencedFilePaths: *excludeReferencedFilePaths,
		ExpectedReferencedFilePaths: *expectedReferencedFilePaths,
		FileDescriptorSetPaths:      *fileDescriptorSetPaths,
		OutputBundlePath:            *outputBundlePath,
	}); err != nil {
		log.Exitf("failed to create Data Asset bundle: %v", err)
	}
}
