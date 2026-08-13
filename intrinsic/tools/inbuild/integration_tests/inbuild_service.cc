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

#include <memory>

#include "absl/log/check.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "intrinsic/icon/release/file_helpers.h"
#include "intrinsic/icon/release/portable/init_intrinsic.h"
#include "intrinsic/resources/proto/runtime_context.pb.h"
#include "intrinsic/tools/inbuild/integration_tests/inbuild_service.pb.h"
#include "intrinsic/util/proto/any.h"
#include "intrinsic/util/status/status_macros.h"

absl::Status MainImpl() {
  constexpr absl::string_view kContextFilePath =
      "/etc/intrinsic/runtime_config.pb";
  INTR_ASSIGN_OR_RETURN(
      const auto context,
      intrinsic::GetBinaryProto<intrinsic_proto::config::RuntimeContext>(
          kContextFilePath),
      _ << "Reading runtime context");

  auto config =
      std::make_unique<intrinsic_proto::services::InbuildServiceConfig>();
  INTR_RETURN_IF_ERROR(intrinsic::UnpackAny(context.config(), *config));

  LOG(INFO) << "Hello from C++ InbuildService: " << config->bar();

  // Sleep forever.
  absl::SleepFor(absl::InfiniteDuration());

  // Never reached
  return absl::OkStatus();
}

int main(int argc, char** argv) {
  InitIntrinsic(argv[0], argc, argv);

  QCHECK_OK(MainImpl());
  return 0;
}
