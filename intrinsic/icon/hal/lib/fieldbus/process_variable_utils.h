// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_FIELDBUS_PROCESS_VARIABLE_UTILS_H_
#define INTRINSIC_ICON_HAL_LIB_FIELDBUS_PROCESS_VARIABLE_UTILS_H_

#include <limits>
#include <type_traits>

#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/utils/realtime_status_or.h"

namespace intrinsic::fieldbus {

template <typename From, typename To>
bool IsOutOfRange(const From& value) {
  // Handle potential loss of information (e.g., denormalized numbers)
  if constexpr (std::is_floating_point<From>::value &&
                std::is_floating_point<To>::value) {
    To temp = static_cast<To>(value);
    if (std::isnan(temp) || std::isinf(temp)) {
      return true;
    }
  }

  // Sign check based on the actual value
  if (value < 0 && std::is_unsigned<To>::value) {
    return true;
  }

  // Check if To type has a wider range than From type
  if constexpr (std::is_arithmetic<From>::value &&
                std::is_arithmetic<To>::value) {
    if (std::numeric_limits<To>::max() >= std::numeric_limits<From>::max() &&
        std::numeric_limits<To>::lowest() <=
            std::numeric_limits<From>::lowest()) {
      return false;  // To type has a wider or equal range, always fits
    }
  }

  // Check if value is within the target type's range
  if (value > static_cast<From>(std::numeric_limits<To>::max()) ||
      value < static_cast<From>(std::numeric_limits<To>::lowest())) {
    return true;
  }

  // Handle potential loss of information for other type combinations
  To temp = static_cast<To>(value);
  From tempBack = static_cast<From>(temp);
  if (!std::is_floating_point<To>::value && tempBack != value) {
    return true;
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
