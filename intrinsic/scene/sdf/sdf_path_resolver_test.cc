// Copyright 2023 Intrinsic Innovation LLC

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
