// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/fieldbus/value_pattern_utils.h"

#include <cstdint>

#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"

namespace intrinsic::fieldbus {

absl::Status ValidateValuePattern(uint32_t signal_mask, uint32_t expected_value,
                                  size_t num_input_bits) {
  if (signal_mask == 0) {
    // Mask if zero means that no bit of the signal is checked.
    return absl::InvalidArgumentError(
        "Signal mask cannot be zero because that would ignore the signal.");
  }

  if (num_input_bits == 0) {
    return absl::InvalidArgumentError("Input size cannot be zero.");
  }

  if (const size_t max_input_size = (sizeof(signal_mask) * 8);
      num_input_bits > max_input_size) {
    // Error, as only the 32 least significant bits of the input would be
    // covered by the mask.
    return absl::InvalidArgumentError(
        absl::StrCat("Only inputs <= ", max_input_size,
                     " bits are supported. Got ", num_input_bits, " bits."));
  }

  const auto max_bit_index = num_input_bits - 1;

  // Perform shift using 64bit to prevent overflow.
  if (signal_mask >= (1ULL << num_input_bits)) {
    return absl::InvalidArgumentError(absl::StrCat(
        "The 'signal_mask' 0x", absl::Hex(signal_mask),
        " has a bit set outside of the valid range [0,", max_bit_index, "]"));
  }

  // Perform shift using 64bit to prevent overflow.
  if (expected_value >= (1ULL << num_input_bits)) {
    return absl::InvalidArgumentError(absl::StrCat(
        "The 'expected_value' 0x", absl::Hex(expected_value),
        " has a bit set outside of the valid range [0,", max_bit_index, "]"));
  }

  return absl::OkStatus();
}

absl::Status ValidateValuePattern(
    uint32_t signal_mask, uint32_t expected_value,
    const fieldbus::ProcessVariable& digital_input_variable) {
  size_t input_size = digital_input_variable.bit_size();
  // Provide an actionable error message for `ProcessVariable`s.
  if (signal_mask == 0) {
    // Mask of zero means that no bit of the signal is checked.
    return absl::InvalidArgumentError(
        "Signal mask cannot be zero. Instead use `BEHAVIOR_IGNORE_ALWAYS` in "
        "the `unexpected_state_transition_configuration`.");
  }

  return ValidateValuePattern(signal_mask, expected_value, input_size);
}

}  // namespace intrinsic::fieldbus
