// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/adio/adio_bus_component.h"

#include <algorithm>
#include <bitset>
#include <cstddef>
#include <cstdint>
#include <limits>
#include <memory>
#include <string>
#include <utility>
#include <vector>

#include "absl/functional/any_invocable.h"
#include "absl/log/log.h"
#include "absl/memory/memory.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "intrinsic/icon/hal/default_hardware_interfaces.h"  // IWYU pragma: keep
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/hardware_interface_registry.h"
#include "intrinsic/icon/hal/interfaces/adio.fbs.h"
#include "intrinsic/icon/hal/interfaces/io_controller.fbs.h"
#include "intrinsic/icon/hal/lib/adio/v1/adio_bus_component_config.pb.h"
#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable_utils.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"
#include "intrinsic/icon/utils/fixed_str_cat.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_macro.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::adio {

AdioBusComponent::AdioBusComponent(
    absl::AnyInvocable<intrinsic::icon::RealtimeStatus()> variable_to_interface,
    absl::AnyInvocable<intrinsic::icon::RealtimeStatus()> interface_to_variable)
    : variable_to_interface_(std::move(variable_to_interface)),
      interface_to_variable_(std::move(interface_to_variable)) {}

template <typename T>
absl::AnyInvocable<intrinsic::icon::RealtimeStatus(double)> CreateAnalogWriterT(
    fieldbus::ProcessVariable analog_output_variable) {
  return [analog_output_variable](
             double value) mutable -> intrinsic::icon::RealtimeStatus {
    if (fieldbus::IsOutOfRange<double, T>(value)) {
      return intrinsic::icon::InvalidArgumentError(
          intrinsic::icon::FixedStrCat<
              intrinsic::icon::RealtimeStatus::kMaxMessageLength>(
              "Analog output value is out of range. Type range is [",
              std::numeric_limits<T>::lowest(), ", ",
              std::numeric_limits<T>::max(), "] and value is ", value));
    }
    analog_output_variable.WriteUnchecked(static_cast<T>(value));
    return intrinsic::icon::OkStatus();
  };
}

absl::StatusOr<absl::AnyInvocable<intrinsic::icon::RealtimeStatus(double)>>
CreateAnalogWriter(fieldbus::ProcessVariable analog_output_variable) {
  if (analog_output_variable.IsCompatibleType<uint8_t>().ok()) {
    return CreateAnalogWriterT<uint8_t>(analog_output_variable);
  } else if (analog_output_variable.IsCompatibleType<uint16_t>().ok()) {
    return CreateAnalogWriterT<uint16_t>(analog_output_variable);
  } else if (analog_output_variable.IsCompatibleType<uint32_t>().ok()) {
    return CreateAnalogWriterT<uint32_t>(analog_output_variable);
  } else if (analog_output_variable.IsCompatibleType<uint64_t>().ok()) {
    return CreateAnalogWriterT<uint64_t>(analog_output_variable);
  } else if (analog_output_variable.IsCompatibleType<int8_t>().ok()) {
    return CreateAnalogWriterT<int8_t>(analog_output_variable);
  } else if (analog_output_variable.IsCompatibleType<int16_t>().ok()) {
    return CreateAnalogWriterT<int16_t>(analog_output_variable);
  } else if (analog_output_variable.IsCompatibleType<int32_t>().ok()) {
    return CreateAnalogWriterT<int32_t>(analog_output_variable);
  } else if (analog_output_variable.IsCompatibleType<int64_t>().ok()) {
    return CreateAnalogWriterT<int64_t>(analog_output_variable);
  } else if (analog_output_variable.IsCompatibleType<double>().ok()) {
    return CreateAnalogWriterT<double>(analog_output_variable);
  } else if (analog_output_variable.IsCompatibleType<float>().ok()) {
    return CreateAnalogWriterT<float>(analog_output_variable);
  } else if (analog_output_variable.IsCompatibleType<bool>().ok()) {
    return CreateAnalogWriterT<bool>(analog_output_variable);
  }
  return absl::InvalidArgumentError(
      "Analog output variable type is not supported.");
}

