// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_ADIO_ADIO_BUS_COMPONENT_H_
#define INTRINSIC_ICON_HAL_LIB_ADIO_ADIO_BUS_COMPONENT_H_

#include <memory>

#include "absl/functional/any_invocable.h"
#include "absl/status/statusor.h"
#include "intrinsic/icon/hal/lib/adio/v1/adio_bus_component_config.pb.h"
#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/hal/lib/fieldbus/bus_component.h"
#include "intrinsic/icon/hal/lib/fieldbus/bus_component_factory.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_or.h"

namespace intrinsic::adio {

// ADIO Bus Device.
// Exposes the a single digital input, analog input or digital output bus
// variable via hardware interface.
class AdioBusComponent : public fieldbus::BusComponent {
 public:
  // Creates an AdioBusComponent and registers its input or output hardware
  // interface.
  // Returns an error if neither input nor output are defined in the config, if
  // no bus variable can be found with the given name or if interface creation
  // fails, i.e. due to duplication.
  static absl::StatusOr<std::unique_ptr<AdioBusComponent>> Create(
      fieldbus::DeviceInitContext& device_init_context,
      const intrinsic_proto::icon::v1::AdioBusComponent& config);

  // Reads and updates the state of any input hardware interface.
  // Returns an error on bit size mismatch.
  // Calling `Read` on an output device is not an error, but will be a no-op.
  intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus> CyclicRead(
      fieldbus::RequestType request_type) override;

  // Writes the values of any output hardware interface to the bus.
  // Returns an error on bit size mismatch.
  // Calling `Write` on an input device is not an error, but will be a no-op.
  intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus> CyclicWrite(
      fieldbus::RequestType request_type) override;

 private:
  // Constructor for an AdioBusComponent.
  explicit AdioBusComponent(
      absl::AnyInvocable<intrinsic::icon::RealtimeStatus()>
          variable_to_interface,
      absl::AnyInvocable<intrinsic::icon::RealtimeStatus()>
          interface_to_variable);

  // Function that reads the bus variable and stores it into the interface.
  absl::AnyInvocable<intrinsic::icon::RealtimeStatus()> variable_to_interface_;

  // Function that writes the value from the interface into the bus variable.
  absl::AnyInvocable<intrinsic::icon::RealtimeStatus()> interface_to_variable_;
};

}  // namespace intrinsic::adio

// Registers the AdioBusComponent and its config type with the EtherCAT bus
// device factory. Allows constructing the device from its config type via
// `CreateBusComponentFromConfig`.
REGISTER_BUS_COMPONENT(adio::AdioBusComponent,
                       intrinsic_proto::icon::v1::AdioBusComponent);

#endif  // INTRINSIC_ICON_HAL_LIB_ADIO_ADIO_BUS_COMPONENT_H_
