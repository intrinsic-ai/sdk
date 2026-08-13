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

#include "intrinsic/util/testing/status_payload_matchers.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/cord.h"
#include "google/protobuf/wrappers.pb.h"
#include "intrinsic/util/proto/type_url.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

using ::intrinsic::testing::EqualsProto;
using ::intrinsic::testing::WhenDeserializedAs;

namespace intrinsic::testing {
namespace {

TEST(StatusPayloadProtoMatcher, StatusMatches) {
  google::protobuf::Int32Value int_value;
  int_value.set_value(123);

  absl::Status s = absl::InvalidArgumentError("Foo");
  s.SetPayload(AddTypeUrlPrefix(int_value), int_value.SerializeAsCord());

  EXPECT_THAT(s, StatusHasProtoPayload<google::protobuf::Int32Value>(
                     EqualsProto(int_value)));
}

TEST(StatusPayloadProtoMatcher, MissingPayloadDoesNotMatch) {
  google::protobuf::Int32Value int_value;
  int_value.set_value(123);

  absl::Status s = absl::InvalidArgumentError("Foo");
  s.SetPayload(AddTypeUrlPrefix(int_value), int_value.SerializeAsCord());

  EXPECT_THAT(s, Not(StatusHasProtoPayload<google::protobuf::DoubleValue>(
                     EqualsProto(int_value))));
}

TEST(StatusPayloadProtoMatcher, OkStatusDoesNotMatch) {
  google::protobuf::Int32Value int_value;
  int_value.set_value(123);

  absl::Status s = absl::OkStatus();
  s.SetPayload(AddTypeUrlPrefix(int_value), int_value.SerializeAsCord());

  EXPECT_THAT(s, Not(StatusHasProtoPayload<google::protobuf::Int32Value>(
                     EqualsProto(int_value))));
}

TEST(StatusPayloadProtoMatcher, StatusOrMatches) {
  google::protobuf::Int32Value int_value;
  int_value.set_value(123);

  absl::Status s = absl::InvalidArgumentError("Foo");
  s.SetPayload(AddTypeUrlPrefix(int_value), int_value.SerializeAsCord());
  absl::StatusOr<std::string> s_or = s;

  EXPECT_THAT(s_or, StatusHasProtoPayload<google::protobuf::Int32Value>(
                        EqualsProto(int_value)));
}

TEST(StatusPayloadProtoMatcher, MissingPayloadDoesNotMatchStatusOr) {
  google::protobuf::Int32Value int_value;
  int_value.set_value(123);

  absl::Status s = absl::InvalidArgumentError("Foo");
  s.SetPayload(AddTypeUrlPrefix(int_value), int_value.SerializeAsCord());
  absl::StatusOr<std::string> s_or = s;

  EXPECT_THAT(s_or, Not(StatusHasProtoPayload<google::protobuf::DoubleValue>(
                        EqualsProto(int_value))));
}

TEST(StatusPayloadProtoMatcher, OkStatusOrDoesNotMatch) {
  absl::StatusOr<std::string> s = "Foo";

  google::protobuf::Int32Value int_value;
  int_value.set_value(123);

  EXPECT_THAT(s, Not(StatusHasProtoPayload<google::protobuf::Int32Value>(
                     EqualsProto(int_value))));
}

TEST(StatusPayloadGenericMatcher, StatusMatchesOnString) {
  absl::Status s = absl::InvalidArgumentError("Foo");
  s.SetPayload("my_value", absl::Cord("Bar"));

  EXPECT_THAT(s, StatusHasGenericPayload("my_value", "Bar"));
}

TEST(StatusPayloadGenericMatcher, StatusMatchesOnProto) {
  google::protobuf::Int32Value int_value;
  int_value.set_value(123);

  absl::Status s = absl::InvalidArgumentError("Foo");
  s.SetPayload(AddTypeUrlPrefix(int_value), int_value.SerializeAsCord());

  EXPECT_THAT(s, StatusHasGenericPayload(
                     AddTypeUrlPrefix<google::protobuf::Int32Value>(),
                     WhenDeserializedAs<google::protobuf::Int32Value>(
                         EqualsProto(int_value))));
}

TEST(StatusPayloadGenericMatcher, StatusMatchesHasPayload) {
  google::protobuf::Int32Value int_value;
  int_value.set_value(123);

  absl::Status s = absl::InvalidArgumentError("Foo");
  s.SetPayload(AddTypeUrlPrefix(int_value), int_value.SerializeAsCord());

  EXPECT_THAT(s, StatusHasGenericPayload(
                     AddTypeUrlPrefix<google::protobuf::Int32Value>()));
}

}  // namespace
}  // namespace intrinsic::testing
