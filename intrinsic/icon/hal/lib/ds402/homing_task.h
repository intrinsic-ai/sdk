// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_DS402_HOMING_TASK_H_
#define INTRINSIC_ICON_HAL_LIB_DS402_HOMING_TASK_H_

#include <stdbool.h>

#include <cstdint>
#include <memory>
#include <optional>
#include <string>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "intrinsic/icon/hal/command_validator.h"
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/hardware_interface_registry.h"
#include "intrinsic/icon/hal/interfaces/electrical_motor.fbs.h"
#include "intrinsic/icon/hal/lib/ds402/v1/ds402_bus_component_config.pb.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"
#include "intrinsic/icon/utils/async_function.h"
#include "intrinsic/icon/utils/bitset.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/platform/common/buffers/rt_promise.h"

namespace intrinsic::ds402 {

// Manages the homing procedure for a DS402-compliant motor drive.
// This class implements a state machine to handle the homing process, which
// involves setting parameters, changing operation modes, and monitoring the
// drive's status. It interacts with the hardware in both real-time and
// non-real-time contexts.
class HomingTask {
 public:
  using ControlWord = intrinsic::bitset<uint16_t>;
  using StatusWord = intrinsic::bitset<uint16_t>;
  // Defines the states of the homing process state machine.
  enum class InternalHomingState : uint8_t {
    kHomingNotConfigured = 0,
    kHomingInit,
    kHomingIdle,
    kHomingCommanded,
    kHomingSettingParameters,
    kHomingParametersSet,
    kHomingSettingModesOfOperation,
    kHomingModesOfOperationReached,
    kHomingAttainedRestoringModesOfOperation,
    kHomingAttained,
    kHomingError,
    kHomingInternalError,
  };

  // Asynchronous function type for setting homing parameters on the drive.
  // This is a non-real-time operation.
  using SetHomingParamsAsyncType = intrinsic::icon::AsyncNonRealtimeFunction<
      absl::StatusOr</*method=*/int8_t>(
          /*method=*/int8_t,
          /*offset=*/int32_t,
          /*search_speed=*/uint32_t,
          /*creep_speed=*/uint32_t,
          /*acceleration=*/uint32_t)>;
  // Asynchronous function type for setting the mode of operation on the drive.
  // This is a non-real-time operation.
  using SetModesOfOperationAsyncType =
      intrinsic::icon::AsyncNonRealtimeFunction<
          absl::StatusOr<intrinsic_fbs::MotorControlMode>(
              /*modes_of_operation=*/intrinsic_fbs::MotorControlMode,
              /*settle_time=*/absl::Duration,
              /*timeout=*/absl::Duration)>;

  // Object dictionary indices for homing-related parameters.
  static constexpr uint32_t kHomingMethodIndex = 0x6098;
  static constexpr uint32_t kHomingOffsetIndex = 0x607C;
  static constexpr uint32_t kHomingSpeedIndex = 0x6099;
  static constexpr uint32_t kHomingAccelerationIndex = 0x609A;
  static constexpr uint32_t kModesOfOperationIndex = 0x6060;
  static constexpr uint32_t kModesOfOperationDisplayIndex = 0x6061;

  // Bit indices within the StatusWord for homing-related flags.
  static constexpr uint32_t kHomingBitIndex = 4;
  static constexpr uint32_t kHomingErrorBitIndex = 13;
  static constexpr uint32_t kHomingAttainedBitIndex = 12;
  static constexpr uint32_t kTargetReachedBitIndex = 10;

  // Error messages.
  static constexpr absl::string_view
      kHomingRequiresOperationEnabledErrorMessage =
          "Homing requires drive to be in OperationEnabled.";
  static constexpr absl::string_view
      kHomingStillProcessingOldCommandErrorMessage =
          "Still processing old homing command.";
  // Timeout for waiting for the mode of operation to be confirmed by the drive.
  static constexpr absl::Duration kWaitForModesOfOperationDisplayTimeout =
      absl::Seconds(1);
  // Time to wait after homing has been reached before restoring the modes of
  // operation. Some drives need this before they update their position values.
  static constexpr absl::Duration
      kSetModesOfOperationSettleTimeAfterHomingAttained =
          absl::Milliseconds(100);

  ~HomingTask() = default;

