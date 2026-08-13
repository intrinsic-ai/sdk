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

#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"

#include <cmath>
#include <cstddef>
#include <cstdint>

namespace intrinsic::fieldbus {

ProcessVariable::ProcessVariable(uint8_t* data, Type type, std::size_t bit_size,
                                 uint8_t bit_offset)
    : data_(data),
      type_(type),
      bit_size_(bit_size),
      byte_size_(std::ceil(bit_size / 8.0)),
      bit_offset_(bit_offset) {}

}  // namespace intrinsic::fieldbus
