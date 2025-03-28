// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/proto/any.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "absl/status/status.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/wrappers.pb.h"
#include "intrinsic/util/proto/testing/param_message.pb.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic {
namespace {

using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;
using ::intrinsic::testing::EqualsProto;
using ::testing::AllOf;
using ::testing::HasSubstr;

TEST(UnpackAny, UnpackAnyWrongTypeFail) {
  google::protobuf::FloatValue float_value;
  google::protobuf::Any any;
  any.PackFrom(float_value);
  EXPECT_THAT(UnpackAny<google::protobuf::DoubleValue>(any),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       AllOf(HasSubstr("google.protobuf.FloatValue"),
                             HasSubstr("google.protobuf.DoubleValue"))));
}

TEST(UnpackAny, UnpackAnyWorks) {
  google::protobuf::FloatValue float_value;
  float_value.set_value(18.0);
  google::protobuf::Any any;
  any.PackFrom(float_value);
  EXPECT_THAT(UnpackAny<google::protobuf::FloatValue>(any),
              IsOkAndHolds(EqualsProto(float_value)));
}

TEST(UnpackAny, UnpackAnyToParamWorks) {
  google::protobuf::FloatValue float_value;
  float_value.set_value(18.0);
  google::protobuf::Any any;
  any.PackFrom(float_value);
  google::protobuf::FloatValue recovered;
  ASSERT_OK(UnpackAny(any, recovered));
  EXPECT_THAT(recovered, EqualsProto(float_value));
}

TEST(UnpackAnyAndMerge, AppliesDefaults) {
  intrinsic_proto::test::ParamMessageDefaultsTestMessage msg;
  google::protobuf::Any msg_any;
  msg.set_my_string("foo");
  msg_any.PackFrom(msg);

  intrinsic_proto::test::ParamMessageDefaultsTestMessage defaults;
  google::protobuf::Any defaults_any;
  defaults.set_my_string("bar");
  defaults.set_maybe_int32(7);
  defaults_any.PackFrom(defaults);

  intrinsic_proto::test::ParamMessageDefaultsTestMessage expected;
  expected.set_my_string("foo");
  expected.set_maybe_int32(7);

  EXPECT_THAT(
      UnpackAnyAndMerge<intrinsic_proto::test::ParamMessageDefaultsTestMessage>(
          msg_any, defaults_any),
      IsOkAndHolds(EqualsProto(expected)));
}

}  // namespace
}  // namespace intrinsic
