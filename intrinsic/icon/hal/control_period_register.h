// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_CONTROL_PERIOD_REGISTER_H_
#define INTRINSIC_ICON_HAL_CONTROL_PERIOD_REGISTER_H_

#include "absl/strings/string_view.h"
#include "intrinsic/icon/hal/hardware_interface_traits.h"
#include "intrinsic/icon/hal/interfaces/control_period.fbs.h"
#include "intrinsic/icon/hal/interfaces/control_period_utils.h"

namespace intrinsic::icon {

// Reserved name of the control period interface.
static constexpr absl::string_view kControlPeriodInterfaceName =
    "control_period";

namespace hardware_interface_traits {

// Registers the ControlPeriod hardware interface.
// Allows transparently depending on ControlPeriod and can be included in
// multiple files.
//
// Usage:
// #include "intrinsic/icon/hal/control_period_register.h"  // IWYU pragma: keep
INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::ControlPeriod,
                                 intrinsic_fbs::BuildControlPeriod,
                                 "intrinsic_fbs.ControlPeriod")
}  // namespace hardware_interface_traits
}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_HAL_CONTROL_PERIOD_REGISTER_H_
