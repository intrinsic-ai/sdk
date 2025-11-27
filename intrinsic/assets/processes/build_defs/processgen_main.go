// Copyright 2023 Intrinsic Innovation LLC

// package main is the entrypoint for the processgen binary.
package main

import (
	"flag"

	"intrinsic/assets/processes/build_defs/processgen"
	"intrinsic/production/intrinsic"
	intrinsicflag "intrinsic/util/flag"

	log "github.com/golang/glog"
)

var (
	manifestPath                    = flag.String("manifest", "", "Path to a ProcessManifest textproto file.")
	behaviorTreePath                = flag.String("behavior_tree", "", "Optional path to a BehaviorTree textproto or binary proto file (if the manifest does not contain a behavior tree).")
	textprotoFileDescriptorSetPaths = intrinsicflag.MultiString("textproto_file_descriptor_set", nil, "Path to a binary file descriptor set proto to be used for parsing expanded Any protos in the input textprotos. Can be repeated.")
	outputBundlePath                = flag.String("output_bundle", "", "Output path for the tar bundle.")
	outputFileDescriptorSetPath     = flag.String("output_file_descriptor_set", "", "Output path to which the parameter file descriptor set of the behavior tree will be written as a binary proto. If the behavior tree has no parameter file descriptor set, an empty file will be written.")
	outputManifestBinaryPath        = flag.String("output_manifest_binary", "", "Output path to which the manifest will be written as a binary proto.")
)

func main() {
	intrinsic.Init()

	if *manifestPath == "" {
		log.Exitf("--manifest is required")
	}
	if *outputBundlePath == "" {
		log.Exitf("--output_bundle is required")
	}
	if *outputFileDescriptorSetPath == "" {
		log.Exitf("--output_file_descriptor_set is required")
	}
	if *outputManifestBinaryPath == "" {
		log.Exitf("--output_manifest_binary is required")
	}

	if err := processgen.CreateProcessAssetBundle(processgen.CreateProcessAssetBundleOptions{
		ManifestPath:                    *manifestPath,
		BehaviorTreePath:                *behaviorTreePath,
		TextprotoFileDescriptorSetPaths: *textprotoFileDescriptorSetPaths,
		OutputBundlePath:                *outputBundlePath,
		OutputFileDescriptorSetPath:     *outputFileDescriptorSetPath,
		OutputManifestBinaryPath:        *outputManifestBinaryPath,
	}); err != nil {
		log.Exitf("creating ProcessManifest: %v", err)
	}
}
