// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_FIELDBUS_PROCESS_VARIABLE_UTILS_H_
#define INTRINSIC_ICON_HAL_LIB_FIELDBUS_PROCESS_VARIABLE_UTILS_H_

#include <cmath>
#include <concepts>
#include <limits>
#include <type_traits>
#include <utility>

#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/utils/realtime_status_or.h"

namespace intrinsic::fieldbus {

// Checks if the given `value` of type `From` is out of range for the target
// type `To`.
//
// This function performs a compile-time dispatched safety check to determine
// if a value can be represented by the target type `To` without overflow.
//
// @param value The value to check.
// @param check_fractional_loss If true (default), additionally checks for
//                              fractional loss when converting from
//                              floating-point to integral types (i.e., requires
//                              exact representation). If false, only checks
//                              if the value falls within the representable
//                              boundaries of the target type.
//
// Returns true if the value is out of range, false otherwise.
template <typename From, typename To>
constexpr bool IsOutOfRange(const From& value,
                            bool check_fractional_loss = true) {
  if constexpr (std::floating_point<From> && std::floating_point<To>) {
    // Floating-point to floating-point conversion.
    // Check for NaN/Infinity in source.
    if (!std::isfinite(value)) {
      return true;
    }
    // Narrowing conversion check (e.g., double to float)
    if constexpr (sizeof(To) < sizeof(From)) {
      // Cast the limits of the narrower type 'To' to the wider type 'From'.
      // This is safe and avoids undefined behavior.
      if (value > static_cast<From>(std::numeric_limits<To>::max()) ||
          value < static_cast<From>(std::numeric_limits<To>::lowest())) {
        return true;
      }
    }
    return false;
  } else if constexpr (std::floating_point<From> && std::integral<To>) {
    // Floating-point to integral conversion.
    // Check for finite values, boundaries (using scalbn to avoid precision loss
    // on limits), and fractional loss for exact representation.
    if (!std::isfinite(value)) {
      return true;
    }
    // max_plus_one is 2^digits, which is exact and safe.
    From max_plus_one =
        std::scalbn(static_cast<From>(1.0), std::numeric_limits<To>::digits);
    From lowest = static_cast<From>(std::numeric_limits<To>::lowest());

    if (value < lowest || value >= max_plus_one) {
      return true;
    }

    // Strict check for fractional loss.
    if (check_fractional_loss) {
      To temp = static_cast<To>(value);
      From tempBack = static_cast<From>(temp);
      if (tempBack != value) {
        return true;
      }
    }
    return false;
  } else if constexpr (std::integral<From> && std::floating_point<To>) {
    // Integral to floating-point conversion.
    // Always fits within the range of standard floating-point types.
    return false;
  } else if constexpr (std::integral<From> && std::integral<To>) {
    // Integral to integral conversion.
    return !std::in_range<To>(value);
  }
  return false;
}

// Reads the value of the bus variable as the given type.
// Template definition only, do not use directly.
template <typename T>
intrinsic::icon::RealtimeStatusOr<T> ReadAs(
    const ProcessVariable& process_variable);

// Template specialization of ReadAs for double.
// The template definition is in the .cc file.
// Returns the value of the bus variable as a double.
// Returns an invalid argument error if the bus variable is not a primitive
// type.
template <>
intrinsic::icon::RealtimeStatusOr<double> ReadAs(
    const ProcessVariable& process_variable);

// Template specialization of ReadAs for float.
// The template definition is in the .cc file.
// Returns the value of the bus variable as a float.
// Returns an invalid argument error if the bus variable is not a primitive type
// or an out of range error if the value is out of range for float.
template <>
intrinsic::icon::RealtimeStatusOr<float> ReadAs(
    const ProcessVariable& process_variable);

// Add other template specializations here, once required.

}  // namespace intrinsic::fieldbus

#endif  // INTRINSIC_ICON_HAL_LIB_FIELDBUS_PROCESS_VARIABLE_UTILS_H_
