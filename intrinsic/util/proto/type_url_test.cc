// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/proto/type_url.h"

#include <gtest/gtest.h>

#include <string>

#include "absl/strings/string_view.h"
#include "google/protobuf/wrappers.pb.h"

namespace intrinsic {

namespace {

TEST(TypeUrl, AddPrefix) {
  EXPECT_EQ(AddTypeUrlPrefix("google.protobuf.Int64Value"),
            "type.googleapis.com/google.protobuf.Int64Value");
}

TEST(TypeUrl, AddPrefixCustomWithSlash) {
  EXPECT_EQ(
      AddTypeUrlPrefix("google.protobuf.Int64Value", "type.intrinsic.ai/test/"),
      "type.intrinsic.ai/test/google.protobuf.Int64Value");
}

TEST(TypeUrl, AddPrefixCustomWithoutSlash) {
  EXPECT_EQ(
      AddTypeUrlPrefix("google.protobuf.Int64Value", "type.intrinsic.ai/test"),
      "type.intrinsic.ai/test/google.protobuf.Int64Value");
}

TEST(TypeUrl, AddPrefixIdempotent) {
  std::string type_url = "type.googleapis.com/google.protobuf.Int64Value";
  EXPECT_EQ(AddTypeUrlPrefix(type_url), type_url);
}

TEST(TypeUrl, AddPrefixType) {
  EXPECT_EQ(AddTypeUrlPrefix<google::protobuf::Int64Value>(),
            "type.googleapis.com/google.protobuf.Int64Value");
}

TEST(TypeUrl, AddPrefixMessageReference) {
  google::protobuf::Int64Value m;
  EXPECT_EQ(AddTypeUrlPrefix(m),
            "type.googleapis.com/google.protobuf.Int64Value");
}

TEST(TypeUrl, AddPrefixMessagePointer) {
  google::protobuf::Int64Value m;
  EXPECT_EQ(AddTypeUrlPrefix(&m),
            "type.googleapis.com/google.protobuf.Int64Value");
}

TEST(TypeUrl, AddPrefixTypeCustomWithSlash) {
  EXPECT_EQ(
      AddTypeUrlPrefix<google::protobuf::Int64Value>("type.intrinsic.ai/test/"),
      "type.intrinsic.ai/test/google.protobuf.Int64Value");
}

TEST(TypeUrl, AddPrefixTypeCustomWithoutSlash) {
  EXPECT_EQ(
      AddTypeUrlPrefix<google::protobuf::Int64Value>("type.intrinsic.ai/test"),
      "type.intrinsic.ai/test/google.protobuf.Int64Value");
}

TEST(TypeUrl, AddPrefixMessageReferenceCustomWithSlash) {
  google::protobuf::Int64Value m;
  EXPECT_EQ(AddTypeUrlPrefix(m, "type.intrinsic.ai/test/"),
            "type.intrinsic.ai/test/google.protobuf.Int64Value");
}

TEST(TypeUrl, AddPrefixMessageReferenceCustomWithoutSlash) {
  google::protobuf::Int64Value m;
  EXPECT_EQ(AddTypeUrlPrefix(m, "type.intrinsic.ai/test"),
            "type.intrinsic.ai/test/google.protobuf.Int64Value");
}

TEST(TypeUrl, AddPrefixMessagePointerCustomWithSlash) {
  google::protobuf::Int64Value m;
  EXPECT_EQ(AddTypeUrlPrefix(&m, "type.intrinsic.ai/test/"),
            "type.intrinsic.ai/test/google.protobuf.Int64Value");
}

TEST(TypeUrl, AddPrefixMessagePointerCustomWithoutSlash) {
  google::protobuf::Int64Value m;
  EXPECT_EQ(AddTypeUrlPrefix(&m, "type.intrinsic.ai/test"),
            "type.intrinsic.ai/test/google.protobuf.Int64Value");
}

TEST(TypeUrl, StripPrefix) {
  EXPECT_EQ(
      StripTypeUrlPrefix("type.googleapis.com/google.protobuf.Int64Value"),
      "google.protobuf.Int64Value");
}

TEST(TypeUrl, StripPrefixIdempotent) {
  std::string proto_type = "google.protobuf.Int64Value";
  EXPECT_EQ(StripTypeUrlPrefix(proto_type), proto_type);
}

TEST(TypeUrl, GenerateIntrinsicTypeUrl) {
  std::string proto_type = "google.protobuf.Int64Value";
  EXPECT_EQ(GenerateIntrinsicTypeUrl("area", "foo", "bar", proto_type),
            "type.intrinsic.ai/area/foo/bar/google.protobuf.Int64Value");
  EXPECT_EQ(GenerateIntrinsicTypeUrl("foo", 25), "type.intrinsic.ai/foo/25");
}

TEST(TypeUrl, GenerateIntrinsicTypeUrlForMessage) {
  EXPECT_EQ(GenerateIntrinsicTypeUrlForMessage<google::protobuf::Int64Value>(
                "area", "foo", "bar"),
            "type.intrinsic.ai/area/foo/bar/google.protobuf.Int64Value");
}

}  // namespace
}  // namespace intrinsic
