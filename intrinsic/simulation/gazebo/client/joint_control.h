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

#ifndef INTRINSIC_SIMULATION_GAZEBO_CLIENT_JOINT_CONTROL_H_
#define INTRINSIC_SIMULATION_GAZEBO_CLIENT_JOINT_CONTROL_H_

#include <cstdint>
#include <memory>
#include <string>
#include <string_view>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "gz/transport/Node.hh"
#include "intrinsic/simulation/gazebo/client/gazebo_client.h"

namespace intrinsic {
namespace simulation {

// JointControl class publishes joint position commands to joint controller
// plugins in Gazebo, e.g. JointPositionController.
//
// Example:
//
//   absl::StatusOr<std::shared_ptr<GazeboClient>> gazebo_client =
//       GazeboClient::Create("localhost"));
//   if (gazebo_client.ok()) {
//     absl::StatusOr<std::unique_ptr<JointControl>> joint_control =
//         JointControl::Create(*gazebo_client, "/joint_cmds");
//     if (joint_control.ok()) {
//       absl::Status result = joint_control->SetJointPosition(1.0);
//     }
//   }
//
// See topic_utils.h for helper methods to retrieve the joint control topic.
class JointControl {
 public:
  // Create a JointControl object for publishing joint position commands to
  // simulation on the specified topic.
  //
  // A kInvalidArgument status error is returned if the `gazebo_client` is null,
  // the topic is empty, or a publisher can not be created for the topic.
  static absl::StatusOr<std::unique_ptr<JointControl>> Create(
      std::shared_ptr<GazeboClient> gazebo_client, std::string_view topic);

  // Set the joint position target by publishing a joint position command to
  // simulation.
  //
  // Joint position is in radians for revolute joints and meters for prismatic
  // joints. This is a position 'target'. Depending on joint configurations
  // (velocity and effort limits, etc) the joint may not move immediately to the
  // target position.
  //
  // A kAborted status error is returned if no subscribers exist on the topic.
  absl::Status SetJointPosition(double position);

  // Get the topic that the joint position commands are published to.
  std::string_view Topic() const;

 private:
  JointControl(std::shared_ptr<GazeboClient> gazebo_client,
               std::string_view topic, gz::transport::Node::Publisher pub);
  std::shared_ptr<GazeboClient> gazebo_client_;
  std::string topic_;
  gz::transport::Node::Publisher pub_;
};

}  // namespace simulation
}  // namespace intrinsic

#endif  // INTRINSIC_SIMULATION_GAZEBO_CLIENT_JOINT_CONTROL_H_
