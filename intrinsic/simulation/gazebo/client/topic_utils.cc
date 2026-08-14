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

#include "intrinsic/simulation/gazebo/client/topic_utils.h"

#include <algorithm>
#include <cstdint>
#include <memory>
#include <string>
#include <string_view>
#include <utility>

#include "absl/base/nullability.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "intrinsic/simulation/gazebo/client/gazebo_client.h"
#include "intrinsic/simulation/gazebo/client/joint_control.h"
#include "intrinsic/simulation/gazebo/client/joint_state.h"
#include "intrinsic/simulation/gazebo/gz_topic_constants.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {
namespace simulation {

namespace {

bool IsJointControlTopicForAxis(
    const intrinsic_proto::simulation::v1::TopicInfo& topic_info,
    std::string_view joint_name, uint32_t axis_index) {
  return topic_info.has_joint_info() &&
         topic_info.joint_info().name() == joint_name &&
         topic_info.joint_info().has_axis_index() &&
         topic_info.joint_info().axis_index() == axis_index &&
         !topic_info.is_sim_server_advertised() &&
         topic_info.metadata() == kGzTopicInfoJointPositionControlMetadata;
}

bool IsJointStateTopicForAxis(
    const intrinsic_proto::simulation::v1::TopicInfo& topic_info,
    std::string_view joint_name, uint32_t axis_index) {
  return topic_info.has_joint_info() &&
         topic_info.joint_info().name() == joint_name &&
         topic_info.joint_info().has_axis_index() &&
         topic_info.joint_info().axis_index() == axis_index &&
         topic_info.is_sim_server_advertised() &&
         topic_info.metadata() == kGzTopicInfoJointStateMetadata;
}

}  // namespace

absl::StatusOr<std::unique_ptr<JointControl>> CreateJointControl(
    std::string_view object_name,
    std::shared_ptr<GazeboClient> absl_nonnull gazebo_client,
    std::string_view joint_name, uint32_t axis_index) {
  if (gazebo_client == nullptr) {
    return absl::InvalidArgumentError("GazeboClient must not be null.");
  }
  INTR_ASSIGN_OR_RETURN(
      const intrinsic_proto::simulation::v1::Topics list_topics_res,
      gazebo_client->ListTopics(object_name,
                                intrinsic_proto::simulation::v1::MODELS));
  const auto& topics = list_topics_res.topics();
  auto match_it = std::ranges::find_if(
      topics, [&joint_name, &axis_index](const auto& topic_info) {
        return IsJointControlTopicForAxis(topic_info, joint_name, axis_index);
      });
  if (match_it != topics.end()) {
    return JointControl::Create(std::move(gazebo_client),
                                match_it->topic_name());
  }
  return absl::InvalidArgumentError(
      absl::StrFormat("Unable to create Joint Control for [%s] on object [%s]. "
                      "Please check that the object is not a hardware module "
                      "controlled by a realtime control service.",
                      joint_name, object_name));
}

absl::StatusOr<std::unique_ptr<JointState>> CreateJointState(
    std::string_view object_name,
    std::shared_ptr<GazeboClient> absl_nonnull gazebo_client,
    std::string_view joint_name, uint32_t axis_index) {
  if (gazebo_client == nullptr) {
    return absl::InvalidArgumentError("GazeboClient must not be null.");
  }
  INTR_ASSIGN_OR_RETURN(
      const intrinsic_proto::simulation::v1::Topics list_topics_res,
      gazebo_client->ListTopics(object_name,
                                intrinsic_proto::simulation::v1::MODELS));
  const auto& topics = list_topics_res.topics();
  auto match_it = std::ranges::find_if(
      topics, [&joint_name, &axis_index](const auto& topic_info) {
        return IsJointStateTopicForAxis(topic_info, joint_name, axis_index);
      });
  if (match_it != topics.end()) {
    return JointState::Create(std::move(gazebo_client), match_it->topic_name());
  }
  return absl::InvalidArgumentError(
      absl::StrFormat("Unable to create Joint State for [%s] on object [%s]. "
                      "Please check that the object is not a hardware module "
                      "controlled by a realtime control service.",
                      joint_name, object_name));
}

}  // namespace simulation
}  // namespace intrinsic
