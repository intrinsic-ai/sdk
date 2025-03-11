// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/proto/parse_text_proto.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "google/protobuf/type.pb.h"
#include "google/protobuf/wrappers.pb.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic {
namespace {

using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;
using ::intrinsic::testing::EqualsProto;
using ::testing::AllOf;
using ::testing::Eq;
using ::testing::HasSubstr;
using ::testing::Property;

TEST(ParseTextProtoTest, ParseTextProtoInto) {
  google::protobuf::Int32Value int32_value;

  EXPECT_THAT(ParseTextProtoInto("value: 1", &int32_value),
              ::absl_testing::IsOk());

  EXPECT_THAT(int32_value, EqualsProto("value: 1"));
}

TEST(ParseTextProtoTest, ParseTextProtoIntoFails) {
  google::protobuf::Int32Value int32_value;

  EXPECT_THAT(
      ParseTextProtoInto("non_existent_field: 1", &int32_value),
      StatusIs(
          absl::StatusCode::kInvalidArgument,
          AllOf(HasSubstr("Cannot parse protobuf google.protobuf.Int32Value"),
                HasSubstr("non_existent_field"))));
}

TEST(ParseTextProtoTest, ParseTextProto) {
  EXPECT_THAT(ParseTextProto<google::protobuf::Int32Value>("value: 1"),
              IsOkAndHolds(EqualsProto("value: 1")));
}

TEST(ParseTextProtoTest, ParseTextProtoWorksWithCustomAnyTypeUrlPrefix) {
  google::protobuf::Int32Value int32_value;
  int32_value.set_value(1);

  // Use Option because it is a well-known type which has an Any field.
  google::protobuf::Option option;
  option.mutable_value()->PackFrom(int32_value);

  absl::StatusOr<google::protobuf::Option> actual_option =
      ParseTextProto<google::protobuf::Option>(R"pb(
        value: {
          [type.intrinsic.ai/google.protobuf.Int32Value] { value: 1 }
        }
      )pb");
  EXPECT_THAT(actual_option, IsOkAndHolds(EqualsProto(option)));
  EXPECT_THAT(
      actual_option,
      IsOkAndHolds(Property(
          &google::protobuf::Option::value,
          Property(&google::protobuf::Any::type_url,
                   Eq("type.intrinsic.ai/google.protobuf.Int32Value")))));

  actual_option = ParseTextProto<google::protobuf::Option>(R"pb(
    value: {
      [type.intrinsic.ai/skills/google.protobuf.Int32Value] { value: 1 }
    }
  )pb");
  EXPECT_THAT(actual_option, IsOkAndHolds(EqualsProto(option)));
  EXPECT_THAT(
      actual_option,
      IsOkAndHolds(Property(
          &google::protobuf::Option::value,
          Property(
              &google::protobuf::Any::type_url,
              Eq("type.intrinsic.ai/skills/google.protobuf.Int32Value")))));

  actual_option = ParseTextProto<google::protobuf::Option>(R"pb(
    value: {
      [type.intrinsic.ai/skills/ai.intrinsic.test/google.protobuf.Int32Value] {
        value: 1
      }
    }
  )pb");
  EXPECT_THAT(actual_option, IsOkAndHolds(EqualsProto(option)));
  EXPECT_THAT(actual_option,
              IsOkAndHolds(Property(
                  &google::protobuf::Option::value,
                  Property(&google::protobuf::Any::type_url,
                           Eq("type.intrinsic.ai/skills/ai.intrinsic.test/"
                              "google.protobuf.Int32Value")))));

  actual_option = ParseTextProto<google::protobuf::Option>(R"pb(
    value: {
      [type.intrinsic.ai/skills/ai.intrinsic.test/0.0.1/
       google.protobuf.Int32Value] { value: 1 }
    }
  )pb");
  EXPECT_THAT(actual_option, IsOkAndHolds(EqualsProto(option)));
  EXPECT_THAT(actual_option,
              IsOkAndHolds(Property(
                  &google::protobuf::Option::value,
                  Property(&google::protobuf::Any::type_url,
                           Eq("type.intrinsic.ai/skills/ai.intrinsic.test/"
                              "0.0.1/google.protobuf.Int32Value")))));

  actual_option = ParseTextProto<google::protobuf::Option>(R"(
    value: {
      [type.intrinsic.ai/skills/ai.intrinsic.test/0.0.1-alpha-0aZ+buildspec/google.protobuf.Int32Value] {
        value: 1
      }
    }
  )");
  EXPECT_THAT(actual_option, IsOkAndHolds(EqualsProto(option)));
  EXPECT_THAT(
      actual_option,
      IsOkAndHolds(Property(
          &google::protobuf::Option::value,
          Property(
              &google::protobuf::Any::type_url,
              Eq("type.intrinsic.ai/skills/ai.intrinsic.test/"
                 "0.0.1-alpha-0aZ+buildspec/google.protobuf.Int32Value")))));

  // Test multiple Any type URLs in one text proto.
  // Use Type because it is a well-known type which has a repeated Option
  // field (and each Option has an Any field).
  google::protobuf::Type type_with_two_options;
  type_with_two_options.add_options()->mutable_value()->PackFrom(int32_value);
  type_with_two_options.add_options()->mutable_value()->PackFrom(int32_value);
  EXPECT_THAT(
      ParseTextProto<google::protobuf::Type>(R"pb(
        options: {
          value: {
            [type.intrinsic.ai/skills/google.protobuf.Int32Value] { value: 1 }
          }
        }
        options: {
          value: {
            [type.intrinsic.ai/skills/0.0.1/google.protobuf.Int32Value] {
              value: 1
            }
          }
        }
      )pb"),
      IsOkAndHolds(EqualsProto(type_with_two_options)));

  // Test that the list syntax ("options: [...]") does not break anything.
  EXPECT_THAT(
      ParseTextProto<google::protobuf::Type>(R"pb(
        options:
        [ {
          value: {
            [type.intrinsic.ai/skills/google.protobuf.Int32Value] { value: 1 }
          }
        }
          , {
            value: {
              [type.intrinsic.ai/skills/0.0.1/google.protobuf.Int32Value] {
                value: 1
              }
            }
          }]
      )pb"),
      IsOkAndHolds(EqualsProto(type_with_two_options)));

  // Test that type URLs of nested Any protos work.
  google::protobuf::Option inner_option, outer_option;
  inner_option.mutable_value()->PackFrom(int32_value);
  outer_option.mutable_value()->PackFrom(inner_option);
  EXPECT_THAT(
      ParseTextProto<google::protobuf::Option>(R"pb(
        value: {
          [type.intrinsic.ai/skills/0.0.1/google.protobuf.Option] {
            value: {
              [type.intrinsic.ai/skills/0.0.1/google.protobuf.Int32Value] {
                value: 1
              }
            }
          }
        }
      )pb"),
      IsOkAndHolds(EqualsProto(outer_option)));

  // Note: Type URLs of nested Anys will be parsed, but contain replacements
  // like _DOT_.
}

TEST(ParseTextProtoTest, ParseTextProtoExtensions) {
  // Test that extensions are parsed correctly and not interfered with by the
  // rewriting of Any type URLs. The following should give an error about the
  // extension not being defined and not a syntax error.
  EXPECT_THAT(ParseTextProto<google::protobuf::Type>(R"pb(
                options: {
                  [com.example.extension_field]: 20
                }
              )pb"),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       AllOf(HasSubstr("com.example.extension_field"),
                             HasSubstr("is not defined"))));
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
