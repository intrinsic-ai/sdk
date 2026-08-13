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

#include "intrinsic/icon/control/logging_mode.h"

#include <cmath>
#include <cstdint>

#include "absl/status/statusor.h"
#include "intrinsic/icon/proto/logging_mode.pb.h"

namespace intrinsic::icon {

LoggingMode FromProto(const intrinsic_proto::icon::LoggingMode& proto) {
  switch (proto) {
    case intrinsic_proto::icon::LoggingMode::LOGGING_MODE_FULL_RATE:
      return LoggingMode::kFullRate;
    case intrinsic_proto::icon::LoggingMode::LOGGING_MODE_UNSPECIFIED:
    case intrinsic_proto::icon::LoggingMode::LOGGING_MODE_THROTTLED:
    default:
      return LoggingMode::kThrottled;
  }
}

intrinsic_proto::icon::LoggingMode ToProto(LoggingMode logging_mode) {
  switch (logging_mode) {
    case LoggingMode::kFullRate:
      return intrinsic_proto::icon::LoggingMode::LOGGING_MODE_FULL_RATE;
    case LoggingMode::kThrottled:
    default:
      return intrinsic_proto::icon::LoggingMode::LOGGING_MODE_THROTTLED;
  }
}

absl::StatusOr<int> CalculateDecimationFactor(double control_frequency_hz,
                                              double target_log_rate_hz) {
  if (control_frequency_hz <= 0.0 || target_log_rate_hz <= 0.0) {
    return absl::InvalidArgumentError(
        "control_frequency_hz and target_rate_hz must be positive.");
  }
  if (control_frequency_hz <= target_log_rate_hz) {
    return 1;
  }
  return static_cast<int>(
      std::round(control_frequency_hz / target_log_rate_hz));
}

bool IsThrottledLogCycle(uint64_t cycle, int decimation_factor) {
  if (decimation_factor <= 1) {
    return true;
  }
  return (cycle % decimation_factor) == 0;
}

}  // namespace intrinsic::icon
