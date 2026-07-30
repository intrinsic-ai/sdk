// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SIMULATION_GAZEBO_TYPE_CONVERSION_H_
#define INTRINSIC_SIMULATION_GAZEBO_TYPE_CONVERSION_H_

#include "absl/status/statusor.h"
#include "gz/math/Pose3.hh"
#include "gz/math/Quaternion.hh"
#include "gz/math/Vector2.hh"
#include "gz/math/Vector3.hh"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"

namespace intrinsic {

eigenmath::Vector2d GzToIntrinsic(const gz::math::Vector2d& vec);

eigenmath::Vector3d GzToIntrinsic(const gz::math::Vector3d& vec);

eigenmath::Quaterniond GzToIntrinsic(const gz::math::Quaterniond& quat);

intrinsic::Pose3d GzToIntrinsic(const gz::math::Pose3d& pose);

// Returns the intrinsic pose for the given gazebo pose.
// Returns an InvalidArgumentError if the input pose has NaNs.
absl::StatusOr<intrinsic::Pose3d> GzToIntrinsicChecked(
    const gz::math::Pose3d& pose);

}  // namespace intrinsic

#endif  // INTRINSIC_SIMULATION_GAZEBO_TYPE_CONVERSION_H_
