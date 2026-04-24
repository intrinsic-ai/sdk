// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_FIELDBUS_FAKE_VARIABLE_REGISTRY_H_
#define INTRINSIC_ICON_HAL_LIB_FIELDBUS_FAKE_VARIABLE_REGISTRY_H_

#include <cstddef>
#include <cstdint>
#include <string>
#include <tuple>
#include <type_traits>
#include <utility>
#include <vector>

#include "absl/container/flat_hash_map.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/service_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"

namespace intrinsic::fieldbus {

// A fake fieldbus variable registry for the testing of bus device
// implementations.
// Bus devices access the variables on the fieldbus bus via instances of
// `intrinsic::fieldbus::ProcessVariable` or
// `intrinsic::fieldbus::ServiceVariable` obtained from an VariableRegistry
// interface. This FakeVariableRegistry  implements the same interface, doesn't
// require a physical bus, and provides functions to register pointers to
// existing variables that may then exposed via the registry interface.
class FakeVariableRegistry : public VariableRegistry {
 public:
  // Registers `value` to be available as an input ProcessVariable under the
  // provided `variable_name`.
  template <typename T>
  void AddInputProcessVariable(absl::string_view variable_name, T* value) {
    input_bus_variables_.emplace(
        variable_name,
        std::vector<ProcessVariable>{CreateProcessVariable(value)});
  }

  // Registers the values from `begin` to `end` to be available as an input
  // ProcessVariable array under the provided `variable_name`.
  template <typename Iterator>
  void AddInputProcessVariableArray(absl::string_view variable_name,
                                    Iterator begin, Iterator end) {
    std::vector<ProcessVariable> variables;
    for (Iterator it = begin; it != end; ++it) {
      variables.emplace_back(CreateProcessVariable(&(*it)));
    }
    input_bus_variables_.emplace(variable_name, std::move(variables));
  }

  // Registers `value` to be available as an output ProcessVariable under the
  // provided `variable_name`.
  template <typename T>
  void AddOutputProcessVariable(absl::string_view variable_name, T* value) {
    output_bus_variables_.emplace(
        variable_name,
        std::vector<ProcessVariable>{CreateProcessVariable(value)});
  }

  // Registers the values from `begin` to `end` to be available as an output
  // ProcessVariable array under the provided `variable_name`.
  template <typename Iterator>
  void AddOutputProcessVariableArray(absl::string_view variable_name,
                                     Iterator begin, Iterator end) {
    std::vector<ProcessVariable> variables;
    for (Iterator it = begin; it != end; ++it) {
      variables.emplace_back(CreateProcessVariable(&(*it)));
    }
    output_bus_variables_.emplace(variable_name, std::move(variables));
  }
  // Registers `value` to be available as a service variable at
  // `index`.`subindex` at the device at `bus_position`.
  template <typename T>
  void AddServiceVariable(uint32_t index, uint32_t subindex,
                          int32_t bus_position, T* value) {
    static_assert(std::is_arithmetic_v<T>);

    auto service_variable_read = [value](uint8_t* data,
                                         std::size_t size) -> absl::Status {
      if (size != sizeof(T)) {
        return absl::InvalidArgumentError(
            absl::StrCat("Size mismatch in service_variable_read. Expected ",
                         sizeof(T), ", but got ", size, "."));
      }
      memcpy(data, value, sizeof(T));
      return absl::OkStatus();
    };

    auto service_variable_write = [value](const uint8_t* data,
                                          std::size_t size) -> absl::Status {
      if (size != sizeof(T)) {
        return absl::InvalidArgumentError(
            absl::StrCat("Size mismatch in service_variable_write. Expected ",
                         sizeof(T), ", but got ", size, "."));
      }
      memcpy(value, data, sizeof(T));
      return absl::OkStatus();
    };

    service_variables_.emplace(
        std::make_tuple(index, subindex, bus_position),
        ServiceVariable(service_variable_read, service_variable_write));
  }

  absl::StatusOr<ProcessVariable> GetInputVariable(
      absl::string_view variable_name) const override;

  absl::StatusOr<ProcessVariable> GetInputArrayFieldVariable(
      absl::string_view variable_name, uint64_t array_index) const override;

  absl::StatusOr<ProcessVariable> GetOutputVariable(
      absl::string_view variable_name) const override;

  absl::StatusOr<ProcessVariable> GetOutputArrayFieldVariable(
      absl::string_view variable_name, uint64_t array_index) const override;

  absl::StatusOr<ServiceVariable> GetServiceVariable(
      uint32_t index, uint32_t subindex, int32_t bus_position) const override;

 private:
  template <typename T>
  ProcessVariable CreateProcessVariable(T* value) {
    static_assert(std::is_arithmetic_v<T>);
    ProcessVariable::Type type = ProcessVariable::kUnknown;
    if constexpr (std::is_floating_point_v<T>) {
      if constexpr (sizeof(T) == 4) {
        type = ProcessVariable::kFloat;
      } else if constexpr (sizeof(T) == 8) {
        type = ProcessVariable::kDouble;
      }
    } else if constexpr (std::is_unsigned_v<T>) {
      type = ProcessVariable::kUnsignedIntegral;
    } else if constexpr (std::is_signed_v<T>) {
      type = ProcessVariable::kSignedIntegral;
    }

    constexpr std::size_t kBitSize =
        std::is_same_v<T, bool> ? 1 : sizeof(T) * 8;
    return ProcessVariable(reinterpret_cast<uint8_t*>(value), type, kBitSize);
  }

  absl::flat_hash_map<std::string, std::vector<ProcessVariable>>
      input_bus_variables_;
  absl::flat_hash_map<std::string, std::vector<ProcessVariable>>
      output_bus_variables_;
  absl::flat_hash_map<std::tuple<uint32_t, uint32_t, int32_t>, ServiceVariable>
      service_variables_;
};

}  // namespace intrinsic::fieldbus

#endif  // INTRINSIC_ICON_HAL_LIB_FIELDBUS_FAKE_VARIABLE_REGISTRY_H_
