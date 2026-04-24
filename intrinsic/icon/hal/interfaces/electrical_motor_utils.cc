// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/interfaces/electrical_motor_utils.h"

#include <cstdint>

#include "absl/container/flat_hash_set.h"
#include "absl/strings/string_view.h"
#include "flatbuffers/detached_buffer.h"
#include "flatbuffers/flatbuffer_builder.h"
#include "intrinsic/icon/flatbuffers/fixed_string.h"
#include "intrinsic/icon/hal/interfaces/electrical_motor.fbs.h"
#include "intrinsic/icon/utils/log.h"

namespace intrinsic_fbs {
flatbuffers::DetachedBuffer BuildElectricalMotorStatus() {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  builder.Finish(CreateElectricalMotorStatus(builder));
  return builder.Release();
}

flatbuffers::DetachedBuffer BuildElectricalMotorCommand(
    const absl::flat_hash_set<MotorControlMode>& supported_modes) {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);

  ElectricalMotorCommandBuilder cmd_builder(builder);

  // Need to add these explicitly, even if they're the default values, even
  // though we set ForceDefaults(true). Otherwise, the builder will not add them
  // to the buffer. This is only true for builders, not for the generated
  // Create* functions.
  cmd_builder.add_control_mode(MotorControlMode::NONE);
  cmd_builder.add_control_word(0);

  for (const auto& mode : supported_modes) {
    switch (mode) {
      case MotorControlMode::NONE: {
        // There used to be a NoControlCommand, but its data (control mode and
        // control word) is now simply part of ElectricalMotorCommand itself.
        continue;
      }
      case MotorControlMode::PROFILE_POSITION: {
        PositionProfileCommand cmd;
        // This is fine, the builder makes a copy.
        cmd_builder.add_position_profile(&cmd);
        break;
      }
      case MotorControlMode::VELOCITY: {
        VelocityCommand cmd;
        cmd_builder.add_velocity(&cmd);
        break;
      }
      case MotorControlMode::VELOCITY_PROFILE: {
        VelocityProfileCommand cmd;
        cmd_builder.add_velocity_profile(&cmd);
        break;
      }
      case MotorControlMode::TORQUE_PROFILE: {
        TorqueProfileCommand cmd;
        cmd_builder.add_torque_profile(&cmd);
        break;
      }
      case MotorControlMode::HOME: {
        HomeCommand cmd;
        cmd_builder.add_home(&cmd);
        break;
      }
      case MotorControlMode::CYCLIC_POSITION: {
        CyclicPositionCommand cmd;
        cmd_builder.add_cyclic_position(&cmd);
        break;
      }
      case MotorControlMode::CYCLIC_VELOCITY: {
        CyclicVelocityCommand cmd;
        cmd_builder.add_cyclic_velocity(&cmd);
        break;
      }
      case MotorControlMode::CYCLIC_TORQUE: {
        CyclicTorqueCommand cmd;
        cmd_builder.add_cyclic_torque(&cmd);
        break;
      }
      case MotorControlMode::CYCLIC_JOINT_IMPEDANCE: {
        CyclicJointImpedanceCommand cmd;
        cmd_builder.add_cyclic_joint_impedance(&cmd);
        break;
      }
      case MotorControlMode::CYCLIC_POSITION_TORQUE: {
        CyclicPositionTorqueCommand cmd;
        cmd_builder.add_cyclic_position_torque(&cmd);
        break;
      }
      case MotorControlMode::MAX_CONTROL_MODES: {
        INTRINSIC_RT_LOG(ERROR)
            << "ElectricalMotorCommand with MAX_CONTROL_MODE support "
               "requested. That's not a real mode";
        break;
      }
    }
  }

  builder.Finish(cmd_builder.Finish());
  return builder.Release();
}

flatbuffers::DetachedBuffer BuildHomeCommand(int8_t method, double offset,
                                             double search_speed,
                                             double creep_speed,
                                             double acceleration) {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  builder.Finish(builder.CreateStruct(HomeCommand(
      /*_offset=*/offset, /*_search_speed=*/search_speed,
      /*_creep_speed=*/creep_speed,
      /*_acceleration=*/acceleration, /*_method=*/method)));
  return builder.Release();
}

