// Copyright 2023 Intrinsic Innovation LLC

// Package skillbundlegen creates a skill bundle.
package main

import (
	"flag"
	log "github.com/golang/glog"
	"intrinsic/assets/bundleio"
	"intrinsic/production/intrinsic"
	"intrinsic/util/proto/protoio"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
)

var (
	flagFileDescriptorSet = flag.String("file_descriptor_set", "", "File descriptor set.")
	flagImageTar          = flag.String("image_tar", "", "Skill image file.")
	flagPBT               = flag.String("pbt", "", "Parameterized behavior tree file.")
	flagManifest          = flag.String("manifest", "", "Skill manifest.")

	flagOutputBundle = flag.String("output_bundle", "", "Output path.")
)

func main() {
	intrinsic.Init()

	fds := &descriptorpb.FileDescriptorSet{}
	if err := protoio.ReadBinaryProto(*flagFileDescriptorSet, fds); err != nil {
		log.Exitf("failed to read file descriptor set: %v", err)
	}

	m := new(smpb.SkillManifest)
	if err := protoio.ReadBinaryProto(*flagManifest, m); err != nil {
		log.Exitf("failed to read manifest: %v", err)
	}

	if err := bundleio.WriteSkill(*flagOutputBundle, bundleio.WriteSkillOpts{
		Manifest:    m,
		Descriptors: fds,
		ImageTar:    *flagImageTar,
		PBT:         *flagPBT,
	}); err != nil {
		log.Exitf("unable to write skill bundle: %v", err)
	}
}
