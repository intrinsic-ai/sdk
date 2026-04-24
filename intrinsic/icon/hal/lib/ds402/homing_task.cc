// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/ds402/homing_task.h"

#include <sys/types.h>

#include <bitset>
#include <cstdint>
#include <memory>
#include <optional>
#include <string>
#include <type_traits>
#include <utility>

#include "absl/functional/any_invocable.h"
#include "absl/log/log.h"
#include "absl/memory/memory.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "intrinsic/icon/flatbuffers/fixed_string.h"
#include "intrinsic/icon/hal/command_validator.h"
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/hardware_interface_registry.h"
#include "intrinsic/icon/hal/interfaces/electrical_motor.fbs.h"
#include "intrinsic/icon/hal/interfaces/electrical_motor_hardware_interfaces.h"  // IWYU pragma: keep (registers HomeCommand hardware interface)
#include "intrinsic/icon/hal/lib/ds402/ds402_driver.h"
#include "intrinsic/icon/hal/lib/ds402/v1/ds402_bus_component_config.pb.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/service_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"
#include "intrinsic/icon/utils/bitset.h"
#include "intrinsic/icon/utils/log.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::ds402 {

template <typename T>
absl::StatusOr<absl::AnyInvocable<absl::StatusOr<T>(T)>>
GetWriteFunctionFromConfig(
    int32_t bus_position,
    const intrinsic::fieldbus::VariableRegistry& variable_registry,
    const intrinsic_proto::icon::v1::BusVariable& config,
    uint32_t fallback_service_variable_index,
    uint32_t fallback_service_variable_subindex,
    bool ignore_signedness = false) {
  if (config.has_pdo_variable()) {
    // Configuration provides a process_variable variable.
    INTR_ASSIGN_OR_RETURN(auto process_variable,
                          variable_registry.GetOutputVariable(
                              config.pdo_variable().variable_name()),
                          _ << "while processing process_variable variable: "
                            << config.pdo_variable().variable_name());
    INTR_RETURN_IF_ERROR(
        process_variable.template IsCompatibleType<T>(ignore_signedness))
        << "while processing process_variable variable: "
        << config.pdo_variable().variable_name();
    LOG(INFO) << "Using output process_variable variable `"
              << config.pdo_variable().variable_name() << "` for homing";
    return [process_variable](T value) mutable -> absl::StatusOr<T> {
      process_variable.WriteUnchecked(value);
      return value;
    };
  } else {
    uint32_t service_variable_index = config.has_sdo_variable()
                                          ? config.sdo_variable().index()
                                          : fallback_service_variable_index;
    uint32_t service_variable_subindex =
        config.has_sdo_variable() ? config.sdo_variable().subindex()
                                  : fallback_service_variable_subindex;
    // Configuration provides an service_variable variable.
    INTR_ASSIGN_OR_RETURN(
        auto service_variable,
        variable_registry.GetServiceVariable(
            service_variable_index, service_variable_subindex, bus_position),
        _ << absl::StrCat("while processing service_variable variable: 0x",
                          absl::Hex(service_variable_index, absl::kZeroPad4),
                          ".", absl::Hex(service_variable_subindex)));
    LOG(INFO) << "Using output service_variable variable 0x"
              << absl::Hex(service_variable_index, absl::kZeroPad4) << "."
              << absl::Hex(service_variable_subindex) << " for homing";
    return [service_variable](T value) mutable -> absl::StatusOr<T> {
      INTR_RETURN_IF_ERROR(service_variable.Write(value));
      return service_variable.template Read<T>();
    };
  }
}

template <typename T>
absl::StatusOr<absl::AnyInvocable<absl::StatusOr<T>(const T& expected_value,
                                                    absl::Duration timeout)>>
