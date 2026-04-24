// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_ROBOT_ARM_ROBOT_ARM_BUS_COMPONENT_H_
#define INTRINSIC_ICON_HAL_LIB_ROBOT_ARM_ROBOT_ARM_BUS_COMPONENT_H_

#include <functional>
#include <memory>
#include <vector>

#include "absl/status/statusor.h"
#include "intrinsic/icon/hal/default_hardware_interfaces.h"  // IWYU pragma: keep
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/interfaces/joint_command.fbs.h"
#include "intrinsic/icon/hal/interfaces/joint_state.fbs.h"
#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/hal/lib/fieldbus/bus_component.h"
#include "intrinsic/icon/hal/lib/fieldbus/bus_component_factory.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/hal/lib/robot_arm/v1/robot_arm_bus_component_config.pb.h"
#include "intrinsic/icon/utils/realtime_status_or.h"

namespace intrinsic::robot_arm {

// Robot Arm bus device.
// A bus device that writes position commands to the EtherCAT bus, reads
// position states and also advertises velocity and acceleration states as well
// as joint limits of a robot arm.
class RobotArmBusComponent : public fieldbus::BusComponent {
 public:
  // Creates a RobotArmBusComponent and registers its hardware interfaces.
  // Returns an error if neither if any of the configured bus variable can be
  // found with the given name or if interface creation fails, i.e. due to
  // duplication.
  static absl::StatusOr<std::unique_ptr<RobotArmBusComponent>> Create(
      fieldbus::DeviceInitContext& init_context,
      const intrinsic_proto::icon::v1::RobotArmBusComponentConfig& config);

  // Reads and updates the state of any input hardware interface.
  intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus> CyclicRead(
      fieldbus::RequestType request_type) override;

  // Writes the values of any output hardware interface to the bus.
  intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus> CyclicWrite(
      fieldbus::RequestType request_type) override;

 private:
  RobotArmBusComponent(
      intrinsic::icon::MutableHardwareInterfaceHandle<
          intrinsic_fbs::JointPositionState>
          position_state,
      intrinsic::icon::MutableHardwareInterfaceHandle<
          intrinsic_fbs::JointVelocityState>
          velocity_state,
      intrinsic::icon::MutableHardwareInterfaceHandle<
          intrinsic_fbs::JointAccelerationState>
          acceleration_state,
      intrinsic::icon::StrictHardwareInterfaceHandle<
          intrinsic_fbs::JointPositionCommand>
          position_command,
      std::vector<std::function<double()>> scaled_position_state_getters,
      std::vector<std::function<double()>> scaled_velocity_state_getters,
      std::vector<std::function<double()>> scaled_acceleration_state_getters,
      std::vector<std::function<void(double)>> scaled_position_command_setters,
      std::vector<std::function<void(double)>>
          scaled_feedforward_velocity_command_setters);

  // Hardware interfaces.
  intrinsic::icon::MutableHardwareInterfaceHandle<
      intrinsic_fbs::JointPositionState>
      position_state_;
  intrinsic::icon::MutableHardwareInterfaceHandle<
      intrinsic_fbs::JointVelocityState>
      velocity_state_;
  intrinsic::icon::MutableHardwareInterfaceHandle<
      intrinsic_fbs::JointAccelerationState>
      acceleration_state_;
  intrinsic::icon::StrictHardwareInterfaceHandle<
      intrinsic_fbs::JointPositionCommand>
      position_command_;

  // List of functions that read and scale (to radian) the joint position values
  // from the bus process data.
  std::vector<std::function<double()>> scaled_position_state_getters_;

  // List of functions that read and scale (to radian/s) the joint velocity
  // values from the bus process data.
  std::vector<std::function<double()>> scaled_velocity_state_getters_;
  // List of functions that read and scale (to radian/s²) the joint acceleration
  // values from the bus process data.
  std::vector<std::function<double()>> scaled_acceleration_state_getters_;
  // List of functions that scaled (from radian) and write the joint command
  // values to the bus process data.
  std::vector<std::function<void(double)>> scaled_position_command_setters_;
  // List of functions that scaled (from radian/s) and write the joint
  // feedforward velocity command values to the bus process data.
  std::vector<std::function<void(double)>>
      scaled_feedforward_velocity_command_setters_;
};

}  // namespace intrinsic::robot_arm

// Registers the RobotArmBusComponent and its config type with the EtherCAT bus
// device factory. Allows constructing the device from its config type via
// `CreateBusComponentFromConfig`.
REGISTER_BUS_COMPONENT(intrinsic::robot_arm::RobotArmBusComponent,
                       intrinsic_proto::icon::v1::RobotArmBusComponentConfig);

#endif  // INTRINSIC_ICON_HAL_LIB_ROBOT_ARM_ROBOT_ARM_BUS_COMPONENT_H_
