// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/simulation/gazebo/type_conversion.h"

#include "absl/status/statusor.h"
#include "gz/math/Pose3.hh"
#include "gz/math/Quaternion.hh"
#include "gz/math/Vector2.hh"
#include "gz/math/Vector3.hh"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/util/status/status_builder.h"

namespace intrinsic {

eigenmath::Vector2d GzToIntrinsic(const gz::math::Vector2d& vec) {
  return eigenmath::Vector2d(vec.X(), vec.Y());
}

eigenmath::Vector3d GzToIntrinsic(const gz::math::Vector3d& vec) {
  return eigenmath::Vector3d(vec.X(), vec.Y(), vec.Z());
}

eigenmath::Quaterniond GzToIntrinsic(const gz::math::Quaterniond& quat) {
  return eigenmath::Quaterniond(quat.W(), quat.X(), quat.Y(), quat.Z());
}

intrinsic::Pose3d GzToIntrinsic(const gz::math::Pose3d& pose) {
  return intrinsic::Pose3d(GzToIntrinsic(pose.Rot()),
                           GzToIntrinsic(pose.Pos()));
}

absl::StatusOr<intrinsic::Pose3d> GzToIntrinsicChecked(
    const gz::math::Pose3d& pose) {
  if (!pose.IsFinite()) {
    return intrinsic::InvalidArgumentErrorBuilder()
           << "Pose is not finite: " << pose;
  }

  return intrinsic::Pose3d(GzToIntrinsic(pose));
}

}  // namespace intrinsic
