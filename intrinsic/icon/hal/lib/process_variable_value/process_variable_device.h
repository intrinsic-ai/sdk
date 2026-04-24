// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_PDO_VALUE_PROCESS_VARIABLE_DEVICE_H_
#define INTRINSIC_ICON_HAL_LIB_PDO_VALUE_PROCESS_VARIABLE_DEVICE_H_

#include <functional>
#include <memory>

#include "absl/status/statusor.h"
#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/hal/lib/fieldbus/bus_component.h"
#include "intrinsic/icon/hal/lib/fieldbus/bus_component_factory.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/hal/lib/process_variable_value/v1/process_variable_config.pb.h"
#include "intrinsic/icon/utils/realtime_status_or.h"

namespace intrinsic::fieldbus {

// PDO Value Device.
// Allows to set a value to a bus variable.
// Example configuration:
// pdo_values: [
//   {
//     variable_name: "Drive 1 (CMMT-AS-MP-S1).Outputs.Modes of operation"
//     value: {
//       int8_value: 8  # CSP, Cyclic Sync Position
//     }
//   }
// ]
class ProcessVariableDevice : public fieldbus::BusComponent {
 public:
  // Creates a ProcessVariableDevice.
  // Returns an error if the a bus variable cannot be found, if the specified
  // value type is incompatible with that of the bus variable or (for 8 and
  // 16bit variables) if the value is out of range.
  // Note that the implementation does not capture the same output variable
  // being written to by multiple devices (this or other devices).
  static absl::StatusOr<std::unique_ptr<ProcessVariableDevice>> Create(
      fieldbus::DeviceInitContext& init_context,
      const intrinsic_proto::fieldbus::v1::ProcessVariableConfig& config);

  // NO-OP, since constant values are only written.
  intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus> CyclicRead(
      fieldbus::RequestType request_type) override;

  // Writes the configured values to their respective bus variables.
  intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus> CyclicWrite(
      fieldbus::RequestType request_type) override;

 private:
  // Constructor for a ProcessVariableDevice.
  // `constants_to_variable_when_enabling` is the function that writes the value
  // from the constants into the bus variable during the enabling state.
  // `constants_to_variable_when_operational` is the function that writes the
  // value from the constants into the bus variable during the operational
  // state.
  // `constants_to_variable_when_disabling` is the function that writes the
  // value from the constants into the bus variable during the disabling state.
  ProcessVariableDevice(
      std::function<void()> constants_to_variable_when_enabling,
      std::function<void()> constants_to_variable_when_operational,
      std::function<void()> constants_to_variable_when_disabling);

  // Function that writes the value from the constants into the bus variable
  // during the enabling state.
  std::function<void()> constants_to_variable_when_enabling_;
  // Function that writes the value from the constants into the bus variable
  // during the operational state.
  std::function<void()> constants_to_variable_when_operational_;
  // Function that writes the value from the constants into the bus variable
  // during the disabling state.
  std::function<void()> constants_to_variable_when_disabling_;
};

}  // namespace intrinsic::fieldbus

// Registers the ProcessVariableDevice and its config type with the fieldbus
// component factory. Allows constructing the device from its config type via
// `CreateBusComponentFromConfig`.
REGISTER_BUS_COMPONENT(fieldbus::ProcessVariableDevice,
                       intrinsic_proto::fieldbus::v1::ProcessVariableConfig);

#endif  // INTRINSIC_ICON_HAL_LIB_PDO_VALUE_PROCESS_VARIABLE_DEVICE_H_
