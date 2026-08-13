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
