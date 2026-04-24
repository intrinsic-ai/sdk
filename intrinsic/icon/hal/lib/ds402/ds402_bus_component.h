// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_DS402_DS402_BUS_COMPONENT_H_
#define INTRINSIC_ICON_HAL_LIB_DS402_DS402_BUS_COMPONENT_H_

#include <sys/types.h>

#include <cstdint>
#include <memory>
#include <optional>
#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/hal/lib/ds402/ds402_bus_component_utils.h"
#include "intrinsic/icon/hal/lib/ds402/ds402_driver.h"
#include "intrinsic/icon/hal/lib/ds402/homing_task.h"
#include "intrinsic/icon/hal/lib/ds402/v1/ds402_bus_component_config.pb.h"
#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/hal/lib/fieldbus/bus_component.h"
#include "intrinsic/icon/hal/lib/fieldbus/bus_component_factory.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"
#include "intrinsic/icon/testing/realtime_annotations.h"
#include "intrinsic/icon/utils/bitset.h"
#include "intrinsic/icon/utils/fixed_string.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/platform/common/buffers/realtime_write_queue.h"
#include "intrinsic/platform/common/buffers/rt_queue.h"
#include "intrinsic/util/thread/thread.h"

namespace intrinsic::ds402 {

// DS402 bus device.
// Controls the state of the corresponding physical device on the bus through
// the various states of the DS402 state machine.
class Ds402BusComponent : public fieldbus::BusComponent {
 public:
  using DigitalOutputsType = uint32_t;
  using ErrorInterpretationString = intrinsic::icon::FixedString<
      intrinsic::icon::RealtimeStatus::kMaxMessageLength>;

  // Result of a process variable interpretation. Passed from the non-rt thread
  // to the rt thread.
  struct ProcessVariableInterpretationResult {
    // The error code that was translated. So that the real time thread can
    // request a new translation when required.
    ProcessVariableType error_code;
    // All matching interpretations, separated by '; '. Potentially truncated.
    ErrorInterpretationString interpretation = ErrorInterpretationString("");
  };

  class ProcessVariableInterpretationProvider {
   public:
    ~ProcessVariableInterpretationProvider();
    // Creates a ProcessVariableInterpretationHelper.
    // Uses a non-rt thread to translate process variable values to fixed
    // strings. Forwards thread creation errors. Skips thread creation and
    // disables interpretation lookup if the config does not contain any error
    // code interpretations.
    static absl::StatusOr<
        std::unique_ptr<ProcessVariableInterpretationProvider>>
    Create(const intrinsic_proto::icon::v1::Ds402BusComponent& config);

    // Returns the interpretation for the given error code.
    //
    // If the current interpretation does not match the given error code, a new
    // interpretation is requested and std::nullopt is returned.
    // Once the interpretation is available, the current interpretation is
    // updated and the interpretation is returned.
    //
    // Returns std::nullopt if interpretation lookup is disabled.
    std::optional<ProcessVariableInterpretationResult>
    GetProcessVariableInterpretation(ProcessVariableType error_code)
        INTRINSIC_CHECK_REALTIME_SAFE;

   private:
    explicit ProcessVariableInterpretationProvider(
        bool interpretation_lookup_enabled);
    intrinsic::icon::RealtimeStatus RequestInterpretation(
        ProcessVariableType error_code);

    const bool interpretation_lookup_enabled_;

    intrinsic::RealtimeWriteQueue<ProcessVariableType> request_queue_;
    intrinsic::RealtimeQueue<ProcessVariableInterpretationResult> result_queue_;
    intrinsic::Thread non_rt_translation_thread_;

    std::optional<ProcessVariableInterpretationResult> current_interpretation_ =
        std::nullopt;
  };

  // The index of the digital output bit that engages/releases the brakes
  // (according the DS402 spec).
  // The osprey in mtv uses bit 18, the one in muc uses bit 0.
  static constexpr uint32_t kBrakeReleaseBitValue = (1 << 18) + 1;

