// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/interfaces/control_period_utils.h"

#include <cmath>

#include "absl/strings/str_format.h"
#include "flatbuffers/flatbuffer_builder.h"
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/interfaces/control_period.fbs.h"
#include "intrinsic/icon/utils/clock.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer BuildControlPeriod() {
  flatbuffers::FlatBufferBuilder builder;
  builder.Finish(builder.CreateStruct(ControlPeriod(0)));
  return builder.Release();
}

}  // namespace intrinsic_fbs

namespace {}

namespace intrinsic::icon {

absl::Status UpdateControlPeriod(
    MutableHardwareInterfaceHandle<intrinsic_fbs::ControlPeriod>& handle,
    intrinsic::Duration duration) {
  if (duration <= intrinsic::ZeroDuration()) {
    return absl::FailedPreconditionError(
        absl::StrCat("Control period must be > 0, got ",
                     intrinsic::ToInt64Nanoseconds(duration),
                     " ns. Check your configuration."));
  }
  handle->mutate_control_period_ns(intrinsic::ToInt64Nanoseconds(duration));
  handle.UpdatedAt(intrinsic::Clock::Now());
  return absl::OkStatus();
}

}  // namespace intrinsic::icon

namespace intrinsic_fbs {

std::string FormatControlPeriodMismatchError(absl::string_view module_name,
                                             intrinsic::Duration expected,
                                             intrinsic::Duration actual) {
  double expected_hz = NAN;
  if (expected > intrinsic::ZeroDuration()) {
    expected_hz = intrinsic::toHertz<double>(expected);
  }

  double actual_hz = NAN;
  if (actual > intrinsic::ZeroDuration()) {
    actual_hz = intrinsic::toHertz<double>(actual);
  }

  return absl::StrFormat(
      "Inconsistent configuration with Hardware Module '%s'."
      " ICON ('control_frequency_hz'): %d ns (%.1f Hz), Hardware Module "
      "reports: %d ns (%.1f Hz). Check your configuration.",
      module_name, intrinsic::ToInt64Nanoseconds(expected), expected_hz,
      intrinsic::ToInt64Nanoseconds(actual), actual_hz);
}

}  // namespace intrinsic_fbs
