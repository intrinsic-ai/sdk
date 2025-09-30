// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/skills/internal/skill_service_config_utils.h"

#include <optional>

#include "absl/log/log.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "google/protobuf/descriptor.pb.h"
#include "intrinsic/assets/proto/id.pb.h"
#include "intrinsic/assets/proto/status_spec.pb.h"
#include "intrinsic/icon/release/file_helpers.h"
#include "intrinsic/skills/internal/skill_proto_utils.h"
#include "intrinsic/skills/proto/skill_manifest.pb.h"
#include "intrinsic/skills/proto/skill_service_config.pb.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::skills {

void LogLimitedSkillServiceConfig(
    const intrinsic_proto::skills::SkillServiceConfig& config) {
  // Copy to modify.
  intrinsic_proto::skills::SkillServiceConfig service_config_copy = config;
  service_config_copy.mutable_skill_description()
      ->mutable_parameter_description()
      ->clear_parameter_descriptor_fileset();
  service_config_copy.mutable_skill_description()
      ->mutable_return_value_description()
      ->clear_descriptor_fileset();
  LOG(INFO) << absl::StrCat(
      "\nUsing skill configuration proto (minus the descriptor filesets):\n",
      service_config_copy);
}

absl::StatusOr<intrinsic_proto::skills::SkillServiceConfig>
GetSkillServiceConfig(absl::string_view skill_service_config_filename) {
  intrinsic_proto::skills::SkillServiceConfig service_config;
  if (!skill_service_config_filename.empty()) {
    LOG(INFO) << "Reading skill configuration proto from "
              << skill_service_config_filename;
    INTR_ASSIGN_OR_RETURN(
        service_config,
        intrinsic::GetBinaryProto<intrinsic_proto::skills::SkillServiceConfig>(
            skill_service_config_filename));
  }
  return service_config;
}

absl::StatusOr<intrinsic_proto::skills::SkillServiceConfig>
GetSkillServiceConfigFromManifest(
    absl::string_view manifest_pbbin_filename,
    absl::string_view file_descriptor_set_pbbin_filename,
    std::optional<absl::string_view> version) {
  intrinsic_proto::skills::SkillServiceConfig service_config;

  LOG(INFO) << "Loading manifest from " << manifest_pbbin_filename;
  INTR_ASSIGN_OR_RETURN(
      auto manifest,
      intrinsic::GetBinaryProto<intrinsic_proto::skills::SkillManifest>(
          manifest_pbbin_filename));

  LOG(INFO) << "Loading FileDescriptorSet from "
            << file_descriptor_set_pbbin_filename;
  INTR_ASSIGN_OR_RETURN(
      auto file_descriptor_set,
      intrinsic::GetBinaryProto<google::protobuf::FileDescriptorSet>(
          file_descriptor_set_pbbin_filename));

  return GetSkillServiceConfigFromManifest(manifest, file_descriptor_set,
                                           version);
}

absl::StatusOr<intrinsic_proto::skills::SkillServiceConfig>
GetSkillServiceConfigFromManifest(
    const intrinsic_proto::skills::SkillManifest& manifest,
    const google::protobuf::FileDescriptorSet& file_descriptor_set,
    std::optional<absl::string_view> version) {
  intrinsic_proto::skills::SkillServiceConfig service_config;

  if (manifest.options().has_cancellation_ready_timeout()) {
    *service_config.mutable_execution_service_options()
         ->mutable_cancellation_ready_timeout() =
        manifest.options().cancellation_ready_timeout();
  }
  if (manifest.options().has_execution_timeout()) {
    *service_config.mutable_execution_service_options()
         ->mutable_execution_timeout() = manifest.options().execution_timeout();
  }

  *service_config.mutable_status_info() = manifest.status_info();

  INTR_ASSIGN_OR_RETURN(
      *service_config.mutable_skill_description(),
      BuildSkillProto(manifest, file_descriptor_set, version));

  return service_config;
}

absl::StatusOr<intrinsic_proto::skills::SkillServiceConfig>
GetSkillServiceConfigFromManifest(
    const intrinsic_proto::skills::SkillManifest& manifest,
    const google::protobuf::FileDescriptorSet& parameter_descriptor_set,
    const google::protobuf::FileDescriptorSet& return_type_descriptor_set,
    std::optional<absl::string_view> version) {
  intrinsic_proto::skills::SkillServiceConfig service_config;

  if (manifest.options().has_cancellation_ready_timeout()) {
    *service_config.mutable_execution_service_options()
         ->mutable_cancellation_ready_timeout() =
        manifest.options().cancellation_ready_timeout();
  }
  if (manifest.options().has_execution_timeout()) {
    *service_config.mutable_execution_service_options()
         ->mutable_execution_timeout() = manifest.options().execution_timeout();
  }

  *service_config.mutable_status_info() = manifest.status_info();

  INTR_ASSIGN_OR_RETURN(
      *service_config.mutable_skill_description(),
      BuildSkillProto(manifest, parameter_descriptor_set,
                      return_type_descriptor_set,
                      version));

  return service_config;
}

}  // namespace intrinsic::skills
