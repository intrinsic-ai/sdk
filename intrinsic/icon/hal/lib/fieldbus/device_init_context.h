// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef INTRINSIC_ICON_HAL_LIB_FIELDBUS_DEVICE_INIT_CONTEXT_H_
#define INTRINSIC_ICON_HAL_LIB_FIELDBUS_DEVICE_INIT_CONTEXT_H_

#include <functional>
#include <utility>

#include "grpcpp/impl/service_type.h"
#include "intrinsic/icon/hal/hardware_interface_registry.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"

namespace intrinsic::fieldbus {

// Provides configuration and functions needed during initialization of a
// fieldbus device such as
// - access to the interface registry
// - access to the variable registry
// - the ability to register a gRPC service
class DeviceInitContext {
 public:
  using RegisterGrpcServiceFn = std::function<void(grpc::Service& service)>;

  DeviceInitContext(
      intrinsic::icon::HardwareInterfaceRegistry& interface_registry,
      const intrinsic::fieldbus::VariableRegistry& variable_registry,
      RegisterGrpcServiceFn register_grpc_service)
      : interface_registry_(interface_registry),
        variable_registry_(variable_registry),
        register_grpc_service_(std::move(register_grpc_service)) {}

  // Constructor for devices that don't need to register a gRPC service. Useful
  // for tests.
  DeviceInitContext(
      intrinsic::icon::HardwareInterfaceRegistry& interface_registry,
      const intrinsic::fieldbus::VariableRegistry& variable_registry)
      : interface_registry_(interface_registry),
        variable_registry_(variable_registry) {}

  // Delete copy and move constructors since this class contains temporary
  // objects which are deleted after the device is initialized.
  // Devices should not be able to copy or move this class.
  DeviceInitContext(const DeviceInitContext&) = delete;
  DeviceInitContext& operator=(const DeviceInitContext&) = delete;
  DeviceInitContext& operator=(DeviceInitContext&&) = delete;

  // Returns the interface registry for this Hardware Module to register
  // interfaces.
  intrinsic::icon::HardwareInterfaceRegistry& GetInterfaceRegistry() const {
    return interface_registry_;
  }

  // Returns the variable registry for this Hardware Module to register
  // fieldbus variables.
  const intrinsic::fieldbus::VariableRegistry& GetVariableRegistry() const {
    return variable_registry_;
  }

  // Registers a gRPC service with the hardware module runtime. The runtime
  // makes this service available to external components some time after the
  // hardware module's `Init()` function returns.
  //
  // Attention: `service` must live until the `Shutdown()` of the Hardware
  // Module is called!
  //
  // The gRPC service will still be served even if
  // HardwareModuleInterface::Init() returns an error.
  //
  // The gRPC service will run on a port that is reachable from external
  // components such as the frontend.
  void RegisterGrpcService(grpc::Service& service) const {
    register_grpc_service_(service);
  }

 private:
  intrinsic::icon::HardwareInterfaceRegistry& interface_registry_;
  const intrinsic::fieldbus::VariableRegistry& variable_registry_;
  RegisterGrpcServiceFn register_grpc_service_;
};

}  // namespace intrinsic::fieldbus

#endif  // INTRINSIC_ICON_HAL_LIB_FIELDBUS_DEVICE_INIT_CONTEXT_H_
