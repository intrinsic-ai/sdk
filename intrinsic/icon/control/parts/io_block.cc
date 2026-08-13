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

#include "intrinsic/icon/control/parts/io_block.h"

#include <stddef.h>

#include "absl/types/span.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_or.h"

namespace intrinsic::icon {

RealtimeStatusOr<DioBlock> DioBlock::Create(size_t size) {
  if (size > kMaxValuesPerBlock) {
    return InvalidArgumentError(icon::RealtimeStatus::StrCat(
        "The parameter size(", size, ") exceeds kMaxValuesPerBlock(",
        kMaxValuesPerBlock, ")"));
  }
  return DioBlock(size);
}

RealtimeStatusOr<AnalogBlock> AnalogBlock::Create(
    absl::Span<const Unit> units) {
  if (units.size() > kMaxValuesPerBlock) {
    return InvalidArgumentError(icon::RealtimeStatus::StrCat(
        "The size of units(", units.size(), ") exceeds kMaxValuesPerBlock(",
        kMaxValuesPerBlock, ")"));
  }
  return AnalogBlock(units);
}

RealtimeStatusOr<AnalogBlock> AnalogBlock::Create(size_t size) {
  if (size > kMaxValuesPerBlock) {
    return InvalidArgumentError(icon::RealtimeStatus::StrCat(
        "The parameter size(", size, ") exceeds kMaxValuesPerBlock(",
        kMaxValuesPerBlock, ")"));
  }
  return AnalogBlock(size);
}

}  // namespace intrinsic::icon
