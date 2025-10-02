// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/control/c_api/convert_c_types.h"

#include <cstring>
#include <optional>

#include "Eigen/Core"
#include "absl/log/check.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/control/c_api/c_types.h"
#include "intrinsic/icon/control/joint_position_command.h"
#include "intrinsic/icon/control/realtime_signal_types.h"
#include "intrinsic/kinematics/types/joint_limits.h"
#include "intrinsic/kinematics/types/joint_state.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/math/twist.h"

namespace intrinsic::icon {
namespace {
static_assert(
    kIntrinsicIconMaxNumberOfJoints == eigenmath::MAX_EIGEN_VECTOR_SIZE,
    "Mismatch between maximum size of C++ (intrinsic::eigenmath) and C "
    "vectors. This breaks the ICON C API!");
}

JointPositionCommand Convert(const IntrinsicIconJointPositionCommand& in) {
  CHECK(in.size < kIntrinsicIconMaxNumberOfJoints)
      << "IntrinsicIconJointPositionCommand has more than the maximum of "
      << kIntrinsicIconMaxNumberOfJoints << " joints.";
  eigenmath::VectorNd position_setpoints(
      Eigen::Map<const eigenmath::VectorNd>(in.position_setpoints, in.size));
  // Return .value() without checking because we guarantee that position,
  // velocity and acceleration are the same size here.
  return JointPositionCommand::Create(
             Eigen::Map<const eigenmath::VectorNd>(in.position_setpoints,
                                                   in.size),
             in.has_velocity_feedforwards
                 ? std::make_optional(eigenmath::VectorNd(
                       Eigen::Map<const eigenmath::VectorNd>(
                           in.velocity_feedforwards, in.size)))
                 : std::nullopt,
             in.has_acceleration_feedforwards
                 ? std::make_optional(eigenmath::VectorNd(
                       Eigen::Map<const eigenmath::VectorNd>(
                           in.acceleration_feedforwards, in.size)))
                 : std::nullopt)
      .value();
}

IntrinsicIconJointPositionCommand Convert(const JointPositionCommand& in) {
  CHECK(in.Size() < kIntrinsicIconMaxNumberOfJoints)
      << "JointPositionCommand has more than the maximum of "
      << kIntrinsicIconMaxNumberOfJoints << " joints.";
  IntrinsicIconJointPositionCommand out{.size = in.Size()};
  for (size_t i = 0; i < out.size; ++i) {
    out.position_setpoints[i] = in.position()(i);
    if (in.velocity_feedforward().has_value()) {
      out.has_velocity_feedforwards = true;
      out.velocity_feedforwards[i] = in.velocity_feedforward().value()(i);
    }
    if (in.acceleration_feedforward().has_value()) {
      out.has_acceleration_feedforwards = true;
      out.acceleration_feedforwards[i] =
          in.acceleration_feedforward().value()(i);
    }
  }
  return out;
}

JointLimits Convert(const IntrinsicIconJointLimits& in) {
  CHECK(in.size < kIntrinsicIconMaxNumberOfJoints)
      << "IntrinsicIconJointLimits has more than the maximum of "
      << kIntrinsicIconMaxNumberOfJoints << " joints.";
  return {
      .min_position =
          Eigen::Map<const eigenmath::VectorNd>(in.min_position, in.size),
      .max_position =
          Eigen::Map<const eigenmath::VectorNd>(in.max_position, in.size),
      .max_velocity =
          Eigen::Map<const eigenmath::VectorNd>(in.max_velocity, in.size),
      .max_acceleration =
          Eigen::Map<const eigenmath::VectorNd>(in.max_acceleration, in.size),
      .max_jerk = Eigen::Map<const eigenmath::VectorNd>(in.max_jerk, in.size),
      .max_torque =
          Eigen::Map<const eigenmath::VectorNd>(in.max_torque, in.size),
  };
}

IntrinsicIconJointLimits Convert(const JointLimits& in) {
  CHECK(in.size() < kIntrinsicIconMaxNumberOfJoints)
      << "JointLimits have more than the maximum of "
      << kIntrinsicIconMaxNumberOfJoints << " joints.";
  IntrinsicIconJointLimits out{.size = static_cast<size_t>(in.size())};
  for (size_t i = 0; i < out.size; ++i) {
    out.min_position[i] = in.min_position(i);
    out.max_position[i] = in.max_position(i);
    out.max_velocity[i] = in.max_velocity(i);
    out.max_acceleration[i] = in.max_acceleration(i);
    out.max_jerk[i] = in.max_jerk(i);
    out.max_torque[i] = in.max_torque(i);
  }
  return out;
}

JointStateP Convert(const IntrinsicIconJointStateP& in) {
  CHECK(in.size < kIntrinsicIconMaxNumberOfJoints)
      << "IntrinsicIconJointStateP has more than the maximum of "
      << kIntrinsicIconMaxNumberOfJoints << " joints.";
  return JointStateP(
      Eigen::Map<const eigenmath::VectorNd>(in.positions, in.size));
}

IntrinsicIconJointStateP Convert(const JointStateP& in) {
  CHECK(in.size() < kIntrinsicIconMaxNumberOfJoints)
      << "JointStateP has more than the maximum of "
      << kIntrinsicIconMaxNumberOfJoints << " joints.";
  IntrinsicIconJointStateP out{.size = static_cast<size_t>(in.size())};
  for (size_t i = 0; i < out.size; ++i) {
    out.positions[i] = in.position(i);
  }
  return out;
}

JointStateV Convert(const IntrinsicIconJointStateV& in) {
  CHECK(in.size < kIntrinsicIconMaxNumberOfJoints)
      << "IntrinsicIconJointStateV has more than the maximum of "
      << kIntrinsicIconMaxNumberOfJoints << " joints.";
  return JointStateV(
      Eigen::Map<const eigenmath::VectorNd>(in.velocities, in.size));
}

IntrinsicIconJointStateV Convert(const JointStateV& in) {
  CHECK(in.size() < kIntrinsicIconMaxNumberOfJoints)
      << "JointStateV has more than the maximum of "
      << kIntrinsicIconMaxNumberOfJoints << " joints.";
  IntrinsicIconJointStateV out{.size = static_cast<size_t>(in.size())};
  for (size_t i = 0; i < out.size; ++i) {
    out.velocities[i] = in.velocity(i);
  }
  return out;
}

JointStateA Convert(const IntrinsicIconJointStateA& in) {
  CHECK(in.size < kIntrinsicIconMaxNumberOfJoints)
      << "IntrinsicIconJointStateA has more than the maximum of "
      << kIntrinsicIconMaxNumberOfJoints << " joints.";
  return JointStateA(
      Eigen::Map<const eigenmath::VectorNd>(in.accelerations, in.size));
}

IntrinsicIconJointStateA Convert(const JointStateA& in) {
  CHECK(in.size() < kIntrinsicIconMaxNumberOfJoints)
      << "JointStateA has more than the maximum of "
      << kIntrinsicIconMaxNumberOfJoints << " joints.";
  IntrinsicIconJointStateA out{.size = static_cast<size_t>(in.size())};
  for (size_t i = 0; i < out.size; ++i) {
    out.accelerations[i] = in.acceleration(i);
  }

  return out;
}

eigenmath::Quaterniond Convert(const IntrinsicIconQuaternion& in) {
  return eigenmath::Quaterniond(/*w=*/in.w, /*x=*/in.x,
                                /*y=*/in.y, /*z=*/in.z);
}

IntrinsicIconQuaternion Convert(const eigenmath::Quaterniond& in) {
  return {.w = in.w(), .x = in.x(), .y = in.y(), .z = in.z()};
}

eigenmath::Vector3d Convert(const IntrinsicIconPoint& in) {
  return eigenmath::Vector3d(/*x=*/in.x, /*y=*/in.y, /*z=*/in.z);
}

IntrinsicIconPoint Convert(const eigenmath::Vector3d& in) {
  return {.x = in.x(), .y = in.y(), .z = in.z()};
}

Pose3d Convert(const IntrinsicIconPose3d& in) {
  return Pose3d(/*rotation=*/Convert(in.rotation),
                /*translation=*/Convert(in.translation));
}

IntrinsicIconPose3d Convert(const Pose3d& in) {
  return {.rotation = Convert(in.quaternion()),
          .translation = Convert(in.translation())};
}

Wrench Convert(const IntrinsicIconWrench& in) {
  return {
      in.x, in.y, in.z, in.rx, in.ry, in.rz,
  };
}

IntrinsicIconWrench Convert(const Wrench& in) {
  return {
      .x = in.x(),
      .y = in.y(),
      .z = in.z(),
      .rx = in.RX(),
      .ry = in.RY(),
      .rz = in.RZ(),
  };
}

eigenmath::Matrix6Nd Convert(const IntrinsicIconMatrix6Nd& in) {
  return Eigen::Map<const eigenmath::Matrix6Nd>(in.data, 6, in.num_cols);
}

IntrinsicIconMatrix6Nd Convert(const eigenmath::Matrix6Nd& in) {
  IntrinsicIconMatrix6Nd out;
  out.num_cols = in.cols();
  std::memcpy(out.data, in.data(), in.size() * sizeof(double));
  return out;
}

SignalValue Convert(const IntrinsicIconSignalValue& in) {
  SignalValue out;
  out.current_value = in.current_value;
  out.previous_value = in.previous_value;
  return out;
}

IntrinsicIconSignalValue Convert(const SignalValue& in) {
  IntrinsicIconSignalValue out;
  out.current_value = in.current_value;
  out.previous_value = in.previous_value;
  return out;
}

}  // namespace intrinsic::icon
