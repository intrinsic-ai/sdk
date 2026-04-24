// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/process_variable_value/process_variable_device.h"

#include <cstdint>
#include <functional>
#include <limits>
#include <memory>
#include <utility>
#include <vector>

#include "absl/memory/memory.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "intrinsic/icon/hal/default_hardware_interfaces.h"  // IWYU pragma: keep
#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"
#include "intrinsic/icon/hal/lib/process_variable_value/v1/process_variable_config.pb.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::fieldbus {

ProcessVariableDevice::ProcessVariableDevice(
    std::function<void()> constants_to_variable_when_enabling,
    std::function<void()> constants_to_variable_when_operational,
    std::function<void()> constants_to_variable_when_disabling)
    : constants_to_variable_when_enabling_(
          std::move(constants_to_variable_when_enabling)),
      constants_to_variable_when_operational_(
          std::move(constants_to_variable_when_operational)),
      constants_to_variable_when_disabling_(
          std::move(constants_to_variable_when_disabling)) {}

template <typename SourceType, typename TargetType>
absl::Status RangeCheck(SourceType value, TargetType min_value,
                        TargetType max_value) {
  if (value < min_value) {
    return absl::OutOfRangeError(absl::StrCat("Value is ", value,
                                              ", but must not be less than ",
                                              static_cast<double>(min_value)));
  }
  if (value > max_value) {
    return absl::OutOfRangeError(absl::StrCat("Value is ", value,
                                              ", but must not be greater than ",
                                              static_cast<double>(max_value)));
  }
  return absl::OkStatus();
}

template <typename SourceType, typename TargetType>
absl::Status CheckCompatibility(
    const fieldbus::ProcessVariable& process_variable,
    SourceType constant_value) {
  INTR_RETURN_IF_ERROR(process_variable.IsCompatibleType<TargetType>());
  INTR_RETURN_IF_ERROR(RangeCheck(constant_value,
                                  std::numeric_limits<TargetType>::min(),
                                  std::numeric_limits<TargetType>::max()));
  return absl::OkStatus();
}

absl::StatusOr<std::function<void()>> GetVariableWriter(
    const intrinsic_proto::fieldbus::v1::ProcessVariableValueType&
        constant_value,
    const fieldbus::ProcessVariable& process_variable) {
  if (constant_value.has_bool_value()) {
    INTR_RETURN_IF_ERROR(process_variable.IsCompatibleType<bool>());
    return [process_variable = process_variable,
            value = constant_value.bool_value()]() mutable {
      process_variable.WriteUnchecked(value);
    };
  } else if (constant_value.has_uint8_value()) {
    INTR_RETURN_IF_ERROR(
        (CheckCompatibility<decltype(constant_value.uint8_value()), uint8_t>(
            process_variable, constant_value.uint8_value())));
    return
        [process_variable = process_variable,
         value = static_cast<uint8_t>(constant_value.uint8_value())]() mutable {
          process_variable.WriteUnchecked(value);
        };
  } else if (constant_value.has_int8_value()) {
    INTR_RETURN_IF_ERROR(
        (CheckCompatibility<decltype(constant_value.int8_value()), int8_t>(
            process_variable, constant_value.int8_value())));
    return
        [process_variable = process_variable,
         value = static_cast<int8_t>(constant_value.int8_value())]() mutable {
          process_variable.WriteUnchecked(value);
        };
  } else if (constant_value.has_uint16_value()) {
    INTR_RETURN_IF_ERROR(
        (CheckCompatibility<decltype(constant_value.uint16_value()), uint16_t>(
            process_variable, constant_value.uint16_value())));
    return [process_variable = process_variable,
            value = static_cast<uint16_t>(
                constant_value.uint16_value())]() mutable {
      process_variable.WriteUnchecked(value);
    };
  } else if (constant_value.has_int16_value()) {
    INTR_RETURN_IF_ERROR(
        (CheckCompatibility<decltype(constant_value.int16_value()), int16_t>(
            process_variable, constant_value.int16_value())));
    return
        [process_variable = process_variable,
         value = static_cast<int16_t>(constant_value.int16_value())]() mutable {
          process_variable.WriteUnchecked(value);
        };
  } else if (constant_value.has_uint32_value()) {
    INTR_RETURN_IF_ERROR(process_variable.IsCompatibleType<uint32_t>());
    return [process_variable = process_variable,
            value = static_cast<uint32_t>(
                constant_value.uint32_value())]() mutable {
      process_variable.WriteUnchecked(value);
    };
  } else if (constant_value.has_int32_value()) {
    INTR_RETURN_IF_ERROR(process_variable.IsCompatibleType<int32_t>());
    return
        [process_variable = process_variable,
         value = static_cast<int32_t>(constant_value.int32_value())]() mutable {
          process_variable.WriteUnchecked(value);
        };
  } else if (constant_value.has_uint64_value()) {
    INTR_RETURN_IF_ERROR(process_variable.IsCompatibleType<uint64_t>());
    return [process_variable = process_variable,
            value = static_cast<uint64_t>(
                constant_value.uint64_value())]() mutable {
      process_variable.WriteUnchecked(value);
    };
  } else if (constant_value.has_int64_value()) {
    INTR_RETURN_IF_ERROR(process_variable.IsCompatibleType<int64_t>());
    return
        [process_variable = process_variable,
         value = static_cast<int64_t>(constant_value.int64_value())]() mutable {
          process_variable.WriteUnchecked(value);
        };
  } else if (constant_value.has_float_value()) {
    INTR_RETURN_IF_ERROR(process_variable.IsCompatibleType<float>());
    return
        [process_variable = process_variable,
         value = static_cast<float>(constant_value.float_value())]() mutable {
          process_variable.WriteUnchecked(value);
        };
  } else if (constant_value.has_double_value()) {
    INTR_RETURN_IF_ERROR(process_variable.IsCompatibleType<double>());
    return
        [process_variable = process_variable,
         value = static_cast<double>(constant_value.double_value())]() mutable {
          process_variable.WriteUnchecked(value);
        };
  }
  return absl::InvalidArgumentError(
      absl::StrCat("Unsupported value type: ", constant_value));
}

