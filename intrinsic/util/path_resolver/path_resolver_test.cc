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

#include "intrinsic/util/path_resolver/path_resolver.h"

#include <gtest/gtest.h>
#include <unistd.h>

#include <string>
#include <string_view>

#include "absl/strings/match.h"

namespace intrinsic {
namespace {

TEST(PathResolver, Resolve) {
  std::string run_file = "intrinsic/util/path_resolver/path_resolver_test.cc";

  std::string full_path =
      PathResolver::ResolveRunfilesPath(std::string_view(run_file));
  EXPECT_NE(full_path.size(), 0);
  EXPECT_TRUE(absl::EndsWith(full_path, run_file));
  EXPECT_EQ(access(full_path.c_str(), F_OK), 0);
}

TEST(PathResolver, Resolve_NonExistentFile) {
  std::string run_file = "intrinsic/util/path_resolver/nonexistent_file.txt";

  std::string full_path =
      PathResolver::ResolveRunfilesPath(std::string_view(run_file));
  // Rlocation might return an empty string if the mapping fails, or a path
  // that doesn't exist on disk.
  EXPECT_NE(full_path, run_file);
  if (!full_path.empty()) {
    EXPECT_NE(access(full_path.c_str(), F_OK), 0);
  }
}

TEST(PathResolver, Resolve_MainModulePrefix) {
  std::string run_file =
      "_main/intrinsic/util/path_resolver/path_resolver_test.cc";

  std::string full_path =
      PathResolver::ResolveRunfilesPath(std::string_view(run_file));
  // Could be empty or a valid path depending on module setup, just ensure it
  // handles the _main/ prefix without crashing.
  EXPECT_NE(full_path, run_file);
}

TEST(PathResolver, Resolve_FallbackPrefix) {
  std::string run_file =
      "some_other_repo/util/path_resolver/path_resolver_test.cc";

  std::string full_path =
      PathResolver::ResolveRunfilesPath(std::string_view(run_file));
  EXPECT_NE(full_path, run_file);
}

TEST(PathResolver, ResolveRunfilesPathForTest) {
  std::string run_file = "intrinsic/util/path_resolver/path_resolver_test.cc";

  std::string full_path =
      PathResolver::ResolveRunfilesPathForTest(std::string_view(run_file));
  EXPECT_NE(full_path.size(), 0);
  EXPECT_TRUE(absl::EndsWith(full_path, run_file));
  EXPECT_EQ(access(full_path.c_str(), F_OK), 0);
}

TEST(PathResolver, ResolveRunfilesPathIfRelative) {
  EXPECT_EQ(PathResolver::ResolveRunfilesPathIfRelative("/absolute/path"),
            "/absolute/path");
  EXPECT_EQ(PathResolver::ResolveRunfilesPathIfRelative("relative/path"),
            PathResolver::ResolveRunfilesPath("relative/path"));
}

}  // namespace
}  // namespace intrinsic
