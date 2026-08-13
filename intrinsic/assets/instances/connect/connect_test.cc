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

#include "intrinsic/assets/instances/connect/connect.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <map>
#include <memory>
#include <string>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "grpc/grpc_security_constants.h"
#include "grpcpp/channel.h"
#include "grpcpp/security/server_credentials.h"
#include "grpcpp/server.h"
#include "grpcpp/server_builder.h"
#include "grpcpp/server_context.h"
#include "intrinsic/assets/instances/connect/testing/test_service.grpc.pb.h"
#include "intrinsic/assets/proto/v1/grpc_connection.pb.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic::assets::instances::connect {
namespace {

using ::intrinsic_proto::assets::testing::TestRequest;
using ::intrinsic_proto::assets::testing::TestResponse;
using ::intrinsic_proto::assets::testing::TestService;
using ::intrinsic_proto::assets::v1::GrpcConnection;

class TestServiceImpl final : public TestService::Service {
 public:
  grpc::Status Test(grpc::ServerContext* context, const TestRequest* request,
                    TestResponse* response) override {
    for (const auto& [key, value] : context->client_metadata()) {
      (*response
            ->mutable_context_metadata())[std::string(key.data(), key.length())]
          .add_values(std::string(value.data(), value.length()));
    }
    return grpc::Status::OK;
  }
};

class ConnectTest : public ::testing::Test {
 protected:
  void SetUp() override {
    grpc::ServerBuilder builder;
    builder.RegisterService(&service_);
    int selected_port = 0;
    builder.AddListeningPort(
        "[::]:0", grpc::experimental::LocalServerCredentials(LOCAL_TCP),
        &selected_port);
    server_ = builder.BuildAndStart();
    ASSERT_NE(server_, nullptr);
    ASSERT_GT(selected_port, 0);
    server_address_ = absl::StrFormat("localhost:%d", selected_port);
  }

  void TearDown() override { server_->Shutdown(); }

  std::unique_ptr<grpc::Server> server_;
  TestServiceImpl service_;
  std::string server_address_;
};

TEST_F(ConnectTest, SuccessWithMetadata) {
  GrpcConnection connection;
  connection.set_address(server_address_);
  auto* meta1 = connection.add_metadata();
  meta1->set_key("test_key");
  meta1->set_value("test_value1");
  auto* meta2 = connection.add_metadata();
  meta2->set_key("test_key");
  meta2->set_value("test_value2");

  absl::StatusOr<std::shared_ptr<grpc::Channel>> channel_or =
      Connect(connection);
  ASSERT_OK(channel_or);

  auto stub = TestService::NewStub(*channel_or);
  TestRequest request;
  TestResponse response;
  grpc::ClientContext context;
  ASSERT_TRUE(stub->Test(&context, request, &response).ok());

  const auto& metadata = response.context_metadata();
  ASSERT_EQ(metadata.count("test_key"), 1);
  const auto& values = metadata.at("test_key").values();
  EXPECT_THAT(std::vector<std::string>(values.begin(), values.end()),
              ::testing::UnorderedElementsAre("test_value1", "test_value2"));
}

TEST_F(ConnectTest, SuccessWithoutMetadata) {
  GrpcConnection connection;
  connection.set_address(server_address_);

  absl::StatusOr<std::shared_ptr<grpc::Channel>> channel_or =
      Connect(connection);
  ASSERT_OK(channel_or);

  auto stub = TestService::NewStub(*channel_or);
  TestRequest request;
  TestResponse response;
  grpc::ClientContext context;
  ASSERT_TRUE(stub->Test(&context, request, &response).ok());
}

}  // namespace
}  // namespace intrinsic::assets::instances::connect