absl::StatusOr<std::unique_ptr<ProcessVariableDevice>>
ProcessVariableDevice::Create(
    fieldbus::DeviceInitContext& init_context,
    const intrinsic_proto::fieldbus::v1::ProcessVariableConfig& config) {
  const fieldbus::VariableRegistry& variable_registry =
      init_context.GetVariableRegistry();

  std::vector<std::function<void()>> operational_variable_writers;
  std::vector<std::function<void()>> enabling_variable_writers;
  std::vector<std::function<void()>> disabling_variable_writers;
  for (const auto& constant_value : config.pdo_writes()) {
    INTR_ASSIGN_OR_RETURN(
        const auto process_variable,
        variable_registry.GetOutputVariable(constant_value.variable_name()),
        _ << " for variable " << constant_value.variable_name());
    INTR_ASSIGN_OR_RETURN(
        auto variable_writer,
        GetVariableWriter(constant_value.value(), process_variable),
        _ << " for variable " << constant_value.variable_name());
    operational_variable_writers.push_back(variable_writer);

    if (constant_value.has_enabling_override_value()) {
      INTR_ASSIGN_OR_RETURN(
          auto enabling_variable_writer,
          GetVariableWriter(constant_value.enabling_override_value(),
                            process_variable),
          _ << " for variable " << constant_value.variable_name());
      enabling_variable_writers.push_back(std::move(enabling_variable_writer));
    } else {
      enabling_variable_writers.push_back(variable_writer);
    }

    if (constant_value.has_disabling_override_value()) {
      INTR_ASSIGN_OR_RETURN(
          auto disabling_variable_writer,
          GetVariableWriter(constant_value.disabling_override_value(),
                            process_variable),
          _ << " for variable " << constant_value.variable_name());
      disabling_variable_writers.push_back(
          std::move(disabling_variable_writer));
    } else {
      disabling_variable_writers.push_back(std::move(variable_writer));
    }
  }

  std::function<void()> constants_to_variable_when_enabling =
      [enabling_variable_writers = std::move(enabling_variable_writers)]() {
        for (const auto& variable_writer : enabling_variable_writers) {
          variable_writer();
        }
      };
  std::function<void()> constants_to_variable_when_operational =
      [operational_variable_writers =
           std::move(operational_variable_writers)]() {
        for (const auto& variable_writer : operational_variable_writers) {
          variable_writer();
        }
      };
  std::function<void()> constants_to_variable_when_disabling =
      [disabling_variable_writers = std::move(disabling_variable_writers)]() {
        for (const auto& variable_writer : disabling_variable_writers) {
          variable_writer();
        }
      };
  return absl::WrapUnique(new ProcessVariableDevice(
      std::move(constants_to_variable_when_enabling),
      std::move(constants_to_variable_when_operational),
      std::move(constants_to_variable_when_disabling)));
}

intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus>
ProcessVariableDevice::CyclicRead(fieldbus::RequestType) {
  return fieldbus::RequestStatus::kDone;
}

intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus>
ProcessVariableDevice::CyclicWrite(fieldbus::RequestType request_type) {
  switch (request_type) {
    case fieldbus::RequestType::kEnableMotion:
      constants_to_variable_when_enabling_();
      break;
    case fieldbus::RequestType::kDisableMotion:
      constants_to_variable_when_disabling_();
      break;
    default:
      constants_to_variable_when_operational_();
      break;
  }
  return fieldbus::RequestStatus::kDone;
}

}  // namespace intrinsic::fieldbus
