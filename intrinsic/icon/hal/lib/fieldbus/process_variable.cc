// Copyright 2023 Intrinsic Innovation LLC

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
