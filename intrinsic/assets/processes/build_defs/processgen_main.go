// Copyright 2023 Intrinsic Innovation LLC

// package main is the entrypoint for creating Process Asset bundles.
package main

import (
	"flag"

	"intrinsic/assets/processes/build_defs/processgen"
	"intrinsic/production/intrinsic"
	intrinsicflag "intrinsic/util/flag"

	log "github.com/golang/glog"
)

var (
	manifestPath                    = flag.String("manifest", "", "Path to the ProcessManifest textproto file.")
	behaviorTreePath                = flag.String("behavior_tree", "", "Optional path to a BehaviorTree textproto or binary proto file (if the manifest does not contain a behavior tree).")
	textprotoFileDescriptorSetPaths = intrinsicflag.MultiString("textproto_file_descriptor_set", nil, "Path to a binary file descriptor set proto to be used for parsing expanded Any protos in the input textprotos. Can be repeated.")
	outputFileDescriptorSetPath     = flag.String("output_file_descriptor_set", "", "Output path to which the parameter file descriptor set of the behavior tree will be written as a binary proto. If the behavior tree has no parameter file descriptor set, an empty file will be written.")
	outputManifestBinaryPath        = flag.String("output_manifest_binary", "", "Output path to which the manifest will be written as a binary proto.")
	outputBundlePath                = flag.String("output_bundle", "", "Output path for the .tar bundle.")
)

func main() {
	intrinsic.Init()

	if err := processgen.CreateProcessBundle(&processgen.CreateProcessBundleOptions{
		ManifestPath:                    *manifestPath,
		BehaviorTreePath:                *behaviorTreePath,
		TextprotoFileDescriptorSetPaths: *textprotoFileDescriptorSetPaths,
		OutputBundlePath:                *outputBundlePath,
		OutputFileDescriptorSetPath:     *outputFileDescriptorSetPath,
		OutputManifestBinaryPath:        *outputManifestBinaryPath,
	}); err != nil {
		log.Exitf("failed to create Process Asset bundle: %v", err)
	}
}
