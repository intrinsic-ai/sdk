// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_INTERFACES_ELECTRICAL_MOTOR_HARDWARE_INTERFACES_H_
#define INTRINSIC_ICON_HAL_INTERFACES_ELECTRICAL_MOTOR_HARDWARE_INTERFACES_H_

#include "intrinsic/icon/hal/hardware_interface_traits.h"
#include "intrinsic/icon/hal/interfaces/electrical_motor.fbs.h"
#include "intrinsic/icon/hal/interfaces/electrical_motor_utils.h"

namespace intrinsic::icon {
namespace hardware_interface_traits {
INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::HomeCommand,
                                 intrinsic_fbs::BuildHomeCommand,
                                 "intrinsic_fbs.HomeCommand");

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::HomingStatus,
                                 intrinsic_fbs::BuildHomingStatus,
                                 "intrinsic_fbs.HomingStatus");
}  // namespace hardware_interface_traits
}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_HAL_INTERFACES_ELECTRICAL_MOTOR_HARDWARE_INTERFACES_H_
