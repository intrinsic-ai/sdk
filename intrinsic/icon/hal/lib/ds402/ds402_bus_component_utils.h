// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_DS402_DS402_BUS_DEVICE_UTILS_H_
#define INTRINSIC_ICON_HAL_LIB_DS402_DS402_BUS_DEVICE_UTILS_H_

#include <cstdint>
#include <optional>
#include <string_view>

#include "absl/functional/any_invocable.h"
#include "absl/status/statusor.h"
#include "intrinsic/icon/hal/lib/ds402/v1/ds402_bus_component_config.pb.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"

namespace intrinsic::ds402 {

// The type of the bus variable used for the DS402 PDOs.
// E.g. status word, or control word.
using ProcessVariableType = uint16_t;

// Default values mean that the device will fault on unexpected state
// transitions.
struct UnexpectedStateTransitionConfiguration {
  intrinsic_proto::icon::v1::Ds402UnexpectedStateTransitionConfiguration::
      Behavior behavior = intrinsic_proto::icon::v1::
          Ds402UnexpectedStateTransitionConfiguration::BEHAVIOR_FAULT;

  // Reads the Signal that indicates that the device may transition into an
  // unexpected (non OE) state. It is the source of the signal for
  // `BEHAVIOR_IGNORE_ON_SIGNAL`.
  //  Returns `true` when the Signal is active.
  std::optional<absl::AnyInvocable<bool()>> signal_reader = std::nullopt;
};

// Parses the `unexpected_state_transition_configuration` from the
// `Ds402BusComponent` config.
//
// Returns InvalidArgumentError when the `behavior` is not supported.
// Forwards errors from `ParseUnexpectedStateTransitionSignal` when the
// `behavior` is `BEHAVIOR_IGNORE_ON_SIGNAL`.
absl::StatusOr<UnexpectedStateTransitionConfiguration>
ParseUnexpectedStateTransitionConfiguration(
    const intrinsic::fieldbus::VariableRegistry& variable_registry,
    const intrinsic_proto::icon::v1::Ds402BusComponent& config);

// Parses and checks the signal for unexpected state transitions.
//
// Returns InvalidArgumentError when the `behavior` is not
// `BEHAVIOR_IGNORE_ON_SIGNAL`.
// Returns InvalidArgumentError when neither `status_word_bit_index` nor
// `digital_input_pattern` are set.
// Returns InvalidArgumentError when the `status_word_bit_index` is invalid.
//    Valid range is bits 8 - 15 (the Manufacturer-specific range).
//    The bits 0 - 7 of the `Statusword` are defined by the DS402 standard.
// Forwards errors from `ValidateValuePattern` when parsing a
// `digital_input_pattern`.
absl::StatusOr<UnexpectedStateTransitionConfiguration>
ParseUnexpectedStateTransitionSignal(
    const intrinsic::fieldbus::VariableRegistry& variable_registry,
    std::string_view signal_variable_name,
    const intrinsic_proto::icon::v1::
        Ds402UnexpectedStateTransitionConfiguration&
            unexpected_state_transition_config);

}  // namespace intrinsic::ds402

#endif  // INTRINSIC_ICON_HAL_LIB_DS402_DS402_BUS_DEVICE_UTILS_H_
