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

#include "intrinsic/icon/hal/lib/fieldbus/service_variable.h"

#include <cstddef>
#include <cstdint>
#include <functional>

#include "absl/status/status.h"

namespace intrinsic::fieldbus {

ServiceVariable::ServiceVariable(
    std::function<absl::Status(uint8_t*, std::size_t)> service_variable_read,
    std::function<absl::Status(const uint8_t*, std::size_t)>
        service_variable_write)
    : service_variable_read_(service_variable_read),
      service_variable_write_(service_variable_write) {}

}  // namespace intrinsic::fieldbus
