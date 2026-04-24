// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/ds402/ds402_driver.h"

#include <cstdint>

#include "absl/strings/string_view.h"
#include "intrinsic/icon/utils/bitset.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_or.h"

namespace intrinsic::ds402 {

// Status word bits to determine the active state from a status word.
// clang-format off
enum class StatusWordBit : uint16_t {
  kReadyToSwitchOn =  0b00000001,  // =  1
  kSwitchedOn =       0b00000010,  // =  2
  kOperationEnabled = 0b00000100,  // =  4
  kFault =            0b00001000,  // =  8
  kQuickStop =        0b00100000,  // = 32
  kSwitchOnDisabled = 0b01000000   // = 64
};
// clang-format on

intrinsic::icon::RealtimeStatusOr<Ds402State> ToDs402State(
    uint16_t status_word) {
  // Clear irrelevant bits in the status word.
  // Only the last 7 bits are relevant to extract the state.
  switch (StateEncodingStatusWord(status_word & 0b1111111)) {
    case StateEncodingStatusWord::kNotReadyToSwitchOn0:
    case StateEncodingStatusWord::kNotReadyToSwitchOn1:
    case StateEncodingStatusWord::kNotReadyToSwitchOn2:
    case StateEncodingStatusWord::kNotReadyToSwitchOn3:
      return Ds402State::kNotReadyToSwitchOn;
    case StateEncodingStatusWord::kSwitchOnDisabled0:
    case StateEncodingStatusWord::kSwitchOnDisabled1:
    case StateEncodingStatusWord::kSwitchOnDisabled2:
    case StateEncodingStatusWord::kSwitchOnDisabled3:
      return Ds402State::kSwitchOnDisabled;
    case StateEncodingStatusWord::kReadyToSwitchOn0:
    case StateEncodingStatusWord::kReadyToSwitchOn1:
      return Ds402State::kReadyToSwitchOn;
    case StateEncodingStatusWord::kSwitchedOn0:
    case StateEncodingStatusWord::kSwitchedOn1:
      return Ds402State::kSwitchedOn;
    case StateEncodingStatusWord::kOperationEnabled0:
    case StateEncodingStatusWord::kOperationEnabled1:
      return Ds402State::kOperationEnabled;
    case StateEncodingStatusWord::kFault0:
    case StateEncodingStatusWord::kFault1:
    case StateEncodingStatusWord::kFault2:
    case StateEncodingStatusWord::kFault3:
      return Ds402State::kFault;
    case StateEncodingStatusWord::kFaultReactionActive0:
    case StateEncodingStatusWord::kFaultReactionActive1:
    case StateEncodingStatusWord::kFaultReactionActive2:
    case StateEncodingStatusWord::kFaultReactionActive3:
      return Ds402State::kFaultReactionActive;
    case StateEncodingStatusWord::kQuickStopActive0:
    case StateEncodingStatusWord::kQuickStopActive1:
      return Ds402State::kQuickStopActive;
      // No default: StateEncodingStatusWord is guaranteed to not be extended.
  }
  return InvalidArgumentError(intrinsic::icon::RealtimeStatus::StrCat(
      "Failed to get state from status word: ", status_word));
}

bool IsOperationEnabled(intrinsic::bitset<uint16_t> status_word) {
  bool is_fault = (status_word & FromEnum(StatusWordBit::kFault)).any();
  if (is_fault) {
    return false;
  }
  return (
      ((status_word & FromEnum(StateEncodingStatusWord::kOperationEnabled0)) ==
       FromEnum(StateEncodingStatusWord::kOperationEnabled0)) ||
      ((status_word & FromEnum(StateEncodingStatusWord::kOperationEnabled1)) ==
       FromEnum(StateEncodingStatusWord::kOperationEnabled1)));
}

intrinsic::icon::RealtimeStatusOr<ControlWord> GetNextControlWord(
    Ds402State prev_state, Ds402State curr_state, Ds402State goal_state,
    ControlWord prev_control_word, FaultHandling fault_handling) {
  // Keep sending the previous command, if we're already at the goal state.
  if (curr_state == goal_state) {
    // Resending `kFaultReset` might re-trigger a fault handling reaction on the
    // device.
    if (prev_control_word == ControlWord::kFaultReset) {
      return ControlWord::kDisableVoltage;
    }
    return prev_control_word;
  }

  // Device is not ready yet or still in transition.
  // Keep sending the previous command.
  if ((curr_state == Ds402State::kNotReadyToSwitchOn) ||
      (curr_state == Ds402State::kFault &&
       fault_handling == FaultHandling::kPreserve) ||
      (curr_state == Ds402State::kFaultReactionActive)) {
    return prev_control_word;
  }

  // When clearing a fault, we need to send disable voltage first, before
  // sending the reset command, to ensure a rising edge on the reset signal.
  // This way we basically toggle between `kDisableVoltage` and `kFaultReset`
  // each cycle, which is good, since we don't know whether the bus is even in
  // OP and thus the device reading.
  ControlWord fault_reset_command =
      (prev_control_word == ControlWord::kDisableVoltage)
          ? ControlWord::kFaultReset
          : ControlWord::kDisableVoltage;

  switch (goal_state) {
      // Not using exhaustive switch statements, since non-reachable goal states
      // and `curr_state == goal_state` are already handled above.
    case Ds402State::kNotReadyToSwitchOn:
    case Ds402State::kFault:
    case Ds402State::kFaultReactionActive:
      return InvalidArgumentError(intrinsic::icon::RealtimeStatus::StrCat(
          "Cannot actively transition to state: ", ToString(goal_state)));
    case Ds402State::kSwitchOnDisabled:
      switch (curr_state) {
        case Ds402State::kReadyToSwitchOn:
          return ControlWord::kDisableVoltage;
        case Ds402State::kSwitchedOn:
          return ControlWord::kDisableVoltage;
        case Ds402State::kOperationEnabled:
          return ControlWord::kQuickStop;
        case Ds402State::kFault:
          return fault_reset_command;
        case Ds402State::kQuickStopActive:
          return ControlWord::kDisableVoltage;
        default:
          break;
      }
      break;
    case Ds402State::kReadyToSwitchOn:
      switch (curr_state) {
        case Ds402State::kSwitchOnDisabled:
          return ControlWord::kShutdown;
        case Ds402State::kSwitchedOn:
          return ControlWord::kShutdown;
        case Ds402State::kOperationEnabled:
          return ControlWord::kQuickStop;
        case Ds402State::kFault:
          return fault_reset_command;
        case Ds402State::kQuickStopActive:
          return ControlWord::kDisableVoltage;
        default:
          break;
      }
      break;
    case Ds402State::kSwitchedOn:
      switch (curr_state) {
        case Ds402State::kSwitchOnDisabled:
          return ControlWord::kShutdown;
        case Ds402State::kReadyToSwitchOn:
          return ControlWord::kSwitchOn;
        case Ds402State::kOperationEnabled:
          return ControlWord::kQuickStop;
        case Ds402State::kFault:
          return fault_reset_command;
        case Ds402State::kQuickStopActive:
          return ControlWord::kDisableVoltage;
        default:
          break;
      }
      break;
    case Ds402State::kOperationEnabled:
      switch (curr_state) {
        case Ds402State::kSwitchOnDisabled:
          return ControlWord::kShutdown;
        case Ds402State::kReadyToSwitchOn:
          return ControlWord::kSwitchOn;
        case Ds402State::kSwitchedOn:
          return ControlWord::kEnableOperation;
        case Ds402State::kFault:
          return fault_reset_command;
        case Ds402State::kQuickStopActive:
          return ControlWord::kEnableOperationAfterQuickStop;
        default:
          break;
      }
      break;
    case Ds402State::kQuickStopActive:
      switch (curr_state) {
        case Ds402State::kSwitchOnDisabled:
          return ControlWord::kShutdown;
        case Ds402State::kReadyToSwitchOn:
          return ControlWord::kSwitchOn;
        case Ds402State::kSwitchedOn:
          return ControlWord::kEnableOperation;
        case Ds402State::kFault:
          return fault_reset_command;
        case Ds402State::kOperationEnabled:
          return ControlWord::kQuickStop;
        default:
          break;
      }
      break;
    default:
      return InternalError(intrinsic::icon::RealtimeStatus::StrCat(
          "Requested unknown DS402 goal state: ", ToString(goal_state)));
  }

  return InternalError(intrinsic::icon::RealtimeStatus::StrCat(
      "Missing DS402 state transition: ", ToString(curr_state), " -> ",
      ToString(goal_state)));
}

absl::string_view ToString(Ds402State state) {
  switch (state) {
    case Ds402State::kNotReadyToSwitchOn:
      return "NotReadyToSwitchOn";
    case Ds402State::kSwitchOnDisabled:
      return "SwitchOnDisabled";
    case Ds402State::kReadyToSwitchOn:
      return "ReadyToSwitchOn";
    case Ds402State::kSwitchedOn:
      return "SwitchedOn";
    case Ds402State::kOperationEnabled:
      return "OperationEnabled";
    case Ds402State::kFault:
      return "Fault";
    case Ds402State::kFaultReactionActive:
      return "FaultReactionActive";
    case Ds402State::kQuickStopActive:
      return "QuickStopActive";
  }
  return "Unknown";
}

absl::string_view ToString(ControlWord control_word) {
  switch (control_word) {
    case ControlWord::kShutdown:
      return "Shutdown";
    case ControlWord::kSwitchOn:  // same as kDisableOperation
      return "SwitchOn/DisableOperation";
    case ControlWord::kDisableVoltage:
      return "DisableVoltage";
    case ControlWord::kEnableOperation:  // same as
                                         // kEnableOperationAfterQuickStop
      return "EnableOperation/AfterQuickStop";
    case ControlWord::kQuickStop:
      return "QuickStop";
    case ControlWord::kFaultReset:
      return "FaultReset";
  }
  return "Unknown";
}

}  // namespace intrinsic::ds402
