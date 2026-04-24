// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/ds402/ds402_bus_component.h"

#include <unistd.h>

#include <bitset>
#include <cmath>
#include <cstddef>
#include <cstdint>
#include <limits>
#include <memory>
#include <optional>
#include <string>
#include <utility>
#include <vector>

#include "absl/log/check.h"
#include "absl/log/log.h"
#include "absl/memory/memory.h"
#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_format.h"
#include "absl/strings/str_join.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/notification.h"
#include "absl/time/time.h"
#include "intrinsic/icon/hal/interfaces/electrical_motor_hardware_interfaces.h"  // IWYU pragma: keep (registers HomeCommand hardware interface)
#include "intrinsic/icon/hal/lib/adio/v1/adio_bus_component_config.pb.h"
#include "intrinsic/icon/hal/lib/ds402/ds402_bus_component_utils.h"
#include "intrinsic/icon/hal/lib/ds402/ds402_driver.h"
#include "intrinsic/icon/hal/lib/ds402/homing_task.h"
#include "intrinsic/icon/hal/lib/ds402/v1/ds402_bus_component_config.pb.h"
#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/service_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/v1/value_parsing.pb.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"
#include "intrinsic/icon/utils/bitset.h"
#include "intrinsic/icon/utils/fixed_str_cat.h"
#include "intrinsic/icon/utils/fixed_string.h"
#include "intrinsic/icon/utils/log.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_macro.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/platform/common/buffers/realtime_write_queue.h"
#include "intrinsic/platform/common/buffers/rt_queue.h"
#include "intrinsic/util/proto_time.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/util/thread/rt_thread.h"
#include "intrinsic/util/thread/stop_token.h"
#include "intrinsic/util/thread/thread.h"
#include "intrinsic/util/thread/thread_options.h"

