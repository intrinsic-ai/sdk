// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_FIELDBUS_BUS_DEVICE_FACTORY_H_
#define INTRINSIC_ICON_HAL_LIB_FIELDBUS_BUS_DEVICE_FACTORY_H_

#include <memory>

#include "absl/log/log.h"
#include "grpcpp/server_builder.h"
#include "intrinsic/icon/hal/hardware_module_init_context.h"
#include "intrinsic/icon/hal/lib/fieldbus/bus_component.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"

namespace intrinsic::fieldbus {

// Defines a factory for creating bus devices.
// Works by bus device implementation calling REGISTER_BUS_COMPONENT with
// their own type and that of their config proto.

// Type trait used to map a device `ConfigType` to the device type.
// Bus device implementations should use the REGISTER_BUS_COMPONENT macro
// below.
template <typename ConfigType>
struct BusComponentFromConfig;

// Creates a bus device from a given its configuration.
// Uses the type trait above to determine the device type corresponding to
// `BusComponentConfig`. Expects the bus device implementation to provide a
// static `Create` function with the same parameters. `args` may be used to pass
// device specific arguments to `Create`.
template <typename BusComponentConfig, typename... Args>
absl::StatusOr<std::unique_ptr<BusComponent>> CreateBusComponentFromConfig(
    intrinsic::icon::HardwareModuleInitContext& init_context,
    const intrinsic::fieldbus::VariableRegistry& variable_registry,
    const BusComponentConfig& config, Args&&... args) {
  DeviceInitContext device_init_context(
      init_context.GetInterfaceRegistry(), variable_registry,
      /*register_grpc_service=*/[&init_context](grpc::Service& service) {
        init_context.RegisterGrpcService(service);
      });
  return BusComponentFromConfig<BusComponentConfig>::Type::Create(
      device_init_context, config, std::forward<Args>(args)...);
}

// Macro to register a bus device implementation.
// Expected to be called outside of any namespace inside the bus device's
// header.
#define REGISTER_BUS_COMPONENT(device_type, config_type) \
  namespace intrinsic::fieldbus {                        \
  template <>                                            \
  struct BusComponentFromConfig<config_type> {           \
    using Type = device_type;                            \
  };                                                     \
  }  // namespace intrinsic::fieldbus

}  // namespace intrinsic::fieldbus

#endif  // INTRINSIC_ICON_HAL_LIB_FIELDBUS_BUS_DEVICE_FACTORY_H_