absl::StatusOr<absl::AnyInvocable<double()>> CreateAnalogReader(
    fieldbus::ProcessVariable analog_input_variable) {
  if (analog_input_variable.IsCompatibleType<uint8_t>().ok()) {
    return [analog_input_variable]() {
      return static_cast<double>(
          analog_input_variable.ReadUnchecked<uint8_t>());
    };
  } else if (analog_input_variable.IsCompatibleType<uint16_t>().ok()) {
    return [analog_input_variable]() {
      return static_cast<double>(
          analog_input_variable.ReadUnchecked<uint16_t>());
    };
  } else if (analog_input_variable.IsCompatibleType<uint32_t>().ok()) {
    return [analog_input_variable]() {
      return static_cast<double>(
          analog_input_variable.ReadUnchecked<uint32_t>());
    };
  } else if (analog_input_variable.IsCompatibleType<uint64_t>().ok()) {
    return [analog_input_variable]() {
      return static_cast<double>(
          analog_input_variable.ReadUnchecked<uint64_t>());
    };
  } else if (analog_input_variable.IsCompatibleType<int8_t>().ok()) {
    return [analog_input_variable]() {
      return static_cast<double>(analog_input_variable.ReadUnchecked<int8_t>());
    };
  } else if (analog_input_variable.IsCompatibleType<int16_t>().ok()) {
    return [analog_input_variable]() {
      return static_cast<double>(
          analog_input_variable.ReadUnchecked<int16_t>());
    };
  } else if (analog_input_variable.IsCompatibleType<int32_t>().ok()) {
    return [analog_input_variable]() {
      return static_cast<double>(
          analog_input_variable.ReadUnchecked<int32_t>());
    };
  } else if (analog_input_variable.IsCompatibleType<int64_t>().ok()) {
    return [analog_input_variable]() {
      return static_cast<double>(
          analog_input_variable.ReadUnchecked<int64_t>());
    };
  } else if (analog_input_variable.IsCompatibleType<double>().ok()) {
    return [analog_input_variable]() {
      return analog_input_variable.ReadUnchecked<double>();
    };
  } else if (analog_input_variable.IsCompatibleType<float>().ok()) {
    return [analog_input_variable]() {
      return static_cast<double>(analog_input_variable.ReadUnchecked<float>());
    };
  } else if (analog_input_variable.IsCompatibleType<bool>().ok()) {
    return [analog_input_variable]() {
      return static_cast<double>(analog_input_variable.ReadUnchecked<bool>());
    };
  } else {
    return absl::InvalidArgumentError(
        "Analog input variable type is not supported.");
  }
}

absl::AnyInvocable<intrinsic::icon::RealtimeStatus()> CreateDigitalReader(
    intrinsic::icon::MutableHardwareInterfaceHandle<intrinsic_fbs::DIOStatus>
        digital_input_interface,
    fieldbus::ProcessVariable digital_input_variable) {
  return [digital_input_interface = std::move(digital_input_interface),
          digital_input_variable]() mutable {
    // Read the bits into a bitset.
    std::bitset<sizeof(digital_input_variable.ReadRawUnchecked()) * 8> bits(
        digital_input_variable.ReadRawUnchecked());

    // Update the `digital_input_interface_` with the values from `bits`.
    for (std::size_t i = 0; i < digital_input_variable.bit_size(); ++i) {
      digital_input_interface->mutable_signals()
          ->GetMutableObject(i)
          ->mutate_bit_number(i);
      digital_input_interface->mutable_signals()
          ->GetMutableObject(i)
          ->mutate_value(bits[i]);
    }
    return intrinsic::icon::OkStatus();
  };
}
absl::AnyInvocable<intrinsic::icon::RealtimeStatus()> CreateDigitalWriter(
    intrinsic::icon::HardwareInterfaceHandle<intrinsic_fbs::DIOCommand>
        digital_output_interface,
    fieldbus::ProcessVariable digital_output_variable) {
  return [digital_output_interface = std::move(digital_output_interface),
          digital_output_variable]() mutable {
    // Read the values from the `digital_output_interface_` and pack them into a
    // bitset. The bitset has the maximum available size. Since it is default
    // constructed, all bits are set to 0.
    std::bitset<sizeof(uint64_t) * 8> bits;
    for (std::size_t i = 0; i < digital_output_interface->signals()->size();
         ++i) {
      bits.set(i, digital_output_interface->signals()->Get(i)->value());
    }

    // Write the bitset to the `digital_output_variable`.
    digital_output_variable.WriteRawUnchecked(bits.to_ulong());
    return intrinsic::icon::OkStatus();
  };
}

