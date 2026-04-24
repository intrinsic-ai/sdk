// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_DS402_DS402_FAKE_DEVICE_H_
#define INTRINSIC_ICON_HAL_LIB_DS402_DS402_FAKE_DEVICE_H_

#include <cstdint>
#include <optional>

#include "absl/strings/string_view.h"
#include "intrinsic/icon/hal/lib/ds402/ds402_driver.h"
#include "intrinsic/icon/hal/lib/fieldbus/fake_variable_registry.h"

namespace intrinsic::ds402 {

// A fake implementation of a Ds402 device.
//
// This class is used for testing the Ds402BusComponent. It creates fake bus
// variables that can be read by the Ds402BusComponent and internally keeps
// track of the DS402 state of the device, reacting to control words and
// updating the status word and error code accordingly. It allows to write tests
// that are very close to real-world operation, but without the need for an
// actual DS402 device. Plus, given a configurable transition delay, it will
// make the behavior of the device deterministic and predictable.
class Ds402FakeDevice {
 public:
  // Creates a new Ds402FakeDevice.
  //
  // Args:
  //   variable_registry: The fake variable registry to use for this device.
  //   transition_delay: The transition delay to use for this device. For most
  //   state transitions, the device will need to see the "right" control word
  //   for this amount of ticks before transitioning to the new state.
  //   status_word_variable_name: The name of the status word variable.
  //   control_word_variable_name: The name of the control word variable.
  //   error_code_variable_name: The name of the error code variable.
  //   modes_of_operation_variable_name: The name of the modes of operation
  //   variable.
  //   modes_of_operation_display_variable_name: The name of the modes of
  //   operation display variable.
  Ds402FakeDevice(
      intrinsic::fieldbus::FakeVariableRegistry& variable_registry,
      uint32_t transisiton_delay, absl::string_view status_word_variable_name,
      absl::string_view control_word_variable_name,
      absl::string_view error_code_variable_name,
      std::optional<absl::string_view> modes_of_operation_variable_name,
      std::optional<absl::string_view>
          modes_of_operation_display_variable_name);

  // Simulates a single tick of the device.
  // This will update the status word and error code variables based on the
  // current state of the device and the control word.
  // The device behavior is simplified but deterministic. Don't expect real
  // devices to behave exactly like this (most importantly w.r.t. timing).
  void Tick();

  // Updates the state of the device.
  // Used internally by Tick() to update the state of the device, but can also
  // be used to set the state of the device directly.
  //
  // Args:
  //   new_state: The new state to set the device to.
  //   reset_transition_delay: Whether to reset the transition delay.
  void UpdateState(Ds402State new_state, bool reset_transition_delay = false);

  // Returns the status word for the given state.
  static uint16_t GetStatusWordForState(Ds402State state);

  // Allows forcing bits in the status word to true when UpdateState() is
  // called.
  // `forced_status_word_values` is OR'ed with the status word returned by
  // GetStatusWordForState(). This means all bits that are `false` are ignored.
  void SetForcedStatusWordValues(uint16_t forced_status_word_values);

  // The transition delay for this device, applied to most state transitions
  // (except for kFault -> kSwitchedOnDisabled and kQuickStop to
  // kOperationEnabled).
  uint32_t transition_delay_;
  // The remaining transition delay for the next state transition.
  uint32_t remaining_transition_delay_;
  // The internal variables that will be read by the Ds402BusComponent through
  // the fake variable registry.
  uint16_t status_word_;
  uint16_t control_word_;
  uint16_t error_code_;
  // The current state of the device.
  Ds402State state_;
  // The last control word received by the device.
  ControlWord last_control_word_;

  // Allows forcing bits in the status word to true.
  uint16_t forced_status_word_values_ = 0x0;

  // Homing variables.
  int8_t homing_mode_;
  int32_t homing_offset_;
  uint32_t homing_search_speed_;
  uint32_t homing_creep_speed_;
  uint32_t homing_acceleration_;

  // Modes of operation variables.
  int8_t modes_of_operation_;
  int8_t modes_of_operation_display_;
};
}  // namespace intrinsic::ds402

#endif  // INTRINSIC_ICON_HAL_LIB_DS402_DS402_FAKE_DEVICE_H_
