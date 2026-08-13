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

#include "intrinsic/scene/sdf/sdf_path_resolver.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <string>

#include "intrinsic/util/path_resolver/path_resolver.h"
#include "intrinsic/util/testing/gtest_wrapper.h"
#include "ortools/base/path.h"

namespace intrinsic {
namespace sdf {
namespace {

using ::absl_testing::IsOkAndHolds;

TEST(SdfPathResolver, Resolve) {
  const std::string run_file =
      "intrinsic/util/path_resolver/path_resolver_test.cc";

  const std::string full_path = PathResolver::ResolveRunfilesPath(run_file);

  EXPECT_THAT(SdfPathResolver("model://" + run_file), IsOkAndHolds(full_path));
  EXPECT_THAT(SdfPathResolver(file::JoinPath("model://", run_file)),
              IsOkAndHolds(PathResolver::ResolveRunfilesPath(run_file)));
  EXPECT_THAT(SdfPathResolver(full_path), IsOkAndHolds(full_path));
  EXPECT_THAT(SdfPathResolver("bypass://some_random_thing.json"),
              IsOkAndHolds("some_random_thing.json"));
}

}  // namespace
}  // namespace sdf
}  // namespace intrinsic