absl::StatusOr<std::unique_ptr<AdioBusComponent>> AdioBusComponent::Create(
    fieldbus::DeviceInitContext& device_init_context,
    const intrinsic_proto::icon::v1::AdioBusComponent& config) {
  // Helper generator to populate the `bit_description` vectors below.
  struct increment {
    std::size_t value;
    increment() : value(0) {}
    // Uses post increment because the bit index is zero indexed.
    std::string operator()() { return std::to_string(value++); }
  };

  // Default (no-op) read and write functions.
  absl::AnyInvocable<intrinsic::icon::RealtimeStatus()> variable_to_interface =
      []() { return intrinsic::icon::OkStatus(); };
  absl::AnyInvocable<intrinsic::icon::RealtimeStatus()> interface_to_variable =
      []() { return intrinsic::icon::OkStatus(); };

  const intrinsic::fieldbus::VariableRegistry& variable_registry =
      device_init_context.GetVariableRegistry();
  intrinsic::icon::HardwareInterfaceRegistry& interface_registry =
      device_init_context.GetInterfaceRegistry();

  if (config.has_digital_input_variable()) {
    // Configure the device for a digital input.
    std::unique_ptr<fieldbus::ProcessVariable> digital_input_variable;
    if (config.digital_input_variable().has_array_index()) {
      // The user wants to configure a specific array field.
      auto process_variable = variable_registry.GetInputArrayFieldVariable(
          config.digital_input_variable().variable_name(),
          config.digital_input_variable().array_index());
      INTR_RETURN_IF_ERROR(process_variable.status());
      digital_input_variable = std::make_unique<fieldbus::ProcessVariable>(
          std::move(process_variable.value()));
    } else {
      // The user wants to configure a single variable.
      auto process_variable = variable_registry.GetInputVariable(
          config.digital_input_variable().variable_name());
      INTR_RETURN_IF_ERROR(process_variable.status());
      digital_input_variable = std::make_unique<fieldbus::ProcessVariable>(
          std::move(process_variable.value()));
    }
    const auto bit_size = digital_input_variable->bit_size();

    std::vector<std::string> bit_descriptions(bit_size);

    // Number the bits sequentially. Some may be overridden by the
    // `bit_index_to_alias` field below.
    std::generate_n(bit_descriptions.begin(), bit_size, increment());
    for (const auto& [bit_index, alias] :
         config.digital_input_variable().bit_index_to_alias()) {
      // 'bit_index' is an unsigned integer and can't be negative.
      if (bit_index >= bit_size) {
        return absl::InvalidArgumentError(
            absl::StrCat("The index '", bit_index, "' for the alias '", alias,
                         "' must be smaller than the number of "
                         "bits '",
                         bit_size, "' in the digital input variable ",
                         config.digital_input_variable().variable_name()));
      }
      bit_descriptions[bit_index] = alias;
    }

    INTR_ASSIGN_OR_RETURN(
        auto digital_input_interface,
        interface_registry.AdvertiseMutableInterface<intrinsic_fbs::DIOStatus>(
            config.interface_name(), bit_descriptions));

    variable_to_interface = CreateDigitalReader(
        std::move(digital_input_interface), *digital_input_variable);

    // Using `new` to access a non-public constructor.
    return absl::WrapUnique(new AdioBusComponent(
        std::move(variable_to_interface), std::move(interface_to_variable)));
  } else if (config.has_analog_input_variables()) {
    // Configure the device for a analog inputs.
    std::vector<std::string> field_description = {};
    std::vector<absl::AnyInvocable<double()>> analog_reader_functions = {};

    for (std::size_t i = 0;
         i < config.analog_input_variables().variables_size(); ++i) {
      // Get the field description aka human readable variable alias.
      if (config.analog_input_variables().variables(i).has_alias()) {
        field_description.push_back(
            config.analog_input_variables().variables(i).alias());
      } else {
        field_description.push_back(absl::StrCat(i));
      }

      // Create the bus variable
      INTR_ASSIGN_OR_RETURN(
          auto analog_input_variable,
          variable_registry.GetInputVariable(
              config.analog_input_variables().variables(i).variable_name()));

      // Create the reader for this bus variable.
      INTR_ASSIGN_OR_RETURN(auto analog_reader_func,
                            CreateAnalogReader(analog_input_variable));
      analog_reader_functions.emplace_back(std::move(analog_reader_func));
    }

    INTR_ASSIGN_OR_RETURN(
        auto analog_input_interface,
        interface_registry.AdvertiseMutableInterface<intrinsic_fbs::AIOStatus>(
            config.interface_name(), field_description));

    variable_to_interface =
        [analog_input_interface = std::move(analog_input_interface),
         analog_reader_functions = std::move(analog_reader_functions)]() mutable
        -> intrinsic::icon::RealtimeStatus {
      for (std::size_t i = 0; i < analog_reader_functions.size(); ++i) {
        auto signal =
            analog_input_interface->mutable_signals()->GetMutableObject(i);
        signal->mutate_value(analog_reader_functions[i]());
      }
      return intrinsic::icon::OkStatus();
    };
    // Using `new` to access a non-public constructor.
    return absl::WrapUnique(new AdioBusComponent(
        std::move(variable_to_interface), std::move(interface_to_variable)));
  } else if (config.has_digital_output_variable()) {
    // Configure the device for a digital output.
    std::unique_ptr<fieldbus::ProcessVariable> digital_output_variable;
    if (config.digital_output_variable().has_array_index()) {
      // The user wants to configure a specific array field.
      auto process_variable = variable_registry.GetOutputArrayFieldVariable(
          config.digital_output_variable().variable_name(),
          config.digital_output_variable().array_index());
      INTR_RETURN_IF_ERROR(process_variable.status());
      digital_output_variable = std::make_unique<fieldbus::ProcessVariable>(
          std::move(process_variable.value()));
    } else {
      // The user wants to configure a single variable.
      auto process_variable = variable_registry.GetOutputVariable(
          config.digital_output_variable().variable_name());
      INTR_RETURN_IF_ERROR(process_variable.status());
      digital_output_variable = std::make_unique<fieldbus::ProcessVariable>(
          std::move(process_variable.value()));
    }
    const auto bit_size = digital_output_variable->bit_size();

    std::vector<std::string> bit_descriptions(bit_size);
    // Number the bits sequentially. Some may be overridden by the
    // `bit_index_to_alias` field below.
    std::generate_n(bit_descriptions.begin(), bit_size, increment());
    for (const auto& [bit_index, alias] :
         config.digital_output_variable().bit_index_to_alias()) {
      // 'bit_index' is an unsigned integer and can't be negative.
      if (bit_index >= bit_size) {
        return absl::InvalidArgumentError(
            absl::StrCat("The index '", bit_index, "' for the alias '", alias,
                         "' must be smaller than the number of "
                         "bits '",
                         bit_size, "' in the digital output variable ",
                         config.digital_output_variable().variable_name()));
      }
      bit_descriptions[bit_index] = alias;
    }

    INTR_ASSIGN_OR_RETURN(
        auto digital_output_interface,
        interface_registry.AdvertiseInterface<intrinsic_fbs::DIOCommand>(
            config.interface_name(), bit_descriptions));

    interface_to_variable = CreateDigitalWriter(
        std::move(digital_output_interface), *digital_output_variable);

    // Using `new` to access a non-public constructor.
    return absl::WrapUnique(new AdioBusComponent(
        std::move(variable_to_interface), std::move(interface_to_variable)));
  } else if (config.has_analog_output_variables()) {
    // Configure the device for a analog outputs.
    std::vector<std::string> field_description = {};
    std::vector<absl::AnyInvocable<intrinsic::icon::RealtimeStatus(double)>>
        analog_writer_functions = {};

    for (std::size_t i = 0;
         i < config.analog_output_variables().variables_size(); ++i) {
      // Get the field description aka human readable variable alias.
      if (config.analog_output_variables().variables(i).has_alias()) {
        field_description.push_back(
            config.analog_output_variables().variables(i).alias());
      } else {
        field_description.push_back(absl::StrCat(i));
      }

      // Create the bus variable
      INTR_ASSIGN_OR_RETURN(
          auto analog_output_variable,
          variable_registry.GetOutputVariable(
              config.analog_output_variables().variables(i).variable_name()));

      // Create the writer for this bus variable.
      INTR_ASSIGN_OR_RETURN(auto analog_writer_func,
                            CreateAnalogWriter(analog_output_variable));
      analog_writer_functions.emplace_back(std::move(analog_writer_func));
    }

    INTR_ASSIGN_OR_RETURN(
        auto analog_output_interface,
        interface_registry.AdvertiseInterface<intrinsic_fbs::AIOCommand>(
            config.interface_name(), field_description));

    interface_to_variable =
        [analog_output_interface = std::move(analog_output_interface),
         analog_writer_functions = std::move(analog_writer_functions)]() mutable
        -> intrinsic::icon::RealtimeStatus {
      intrinsic::icon::RealtimeStatus status = intrinsic::icon::OkStatus();
      for (std::size_t i = 0; i < analog_writer_functions.size(); ++i) {
        status = intrinsic::icon::OverwriteIfNotInError(
            status, analog_writer_functions[i](
                        analog_output_interface->signals()->Get(i)->value()));
      }
      return status;
    };
    // Using `new` to access a non-public constructor.
    return absl::WrapUnique(new AdioBusComponent(
        std::move(variable_to_interface), std::move(interface_to_variable)));
  }
  return absl::InvalidArgumentError(
      "Configuration does not provide an analog or digital input or output "
      "variable name.");
}

intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus>
AdioBusComponent::CyclicRead(fieldbus::RequestType) {
  INTRINSIC_RT_RETURN_IF_ERROR(variable_to_interface_());
  return fieldbus::RequestStatus::kDone;
}

intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus>
AdioBusComponent::CyclicWrite(fieldbus::RequestType) {
  INTRINSIC_RT_RETURN_IF_ERROR(interface_to_variable_());
  return fieldbus::RequestStatus::kDone;
}

}  // namespace intrinsic::adio
