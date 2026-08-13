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