flatbuffers::DetachedBuffer BuildHomingStatus(
    intrinsic_fbs::HomingStatusFlag state, absl::string_view error_message) {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  auto homing_status = HomingStatus();
  homing_status.mutate_state(state);
  intrinsic_fbs::StringCopy(&homing_status.mutable_error_message(),
                            error_message);
  builder.Finish(builder.CreateStruct(homing_status));
  return builder.Release();
}

bool SetMotorCommandCyclicPosition(ElectricalMotorCommand* command,
                                   double target_position,
                                   double velocity_feedforward,
                                   double acceleration_feedforward,
                                   double torque_feedforward) {
  auto cyclic_position = command->mutable_cyclic_position();
  if (cyclic_position == nullptr ||
      !command->mutate_control_mode(MotorControlMode::CYCLIC_POSITION)) {
    return false;
  }
  cyclic_position->mutate_target_position(target_position);
  cyclic_position->mutate_velocity_feedforward(velocity_feedforward);
  cyclic_position->mutate_acceleration_feedforward(acceleration_feedforward);
  cyclic_position->mutate_torque_feedforward(torque_feedforward);
  return true;
}

bool SetMotorCommandCyclicVelocity(ElectricalMotorCommand* command,
                                   double target_velocity,
                                   double torque_offset) {
  auto cyclic_velocity = command->mutable_cyclic_velocity();
  if (cyclic_velocity == nullptr ||
      !command->mutate_control_mode(MotorControlMode::CYCLIC_VELOCITY)) {
    return false;
  }
  cyclic_velocity->mutate_target_velocity(target_velocity);
  cyclic_velocity->mutate_torque_offset(torque_offset);
  return true;
}

bool SetMotorCommandCyclicTorque(ElectricalMotorCommand* command,
                                 double target_torque) {
  auto cyclic_torque = command->mutable_cyclic_torque();
  if (cyclic_torque == nullptr ||
      !command->mutate_control_mode(MotorControlMode::CYCLIC_TORQUE)) {
    return false;
  }
  cyclic_torque->mutate_target_torque(target_torque);
  return true;
}

bool SetMotorCommandCyclicJointImpedance(ElectricalMotorCommand* command,
                                         double target_stiffness,
                                         double target_damping,
                                         double target_position,
                                         double target_velocity,
                                         double target_torque) {
  auto cyclic_joint_impedance = command->mutable_cyclic_joint_impedance();
  if (cyclic_joint_impedance == nullptr ||
      !command->mutate_control_mode(MotorControlMode::CYCLIC_JOINT_IMPEDANCE)) {
    return false;
  }
  cyclic_joint_impedance->mutate_target_stiffness(target_stiffness);
  cyclic_joint_impedance->mutate_target_damping(target_damping);
  cyclic_joint_impedance->mutate_target_position(target_position);
  cyclic_joint_impedance->mutate_target_velocity(target_velocity);
  cyclic_joint_impedance->mutate_target_torque(target_torque);
  return true;
}

bool SetMotorCommandCyclicPositionTorque(ElectricalMotorCommand* command,
                                         double target_position,
                                         double velocity_offset,
                                         double target_torque) {
  auto cyclic_position_torque = command->mutable_cyclic_position_torque();
  if (cyclic_position_torque == nullptr ||
      !command->mutate_control_mode(MotorControlMode::CYCLIC_POSITION_TORQUE)) {
    return false;
  }
  cyclic_position_torque->mutate_target_position(target_position);
  cyclic_position_torque->mutate_velocity_offset(velocity_offset);
  cyclic_position_torque->mutate_target_torque(target_torque);

  return true;
}

