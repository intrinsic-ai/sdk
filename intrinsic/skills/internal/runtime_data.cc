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

#include "intrinsic/skills/internal/runtime_data.h"

#include <iterator>
#include <optional>
#include <string>
#include <utility>
#include <vector>

#include "absl/algorithm/container.h"
#include "absl/container/flat_hash_map.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "absl/types/span.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/descriptor.h"
#include "google/protobuf/descriptor.pb.h"
#include "google/protobuf/descriptor_database.h"
#include "intrinsic/assets/proto/status_spec.pb.h"
#include "intrinsic/skills/proto/equipment.pb.h"
#include "intrinsic/skills/proto/skill_service_config.pb.h"
#include "intrinsic/skills/proto/skills.pb.h"
#include "intrinsic/skills/skill_pub_topic.h"  // intrinsic:skills_pubsub:strip(b/224840414)
#include "intrinsic/util/proto_time.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::skills::internal {
namespace {
// intrinsic:skills_pubsub:strip_begin(b/224840414)
absl::StatusOr<std::vector<SkillPubTopic>> GetSkillPubTopic(
    const intrinsic_proto::skills::PubTopicDescription& description) {
  google::protobuf::SimpleDescriptorDatabase descriptor_database;
  google::protobuf::DescriptorPool descriptor_pool(&descriptor_database);
  for (const auto& file : description.file_descriptor_set().file()) {
    if (!descriptor_database.Add(file)) {
      return absl::InvalidArgumentError(
          "`file_descriptor_set` contains duplicate files.");
    }
  }

  std::vector<SkillPubTopic> pub_topics;
  pub_topics.reserve(description.pub_topics().size());
  for (const auto& pub_topic_proto : description.pub_topics()) {
    if (!descriptor_pool.FindMessageTypeByName(
            pub_topic_proto.message_full_name())) {
      return absl::InvalidArgumentError(absl::StrFormat(
          "missing pub_topic_descriptor for topic with data_id: %s. expected "
          "descriptor for message: %s",
          pub_topic_proto.data_id(), pub_topic_proto.message_full_name()));
    }
    pub_topics.push_back(SkillPubTopic{
        .data_id = pub_topic_proto.data_id(),
        .description = pub_topic_proto.description(),
        .message_full_name = pub_topic_proto.message_full_name(),
    });
  }
  return pub_topics;
}
// intrinsic:skills_pubsub:strip_end
}  // namespace

ParameterData::ParameterData(const google::protobuf::Any& default_value)
    : default_(default_value) {}

ExecutionOptions::ExecutionOptions(
    bool supports_cancellation,
    std::optional<absl::Duration> cancellation_ready_timeout,
    std::optional<absl::Duration> execution_timeout)
    : supports_cancellation_(supports_cancellation) {
  if (cancellation_ready_timeout.has_value()) {
    cancellation_ready_timeout_ = *cancellation_ready_timeout;
  }
  if (execution_timeout.has_value()) {
    execution_timeout_ = *execution_timeout;
  }
}

ResourceData::ResourceData(
    const absl::flat_hash_map<std::string,
                              intrinsic_proto::skills::ResourceSelector>&
        resources_required)
    : resources_required_(resources_required) {}

StatusSpecs::StatusSpecs(
    const std::vector<intrinsic_proto::assets::StatusSpec>& specs) {
  absl::c_transform(specs, std::inserter(specs_, specs_.end()),
                    [](const intrinsic_proto::assets::StatusSpec& spec) {
                      return std::make_pair(spec.code(), spec);
                    });
}

// intrinsic:skills_pubsub:strip_begin(b/224840414)
TopicData::TopicData(absl::Span<const SkillPubTopic> pub_topics)
    : pub_topics_(pub_topics.begin(), pub_topics.end()) {}
// intrinsic:skills_pubsub:strip_end

SkillRuntimeData::SkillRuntimeData(
    const ParameterData& parameter_data,
    const ExecutionOptions& execution_options,
    const ResourceData& resource_data, const StatusSpecs& status_specs,
    const TopicData& topic_data,  // intrinsic:skills_pubsub:strip(b/224840414)
    absl::string_view id)
    : parameter_data_(parameter_data),
      execution_options_(execution_options),
      resource_data_(resource_data),
      status_specs_(status_specs),
      topic_data_(topic_data),  // intrinsic:skills_pubsub:strip(b/224840414)
      id_(id) {}

absl::StatusOr<SkillRuntimeData> GetRuntimeDataFrom(
    const intrinsic_proto::skills::SkillServiceConfig& skill_service_config) {
  // intrinsic:skills_pubsub:strip_begin(b/224840414)
  INTR_ASSIGN_OR_RETURN(
      std::vector<SkillPubTopic> pub_topics,
      GetSkillPubTopic(
          skill_service_config.skill_description().pub_topic_description()));
  // intrinsic:skills_pubsub:strip_end

  std::optional<absl::Duration> cancellation_ready_timeout;
  if (skill_service_config.execution_service_options()
          .has_cancellation_ready_timeout()) {
    cancellation_ready_timeout = ToAbslDurationNoValidation(
        skill_service_config.execution_service_options()
            .cancellation_ready_timeout());
  }

  std::optional<absl::Duration> execution_timeout;
  if (skill_service_config.execution_service_options()
          .has_execution_timeout()) {
    execution_timeout = ToAbslDurationNoValidation(
        skill_service_config.execution_service_options().execution_timeout());
  }

  return SkillRuntimeData(
      skill_service_config.skill_description()
              .parameter_description()
              .has_default_value()
          ? ParameterData(skill_service_config.skill_description()
                              .parameter_description()
                              .default_value())
          : ParameterData(),
      ExecutionOptions(skill_service_config.skill_description()
                           .execution_options()
                           .supports_cancellation(),
                       cancellation_ready_timeout, execution_timeout),
      ResourceData({skill_service_config.skill_description()
                        .resource_selectors()
                        .begin(),
                    skill_service_config.skill_description()
                        .resource_selectors()
                        .end()}),
      StatusSpecs({skill_service_config.status_info().begin(),
                   skill_service_config.status_info().end()}),
      // intrinsic:skills_pubsub:strip_begin(b/224840414)
      TopicData(pub_topics),
      // intrinsic:skills_pubsub:strip_end
      skill_service_config.skill_description().id());
}

}  // namespace intrinsic::skills::internal
