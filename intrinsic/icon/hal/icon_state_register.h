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

#ifndef INTRINSIC_ICON_HAL_ICON_STATE_REGISTER_H_
#define INTRINSIC_ICON_HAL_ICON_STATE_REGISTER_H_

#include "absl/strings/string_view.h"
#include "intrinsic/icon/hal/hardware_interface_traits.h"
#include "intrinsic/icon/hal/interfaces/icon_state.fbs.h"
#include "intrinsic/icon/hal/interfaces/icon_state_utils.h"

namespace intrinsic::icon {

// Reserved name of the ICON state interface.
static constexpr absl::string_view kIconStateInterfaceName = "icon_state";

namespace hardware_interface_traits {

// Registers the IconState hardware interface.
// Allows transparently depending on IconState and can be included in multiple
// files.
//
// Usage:
// #include "intrinsic/icon/hal/icon_state_register.h"  // IWYU pragma: keep
INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::IconState,
                                 intrinsic_fbs::BuildIconState,
                                 "intrinsic_fbs.IconState")
}  // namespace hardware_interface_traits
}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_HAL_ICON_STATE_REGISTER_H_
