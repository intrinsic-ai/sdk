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

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <cstdint>
#include <functional>
#include <memory>
#include <optional>
#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/synchronization/notification.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "gz/msgs/MessageTypes.hh"
#include "gz/transport/Node.hh"
#include "intrinsic/simulation/gazebo/client/gazebo_client.h"
#include "intrinsic/simulation/gazebo/client/gazebo_client_test_env.h"
#include "intrinsic/util/testing/gtest_wrapper.h"
#include "intrinsic/util/thread/stop_token.h"
#include "intrinsic/util/thread/thread.h"

using ::absl_testing::StatusIs;
using ::testing::FieldsAre;
using ::testing::Optional;

namespace intrinsic {
namespace simulation {
namespace {

constexpr double kPublishedJointPos[2]{0.2, -1.2};
constexpr double kPublishedJointVel[2]{3.3, -1.4};
constexpr uint32_t kTargetPubMsgCount = 5u;

TEST(JointState, CreateFailsWithNullGazeboClient) {
  EXPECT_THAT(JointState::Create(nullptr, ""),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(JointState, GetJointPosition) {
  const std::string topic = "test_state_topic";
  GazeboClientTestEnv env;
  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());
  ASSERT_OK_AND_ASSIGN(std::unique_ptr<JointState> joint_state,
                       JointState::Create(gazebo_client, topic));
  EXPECT_EQ(joint_state->Topic(), topic);
  // No messages are published so there should not be any messages received.
  EXPECT_THAT(joint_state->GetJointPosition(absl::Milliseconds(100), 0u),
              StatusIs(absl::StatusCode::kDeadlineExceeded));

  // Test publishing a joint state msg and verify that JointState can retrieve
  // the joint position
  // Wait for connection from `joint_state` and publish model msg.
  intrinsic::Thread publisher_thread([&gazebo_client, topic](StopToken st) {
    gz::transport::Node::Publisher pub =
        gazebo_client->Node().Advertise<gz::msgs::Model>(topic);
    gz::msgs::Model msg;
    msg.add_joint()->mutable_axis1()->set_position(kPublishedJointPos[0]);
    while (!st.stop_requested() && !pub.HasConnections()) {
      absl::SleepFor(absl::Milliseconds(10));
    }
    if (pub.HasConnections()) {
      ASSERT_TRUE(pub.Publish(msg));
    }
  });
  ASSERT_OK_AND_ASSIGN(double joint_pos_received,
                       joint_state->GetJointPosition(absl::Seconds(5), 0u));
  EXPECT_DOUBLE_EQ(joint_pos_received, kPublishedJointPos[0]);
}

TEST(JointState, GetJointPositionFromEmptyJointStateMessage) {
  const std::string topic = "test_state_topic_joint_pos";
  GazeboClientTestEnv env;
  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());
  ASSERT_OK_AND_ASSIGN(std::unique_ptr<JointState> joint_state,
                       JointState::Create(gazebo_client, topic));
  // Test publishing an empty joint state message and verify
  // GetJointPosition fails.
  // Wait for connection from `joint_state` and publish empty model msg.
  intrinsic::Thread publisher_thread([&gazebo_client, topic](StopToken st) {
    gz::transport::Node::Publisher pub =
        gazebo_client->Node().Advertise<gz::msgs::Model>(topic);
    gz::msgs::Model msg;
    while (!st.stop_requested() && !pub.HasConnections()) {
      absl::SleepFor(absl::Milliseconds(10));
    }
    if (pub.HasConnections()) {
      ASSERT_TRUE(pub.Publish(msg));
    }
  });
  EXPECT_THAT(joint_state->GetJointPosition(absl::Seconds(5), 0u),
              StatusIs(absl::StatusCode::kNotFound));
}

TEST(JointState, GetJointPositionWithInvalidAxisIndex) {
  const std::string topic = "test_state_topic_invalid_axis_index";
  GazeboClientTestEnv env;
  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());
  ASSERT_OK_AND_ASSIGN(std::unique_ptr<JointState> joint_state,
                       JointState::Create(gazebo_client, topic));
  // Test getting joint position with invalid axis index
  EXPECT_THAT(joint_state->GetJointPosition(absl::Seconds(1), 3u),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(JointState, GetJointState) {
  const std::string topic = "test_state_topic_joint_state";
  GazeboClientTestEnv env;
  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());
  ASSERT_OK_AND_ASSIGN(std::unique_ptr<JointState> joint_state,
                       JointState::Create(gazebo_client, topic));
  EXPECT_EQ(joint_state->Topic(), topic);

  // Test publishing a joint state msg and verify that JointState can retrieve
  // the joint position and velocity data via the GetJointState function
  // Wait for connection from `joint_state` and publish model msg.
  auto joint_state_publish_func = [&gazebo_client, topic](StopToken st) {
    gz::transport::Node::Publisher pub =
        gazebo_client->Node().Advertise<gz::msgs::Model>(topic);
    gz::msgs::Model msg;
    gz::msgs::Joint* joint_msg = msg.add_joint();
    joint_msg->mutable_axis1()->set_position(kPublishedJointPos[0]);
    joint_msg->mutable_axis1()->set_velocity(kPublishedJointVel[0]);
    joint_msg->mutable_axis2()->set_position(kPublishedJointPos[1]);
    joint_msg->mutable_axis2()->set_velocity(kPublishedJointVel[1]);
    while (!st.stop_requested() && !pub.HasConnections()) {
      absl::SleepFor(absl::Milliseconds(10));
    }
    if (pub.HasConnections()) {
      ASSERT_TRUE(pub.Publish(msg));
    }
  };
  intrinsic::Thread publisher_thread_axis_1(joint_state_publish_func);
  ASSERT_OK_AND_ASSIGN(JointState::State joint_state_received,
                       joint_state->GetJointState(absl::Seconds(5)));
  EXPECT_THAT(
      joint_state_received.axis_0,
      Optional(FieldsAre(kPublishedJointPos[0], kPublishedJointVel[0])));
  EXPECT_THAT(
      joint_state_received.axis_1,
      Optional(FieldsAre(kPublishedJointPos[1], kPublishedJointVel[1])));
}

TEST(JointState, SubscribeToJointStates) {
  const std::string topic = "test_state_topic_joint_state";
  GazeboClientTestEnv env;
  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());
  ASSERT_OK_AND_ASSIGN(std::unique_ptr<JointState> joint_state,
                       JointState::Create(gazebo_client, topic));
  EXPECT_EQ(joint_state->Topic(), topic);