GetWaitFunctionFromConfig(
    int32_t bus_position,
    const intrinsic::fieldbus::VariableRegistry& variable_registry,
    const intrinsic_proto::icon::v1::BusVariable& config,
    uint32_t fallback_service_variable_index,
    uint32_t fallback_service_variable_subindex,
    bool ignore_signedness = false) {
  if (config.has_pdo_variable()) {
    INTR_ASSIGN_OR_RETURN(auto process_variable,
                          variable_registry.GetInputVariable(
                              config.pdo_variable().variable_name()),
                          _ << "while processing process_variable variable: "
                            << config.pdo_variable().variable_name());
    INTR_RETURN_IF_ERROR(
        process_variable.template IsCompatibleType<T>(ignore_signedness))
        << "while processing process_variable variable: "
        << config.pdo_variable().variable_name();
    LOG(INFO) << "Using input process_variable variable `"
              << config.pdo_variable().variable_name() << "` for homing";
    return [process_variable](const T& expected_value,
                              absl::Duration timeout) -> absl::StatusOr<T> {
      absl::Time start = absl::Now();
      while (true) {
        auto value = process_variable.template ReadUnchecked<T>();
        if (value == expected_value) {
          return value;
        }
        if (absl::Now() - start > timeout) {
          return absl::DeadlineExceededError(
              "Deadline exceeded while waiting for process_variable variable "
              "to have value "
              "");
        }
        absl::SleepFor(absl::Milliseconds(100));
      }
    };
  } else {
    uint32_t service_variable_index = config.has_sdo_variable()
                                          ? config.sdo_variable().index()
                                          : fallback_service_variable_index;
    uint32_t service_variable_subindex =
        config.has_sdo_variable() ? config.sdo_variable().subindex()
                                  : fallback_service_variable_subindex;

    INTR_ASSIGN_OR_RETURN(
        auto service_variable,
        variable_registry.GetServiceVariable(
            service_variable_index, service_variable_subindex, bus_position),
        _ << absl::StrCat("while processing service_variable variable: 0x",
                          absl::Hex(service_variable_index, absl::kZeroPad4),
                          ".", absl::Hex(service_variable_subindex)));
    LOG(INFO) << "Using input service_variable variable 0x"
              << absl::Hex(service_variable_index, absl::kZeroPad4) << "."
              << absl::Hex(service_variable_subindex) << " for homing";
    return [service_variable](
               const T& expected_value,
               absl::Duration timeout) mutable -> absl::StatusOr<T> {
      absl::Time start = absl::Now();
      while (true) {
        INTR_ASSIGN_OR_RETURN(auto value, service_variable.template Read<T>());
        if (value == expected_value) {
          return value;
        }
        if (absl::Now() - start > timeout) {
          return absl::DeadlineExceededError(
              "Deadline exceeded while waiting for service_variable variable "
              "to have "
              "value "
              "");
        }
        absl::SleepFor(absl::Milliseconds(100));
      }
    };
  }
}

