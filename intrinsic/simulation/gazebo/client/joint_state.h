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

#ifndef INTRINSIC_SIMULATION_GAZEBO_CLIENT_JOINT_STATE_H_
#define INTRINSIC_SIMULATION_GAZEBO_CLIENT_JOINT_STATE_H_

#include <cstdint>
#include <functional>
#include <memory>
#include <optional>
#include <string>
#include <string_view>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/time/time.h"
#include "gz/transport/Node.hh"
#include "intrinsic/simulation/gazebo/client/gazebo_client.h"

namespace intrinsic {
namespace simulation {

// JointState class subscribes to joint states published by
// JointStatePublisher plugins in Gazebo.
//
// This class is not thread-safe.
//
// Example:
//
//   absl::StatusOr<std::shared_ptr<GazeboClient>> gazebo_client =
//       GazeboClient::Create("localhost"));
//   if (gazebo_client.ok()) {
//     absl::StatusOr<std::unique_ptr<JointState>> joint_state =
//         JointState::Create(*gazebo_client, "/joint_states");
//     if (joint_state.ok()) {
//       absl::StatusOr<double> joint_pos =
//           joint_state->GetJointPosition(absl::Milliseconds(100), 0u);
//     }
//   }
//
// See topic_utils.h for helper methods to retrieve the joint state topic.
class JointState {
 public:
  // Joint state data for joints with 0 or more axes.
  struct State {
    // Axis state consisting of joint position and joint velocity data.
    // The units for the different types of joints are:
    //   * revolute joint: radians for position and radians per second for
    //   velocity.
    //   * prismatic joint: meters for position and meters per second for
    //   velocity.
    struct Axis {
      double position{0.0};
      double velocity{0.0};
    };
    std::optional<Axis> axis_0;
    std::optional<Axis> axis_1;
  };

  // Creates a JointState object for listening to joint states in simulation
  // on the specified topic.
  //
  // A kInvalidArgument status error is returned if `gazebo_client` is null or
  // the topic is either empty or invalid.
  static absl::StatusOr<std::unique_ptr<JointState>> Create(
      std::shared_ptr<GazeboClient> gazebo_client, std::string_view topic);

  ~JointState();

  // Gets the current position of the joint in simulation.
  //
  // The function subscribes to new joint states and waits for a new message
  // for the specified timeout duration. If a joint state is successfully
  // received, it returns a double value representing the joint position which
  // is in radians for revolute joints and meters for prismatic joints.
  //
  // The axis_index should be 0 for 1-DOF joints and either 0 or 1 for 2-DOF
  // joints.
  //
  // Error statuses returned are:
  //   * kNotFound: Joint states are not found on this joint state topic.
  //   * kInvalidArgument: Specified joint axis index is out of range.
  absl::StatusOr<double> GetJointPosition(absl::Duration timeout,
                                          uint32_t axis_index = 0u);

  // Gets the current state of the joint in simulation.
  //
  // The function subscribes to new joint states and waits for a new message
  // for the specified timeout duration.
  //
  // Error statuses returned are:
  //   * kNotFound: Joint states are not found on this joint state topic.
  absl::StatusOr<State> GetJointState(absl::Duration timeout);

  // Subscribes to joint states from simulation.
  //
  // This is a non-blocking call to stream joint states. The data received from
  // simulation will be passed to the registered callback function in the form
  // of a `JointState::State`.
  //
  // If no joint states are available in the data message received from
  // simulation, the callback argument is set to kNotFound status.
  //
  // A kFailedPrecondition error is returned if subscription failed, which
  // indicates that there is already an existing subscription.
  absl::Status Subscribe(
      std::function<void(const absl::StatusOr<JointState::State>)> callback);

  // Unsubscribes from the joint state topic.
  //
  // This terminates all subscriptions in progress from previous Subscribe
  // calls.
  absl::Status Unsubscribe();

  // Gets the topic that is being subscribed to for joint states.
  std::string_view Topic() const;

 private:
  JointState(std::shared_ptr<GazeboClient> gazebo_client,
             std::string_view topic);
  std::shared_ptr<GazeboClient> gazebo_client_;
  std::string topic_;
  gz::transport::Node::Subscriber subscriber_;
};

}  // namespace simulation
}  // namespace intrinsic

#endif  // INTRINSIC_SIMULATION_GAZEBO_CLIENT_JOINT_STATE_H_