namespace intrinsic::ds402 {

// Ds402 devices have an internal position limit, to limit integer overflow. The
// default value is more conservative than it needs to be, so we set it to the
// max value allowed.
constexpr unsigned int kMaxPositionLimit = (1u << 30) - 1;

constexpr size_t kWarningBitIndex = 7;
constexpr size_t kInternalLimitsActiveBitIndex = 11;

using ::intrinsic::icon::RealtimeStatus;
using ::intrinsic_proto::icon::v1::Ds402UnexpectedStateTransitionConfiguration;

Ds402BusComponent::ProcessVariableInterpretationResult
Ds402BusComponent::GetProcessVariableInterpretationFromConfig(
    ProcessVariableType process_variable_value,
    const intrinsic_proto::icon::v1::Ds402BusComponent&
        config_with_interpretations) {
  Ds402BusComponent::ProcessVariableInterpretationResult interpretation{
      .error_code = process_variable_value};
  std::vector<std::string> result;

  for (const auto& [error_interpretation, interpretations] :
       config_with_interpretations.error_code_interpretations()) {
    for (const auto& pattern : interpretations.patterns())
      if ((process_variable_value & pattern.signal_mask()) ==
          pattern.expected_value()) {
        result.push_back(error_interpretation);
      }
  }
  // Resulting string is automatically truncated to the max supported length.
  interpretation.interpretation =
      Ds402BusComponent::ErrorInterpretationString(absl::StrJoin(result, "; "));
  return interpretation;
}

std::optional<Ds402BusComponent::ProcessVariableInterpretationResult>
Ds402BusComponent::ProcessVariableInterpretationProvider::
    GetProcessVariableInterpretation(ProcessVariableType error_code) {
  if (!interpretation_lookup_enabled_) {
    return std::nullopt;
  }

  // Read translation result from the queue.
  // Keeps the queue empty.
  if (const auto result = result_queue_.reader()->Pop(); result.has_value()) {
    // Ignore result if it's not for the current error code.
    if (result->error_code == error_code) {
      current_interpretation_ = std::move(result.value());
    }
  }

  if (!current_interpretation_.has_value() ||
      current_interpretation_->error_code != error_code) {
    // Stored interpretation doesn't match. Request new interpretation.
    if (const auto result = RequestInterpretation(error_code); !result.ok()) {
      // Request failed.
      current_interpretation_ = std::nullopt;
      Ds402BusComponent::ProcessVariableInterpretationResult error_code_result{
          .error_code = 0x00,
          .interpretation =
              Ds402BusComponent::ErrorInterpretationString(result.message())};
      return error_code_result;
    }
    current_interpretation_ = std::nullopt;
    return std::nullopt;
  }
  // current_interpretation_->error_code == error_code
  return current_interpretation_;
}

RealtimeStatus
Ds402BusComponent::ProcessVariableInterpretationProvider::RequestInterpretation(
    ProcessVariableType error_code) {
  if (!interpretation_lookup_enabled_) {
    return intrinsic::icon::OkStatus();
  }

  if (bool write_result = request_queue_.Writer().Write(error_code);
      !write_result) {
    return intrinsic::icon::InternalError(
        intrinsic::icon::FixedStrCat<RealtimeStatus::kMaxMessageLength>(
            "Failed to request interpretations of error code 0x",
            absl::Hex(error_code)));
  }
  return intrinsic::icon::OkStatus();
}

Ds402BusComponent::ProcessVariableInterpretationProvider::
    ProcessVariableInterpretationProvider(bool interpretation_lookup_enabled)
    : interpretation_lookup_enabled_(interpretation_lookup_enabled) {
  if (!interpretation_lookup_enabled_) {
    LOG(INFO) << "Skipping error code interpretation, because no "
                 "interpretations are provided.";
  }
};

Ds402BusComponent::ProcessVariableInterpretationProvider::
    ~ProcessVariableInterpretationProvider() {
  request_queue_.Writer().Close();
  non_rt_translation_thread_.request_stop();
}

absl::StatusOr<
    std::unique_ptr<Ds402BusComponent::ProcessVariableInterpretationProvider>>
Ds402BusComponent::ProcessVariableInterpretationProvider::Create(
    const intrinsic_proto::icon::v1::Ds402BusComponent& config) {
  bool interpretation_lookup_enabled =
      !config.error_code_interpretations().empty();
  auto process_variable_interpretation = absl::WrapUnique(
      new ProcessVariableInterpretationProvider(interpretation_lookup_enabled));

  if (!interpretation_lookup_enabled) {
    return process_variable_interpretation;
  }

  ProcessVariableInterpretationProvider* process_variable_interpretation_ptr =
      process_variable_interpretation.get();
  absl::Notification interpretation_thread_started;

  INTR_ASSIGN_OR_RETURN(
      process_variable_interpretation->non_rt_translation_thread_,
      intrinsic::CreateRealtimeCapableThread(
          intrinsic::ThreadOptions().SetName("process_variable_interpr"),
          [process_variable_interpretation_ptr, &interpretation_thread_started,
           config](intrinsic::StopToken st) {
            interpretation_thread_started.Notify();
            ProcessVariableType request;
            while (!st.stop_requested() &&
                   process_variable_interpretation_ptr->request_queue_.Reader()
                           .Read(request) == intrinsic::ReadResult::kConsumed) {
              const auto interpretation =
                  GetProcessVariableInterpretationFromConfig(request, config);
              if (bool write_result =
                      process_variable_interpretation_ptr->result_queue_
                          .writer()
                          ->Insert(interpretation);
                  !write_result) {
                LOG_EVERY_N_SEC(WARNING, 10)
                    << "Failed to write process variable interpretation result "
                       "for error "
                       "code 0x"
                    << absl::Hex(request) << " interpretation ["
                    << interpretation.interpretation
                    << "]. Is the realtime thread stuck?";
              }
            }
          }));

  interpretation_thread_started.WaitForNotification();

  return process_variable_interpretation;
}

Ds402BusComponent::Ds402BusComponent(
    absl::string_view device_name, fieldbus::ProcessVariable status_word,
    std::optional<fieldbus::ProcessVariable> error_code, fieldbus::ProcessVariable control_word,
    std::optional<fieldbus::ProcessVariable> digital_outputs,
    bool treat_internal_limits_active_as_error,
    UnexpectedStateTransitionConfiguration
        unexpected_state_transition_configuration,
    int64_t enable_delay_cycles, Ds402State initial_goal_state,
    Ds402State enabled_goal_state,
    std::unique_ptr<ProcessVariableInterpretationProvider> process_variable_interpretation,
    std::unique_ptr<HomingTask> homing_task, uint32_t brake_release_bit_value)
    : device_state_{
          .status_word = 0,
          .error_code = 0,
          .control_word = intrinsic::FromEnum(ControlWord::kDisableVoltage),
          .digital_outputs = 0,
      },
      goal_state_(initial_goal_state),
      ds402_state_(Ds402State::kNotReadyToSwitchOn),
      prev_state_(Ds402State::kNotReadyToSwitchOn),
      initial_goal_state_(initial_goal_state),
      enabled_goal_state_(enabled_goal_state),
      is_faulted_(false),
      fault_message_(""),
      device_name_(device_name),
      treat_internal_limits_active_as_error_(
          treat_internal_limits_active_as_error),
      unexpected_state_transition_configuration_(
          std::move(unexpected_state_transition_configuration)),
      status_word_(status_word),
      error_code_(error_code),
      control_word_(control_word),
      digital_outputs_(digital_outputs),
      enable_delay_cycles_(enable_delay_cycles),
      remaining_enable_delay_cycles_(enable_delay_cycles),
      process_variable_interpretation_provider_(std::move(process_variable_interpretation)),
      homing_task_(std::move(homing_task)),
      brake_release_bit_value_(brake_release_bit_value) {}

absl::StatusOr<std::unique_ptr<Ds402BusComponent>> Ds402BusComponent::Create(
    fieldbus::DeviceInitContext& device_init_context,
    const intrinsic_proto::icon::v1::Ds402BusComponent& config,
    double frequency, bool distributed_clock_and_bus_shift_sync_are_enabled) {
  if (!distributed_clock_and_bus_shift_sync_are_enabled) {
    return absl::InvalidArgumentError(
        "Distributed clock and the bus shift DCM mode must be enabled for "
        "DS402. Please check the bus.distributed_clock settings of the "
        "EtherCAT hardware module.");
  }

  const intrinsic::fieldbus::VariableRegistry& variable_registry =
      device_init_context.GetVariableRegistry();

  INTR_ASSIGN_OR_RETURN(
      auto status_word_variable,
      variable_registry.GetInputVariable(config.status_word_variable_name()),
      _ << "while processing `status_word` variable: "
        << config.status_word_variable_name());
  INTR_RETURN_IF_ERROR(
      status_word_variable.IsCompatibleType<ProcessVariableType>())
      << "while processing `status_word` variable: "
      << config.status_word_variable_name();

  std::optional<fieldbus::ProcessVariable> error_code_variable;
  if (!config.error_code_variable_name().empty()) {
    INTR_ASSIGN_OR_RETURN(
        error_code_variable,
        variable_registry.GetInputVariable(config.error_code_variable_name()),
        _ << "while processing `error_code` variable: "
          << config.error_code_variable_name());
    INTR_RETURN_IF_ERROR(
        error_code_variable->IsCompatibleType<ProcessVariableType>())
        << "while processing `error_code` variable name: "
        << config.error_code_variable_name();
  }

  INTR_ASSIGN_OR_RETURN(
      auto control_word_variable,
      variable_registry.GetOutputVariable(config.control_word_variable_name()),
      _ << "while processing `control_word` variable: "
        << config.control_word_variable_name());
  INTR_RETURN_IF_ERROR(
      control_word_variable.IsCompatibleType<ProcessVariableType>())
      << "while processing `control_word` variable: "
      << config.control_word_variable_name();

  std::optional<fieldbus::ProcessVariable> digital_outputs_variable;
  if (!config.digital_outputs_variable_name().empty()) {
    INTR_ASSIGN_OR_RETURN(digital_outputs_variable,
                          variable_registry.GetOutputVariable(
                              config.digital_outputs_variable_name()),
                          _ << "while processing `digital_outputs` variable: "
                            << config.digital_outputs_variable_name());
    INTR_RETURN_IF_ERROR(
        digital_outputs_variable->IsCompatibleType<DigitalOutputsType>())
        << "while processing `digital_outputs` variable: "
        << config.digital_outputs_variable_name();
  }

  std::unique_ptr<HomingTask> homing_task{nullptr};
  if (config.has_homing_config()) {
    INTR_ASSIGN_OR_RETURN(
        homing_task,
        HomingTask::Create(config.bus_position(), variable_registry,
                           device_init_context.GetInterfaceRegistry(), config));
  }

  if (config.enable_automatic_interpolation_window()) {
    INTR_RETURN_IF_ERROR(SetInterpolationWindow(
        variable_registry, config.bus_position(), frequency));
  }

  INTR_ASSIGN_OR_RETURN(absl::Duration enable_delay,
                        intrinsic::ToAbslDuration(config.enable_delay()));
  int64_t enable_delay_ms = absl::ToInt64Milliseconds(enable_delay);
  if (enable_delay_ms < 0) {
    return absl::InvalidArgumentError(absl::StrFormat(
        "Enable delay must be positive, got %d ms.", enable_delay_ms));
  }

  int64_t enable_delay_cycles = static_cast<int64_t>(
      static_cast<double>(enable_delay_ms) / 1000.0 * frequency);

  Ds402State initial_goal_state = Ds402State::kSwitchOnDisabled;
  Ds402State enabled_goal_state = Ds402State::kOperationEnabled;
  if (config.has_initial_target_state()) {
    initial_goal_state = Ds402State(config.initial_target_state());
    INTRINSIC_RT_LOG(INFO) << "Using non-default initial target state: "
                           << ToString(initial_goal_state);
  }
  if (config.has_enabled_target_state()) {
    enabled_goal_state = Ds402State(config.enabled_target_state());
    INTRINSIC_RT_LOG(INFO) << "Using non-default enabled target state: "
                           << ToString(enabled_goal_state);
  }

  INTR_ASSIGN_OR_RETURN(auto process_variable_interpretation,
                        ProcessVariableInterpretationProvider::Create(config));

  INTR_ASSIGN_OR_RETURN(
      auto unexpected_state_transition_configuration,
      ParseUnexpectedStateTransitionConfiguration(variable_registry, config));

  uint32_t brake_release_bit_value = Ds402BusComponent::kBrakeReleaseBitValue;
  if (config.has_brake_release_bit_value()) {
    // The proto field is a signed int64, but we only accept uint32_t values.
    auto signed_brake_release_bit_value = config.brake_release_bit_value();
    if (signed_brake_release_bit_value < 0) {
      return absl::InvalidArgumentError(
          absl::StrFormat("brake_release_bit_value must be positive, got %d.",
                          signed_brake_release_bit_value));
    }
    brake_release_bit_value =
        static_cast<uint32_t>(config.brake_release_bit_value());
  }
  INTRINSIC_RT_LOG(INFO) << "Using brake_release_bit_value: "
                         << brake_release_bit_value << " (0x"
                         << absl::Hex(brake_release_bit_value) << ")";

  return absl::WrapUnique(new Ds402BusComponent(
      config.device_name(), status_word_variable, error_code_variable,
      control_word_variable, digital_outputs_variable,
      config.treat_internal_limits_active_as_error(),
      std::move(unexpected_state_transition_configuration), enable_delay_cycles,
      initial_goal_state, enabled_goal_state,
      std::move(process_variable_interpretation), std::move(homing_task),
      brake_release_bit_value));
}

intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus>
Ds402BusComponent::CyclicRead(fieldbus::RequestType request_type) {
  // Remember the previous state.
  prev_state_ = ds402_state_;

  // Read status word and error code.
  device_state_.status_word = status_word_.ReadUnchecked<ProcessVariableType>();
  if (error_code_.has_value()) {
    device_state_.error_code =
        error_code_->ReadUnchecked<ProcessVariableType>();
  } else {
    device_state_.error_code = 0;
  }

  // Update `ds402_state_` to the current state from the status word.
  INTRINSIC_RT_ASSIGN_OR_RETURN(ds402_state_,
                                ToDs402State(device_state_.status_word));

  switch (request_type) {
    case fieldbus::RequestType::kEnableMotion:
      // Just stay in the enabled goal state if we're already there.
      if (ds402_state_ == enabled_goal_state_) {
        goal_state_ = enabled_goal_state_;
        break;
      }
      // Only go to kOperationEnabled if the enable delay has passed, instead
      // go to kReadyToSwitchOn.
      if (--remaining_enable_delay_cycles_ < 0) {
        remaining_enable_delay_cycles_ = 0;
        goal_state_ = enabled_goal_state_;
      } else {
        goal_state_ = initial_goal_state_;
      }
      break;
    case fieldbus::RequestType::kClearFaults:
      [[fallthrough]];
    case fieldbus::RequestType::kDisableMotion:
      // Reset the enable delay.
      remaining_enable_delay_cycles_ = enable_delay_cycles_;
      goal_state_ = Ds402State::kSwitchOnDisabled;
      break;
    default:
      // Reset the enable delay.
      remaining_enable_delay_cycles_ = enable_delay_cycles_;
      // Maintain the current goal_state_.
      break;
  }

  OptionalStatusBits optional_status_bits =
      ParseOptionalStatusBits(device_state_.status_word);

  bool new_fault = false;
  const bool prev_state_was_enabled = prev_state_ == enabled_goal_state_;
  const bool prev_state_was_goal_state = prev_state_ == goal_state_;

  if ((prev_state_was_enabled) && (prev_state_was_goal_state) &&
      (ds402_state_ != goal_state_)) {
    fault_message_ = RealtimeStatus::StrCat(
        "Device `", device_name_, "` unexpectedly transitioned from state `",
        ToString(enabled_goal_state_), "` to state `", ToString(ds402_state_),
        "`.");
    new_fault = true;
    // When the device switches from the enabled state to a non-enabled state,
    // CycleRead will typically fault.
    // Potentially override the fault based on the configured behavior and
    // signal.
    // When the fault is overridden, CyclicRead returns `kProcessing`, but this
    // is currently ignored by the EtherCAT Hardware Module.
    // This means there is no indication to the user that the device is in an
    // unexpected state and ICON can still send commands to the device,
    // even though the hardware is not in the enabled state and will most likely
    // ignore the command.
    // This can lead to a fault once the device transitions back to the enabled
    // state.
    switch (unexpected_state_transition_configuration_.behavior) {
      case Ds402UnexpectedStateTransitionConfiguration::BEHAVIOR_UNSPECIFIED:
        [[fallthrough]];
      case Ds402UnexpectedStateTransitionConfiguration::BEHAVIOR_FAULT: {
        break;
      }
      case Ds402UnexpectedStateTransitionConfiguration::
          BEHAVIOR_IGNORE_ON_SIGNAL: {
        if (optional_status_bits.ignore_unexpected_state_transition_signal) {
          // switch (request_type) below sets current_request_status_ =
          // kProcessing and not fault.
          new_fault = false;
          INTRINSIC_RT_LOG_THROTTLED(WARNING)
              << "`ignore_unexpected_state_transition_signal` "
                 "active. Ignoring: "
              << fault_message_;
        }
        break;
      }
      case Ds402UnexpectedStateTransitionConfiguration::
          BEHAVIOR_IGNORE_ALWAYS: {
        new_fault = false;
        INTRINSIC_RT_LOG_THROTTLED(WARNING)
            << "Configured to always ignore unexpected state transition. "
            << "Ignoring: " << fault_message_;
        break;
      }
      default:
        INTRINSIC_RT_LOG_THROTTLED(ERROR)
            << "Undefined configuration for unexpected state transition. Got: "
            << unexpected_state_transition_configuration_.behavior;
    }
  }

  if (const bool has_device_error = device_state_.error_code != 0;
      has_device_error) {
    new_fault = true;
    fault_message_ = RealtimeStatus::StrCat(
        "Device `", device_name_, "` reported error code 0x",
        absl::Hex(device_state_.error_code), ".");

    if (const auto& interpretation =
            process_variable_interpretation_provider_
                ->GetProcessVariableInterpretation(device_state_.error_code);
        interpretation.has_value()) {
      // Error matches interpretation.
      fault_message_ = RealtimeStatus::StrCat(
          "Device `", device_name_, "` reported error code 0x",
          absl::Hex(device_state_.error_code), ": ",
          interpretation->interpretation);
    }
  }

  if (optional_status_bits.has_internal_limits_active) {
    intrinsic::icon::FixedString<
        intrinsic::icon::RealtimeStatus::kMaxMessageLength>
        message = RealtimeStatus::StrCat(
            "Setpoint for device `", device_name_,
            "` can't be reached because of internal limits (e.g. hardware "
            "position switches, current limiter or thermal overload).");

    if (treat_internal_limits_active_as_error_) {
      new_fault = true;
      fault_message_ = message;
    }
    INTRINSIC_RT_LOG_THROTTLED(WARNING) << message;
  }

  if (optional_status_bits.has_warning) {
    INTRINSIC_RT_LOG_THROTTLED(WARNING)
        << "Device `" << device_name_ << "` reported Warning 0x"
        << absl::Hex(device_state_.error_code) << ".";
  }

  if (homing_task_ != nullptr) {
    auto status = homing_task_->CyclicRead(device_state_.status_word);
    if (!status.ok()) {
      new_fault = true;
      fault_message_ = RealtimeStatus::StrCat(
          "`", device_name_, "` homing error: ", status.message());
    }
  }

  switch (request_type) {
    case fieldbus::RequestType::kDisableMotion:
      [[fallthrough]];
    case fieldbus::RequestType::kNormalOperation:
      [[fallthrough]];
    default:
      if (new_fault || is_faulted_) {
        goal_state_ = initial_goal_state_;
        is_faulted_ = true;
        current_request_status_ = fieldbus::RequestStatus::kProcessing;
        return UnavailableError(fault_message_);
      }
      if (ds402_state_ == goal_state_) {
        current_request_status_ = fieldbus::RequestStatus::kDone;
      } else {
        current_request_status_ = fieldbus::RequestStatus::kProcessing;
      }
      break;
    case fieldbus::RequestType::kEnableMotion:
      if (new_fault || is_faulted_) {
        goal_state_ = initial_goal_state_;
        is_faulted_ = true;
        current_request_status_ = fieldbus::RequestStatus::kProcessing;
        return UnavailableError(fault_message_);
      }
      if (ds402_state_ == enabled_goal_state_) {
        current_request_status_ = fieldbus::RequestStatus::kDone;
      } else {
        current_request_status_ = fieldbus::RequestStatus::kProcessing;
      }
      break;
    case fieldbus::RequestType::kClearFaults:
      if (ds402_state_ == goal_state_ && !new_fault) {
        is_faulted_ = false;
        current_request_status_ = fieldbus::RequestStatus::kDone;
      } else {
        current_request_status_ = fieldbus::RequestStatus::kProcessing;
      }
      break;
  }
  return current_request_status_;
}

intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus>
Ds402BusComponent::CyclicWrite(fieldbus::RequestType request_type) {
  FaultHandling fault_handling_behavior = FaultHandling::kPreserve;
  if (request_type == fieldbus::RequestType::kClearFaults) {
    fault_handling_behavior = FaultHandling::kClear;
  }

  INTRINSIC_RT_ASSIGN_OR_RETURN(
      auto next_control_word,
      GetNextControlWord(
          prev_state_, ds402_state_, goal_state_,
          intrinsic::ToEnum<ControlWord>(device_state_.control_word),
          fault_handling_behavior));
  device_state_.control_word = intrinsic::FromEnum(next_control_word);

  // Engage the brakes if we're not enabling and (faulted or not (in
  // kOperationEnabled or kQuickStopActive)), release otherwise.
  if ((request_type != fieldbus::RequestType::kEnableMotion) &&
      (is_faulted_ || !(ds402_state_ == Ds402State::kOperationEnabled ||
                        ds402_state_ == Ds402State::kQuickStopActive))) {
    device_state_.digital_outputs = 0;
  } else {
    device_state_.digital_outputs = brake_release_bit_value_;
  }

  if (homing_task_ != nullptr) {
    INTRINSIC_RT_ASSIGN_OR_RETURN(
        device_state_.control_word,
        homing_task_->CyclicWrite(device_state_.control_word));
  }
  // Write control word and digital outputs.
  control_word_.WriteUnchecked(device_state_.control_word);
  if (digital_outputs_.has_value()) {
    digital_outputs_->WriteUnchecked(device_state_.digital_outputs);
  }

  return current_request_status_;
}

absl::Status Ds402BusComponent::SetInterpolationWindow(
    const intrinsic::fieldbus::VariableRegistry& variable_registry,
    int32_t bus_position, double frequency) {
  if (frequency <= 0) {
    return absl::InvalidArgumentError(
        absl::StrFormat("[DS402] Requested frequency (%.1f) "
                        "must be positive.",
                        frequency));
  } else {
    INTRINSIC_RT_LOG(INFO)
        << "Ds402BusComponent configured with a frequency of: " << frequency
        << " Hz.";
  }
  // Transform the init time to a power of 10. We have -8 on the exponent to
  // scale to use the full 8 bits integer range.
  int exp[2], val[2];
  exp[0] = std::ceil(std::log10(pow(2, -8)));
  exp[1] = std::floor(std::log10(pow(2, -8)));
  const double init_time = 1.0 / frequency;
  // clang-format off


  // clang-format on

  // Attempt to choose the maximum value for v, and ensure it is positive.
  const int max_index = val[0] > val[1] ? 0 : 1;
  const int min_index = val[0] <= val[1] ? 0 : 1;
  constexpr int32_t max_u8 = std::numeric_limits<uint8_t>::max();
  constexpr int32_t max_exp = 63;
  constexpr int32_t min_exp = -128;
  int current_index = max_index;
  if ((val[max_index] > max_u8) || (exp[max_index] > max_exp) ||
      (exp[max_index] < min_exp)) {
    current_index = min_index;
  }

  if ((val[current_index] > max_u8) || (exp[current_index] > max_exp) ||
      (exp[current_index] < min_exp)) {
    // Max and min values are invalid.
    return absl::InvalidArgumentError(absl::StrFormat(
        "The requested frequency (%.3f) cannot fit into a 2 byte register.",
        frequency));
  }

  // Try to write the value. Since not all drives support the value, we will
  // log a warning if the write fails.
  auto interpolation_value = variable_registry.GetServiceVariable(
      kInterpolationPeriodIndex, kInterpolationPeriodValueSubindex,
      bus_position);
  auto status = interpolation_value.status();
  if (status.ok()) {
    status =
        interpolation_value->Write(static_cast<uint8_t>(val[current_index]));
  }
  if (!status.ok()) {
    INTRINSIC_RT_LOG(INFO)
        << "Device at bus position " << bus_position
        << " does not support writing the interpolation value: " << status;
  }

  // Try to write the exponent. Since not all drives support the exponent,
  // we will log a warning if the write fails.
  auto interpolation_exponent = variable_registry.GetServiceVariable(
      kInterpolationPeriodIndex, kInterpolationPeriodExponentSubindex,
      bus_position);
  status = interpolation_exponent.status();
  if (status.ok()) {
    status =
        interpolation_exponent->Write(static_cast<int8_t>(exp[current_index]));
  }
  if (!status.ok()) {
    INTRINSIC_RT_LOG(INFO)
        << "Device at bus position " << bus_position
        << " does not support writing the interpolation exponent: " << status;
  }

  // Try to extend the position limit. Since not all drives support this, we
  // will log a warning if the write fails.
  auto position_limit = variable_registry.GetServiceVariable(
      kPositionLimitIndex, kPositionLimitSubindex, bus_position);
  status = position_limit.status();
  if (status.ok()) {
    status = position_limit->Write(kMaxPositionLimit);
  }
  if (!status.ok()) {
    INTRINSIC_RT_LOG(INFO) << "Device at bus position " << bus_position
                           << " does not support writing the position limit: "
                           << status;
  }

  return absl::OkStatus();
}

Ds402BusComponent::OptionalStatusBits
Ds402BusComponent::ParseOptionalStatusBits(ProcessVariableType status_word) {
  OptionalStatusBits optional_status_bits;
  const std::bitset<sizeof(ProcessVariableType) * 8> status_bits(status_word);

  // Bit 7 of the `Statusword`.
  optional_status_bits.has_warning = status_bits[kWarningBitIndex];

  // Bit 11 of the `Statusword`
  optional_status_bits.has_internal_limits_active =
      status_bits[kInternalLimitsActiveBitIndex];

  // The ignore_unexpected_state_transition_signal can be read from either the
  // `Statusword` or a generic digital input using the same logic.
  if (unexpected_state_transition_configuration_.signal_reader.has_value()) {
    // Sets the ignore_unexpected_state_transition_signal by reading a
    // variable. Does not have to be the `Statusword`, this
    // is a convenient place to check if the signal is active.
    optional_status_bits.ignore_unexpected_state_transition_signal =
        unexpected_state_transition_configuration_.signal_reader->operator()();
  }

  // Log changes to the ignore_unexpected_state_transition_signal.
  if (const bool signal =
          optional_status_bits.ignore_unexpected_state_transition_signal;
      ignore_unexpected_state_transition_signal_previous_cycle_ != signal) {
    INTRINSIC_RT_LOG(INFO)
        << "`ignore_unexpected_state_transition_signal` changed from "
        << ignore_unexpected_state_transition_signal_previous_cycle_ << " to "
        << signal;
    ignore_unexpected_state_transition_signal_previous_cycle_ = signal;
  }

  return optional_status_bits;
}

}  // namespace intrinsic::ds402
