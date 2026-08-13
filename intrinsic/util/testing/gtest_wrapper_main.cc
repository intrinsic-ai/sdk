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

// A custom main for running unit tests and microbenchmarks.

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <cstdio>

#include "absl/debugging/failure_signal_handler.h"
#include "absl/debugging/symbolize.h"
#include "absl/flags/parse.h"
#include "benchmark/benchmark.h"
#include "fuzztest/init_fuzztest.h"

int main(int argc, char** argv) {
  printf("Running main() from %s\n", __FILE__);

  // Run benchmarks if there are any and --benchmark_filter is provided.
  benchmark::Initialize(&argc, argv);
  if (!benchmark::GetBenchmarkFilter().empty()) {
    benchmark::RunSpecifiedBenchmarks();
    benchmark::Shutdown();
    return 0;
  }

  // Before usage, the fuzzy tests need to be initialized.
  fuzztest::ParseAbslFlags(argc, argv);
  fuzztest::InitFuzzTest(&argc, &argv);

  // Run unit tests if there are any. Since Google Mock depends on Google
  // Test, InitGoogleMock() is also responsible for initializing Google Test.
  // Thus, there is no need to call InitGoogleTest() separately. Note that
  // not calling InitGoogleTest() before calling RUN_ALL_TESTS() is INVALID.
  // Valid usage will be enforced in the future.
  testing::InitGoogleMock(&argc, argv);

  // absl::ParseCommandLine to return error on invalid flags.
  absl::ParseCommandLine(argc, argv);

  // Provide stack traces on SIGSEGV and other signals.
  absl::InitializeSymbolizer(argv[0]);
  absl::FailureSignalHandlerOptions options;
  options.call_previous_handler = true;
  absl::InstallFailureSignalHandler(options);

  return RUN_ALL_TESTS();
}