absl::StatusOr<std::unique_ptr<HomingTask>> HomingTask::Create(
    int32_t bus_position,
    const intrinsic::fieldbus::VariableRegistry& variable_registry,
    intrinsic::icon::HardwareInterfaceRegistry& interface_registry,
    const intrinsic_proto::icon::v1::Ds402BusComponent& config) {
  if (!config.has_homing_config()) {
    // Homing is not configured.
    return nullptr;
  }

  INTR_ASSIGN_OR_RETURN(
      absl::AnyInvocable<absl::StatusOr<int8_t>(int8_t)> set_homing_method,
      GetWriteFunctionFromConfig<int8_t>(bus_position, variable_registry,
                                         config.homing_config().homing_method(),
                                         kHomingMethodIndex, 0));

  INTR_ASSIGN_OR_RETURN(
      absl::AnyInvocable<absl::StatusOr<int32_t>(int32_t)> set_homing_offset,
      GetWriteFunctionFromConfig<int32_t>(
          bus_position, variable_registry,
          config.homing_config().homing_offset(), kHomingOffsetIndex, 0));

  INTR_ASSIGN_OR_RETURN(
      absl::AnyInvocable<absl::StatusOr<uint32_t>(uint32_t)>
          set_homing_search_speed,
      GetWriteFunctionFromConfig<uint32_t>(
          bus_position, variable_registry,
          config.homing_config().homing_search_speed(), kHomingSpeedIndex, 1));

  INTR_ASSIGN_OR_RETURN(
      absl::AnyInvocable<absl::StatusOr<uint32_t>(uint32_t)>
          set_homing_creep_speed,
      GetWriteFunctionFromConfig<uint32_t>(
          bus_position, variable_registry,
          config.homing_config().homing_creep_speed(), kHomingSpeedIndex, 2));

  INTR_ASSIGN_OR_RETURN(absl::AnyInvocable<absl::StatusOr<uint32_t>(uint32_t)>
                            set_homing_acceleration,
                        GetWriteFunctionFromConfig<uint32_t>(
                            bus_position, variable_registry,
                            config.homing_config().homing_acceleration(),
                            kHomingAccelerationIndex, 0));

  SetHomingParamsAsyncType set_homing_params_async(
      [set_homing_method = std::move(set_homing_method),
       set_homing_offset = std::move(set_homing_offset),
       set_homing_search_speed = std::move(set_homing_search_speed),
       set_homing_creep_speed = std::move(set_homing_creep_speed),
       set_homing_acceleration = std::move(set_homing_acceleration)](
          int8_t method, int32_t offset, uint32_t search_speed,
          uint32_t creep_speed,
          uint32_t acceleration) mutable -> absl::StatusOr<int8_t> {
        INTR_ASSIGN_OR_RETURN(auto method_result, set_homing_method(method));
        if (method_result != method) {
          return intrinsic::icon::InternalError(
              absl::StrCat("SetHomingMethod function returned unexpected "
                           "value. Expected: ",
                           method, ", but got: ", method_result));
        }
        INTR_ASSIGN_OR_RETURN(auto offset_result, set_homing_offset(offset));
        if (offset_result != offset) {
          return intrinsic::icon::InternalError(
              absl::StrCat("SetHomingOffset function returned unexpected "
                           "value. Expected: ",
                           offset, ", but got: ", offset_result));
        }
        INTR_ASSIGN_OR_RETURN(auto search_speed_result,
                              set_homing_search_speed(search_speed));
        if (search_speed_result != search_speed) {
          return intrinsic::icon::InternalError(
              absl::StrCat("SetHomingSearchSpeed function returned unexpected "
                           "value. Expected: ",
                           search_speed, ", but got: ", search_speed_result));
        }
        INTR_ASSIGN_OR_RETURN(auto creep_speed_result,
                              set_homing_creep_speed(creep_speed));
        if (creep_speed_result != creep_speed) {
          return intrinsic::icon::InternalError(
              absl::StrCat("SetHomingCreepSpeed function returned unexpected "
                           "value. Expected: ",
                           creep_speed, ", but got: ", creep_speed_result));
        }
        INTR_ASSIGN_OR_RETURN(auto acceleration_result,
                              set_homing_acceleration(acceleration));
        if (acceleration_result != acceleration) {
          return intrinsic::icon::InternalError(
              absl::StrCat("SetHomingAcceleration function returned unexpected "
                           "value. Expected: ",
                           acceleration, ", but got: ", acceleration_result));
        }
        return method_result;
      });

  ::intrinsic_proto::icon::v1::ModesOfOperationConfig modes_of_operation_config;
  if (config.has_modes_of_operation_config()) {
    modes_of_operation_config = config.modes_of_operation_config();
  }
  INTR_ASSIGN_OR_RETURN(
      absl::AnyInvocable<absl::StatusOr<intrinsic_fbs::MotorControlMode>(
          intrinsic_fbs::MotorControlMode)>
          set_modes_of_operation,
      GetWriteFunctionFromConfig<intrinsic_fbs::MotorControlMode>(
          bus_position, variable_registry,
          modes_of_operation_config.output_variable(), kModesOfOperationIndex,
          0, /*ignore_signedness=*/true));

  // Call the set_modes_of_operation function once to initialize the modes of
  // operation.
  INTR_ASSIGN_OR_RETURN(
      auto initial_modes_of_operation,
      set_modes_of_operation(intrinsic_fbs::MotorControlMode::CYCLIC_POSITION));

  if (initial_modes_of_operation !=
      intrinsic_fbs::MotorControlMode::CYCLIC_POSITION) {
    return intrinsic::icon::InternalError(absl::StrCat(
        "Initial call to SetModesOfOperation function returned unexpected "
        "value. Expected: ",
        intrinsic_fbs::MotorControlMode::CYCLIC_POSITION,
        ", but got: ", initial_modes_of_operation));
  }

  INTR_ASSIGN_OR_RETURN(
      absl::AnyInvocable<absl::StatusOr<intrinsic_fbs::MotorControlMode>(
          intrinsic_fbs::MotorControlMode, absl::Duration)>
          wait_for_modes_of_operation_display,
      GetWaitFunctionFromConfig<intrinsic_fbs::MotorControlMode>(
          bus_position, variable_registry,
          modes_of_operation_config.display_variable(),
          kModesOfOperationDisplayIndex, 0, /*ignore_signedness=*/true));

  SetModesOfOperationAsyncType set_modes_of_operation_async(
      [set_modes_of_operation = std::move(set_modes_of_operation),
       wait_for_modes_of_operation_display =
           std::move(wait_for_modes_of_operation_display)](
          intrinsic_fbs::MotorControlMode modes_of_operation,
          absl::Duration settle_time, absl::Duration timeout) mutable
          -> absl::StatusOr<intrinsic_fbs::MotorControlMode> {
        absl::SleepFor(settle_time);
        INTR_ASSIGN_OR_RETURN(auto set_modes_of_operation_result,
                              set_modes_of_operation(modes_of_operation));
        if (set_modes_of_operation_result != modes_of_operation) {
          return intrinsic::icon::InternalError(
              absl::StrCat("SetModesOfOperation function returned unexpected "
                           "value. Expected: ",
                           modes_of_operation,
                           ", but got: ", set_modes_of_operation_result));
        }
        return wait_for_modes_of_operation_display(modes_of_operation, timeout);
      });

  INTR_ASSIGN_OR_RETURN(auto command_validator,
                        intrinsic::icon::Validator::Create(interface_registry));

  INTR_ASSIGN_OR_RETURN(
      auto homing_command,
      interface_registry.AdvertiseInterface<::intrinsic_fbs::HomeCommand>(
          std::string(config.device_name()) + "_homing_command", /*method=*/0,
          /*offset=*/0.0, /*search_speed=*/0.0, /*creep_speed=*/0.0,
          /*acceleration=*/0.0));
  INTR_ASSIGN_OR_RETURN(
      auto homing_status,
      interface_registry
          .AdvertiseMutableInterface<::intrinsic_fbs::HomingStatus>(
              std::string(config.device_name()) + "_homing_status",
              /*state=*/intrinsic_fbs::HomingStatusFlag::HomingUnknown,
              /*error_message=*/""));

  return absl::WrapUnique(new HomingTask(
      config.device_name(), std::move(set_homing_params_async),
      std::move(set_modes_of_operation_async), std::move(command_validator),
      std::move(homing_command), std::move(homing_status)));
}

