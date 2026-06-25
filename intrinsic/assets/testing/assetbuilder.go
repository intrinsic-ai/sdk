// Copyright 2023 Intrinsic Innovation LLC

// Package assetbuilder contains test utilities for building the Asset proto.
package assetbuilder

import (
	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	hdmpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"
	processassetpb "intrinsic/assets/processes/proto/process_asset_go_proto"
	assetpb "intrinsic/assets/proto/v1/asset_go_proto"
	processedassetpb "intrinsic/assets/proto/v1/processed_asset_go_proto"
	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	psmpb "intrinsic/skills/proto/processed_skill_manifest_go_proto"
)

// FromDataAsset creates a local variant of an Asset from a DataAsset.
func FromDataAsset(d *dapb.DataAsset) *assetpb.Asset {
	return &assetpb.Asset{
		Source: &assetpb.Asset_Local{
			Local: &processedassetpb.ProcessedAsset{
				Variant: &processedassetpb.ProcessedAsset_Data{
					Data: d,
				},
			},
		},
	}
}

// FromProcessedHardwareDeviceManifest creates a local variant of an Asset from a ProcessedHardwareDeviceManifest.
func FromProcessedHardwareDeviceManifest(h *hdmpb.ProcessedHardwareDeviceManifest) *assetpb.Asset {
	return &assetpb.Asset{
		Source: &assetpb.Asset_Local{
			Local: &processedassetpb.ProcessedAsset{
				Variant: &processedassetpb.ProcessedAsset_HardwareDevice{
					HardwareDevice: h,
				},
			},
		},
	}
}

// FromProcessAsset creates a local variant of an Asset from a ProcessAsset.
func FromProcessAsset(p *processassetpb.ProcessAsset) *assetpb.Asset {
	return &assetpb.Asset{
		Source: &assetpb.Asset_Local{
			Local: &processedassetpb.ProcessedAsset{
				Variant: &processedassetpb.ProcessedAsset_Process{
					Process: p,
				},
			},
		},
	}
}

// FromProcessedSceneObjectManifest creates a local variant of an Asset from a ProcessedSceneObjectManifest.
func FromProcessedSceneObjectManifest(s *sompb.ProcessedSceneObjectManifest) *assetpb.Asset {
	return &assetpb.Asset{
		Source: &assetpb.Asset_Local{
			Local: &processedassetpb.ProcessedAsset{
				Variant: &processedassetpb.ProcessedAsset_SceneObject{
					SceneObject: s,
				},
			},
		},
	}
}

// FromProcessedServiceManifest creates a local variant of an Asset from a ProcessedServiceManifest.
func FromProcessedServiceManifest(s *smpb.ProcessedServiceManifest) *assetpb.Asset {
	return &assetpb.Asset{
		Source: &assetpb.Asset_Local{
			Local: &processedassetpb.ProcessedAsset{
				Variant: &processedassetpb.ProcessedAsset_Service{
					Service: s,
				},
			},
		},
	}
}

// FromProcessedSkillManifest creates a local variant of an Asset from a ProcessedSkillManifest.
func FromProcessedSkillManifest(s *psmpb.ProcessedSkillManifest) *assetpb.Asset {
	return &assetpb.Asset{
		Source: &assetpb.Asset_Local{
			Local: &processedassetpb.ProcessedAsset{
				Variant: &processedassetpb.ProcessedAsset_Skill{
					Skill: s,
				},
			},
		},
	}
}