  // service variable Indices for interpolation and position limits.
  static constexpr int32_t kInterpolationPeriodIndex = 0x60C2;
  static constexpr int32_t kInterpolationPeriodValueSubindex = 1;
  static constexpr int32_t kInterpolationPeriodExponentSubindex = 2;
  static constexpr int32_t kPositionLimitIndex = 0x2520;
  static constexpr int32_t kPositionLimitSubindex = 0;

  // Move constructor and assign.
  Ds402BusComponent(Ds402BusComponent&& other) = delete;
  Ds402BusComponent& operator=(Ds402BusComponent&& other) = delete;

  // Creates a Ds402BusComponent.
  // Returns
  //  * a NotFoundError if the required bus variables cannot be found.
  //  * an InvalidArgumentError if the bus variables are of incompatible types.
  //  * an InvalidArgumentError if setting the interpolation window fails.
  //  * an InvalidArgumentError if the enable_delay is negative.
  //  * an InvalidArgumentError if
  //  `distributed_clock_and_bus_shift_sync_are_enabled` is false.
  // The value of `distributed_clock_and_bus_shift_sync_are_enabled` needs to be
  // determined outside of this factory function, i.e. from the fieldbus
  // implementation.
  static absl::StatusOr<std::unique_ptr<Ds402BusComponent>> Create(
      fieldbus::DeviceInitContext& device_init_context,
      const intrinsic_proto::icon::v1::Ds402BusComponent& config,
      double frequency, bool distributed_clock_and_bus_shift_sync_are_enabled);

  // Reads and updates the state of any input hardware interface.
  // Returns the status of the request or an error if...
  //  * it fails to decode the DS402 status word.
  //  * the device reports an error_code unless request_type is kClearFaults.
  //  * the device falls out of OperationEnabled unless request_type is
  //    kClearFaults.
  intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus> CyclicRead(
      fieldbus::RequestType request_type) override;

  // Writes the values of any output hardware interface to the bus.
  // Returns the status of the request or an error if it fails to derive the
  // next control word.
  intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus> CyclicWrite(
      fieldbus::RequestType request_type) override;

  // Returns empty `ProcessVariableInterpretationResult.interpretation` if no
  // interpretations are defined, or match.
  static ProcessVariableInterpretationResult
  GetProcessVariableInterpretationFromConfig(
      ProcessVariableType process_variable_value,
      const intrinsic_proto::icon::v1::Ds402BusComponent&
          config_with_interpretations);

 private:
  // Constructor.
  // `device_name`: The name of the device, used for logging.
  // `status_word`: The status word ProcessVariable.
  // `error_code`: The optional error code on ProcessVariable.
  // `control_word`: The control word on ProcessVariable.
  // `digital_outputs`: The optional digital outputs on the bus, used for brake
  //                    release.
  // `treat_internal_limits_active_as_error`: When true, the bit
  //                                          internal_limits_active triggers a
  //                                          fault.
  // `unexpected_state_transition_configuration`: Configuration of how the
  //                        device should handle unexpected state transitions
  //                        from the `OperationEnabled` state.
  // `enable_delay_cycles`: Delay in cycles before enabling the device during
  //                        EnableMotion.
  // `initial_goal_state`: The goal state of the device when it is turned on or
  //                       disabled.
  // `enabled_goal_state`: The goal state for the device when it is enabled
  //                       (during EnableMotion).
  // `process_variable_interpretation`: Helper to translate the error_code to a
  // string. `homing_command`: The homing command interface.
  Ds402BusComponent(absl::string_view device_name,
                    fieldbus::ProcessVariable status_word,
                    std::optional<fieldbus::ProcessVariable> error_code,
                    fieldbus::ProcessVariable control_word,
                    std::optional<fieldbus::ProcessVariable> digital_outputs,
                    bool treat_internal_limits_active_as_error,
                    UnexpectedStateTransitionConfiguration
                        unexpected_state_transition_configuration,
                    int64_t enable_delay_cycles, Ds402State initial_goal_state,
                    Ds402State enabled_goal_state,
                    std::unique_ptr<ProcessVariableInterpretationProvider>
                        process_variable_interpretation,
                    std::unique_ptr<HomingTask> homing_task,
                    uint32_t brake_release_bit_value);
  struct Ds402DeviceState {
    ProcessVariableType status_word;
    ProcessVariableType error_code;
    intrinsic::bitset<uint16_t> control_word;
    DigitalOutputsType digital_outputs;
  };

