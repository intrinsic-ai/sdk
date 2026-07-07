// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/proto/parsed_type_url.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "absl/status/status.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

using ::absl_testing::StatusIs;
using ::testing::HasSubstr;

namespace intrinsic {

namespace {

TEST(ParsedTypeUrl, ParseTypeUrlGoogleSpec) {
  ASSERT_OK_AND_ASSIGN(
      ParsedUrl pu,
      ParseTypeUrl("type.googleapis.com/google.protobuf.Int64Value"));

  EXPECT_EQ(pu.prefix, "type.intrinsic.ai/");
  EXPECT_EQ(pu.area, "common");
  EXPECT_EQ(pu.path, "");
  EXPECT_EQ(pu.message_type, "google.protobuf.Int64Value");
}

TEST(ParsedTypeUrl, ParseTypeUrlCommon) {
  ASSERT_OK_AND_ASSIGN(
      ParsedUrl pu,
      ParseTypeUrl("type.intrinsic.ai/common/google.protobuf.Int64Value"));

  EXPECT_EQ(pu.prefix, "type.intrinsic.ai/");
  EXPECT_EQ(pu.area, "common");
  EXPECT_EQ(pu.path, "");
  EXPECT_EQ(pu.message_type, "google.protobuf.Int64Value");
}

TEST(ParsedTypeUrl, ParseTypeUrlWellKnownAlias) {
  ASSERT_OK_AND_ASSIGN(
      ParsedUrl pu,
      ParseTypeUrl("type.intrinsic.ai/well-known/google.protobuf.Int64Value"));

  EXPECT_EQ(pu.prefix, "type.intrinsic.ai/");
  EXPECT_EQ(pu.area, "common");
  EXPECT_EQ(pu.path, "");
  EXPECT_EQ(pu.message_type, "google.protobuf.Int64Value");
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

TEST(TypeUrl, ParseTypeUrlEmptyPath) {
  ASSERT_OK_AND_ASSIGN(
      ParsedUrl pu,
      ParseTypeUrl("type.intrinsic.ai/area/google.protobuf.Int64Value"));

  EXPECT_EQ(pu.area, "area");
  EXPECT_EQ(pu.path, "");
  EXPECT_EQ(pu.message_type, "google.protobuf.Int64Value");
}

TEST(TypeUrl, ParseTypeUrlInvalid) {
  EXPECT_THAT(
      ParseTypeUrl("type.intrinsic.ai///"),
      StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("missing area")));
  EXPECT_THAT(ParseTypeUrl("type.intrinsic.ai/area//"),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("missing message type")));
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

TEST(TypeUrl, ParseTypeUrlPrefixEmptyPath) {
  ASSERT_OK_AND_ASSIGN(ParsedUrl pu,
                       ParseTypeUrlPrefix("type.intrinsic.ai/area/"));

  EXPECT_EQ(pu.area, "area");
  EXPECT_EQ(pu.path, "");
  EXPECT_EQ(pu.message_type, "");
}

TEST(TypeUrl, ParseTypeUrlPrefixNoTrailingSlashEmptyPath) {
  ASSERT_OK_AND_ASSIGN(ParsedUrl pu,
                       ParseTypeUrlPrefix("type.intrinsic.ai/area"));

  EXPECT_EQ(pu.area, "area");
  EXPECT_EQ(pu.path, "");
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
      ParseTypeUrlPrefix("type.intrinsic.ai//asd/"),
      StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("missing area")));

  ASSERT_OK_AND_ASSIGN(ParsedUrl pu,
                       ParseTypeUrlPrefix("type.intrinsic.ai/area//"));
  EXPECT_EQ(pu.area, "area");
  EXPECT_EQ(pu.path, "");
}

}  // namespace
}  // namespace intrinsic
