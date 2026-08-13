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

#include "intrinsic/util/status/status_conversion_grpc.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <optional>

#include "absl/status/status.h"
#include "absl/strings/cord.h"
#include "google/protobuf/wrappers.pb.h"
#include "google/rpc/code.pb.h"
#include "grpcpp/support/status.h"
#include "intrinsic/util/proto/type_url.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic {
namespace {

using ::testing::TestWithParam;
using ::testing::Values;

class AbslStatusConversionRoundTripTest
    : public TestWithParam<absl::StatusCode> {};

INSTANTIATE_TEST_SUITE_P(
    AbslCodes, AbslStatusConversionRoundTripTest,
    Values(absl::StatusCode::kOk, absl::StatusCode::kCancelled,
           absl::StatusCode::kUnknown, absl::StatusCode::kInvalidArgument,
           absl::StatusCode::kDeadlineExceeded, absl::StatusCode::kNotFound,
           absl::StatusCode::kAlreadyExists,
           absl::StatusCode::kPermissionDenied,
           absl::StatusCode::kResourceExhausted,
           absl::StatusCode::kFailedPrecondition, absl::StatusCode::kAborted,
           absl::StatusCode::kOutOfRange, absl::StatusCode::kUnimplemented,
           absl::StatusCode::kInternal, absl::StatusCode::kUnavailable,
           absl::StatusCode::kDataLoss, absl::StatusCode::kUnauthenticated));

TEST_P(AbslStatusConversionRoundTripTest, CodeRoundTripsThroughGrpcStatus) {
  absl::Status absl_status = absl::Status(GetParam(), "");
  absl::Status returned_status = ToAbslStatus(ToGrpcStatus(absl_status));
  EXPECT_EQ(returned_status.code(), absl_status.code());
}

class GrpcStatusConversionRoundTripTest
    : public TestWithParam<grpc::StatusCode> {};

INSTANTIATE_TEST_SUITE_P(
    GrpcCodes, GrpcStatusConversionRoundTripTest,
    Values(grpc::StatusCode::OK, grpc::StatusCode::CANCELLED,
           grpc::StatusCode::UNKNOWN, grpc::StatusCode::INVALID_ARGUMENT,
           grpc::StatusCode::DEADLINE_EXCEEDED, grpc::StatusCode::NOT_FOUND,
           grpc::StatusCode::ALREADY_EXISTS,
           grpc::StatusCode::PERMISSION_DENIED,
           grpc::StatusCode::RESOURCE_EXHAUSTED,
           grpc::StatusCode::FAILED_PRECONDITION, grpc::StatusCode::ABORTED,
           grpc::StatusCode::OUT_OF_RANGE, grpc::StatusCode::UNIMPLEMENTED,
           grpc::StatusCode::INTERNAL, grpc::StatusCode::UNAVAILABLE,
           grpc::StatusCode::DATA_LOSS, grpc::StatusCode::UNAUTHENTICATED));

TEST_P(GrpcStatusConversionRoundTripTest, CodeRoundTripsThroughAbslStatus) {
  grpc::Status grpc_status = grpc::Status(GetParam(), "");
  grpc::Status returned_status = ToGrpcStatus(ToAbslStatus(grpc_status));
  EXPECT_EQ(returned_status.error_code(), grpc_status.error_code());
}

TEST(StatusConversionGrpcTest, AbslStatusToGrpcStatusRoundTrip) {
  absl::Status absl_status = absl::InternalError("A terrible thing happened!");
  google::protobuf::StringValue value;
  value.set_value("Foo");
  absl_status.SetPayload(AddTypeUrlPrefix(value.GetDescriptor()->full_name()),
                         value.SerializeAsCord());

  // Round trip through grpc status.
  grpc::Status grpc_status = ToGrpcStatus(absl_status);
  absl::Status returned_status = ToAbslStatus(grpc_status);

  EXPECT_EQ(returned_status.code(), absl_status.code());
  EXPECT_EQ(returned_status.message(), absl_status.message());

  google::protobuf::StringValue read_value;

  std::optional<absl::Cord> read_payload = returned_status.GetPayload(
      AddTypeUrlPrefix(value.GetDescriptor()->full_name()));
  ASSERT_TRUE(read_payload.has_value());
  ASSERT_TRUE(read_value.ParseFromString(*read_payload));
  EXPECT_THAT(read_value.value(), ::testing::Eq(value.value()));
}

}  // namespace
}  // namespace intrinsic
