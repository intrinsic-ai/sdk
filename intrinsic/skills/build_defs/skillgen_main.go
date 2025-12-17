// Copyright 2023 Intrinsic Innovation LLC

// Package main is the entrypoint for creating Skill Asset bundles.
package main

import (
	"flag"

	"intrinsic/production/intrinsic"
	"intrinsic/skills/build_defs/skillgen"

	log "github.com/golang/glog"
)

var (
	manifestPath          = flag.String("manifest", "", "Path to the SkillManifest textproto file.")
	fileDescriptorSetPath = flag.String("file_descriptor_set", "", "Path to the binary file descriptor set.")
	imageTarPath          = flag.String("image_tar", "", "Path to the Skill image file.")
	outputBundlePath      = flag.String("output_bundle", "", "Output path for the Skill Asset bundle.")
)

func main() {
	intrinsic.Init()

	if err := skillgen.CreateSkillBundle(&skillgen.CreateSkillBundleOptions{
		ManifestPath:          *manifestPath,
		FileDescriptorSetPath: *fileDescriptorSetPath,
		ImageTarPath:          *imageTarPath,
		OutputBundlePath:      *outputBundlePath,
	}); err != nil {
		log.Exitf("failed to create Skill Asset bundle: %v", err)
	}
}
