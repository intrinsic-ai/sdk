// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/fieldbus/fake_variable_registry.h"

#include <cstdint>
#include <tuple>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/strings/substitute.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/service_variable.h"

namespace intrinsic::fieldbus {

absl::StatusOr<ProcessVariable> FakeVariableRegistry::GetInputVariable(
    absl::string_view variable_name) const {
  auto variable = input_bus_variables_.find(variable_name);
  if (variable == input_bus_variables_.end()) {
    return absl::NotFoundError(absl::StrCat(
        "Failed to find input variable with the name: ", variable_name));
  }
  if (variable->second.size() != 1) {
    return absl::NotFoundError(absl::StrCat("Input variable ", variable_name,
                                            " is not a scalar variable."));
  }
  return variable->second.front();
}

absl::StatusOr<ProcessVariable>
FakeVariableRegistry::GetInputArrayFieldVariable(
    absl::string_view variable_name, uint64_t array_index) const {
  auto variable = input_bus_variables_.find(variable_name);
  if (variable == input_bus_variables_.end()) {
    return absl::NotFoundError(absl::StrCat(
        "Failed to find input variable with the name: ", variable_name));
  }
  if (array_index >= variable->second.size()) {
    return absl::OutOfRangeError(absl::StrCat(
        "Array index ", array_index, " is out of bounds for variable ",
        variable_name, " with size ", variable->second.size()));
  }
  return variable->second[array_index];
}

absl::StatusOr<ProcessVariable> FakeVariableRegistry::GetOutputVariable(
    absl::string_view variable_name) const {
  auto variable = output_bus_variables_.find(variable_name);
  if (variable == output_bus_variables_.end()) {
    return absl::NotFoundError(absl::StrCat(
        "Failed to find output variable with the name: ", variable_name));
  }
  if (variable->second.size() != 1) {
    return absl::NotFoundError(absl::StrCat("Output variable ", variable_name,
                                            " is not a scalar variable."));
  }
  return variable->second.front();
}

absl::StatusOr<ProcessVariable>
FakeVariableRegistry::GetOutputArrayFieldVariable(
    absl::string_view variable_name, uint64_t array_index) const {
  auto variable = output_bus_variables_.find(variable_name);
  if (variable == output_bus_variables_.end()) {
    return absl::NotFoundError(absl::StrCat(
        "Failed to find output variable with the name: ", variable_name));
  }
  if (array_index >= variable->second.size()) {
    return absl::OutOfRangeError(absl::StrCat(
        "Array index ", array_index, " is out of bounds for variable ",
        variable_name, " with size ", variable->second.size()));
  }
  return variable->second[array_index];
}

absl::StatusOr<ServiceVariable> FakeVariableRegistry::GetServiceVariable(
    uint32_t index, uint32_t subindex, int32_t bus_position) const {
  auto key = std::make_tuple(index, subindex, bus_position);
  if (!service_variables_.contains(key)) {
    return absl::NotFoundError(absl::Substitute(
        "Failed to find service variable with [$0.$1] at bus position: $2",
        index, subindex, bus_position));
  }
  return service_variables_.at(key);
}

}  // namespace intrinsic::fieldbus
