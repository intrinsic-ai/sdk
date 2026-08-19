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

#include "intrinsic/platform/pubsub/zenoh_util/zenoh_helpers.h"

#include <cstdlib>
#include <string>
#include <vector>

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/strings/match.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"

namespace intrinsic {

bool RunningUnderTest() {
  return (getenv("TEST_TMPDIR") != nullptr) ||
         (getenv("PORTSERVER_ADDRESS") != nullptr);
}

bool RunningInKubernetes() {
  return getenv("KUBERNETES_SERVICE_HOST") != nullptr;
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
