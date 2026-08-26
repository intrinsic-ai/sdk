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

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <functional>
#include <memory>
#include <string>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "absl/synchronization/notification.h"
#include "absl/time/time.h"
#include "gz/msgs/MessageTypes.hh"
#include "gz/transport/Node.hh"
#include "internal/testing.h"
#include "intrinsic/simulation/gazebo/client/gazebo_client.h"
#include "intrinsic/simulation/gazebo/client/gazebo_client_test_env.h"

using ::absl_testing::StatusIs;

namespace intrinsic {
namespace simulation {
namespace {

TEST(JointControl, CreateFailsWithNullGazeboClient) {
  EXPECT_THAT(JointControl::Create(nullptr, ""),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(JointControl, SetJointPosition) {
  GazeboClientTestEnv env;

  const std::string topic = "test_cmd_topic";
  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());
  gz::transport::Node& node = gazebo_client->Node();
  ASSERT_OK_AND_ASSIGN(std::unique_ptr<JointControl> joint_control,
                       JointControl::Create(gazebo_client, topic));
  EXPECT_EQ(joint_control->Topic(), topic);

  // Test sending a joint position command from JointControl and verify
  // that a msg is received on the associated gz topic.
  absl::Notification received_msg;
  double received_joint_pos_cmd = 0.0;
  std::function<void(const gz::msgs::Double& msg)> cb =
      [&received_joint_pos_cmd, &received_msg](const gz::msgs::Double& msg) {
        if (!received_msg.HasBeenNotified()) {
          received_joint_pos_cmd = msg.data();
          received_msg.Notify();
        }
      };
  node.Subscribe(std::string(joint_control->Topic()), cb);

  double joint_pos_cmd = 0.1;
  EXPECT_OK(joint_control->SetJointPosition(joint_pos_cmd));

  absl::Duration timeout = absl::Seconds(5);
  received_msg.WaitForNotificationWithTimeout(timeout);
  EXPECT_DOUBLE_EQ(received_joint_pos_cmd, joint_pos_cmd);
}

}  // namespace
}  // namespace simulation
}  // namespace intrinsic
