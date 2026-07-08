// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CONTROL_LOGGING_MODE_H_
#define INTRINSIC_ICON_CONTROL_LOGGING_MODE_H_

#include <cstdint>

#include "absl/status/statusor.h"
#include "absl/time/time.h"
#include "intrinsic/icon/proto/logging_mode.pb.h"

namespace intrinsic::icon {

enum class LoggingMode : char {
  // Log at a throttled rate.
  kThrottled,
  // Log every cycle.
  kFullRate,
};

// Regardless of the actual controller frequency, we will publish a throttled
// version on a different topic that is always 50 Hz to simplify its use in
// subscribers that do not need or want the full-rate (often >= 500 Hz) stream
// of status messages. Details in go/intrinsic-robot-status-rate
constexpr double kThrottledStatusRate = 50.0;

LoggingMode FromProto(const intrinsic_proto::icon::LoggingMode& proto);

intrinsic_proto::icon::LoggingMode ToProto(LoggingMode mode);

// Calculates the factor with which to downsample cloud logs in ICON.
absl::StatusOr<int> CalculateDecimationFactor(double control_frequency_hz,
                                              double target_log_rate_hz);

// Calulates whether the current `cycle` should be published while throttling
// the logging rate.
bool IsThrottledLogCycle(uint64_t cycle, int decimation_factor);

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_CONTROL_LOGGING_MODE_H_
