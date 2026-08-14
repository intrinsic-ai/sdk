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

#include "intrinsic/simulation/gazebo/client/joint_control.h"

#include <memory>
#include <string>
#include <string_view>
#include <utility>

#include "absl/log/log.h"
#include "absl/memory/memory.h"
#include "absl/status/status.h"
#include "gz/msgs/MessageTypes.hh"
#include "gz/transport/Node.hh"
#include "intrinsic/simulation/gazebo/client/gazebo_client.h"
#include "intrinsic/util/status/status_builder.h"

namespace intrinsic {
namespace simulation {

absl::StatusOr<std::unique_ptr<JointControl>> JointControl::Create(
    std::shared_ptr<GazeboClient> gazebo_client, std::string_view topic) {
  if (topic.empty()) {
    return absl::InvalidArgumentError("Joint control topic cannot be empty.");
  }
  if (gazebo_client == nullptr) {
    return absl::InvalidArgumentError("Gazebo client is null");
  }
  gz::transport::Node::Publisher pub =
      gazebo_client->Node().Advertise<gz::msgs::Double>(std::string(topic));
  if (!pub.Valid()) {
    return InvalidArgumentErrorBuilder()
           << "Failed to create publisher for joint control topic " << topic;
  }
  return absl::WrapUnique(
      new JointControl(std::move(gazebo_client), topic, std::move(pub)));
}

JointControl::JointControl(std::shared_ptr<GazeboClient> gazebo_client,
                           std::string_view topic,
                           gz::transport::Node::Publisher pub)
    // keep a shared pointer to gazebo_client_ as a member variable to ensure
    // that the publisher remains valid
    : gazebo_client_(std::move(gazebo_client)),
      topic_(topic),
      pub_(std::move(pub)) {}

std::string_view JointControl::Topic() const { return topic_; }

absl::Status JointControl::SetJointPosition(double position) {
  if (!pub_.HasConnections()) {
    return AbortedErrorBuilder()
           << "Joint control topic " << topic_
           << " has no connections. Ensure that gzserver is running with a "
              "position controller plugin subscribing to this topic.";
  }
  gz::msgs::Double joint_pos_msg;
  joint_pos_msg.set_data(position);
  return (pub_.Publish(joint_pos_msg))
             ? absl::OkStatus()
             : InternalErrorBuilder()
                   << "Unable to publish joint control command to topic "
                   << topic_;
}

}  // namespace simulation
}  // namespace intrinsic
