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

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <memory>

#include "absl/status/status.h"
#include "intrinsic/simulation/gazebo/client/gazebo_client.h"
#include "intrinsic/simulation/gazebo/client/gazebo_client_test_env.h"
#include "intrinsic/simulation/gazebo/client/joint_control.h"
#include "intrinsic/simulation/gazebo/client/joint_state.h"
#include "intrinsic/simulation/gazebo/gz_topic_constants.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

using ::absl_testing::StatusIs;

namespace intrinsic {
namespace simulation {
namespace {

using ::testing::_;
using ::testing::DoAll;
using ::testing::Return;
using ::testing::SetArgPointee;

TEST(TopicUtils, CanCreateJointControl) {
  intrinsic_proto::simulation::v1::Topics mock_response;
  auto* topic_info = mock_response.add_topics();
  topic_info->set_topic_name("custom_control_topic");
  topic_info->set_is_sim_server_advertised(false);
  topic_info->set_metadata(
      std::string(kGzTopicInfoJointPositionControlMetadata));
  topic_info->mutable_joint_info()->set_name("joint_control_test");
  topic_info->mutable_joint_info()->set_axis_index(0);

  GazeboClientTestEnv env;
  EXPECT_CALL(env.service(), ListTopics(_, _, _))
      .WillOnce(
          DoAll(SetArgPointee<2>(mock_response), Return(::grpc::Status::OK)));

  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());

  ASSERT_OK_AND_ASSIGN(std::unique_ptr<JointControl> joint_control,
                       CreateJointControl("some_name", gazebo_client,
                                          /*joint_name=*/"joint_control_test",
                                          /*axis_index=*/0));
  ASSERT_NE(joint_control, nullptr);
  EXPECT_EQ(joint_control->Topic(), "custom_control_topic");
}

TEST(TopicUtils, CanCreateJointState) {
  intrinsic_proto::simulation::v1::Topics mock_response;
  auto* topic_info = mock_response.add_topics();
  topic_info->set_topic_name("custom_state_topic");
  topic_info->set_is_sim_server_advertised(true);
  topic_info->set_metadata(std::string(kGzTopicInfoJointStateMetadata));
  topic_info->mutable_joint_info()->set_name("joint_state_test");
  topic_info->mutable_joint_info()->set_axis_index(0);

  GazeboClientTestEnv env;
  EXPECT_CALL(env.service(), ListTopics(_, _, _))
      .WillOnce(
          DoAll(SetArgPointee<2>(mock_response), Return(::grpc::Status::OK)));

  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());

  ASSERT_OK_AND_ASSIGN(std::unique_ptr<JointState> joint_state,
                       CreateJointState("some_name", gazebo_client,
                                        /*joint_name=*/"joint_state_test",
                                        /*axis_index=*/0));
  ASSERT_NE(joint_state, nullptr);
  EXPECT_EQ(joint_state->Topic(), "custom_state_topic");
}

TEST(TopicUtils, CannotCreateJointControlIfTopicNotFound) {
  GazeboClientTestEnv env;
  // Return empty topic list.
  EXPECT_CALL(env.service(), ListTopics(_, _, _))
      .WillOnce(Return(::grpc::Status::OK));

  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());

  EXPECT_THAT(CreateJointControl("some_name", gazebo_client,
                                 /*joint_name=*/"joint_control_test",
                                 /*axis_index=*/0),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(TopicUtils, CannotCreateJointStateIfTopicNotFound) {
  GazeboClientTestEnv env;
  // Return empty topic list.
  EXPECT_CALL(env.service(), ListTopics(_, _, _))
      .WillOnce(Return(::grpc::Status::OK));

  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());

  EXPECT_THAT(CreateJointState("some_name", gazebo_client,
                               /*joint_name=*/"joint_state_test",
                               /*axis_index=*/0),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

}  // namespace
}  // namespace simulation
}  // namespace intrinsic
