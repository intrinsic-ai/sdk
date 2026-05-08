// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/proto/any.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <string>

#include "absl/status/status.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/wrappers.pb.h"
#include "intrinsic/util/proto/parse_text_proto.h"
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

struct UnpackAnyAndMergeTestCase {
  std::string name;
  std::string defaults;
  std::string params;
  std::string expected;
};

using UnpackAnyAndMergeTest =
    ::testing::TestWithParam<UnpackAnyAndMergeTestCase>;

TEST_P(UnpackAnyAndMergeTest, MergeBehavior) {
  const auto& test_case = GetParam();
  auto defaults_msg =
      ParseTextOrDie<intrinsic_proto::test::ParamMessageDefaultsTestMessage>(
          test_case.defaults);
  auto params_msg =
      ParseTextOrDie<intrinsic_proto::test::ParamMessageDefaultsTestMessage>(
          test_case.params);
  auto expected_msg =
      ParseTextOrDie<intrinsic_proto::test::ParamMessageDefaultsTestMessage>(
          test_case.expected);

  google::protobuf::Any defaults_any;
  defaults_any.PackFrom(defaults_msg);
  google::protobuf::Any params_any;
  params_any.PackFrom(params_msg);

  EXPECT_THAT(
      UnpackAnyAndMerge<intrinsic_proto::test::ParamMessageDefaultsTestMessage>(
          params_any, defaults_any),
      IsOkAndHolds(EqualsProto(expected_msg)));
}

INSTANTIATE_TEST_SUITE_P(
    UnpackAnyAndMergeTests, UnpackAnyAndMergeTest,
    ::testing::ValuesIn<UnpackAnyAndMergeTestCase>({
        {.name = "AppliesDefaults",
         .defaults = "my_string: 'bar' maybe_int32: 7",
         .params = "my_string: 'foo'",
         .expected = "my_string: 'foo' maybe_int32: 7"},
        {.name = "UnsetOptionalOverwrittenByDefault",
         .defaults = "my_string: 'default_value'",
         .params = "",
         .expected = "my_string: 'default_value'"},
        {.name = "SetOptionalNotOverwrittenByDefault",
         .defaults = "my_string: 'default_value'",
         .params = "my_string: ''",
         .expected = "my_string: ''"},
        {.name = "NonOptionalSetToDefaultOverwrittenByDefault",
         .defaults = "my_non_optional_int: 42",
         .params = "my_non_optional_int: 0",
         .expected = "my_non_optional_int: 42"},
        {.name = "UnsetOneofOverwrittenByDefault",
         .defaults = "maybe_int32: 100",
         .params = "",
         .expected = "maybe_int32: 100"},
        {.name = "SetOneofNotOverwrittenByDefault",
         .defaults = "maybe_int32: 100",
         .params = "maybe_int64: 0",
         .expected = "maybe_int64: 0"},
        {.name = "UnsetOptionalBoolOverwrittenByDefault",
         .defaults = "my_optional_bool: false",
         .params = "",
         .expected = "my_optional_bool: false"},
        {.name = "SetOptionalBoolNotOverwrittenByDefault",
         .defaults = "my_optional_bool: true",
         .params = "my_optional_bool: false",
         .expected = "my_optional_bool: false"},
        {.name = "NonOptionalBoolSetToDefaultOverwrittenByDefault",
         .defaults = "my_non_optional_bool: true",
         .params = "my_non_optional_bool: false",
         .expected = "my_non_optional_bool: true"},
    }),
    [](const ::testing::TestParamInfo<UnpackAnyAndMergeTestCase>& info) {
      return info.param.name;
    });

}  // namespace
}  // namespace intrinsic
