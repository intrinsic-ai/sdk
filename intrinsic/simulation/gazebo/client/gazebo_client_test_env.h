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

#ifndef INTRINSIC_SIMULATION_GAZEBO_CLIENT_GAZEBO_CLIENT_TEST_ENV_H_
#define INTRINSIC_SIMULATION_GAZEBO_CLIENT_GAZEBO_CLIENT_TEST_ENV_H_

#include <memory>
#include <string>
#include <string_view>
#include <utility>

#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "gmock/gmock.h"
#include "grpcpp/security/server_credentials.h"
#include "grpcpp/server.h"
#include "grpcpp/server_builder.h"
#include "grpcpp/server_context.h"
#include "gtest/gtest.h"
#include "intrinsic/simulation/gazebo/client/gazebo_client.h"
#include "intrinsic/simulation/gazebo/proto/v1/gazebo_service.grpc.pb.h"
#include "intrinsic/simulation/gazebo/proto/v1/gazebo_service.pb.h"

namespace intrinsic::simulation {

class MockGazeboService
    : public intrinsic_proto::simulation::v1::GazeboService::Service {
 public:
  MOCK_METHOD(::grpc::Status, GetIPForTransportRelay,
              (::grpc::ServerContext*,
               const intrinsic_proto::simulation::v1::GetIPRequest*,
               intrinsic_proto::simulation::v1::GetIPResponse*),
              (override));
  MOCK_METHOD(::grpc::Status, ListTopics,
              (::grpc::ServerContext*,
               const intrinsic_proto::simulation::v1::ListTopicsRequest*,
               intrinsic_proto::simulation::v1::Topics*),
              (override));
};

// Test environment helper for Gazebo client unit tests.
// Manages an in-process gRPC server that serves a `MockGazeboService`,
// and provides convenience methods to access the mock service and create
// connected `GazeboClient` clients.
class GazeboClientTestEnv {
 public:
  explicit GazeboClientTestEnv(std::string_view ip_address = "127.0.0.1") {
    intrinsic_proto::simulation::v1::GetIPResponse mock_response;
    mock_response.set_ip_address(std::string(ip_address));
    EXPECT_CALL(service_, GetIPForTransportRelay(::testing::_, ::testing::_,
                                                 ::testing::_))
        .WillRepeatedly(
            ::testing::DoAll(::testing::SetArgPointee<2>(mock_response),
                             ::testing::Return(::grpc::Status::OK)));

    int selected_port = 0;
    grpc::ServerBuilder builder;
    builder.RegisterService(&service_);
    builder.SetMaxReceiveMessageSize(-1);
    builder.SetMaxSendMessageSize(-1);
    builder.AddListeningPort("dns:///localhost:0",
                             grpc::InsecureServerCredentials(), &selected_port);
    server_ = builder.BuildAndStart();
    address_ = absl::StrCat("dns:///localhost:", selected_port);
  }

  ~GazeboClientTestEnv() {
    if (server_) {
      server_->Shutdown();
    }
  }

  MockGazeboService& service() { return service_; }
  std::string address() const { return address_; }

  absl::StatusOr<std::shared_ptr<GazeboClient>> CreateClient() const {
    return GazeboClient::Create(address_);
  }

 private:
  MockGazeboService service_;
  std::string address_;
  std::unique_ptr<grpc::Server> server_;
};

}  // namespace intrinsic::simulation

#endif  // INTRINSIC_SIMULATION_GAZEBO_CLIENT_GAZEBO_CLIENT_TEST_ENV_H_
