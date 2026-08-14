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

#include "intrinsic/simulation/gazebo/client/gazebo_client.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "absl/strings/str_replace.h"
#include "absl/time/time.h"
#include "grpcpp/server.h"
#include "grpcpp/server_builder.h"
#include "grpcpp/server_context.h"
#include "gz/transport/Node.hh"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"
#include "intrinsic/simulation/gazebo/client/gazebo_client_test_env.h"
#include "intrinsic/util/proto/parse_text_proto.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic {
namespace simulation {
namespace {

using ::intrinsic::testing::EqualsProto;
using ::testing::_;
using ::testing::DoAll;
using ::testing::Pointee;
using ::testing::Return;
using ::testing::SetArgPointee;

TEST(GazeboClient, CreateClient) {
  GazeboClientTestEnv env;

  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());
  ASSERT_FALSE(gazebo_client->Node().GlobalRelays().empty());
  EXPECT_EQ(gazebo_client->Node().GlobalRelays()[0], "127.0.0.1");
  EXPECT_FALSE(gazebo_client->Node().Options().Partition().empty());
}

TEST(GazeboClient, PartitionSetFromTransportRelay) {
  GazeboClientTestEnv env;

  const std::string kCustomPartition = "custom_sim_partition";
  intrinsic_proto::simulation::v1::GetIPResponse mock_response;
  mock_response.set_ip_address("127.0.0.1");
  mock_response.set_partition(kCustomPartition);
  EXPECT_CALL(env.service(), GetIPForTransportRelay(_, _, _))
      .WillRepeatedly(
          DoAll(SetArgPointee<2>(mock_response), Return(::grpc::Status::OK)));

  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());
  EXPECT_EQ(gazebo_client->Node().Options().Partition(), kCustomPartition);
}

TEST(GazeboClient, CreateClientFromResolvedDependency) {
  GazeboClientTestEnv env;

  const intrinsic_proto::assets::v1::ResolvedDependency resolved_dep =
      ParseTextProtoOrDie(absl::StrReplaceAll(
          R"pb(interfaces: {
                 key: "grpc://intrinsic_proto.simulation.v1.GazeboService"
                 value: { grpc: { connection: { address: "%s" } } }
               })pb",
          {{"%s", env.address()}}));

  ASSERT_OK_AND_ASSIGN(auto gazebo_client,
                       GazeboClient::CreateFromResolvedDependency(
                           resolved_dep, absl::InfiniteDuration()));
  ASSERT_FALSE(gazebo_client->Node().GlobalRelays().empty());
  EXPECT_EQ(gazebo_client->Node().GlobalRelays()[0], "127.0.0.1");
  EXPECT_FALSE(gazebo_client->Node().Options().Partition().empty());
}

TEST(GazeboClient, ListTopics) {
  GazeboClientTestEnv env;

  ASSERT_OK_AND_ASSIGN(auto gazebo_client, env.CreateClient());

  intrinsic_proto::simulation::v1::Topics list_resp;
  list_resp.add_topics()->set_topic_name("/model/robot/joint/0/cmd_pos");

  EXPECT_CALL(env.service(), ListTopics(_, Pointee(EqualsProto(R"pb(
                                          object_name: "robot"
                                          entity_type_filter: MODELS
                                        )pb")),
                                        _))
      .WillOnce(
          [&list_resp](
              ::grpc::ServerContext* context,
              const intrinsic_proto::simulation::v1::ListTopicsRequest* request,
              intrinsic_proto::simulation::v1::Topics* response) {
            *response = list_resp;
            return ::grpc::Status::OK;
          });

  ASSERT_OK_AND_ASSIGN(auto resp,
                       gazebo_client->ListTopics(
                           "robot", intrinsic_proto::simulation::v1::MODELS));
  EXPECT_EQ(resp.topics(0).topic_name(), "/model/robot/joint/0/cmd_pos");
}

}  // namespace
}  // namespace simulation
}  // namespace intrinsic
