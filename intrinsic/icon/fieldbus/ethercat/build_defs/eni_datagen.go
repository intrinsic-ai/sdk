// Copyright 2023 Intrinsic Innovation LLC

package main

import (
	"flag"
	"log"
	"os"
	"unicode/utf8"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/known/anypb"

	dmpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	vpb "intrinsic/assets/proto/vendor_go_proto"
	enipb "intrinsic/icon/fieldbus/ethercat/device_service/v1/eni_go_proto"
)

var (
	eniFile           = flag.String("eni_file", "", "Path to the ENI file.")
	outputFile        = flag.String("output_file", "", "Path to the output file.")
	assetPackage      = flag.String("asset_package", "", "The package of the Data asset ID.")
	assetName         = flag.String("asset_name", "", "The name of the Data asset ID.")
	displayName       = flag.String("display_name", "", "The display name of the Data asset.")
	vendorDisplayName = flag.String("vendor_display_name", "", "The display name of the vendor.")
)

func main() {
	flag.Parse()

	if *eniFile == "" {
		log.Fatal("Missing --eni_file flag.")
	}
	if *outputFile == "" {
		log.Fatal("Missing --output_file flag.")
	}

	eniContent, err := os.ReadFile(*eniFile)
	if err != nil {
		log.Fatalf("Failed to read ENI file: %v", err)
	}

	if !utf8.Valid(eniContent) {
		log.Fatalf("ENI file %s is not valid UTF-8.", *eniFile)
	}

	eniMsg := &enipb.Eni{
		Data: string(eniContent),
	}

	any, err := anypb.New(eniMsg)
	if err != nil {
		log.Fatalf("Failed to create Any proto: %v", err)
	}

	manifest := &dmpb.DataManifest{
		Metadata: &dmpb.DataManifest_Metadata{
			Id: &idpb.Id{
				Package: *assetPackage,
				Name:    *assetName,
			},
			Vendor: &vpb.Vendor{
				DisplayName: *vendorDisplayName,
			},
			DisplayName: *displayName,
		},
		Data: any,
	}

	protoText, err := prototext.MarshalOptions{Multiline: true}.Marshal(manifest)
	if err != nil {
		log.Fatalf("Failed to marshal proto to text: %v", err)
	}

	// R/W for owner, R for group, R for others
	err = os.WriteFile(*outputFile, protoText, 0644)
	if err != nil {
		log.Fatalf("Failed to write to output file: %v", err)
	}
}