intrinsic::icon::RealtimeStatus HomingTask::CyclicRead(
    const StatusWord& status_word) {
  switch (homing_state_) {
    case InternalHomingState::kHomingNotConfigured: {
      // Nothing to do, if homing is not configured.
      break;
    }
    case InternalHomingState::kHomingInternalError:
      [[fallthrough]];
    case InternalHomingState::kHomingError: {
      INTRINSIC_RT_LOG_THROTTLED(WARNING)
          << "[Homing: " << device_name_
          << "] homing error, state: " << homing_state_;
      // Detach the result future, in case it's still attached.
      // This could happen, if the homing task was interrupted middle of a
      // previous homing procedure.
      set_modes_of_operation_future_ =
          SetModesOfOperationAsyncType::FutureType::GetDetachedFuture();
      set_homing_params_future_ =
          SetHomingParamsAsyncType::FutureType::GetDetachedFuture();
    }
      [[fallthrough]];
    case InternalHomingState::kHomingInit: {
      if (!IsOperationEnabled(status_word)) {
        homing_state_ = InternalHomingState::kHomingInit;
        break;
      }
      // Trigger the reset sequence that sets the modes of operation to
      // kCyclicSynchronousPositioning and the homing method to 0.
      auto future_or_status = set_modes_of_operation_(
          /*modes_of_operation=*/intrinsic_fbs::MotorControlMode::
              CYCLIC_POSITION,
          /*settle_time=*/absl::ZeroDuration(),
          /*timeout=*/absl::Duration(kWaitForModesOfOperationDisplayTimeout));
      if (!future_or_status.ok()) {
        INTRINSIC_RT_LOG(ERROR) << "[Homing: " << device_name_
                                << "] set modes of operation failed: "
                                << future_or_status.status().message();
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  future_or_status.status().message());
        break;
      }
      set_modes_of_operation_future_ = std::move(future_or_status.value());
      homing_state_ = InternalHomingState::kHomingSettingModesOfOperation;
      break;
    }
    case InternalHomingState::kHomingIdle: {
      if (IsHomingCommanded()) {
        // We received a homing command.
        INTRINSIC_RT_LOG(INFO)
            << "[Homing: " << device_name_ << "] homing is commanded.";
        homing_state_ = InternalHomingState::kHomingCommanded;
      } else {
        // Nothing to do, if homing is not commanded.
        homing_state_ = InternalHomingState::kHomingIdle;
      }
      break;
    }
    case InternalHomingState::kHomingCommanded: {
      if (!IsOperationEnabled(status_word)) {
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  kHomingRequiresOperationEnabledErrorMessage);
        break;
      }
      // Check if the command is new.
      if (IsNewHomingCommand()) {
        INTRINSIC_RT_LOG(INFO) << "[Homing: " << device_name_
                               << "] new homing command with method: "
                               << homing_command_->method();
        last_homing_command_ = **homing_command_;
        // Trigger setting the homing params.
        int8_t method = last_homing_command_->method();
        int32_t offset = last_homing_command_->offset();
        uint32_t search_speed = last_homing_command_->search_speed();
        uint32_t creep_speed = last_homing_command_->creep_speed();
        uint32_t acceleration = last_homing_command_->acceleration();
        // Call the async function to set the homing params.
        auto future_or_status = set_homing_params_(
            std::move(method), std::move(offset), std::move(search_speed),
            std::move(creep_speed), std::move(acceleration));
        if (!future_or_status.ok()) {
          INTRINSIC_RT_LOG(ERROR)
              << "[Homing: " << device_name_ << "] set homing params failed: "
              << future_or_status.status().message();
          homing_state_ = InternalHomingState::kHomingError;
          intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                    future_or_status.status().message());
          break;
        }
        set_homing_params_future_ = std::move(future_or_status.value());
        homing_state_ = InternalHomingState::kHomingSettingParameters;
      } else {
        // No new homing command, nothing to do.
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  "Not a new homing command.");
        break;
      }
      break;
    }
    case InternalHomingState::kHomingSettingParameters: {
      // Ensure the drive is still in kOperationEnabled.
      if (!IsOperationEnabled(status_word)) {
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  kHomingRequiresOperationEnabledErrorMessage);
        break;
      }

      auto is_cancelled_or_status = set_homing_params_future_.IsCancelled();
      if (!is_cancelled_or_status.ok()) {
        INTRINSIC_RT_LOG(ERROR)
            << "[Homing: " << device_name_
            << "] set homing params failed to check if it was "
               "cancelled!";
        homing_state_ = InternalHomingState::kHomingInternalError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  "Internal error: set homing params failed to "
                                  "check if it was cancelled!");
        break;
      }

      if (is_cancelled_or_status.value()) {
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  "Set homing params has been cancelled.");
        break;
      }

      auto result = set_homing_params_future_.Get();
      if (!result.ok()) {
        // If it's a not available error, we'll try again next cycle.
        if (result.status().code() == absl::StatusCode::kUnavailable) {
          homing_state_ = InternalHomingState::kHomingSettingParameters;
          break;
        }
        INTRINSIC_RT_LOG(ERROR)
            << "[Homing: " << device_name_
            << "] set homing params failed: " << result.status().message();
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  result.status().message());
        break;
      } else {
        if (result.value().ok()) {
          if (result.value().value() == 0) {
            homing_state_ = InternalHomingState::kHomingIdle;
          } else {
            homing_state_ = InternalHomingState::kHomingParametersSet;
          }

        } else {
          INTRINSIC_RT_LOG(ERROR)
              << "[Homing: " << device_name_ << "] set homing params failed: "
              << result.value().status().message();
          homing_state_ = InternalHomingState::kHomingError;
          intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                    result.value().status().message());
          break;
        }
      }
    } break;
    case InternalHomingState::kHomingParametersSet: {
      // Ensure the drive is still in kOperationEnabled.
      if (!IsOperationEnabled(status_word)) {
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  kHomingRequiresOperationEnabledErrorMessage);
        break;
      }

      // We might still have a future from the init-step, so we need to detach
      // it.
      set_modes_of_operation_future_ =
          SetModesOfOperationAsyncType::FutureType::GetDetachedFuture();
      auto future_or_status = set_modes_of_operation_(
          /*modes_of_operation=*/intrinsic_fbs::MotorControlMode::HOME,
          /*settle_time=*/absl::ZeroDuration(),
          /*timeout=*/
          absl::Duration(kWaitForModesOfOperationDisplayTimeout));
      if (!future_or_status.ok()) {
        INTRINSIC_RT_LOG(ERROR) << "[Homing: " << device_name_
                                << "] set modes of operation failed: "
                                << future_or_status.status().message();
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  future_or_status.status().message());
        break;
      }
      set_modes_of_operation_future_ = std::move(future_or_status.value());
      homing_state_ = InternalHomingState::kHomingSettingModesOfOperation;
      break;
    }
    case InternalHomingState::kHomingSettingModesOfOperation: {
      // Ensure the drive is still in kOperationEnabled.
      if (!IsOperationEnabled(status_word)) {
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  kHomingRequiresOperationEnabledErrorMessage);
        break;
      }

      auto is_cancelled_or_status =
          set_modes_of_operation_future_.IsCancelled();
      if (!is_cancelled_or_status.ok()) {
        INTRINSIC_RT_LOG(ERROR)
            << "[Homing: " << device_name_
            << "] set modes of operation failed to check if "
               "it was cancelled!";
        homing_state_ = InternalHomingState::kHomingInternalError;
        intrinsic_fbs::StringCopy(
            &homing_status_->mutable_error_message(),
            "Internal error: set modes of operation failed to "
            "check if it was cancelled!");
        break;
      }

      if (is_cancelled_or_status.value()) {
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  "Set modes of operation has been cancelled.");
        break;
      }

      auto result = set_modes_of_operation_future_.Get();
      if (!result.ok()) {
        // If it's a not available error, we'll try again next cycle.
        if (result.status().code() == absl::StatusCode::kUnavailable) {
          homing_state_ = InternalHomingState::kHomingSettingModesOfOperation;
          break;
        }
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  result.status().message());
        break;
      } else {
        if (result.value().ok()) {
          if (result.value().value() == intrinsic_fbs::MotorControlMode::HOME) {
            homing_state_ = InternalHomingState::kHomingModesOfOperationReached;
            break;
          } else {
            homing_state_ = InternalHomingState::kHomingIdle;
            break;
          }
          break;
        } else {
          INTRINSIC_RT_LOG(ERROR) << "[Homing: " << device_name_
                                  << "] set modes of operation failed: "
                                  << result.value().status().message();
          homing_state_ = InternalHomingState::kHomingError;
          intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                    result.value().status().message());
          break;
        }
      }
      break;
    }
    case InternalHomingState::kHomingModesOfOperationReached: {
      // Ensure the drive is still in kOperationEnabled.
      if (!IsOperationEnabled(status_word)) {
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  kHomingRequiresOperationEnabledErrorMessage);
        break;
      }

      // Get the 3 homing related status flags from the status_word.
      bool target_reached = status_word[kTargetReachedBitIndex];
      bool homing_attained = status_word[kHomingAttainedBitIndex];
      bool homing_error = status_word[kHomingErrorBitIndex];

      if (homing_error) {
        INTRINSIC_RT_LOG(ERROR) << "[Homing: " << device_name_
                                << "] Device reported a homing error.";
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  "Device reported a homing error.");
        break;
      }

      if (target_reached && homing_attained) {
        // Return the futures, so that we're resetting the modes of operation
        // and homing method in kHomingIdle.
        auto future_or_status = set_modes_of_operation_(
            /*modes_of_operation=*/intrinsic_fbs::MotorControlMode::
                CYCLIC_POSITION,
            /*settle_time=*/
            absl::Duration(kSetModesOfOperationSettleTimeAfterHomingAttained),
            /*timeout=*/
            absl::Duration(kWaitForModesOfOperationDisplayTimeout));
        if (!future_or_status.ok()) {
          INTRINSIC_RT_LOG(ERROR) << "[Homing: " << device_name_
                                  << "] set modes of operation failed: "
                                  << future_or_status.status().message();
          homing_state_ = InternalHomingState::kHomingError;
          intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                    future_or_status.status().message());
          break;
        }
        set_modes_of_operation_future_ = std::move(future_or_status.value());
        homing_state_ =
            InternalHomingState::kHomingAttainedRestoringModesOfOperation;
        break;
      }
      homing_state_ = InternalHomingState::kHomingModesOfOperationReached;
      break;
    }
    case InternalHomingState::kHomingAttainedRestoringModesOfOperation: {
      // Ensure the drive is still in kOperationEnabled.
      if (!IsOperationEnabled(status_word)) {
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  kHomingRequiresOperationEnabledErrorMessage);
        break;
      }
      auto has_value_or = set_modes_of_operation_future_.HasValue();
      if (has_value_or.ok()) {
        if (has_value_or.value()) {
          homing_state_ = InternalHomingState::kHomingAttained;
        } else {
          homing_state_ =
              InternalHomingState::kHomingAttainedRestoringModesOfOperation;
        }
        break;
      } else {
        INTRINSIC_RT_LOG(ERROR)
            << "[Homing: " << device_name_
            << "] failed to get the value of the set modes of operation "
               "future: "
            << has_value_or.status().message();
        homing_state_ = InternalHomingState::kHomingError;
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  has_value_or.status().message());
        break;
      }
    }
    case InternalHomingState::kHomingAttained: {
      INTRINSIC_RT_LOG(INFO)
          << "[Homing: " << device_name_ << "] homing attained.";
      homing_state_ = InternalHomingState::kHomingIdle;
    }
  }
  return UpdateHomingStatus();
}

