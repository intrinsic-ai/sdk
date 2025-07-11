// Copyright 2023 Intrinsic Innovation LLC

// Package assetutils provides utility functions for working with the catalog's Asset proto.
package assetutils

import (
	log "github.com/golang/glog"

	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	hdmpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"
	servicempb "intrinsic/assets/services/proto/service_manifest_go_proto"
	ipb "intrinsic/kubernetes/workcell_spec/proto/image_go_proto"
)

// CollectImages returns the images included in the given Asset.
//
// NOTE that it does not include images from catalog-referenced assets.
func CollectImages(asset *acpb.Asset) []*ipb.Image {
	var images []*ipb.Image
	switch dd := asset.GetDeploymentData().GetAssetSpecificDeploymentData().(type) {
	case *acpb.Asset_AssetDeploymentData_HardwareDeviceSpecificDeploymentData:
		for _, asset := range dd.HardwareDeviceSpecificDeploymentData.GetManifest().GetAssets() {
			switch asset.GetVariant().(type) {
			case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Service:
				images = append(images, collectServiceImages(asset.GetService())...)
			case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Catalog:
			case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Data:
			case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_SceneObject:
			default:
				log.Fatalf("unknown asset type in HardwareDevice: %v", asset.GetVariant())
			}
		}
	case *acpb.Asset_AssetDeploymentData_ServiceSpecificDeploymentData:
		images = collectServiceImages(dd.ServiceSpecificDeploymentData.GetManifest())
	case *acpb.Asset_AssetDeploymentData_SkillSpecificDeploymentData:
		m := dd.SkillSpecificDeploymentData.GetManifest()
		if m.GetAssets().GetImage() != nil {
			images = append(images, m.GetAssets().GetImage())
		}
	case *acpb.Asset_AssetDeploymentData_DataSpecificDeploymentData:
	case *acpb.Asset_AssetDeploymentData_SceneObjectSpecificDeploymentData:
	default:
		log.Fatalf("unknown asset type: %v", asset.GetMetadata().GetAssetType())
	}
	return images
}

func collectServiceImages(m *servicempb.ProcessedServiceManifest) []*ipb.Image {
	images := make([]*ipb.Image, 0, len(m.GetAssets().GetImages()))
	for _, v := range m.GetAssets().GetImages() {
		images = append(images, v)
	}
	return images
}
