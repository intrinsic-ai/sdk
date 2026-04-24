// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/robot_arm/robot_arm_bus_component.h"

#include <cstddef>
#include <cstdint>
#include <functional>
#include <memory>
#include <string>
#include <utility>
#include <vector>

#include "absl/memory/memory.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "intrinsic/icon/hal/default_hardware_interfaces.h"  // IWYU pragma: keep
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/hardware_interface_registry.h"
#include "intrinsic/icon/hal/interfaces/joint_command.fbs.h"
#include "intrinsic/icon/hal/interfaces/joint_state.fbs.h"
#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"
#include "intrinsic/icon/hal/lib/robot_arm/v1/robot_arm_bus_component_config.pb.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_macro.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::robot_arm {

/* static */
absl::StatusOr<std::unique_ptr<RobotArmBusComponent>>
RobotArmBusComponent::Create(
    fieldbus::DeviceInitContext& init_context,
    const intrinsic_proto::icon::v1::RobotArmBusComponentConfig& config) {
  const auto number_of_joints = config.position_command_variables_size();
  if (number_of_joints != config.position_state_variables_size() ||
      (config.velocity_state_variables_size() > 0 &&
       number_of_joints != config.velocity_state_variables_size()) ||
      (config.acceleration_state_variables_size() > 0 &&
       number_of_joints != config.acceleration_state_variables_size()) ||
      (config.feedforward_velocity_command_variables_size() > 0 &&
       number_of_joints !=
           config.feedforward_velocity_command_variables_size())) {
    return absl::InvalidArgumentError(absl::StrCat(
        "Number of joints mismatch: Size of `velocity_state_variables` (",
        config.velocity_state_variables_size(),
        ") and `acceleration_state_variables` (",
        config.acceleration_state_variables_size(),
        ") and `feedforward_velocity_command_variables` (",
        config.feedforward_velocity_command_variables_size(),
        ") must either be zero or identical to `position_state_variables` (",
        config.position_state_variables_size(), ")."));
  }

  std::string joint_position_command_name = "joint_position_command";
  std::string joint_position_state_name = "joint_position_state";
  std::string joint_velocity_state_name = "joint_velocity_state";
  std::string joint_acceleration_state_name = "joint_acceleration_state";

  if (!config.device_prefix().empty()) {
    joint_position_command_name =
        absl::StrCat(config.device_prefix(), "_", joint_position_command_name);
    joint_position_state_name =
        absl::StrCat(config.device_prefix(), "_", joint_position_state_name);
    joint_velocity_state_name =
        absl::StrCat(config.device_prefix(), "_", joint_velocity_state_name);
    joint_acceleration_state_name = absl::StrCat(config.device_prefix(), "_",
                                                 joint_acceleration_state_name);
  }

  // Advertise the position command and state interfaces. The velocity and
  // acceleration interfaces are optional (see below).
  intrinsic::icon::HardwareInterfaceRegistry& interface_registry =
      init_context.GetInterfaceRegistry();
  INTR_ASSIGN_OR_RETURN(
      auto position_command,
      interface_registry
          .AdvertiseStrictInterface<intrinsic_fbs::JointPositionCommand>(
              joint_position_command_name, number_of_joints));
  INTR_ASSIGN_OR_RETURN(
      auto position_state,
      interface_registry
          .AdvertiseMutableInterface<intrinsic_fbs::JointPositionState>(
              joint_position_state_name, number_of_joints));

  intrinsic::icon::MutableHardwareInterfaceHandle<
      intrinsic_fbs::JointVelocityState>
      velocity_state;
  intrinsic::icon::MutableHardwareInterfaceHandle<
      intrinsic_fbs::JointAccelerationState>
      acceleration_state;

  // Create scaled position state getters.
  const auto& variable_registry = init_context.GetVariableRegistry();
  std::vector<std::function<double()>> scaled_position_state_getters;
  for (const auto& state_var : config.position_state_variables()) {
    INTR_ASSIGN_OR_RETURN(
        intrinsic::fieldbus::ProcessVariable var,
        variable_registry.GetInputVariable(state_var.variable_name()));
    const double scale = state_var.scale();

    // Try the different supported bus variable data types.
    if (var.IsCompatibleType<int32_t>().ok()) {
      scaled_position_state_getters.emplace_back([var, scale]() -> double {
        return var.ReadUnchecked<int32_t>() * scale;
      });
    } else if (var.IsCompatibleType<double>().ok()) {
      scaled_position_state_getters.emplace_back([var, scale]() -> double {
        return var.ReadUnchecked<double>() * scale;
      });
    } else {
      return absl::InvalidArgumentError(absl::StrCat(
          "Unsupported bus variable data type for variable: ",
          state_var.variable_name(), " (expected double or int32_t)"));
    }
  }
  // Create scaled velocity state getters.
  std::vector<std::function<double()>> scaled_velocity_state_getters;
  for (const auto& state_var : config.velocity_state_variables()) {
    INTR_ASSIGN_OR_RETURN(
        intrinsic::fieldbus::ProcessVariable var,
        variable_registry.GetInputVariable(state_var.variable_name()));
    const double scale = state_var.scale();

    // Try the different supported bus variable data types.
    if (var.IsCompatibleType<int32_t>().ok()) {
      scaled_velocity_state_getters.emplace_back([var, scale]() -> double {
        return var.ReadUnchecked<int32_t>() * scale;
      });
    } else if (var.IsCompatibleType<double>().ok()) {
      scaled_velocity_state_getters.emplace_back([var, scale]() -> double {
        return var.ReadUnchecked<double>() * scale;
      });
    } else {
      return absl::InvalidArgumentError(absl::StrCat(
          "Unsupported bus variable data type for variable: ",
          state_var.variable_name(), " (expected double or int32_t)"));
    }
  }
  // Create scaled acceleration state getters.
  std::vector<std::function<double()>> scaled_acceleration_state_getters;
  for (const auto& state_var : config.acceleration_state_variables()) {
    INTR_ASSIGN_OR_RETURN(
        intrinsic::fieldbus::ProcessVariable var,
        variable_registry.GetInputVariable(state_var.variable_name()));
    const double scale = state_var.scale();

    // Try the different supported bus variable data types.
    if (var.IsCompatibleType<int32_t>().ok()) {
      scaled_acceleration_state_getters.emplace_back([var, scale]() -> double {
        return var.ReadUnchecked<int32_t>() * scale;
      });
    } else if (var.IsCompatibleType<double>().ok()) {
      scaled_acceleration_state_getters.emplace_back([var, scale]() -> double {
        return var.ReadUnchecked<double>() * scale;
      });
    } else {
      return absl::InvalidArgumentError(absl::StrCat(
          "Unsupported bus variable data type for variable: ",
          state_var.variable_name(), " (expected double or int32_t)"));
    }
  }

  // If velocity state variables are provided, advertise the velocity interface.
  if (!scaled_velocity_state_getters.empty()) {
    INTR_ASSIGN_OR_RETURN(
        velocity_state,
        interface_registry
            .AdvertiseMutableInterface<intrinsic_fbs::JointVelocityState>(
                joint_velocity_state_name, number_of_joints));
  }

  // If acceleration state variables are provided, advertise the acceleration
  // interface.
  if (!scaled_acceleration_state_getters.empty()) {
    INTR_ASSIGN_OR_RETURN(
        acceleration_state,
        interface_registry
            .AdvertiseMutableInterface<intrinsic_fbs::JointAccelerationState>(
                joint_acceleration_state_name, number_of_joints));
  }

  // Create scaled position command setters.
  std::vector<std::function<void(double)>> scaled_position_command_setters;
  for (const auto& state_var : config.position_command_variables()) {
    INTR_ASSIGN_OR_RETURN(
        intrinsic::fieldbus::ProcessVariable var,
        variable_registry.GetOutputVariable(state_var.variable_name()));
    const double scale = state_var.scale();

    // Try the different supported bus variable data types.
    if (var.IsCompatibleType<int32_t>().ok()) {
      scaled_position_command_setters.emplace_back(
          [var, scale](double value) mutable -> void {
            var.WriteUnchecked(static_cast<int32_t>(value * scale));
          });

    } else if (var.IsCompatibleType<double>().ok()) {
      scaled_position_command_setters.emplace_back(
          [var, scale](double value) mutable -> void {
            var.WriteUnchecked(static_cast<double>(value * scale));
          });
    } else {
      return absl::InvalidArgumentError(
          absl::StrCat("Unsupported bus variable data type for variable "
                       "(expected double or int32_t): ",
                       state_var.variable_name()));
    }
  }

  // Create scaled feedforward velocity command setters.
  std::vector<std::function<void(double)>>
      scaled_feedforward_velocity_command_setters;
  for (const auto& state_var :
       config.feedforward_velocity_command_variables()) {
    INTR_ASSIGN_OR_RETURN(
        intrinsic::fieldbus::ProcessVariable var,
        variable_registry.GetOutputVariable(state_var.variable_name()));
    const double scale = state_var.scale();

    // Try the different supported bus variable data types.
    if (var.IsCompatibleType<int32_t>().ok()) {
      scaled_feedforward_velocity_command_setters.emplace_back(
          [var, scale](double value) mutable -> void {
            var.WriteUnchecked(static_cast<int32_t>(value * scale));
          });

    } else if (var.IsCompatibleType<double>().ok()) {
      scaled_feedforward_velocity_command_setters.emplace_back(
          [var, scale](double value) mutable -> void {
            var.WriteUnchecked(static_cast<double>(value * scale));
          });
    } else {
      return absl::InvalidArgumentError(
          absl::StrCat("Unsupported bus variable data type for variable "
                       "(expected double or int32_t): ",
                       state_var.variable_name()));
    }
  }

  return absl::WrapUnique(new RobotArmBusComponent(
      std::move(position_state), std::move(velocity_state),
      std::move(acceleration_state), std::move(position_command),
      scaled_position_state_getters, scaled_velocity_state_getters,
      scaled_acceleration_state_getters, scaled_position_command_setters,
      scaled_feedforward_velocity_command_setters));
}

