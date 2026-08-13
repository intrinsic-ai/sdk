// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// assetcatalogrefinfogen creates an AssetCatalogRefInfo proto which contains info about a catalog asset.
package main

import (
	"flag"

	"intrinsic/assets/idutils"
	"intrinsic/assets/typeutils"
	"intrinsic/production/intrinsic"
	intrinsicflag "intrinsic/util/flag"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"

	log "github.com/golang/glog"

	assetpb "intrinsic/assets/build_defs/asset_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"

	dpb "google.golang.org/protobuf/types/descriptorpb"
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
	case atypepb.AssetType_ASSET_TYPE_PROCESS:
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
