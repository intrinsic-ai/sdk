// Copyright 2023 Intrinsic Innovation LLC

// Package platform contains utilities that list the interfaces an Asset provides to the
// platform. These interfaces can be used to determine whether an Asset is compatible with
// a given platform version.
package platform

import (
	"intrinsic/assets/interfaceutils"

	hdmpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	drpb "intrinsic/assets/services/proto/v1/dynamic_reconfiguration_go_proto"
	sspb "intrinsic/assets/services/proto/v1/service_state_go_proto"
	pskmpb "intrinsic/skills/proto/processed_skill_manifest_go_proto"
	skmpb "intrinsic/skills/proto/skill_manifest_go_proto"
)

// ProvidedByProcessedSkillManifest lists the interfaces the Skill provides to the platform.
func ProvidedByProcessedSkillManifest(manifest *pskmpb.ProcessedSkillManifest) []*metadatapb.Interface {
	if manifest == nil {
		return nil
	}
	return providedBySkillOptions(manifest.GetDetails().GetOptions())
}

// ProvidedBySkillManifest lists the interfaces the Skill provides to the platform.
func ProvidedBySkillManifest(manifest *skmpb.SkillManifest) []*metadatapb.Interface {
	if manifest == nil {
		return nil
	}
	return providedBySkillOptions(manifest.GetOptions())
}

func providedBySkillOptions(options *skmpb.Options) []*metadatapb.Interface {
	if options == nil {
		return nil
	}
	var interfaces []*metadatapb.Interface
	for _, v := range options.GetSkillServicesConfig().GetServiceVersions() {
		var name string
		switch v {
		case skmpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_PROJECTOR:
			name = interfaceutils.GRPCURIPrefix + "intrinsic_proto.skills.Projector"
		case skmpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_EXECUTOR:
			name = interfaceutils.GRPCURIPrefix + "intrinsic_proto.skills.Executor"
		case skmpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_SKILL_INFORMATION:
			name = interfaceutils.GRPCURIPrefix + "intrinsic_proto.skills.SkillInformation"
		default:
			continue
		}
		interfaces = append(interfaces, &metadatapb.Interface{
			Uri: name,
		})
	}
	return interfaces
}

// ProvidedByProcessedServiceManifest lists the interfaces the Service provides to the platform.
func ProvidedByProcessedServiceManifest(manifest *smpb.ProcessedServiceManifest) []*metadatapb.Interface {
	if manifest == nil {
		return nil
	}
	return providedByServiceDef(manifest.GetServiceDef())
}

// ProvidedByServiceManifest lists the interfaces the Service provides to the platform.
func ProvidedByServiceManifest(manifest *smpb.ServiceManifest) []*metadatapb.Interface {
	if manifest == nil {
		return nil
	}
	return providedByServiceDef(manifest.GetServiceDef())
}

func providedByServiceDef(serviceDef *smpb.ServiceDef) []*metadatapb.Interface {
	if serviceDef == nil {
		return nil
	}
	var interfaces []*metadatapb.Interface
	for _, v := range serviceDef.GetDynamicReconfigurationConfig().GetServiceVersions() {
		if v == drpb.DynamicReconfigurationConfig_INTRINSIC_PROTO_SERVICES_V1_DYNAMIC_RECONFIGURATION {
			interfaces = append(interfaces, &metadatapb.Interface{
				Uri: interfaceutils.GRPCURIPrefix + "intrinsic_proto.services.v1.DynamicReconfiguration",
			})
		}
	}
	for _, v := range serviceDef.GetServiceStateConfig().GetServiceVersions() {
		if v == sspb.ServiceStateConfig_INTRINSIC_PROTO_SERVICES_V1_SERVICE_STATE {
			interfaces = append(interfaces, &metadatapb.Interface{
				Uri: interfaceutils.GRPCURIPrefix + "intrinsic_proto.services.v1.ServiceState",
			})
		}
	}
	return interfaces
}

// ProvidedByProcessedHardwareDeviceManifest lists the interfaces a Hardware Device provides to the
// platform. Note that interfaces from non-inlined service Assets are excluded. Prefer
// ProvidedByCollectedAssets when possible.
func ProvidedByProcessedHardwareDeviceManifest(manifest *hdmpb.ProcessedHardwareDeviceManifest) []*metadatapb.Interface {
	if manifest == nil {
		return nil
	}
	var interfaces []*metadatapb.Interface
	for _, pa := range manifest.GetAssets() {
		// Note that non-inlined Services that are stored in the catalog will not have their interfaces
		// included in the output. Please use ProvidedByCollectedAssets instead.
		if m := pa.GetService(); m != nil {
			interfaces = append(interfaces, ProvidedByProcessedServiceManifest(m)...)
		}
	}
	return interfaces
}
