// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CONTROL_SAFETY_SAFETY_MESSAGES_UTILS_H_
#define INTRINSIC_ICON_CONTROL_SAFETY_SAFETY_MESSAGES_UTILS_H_

#include <bitset>
#include <cstdint>
#include <string>
#include <type_traits>

#include "absl/strings/string_view.h"
#include "flatbuffers/detached_buffer.h"
#include "intrinsic/icon/control/safety/extern/safety_status.fbs.h"
#include "intrinsic/icon/control/safety/safety_messages.fbs.h"

namespace intrinsic_fbs {

// Build a SafetyStatusMessage.
flatbuffers::DetachedBuffer BuildSafetyStatusMessage(
    ModeOfSafeOperation mode_of_safe_operation = ModeOfSafeOperation::UNKNOWN,
    ButtonStatus estop_button_status = ButtonStatus::UNKNOWN,
    ButtonStatus enable_button_status = ButtonStatus::UNKNOWN,
    RequestedBehavior requested_behavior = RequestedBehavior::UNKNOWN);

template <typename EnumType>
constexpr auto AsIndex(const EnumType value) ->
    typename std::underlying_type_t<EnumType> {
  return static_cast<typename std::underlying_type_t<EnumType>>(value);
}

void SetSafetyStatusMessage(
    ::intrinsic_fbs::ModeOfSafeOperation mode_of_safe_operation,
    ::intrinsic_fbs::ButtonStatus estop_button_status,
    ::intrinsic_fbs::ButtonStatus enable_button_status,
    ::intrinsic_fbs::RequestedBehavior requested_behavior,
    ::intrinsic_fbs::SafetyStatusMessage& message);

// Extract ModeOfSafeOperation from safety inputs.
// The safety inputs are expected to follow the order as in
// intrinsic_fbs::SafetyStatusBit.
ModeOfSafeOperation ExtractModeOfSafeOperation(
    const std::bitset<8>& safety_inputs);

// Extract the status of the e-stop button from safety inputs.
// The safety inputs are expected to follow the order as in
// intrinsic_fbs::SafetyStatusBit.
ButtonStatus ExtractEStopButtonStatus(const std::bitset<8>& safety_inputs);

// Extract the status of the enable button from safety inputs.
// The safety inputs are expected to follow the order as in
// intrinsic_fbs::SafetyStatusBit.
ButtonStatus ExtractEnableButtonStatus(const std::bitset<8>& safety_inputs);

// Extract the requested behavior from safety inputs.
// The safety inputs are expected to follow the order as in
// intrinsic_fbs::SafetyStatusBit.
RequestedBehavior ExtractRequestedBehavior(const std::bitset<8>& safety_inputs);

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_CONTROL_SAFETY_SAFETY_MESSAGES_UTILS_H_
