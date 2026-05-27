// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_INTERFACES_CONTROL_PERIOD_UTILS_H_
#define INTRINSIC_ICON_HAL_INTERFACES_CONTROL_PERIOD_UTILS_H_

#include <string>

#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "flatbuffers/detached_buffer.h"
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/utils/duration.h"

namespace intrinsic_fbs {

struct ControlPeriod;

// Creates the `ControlPeriod` flatbuffer initialized to an invalid value.
// Use `UpdateControlPeriod` to set a value.
flatbuffers::DetachedBuffer BuildControlPeriod();

}  // namespace intrinsic_fbs

namespace intrinsic::icon {

// Updates the control period in the given ControlPeriod hardware interface
// handle.
// Returns FailedPrecondition when the Duration is negative, or Zero.
absl::Status UpdateControlPeriod(
    MutableHardwareInterfaceHandle<intrinsic_fbs::ControlPeriod>& handle,
    intrinsic::Duration duration);

}  // namespace intrinsic::icon

namespace intrinsic_fbs {

// Returns a canonical user friendly error message containing the periods in ns,
// as well is the respective frequency in Hertz.
std::string FormatControlPeriodMismatchError(absl::string_view module_name,
                                             intrinsic::Duration expected,
                                             intrinsic::Duration actual);

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_HAL_INTERFACES_CONTROL_PERIOD_UTILS_H_
