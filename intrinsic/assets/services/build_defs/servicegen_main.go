// Copyright 2023 Intrinsic Innovation LLC

// package main implements creation of the service type bundle.
package main

import (
	"strings"

	"flag"
	log "github.com/golang/glog"
	"intrinsic/assets/services/build_defs/servicegen"
	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	intrinsic "intrinsic/production/intrinsic"
	"intrinsic/util/proto/protoio"
)

var (
	flagDefaultConfig      = flag.String("default_config", "", "Optional path to default config proto.")
	flagFileDescriptorSets = flag.String("file_descriptor_sets", "", "Comma separated paths to binary file descriptor set protos associated with the manifest.")
	flagImageTars          = flag.String("image_tars", "", "Comma separated full paths to tar archives for images.")
	flagManifest           = flag.String("manifest", "", "Path to a ServiceManifest pbtxt file.")
	flagOutputBundle       = flag.String("output_bundle", "", "Bundle tar path.")
)

func main() {
	intrinsic.Init()

	var fds []string
	if *flagFileDescriptorSets != "" {
		fds = strings.Split(*flagFileDescriptorSets, ",")
	}

	var imageTarsList []string
	if *flagImageTars != "" {
		imageTarsList = strings.Split(*flagImageTars, ",")
	}

	m := new(smpb.ServiceManifest)
	if err := protoio.ReadTextProto(*flagManifest, m); err != nil {
		log.Exitf("Failed to read manifest: %v", err)
	}

	data := servicegen.ServiceData{
		DefaultConfig:      *flagDefaultConfig,
		FileDescriptorSets: fds,
		ImageTars:          imageTarsList,
		Manifest:           m,
		OutputBundle:       *flagOutputBundle,
	}
	if err := servicegen.CreateService(&data); err != nil {
		log.Exitf("Couldn't create service type: %v", err)
	}
}
