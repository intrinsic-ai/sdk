// Copyright 2023 Intrinsic Innovation LLC

// main validates a skill manifest text proto and builds the binary.
package main

import (
	"flag"
	"fmt"
	"strings"

	"intrinsic/production/intrinsic"
	"intrinsic/skills/skillmanifest"
	"intrinsic/skills/skillvalidate"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"

	log "github.com/golang/glog"

	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
)

var (
	flagManifest                                 = flag.String("manifest", "", "Path to a SkillManifest pbtxt file.")
	flagOutput                                   = flag.String("output", "", "Output path.")
	flagFileDescriptorSetOut                     = flag.String("file_descriptor_set_out", "", "Output path for the file descriptor set.")
	flagFileDescriptorSets                       = flag.String("file_descriptor_sets", "", "Comma separated paths to binary file descriptor set protos.")
	flagIncompatibleDisallowManifestDependencies = flag.Bool("incompatible_disallow_manifest_dependencies", false, "Whether to prevent the skill from declaring dependencies in the manifest.")
)

func createSkillManifest() error {
	var fds []string
	if *flagFileDescriptorSets != "" {
		fds = strings.Split(*flagFileDescriptorSets, ",")
	}
	set, err := registryutil.LoadFileDescriptorSets(fds)
	if err != nil {
		return fmt.Errorf("unable to build FileDescriptorSet: %v", err)
	}

	types, err := registryutil.NewTypesFromFileDescriptorSet(set)
	if err != nil {
		return fmt.Errorf("failed to populate the registry: %v", err)
	}

	m := new(smpb.SkillManifest)
	if err := protoio.ReadTextProto(*flagManifest, m, protoio.WithResolver(types)); err != nil {
		return fmt.Errorf("failed to read manifest: %v", err)
	}
	if err := skillvalidate.SkillManifest(m,
		skillvalidate.WithTypes(types),
		skillvalidate.WithIncompatibleDisallowManifestDependencies(*flagIncompatibleDisallowManifestDependencies),
	); err != nil {
		return err
	}
	if err := protoio.WriteBinaryProto(*flagOutput, m, protoio.WithDeterministic(true)); err != nil {
		return fmt.Errorf("could not write skill manifest proto: %v", err)
	}

	skillmanifest.PruneSourceCodeInfo(m, set)
	if err := protoio.WriteBinaryProto(*flagFileDescriptorSetOut, set, protoio.WithDeterministic(true)); err != nil {
		return fmt.Errorf("could not write file descriptor set proto: %v", err)
	}
	return nil
}

func main() {
	intrinsic.Init()
	if err := createSkillManifest(); err != nil {
		log.Exitf("Failed to create skill manifest: %v", err)
	}
}