bool SupportsMode(const ElectricalMotorCommand& command,
                  MotorControlMode mode) {
  switch (mode) {
    case MotorControlMode::NONE: {
      // There used to be a NoControlCommand, but its data (control mode and
      // control word) is now simply part of ElectricalMotorCommand itself.
      return true;
    }
    case MotorControlMode::PROFILE_POSITION: {
      return command.position_profile() != nullptr;
    }
    case MotorControlMode::VELOCITY: {
      return command.velocity() != nullptr;
    }
    case MotorControlMode::VELOCITY_PROFILE: {
      return command.velocity_profile() != nullptr;
    }
    case MotorControlMode::TORQUE_PROFILE: {
      return command.torque_profile() != nullptr;
    }
    case MotorControlMode::HOME: {
      return command.home() != nullptr;
    }
    case MotorControlMode::CYCLIC_POSITION: {
      return command.cyclic_position() != nullptr;
    }
    case MotorControlMode::CYCLIC_VELOCITY: {
      return command.cyclic_velocity() != nullptr;
    }
    case MotorControlMode::CYCLIC_TORQUE: {
      return command.cyclic_torque() != nullptr;
    }
    case MotorControlMode::CYCLIC_JOINT_IMPEDANCE: {
      return command.cyclic_joint_impedance() != nullptr;
    }
    case MotorControlMode::CYCLIC_POSITION_TORQUE: {
      return command.cyclic_position_torque() != nullptr;
    }
    case MotorControlMode::MAX_CONTROL_MODES: {
      INTRINSIC_RT_LOG(ERROR)
          << "ElectricalMotorCommand with MAX_CONTROL_MODE support "
             "requested. That's not a real mode";
      return false;
    }
    default: {
      return false;
    }
  }
}

// Sets the boolean values in ElectricalMotorStatus based on the status word.
void SetElectricalMotorStatus(ElectricalMotorStatus* motor_status,
                              uint16_t status_word) {
  for (const auto& flag : EnumValuesOperationModeStatusFlags()) {
    bool flag_value = ToFlag(flag) & status_word;

    switch (flag) {
      case OperationModeStatusFlags::READY_TO_SWITCH_ON:
        motor_status->mutate_ready_to_switch_on(flag_value);
        break;
      case OperationModeStatusFlags::SWITCHED_ON:
        motor_status->mutate_switched_on(flag_value);
        break;
      case OperationModeStatusFlags::OPERATION_ENABLED:
        motor_status->mutate_operation_enabled(flag_value);
        break;
      case OperationModeStatusFlags::FAULT:
        motor_status->mutate_in_fault(flag_value);
        break;
      case OperationModeStatusFlags::VOLTAGE_ENABLED:
        motor_status->mutate_voltage_enabled(flag_value);
        break;
      case OperationModeStatusFlags::QUICK_STOP:
        motor_status->mutate_in_quick_stop(flag_value);
        break;
      case OperationModeStatusFlags::TARGET_REACHED:
        motor_status->mutate_target_reached(flag_value);
        break;
      case OperationModeStatusFlags::INTERNAL_LIMIT_ACTIVE:
        motor_status->mutate_internal_limit_active(flag_value);
        break;
      case OperationModeStatusFlags::SWITCH_ON_DISABLED:
      case OperationModeStatusFlags::WARNING:
      case OperationModeStatusFlags::REMOTE:
        // Ignore these bits for now. There's no code that reads them.
        break;
    }
  }
}
bool CommandCorrectlySet(const ElectricalMotorCommand& motor_command) {
  switch (motor_command.control_mode()) {
    case MotorControlMode::NONE:
      return true;
    case MotorControlMode::PROFILE_POSITION:
      return motor_command.position_profile() != nullptr;
    case MotorControlMode::VELOCITY:
      return motor_command.velocity() != nullptr;
    case MotorControlMode::VELOCITY_PROFILE:
      return motor_command.velocity_profile() != nullptr;
    case MotorControlMode::TORQUE_PROFILE:
      return motor_command.torque_profile() != nullptr;
    case MotorControlMode::HOME:
      return motor_command.home() != nullptr;
    case MotorControlMode::CYCLIC_POSITION:
      return motor_command.cyclic_position() != nullptr;
    case MotorControlMode::CYCLIC_VELOCITY:
      return motor_command.cyclic_velocity() != nullptr;
    case MotorControlMode::CYCLIC_TORQUE:
      return motor_command.cyclic_torque() != nullptr;
    case MotorControlMode::CYCLIC_JOINT_IMPEDANCE:
      return motor_command.cyclic_joint_impedance() != nullptr;
    case MotorControlMode::CYCLIC_POSITION_TORQUE:
      return motor_command.cyclic_position_torque() != nullptr;
    default:
      INTRINSIC_RT_LOG(ERROR)
          << "Unsupported control mode: "
          << EnumNameMotorControlMode(motor_command.control_mode());
      return false;
  }
}
}  // namespace intrinsic_fbs
