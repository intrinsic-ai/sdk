// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SKILLS_INTERNAL_SKILL_SERVICE_CONFIG_UTILS_H_
#define INTRINSIC_SKILLS_INTERNAL_SKILL_SERVICE_CONFIG_UTILS_H_

#include <optional>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "google/protobuf/descriptor.pb.h"
#include "intrinsic/skills/proto/skill_manifest.pb.h"
#include "intrinsic/skills/proto/skill_service_config.pb.h"

namespace intrinsic::skills {

// Logs a limited subset of the SkillServiceConfig.
//
// In practice, the proto file descriptors can be very large and this makes it
// hard to scroll through the logs. This function removes the descriptors and
// logs the remainder of the proto.
void LogLimitedSkillServiceConfig(
    const intrinsic_proto::skills::SkillServiceConfig& config);

// Returns the SkillServiceConfig at `skill_service_config_filename`.
//
// Reads the proto binary file at `skill_service_config_filename` and returns
// the contents. This file must contain a proto binary
// intrinsic_proto.skills.SkillServiceConfig message.
absl::StatusOr<intrinsic_proto::skills::SkillServiceConfig>
GetSkillServiceConfig(absl::string_view skill_service_config_filename);

// Returns the SkillServiceConfig for the skill with the given manifest.
absl::StatusOr<intrinsic_proto::skills::SkillServiceConfig>
GetSkillServiceConfigFromManifest(
    absl::string_view manifest_pbbin_filename,
    absl::string_view file_descriptor_set_pbbin_filename,
    std::optional<absl::string_view> version = std::nullopt);
absl::StatusOr<intrinsic_proto::skills::SkillServiceConfig>
GetSkillServiceConfigFromManifest(
    const intrinsic_proto::skills::SkillManifest& manifest,
    const google::protobuf::FileDescriptorSet& file_descriptor_set,
    std::optional<absl::string_view> version = std::nullopt);
absl::StatusOr<intrinsic_proto::skills::SkillServiceConfig>
GetSkillServiceConfigFromManifest(
    const intrinsic_proto::skills::SkillManifest& manifest,
    const google::protobuf::FileDescriptorSet& parameter_descriptor_set,
    const google::protobuf::FileDescriptorSet& return_type_descriptor_set,
    std::optional<absl::string_view> version = std::nullopt);

}  // namespace intrinsic::skills

#endif  // INTRINSIC_SKILLS_INTERNAL_SKILL_SERVICE_CONFIG_UTILS_H_
