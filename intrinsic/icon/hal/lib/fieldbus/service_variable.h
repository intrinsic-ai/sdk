// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_FIELDBUS_SERVICE_VARIABLE_H_
#define INTRINSIC_ICON_HAL_LIB_FIELDBUS_SERVICE_VARIABLE_H_

#include <cstddef>
#include <cstdint>
#include <cstring>
#include <functional>
#include <type_traits>
#include <utility>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::fieldbus {

// Provides `Read` and `Write` access to a service variable.
// Decouples the usage of service variables from the actual bus implementation.
class ServiceVariable {
 public:
  // Creates a service variable.
  // `service_variable_read` and `service_variable_write` are bus-specific
  // functions used to access the actual value via the bus.
  ServiceVariable(
      std::function<absl::Status(uint8_t*, std::size_t)> service_variable_read,
      std::function<absl::Status(const uint8_t*, std::size_t)>
          service_variable_write);

  // Read the value of a service variable as type T, where T must be arithmetic.
  // Not realtime safe, since the read operation may be blocking.
  // Might return an error (depending on the injected bus implementation), i.e.
  // when T doesn't match the variable type.
  template <typename T>
  absl::StatusOr<T> Read() {
    static_assert(std::is_arithmetic_v<T> || std::is_enum_v<T>);
    constexpr std::size_t kBufferSize = sizeof(T);
    static T value;
    INTR_RETURN_IF_ERROR(service_variable_read_(
        reinterpret_cast<uint8_t*>(&value), kBufferSize));
    return value;
  }

  // Write to a service variable from a value of type T, where T must be
  // arithmetic. Not realtime safe, since the write operation may be blocking.
  // Might return an error (depending on the injected bus implementation), i.e.
  // when T doesn't match the variable type.
  template <typename T>
  absl::Status Write(const T& value) {
    static_assert(std::is_arithmetic_v<T> || std::is_enum_v<T>);
    constexpr std::size_t kBufferSize = sizeof(T);
    return service_variable_write_(reinterpret_cast<const uint8_t*>(&value),
                                   kBufferSize);
  }

 private:
  std::function<absl::Status(uint8_t*, std::size_t)> service_variable_read_;
  std::function<absl::Status(const uint8_t*, std::size_t)>
      service_variable_write_;
};

// Reads and writes a value to a service variable and returns the current and
// updated value.
// Reads the current value, writes the given value and reads the updated value.
// Returns the current and updated value as a pair.
// Not realtime safe, since the read and write operations may be blocking.
template <typename T>
absl::StatusOr<std::pair<T, T>> ReadWriteRead(ServiceVariable variable,
                                              T value) {
  INTR_ASSIGN_OR_RETURN(auto current_value, variable.Read<T>());
  INTR_RETURN_IF_ERROR(variable.Write(value));
  INTR_ASSIGN_OR_RETURN(auto updated_value, variable.Read<T>());
  return std::make_pair(current_value, updated_value);
}

// Reads and writes a value to a service variable and confirms that the written
// value is read back.
// Reads the current value, writes the given value and reads the updated value.
// Returns the current and updated value as a pair.
// Not realtime safe, since the read and write operations may be blocking.
// Returns an error if the written value is not read back.
template <typename T>
absl::StatusOr<std::pair<T, T>> ReadWriteReadConfirm(ServiceVariable variable,
                                                     T value) {
  INTR_ASSIGN_OR_RETURN(auto current_value, variable.Read<T>());
  INTR_RETURN_IF_ERROR(variable.Write(value));
  INTR_ASSIGN_OR_RETURN(auto updated_value, variable.Read<T>());
  if (value != updated_value) {
    return absl::InternalError(
        "Unexpected value after write. Was: " + std::to_string(current_value) +
        ", want: " + std::to_string(value) +
        " but got: " + std::to_string(updated_value) +
        ". Are you writing the same value via process variable?");
  }
  return std::make_pair(current_value, updated_value);
}

}  // namespace intrinsic::fieldbus

#endif  // INTRINSIC_ICON_HAL_LIB_FIELDBUS_SERVICE_VARIABLE_H_
