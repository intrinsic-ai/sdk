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

#include <cstdlib>
#include <string>

#include "absl/debugging/leak_check.h"
#include "absl/log/log.h"
#include "absl/strings/match.h"
#include "absl/strings/string_view.h"
#include "ortools/base/path.h"
#include "tools/cpp/runfiles/runfiles.h"

namespace intrinsic {

using bazel::tools::cpp::runfiles::Runfiles;

std::string BuildFullyQualifiedRunfilesPath(absl::string_view path) {
  constexpr absl::string_view kRepoName =
  "ai_intrinsic_sdks";
  return file::JoinPath(kRepoName, path);
}

std::string PathResolver::ResolveRunfilesPath(absl::string_view path) {
  std::string error;
  auto runfiles = std::unique_ptr<Runfiles>(
      Runfiles::Create(program_invocation_name, &error));
  if (runfiles == nullptr) {
    LOG(ERROR) << "Error creating Runfiles object: " << error;
    return "";
  }

  std::string full_path = BuildFullyQualifiedRunfilesPath(path);
  return runfiles->Rlocation(full_path);
}

std::string PathResolver::ResolveRunfilesPathForTest(absl::string_view path) {
  if (!getenv("BAZEL_TEST")) {
    return ResolveRunfilesPath(path);
  }

  std::string error;
  auto runfiles = Runfiles::CreateForTest(&error);
  if (runfiles == nullptr) {
    LOG(ERROR) << "Error creating Runfiles object for test: " << error;
    return "";
  }

  absl::IgnoreLeak(runfiles);

  std::string full_path = BuildFullyQualifiedRunfilesPath(path);
  return runfiles->Rlocation(full_path);
}

std::string PathResolver::ResolveRunfilesPathIfRelative(
    absl::string_view path) {
  if (absl::StartsWith(path, "/")) return std::string(path);
  return ResolveRunfilesPath(path);
}

}  // namespace intrinsic
