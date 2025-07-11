// Copyright 2023 Intrinsic Innovation LLC

// Package skillserviceconfiggen defines a tool to serialize ServiceConfig proto from a skill manifest.
package skillserviceconfiggen

import (
	"google.golang.org/protobuf/proto"
	"intrinsic/assets/idutils"
	sciv "intrinsic/util/proto/sourcecodeinfoview"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
	sscpb "intrinsic/skills/proto/skill_service_config_go_proto"
	spb "intrinsic/skills/proto/skills_go_proto"
)

func stripSourceCodeInfo(fds *dpb.FileDescriptorSet) {
	for _, file := range fds.GetFile() {
		file.SourceCodeInfo = nil
	}
}

func createParameterDescription(metadata *smpb.ParameterMetadata, skillProtoDescriptor *dpb.FileDescriptorSet) (*spb.ParameterDescription, error) {
	if metadata == nil {
		return nil, nil
	}
	description := new(spb.ParameterDescription)
	description.ParameterMessageFullName = metadata.GetMessageFullName()
	description.DefaultValue = metadata.GetDefaultValue()
	description.ParameterDescriptorFileset = proto.Clone(skillProtoDescriptor).(*dpb.FileDescriptorSet)

	var err error
	description.ParameterFieldComments, err = sciv.NestedFieldCommentMap(skillProtoDescriptor, metadata.GetMessageFullName())
	if err != nil {
		return nil, err
	}

	stripSourceCodeInfo(description.ParameterDescriptorFileset)
	return description, nil
}

func createReturnDescription(metadata *smpb.ReturnMetadata, skillProtoDescriptor *dpb.FileDescriptorSet) (*spb.ReturnValueDescription, error) {
	if metadata == nil {
		return nil, nil
	}
	description := new(spb.ReturnValueDescription)
	description.ReturnValueMessageFullName = metadata.GetMessageFullName()

	description.DescriptorFileset = proto.Clone(skillProtoDescriptor).(*dpb.FileDescriptorSet)
	var err error
	description.ReturnValueFieldComments, err = sciv.NestedFieldCommentMap(skillProtoDescriptor, metadata.GetMessageFullName())
	if err != nil {
		return nil, err
	}

	stripSourceCodeInfo(description.DescriptorFileset)
	return description, nil
}

func buildSkillProto(manifest *smpb.SkillManifest, skillProtoDescriptor *dpb.FileDescriptorSet) (*spb.Skill, error) {
	skill := new(spb.Skill)
	skill.SkillName = manifest.GetId().GetName()
	var err error
	skill.Id, err = idutils.IDFromProto(manifest.GetId())
	if err != nil {
		return nil, err
	}
	skill.PackageName = manifest.GetId().GetPackage()

	// C++ BuildSkillProto API accepts a version argument, but it is never used when generating a
	// SkillServiceConfig. The version is not in the manifest at compile time. If this function is
	// ever promoted to a golang library then a version attribute should be added and
	// skill.IdVersion should be set accordingly.
	skill.IdVersion = skill.GetId()

	skill.Description = manifest.GetDocumentation().GetDescription()
	skill.DisplayName = manifest.GetDisplayName()
	skill.ResourceSelectors = manifest.GetDependencies().GetRequiredEquipment()
	skill.ExecutionOptions = new(spb.ExecutionOptions)
	skill.ExecutionOptions.SupportsCancellation = manifest.GetOptions().GetSupportsCancellation()

	skill.ParameterDescription, err = createParameterDescription(manifest.GetParameter(), skillProtoDescriptor)
	if err != nil {
		return nil, err
	}
	skill.ReturnValueDescription, err = createReturnDescription(manifest.GetReturnType(), skillProtoDescriptor)
	if err != nil {
		return nil, err
	}
	return skill, nil
}

// ExtractSkillServiceConfigFromManifest extracts a SkillServiceConfig from a SkillManifest.
func ExtractSkillServiceConfigFromManifest(manifest *smpb.SkillManifest, skillProtoDescriptor *dpb.FileDescriptorSet) (*sscpb.SkillServiceConfig, error) {
	config := new(sscpb.SkillServiceConfig)

	if manifest.GetOptions().GetCancellationReadyTimeout() != nil {
		config.ExecutionServiceOptions = new(sscpb.ExecutionServiceOptions)
		config.ExecutionServiceOptions.CancellationReadyTimeout = manifest.GetOptions().GetCancellationReadyTimeout()
	}
	if manifest.GetOptions().GetExecutionTimeout() != nil {
		if config.ExecutionServiceOptions == nil {
			config.ExecutionServiceOptions = new(sscpb.ExecutionServiceOptions)
		}
		config.ExecutionServiceOptions.ExecutionTimeout = manifest.GetOptions().GetExecutionTimeout()
	}
	config.StatusInfo = manifest.GetStatusInfo()

	skillDescription, err := buildSkillProto(manifest, skillProtoDescriptor)
	if err != nil {
		return nil, err
	}
	config.SkillDescription = skillDescription
	return config, nil
}
