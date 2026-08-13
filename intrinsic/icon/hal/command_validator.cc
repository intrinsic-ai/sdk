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

#include "intrinsic/icon/hal/command_validator.h"

#include <utility>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/hardware_interface_registry.h"
#include "intrinsic/icon/hal/icon_state_register.h"
#include "intrinsic/icon/hal/interfaces/icon_state.fbs.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::icon {

// static
absl::StatusOr<Validator> Validator::Create(
    HardwareInterfaceRegistry& interface_registry) {
  INTR_ASSIGN_OR_RETURN(
      auto icon_state,
      interface_registry.GetInterfaceHandle<intrinsic_fbs::IconState>(
          /*interface_name=*/kIconStateInterfaceName));

  return Validator(std::move(icon_state));
}

}  // namespace intrinsic::icon
