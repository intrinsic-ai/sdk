// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_DEFAULT_HARDWARE_INTERFACES_H_
#define INTRINSIC_ICON_HAL_DEFAULT_HARDWARE_INTERFACES_H_

#include "intrinsic/icon/control/safety/safety_messages.fbs.h"
#include "intrinsic/icon/control/safety/safety_messages_utils.h"
#include "intrinsic/icon/flatbuffers/transform_types.fbs.h"
#include "intrinsic/icon/flatbuffers/transform_types.h"
#include "intrinsic/icon/hal/hardware_interface_traits.h"
#include "intrinsic/icon/hal/interfaces/control_mode.fbs.h"
#include "intrinsic/icon/hal/interfaces/control_mode_utils.h"
#include "intrinsic/icon/hal/interfaces/electrical_motor.fbs.h"
#include "intrinsic/icon/hal/interfaces/electrical_motor_utils.h"
#include "intrinsic/icon/hal/interfaces/force_torque.fbs.h"
#include "intrinsic/icon/hal/interfaces/force_torque_utils.h"
#include "intrinsic/icon/hal/interfaces/hardware_module_state.fbs.h"
#include "intrinsic/icon/hal/interfaces/hardware_module_state_utils.h"
#include "intrinsic/icon/hal/interfaces/io_controller.fbs.h"
#include "intrinsic/icon/hal/interfaces/io_controller_utils.h"
#include "intrinsic/icon/hal/interfaces/joint_command.fbs.h"
#include "intrinsic/icon/hal/interfaces/joint_command_utils.h"
#include "intrinsic/icon/hal/interfaces/joint_limits.fbs.h"
#include "intrinsic/icon/hal/interfaces/joint_limits_utils.h"
#include "intrinsic/icon/hal/interfaces/joint_state.fbs.h"
#include "intrinsic/icon/hal/interfaces/joint_state_utils.h"
#include "intrinsic/icon/hal/interfaces/payload_command.fbs.h"
#include "intrinsic/icon/hal/interfaces/payload_command_utils.h"
#include "intrinsic/icon/hal/interfaces/payload_state.fbs.h"
#include "intrinsic/icon/hal/interfaces/payload_state_utils.h"
namespace intrinsic::icon {
namespace hardware_interface_traits {
INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::HardwareModuleState,
                                 intrinsic_fbs::BuildHardwareModuleState,
                                 "intrinsic_fbs.HardwareModuleState")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::ElectricalMotorCommand,
                                 intrinsic_fbs::BuildElectricalMotorCommand,
                                 "intrinsic_fbs.ElectricalMotorCommand")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::ElectricalMotorStatus,
                                 intrinsic_fbs::BuildElectricalMotorStatus,
                                 "intrinsic_fbs.ElectricalMotorStatus");

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::JointPositionCommand,
                                 intrinsic_fbs::BuildJointPositionCommand,
                                 "intrinsic_fbs.JointPositionCommand")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::JointCommandedPosition,
                                 intrinsic_fbs::BuildJointCommandedPosition,
                                 "intrinsic_fbs.JointCommandedPosition")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::JointVelocityCommand,
                                 intrinsic_fbs::BuildJointVelocityCommand,
                                 "intrinsic_fbs.JointVelocityCommand")

INTRINSIC_ADD_HARDWARE_INTERFACE(
    intrinsic_fbs::JointAccelerationAndTorqueCommand,
    intrinsic_fbs::BuildJointAccelerationAndTorqueCommand,
    "intrinsic_fbs.JointAccelerationAndTorqueCommand")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::JointTorqueCommand,
                                 intrinsic_fbs::BuildJointTorqueCommand,
                                 "intrinsic_fbs.JointTorqueCommand")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::JointPositionState,
                                 intrinsic_fbs::BuildJointPositionState,
                                 "intrinsic_fbs.JointPositionState")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::JointVelocityState,
                                 intrinsic_fbs::BuildJointVelocityState,
                                 "intrinsic_fbs.JointVelocityState")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::JointAccelerationState,
                                 intrinsic_fbs::BuildJointAccelerationState,
                                 "intrinsic_fbs.JointAccelerationState")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::JointTorqueState,
                                 intrinsic_fbs::BuildJointTorqueState,
                                 "intrinsic_fbs.JointTorqueState")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::JointLimits,
                                 intrinsic_fbs::BuildJointLimits,
                                 "intrinsic_fbs.JointLimits")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::ForceTorqueCommand,
                                 intrinsic_fbs::CreateFbsForceTorqueCommand,
                                 "intrinsic_fbs.ForceTorqueCommand")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::ForceTorqueStatus,
                                 intrinsic_fbs::CreateFbsForceTorqueStatus,
                                 "intrinsic_fbs.ForceTorqueStatus")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::SafetyStatusMessage,
                                 intrinsic_fbs::BuildSafetyStatusMessage,
                                 "intrinsic_fbs.SafetyStatus")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::DIOStatus,
                                 intrinsic_fbs::BuildDIOStatus,
                                 "intrinsic_fbs.DigitalInputStatus")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::AIOStatus,
                                 intrinsic_fbs::BuildAIOStatus,
                                 "intrinsic_fbs.AnalogInputStatus")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::DIOCommand,
                                 intrinsic_fbs::BuildDIOCommand,
                                 "intrinsic_fbs.DigitalOutputCommand")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::AIOCommand,
                                 intrinsic_fbs::BuildAIOCommand,
                                 "intrinsic_fbs.AnalogOutputCommand")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::Wrench,
                                 intrinsic_fbs::CreateWrenchBuffer,
                                 "intrinsic_fbs.Wrench");

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::HandGuidingCommand,
                                 intrinsic_fbs::BuildHandGuidingCommand,
                                 "intrinsic_fbs.HandGuidingCommand")

INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::ControlModeStatus,
                                 intrinsic_fbs::BuildControlModeStatus,
                                 "intrinsic_fbs.ControlModeStatus")

INTRINSIC_ADD_HARDWARE_INTERFACE(::intrinsic_fbs::PayloadCommand,
                                 intrinsic_fbs::BuildPayloadCommand,
                                 "intrinsic_fbs.PayloadCommand")

INTRINSIC_ADD_HARDWARE_INTERFACE(::intrinsic_fbs::PayloadState,
                                 intrinsic_fbs::BuildPayloadState,
                                 "intrinsic_fbs.PayloadState")
}  // namespace hardware_interface_traits
}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_HAL_DEFAULT_HARDWARE_INTERFACES_H_
