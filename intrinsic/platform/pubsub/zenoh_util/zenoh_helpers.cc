// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/platform/pubsub/zenoh_util/zenoh_helpers.h"

#include <cstdlib>
#include <string>
#include <vector>

#include "absl/status/status.h"
#include "absl/strings/match.h"
#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"
#include "ortools/base/path.h"
#include "rules_cc/cc/runfiles/runfiles.h"

namespace intrinsic {

using ::rules_cc::cc::runfiles::Runfiles;

bool RunningUnderTest() {
  return (getenv("TEST_TMPDIR") != nullptr) ||
         (getenv("PORTSERVER_ADDRESS") != nullptr);
}

bool RunningInKubernetes() {
  return getenv("KUBERNETES_SERVICE_HOST") != nullptr;
}

std::string GetZenohRunfilesPath(absl::string_view file_path) {
  std::string error;
  std::string path = std::string(file_path);
  std::unique_ptr<Runfiles> runfiles;
  std::string repository = BAZEL_CURRENT_REPOSITORY;
  if (repository == "") {
    // If empty, running in current repository.
    repository = "_main";
  }

  if (RunningUnderTest()) {
    std::string error;
    runfiles = std::unique_ptr<Runfiles>(Runfiles::CreateForTest(&error));
    path = runfiles->Rlocation(file::JoinPath(repository, file_path));
  } else if (!RunningInKubernetes()) {
    runfiles = std::unique_ptr<Runfiles>(
        Runfiles::Create(program_invocation_name, &error));
    path = runfiles->Rlocation(file::JoinPath(repository, file_path));
  }
  return path;
}

absl::Status ValidZenohKeyexpr(absl::string_view keyexpr) {
  if (keyexpr.empty()) {
    return absl::InvalidArgumentError("Keyexpr must not be empty");
  }
  if (absl::StartsWith(keyexpr, "/")) {
    return absl::InvalidArgumentError("Keyexpr must not start with /");
  }
  if (absl::EndsWith(keyexpr, "/")) {
    return absl::InvalidArgumentError("Keyexpr must not end with /");
  }
  std::vector<std::string> parts = absl::StrSplit(keyexpr, '/');
  for (absl::string_view part : parts) {
    if (part.empty()) {
      return absl::InvalidArgumentError("Keyexpr must not contain empty parts");
    }
    if (part == "*" || part == "$*" || part == "**") {
      continue;
    }
  }
  return absl::OkStatus();
}

absl::Status ValidZenohKey(absl::string_view key) {
  if (key.empty()) {
    return absl::InvalidArgumentError("Keyexpr must not be empty");
  }
  if (absl::StartsWith(key, "/")) {
    return absl::InvalidArgumentError("Keyexpr must not start with /");
  }
  if (absl::EndsWith(key, "/")) {
    return absl::InvalidArgumentError("Keyexpr must not end with /");
  }
  std::vector<std::string> parts = absl::StrSplit(key, '/');
  for (absl::string_view part : parts) {
    if (part.empty()) {
      return absl::InvalidArgumentError("Keyexpr must not contain empty parts");
    }
  }
  return absl::OkStatus();
}

}  // namespace intrinsic
