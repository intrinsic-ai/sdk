// Copyright 2023 Intrinsic Innovation LLC

// assetcatalogrefinfogen creates an AssetCatalogRefInfo proto which contains info about a catalog asset.
package main

import (
	"flag"
	log "github.com/golang/glog"
	"intrinsic/assets/idutils"
	"intrinsic/assets/typeutils"
	intrinsic "intrinsic/production/intrinsic"
	intrinsicflag "intrinsic/util/flag"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	assetpb "intrinsic/assets/build_defs/asset_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
)

var (
	assetType                 = flag.String("asset_type", "", "The type of asset.")
	id                        = flag.String("id", "", "The id of the catalog asset.")
	version                   = flag.String("version", "", "The version of the catalog asset.")
	fileDescriptorSets        = intrinsicflag.MultiString("file_descriptor_set", nil, "Path to a binary file descriptor set proto to be used to resolve the data payload. Can be repeated.")
	outputAssetInfo           = flag.String("output_asset_info", "", "Output AssetInfo proto path.")
	outputAssetCatalogRefInfo = flag.String("output_asset_catalog_ref_info", "", "Output AssetCatalogRefInfo proto path.")
)

func writeAsset(idVersion *idpb.IdVersion, fds *dpb.FileDescriptorSet) {
	atype := typeutils.AssetTypeFromCodeName(*assetType)
	if atype == atypepb.AssetType_ASSET_TYPE_UNSPECIFIED {
		log.Exitf("unknown asset type %q", *assetType)
	}

	switch atype {
	case atypepb.AssetType_ASSET_TYPE_DATA:
	case atypepb.AssetType_ASSET_TYPE_SCENE_OBJECT:
	case atypepb.AssetType_ASSET_TYPE_SERVICE:
	case atypepb.AssetType_ASSET_TYPE_SKILL:
	case atypepb.AssetType_ASSET_TYPE_HARDWARE_DEVICE:
	default:
		log.Exitf("unsupported asset type %q", *assetType)
	}

	if err := protoio.WriteBinaryProto(*outputAssetInfo, &assetpb.AssetInfo{
		AssetType:         atype,
		Id:                idVersion.GetId(),
		FileDescriptorSet: fds,
	}, protoio.WithDeterministic(true)); err != nil {
		log.Exitf("Could not write asset info: %v", err)
	}
	if err := protoio.WriteBinaryProto(*outputAssetCatalogRefInfo, &assetpb.AssetCatalogRefInfo{
		AssetType: atype,
		IdVersion: idVersion,
	}, protoio.WithDeterministic(true)); err != nil {
		log.Exitf("Could not write asset catalog ref info: %v", err)
	}
}

func main() {
	intrinsic.Init()

	fds, err := registryutil.LoadFileDescriptorSets(*fileDescriptorSets)
	if err != nil {
		log.Exitf("cannot build file descriptor set for asset: %v", err)
	}
	pkg, err := idutils.PackageFrom(*id)
	if err != nil {
		log.Exitf("invalid asset idversion: %v", err)
	}
	name, err := idutils.NameFrom(*id)
	if err != nil {
		log.Exitf("invalid asset idversion: %v", err)
	}
	idVersion, err := idutils.IDVersionProtoFrom(pkg, name, *version)
	if err != nil {
		log.Exitf("invalid asset idversion: %v", err)
	}
	writeAsset(idVersion, fds)
}
