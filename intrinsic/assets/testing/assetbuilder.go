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