RobotArmBusComponent::RobotArmBusComponent(
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
        scaled_feedforward_velocity_command_setters)
    : position_state_(std::move(position_state)),
      velocity_state_(std::move(velocity_state)),
      acceleration_state_(std::move(acceleration_state)),
      position_command_(std::move(position_command)),
      scaled_position_state_getters_(scaled_position_state_getters),
      scaled_velocity_state_getters_(scaled_velocity_state_getters),
      scaled_acceleration_state_getters_(scaled_acceleration_state_getters),
      scaled_position_command_setters_(scaled_position_command_setters),
      scaled_feedforward_velocity_command_setters_(
          scaled_feedforward_velocity_command_setters) {}

intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus>
RobotArmBusComponent::CyclicRead(fieldbus::RequestType) {
  for (std::size_t i = 0; i < scaled_position_state_getters_.size(); ++i) {
    position_state_->mutable_position()->Mutate(
        i, scaled_position_state_getters_[i]());
  }
  // Only "read" the velocity state if the velocity state interfaces has been
  // advertised.
  if (*velocity_state_ != nullptr) {
    for (std::size_t i = 0; i < scaled_velocity_state_getters_.size(); ++i) {
      velocity_state_->mutable_velocity()->Mutate(
          i, scaled_velocity_state_getters_[i]());
    }
  }
  // Only "read" the acceleration state if the acceleration state interfaces has
  // been advertised.
  if (*acceleration_state_ != nullptr) {
    for (std::size_t i = 0; i < scaled_acceleration_state_getters_.size();
         ++i) {
      acceleration_state_->mutable_acceleration()->Mutate(
          i, scaled_acceleration_state_getters_[i]());
    }
  }
  // Set the position command to the sensed position as a precaution, for when
  // `Write` is not (yet) called (b/297333030).
  for (std::size_t i = 0; i < scaled_position_command_setters_.size(); ++i) {
    scaled_position_command_setters_[i](position_state_->position()->Get(i));
  }
  return fieldbus::RequestStatus::kDone;
}

intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus>
RobotArmBusComponent::CyclicWrite(fieldbus::RequestType) {
  if (!position_command_.Value().ok()) {
    // Command the sensed position as a precaution, for when `Write` is called
    // but the position command has not been updated or was never set.
    // (b/320738872).
    for (std::size_t i = 0; i < scaled_position_command_setters_.size(); ++i) {
      scaled_position_command_setters_[i](position_state_->position()->Get(i));
    }
  } else {
    INTRINSIC_RT_ASSIGN_OR_RETURN(const auto joint_position_command,
                                  position_command_.Value());
    for (std::size_t i = 0; i < scaled_position_command_setters_.size(); ++i) {
      scaled_position_command_setters_[i](
          joint_position_command->position()->Get(i));
    }
    for (std::size_t i = 0;
         i < scaled_feedforward_velocity_command_setters_.size(); ++i) {
      scaled_feedforward_velocity_command_setters_[i](
          joint_position_command->velocity_feedforward()->Get(i));
    }
  }

  return fieldbus::RequestStatus::kDone;
}

}  // namespace intrinsic::robot_arm
