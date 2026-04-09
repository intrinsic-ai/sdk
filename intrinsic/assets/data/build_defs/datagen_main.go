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
	manifestPath                = flag.String("manifest_path", "", "Path to the DataManifest textproto file.")
	referencedFilePaths         = intrinsicflag.MultiString("referenced_file_path", nil, "Path to a file that is referenced in the Data Asset. Can be repeated.")
	excludedReferencedFilePaths = intrinsicflag.MultiString("excluded_referenced_file_path", nil, "Path to a referenced file that should not be copied into the bundle. Can be repeated.")
	remappedReferencedFilePaths = intrinsicflag.MultiString("remapped_referenced_file_path", nil, "For each excluded referenced path, a remapped referenced path for the output .tar bundle. Can be repeated.")
	fileDescriptorSetPaths      = intrinsicflag.MultiString("file_descriptor_set_path", nil, "Path to a binary file descriptor set proto to be used to resolve the data payload. Can be repeated.")
	outputBundlePath            = flag.String("output_bundle_path", "", "Output path for the .tar bundle.")
)

func main() {
	intrinsic.Init()

	externalReferencedFilePaths := make(map[string]string)
	for i, excludedFilePath := range *excludedReferencedFilePaths {
		externalReferencedFilePaths[excludedFilePath] = (*remappedReferencedFilePaths)[i]
	}

	if err := datagen.CreateDataBundle(&datagen.CreateDataBundleOptions{
		ManifestPath:                *manifestPath,
		ReferencedFilePaths:         *referencedFilePaths,
		ExternalReferencedFilePaths: externalReferencedFilePaths,
		FileDescriptorSetPaths:      *fileDescriptorSetPaths,
		OutputBundlePath:            *outputBundlePath,
	}); err != nil {
		log.Exitf("failed to create Data Asset bundle: %v", err)
	}
}
