// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/ds402/ds402_fake_device.h"

#include <cstdint>
#include <optional>
#include <type_traits>

#include "absl/strings/string_view.h"
#include "intrinsic/icon/hal/lib/ds402/ds402_driver.h"
#include "intrinsic/icon/hal/lib/ds402/homing_task.h"
#include "intrinsic/icon/hal/lib/fieldbus/fake_variable_registry.h"

namespace intrinsic::ds402 {

Ds402FakeDevice::Ds402FakeDevice(
    intrinsic::fieldbus::FakeVariableRegistry& variable_registry,
    uint32_t transisiton_delay, absl::string_view status_word_variable_name,
    absl::string_view control_word_variable_name,
    absl::string_view error_code_variable_name,
    std::optional<absl::string_view> modes_of_operation_variable_name,
    std::optional<absl::string_view> modes_of_operation_display_variable_name)
    : transition_delay_(transisiton_delay),
      remaining_transition_delay_(transition_delay_),
      status_word_(0),
      control_word_(0),
      error_code_(0),
      state_(Ds402State::kNotReadyToSwitchOn),
      last_control_word_(ControlWord::kDisableVoltage),
      homing_mode_(0),
      homing_offset_(0),
      homing_search_speed_(0),
      homing_creep_speed_(0),
      homing_acceleration_(0),
      modes_of_operation_(0),
      modes_of_operation_display_(0) {
  variable_registry.AddInputProcessVariable(status_word_variable_name,
                                            &status_word_);
  variable_registry.AddInputProcessVariable(error_code_variable_name,
                                            &error_code_);
  variable_registry.AddOutputProcessVariable(control_word_variable_name,
                                             &control_word_);
  variable_registry.AddServiceVariable(HomingTask::kHomingMethodIndex,
                                       /*subindex=*/0, /*bus_position=*/0,
                                       &homing_mode_);
  variable_registry.AddServiceVariable(HomingTask::kHomingOffsetIndex,
                                       /*subindex=*/0, /*bus_position=*/0,
                                       &homing_offset_);
  variable_registry.AddServiceVariable(HomingTask::kHomingSpeedIndex,
                                       /*subindex=*/1, /*bus_position=*/0,
                                       &homing_search_speed_);
  variable_registry.AddServiceVariable(HomingTask::kHomingSpeedIndex,
                                       /*subindex=*/2, /*bus_position=*/0,
                                       &homing_creep_speed_);
  variable_registry.AddServiceVariable(HomingTask::kHomingAccelerationIndex,
                                       /*subindex=*/0, /*bus_position=*/0,
                                       &homing_acceleration_);
  if (modes_of_operation_variable_name.has_value()) {
    variable_registry.AddOutputProcessVariable(
        *modes_of_operation_variable_name, &modes_of_operation_);
  }
  if (modes_of_operation_display_variable_name.has_value()) {
    variable_registry.AddInputProcessVariable(
        *modes_of_operation_display_variable_name,
        &modes_of_operation_display_);
  }
  variable_registry.AddServiceVariable(HomingTask::kModesOfOperationIndex,
                                       /*subindex=*/0, /*bus_position=*/0,
                                       &modes_of_operation_);
  variable_registry.AddServiceVariable(
      HomingTask::kModesOfOperationDisplayIndex,
      /*subindex=*/0, /*bus_position=*/0, &modes_of_operation_display_);
  Tick();
  remaining_transition_delay_ = transition_delay_;
}

void Ds402FakeDevice::SetForcedStatusWordValues(
    uint16_t forced_status_word_values) {
  forced_status_word_values_ = forced_status_word_values;
}

void Ds402FakeDevice::Tick() {
  ControlWord control_word = ControlWord(control_word_);

  // Reset the transition delay, if the control word changed.
  if (control_word != last_control_word_) {
    remaining_transition_delay_ = transition_delay_;
  }
  auto new_state = state_;

  switch (state_) {
    case Ds402State::kNotReadyToSwitchOn:
      // Regardless of the control word, we will transition to
      // kSwitchOnDisabled.
      if (--remaining_transition_delay_ <= 0) {
        new_state = Ds402State::kSwitchOnDisabled;
      }
      break;
    case Ds402State::kSwitchOnDisabled:
      switch (control_word) {
        case ControlWord::kShutdown:
          // kShutdown will transition to kReadyToSwitchOn.
          if (--remaining_transition_delay_ <= 0) {
            new_state = Ds402State::kReadyToSwitchOn;
          }
          break;
        case ControlWord::kDisableVoltage:
          // kDisableVoltage will let us remain in kSwitchOnDisabled.
          break;
        case ControlWord::kSwitchOn:  // same as kDisableOperation
          [[fallthrough]];
        case ControlWord::kEnableOperation:  // same as
                                             // kEnableOperationAfterQuickStop
          [[fallthrough]];
        case ControlWord::kQuickStop:
          [[fallthrough]];
        case ControlWord::kFaultReset:
          // All other control words will transition to kFaultReactionActive
          // (immediately).
          new_state = Ds402State::kFaultReactionActive;
          break;
      }
      break;
    case Ds402State::kReadyToSwitchOn:
      switch (control_word) {
        case ControlWord::kSwitchOn:  // same as kDisableOperation
          // kSwitchOn/kDisableOperation will transition to kSwitchedOn.
          if (--remaining_transition_delay_ <= 0) {
            new_state = Ds402State::kSwitchedOn;
          }
          break;
        case ControlWord::kDisableVoltage:
          // kDisableVoltage will transition to kSwitchOnDisabled.
          if (--remaining_transition_delay_ <= 0) {
            new_state = Ds402State::kSwitchOnDisabled;
          }
          break;
        case ControlWord::kShutdown:
          // kShutdown will let us remain in kReadyToSwitchOn.
          break;
        case ControlWord::kEnableOperation:  // same as
                                             // kEnableOperationAfterQuickStop
          [[fallthrough]];
        case ControlWord::kQuickStop:
          [[fallthrough]];
        case ControlWord::kFaultReset:
          // All other control words will transition to kFaultReactionActive
          // (immediately).
          new_state = Ds402State::kFaultReactionActive;
          break;
      }
      break;
    case Ds402State::kSwitchedOn:
      switch (control_word) {
        case ControlWord::kEnableOperation:  // same as
                                             // kEnableOperationAfterQuickStop
          // kEnableOperation(AfterQuickStop) will transition to
          // kOperationEnabled.
          if (--remaining_transition_delay_ <= 0) {
            new_state = Ds402State::kOperationEnabled;
          }
          break;
        case ControlWord::kDisableVoltage:
          // kDisableVoltage will transition to kSwitchOnDisabled.
          if (--remaining_transition_delay_ <= 0) {
            new_state = Ds402State::kSwitchOnDisabled;
          }
          break;
        case ControlWord::kShutdown:
          // kShutdown will transition to kReadyToSwitchOn.
          if (--remaining_transition_delay_ <= 0) {
            new_state = Ds402State::kReadyToSwitchOn;
          }
          break;
        case ControlWord::kSwitchOn:  // same as kDisableOperation
          // kSwitchOn/kDisableOperation will let us remain in to kSwitchedOn.
          break;
        case ControlWord::kQuickStop:
          [[fallthrough]];
        case ControlWord::kFaultReset:
          // All other control words will transition to kFaultReactionActive
          // (immediately).
          new_state = Ds402State::kFaultReactionActive;
          break;
      }
      break;
    case Ds402State::kOperationEnabled:
      switch (control_word) {
        case ControlWord::kEnableOperation:  // same as
                                             // kEnableOperationAfterQuickStop
          // kEnableOperation will let us remain in to kOperationEnabled.
          break;
        case ControlWord::kSwitchOn:  // same as kDisableOperation
          // kSwitchOn/kDisableOperation will transition to kSwitchedOn.
          if (--remaining_transition_delay_ <= 0) {
            new_state = Ds402State::kSwitchedOn;
          }
          break;
        case ControlWord::kShutdown:
          // kShutdown will transition to kReadyToSwitchOn.
          if (--remaining_transition_delay_ <= 0) {
            new_state = Ds402State::kReadyToSwitchOn;
          }
          break;
        case ControlWord::kDisableVoltage:
          // kDisableVoltage will transition to kSwitchOnDisabled.
          if (--remaining_transition_delay_ <= 0) {
            new_state = Ds402State::kSwitchOnDisabled;
          }
          break;
        case ControlWord::kQuickStop:
          // kQuickStop will transition to kQuickStopActive.
          if (--remaining_transition_delay_ <= 0) {
            new_state = Ds402State::kQuickStopActive;
          }
          break;
        case ControlWord::kFaultReset:
          // All other control words will transition to kFaultReactionActive
          // (immediately).
          new_state = Ds402State::kFaultReactionActive;
          break;
      }
      break;
    case Ds402State::kFault:
      switch (control_word) {
        case ControlWord::kFaultReset:
          // kFaultReset will transition to kSwitchOnDisabled (immediately).
          new_state = Ds402State::kSwitchOnDisabled;
          break;
        case ControlWord::kShutdown:
          [[fallthrough]];
        case ControlWord::kSwitchOn:  // same as kDisableOperation
          [[fallthrough]];
        case ControlWord::kEnableOperation:  // same as
                                             // kEnableOperationAfterQuickStop
          [[fallthrough]];
        case ControlWord::kDisableVoltage:
          [[fallthrough]];
        case ControlWord::kQuickStop:
          // All other control words will let us remain in kFault
          break;
      }
      break;
    case Ds402State::kFaultReactionActive:
      // Regardless of the control word, we will transition to kFault.
      if (--remaining_transition_delay_ <= 0) {
        new_state = Ds402State::kFault;
      }
      break;
    case Ds402State::kQuickStopActive:
      switch (control_word) {
        case ControlWord::kEnableOperation:  // same as
                                             // kEnableOperationAfterQuickStop
          // kEnableOperation will transition to kOperationEnabled
          // (immediately).
          new_state = Ds402State::kSwitchOnDisabled;
          break;
        case ControlWord::kDisableVoltage:
          // kDisableVoltage will transition to kSwitchOnDisabled.
          if (--remaining_transition_delay_ <= 0) {
            new_state = Ds402State::kSwitchOnDisabled;
          }
          break;
        case ControlWord::kShutdown:
          [[fallthrough]];
        case ControlWord::kSwitchOn:  // same as kDisableOperation
          [[fallthrough]];
        case ControlWord::kFaultReset:
          [[fallthrough]];
        case ControlWord::kQuickStop:
          // All other control words will transition to kFaultReactionActive
          // (immediately).
          new_state = Ds402State::kFaultReactionActive;
          break;
      }
      break;
  }

  // Reset the transition delay, if the state changed.
  UpdateState(new_state, /*reset_transition_delay=*/state_ != new_state);
}

void Ds402FakeDevice::UpdateState(Ds402State new_state,
                                  bool reset_transition_delay) {
  remaining_transition_delay_ =
      reset_transition_delay ? transition_delay_ : remaining_transition_delay_;
  last_control_word_ = ControlWord(control_word_);
  state_ = new_state;
  if (state_ == Ds402State::kFaultReactionActive ||
      state_ == Ds402State::kFault) {
    error_code_ = 1;
  } else {
    error_code_ = 0;
  }
  status_word_ = GetStatusWordForState(state_) | forced_status_word_values_;
}

/*static*/
uint16_t Ds402FakeDevice::GetStatusWordForState(Ds402State state) {
  switch (state) {
    case Ds402State::kNotReadyToSwitchOn:
      return std::underlying_type_t<StateEncodingStatusWord>(
          StateEncodingStatusWord::kNotReadyToSwitchOn0);
    case Ds402State::kSwitchOnDisabled:
      return std::underlying_type_t<StateEncodingStatusWord>(
          StateEncodingStatusWord::kSwitchOnDisabled0);
    case Ds402State::kReadyToSwitchOn:
      return std::underlying_type_t<StateEncodingStatusWord>(
          StateEncodingStatusWord::kReadyToSwitchOn0);
    case Ds402State::kSwitchedOn:
      return std::underlying_type_t<StateEncodingStatusWord>(
          StateEncodingStatusWord::kSwitchedOn0);
    case Ds402State::kOperationEnabled:
      return std::underlying_type_t<StateEncodingStatusWord>(
          StateEncodingStatusWord::kOperationEnabled0);
    case Ds402State::kFault:
      return std::underlying_type_t<StateEncodingStatusWord>(
          StateEncodingStatusWord::kFault0);
    case Ds402State::kFaultReactionActive:
      return std::underlying_type_t<StateEncodingStatusWord>(
          StateEncodingStatusWord::kFaultReactionActive0);
    case Ds402State::kQuickStopActive:
      return std::underlying_type_t<StateEncodingStatusWord>(
          StateEncodingStatusWord::kQuickStopActive0);
  }
}

}  // namespace intrinsic::ds402