  // Test publishing a joint state msg and verify that the Subscribe function
  // can receive joint states
  enum AxisMode { kEmpty = 0x00, kAxis1 = 0x01, kAxis2 = 0x02 };
  // Wait for connection from `joint_state` and publish model msg.
  gz::transport::Node::Publisher pub =
      gazebo_client->Node().Advertise<gz::msgs::Model>(topic);
  auto joint_state_publish_func = [topic, &pub](StopToken st,
                                                uint32_t axis_flags) {
    gz::msgs::Model msg;
    gz::msgs::Joint* joint_msg = msg.add_joint();
    if (axis_flags & kAxis1) {
      joint_msg->mutable_axis1()->set_position(kPublishedJointPos[0]);
      joint_msg->mutable_axis1()->set_velocity(kPublishedJointVel[0]);
    }
    if (axis_flags & kAxis2) {
      joint_msg->mutable_axis2()->set_position(kPublishedJointPos[1]);
      joint_msg->mutable_axis2()->set_velocity(kPublishedJointVel[1]);
    }
    while (!st.stop_requested()) {
      if (pub.HasConnections()) {
        ASSERT_TRUE(pub.Publish(msg));
      }
      absl::SleepFor(absl::Milliseconds(100));
    }
  };

  // Test subscribing to joint with single axis
  uint32_t axis_flags = kAxis1;
  intrinsic::Thread publisher_thread_single_axis(joint_state_publish_func,
                                                 axis_flags);
  absl::Notification received_msg_single_axis;
  uint32_t joint_state_single_axis_received_count = 0u;
  auto joint_state_single_axis_callback =
      [&joint_state_single_axis_received_count, &received_msg_single_axis](
          const absl::StatusOr<JointState::State> state) {
        ASSERT_OK(state);
        EXPECT_THAT(state->axis_0, Optional(FieldsAre(kPublishedJointPos[0],
                                                      kPublishedJointVel[0])));
        EXPECT_EQ(state->axis_1, std::nullopt);
        // Wait to make sure we can receive a few messages before notifying
        if (++joint_state_single_axis_received_count == kTargetPubMsgCount &&
            !received_msg_single_axis.HasBeenNotified()) {
          received_msg_single_axis.Notify();
        }
      };
  EXPECT_OK(joint_state->Subscribe(joint_state_single_axis_callback));

