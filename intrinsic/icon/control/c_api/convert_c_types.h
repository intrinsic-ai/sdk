// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CONTROL_C_API_CONVERT_C_TYPES_H_
#define INTRINSIC_ICON_CONTROL_C_API_CONVERT_C_TYPES_H_

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/control/c_api/c_types.h"
#include "intrinsic/icon/control/joint_position_command.h"
#include "intrinsic/icon/control/realtime_signal_types.h"
#include "intrinsic/kinematics/types/joint_limits.h"
#include "intrinsic/kinematics/types/joint_state.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/math/twist.h"

namespace intrinsic::icon {

// These helpers convert back and forth between ICON C++ vocabulary types and
// their C API equivalents.
//
// Note that these functions do CHECKs to ensure that sizes are compatible and
// will crash if the input has (or claims to have) more values than the output
// type can handle. These checks are *not* exposed via (Realtime)Status return
// values because a static_assert in convert_c_types.cc ensures that the maximum
// size of C and C++ vectors is the same.
//
// If a user puts an invalid value into the `size` member of one of the C data
// structs, that's on them.

JointPositionCommand Convert(const IntrinsicIconJointPositionCommand& in);
IntrinsicIconJointPositionCommand Convert(const JointPositionCommand& in);

JointLimits Convert(const IntrinsicIconJointLimits& in);
IntrinsicIconJointLimits Convert(const JointLimits& in);

JointStateP Convert(const IntrinsicIconJointStateP& in);
IntrinsicIconJointStateP Convert(const JointStateP& in);

JointStateV Convert(const IntrinsicIconJointStateV& in);
IntrinsicIconJointStateV Convert(const JointStateV& in);

JointStateA Convert(const IntrinsicIconJointStateA& in);
IntrinsicIconJointStateA Convert(const JointStateA& in);

eigenmath::Quaterniond Convert(const IntrinsicIconQuaternion& in);
IntrinsicIconQuaternion Convert(const eigenmath::Quaterniond& in);

eigenmath::Vector3d Convert(const IntrinsicIconPoint& in);
IntrinsicIconPoint Convert(const eigenmath::Vector3d& in);

Pose3d Convert(const IntrinsicIconPose3d& in);
IntrinsicIconPose3d Convert(const Pose3d& in);

Wrench Convert(const IntrinsicIconWrench& in);
IntrinsicIconWrench Convert(const Wrench& in);

eigenmath::Matrix6Nd Convert(const IntrinsicIconMatrix6Nd& in);
IntrinsicIconMatrix6Nd Convert(const eigenmath::Matrix6Nd& in);

SignalValue Convert(const IntrinsicIconSignalValue& in);
IntrinsicIconSignalValue Convert(const SignalValue& in);

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_CONTROL_C_API_CONVERT_C_TYPES_H_
