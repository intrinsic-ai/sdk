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

#include "intrinsic/util/testing/runfiles.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include <unistd.h>

#include <string>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "absl/status/statusor.h"
#include "internal/testing.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::testing {
namespace {

using ::absl_testing::StatusIs;
using ::testing::EndsWith;

absl::StatusOr<std::string> ResolveTestFileViaStatusMacro() {
  INTR_ASSIGN_OR_RETURN(std::string path,
                        ResolveRlocation(BAZEL_CURRENT_REPOSITORY,
                                         "ai_intrinsic_sdks/intrinsic/"
                                         "util/testing/runfiles_test.cc"));
  return path;
}

TEST(ResolveRlocationTest, ResolvesRunfileWithAssertOkAndAssign) {
  ASSERT_OK_AND_ASSIGN(std::string path,
                       ResolveRlocation(BAZEL_CURRENT_REPOSITORY,
                                        "ai_intrinsic_sdks/intrinsic/"
                                        "util/testing/runfiles_test.cc"));
  EXPECT_THAT(path, EndsWith("intrinsic/util/testing/runfiles_test.cc"));
  EXPECT_EQ(access(path.c_str(), F_OK), 0);
}

TEST(ResolveRlocationTest, ResolvesRunfileWithIntrAssignOrReturn) {
  ASSERT_OK_AND_ASSIGN(std::string path, ResolveTestFileViaStatusMacro());
  EXPECT_THAT(path, EndsWith("intrinsic/util/testing/runfiles_test.cc"));
  EXPECT_EQ(access(path.c_str(), F_OK), 0);
}

TEST(ResolveRlocationTest, InvalidRlocationReturnsError) {
  EXPECT_THAT(ResolveRlocation(BAZEL_CURRENT_REPOSITORY, ""),
              StatusIs(absl::StatusCode::kNotFound));
  EXPECT_THAT(ResolveRlocation(BAZEL_CURRENT_REPOSITORY, "../invalid"),
              StatusIs(absl::StatusCode::kNotFound));
  EXPECT_THAT(ResolveRlocation(BAZEL_CURRENT_REPOSITORY, "repo//double_slash"),
              StatusIs(absl::StatusCode::kNotFound));
}

}  // namespace
}  // namespace intrinsic::testing