  received_msg_single_axis.WaitForNotificationWithTimeout(
      absl::Milliseconds(1000));

  EXPECT_EQ(joint_state_single_axis_received_count, kTargetPubMsgCount);

  // Unsubscribe and make sure publisher has no more connections
  EXPECT_OK(joint_state->Unsubscribe());
  EXPECT_FALSE(pub.HasConnections());

  // Stop publishing messages
  publisher_thread_single_axis.request_stop();
  publisher_thread_single_axis.join();

  // Test subscribing to joint with two axes
  axis_flags = kAxis1 | kAxis2;
  intrinsic::Thread publisher_thread_two_axes(joint_state_publish_func,
                                              axis_flags);
  absl::Notification received_msg_two_axes;
  uint32_t joint_state_two_axes_received_count = 0u;
  auto joint_state_two_axes_callback =
      [&joint_state_two_axes_received_count,
       &received_msg_two_axes](const absl::StatusOr<JointState::State> state) {
        ASSERT_OK(state);
        EXPECT_THAT(state->axis_0, Optional(FieldsAre(kPublishedJointPos[0],
                                                      kPublishedJointVel[0])));
        EXPECT_THAT(state->axis_1, Optional(FieldsAre(kPublishedJointPos[1],
                                                      kPublishedJointVel[1])));
        // Wait to make sure we can receive a few messages before notifying
        if (++joint_state_two_axes_received_count == kTargetPubMsgCount &&
            !received_msg_two_axes.HasBeenNotified()) {
          received_msg_two_axes.Notify();
        }
      };
  EXPECT_OK(joint_state->Subscribe(joint_state_two_axes_callback));

  received_msg_two_axes.WaitForNotificationWithTimeout(
      absl::Milliseconds(1000));
  EXPECT_EQ(joint_state_two_axes_received_count, kTargetPubMsgCount);

  // Unsubscribe and make sure publisher has no more connections
  EXPECT_OK(joint_state->Unsubscribe());
  EXPECT_FALSE(pub.HasConnections());
}

TEST(JointState, SubscribeToJointStatesEmptyJointStateMessage) {
  const std::string topic = "test_state_topic_joint_state";
  GazeboClientTestEnv env;
  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());
  ASSERT_OK_AND_ASSIGN(std::unique_ptr<JointState> joint_state,
                       JointState::Create(gazebo_client, topic));
  EXPECT_EQ(joint_state->Topic(), topic);

  // Test publishing an empty joint state message and verify
  // joint states are not found in the subscription callback
  gz::transport::Node::Publisher pub =
      gazebo_client->Node().Advertise<gz::msgs::Model>(topic);
  intrinsic::Thread publisher_thread([topic, &pub](StopToken st) {
    gz::msgs::Model msg;
    while (!st.stop_requested()) {
      if (pub.HasConnections()) {
        ASSERT_TRUE(pub.Publish(msg));
      }
      absl::SleepFor(absl::Milliseconds(100));
    }
  });

