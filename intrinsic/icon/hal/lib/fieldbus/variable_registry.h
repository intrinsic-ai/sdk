// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_FIELDBUS_VARIABLE_REGISTRY_H_
#define INTRINSIC_ICON_HAL_LIB_FIELDBUS_VARIABLE_REGISTRY_H_

#include <cstdint>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/service_variable.h"

namespace intrinsic::fieldbus {

// Interface for fieldbus bus variables.
// Provides access to input and output bus variables.
class VariableRegistry {
 public:
  virtual ~VariableRegistry() = default;
  VariableRegistry() = default;
  VariableRegistry(const VariableRegistry&) = default;
  VariableRegistry& operator=(const VariableRegistry&) = default;

  // Returns a `ProcessVariable` instance to a bus variable inside the input
  // process image. Returns an error if no variable can be found under the given
  // `variable_name`.
  virtual absl::StatusOr<ProcessVariable> GetInputVariable(
      absl::string_view variable_name) const = 0;
  // Returns a `ProcessVariable` instance to the field of a bus variable array
  // at index `array_index` inside the input process image. Returns an error if
  // no variable can be found under the given `variable_name` or if the variable
  // is not an array.
  virtual absl::StatusOr<ProcessVariable> GetInputArrayFieldVariable(
      absl::string_view variable_name, uint64_t array_index) const = 0;
  // Returns a `ProcessVariable` instance to a bus variable inside the output
  // process image.
  // Returns an error if no variable can be found under the given
  // `variable_name`.
  virtual absl::StatusOr<ProcessVariable> GetOutputVariable(
      absl::string_view variable_name) const = 0;
  // Returns a `ProcessVariable` instance to the field of a bus variable array
  // at index `array_index` inside the output process image. Returns an error if
  // no variable can be found under the given `variable_name` or if the variable
  // is not an array.
  virtual absl::StatusOr<ProcessVariable> GetOutputArrayFieldVariable(
      absl::string_view variable_name, uint64_t array_index) const = 0;

  // Returns a `ServiceVariable` instance to a service variable.
  // Access to the variable is blocking and thus not realtime safe.
  virtual absl::StatusOr<ServiceVariable> GetServiceVariable(
      uint32_t index, uint32_t subindex, int32_t bus_position) const = 0;
};

}  // namespace intrinsic::fieldbus

#endif  // INTRINSIC_ICON_HAL_LIB_FIELDBUS_VARIABLE_REGISTRY_H_
