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

#ifndef INTRINSIC_SIMULATION_GAZEBO_GZ_TOPIC_CONSTANTS_H_
#define INTRINSIC_SIMULATION_GAZEBO_GZ_TOPIC_CONSTANTS_H_

#include <string_view>

namespace intrinsic {
namespace simulation {

// Metadata populated in `GazeboService::ListTopics` gRPC response.
// See intrinsic/simulation/gazebo/proto/v1/gazebo_service.proto
inline constexpr std::string_view kGzTopicInfoJointPositionControlMetadata =
    "joint_position_control";
inline constexpr std::string_view kGzTopicInfoJointStateMetadata =
    "joint_state";

}  // namespace simulation
}  // namespace intrinsic

#endif  // INTRINSIC_SIMULATION_GAZEBO_GZ_TOPIC_CONSTANTS_H_