  // Status bits of the `Statusword` that may be populated by a DS402 device.
  struct OptionalStatusBits {
    // True when the optional bit 7 of the `Statusword` is set.
    // It indicates that an internal state e.g. temperature, is close to the
    // limit. The cause may be found by reading the fault code parameter.
    bool has_warning = false;
    // True when the optional bit 11 of the `Statusword` is set.
    // Indicates that the internal limits are exceeded and the target/set-point
    // values can't be reached. This can be due to e.g. hardware position
    // switches, current limiter or thermal overload.
    bool has_internal_limits_active = false;
    // Runtime signal indicating that the device may transition into an
    // unexpected (non OE) state.
    // Can be configured to be a bit in the `Statusword` or a generic digital
    // input.
    // `True` when the signal is active. Some drives use bit 8 of the
    // `Statusword`.
    bool ignore_unexpected_state_transition_signal = false;
  };

  // Configures the interpolation window on the bus device via service variable
  // commands.
  static absl::Status SetInterpolationWindow(
      const intrinsic::fieldbus::VariableRegistry& variable_registry,
      int32_t bus_position, double frequency);

  OptionalStatusBits ParseOptionalStatusBits(ProcessVariableType status_word);

  // The device state.
  Ds402DeviceState device_state_;

  // The goal state for the device.
  Ds402State goal_state_;
  // Holds the DS402 state based on the latest status word.
  Ds402State ds402_state_;
  // Holds the previous state of the device.
  Ds402State prev_state_;

  const Ds402State initial_goal_state_;
  const Ds402State enabled_goal_state_;

  // Indicates a fault of the device.
  bool is_faulted_;

  // When true, the signal was active in the previous cycle.
  bool ignore_unexpected_state_transition_signal_previous_cycle_ = false;

  // Contains the fault message.
  intrinsic::icon::FixedString<
      intrinsic::icon::RealtimeStatus::kMaxMessageLength>
      fault_message_;

  // The device name.
  std::string device_name_;
  // When true, the bit internal_limits_active triggers a fault.
  bool treat_internal_limits_active_as_error_;
  // Configuration and readers for handling unexpected state transitions.
  UnexpectedStateTransitionConfiguration
      unexpected_state_transition_configuration_;

  // The status word on the bus.
  fieldbus::ProcessVariable status_word_;
  // The error code on the bus.
  std::optional<fieldbus::ProcessVariable> error_code_;
  // The control word on the bus.
  fieldbus::ProcessVariable control_word_;
  // The digital outputs on the bus.
  std::optional<fieldbus::ProcessVariable> digital_outputs_;
  // Hold the request_status from `CyclicRead`.
  fieldbus::RequestStatus current_request_status_ =
      fieldbus::RequestStatus::kDone;
  // Enable delay in cycles.
  const int64_t enable_delay_cycles_;
  // Remaining enable delay in cycles.
  int64_t remaining_enable_delay_cycles_;

  // Helper to translate the error_code to a string.
  std::unique_ptr<ProcessVariableInterpretationProvider>
      process_variable_interpretation_provider_;
  // Helper to perform the homing task.
  std::unique_ptr<HomingTask> homing_task_;

  // The bit value of the digital output bit that engages/releases the brakes
  // (according the DS402 spec).
  // Defaults to bits 18 and 0, thus 262145 (decimal), 0x40001 (hex).
  const uint32_t brake_release_bit_value_;
};
}  // namespace intrinsic::ds402

// Registers the Ds402BusComponent and its config type with the EtherCAT bus
// device factory. Allows constructing the device from its config type via
// `CreateBusComponentFromConfig`.
REGISTER_BUS_COMPONENT(ds402::Ds402BusComponent,
                       intrinsic_proto::icon::v1::Ds402BusComponent);

#endif  // INTRINSIC_ICON_HAL_LIB_DS402_DS402_BUS_COMPONENT_H_
