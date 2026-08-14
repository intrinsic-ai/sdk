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

#ifndef INTRINSIC_SIMULATION_GAZEBO_CLIENT_TOPIC_UTILS_H_
#define INTRINSIC_SIMULATION_GAZEBO_CLIENT_TOPIC_UTILS_H_

#include <cstdint>
#include <memory>
#include <string_view>

#include "absl/base/nullability.h"
#include "absl/status/statusor.h"
#include "intrinsic/simulation/gazebo/client/gazebo_client.h"
#include "intrinsic/simulation/gazebo/client/joint_control.h"
#include "intrinsic/simulation/gazebo/client/joint_state.h"

namespace intrinsic {
namespace simulation {

// Creates `JointControl` for a world object.
//
// Discovers the control topic name by querying the Gazebo server, then creates
// a `JointControl` instance on that topic, to command the specified joint axis
// in simulation.
//
// Returns the following errors if unsuccessful.
// - kInvalidArgument if a matching joint topic cannot be found or if
//   `gazebo_client` is null.
// - kUnavailable if the Gazebo server cannot be reached.
absl::StatusOr<std::unique_ptr<JointControl>> CreateJointControl(
    std::string_view object_name,
    std::shared_ptr<GazeboClient> absl_nonnull gazebo_client,
    std::string_view joint_name, uint32_t axis_index);

// Creates `JointState` for a world object.
//
// Discovers the state topic name by querying the Gazebo server, then creates
// a `JointState` instance on that topic, to receive joint state for the
// specified joint axis in simulation.
//
// Returns the following errors if unsuccessful.
// - kInvalidArgument if a matching joint topic cannot be found or if
//   `gazebo_client` is null.
// - kUnavailable if the Gazebo server cannot be reached.
absl::StatusOr<std::unique_ptr<JointState>> CreateJointState(
    std::string_view object_name,
    std::shared_ptr<GazeboClient> absl_nonnull gazebo_client,
    std::string_view joint_name, uint32_t axis_index);

}  // namespace simulation
}  // namespace intrinsic

#endif  // INTRINSIC_SIMULATION_GAZEBO_CLIENT_TOPIC_UTILS_H_
