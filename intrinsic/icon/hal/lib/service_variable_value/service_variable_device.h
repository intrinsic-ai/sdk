// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_SERVICE_VARIABLE_VALUE_SERVICE_VARIABLE_DEVICE_H_
#define INTRINSIC_ICON_HAL_LIB_SERVICE_VARIABLE_VALUE_SERVICE_VARIABLE_DEVICE_H_

#include <concepts>
#include <memory>
#include <string>
#include <vector>

#include "absl/functional/any_invocable.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/interfaces/io_controller.fbs.h"
#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/hal/lib/fieldbus/bus_component.h"
#include "intrinsic/icon/hal/lib/fieldbus/bus_component_factory.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/hal/lib/fieldbus/v1/value_parsing.pb.h"
#include "intrinsic/icon/hal/lib/service_variable_value/v1/service_variable_config.pb.h"
#include "intrinsic/icon/utils/async_buffer.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/util/thread/thread.h"

namespace intrinsic::fieldbus {

// Service variable Value Device.
// Allows to set/get service variable values of a particular bus participant.
class ServiceVariableDevice : public fieldbus::BusComponent {
 public:
  // The minimum number of cycles per service variable read.
  // This is an estimate, the bus might be able to handle more reads, but we
  // don't want to investigate the limits of the bus here, as we believe this
  // feature is rarely pushed to the this limit.
  static constexpr int kMinCyclesPerSdoRead = 10;

  // Creates a ServiceVariableDevice.
  // Returns an error if a service variable variable for writing cannot be
  // found, if the specified value type is incompatible with that of the service
  // variable variable or (for 8 and 16bit variables) if the value is out of
  // range. Note that the implementation does not capture the same service
  // variable variable being written to by multiple devices (this or other
  // devices). service variable values are written during Create(). Reading
  // service variable values is done a separate thread. By default the read
  // values are exported via PubSub. Incompatible service variable values are
  // ignored and an error is logged.
  static absl::StatusOr<std::unique_ptr<ServiceVariableDevice>> Create(
      fieldbus::DeviceInitContext& init_context,
      const intrinsic_proto::fieldbus::v1::ServiceVariableConfig& config,
      absl::string_view context_name, absl::Duration cycle_time);

  ~ServiceVariableDevice() override;

  // Returns all strings (keys) of the `interpretations_of_values` map that
  // match the `value` and `signal_mask` of the `sdo_read`.
  // Returns an empty vector if no string matches. Currently only int types are
  // supported.
  template <std::integral T>
  static std::vector<std::string> InterpretationsOfValue(
      const intrinsic_proto::fieldbus::v1::ServiceVariableRead& sdo_read,
      T value) {
    std::vector<std::string> result;

    if (sdo_read.interpretations_of_values().empty()) {
      return result;
    }
    for (const auto& [value_interpretation, interpretations] :
         sdo_read.interpretations_of_values()) {
      for (const auto& pattern : interpretations.patterns()) {
        if ((value & pattern.signal_mask()) == pattern.expected_value()) {
          result.push_back(value_interpretation);
        }
      }
    }
    return result;
  }

  // NO-OP, since service variable values are only written during Create().
  intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus> CyclicRead(
      fieldbus::RequestType request_type) override;

  // NO-OP, since service variable values are only written during Create().
  intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus> CyclicWrite(
      fieldbus::RequestType request_type) override;

 private:
  explicit ServiceVariableDevice(
      intrinsic::Thread sdo_read_thread,
      std::vector<absl::AnyInvocable<intrinsic::icon::RealtimeStatus()>>
          interface_updaters);

  // The interface_updaters_ hold pointers to memory owned by the
  // sdo_read_thread_. Thus, the thread must be destroyed after the updaters
  // vector to prevent dangling pointer issues.
  intrinsic::Thread sdo_read_thread_;
  std::vector<absl::AnyInvocable<intrinsic::icon::RealtimeStatus()>>
      interface_updaters_;
};

}  // namespace intrinsic::fieldbus

// Registers the ServiceVariableDevice and its config type with the fieldbus
// component factory. Allows constructing the device from its config type via
// `CreateBusComponentFromConfig`.
REGISTER_BUS_COMPONENT(fieldbus::ServiceVariableDevice,
                       intrinsic_proto::fieldbus::v1::ServiceVariableConfig);

#endif  // INTRINSIC_ICON_HAL_LIB_SERVICE_VARIABLE_VALUE_SERVICE_VARIABLE_DEVICE_H_
