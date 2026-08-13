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

#include "intrinsic/icon/release/portable/init_intrinsic.h"

#include <cstdint>
#include <cstdlib>

#include "absl/base/log_severity.h"
#include "absl/flags/flag.h"
#include "absl/log/globals.h"
#include "absl/strings/string_view.h"

ABSL_FLAG(int64_t, int_flag, 0, "integer value for testing");

// Note: We can't use gtest as it's incompatible with the log init in
// InitIntrinsic().
int main() {
  // Check the command line flags are parsed properly.
  constexpr const char* argv[] = {"init_xfa_test", "--int_flag=10"};
  constexpr int argc = sizeof(argv) / sizeof(argv[0]);
  InitIntrinsic(nullptr, argc, const_cast<char**>(argv));
  if (absl::GetFlag(FLAGS_int_flag) != 10) {
    return EXIT_FAILURE;
  }
  // Check that LOG(INFO) statements are sent to stderr.
  if (absl::StderrThreshold() != absl::LogSeverityAtLeast::kInfo) {
    return EXIT_FAILURE;
  }
  return EXIT_SUCCESS;
}
