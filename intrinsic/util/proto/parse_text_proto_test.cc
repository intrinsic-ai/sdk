// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/proto/parse_text_proto.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "absl/status/status.h"
#include "google/protobuf/type.pb.h"
#include "google/protobuf/wrappers.pb.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic {
namespace {

using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;
using ::intrinsic::testing::EqualsProto;
using ::testing::AllOf;
using ::testing::HasSubstr;

TEST(ParseTextProtoTest, ParseTextProto) {
  EXPECT_THAT(ParseTextProto<google::protobuf::Int32Value>("value: 1"),
              IsOkAndHolds(EqualsProto("value: 1")));
}

TEST(ParseTextProtoTest, ParseTextProtoWorksWithCustomAnyTypeUrlPrefix) {
  google::protobuf::Int32Value int32_value;
  int32_value.set_value(1);
  google::protobuf::Option option;
  option.mutable_value()->PackFrom(int32_value);

  EXPECT_THAT(
      ParseTextProto<google::protobuf::Option>(R"pb(
        value {
          [custom.type.prefix.com/google.protobuf.Int32Value] { value: 1 }
        })pb"),
      IsOkAndHolds(EqualsProto(option)));
}

TEST(ParseTextProtoTest, ParseTextProtoInvalidInput) {
  EXPECT_THAT(
      ParseTextProto<google::protobuf::Int32Value>("non_existent_field: 1"),
      StatusIs(
          absl::StatusCode::kInvalidArgument,
          AllOf(HasSubstr("Cannot parse protobuf google.protobuf.Int32Value"),
                HasSubstr("non_existent_field"))));
}

TEST(ParseTextProtoTest, ParseTextProtoOrDie) {
  google::protobuf::Int32Value result = ParseTextProtoOrDie("value: 1");

  EXPECT_THAT(result, EqualsProto("value: 1"));
}

TEST(ParseTextProtoTest, ParseTextOrDie) {
  EXPECT_THAT(ParseTextOrDie<google::protobuf::Int32Value>("value: 1"),
              EqualsProto("value: 1"));
}

}  // namespace
}  // namespace intrinsic
