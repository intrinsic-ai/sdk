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

#ifndef INTRINSIC_ICON_HAL_INTERFACES_HARDWARE_MODULE_STATE_UTILS_H_
#define INTRINSIC_ICON_HAL_INTERFACES_HARDWARE_MODULE_STATE_UTILS_H_

#include <string_view>

#include "flatbuffers/detached_buffer.h"
#include "intrinsic/icon/hal/interfaces/hardware_module_state.fbs.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer BuildHardwareModuleState();

void SetState(HardwareModuleState* hardware_module_state, StateCode code,
              std::string_view message);

// Returns the message associated with the given state.
//
// Returns an empty string if `hardware_module_state` or
// `hardware_module_state->message()` is null.
std::string_view GetMessage(const HardwareModuleState* hardware_module_state);

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_HAL_INTERFACES_HARDWARE_MODULE_STATE_UTILS_H_
