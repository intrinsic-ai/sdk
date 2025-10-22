// Copyright 2023 Intrinsic Innovation LLC

// package main is the entrypoint to Data asset bundle creation.
package main

import (
	"flag"
	log "github.com/golang/glog"
	"intrinsic/assets/data/build_defs/datagen"
	"intrinsic/production/intrinsic"
	intrinsicflag "intrinsic/util/flag"
)

var (
	excludeReferencedFile  = intrinsicflag.MultiString("exclude_referenced_file", nil, "Path to a referenced data file that should not be copied into the bundle. Can be repeated.")
	expectedReferencedFile = intrinsicflag.MultiString("expected_referenced_file", nil, "Path to a file that is expected to be referenced by the Data asset. Can be repeated.")
	fileDescriptorSets     = intrinsicflag.MultiString("file_descriptor_set", nil, "Path to a binary file descriptor set proto to be used to resolve the data payload. Can be repeated.")
	manifest               = flag.String("manifest", "", "Path to a DataManifest textproto file.")
	outputBundle           = flag.String("output_bundle", "", "Output path for the tar bundle.")
)

func main() {
	intrinsic.Init()

	if err := datagen.CreateDataAssetBundle(datagen.CreateDataAssetBundleOptions{
		ExcludedReferencedFilePaths: *excludeReferencedFile,
		ExpectedReferencedFilePaths: *expectedReferencedFile,
		FileDescriptorSetPaths:      *fileDescriptorSets,
		ManifestPath:                *manifest,
		OutputBundlePath:            *outputBundle,
	}); err != nil {
		log.Exitf("Could not create Data asset bundle: %v", err)
	}
}
