// Copyright 2023 Intrinsic Innovation LLC

// Main implements the main CLI for generating a SkillServiceConfig file
package main

import (
	"fmt"

	"flag"
	log "github.com/golang/glog"
	"intrinsic/production/intrinsic"
	sscg "intrinsic/skills/build_defs/skillserviceconfiggen"
	"intrinsic/util/proto/protoio"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
)

var (
	flagManifestPbbinFilename   = flag.String("manifest_pbbin_filename", "", "Filename for the binary skill manifest proto.")
	flagProtoDescriptorFilename = flag.String("proto_descriptor_filename", "", "Filename for FileDescriptorSet for skill parameter, return value and published topic protos.")
	flagOutputConfigFilename    = flag.String("output_config_filename", "", "Output filename.")
)

func checkArguments() error {
	if len(*flagManifestPbbinFilename) == 0 {
		return fmt.Errorf("--manifest_pbbin_filename is required")
	}
	if len(*flagProtoDescriptorFilename) == 0 {
		return fmt.Errorf("--proto_descriptor_filename is required")
	}
	if len(*flagProtoDescriptorFilename) == 0 {
		return fmt.Errorf("--output_config_filename is required")
	}
	return nil
}

func main() {
	intrinsic.Init()
	// Fail fast if CLI arguments are invalid.
	if err := checkArguments(); err != nil {
		log.Exitf("Invalid arguments: %v", err)
	}

	fileDescriptorSet := new(dpb.FileDescriptorSet)
	if err := protoio.ReadBinaryProto(*flagProtoDescriptorFilename, fileDescriptorSet); err != nil {
		log.Exitf("Unable to read FileDescriptorSet: %v", err)
	}

	manifest := new(smpb.SkillManifest)
	if err := protoio.ReadBinaryProto(*flagManifestPbbinFilename, manifest); err != nil {
		log.Exitf("Unable to read manifest: %v", err)
	}

	skillServiceConfig, err := sscg.ExtractSkillServiceConfigFromManifest(manifest, fileDescriptorSet)
	if err != nil {
		log.Exitf("Unable to extract SkillServiceConfig: %v", err)
	}

	if err := protoio.WriteBinaryProto(*flagOutputConfigFilename, skillServiceConfig, protoio.WithDeterministic(true)); err != nil {
		log.Exitf("Unable to generate SkillServiceConfig: %v", err)
	}
}
