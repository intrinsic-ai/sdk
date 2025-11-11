// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/proto/parsed_type_url.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "absl/status/status.h"
#include "google/protobuf/wrappers.pb.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

using ::absl_testing::StatusIs;
using ::testing::HasSubstr;

namespace intrinsic {

namespace {

TEST(ParsedTypeUrl, ParseTypeUrlGoogleSpec) {
  EXPECT_THAT(ParseTypeUrl("type.googleapis.com/google.protobuf.Int64Value"),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(TypeUrl, ParseTypeUrl) {
  ASSERT_OK_AND_ASSIGN(
      ParsedUrl pu,
      ParseTypeUrl(
          "type.intrinsic.ai/area/foo/bar/google.protobuf.Int64Value"));

  EXPECT_EQ(pu.area, "area");
  EXPECT_EQ(pu.path, "foo/bar");
  EXPECT_EQ(pu.message_type, "google.protobuf.Int64Value");
}

TEST(ParsedTypeUrl, ParseTypeUrlPrefixGoogleSpec) {
  EXPECT_THAT(
      ParseTypeUrlPrefix("type.googleapis.com/google.protobuf.Int64Value"),
      StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(TypeUrl, ParseTypeUrlInvalid) {
  EXPECT_THAT(
      ParseTypeUrl("type.intrinsic.ai///"),
      StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("missing area")));
  EXPECT_THAT(
      ParseTypeUrl("type.intrinsic.ai/area//"),
      StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("missing path")));
  EXPECT_THAT(
      ParseTypeUrl("type.intrinsic.ai/area/path/"),
      StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("message type")));
  EXPECT_THAT(
      ParseTypeUrl("type.intrinsic.ai//asd/"),
      StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("missing area")));
}

TEST(TypeUrl, ParseTypeUrlPrefix) {
  ASSERT_OK_AND_ASSIGN(ParsedUrl pu,
                       ParseTypeUrlPrefix("type.intrinsic.ai/area/foo/bar"));

  EXPECT_EQ(pu.area, "area");
  EXPECT_EQ(pu.path, "foo/bar");
  EXPECT_EQ(pu.message_type, "");
}

TEST(TypeUrl, ParseTypeUrlPrefixWithSlash) {
  ASSERT_OK_AND_ASSIGN(ParsedUrl pu,
                       ParseTypeUrlPrefix("type.intrinsic.ai/area/foo/bar/"));

  EXPECT_EQ(pu.area, "area");
  EXPECT_EQ(pu.path, "foo/bar");
  EXPECT_EQ(pu.message_type, "");
}

TEST(TypeUrl, ParseTypeUrlPrefixAmbiguousMessage) {
  ASSERT_OK_AND_ASSIGN(
      ParsedUrl pu,
      ParseTypeUrlPrefix(
          // This could be a type URL as well as a type URL prefix. There is no
          // way to differentiate the last part. Thus the specific
          // ParseTypeUrlPrefix call is necessary and callers must know what
          // they are querying for.
          "type.intrinsic.ai/area/foo/bar/google.protobuf.Int64Value"));

  EXPECT_EQ(pu.area, "area");
  EXPECT_EQ(pu.path, "foo/bar/google.protobuf.Int64Value");
  EXPECT_EQ(pu.message_type, "");
}

TEST(TypeUrl, ParseTypeUrlPrefixWithoutArea) {
  EXPECT_THAT(
      ParseTypeUrlPrefix("type.intrinsic.ai///"),
      StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("missing area")));
  EXPECT_THAT(
      ParseTypeUrlPrefix("type.intrinsic.ai/area//"),
      StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("missing path")));
  EXPECT_THAT(
      ParseTypeUrlPrefix("type.intrinsic.ai//asd/"),
      StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("missing area")));
}

}  // namespace
}  // namespace intrinsic