  // Factory function to create a new HomingTask instance.
  // Use `bus_position` to specify the position of the device on the bus.
  // `variable_registry` is used to access both process and service variables.
  // The homing-related hardware interfaces are registered in the
  // `interface_registry`. `config` contains the
  // configuration for the homing procedure. Returns a StatusOr containing a
  // unique_ptr to the HomingTask on success, or an error status on failure.
  static absl::StatusOr<std::unique_ptr<HomingTask>> Create(
      int32_t bus_position,
      const intrinsic::fieldbus::VariableRegistry& variable_registry,
      intrinsic::icon::HardwareInterfaceRegistry& interface_registry,

      const intrinsic_proto::icon::v1::Ds402BusComponent& config);

  // Performs the real-time read operations for the homing task.
  // This method should be called in every real-time cycle. It reads the drive's
  // `status_word` and updates the internal homing state machine. Returns a
  // RealtimeStatus indicating success or failure.
  intrinsic::icon::RealtimeStatus CyclicRead(const StatusWord& status_word);

  // Performs the real-time write operations for the homing task.
  // This method should be called in every real-time cycle. It generates the
  // appropriate `control_word` based on the current homing state. Returns a
  // RealtimeStatusOr containing the modified control word on success, or an
  // error status on failure.
  intrinsic::icon::RealtimeStatusOr<ControlWord> CyclicWrite(
      const ControlWord& control_word);

 private:
  // Private constructor. Use the `Create` factory function instead.
  // Expects `device_name` to be the name of the device used for logging.
  // `set_homing_params` and `set_modes_of_operation` are a non-real-time
  // function to set homing parameters on the drive. `command_validator` is a
  // validator for the homing command. `homing_command` is a handle to the
  // homing command interface. `homing_status` is a handle to the homing status
  // interface.
  HomingTask(
      absl::string_view device_name,
      SetHomingParamsAsyncType&& set_homing_params,
      SetModesOfOperationAsyncType&& set_modes_of_operation,
      intrinsic::icon::Validator command_validator,
      intrinsic::icon::HardwareInterfaceHandle<::intrinsic_fbs::HomeCommand>
          homing_command,
      intrinsic::icon::MutableHardwareInterfaceHandle<
          ::intrinsic_fbs::HomingStatus>
          homing_status);

  HomingTask(const HomingTask&) = delete;
  HomingTask& operator=(const HomingTask&) = delete;

  HomingTask(HomingTask&& other);
  HomingTask& operator=(HomingTask&&) = delete;

  // Checks if a homing command has been issued.
  // Returns true if a homing command has been issued.
  bool IsHomingCommanded() const;

  // Checks if a new homing command has been issued since the last one.
  // Returns true if a new homing command was issued in this cycle.
  bool IsNewHomingCommand() const;

  // Updates the homing status interface based on the internal state.
  intrinsic::icon::RealtimeStatus UpdateHomingStatus();

  // The name of the device associated with this task.
  const std::string device_name_;
  // Asynchronous function to set homing parameters.
  SetHomingParamsAsyncType set_homing_params_;
  // Asynchronous function to set the mode of operation.
  SetModesOfOperationAsyncType set_modes_of_operation_;
  // Hardware interface handle for receiving homing commands.
  intrinsic::icon::HardwareInterfaceHandle<::intrinsic_fbs::HomeCommand>
      homing_command_;
  // Hardware interface handle for publishing homing status.
  intrinsic::icon::MutableHardwareInterfaceHandle<::intrinsic_fbs::HomingStatus>
      homing_status_;
  // The current state of the homing state machine.
  InternalHomingState homing_state_{InternalHomingState::kHomingNotConfigured};
  // Validator to ensure homing commands are not issued while another is active.
  intrinsic::icon::Validator command_validator_;
  // Future for the result of the `set_homing_params_` async call.
  intrinsic::RealtimeFuture<absl::StatusOr<int8_t>> set_homing_params_future_;
  // Future for the result of the `set_modes_of_operation_` async call.

  intrinsic::RealtimeFuture<absl::StatusOr<intrinsic_fbs::MotorControlMode>>
      set_modes_of_operation_future_;
  // Stores the last received homing command to detect new commands.
  std::optional<::intrinsic_fbs::HomeCommand> last_homing_command_;
};

}  // namespace intrinsic::ds402

#endif  // INTRINSIC_ICON_HAL_LIB_DS402_HOMING_TASK_H_