  absl::Notification received_msg;
  bool joint_state_received = false;
  auto joint_state_callback =
      [&joint_state_received,
       &received_msg](const absl::StatusOr<JointState::State> state) {
        if (!received_msg.HasBeenNotified()) {
          EXPECT_THAT(state, StatusIs(absl::StatusCode::kNotFound));
          joint_state_received = true;
          received_msg.Notify();
        }
      };
  EXPECT_OK(joint_state->Subscribe(joint_state_callback));

  EXPECT_TRUE(
      received_msg.WaitForNotificationWithTimeout(absl::Milliseconds(1000)));

  // Unsubscribe and make sure publisher has no more connections
  EXPECT_OK(joint_state->Unsubscribe());
  EXPECT_FALSE(pub.HasConnections());
}

void OnJointStateCallback(const absl::StatusOr<JointState::State> state,
                          uint32_t* joint_state_received_count,
                          absl::Notification* received_msg) {
  ASSERT_OK(state);
  EXPECT_THAT(state->axis_0, Optional(FieldsAre(kPublishedJointPos[0],
                                                kPublishedJointVel[0])));
  EXPECT_THAT(state->axis_1, Optional(FieldsAre(kPublishedJointPos[1],
                                                kPublishedJointVel[1])));
  if (++(*joint_state_received_count) == kTargetPubMsgCount &&
      !received_msg->HasBeenNotified()) {
    received_msg->Notify();
  }
}

TEST(JointState, SubscribeToJointStatesFreeFunction) {
  const std::string topic = "test_state_topic_joint_state";
  GazeboClientTestEnv env;
  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());
  ASSERT_OK_AND_ASSIGN(std::unique_ptr<JointState> joint_state,
                       JointState::Create(gazebo_client, topic));
  EXPECT_EQ(joint_state->Topic(), topic);

  gz::transport::Node::Publisher pub =
      gazebo_client->Node().Advertise<gz::msgs::Model>(topic);
  intrinsic::Thread publisher_thread([topic, &pub](StopToken st) {
    gz::msgs::Model msg;
    gz::msgs::Joint* joint_msg = msg.add_joint();
    joint_msg->mutable_axis1()->set_position(kPublishedJointPos[0]);
    joint_msg->mutable_axis1()->set_velocity(kPublishedJointVel[0]);
    joint_msg->mutable_axis2()->set_position(kPublishedJointPos[1]);
    joint_msg->mutable_axis2()->set_velocity(kPublishedJointVel[1]);
    while (!st.stop_requested()) {
      if (pub.HasConnections()) {
        ASSERT_TRUE(pub.Publish(msg));
      }
      absl::SleepFor(absl::Milliseconds(100));
    }
  });

  uint32_t joint_state_free_function_received_count = 0u;
  absl::Notification received_msg_free_function;
  EXPECT_OK(joint_state->Subscribe(std::bind(
      &OnJointStateCallback, std::placeholders::_1,
      &joint_state_free_function_received_count, &received_msg_free_function)));
  received_msg_free_function.WaitForNotificationWithTimeout(
      absl::Milliseconds(1000));
  EXPECT_EQ(joint_state_free_function_received_count, kTargetPubMsgCount);

  EXPECT_OK(joint_state->Unsubscribe());
  EXPECT_FALSE(pub.HasConnections());
}

