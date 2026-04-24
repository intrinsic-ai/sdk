// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/ds402/ds402_bus_component_utils.h"

#include <cstddef>
#include <cstdint>
#include <memory>
#include <optional>
#include <string_view>
#include <utility>

#include "absl/functional/any_invocable.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "intrinsic/icon/hal/lib/adio/v1/adio_bus_component_config.pb.h"
#include "intrinsic/icon/hal/lib/ds402/v1/ds402_bus_component_config.pb.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/v1/value_parsing.pb.h"
#include "intrinsic/icon/hal/lib/fieldbus/value_pattern_utils.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::ds402 {

using ::intrinsic_proto::icon::v1::Ds402UnexpectedStateTransitionConfiguration;

absl::StatusOr<UnexpectedStateTransitionConfiguration>
ParseUnexpectedStateTransitionSignal(
    const intrinsic::fieldbus::VariableRegistry& variable_registry,
    std::string_view signal_variable_name,
    const intrinsic_proto::icon::v1::
        Ds402UnexpectedStateTransitionConfiguration&
            unexpected_state_transition_config) {
  absl::AnyInvocable<bool()> input_reader;

  if (unexpected_state_transition_config.behavior() !=
      Ds402UnexpectedStateTransitionConfiguration::BEHAVIOR_IGNORE_ON_SIGNAL) {
    return absl::InvalidArgumentError(
        "ParseUnexpectedStateTransitionSignal should only be called when the "
        "behavior is `BEHAVIOR_IGNORE_ON_SIGNAL`.");
  }

  if (!unexpected_state_transition_config.has_status_word_bit_index() &&
      !unexpected_state_transition_config.has_digital_input_pattern()) {
    // ParseUnexpectedStateTransitionSignal is only called when the behavior is
    // `BEHAVIOR_IGNORE_ON_SIGNAL`.
    return absl::InvalidArgumentError(
        "`BEHAVIOR_IGNORE_ON_SIGNAL` requires 'status_word_bit_index' or "
        "'digital_input_pattern'.");
  }

  UnexpectedStateTransitionConfiguration configuration;
  configuration.behavior = unexpected_state_transition_config.behavior();

  // Mutable copy of the unexpected_state_transition_config so that we can add
  // the `statusword` as a `digital_input_pattern`.
  Ds402UnexpectedStateTransitionConfiguration config =
      unexpected_state_transition_config;
  // Parse the `status_word_bit_index` into a `digital_input_pattern`.
  if (config.has_status_word_bit_index()) {
    size_t status_word_bit_index;
    status_word_bit_index = config.status_word_bit_index();
    // Valid bits are the Manufacturer-specific range of bit 8 - 15.
    // The bits 0 - 7 of the `Statusword` are defined by the DS402 standard.
    const size_t kMinBitIndex = 8;
    const size_t kMaxBitIndex = (sizeof(ProcessVariableType) * 8) - 1;

    if (status_word_bit_index < kMinBitIndex) {
      return absl::InvalidArgumentError(
          absl::StrCat("`status_word_bit_index` must be >= ", kMinBitIndex,
                       ", got ", status_word_bit_index));
    }

    if (status_word_bit_index > kMaxBitIndex) {
      return absl::InvalidArgumentError(
          absl::StrCat("`status_word_bit_index` must be <= ", kMaxBitIndex,
                       "got ", status_word_bit_index));
    }

    config.mutable_digital_input_pattern()
        ->mutable_digital_variable()
        ->set_variable_name(signal_variable_name);
    // The Signal is considered active when the bit is set.
    config.mutable_digital_input_pattern()
        ->mutable_value_pattern()
        ->set_signal_mask(1 << status_word_bit_index);
    config.mutable_digital_input_pattern()
        ->mutable_value_pattern()
        ->set_expected_value(1 << status_word_bit_index);
  }

  if (!config.has_digital_input_pattern()) {
    return absl::InvalidArgumentError(
        "Cannot configure `BEHAVIOR_IGNORE_ON_SIGNAL` without "
        "'status_word_bit_index' or 'digital_input_pattern'.");
  }

  const auto& pattern = config.digital_input_pattern();
  std::unique_ptr<fieldbus::ProcessVariable> digital_input_variable;
  if (pattern.digital_variable().has_array_index()) {
    // The user wants to configure a specific array field.
    auto process_variable = variable_registry.GetInputArrayFieldVariable(
        pattern.digital_variable().variable_name(),
        pattern.digital_variable().array_index());
    INTR_RETURN_IF_ERROR(process_variable.status());
    digital_input_variable = std::make_unique<fieldbus::ProcessVariable>(
        std::move(process_variable.value()));
  } else {
    // The user wants to configure a single variable.
    auto process_variable = variable_registry.GetInputVariable(
        pattern.digital_variable().variable_name());
    INTR_RETURN_IF_ERROR(process_variable.status());
    digital_input_variable = std::make_unique<fieldbus::ProcessVariable>(
        std::move(process_variable.value()));
  }

  const uint32_t signal_mask = pattern.value_pattern().signal_mask();
  const uint32_t expected_value = pattern.value_pattern().expected_value();

  INTR_RETURN_IF_ERROR(ValidateValuePattern(
      /*signal_mask=*/signal_mask, /*expected_value=*/expected_value,
      *digital_input_variable));
  configuration.signal_reader = [digital_input_variable =
                                     std::move(digital_input_variable),
                                 signal_mask, expected_value]() mutable {
    uint64_t bits = digital_input_variable->ReadRawUnchecked();
    return (bits & signal_mask) == expected_value;
  };

  return configuration;
}

absl::StatusOr<UnexpectedStateTransitionConfiguration>
ParseUnexpectedStateTransitionConfiguration(
    const intrinsic::fieldbus::VariableRegistry& variable_registry,
    const intrinsic_proto::icon::v1::Ds402BusComponent& config) {
  // Default initialization is to fault on unexpected state transitions.
  UnexpectedStateTransitionConfiguration default_configuration;

  if (!config.has_unexpected_state_transition_configuration()) {
    LOG(INFO) << "No unexpected state transition configuration "
                 "provided. Device will fault on unexpected state transitions.";
    return default_configuration;
  }

  switch (config.unexpected_state_transition_configuration().behavior()) {
    case Ds402UnexpectedStateTransitionConfiguration::BEHAVIOR_UNSPECIFIED:
      [[fallthrough]];
    case Ds402UnexpectedStateTransitionConfiguration::BEHAVIOR_FAULT: {
      LOG(INFO) << "Will fault on unexpected state transitions.";
      return default_configuration;
    }
    case Ds402UnexpectedStateTransitionConfiguration::
        BEHAVIOR_IGNORE_ON_SIGNAL: {
      INTR_ASSIGN_OR_RETURN(
          auto configuration,
          ParseUnexpectedStateTransitionSignal(
              variable_registry, config.status_word_variable_name(),
              config.unexpected_state_transition_configuration()));
      LOG(INFO) << "Will ignore unexpected state transitions when configured "
                   "signal is active.";

      return std::move(configuration);
    }
    case Ds402UnexpectedStateTransitionConfiguration::BEHAVIOR_IGNORE_ALWAYS: {
      LOG(INFO) << "Will ignore all unexpected state transitions.";
      return UnexpectedStateTransitionConfiguration{
          .behavior = Ds402UnexpectedStateTransitionConfiguration::
              BEHAVIOR_IGNORE_ALWAYS};
    }
    default:
      return absl::InvalidArgumentError(
          "Unsupported behavior for unexpected state transition "
          "configuration.");
  }
}

}  // namespace intrinsic::ds402
