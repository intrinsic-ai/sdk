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

#ifndef INTRINSIC_EIGENMATH_CERES_POSES_H_
#define INTRINSIC_EIGENMATH_CERES_POSES_H_

#include "Eigen/Core"
#include "Eigen/Geometry"
#include "ceres/ceres.h"
#include "intrinsic/eigenmath/so2.h"
#include "intrinsic/math/pose3.h"

namespace intrinsic {
namespace eigenmath {

namespace details {
template <typename Scalar, int N>
struct IsFinite<ceres::Jet<Scalar, N>> {
  bool operator()(const ceres::Jet<Scalar, N>& value) const {
    using ceres::IsFinite;
    return IsFinite(value);
  }
};
}  // namespace details

namespace ceresutils {

// Mutable access to internal quaternion storage of Pose3.
//
// The user (e.g. Ceres) must take precaution that the quaternion
// stays normalized.
//
// `a_pose_b`: A 3D pose.
// Returns pointer to mutable data.
template <class T, int Options>
double* MutableRotationData(Pose3<T, Options>& a_pose_b) {
  // This const-cast is allowed, since pose is non-const.
  return const_cast<double*>(a_pose_b.so3().quaternion().coeffs().data());
}

// Generates a 3d pose given views to a unit quaternion and translation.
//
// `a_q_b`: quaternion view wrapped around raw data.
// `t_a`: translation view wrapped around raw data.
// Returns a 3d pose.
template <class T, int OptionsQuaternion, int OptionsTranslation>
Pose3<T> MakePose3FromViews(
    const Eigen::Map<const Vector4<T, OptionsQuaternion>>& a_q_b,
    const Eigen::Map<const Vector3<T, OptionsTranslation>>& t_a) {
  return Pose3<T>(SO3<T>(Eigen::Quaternion<T>(a_q_b)), t_a.eval());
}

}  // namespace ceresutils
}  // namespace eigenmath
}  // namespace intrinsic

#endif  // INTRINSIC_EIGENMATH_CERES_POSES_H_