intrinsic::icon::RealtimeStatusOr<HomingTask::ControlWord>
HomingTask::CyclicWrite(const HomingTask::ControlWord& control_word) {
  HomingTask::ControlWord updated_control_word = control_word;

  if (homing_state_ == InternalHomingState::kHomingModesOfOperationReached) {
    // Make sure the homing control word bit is set.
    updated_control_word[kHomingBitIndex] = true;
  } else {
    // Make sure the homing control word bit is cleared.
    updated_control_word[kHomingBitIndex] = false;
  }
  return updated_control_word;
}

HomingTask::HomingTask(
    absl::string_view device_name, SetHomingParamsAsyncType&& set_homing_params,
    SetModesOfOperationAsyncType&& set_modes_of_operation,
    intrinsic::icon::Validator command_validator,
    intrinsic::icon::HardwareInterfaceHandle<::intrinsic_fbs::HomeCommand>
        homing_command,
    intrinsic::icon::MutableHardwareInterfaceHandle<
        ::intrinsic_fbs::HomingStatus>
        homing_status)
    : device_name_(device_name),
      set_homing_params_(std::move(set_homing_params)),
      set_modes_of_operation_(std::move(set_modes_of_operation)),
      homing_command_(std::move(homing_command)),
      homing_status_(std::move(homing_status)),
      homing_state_(InternalHomingState::kHomingInit),
      command_validator_(std::move(command_validator)) {}

