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

#include "intrinsic/simulation/gazebo/client/joint_state.h"

#include <cstdint>
#include <functional>
#include <memory>
#include <string>
#include <string_view>
#include <utility>

#include "absl/base/nullability.h"
#include "absl/log/log.h"
#include "absl/memory/memory.h"
#include "absl/status/status.h"
#include "absl/synchronization/notification.h"
#include "absl/time/time.h"
#include "gz/msgs/MessageTypes.hh"
#include "gz/transport/Node.hh"
#include "gz/transport/TopicUtils.hh"
#include "intrinsic/simulation/gazebo/client/gazebo_client.h"
#include "intrinsic/util/status/ret_check.h"
#include "intrinsic/util/status/status_builder.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {
namespace simulation {
namespace {

absl::StatusOr<double> GetJointPositionFromJointState(
    const JointState::State& state, uint32_t axis_index) {
  switch (axis_index) {
    case 0u: {
      if (!state.axis_0.has_value()) {
        return NotFoundErrorBuilder()
               << "Joint position for axis index 0 is not found.";
      }
      return state.axis_0->position;
    }
    case 1u: {
      if (!state.axis_1.has_value()) {
        return NotFoundErrorBuilder()
               << "Joint position for axis index 1 is not found.";
      }
      return state.axis_1->position;
    }
    default:
      return InvalidArgumentErrorBuilder()
             << "Joint axis index out of range: " << axis_index;
  }
}

absl::StatusOr<JointState::State> GetJointStateFromMessage(
    const gz::msgs::Model& model_msg) {
  if (model_msg.joint_size() < 1) {
    return NotFoundErrorBuilder()
           << "Joint states are not found on model: " << model_msg.name();
  }

  JointState::State state;
  if (model_msg.joint(0).has_axis1()) {
    state.axis_0 = JointState::State::Axis{
        .position = model_msg.joint(0).axis1().position(),
        .velocity = model_msg.joint(0).axis1().velocity()};
  }
  if (model_msg.joint(0).has_axis2()) {
    state.axis_1 = JointState::State::Axis{
        .position = model_msg.joint(0).axis2().position(),
        .velocity = model_msg.joint(0).axis2().velocity()};
  }
  return state;
}

absl::StatusOr<gz::msgs::Model> WaitForNewMessage(
    GazeboClient* absl_nonnull gazebo_client, const std::string& topic,
    absl::Duration timeout) {
  INTR_RET_CHECK(!topic.empty());

  // Wait for a new message, parse and return the joint position value
  gz::msgs::Model model_msg;
  absl::Notification received_msg;
  std::function<void(const gz::msgs::Model&)> on_joint_state =
      [&model_msg, &received_msg](const gz::msgs::Model& msg) {
        if (!received_msg.HasBeenNotified()) {
          model_msg = msg;
          received_msg.Notify();
        }
      };
  gz::transport::Node::Subscriber sub =
      gazebo_client->Node().CreateSubscriber(topic, on_joint_state);
  if (!sub) {
    return InvalidArgumentErrorBuilder()
           << "Unable to subscribe to joint state topic " << topic;
  }
  received_msg.WaitForNotificationWithTimeout(timeout);
  sub.Unsubscribe();

  if (!received_msg.HasBeenNotified()) {
    return DeadlineExceededErrorBuilder() << "No messages received on "
                                          << "joint state topic: " << topic;
  }

  return model_msg;
}

}  // namespace

absl::StatusOr<std::unique_ptr<JointState>> JointState::Create(
    std::shared_ptr<GazeboClient> gazebo_client, std::string_view topic) {
  if (topic.empty()) {
    return absl::InvalidArgumentError("Joint state topic cannot be empty.");
  }
  if (!gz::transport::TopicUtils::IsValidTopic(std::string(topic))) {
    return InvalidArgumentErrorBuilder()
           << "Joint state topic is invalid" << topic;
  }
  if (gazebo_client == nullptr) {
    return absl::InvalidArgumentError("Gazebo client is null.");
  }
  return absl::WrapUnique(new JointState(std::move(gazebo_client), topic));
}

JointState::JointState(std::shared_ptr<GazeboClient> gazebo_client,
                       std::string_view topic)
    : gazebo_client_(std::move(gazebo_client)), topic_(topic) {}

JointState::~JointState() {
  if (subscriber_) {
    if (absl::Status unsub = Unsubscribe(); !unsub.ok()) {
      LOG(ERROR) << unsub << " in destructor.";
    }
  }
}

std::string_view JointState::Topic() const { return topic_; }

absl::StatusOr<double> JointState::GetJointPosition(absl::Duration timeout,
                                                    uint32_t axis_index) {
  if (axis_index > 1u) {
    return InvalidArgumentErrorBuilder()
           << "Joint axis index out of range: " << axis_index;
  }

  INTR_ASSIGN_OR_RETURN(State state, GetJointState(timeout));
  return GetJointPositionFromJointState(state, axis_index);
}

absl::StatusOr<JointState::State> JointState::GetJointState(
    absl::Duration timeout) {
  INTR_ASSIGN_OR_RETURN(
      gz::msgs::Model model_msg,
      WaitForNewMessage(gazebo_client_.get(), topic_, timeout));
  return GetJointStateFromMessage(model_msg);
}

absl::Status JointState::Subscribe(
    std::function<void(const absl::StatusOr<JointState::State>)> callback) {
  if (subscriber_) {
    return FailedPreconditionErrorBuilder()
           << "An existing subscription on the topic " << topic_
           << "already exists. Please unsubscribe first or subscribe using a "
              "separate JointState object.";
  }

  std::function<void(const gz::msgs::Model&)> on_joint_state =
      [callback](const gz::msgs::Model& msg) {
        callback(GetJointStateFromMessage(msg));
      };

  subscriber_ = gazebo_client_->Node().CreateSubscriber(topic_, on_joint_state);
  if (!subscriber_) {
    return InternalErrorBuilder()
           << "Unable to subscribe to joint state topic " << topic_;
  }
  return absl::OkStatus();
}

absl::Status JointState::Unsubscribe() {
  if (subscriber_ && !subscriber_.Unsubscribe()) {
    return InternalErrorBuilder()
           << "Unable to unsubscribe from joint state topic " << topic_;
  }
  return absl::OkStatus();
}

}  // namespace simulation
}  // namespace intrinsic
