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

#include "intrinsic/platform/pubsub/fake_kvstore.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/wrappers.pb.h"
#include "internal/testing.h"

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
