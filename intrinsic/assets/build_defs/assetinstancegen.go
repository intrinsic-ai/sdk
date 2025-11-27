// Copyright 2023 Intrinsic Innovation LLC

// assetinstancegen creates an AssetInstanceInfo proto which configures an asset instance.
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"intrinsic/assets/idutils"
	"intrinsic/production/intrinsic"
	"intrinsic/util/proto/protoio"

	log "github.com/golang/glog"

	assetpb "intrinsic/assets/build_defs/asset_go_proto"
)

var (
	configPath           = flag.String("config_path", "", "The asset's configuration file path.")
	id                   = flag.String("id", "", "The id of the asset.")
	instanceName         = flag.String("instance_name", "", "The asset's instance name.")
	requiredNodeHostname = flag.String("required_node_hostname", "", "The node's hostname where asset is required to be run (only applicable for services).")
	outputAssetInstance  = flag.String("output_asset_instance", "", "Output AssetInstance proto path.")
)

var validInstanceNameRegexp = regexp.MustCompile(`^[a-z]([a-z0-9_]*[a-z0-9])?$`)

func validateInstanceName(name string) error {
	if !validInstanceNameRegexp.MatchString(name) {
		return fmt.Errorf("name must start with a lowercase letter, must use only lowercase letters, numbers and underscores, and must not end with an underscore (got: %q)", name)
	}
	return nil
}

func main() {
	intrinsic.Init()
	if *outputAssetInstance == "" {
		log.Exitf("--output_asset_instance is required")
	}

	idp, err := idutils.NewIDProto(*id)
	if err != nil {
		log.Exitf("invalid asset id: %v", err)
	}
	if err := validateInstanceName(*instanceName); err != nil {
		log.Exitf("invalid asset instance name: %v", err)
	}

	var config *assetpb.AssetInstanceInfo_TextProto
	if *configPath != "" {
		b, err := os.ReadFile(*configPath)
		if err != nil {
			log.Exitf("could not read asset instance config for %q: %v", *id, err)
		}
		config = &assetpb.AssetInstanceInfo_TextProto{
			TextProto: string(b),
		}
	}
	// Presume empty should remain unset.
	if *requiredNodeHostname == "" {
		requiredNodeHostname = nil
	}

	instance := &assetpb.AssetInstanceInfo{
		Id:                   idp,
		InstanceName:         *instanceName,
		Config:               config,
		RequiredNodeHostname: requiredNodeHostname,
	}
	if err := protoio.WriteBinaryProto(*outputAssetInstance, instance, protoio.WithDeterministic(true)); err != nil {
		log.Exitf("could not write instance for %q: %v", *instanceName, err)
	}
}
