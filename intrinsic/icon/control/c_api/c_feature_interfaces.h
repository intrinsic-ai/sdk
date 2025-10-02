// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CONTROL_C_API_C_FEATURE_INTERFACES_H_
#define INTRINSIC_ICON_CONTROL_C_API_C_FEATURE_INTERFACES_H_

#include "intrinsic/icon/control/c_api/c_realtime_status.h"
#include "intrinsic/icon/control/c_api/c_types.h"

#ifdef __cplusplus
extern "C" {
#endif

/////////////////////////////////////////////////
// JointPositionCommandInterface FeatureInterface
/////////////////////////////////////////////////
struct IntrinsicIconFeatureInterfaceJointPositionCommandInterface;

struct IntrinsicIconFeatureInterfaceJointPositionCommandInterfaceVtable {
  // Sends the given position setpoints with optional velocity and acceleration
  // feedforward values. The caller is responsible for ensuring that velocity
  // and acceleration values (if provided) are kinematically consistent.
  //
  // Returns an error if the setpoints are invalid, i.e. they contain the wrong
  // number of values or violate position, velocity or acceleration limits.
  //
  // Returns an error if the part is currently not in position mode.
  IntrinsicIconRealtimeStatus (*set_position_setpoints)(
      IntrinsicIconFeatureInterfaceJointPositionCommandInterface* self,
      const IntrinsicIconJointPositionCommand* const setpoints);

  // Returns the setpoints from the previous control cycle.
  IntrinsicIconJointPositionCommand (*previous_position_setpoints)(
      const IntrinsicIconFeatureInterfaceJointPositionCommandInterface* self);
};

/////////////////////////////////////////////////
// JointPositionSensor FeatureInterface
/////////////////////////////////////////////////
struct IntrinsicIconFeatureInterfaceJointPositionSensor;

struct IntrinsicIconFeatureInterfaceJointPositionSensorVtable {
  // Returns the current joint positions in radians.
  IntrinsicIconJointStateP (*get_sensed_position)(
      const IntrinsicIconFeatureInterfaceJointPositionSensor* self);
};

/////////////////////////////////////////////////
// JointVelocityEstimator FeatureInterface
/////////////////////////////////////////////////
struct IntrinsicIconFeatureInterfaceJointVelocityEstimator;

struct IntrinsicIconFeatureInterfaceJointVelocityEstimatorVtable {
  // Returns the current joint velocity estimates in radians per second.
  IntrinsicIconJointStateV (*get_velocity_estimate)(
      const IntrinsicIconFeatureInterfaceJointVelocityEstimator* self);
};

/////////////////////////////////////////////////
// JointLimits FeatureInterface
/////////////////////////////////////////////////
struct IntrinsicIconFeatureInterfaceJointLimits;

struct IntrinsicIconFeatureInterfaceJointLimitsVtable {
  // Returns the application limits. Actions should use these joint limits by
  // default and reject any user commands that exceed them.
  IntrinsicIconJointLimits (*get_application_limits)(
      const IntrinsicIconFeatureInterfaceJointLimits* self);
  // Returns the system limits. An action must not command motions beyond these
  // limits. ICON monitors this at a low level and faults if the maximum limits
  // are violated.
  IntrinsicIconJointLimits (*get_system_limits)(
      const IntrinsicIconFeatureInterfaceJointLimits* self);
};

/////////////////////////////////////////////////
// ForceTorqueSensor FeatureInterface
/////////////////////////////////////////////////
struct IntrinsicIconFeatureInterfaceForceTorqueSensor;

struct IntrinsicIconFeatureInterfaceForceTorqueSensorVtable {
  // Filtered wrench at the tip compensated for support mass and bias.
  IntrinsicIconWrench (*wrench_at_tip)(
      const IntrinsicIconFeatureInterfaceForceTorqueSensor* self);
  // Requests a taring of the sensor. ICON will apply the current filtered
  // sensor reading as a bias to all future readings. That is, if the forces
  // acting on the sensor do not change, a call to WrenchAtTip() in the next
  // control cycle will return an all-zero Wrench.
  IntrinsicIconRealtimeStatus (*tare)(
      IntrinsicIconFeatureInterfaceForceTorqueSensor* self);
};

/////////////////////////////////////////////////
// ManipulatorKinematics FeatureInterface
/////////////////////////////////////////////////
struct IntrinsicIconFeatureInterfaceManipulatorKinematics;

struct IntrinsicIconFeatureInterfaceManipulatorKinematicsVtable {
  // Writes the base to tip Jacobian for the given `dof_positions` into
  // `jacobian_out`. Assumes the kinematic model is a chain and returns an error
  // otherwise.
  // Caller owns `jacobian_out`.
  IntrinsicIconRealtimeStatus (*compute_chain_jacobian)(
      const IntrinsicIconFeatureInterfaceManipulatorKinematics* self,
      const IntrinsicIconJointStateP* dof_positions,
      IntrinsicIconMatrix6Nd* jacobian_out);
  // Writes the base to tip transform for the given `dof_positions` into
  // `pose_out`. Assumes the kinematic model is a chain and returns an error
  // otherwise.
  IntrinsicIconRealtimeStatus (*compute_chain_fk)(
      const IntrinsicIconFeatureInterfaceManipulatorKinematics* self,
      const IntrinsicIconJointStateP* dof_positions,
      IntrinsicIconPose3d* pose_out);
};

// Holds pointers to the feature interfaces for a given Slot. If the Part
// assigned to that Slot does not implement an interface, the corresponding
// member is set to nullptr.
//
struct IntrinsicIconFeatureInterfacesForSlot {
  IntrinsicIconFeatureInterfaceJointPositionCommandInterface* joint_position;
  IntrinsicIconFeatureInterfaceJointPositionSensor* joint_position_sensor;
  IntrinsicIconFeatureInterfaceJointVelocityEstimator* joint_velocity_estimator;
  IntrinsicIconFeatureInterfaceJointLimits* joint_limits;
  IntrinsicIconFeatureInterfaceManipulatorKinematics* manipulator_kinematics;
  IntrinsicIconFeatureInterfaceForceTorqueSensor* force_torque_sensor;
};

// Same as above, but holds const pointers. This prevents Actions from sending
// commands to a Feature Interface when they should not be able to do so.
struct IntrinsicIconConstFeatureInterfacesForSlot {
  const IntrinsicIconFeatureInterfaceJointPositionCommandInterface*
      joint_position;
  const IntrinsicIconFeatureInterfaceJointPositionSensor* joint_position_sensor;
  const IntrinsicIconFeatureInterfaceJointVelocityEstimator*
      joint_velocity_estimator;
  const IntrinsicIconFeatureInterfaceJointLimits* joint_limits;
  const IntrinsicIconFeatureInterfaceManipulatorKinematics*
      manipulator_kinematics;
  const IntrinsicIconFeatureInterfaceForceTorqueSensor* force_torque_sensor;
};

// Holds function pointers to the functions for each FeatureInterface. Plugin
// code can then call those functions like this
//
// IntrinsicIconFeatureInterfaceVtable feature_interfaces =
//   server_functions.feature_interfaces;
// IntrinsicIconFeatureInterfaceJointPositionCommandInterface*
// joint_position_interface =
//   GetJointPositionCommandInterfaceFromSomewhere();
// IntrinsicIconJointPositionCommand command;
// command.position_setpoints[0]= 1.0;
// command.position_setpoints[1]= 1.5;
// command.position_setpoints[2]= 0.7;
// command.position_setpoints_size = 0;
// command.velocity_feedforwards_size = 0;
// command.acceleration_feedforwards_size = 0;
//
// IntrinsicIconRealtimeStatus result =
// feature_interfaces.joint_position.set_position_setpoints(
//    joint_position_interface, &command);
struct IntrinsicIconFeatureInterfaceVtable {
  IntrinsicIconFeatureInterfaceJointPositionCommandInterfaceVtable
      joint_position;
  IntrinsicIconFeatureInterfaceJointPositionSensorVtable joint_position_sensor;
  IntrinsicIconFeatureInterfaceJointVelocityEstimatorVtable
      joint_velocity_estimator;
  IntrinsicIconFeatureInterfaceJointLimitsVtable joint_limits;
  IntrinsicIconFeatureInterfaceManipulatorKinematicsVtable
      manipulator_kinematics;
  IntrinsicIconFeatureInterfaceForceTorqueSensorVtable force_torque_sensor;
};
#ifdef __cplusplus
}
#endif

#endif  // INTRINSIC_ICON_CONTROL_C_API_C_FEATURE_INTERFACES_H_