HomingTask::HomingTask(HomingTask&& other)
    : device_name_(std::move(other.device_name_)) {
  set_homing_params_ = std::move(other.set_homing_params_);
  set_modes_of_operation_ = std::move(other.set_modes_of_operation_);
  homing_command_ = std::move(other.homing_command_);
  homing_status_ = std::move(other.homing_status_);
  homing_state_ = std::move(other.homing_state_);
  command_validator_ = std::move(other.command_validator_);
  set_homing_params_future_ = std::move(other.set_homing_params_future_);
  set_modes_of_operation_future_ =
      std::move(other.set_modes_of_operation_future_);
  last_homing_command_ = std::move(other.last_homing_command_);
}

bool HomingTask::IsHomingCommanded() const {
  auto command_update_status =
      command_validator_.WasUpdatedThisCycle(homing_command_);
  if (!command_update_status.ok()) {
    return false;
  }
  // The command is valid if it's not NONE.
  return homing_command_->method() != 0;
}

bool HomingTask::IsNewHomingCommand() const {
  if (!last_homing_command_.has_value()) {
    return true;
  }
  if (*last_homing_command_ != **homing_command_) {
    return true;
  }
  return false;
}

intrinsic::icon::RealtimeStatus HomingTask::UpdateHomingStatus() {
  switch (homing_state_) {
    case InternalHomingState::kHomingParametersSet:
      [[fallthrough]];
    case InternalHomingState::kHomingSettingParameters:
      [[fallthrough]];
    case InternalHomingState::kHomingSettingModesOfOperation:
      [[fallthrough]];
    case InternalHomingState::kHomingAttainedRestoringModesOfOperation:
      [[fallthrough]];
    case InternalHomingState::kHomingModesOfOperationReached:
      homing_status_->mutate_state(
          intrinsic_fbs::HomingStatusFlag::HomingInProgress);
      return intrinsic::icon::OkStatus();
      break;
    case InternalHomingState::kHomingAttained:
      homing_status_->mutate_state(
          intrinsic_fbs::HomingStatusFlag::HomingAttained);
      return intrinsic::icon::OkStatus();
      break;
    case InternalHomingState::kHomingCommanded:
      homing_status_->mutate_state(
          intrinsic_fbs::HomingStatusFlag::HomingNotStarted);
      return intrinsic::icon::OkStatus();
      break;
    case InternalHomingState::kHomingError:
      if (homing_status_->error_message().size() == 0) {
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  "Internal error: Unknown error.");
      }
      [[fallthrough]];
    case InternalHomingState::kHomingInternalError:
      if (homing_status_->error_message().size() == 0) {
        intrinsic_fbs::StringCopy(&homing_status_->mutable_error_message(),
                                  "Internal error: Unknown internal error.");
      }
      homing_status_->mutate_state(
          intrinsic_fbs::HomingStatusFlag::HomingError);
      last_homing_command_ = std::nullopt;
      return intrinsic::icon::OkStatus();
      break;
    case InternalHomingState::kHomingInit:
      [[fallthrough]];
    case InternalHomingState::kHomingNotConfigured:
      [[fallthrough]];
    case InternalHomingState::kHomingIdle:
      last_homing_command_ = std::nullopt;
      homing_status_->mutate_state(
          intrinsic_fbs::HomingStatusFlag::HomingUnknown);
      return intrinsic::icon::OkStatus();
      break;
  }

  return intrinsic::icon::InternalError(
      "Internal error: Unknown homing state.");
}

}  // namespace intrinsic::ds402
