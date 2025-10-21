// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/assets/dependencies/utils.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <map>
#include <memory>
#include <string>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "absl/strings/str_replace.h"
#include "grpc/grpc_security_constants.h"
#include "grpcpp/channel.h"
#include "grpcpp/security/server_credentials.h"
#include "grpcpp/server.h"
#include "grpcpp/server_builder.h"
#include "grpcpp/server_context.h"
#include "intrinsic/assets/dependencies/testing/test_service.grpc.pb.h"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"
#include "intrinsic/util/proto/parse_text_proto.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic::assets::dependencies {
namespace {

using ::absl_testing::StatusIs;
using ::intrinsic::ParseTextProtoOrDie;
using ::intrinsic_proto::assets::dependencies::testing::TestRequest;
using ::intrinsic_proto::assets::dependencies::testing::TestResponse;
using ::intrinsic_proto::assets::dependencies::testing::TestService;
using ::intrinsic_proto::assets::v1::ResolvedDependency;

// A test gRPC service that returns the metadata from the incoming context.
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

// ... existing code ...
class UtilsTest : public ::testing::Test {
 protected:
  void SetUp() override {
    grpc::ServerBuilder builder;
    builder.RegisterService(&service_);
    int selected_port = 0;
    // Use "[::]:0" to bind to any available port.
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

struct ConnectTestParam {
  std::string test_name;
  std::string dep_textproto;
  std::string iface;
  std::map<std::string, std::vector<std::string>> expected_metadata;
  absl::StatusCode expected_code;
};

class ParameterizedUtilsTest
    : public UtilsTest,
      public ::testing::WithParamInterface<ConnectTestParam> {};

TEST_P(ParameterizedUtilsTest, Connect) {
  const ConnectTestParam& param = GetParam();
  const std::string& server_address = server_address_;
  const ResolvedDependency dep = ParseTextProtoOrDie(
      absl::StrReplaceAll(param.dep_textproto, {{"%s", server_address}}));

  grpc::ClientContext context;
  absl::StatusOr<std::shared_ptr<grpc::Channel>> channel_or =
      Connect(context, dep, param.iface);

  if (param.expected_code != absl::StatusCode::kOk) {
    EXPECT_THAT(channel_or.status(), StatusIs(param.expected_code));
  } else {
    ASSERT_OK(channel_or);
    auto stub = TestService::NewStub(*channel_or);
    TestRequest request;
    TestResponse response;

    ASSERT_TRUE(stub->Test(&context, request, &response).ok());

    const auto& metadata = response.context_metadata();
    for (const auto& [key, expected_values] : param.expected_metadata) {
      ASSERT_EQ(metadata.count(key), 1);
      const auto& values = metadata.at(key).values();
      EXPECT_THAT(std::vector<std::string>(values.begin(), values.end()),
                  ::testing::UnorderedElementsAreArray(expected_values));
    }
  }
}

INSTANTIATE_TEST_SUITE_P(
    ConnectTests, ParameterizedUtilsTest,
    ::testing::Values(
        ConnectTestParam{
            "Success",
            R"pb(interfaces: {
                   key: "grpc://intrinsic_proto.assets.dependencies.testing.TestService"
                   value: {
                     grpc_connection: {
                       address: "%s"
                       metadata: { key: "test_key", value: "test_value1" }
                       metadata: { key: "test_key", value: "test_value2" }
                     }
                   }
                 })pb",
            "grpc://intrinsic_proto.assets.dependencies.testing.TestService",
            {{"test_key", {"test_value1", "test_value2"}}},
            absl::StatusCode::kOk},
        ConnectTestParam{"MissingInterface",
                         "",
                         "grpc://intrinsic_proto.assets.dependencies.testing."
                         "TestService",
                         {},
                         absl::StatusCode::kNotFound},
        ConnectTestParam{
            "NotGrpc",
            R"pb(interfaces: {
                   key: "data://foo"
                   value: { data: { id: { package: "foo", name: "bar" } } }
                 })pb",
            "data://foo",
            {},
            absl::StatusCode::kInvalidArgument}),
    [](const ::testing::TestParamInfo<ConnectTestParam>& info) {
      return info.param.test_name;
    });

}  // namespace
}  // namespace intrinsic::assets::dependencies
