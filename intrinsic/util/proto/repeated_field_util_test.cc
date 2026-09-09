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

#include "intrinsic/util/proto/repeated_field_util.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <string>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "google/protobuf/repeated_ptr_field.h"
#include "intrinsic/util/proto/testing/test_message.pb.h"

namespace my_test {
enum MyProtoEnum { kVal1 = 1, kVal2 = 2 };
enum class MyEnum { kVal1 = 1, kVal2 = 2 };

inline MyEnum FromProto(MyProtoEnum val) { return static_cast<MyEnum>(val); }
inline MyProtoEnum ToProto(MyEnum val) { return static_cast<MyProtoEnum>(val); }

enum MyProtoStatusEnum { kValid = 1, kInvalid = 2 };
enum class MyStatusEnum { kValid = 1, kInvalid = 2 };

inline absl::StatusOr<MyStatusEnum> FromProto(MyProtoStatusEnum val) {
  if (val == kValid) {
    return MyStatusEnum::kValid;
  }
  return absl::InvalidArgumentError("Invalid enum value");
}

inline MyProtoStatusEnum ToProto(MyStatusEnum val) {
  return static_cast<MyProtoStatusEnum>(val);
}

}  // namespace my_test

template <>
struct google::protobuf::is_proto_enum<my_test::MyProtoEnum> : std::true_type {
};
template <>
struct google::protobuf::is_proto_enum<my_test::MyProtoStatusEnum>
    : std::true_type {};

namespace intrinsic_proto::util::proto::testing {

inline std::string FromProto(const NestedMessage& msg) { return msg.value(); }

inline NestedMessage ToProto(const std::string& str) {
  NestedMessage msg;
  msg.set_value(str);
  return msg;
}

}  // namespace intrinsic_proto::util::proto::testing

namespace intrinsic {
namespace {

using ::absl_testing::StatusIs;
using ::testing::ElementsAreArray;

TEST(RepeatedFieldUtilTest, RemoveIf) {
  ::google::protobuf::RepeatedPtrField<std::string> strings;
  strings.Add("apple");
  strings.Add("banana");
  strings.Add("cherry");
  strings.Add("avocado");

  int removed =
      RemoveIf(&strings, [](const std::string* s) { return (*s)[0] == 'a'; });
  EXPECT_EQ(removed, 2);
  EXPECT_THAT(strings, ElementsAreArray({"banana", "cherry"}));
}

TEST(RepeatedFieldUtilTest, Sort) {
  ::google::protobuf::RepeatedPtrField<std::string> strings;
  strings.Add("cherry");
  strings.Add("apple");
  strings.Add("banana");

  Sort(&strings);
  EXPECT_THAT(strings, ElementsAreArray({"apple", "banana", "cherry"}));
}

TEST(RepeatedFieldUtilTest, RepeatedFieldToVector) {
  ::google::protobuf::RepeatedField<int> repeated_ints;
  repeated_ints.Add(1);
  repeated_ints.Add(2);
  repeated_ints.Add(3);
  repeated_ints.Add(5);
  repeated_ints.Add(8);

  const auto ints = AsVector(repeated_ints);
  EXPECT_THAT(repeated_ints, ElementsAreArray(ints));
}

TEST(RepeatedFieldUtilTest, EmptyRepeatedFieldToVector) {
  ::google::protobuf::RepeatedField<int> repeated_ints;
  const auto ints = AsVector(repeated_ints);
  EXPECT_TRUE(ints.empty());
  EXPECT_THAT(repeated_ints, ElementsAreArray(ints));
}

TEST(RepeatedFieldUtilTest, RepeatedPtrFieldToVector) {
  ::google::protobuf::RepeatedPtrField<std::string> repeated_strings;
  repeated_strings.Add("1");
  repeated_strings.Add("2");
  repeated_strings.Add("3");
  repeated_strings.Add("5");
  repeated_strings.Add("8");

  const auto strings = AsVector(repeated_strings);
  EXPECT_THAT(repeated_strings, ElementsAreArray(strings));
}

TEST(RepeatedFieldUtilTest, EmptyRepeatedPtrFieldToVector) {
  ::google::protobuf::RepeatedPtrField<std::string> repeated_strings;
  const auto strings = AsVector(repeated_strings);
  EXPECT_TRUE(strings.empty());
  EXPECT_THAT(repeated_strings, ElementsAreArray(strings));
}

TEST(RepeatedFieldUtilTest, VectorToRepeatedField) {
  std::vector<int> ints = {1, 2, 3, 5, 8};
  const auto repeated_ints = AsRepeatedField(ints);
  EXPECT_THAT(ints, ElementsAreArray(repeated_ints));
}

TEST(RepeatedFieldUtilTest, EmptyVectorToRepeatedField) {
  std::vector<int> ints = {};
  const auto repeated_ints = AsRepeatedField(ints);
  EXPECT_TRUE(repeated_ints.empty());
  EXPECT_THAT(ints, ElementsAreArray(repeated_ints));
}

TEST(RepeatedFieldUtilTest, VectorToRepeatedPtrField) {
  std::vector<std::string> strings = {"1", "2", "3", "5", "8"};
  const auto repeated_strings = AsRepeatedPtrField(strings);
  EXPECT_THAT(strings, ElementsAreArray(repeated_strings));
}

TEST(RepeatedFieldUtilTest, EmptyVectorToRepeatedPtrField) {
  std::vector<std::string> strings = {};
  const auto repeated_strings = AsRepeatedPtrField(strings);
  EXPECT_TRUE(repeated_strings.empty());
  EXPECT_THAT(strings, ElementsAreArray(repeated_strings));
}

TEST(RepeatedFieldUtilTest, FromRepeatedFieldMessage) {
  const std::vector<std::string> expected{"hello", "world", "test"};
  ::google::protobuf::RepeatedPtrField<
      intrinsic_proto::util::proto::testing::NestedMessage>
      repeated_field;
  for (const auto& s : expected) {
    repeated_field.Add(intrinsic_proto::util::proto::testing::ToProto(s));
  }
  EXPECT_EQ(FromProto(repeated_field), expected);
}

TEST(RepeatedFieldUtilTest, FromRepeatedFieldEnum) {
  const std::vector<my_test::MyEnum> expected{
      my_test::MyEnum::kVal1,
      my_test::MyEnum::kVal2,
  };
  ::google::protobuf::RepeatedField<my_test::MyProtoEnum> repeated_field;
  for (const auto& type : expected) {
    repeated_field.Add(my_test::ToProto(type));
  }
  EXPECT_EQ(FromProto(repeated_field), expected);
}

TEST(RepeatedFieldUtilTest, FromRepeatedFieldWithStatus) {
  const std::vector<my_test::MyStatusEnum> expected{
      my_test::MyStatusEnum::kValid,
  };
  ::google::protobuf::RepeatedField<my_test::MyProtoStatusEnum> repeated_field;
  for (const auto& type : expected) {
    repeated_field.Add(my_test::ToProto(type));
  }
  EXPECT_THAT(FromProto(repeated_field),
              ::absl_testing::IsOkAndHolds(expected));

  repeated_field.Add(my_test::ToProto(my_test::MyStatusEnum::kInvalid));
  EXPECT_THAT(FromProto(repeated_field),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

}  // namespace
}  // namespace intrinsic