TEST(JointState, SubscribeToJointStatesMemberFunction) {
  // Test class with member function for receiving Subscribe callbacks
  class MyTestClass {
   public:
    void OnJointStateMemberCallback(
        const absl::StatusOr<JointState::State> state) {
      ASSERT_OK(state);
      EXPECT_THAT(state->axis_0, Optional(FieldsAre(kPublishedJointPos[0],
                                                    kPublishedJointVel[0])));
      EXPECT_THAT(state->axis_1, Optional(FieldsAre(kPublishedJointPos[1],
                                                    kPublishedJointVel[1])));
      if (++joint_state_member_function_received_count_ == kTargetPubMsgCount &&
          !received_msg_member_function_.HasBeenNotified()) {
        received_msg_member_function_.Notify();
      }
    }
    uint32_t joint_state_member_function_received_count_ = 0u;
    absl::Notification received_msg_member_function_;
  };

  const std::string topic = "test_state_topic_joint_state";
  GazeboClientTestEnv env;
  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());
  ASSERT_OK_AND_ASSIGN(std::unique_ptr<JointState> joint_state,
                       JointState::Create(gazebo_client, topic));
  EXPECT_EQ(joint_state->Topic(), topic);

  gz::transport::Node::Publisher pub =
      gazebo_client->Node().Advertise<gz::msgs::Model>(topic);
  intrinsic::Thread publisher_thread([topic, &pub](StopToken st) {
    gz::msgs::Model msg;
    gz::msgs::Joint* joint_msg = msg.add_joint();
    joint_msg->mutable_axis1()->set_position(kPublishedJointPos[0]);
    joint_msg->mutable_axis1()->set_velocity(kPublishedJointVel[0]);
    joint_msg->mutable_axis2()->set_position(kPublishedJointPos[1]);
    joint_msg->mutable_axis2()->set_velocity(kPublishedJointVel[1]);
    while (!st.stop_requested()) {
      if (pub.HasConnections()) {
        ASSERT_TRUE(pub.Publish(msg));
      }
      absl::SleepFor(absl::Milliseconds(100));
    }
  });

  MyTestClass test_class;
  EXPECT_OK(
      joint_state->Subscribe(std::bind(&MyTestClass::OnJointStateMemberCallback,
                                       &test_class, std::placeholders::_1)));
  test_class.received_msg_member_function_.WaitForNotificationWithTimeout(
      absl::Milliseconds(1000));
  EXPECT_EQ(test_class.joint_state_member_function_received_count_,
            kTargetPubMsgCount);

  EXPECT_OK(joint_state->Unsubscribe());
  EXPECT_FALSE(pub.HasConnections());
}

TEST(JointState, MultipleSubscriptionOnSameTopic) {
  const std::string topic = "test_state_topic_joint_state";
  GazeboClientTestEnv env;
  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());
  ASSERT_OK_AND_ASSIGN(std::unique_ptr<JointState> joint_state,
                       JointState::Create(gazebo_client, topic));
  ASSERT_OK_AND_ASSIGN(std::unique_ptr<JointState> joint_state_duplicate,
                       JointState::Create(gazebo_client, topic));

  EXPECT_EQ(joint_state->Topic(), topic);
  EXPECT_EQ(joint_state_duplicate->Topic(), topic);

  auto joint_state_callback =
      [](const absl::StatusOr<JointState::State> state) {};

  // First subscription should work without any issues.
  EXPECT_OK(joint_state->Subscribe(joint_state_callback));
  // The second subscription from the same joint state object fails because
  // it already has an existing subscription.
  EXPECT_THAT(joint_state->Subscribe(joint_state_callback),
              StatusIs(absl::StatusCode::kFailedPrecondition));
  // Subscription on the same topic from another joint state object should also
  // work fine
  EXPECT_OK(joint_state_duplicate->Subscribe(joint_state_callback));
  EXPECT_OK(joint_state->Unsubscribe());
  EXPECT_OK(joint_state_duplicate->Unsubscribe());

  // Test subscription again after unsubscribing
  EXPECT_OK(joint_state->Subscribe(joint_state_callback));
  EXPECT_OK(joint_state->Unsubscribe());
}

}  // namespace
}  // namespace simulation
}  // namespace intrinsic
