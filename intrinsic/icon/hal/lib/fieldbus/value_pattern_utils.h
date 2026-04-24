// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_FIELDBUS_VALUE_PATTERN_UTILS_H_
#define INTRINSIC_ICON_HAL_LIB_FIELDBUS_VALUE_PATTERN_UTILS_H_

#include <cstdint>

#include "absl/status/status.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"

namespace intrinsic::fieldbus {

// Returns InvalidArgumentError if the signal_mask or expected_value has a bit
// set outside of the valid range of the digital_input_variable.
// Returns OkStatus when the variable exists and the mask and value are within
// the size of the variable.
absl::Status ValidateValuePattern(
    uint32_t signal_mask, uint32_t expected_value,
    const fieldbus::ProcessVariable& digital_input_variable);

// Returns InvalidArgumentError if the signal_mask or expected_value has a bit
// set outside of the valid range of the `num_input_bits`.
// Returns OkStatus when the variable exists and the mask and value are
// within range of `num_input_bits`.
// `num_input_bits` is the size of the input in bits.
absl::Status ValidateValuePattern(uint32_t signal_mask, uint32_t expected_value,
                                  size_t num_input_bits);

}  // namespace intrinsic::fieldbus

#endif  // INTRINSIC_ICON_HAL_LIB_FIELDBUS_VALUE_PATTERN_UTILS_H_
