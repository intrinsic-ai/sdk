// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/proto/parsed_type_url.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "absl/status/status.h"
#include "google/protobuf/wrappers.pb.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

using ::absl_testing::StatusIs;

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

}  // namespace
}  // namespace intrinsic
