// Copyright 2023 Intrinsic Innovation LLC

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
