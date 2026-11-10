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

#include "intrinsic/util/proto/descriptor_pools.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <string>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "google/protobuf/descriptor.h"

namespace intrinsic {
namespace {

using ::absl_testing::IsOk;
using ::absl_testing::StatusIs;
using ::google::protobuf::DescriptorPool;
using ::testing::Eq;
using ::testing::IsEmpty;

TEST(DescriptorPoolErrorCollectorTest, StatusIsOkWithoutErrors) {
  DescriptorPoolErrorCollector error_collector;

  EXPECT_THAT(error_collector.GetErrorsJoined(), IsEmpty());
  EXPECT_THAT(error_collector.GetStatus(), IsOk());
}

TEST(DescriptorPoolErrorCollectorTest, StatusContainsAllErrors) {
  DescriptorPoolErrorCollector error_collector;
  error_collector.RecordError("a.proto", "pkg.Enum", /*descriptor=*/nullptr,
                              DescriptorPool::ErrorCollector::NAME,
                              "Enums must contain at least one value.");
  error_collector.RecordError("b.proto", "pkg.Message.field",
                              /*descriptor=*/nullptr,
                              DescriptorPool::ErrorCollector::TYPE,
                              "\"Undefined\" is not defined.");

  std::string expected_message =
      "a.proto: pkg.Enum: Enums must contain at least one value.; b.proto: "
      "pkg.Message.field: \"Undefined\" is not defined.";
  EXPECT_THAT(error_collector.GetErrorsJoined(), Eq(expected_message));
  EXPECT_THAT(error_collector.GetStatus(absl::StatusCode::kInternal),
              StatusIs(absl::StatusCode::kInternal, expected_message));
}

}  // namespace
}  // namespace intrinsic
