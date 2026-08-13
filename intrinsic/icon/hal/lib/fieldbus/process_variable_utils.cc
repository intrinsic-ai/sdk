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

#include "intrinsic/icon/hal/lib/fieldbus/process_variable_utils.h"

#include <cstdint>

#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_or.h"

namespace intrinsic::fieldbus {

using ::intrinsic::icon::RealtimeStatus;
using ::intrinsic::icon::RealtimeStatusOr;

template <>
RealtimeStatusOr<double> ReadAs(const ProcessVariable& process_variable) {
  if (process_variable.IsCompatibleType<uint8_t>().ok()) {
    return static_cast<double>(process_variable.ReadUnchecked<uint8_t>());
  } else if (process_variable.IsCompatibleType<uint16_t>().ok()) {
    return static_cast<double>(process_variable.ReadUnchecked<uint16_t>());
  } else if (process_variable.IsCompatibleType<uint32_t>().ok()) {
    return static_cast<double>(process_variable.ReadUnchecked<uint32_t>());
  } else if (process_variable.IsCompatibleType<uint64_t>().ok()) {
    return static_cast<double>(process_variable.ReadUnchecked<uint64_t>());
  } else if (process_variable.IsCompatibleType<int8_t>().ok()) {
    return static_cast<double>(process_variable.ReadUnchecked<int8_t>());
  } else if (process_variable.IsCompatibleType<int16_t>().ok()) {
    return static_cast<double>(process_variable.ReadUnchecked<int16_t>());
  } else if (process_variable.IsCompatibleType<int32_t>().ok()) {
    return static_cast<double>(process_variable.ReadUnchecked<int32_t>());
  } else if (process_variable.IsCompatibleType<int64_t>().ok()) {
    return static_cast<double>(process_variable.ReadUnchecked<int64_t>());
  } else if (process_variable.IsCompatibleType<double>().ok()) {
    return process_variable.ReadUnchecked<double>();
  } else if (process_variable.IsCompatibleType<float>().ok()) {
    return static_cast<double>(process_variable.ReadUnchecked<float>());
  } else if (process_variable.IsCompatibleType<bool>().ok()) {
    return static_cast<double>(process_variable.ReadUnchecked<bool>());
  } else {
    return intrinsic::icon::InvalidArgumentError(
        "Variable type is not supported.");
  }
}

template <>
RealtimeStatusOr<float> ReadAs(const ProcessVariable& process_variable) {
  if (process_variable.IsCompatibleType<uint8_t>().ok()) {
    return static_cast<float>(process_variable.ReadUnchecked<uint8_t>());
  } else if (process_variable.IsCompatibleType<uint16_t>().ok()) {
    return static_cast<float>(process_variable.ReadUnchecked<uint16_t>());
  } else if (process_variable.IsCompatibleType<uint32_t>().ok()) {
    return static_cast<float>(process_variable.ReadUnchecked<uint32_t>());
  } else if (process_variable.IsCompatibleType<uint64_t>().ok()) {
    return static_cast<float>(process_variable.ReadUnchecked<uint64_t>());
  } else if (process_variable.IsCompatibleType<int8_t>().ok()) {
    return static_cast<float>(process_variable.ReadUnchecked<int8_t>());
  } else if (process_variable.IsCompatibleType<int16_t>().ok()) {
    return static_cast<float>(process_variable.ReadUnchecked<int16_t>());
  } else if (process_variable.IsCompatibleType<int32_t>().ok()) {
    return static_cast<float>(process_variable.ReadUnchecked<int32_t>());
  } else if (process_variable.IsCompatibleType<int64_t>().ok()) {
    return static_cast<float>(process_variable.ReadUnchecked<int64_t>());
  } else if (process_variable.IsCompatibleType<double>().ok()) {
    double value = process_variable.ReadUnchecked<double>();
    if (IsOutOfRange<double, float>(value)) {
      return intrinsic::icon::OutOfRangeError(RealtimeStatus::StrCat(
          "Value '", value, "' is out of range for float."));
    }
    return static_cast<float>(value);
  } else if (process_variable.IsCompatibleType<float>().ok()) {
    return static_cast<float>(process_variable.ReadUnchecked<float>());
  } else if (process_variable.IsCompatibleType<bool>().ok()) {
    return static_cast<float>(process_variable.ReadUnchecked<bool>());
  } else {
    return intrinsic::icon::InvalidArgumentError(
        "Variable type is not supported.");
  }
}

}  // namespace intrinsic::fieldbus
