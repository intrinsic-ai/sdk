// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/platform/pubsub/fake_kvstore.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "absl/status/status.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/wrappers.pb.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic {
namespace {

using ::testing::HasSubstr;
using ::testing::Test;

class FakeKeyValueStoreTest : public Test {
 protected:
  FakeKeyValueStore kv_store_;
};

TEST_F(FakeKeyValueStoreTest, SetAndGetSuccess) {
  google::protobuf::StringValue val;
  val.set_value("test_value");

  EXPECT_OK(kv_store_.Set("my/key", val));

  ASSERT_OK_AND_ASSIGN(auto retrieved,
                       kv_store_.Get<google::protobuf::StringValue>("my/key"));
  EXPECT_EQ(retrieved.value(), "test_value");
}

TEST_F(FakeKeyValueStoreTest, GetNotFound) {
  auto status =
      kv_store_.Get<google::protobuf::StringValue>("nonexistent").status();
  EXPECT_EQ(status.code(), absl::StatusCode::kNotFound);
  EXPECT_THAT(status.message(), HasSubstr("Key not found: nonexistent"));
}

TEST_F(FakeKeyValueStoreTest, DeleteSuccess) {
  google::protobuf::StringValue val;
  val.set_value("test_value");

  EXPECT_OK(kv_store_.Set("my/key", val));
  EXPECT_OK(kv_store_.Delete("my/key"));

  auto status = kv_store_.Get<google::protobuf::StringValue>("my/key").status();
  EXPECT_EQ(status.code(), absl::StatusCode::kNotFound);
}

TEST_F(FakeKeyValueStoreTest, DeleteNotFound) {
  auto status = kv_store_.Delete("nonexistent");
  EXPECT_EQ(status.code(), absl::StatusCode::kNotFound);
  EXPECT_THAT(status.message(), HasSubstr("Key not found: nonexistent"));
}

}  // namespace
}  // namespace intrinsic
