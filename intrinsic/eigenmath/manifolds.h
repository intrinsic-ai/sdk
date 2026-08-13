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

#ifndef INTRINSIC_EIGENMATH_MANIFOLDS_H_
#define INTRINSIC_EIGENMATH_MANIFOLDS_H_

#include "Eigen/Core"
#include "intrinsic/eigenmath/so3.h"
#include "intrinsic/eigenmath/su2_impl.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"

namespace intrinsic {
namespace eigenmath {

/**
 * @brief Exponential map of the special orthogonal group SO(3)
 *
 * It maps a tangent element of SO(3) at the identity, which is equivalent to
 * the rotation vector (axis times angle), to the corresponing rotation in
 * SO(3). Here we choose to represent the rotation as a unit quaternion instead
 * of a rotation matrix.
 *
 * Let \f$ \exp \f$ be the matrix exponential, \f$ R \f$ be a function which
 * maps a quaternion to the corresponing rotation matrix and \f$ \hat{\cdot} \f$
 * be the function which maps a 3-vector the corresponding skew-symmetric (3x3)
 * matrix representation. It holds hat:
 *
 * \f$ R(\exp_{SO3}(\delta)) = exp(\hat{delta}) \f$
 *
 * @param delta rotation vector (axis times angle)
 * @return corresponding unit quaternion representing rotation in 3D.
 */
template <class T, int Options = Eigen::AutoAlign>
SO3<T, Options> expSO3(const Vector3<T, Options>& delta) {
  //  SU(2) is a double cover of SO(3), thus we have to half the tangent vector
  //  delta
  const Vector3<T, Options> half_delta = T(0.5) * delta;
  return SO3<T, Options>(details::expSU2Impl(half_delta), false);
}

/**
 * @brief Logarithmic map of the special orthogonal group SO(3)
 *
 * This is the inverse of exp.
 *
 * @param so3   3d rotation (internally represented as unit quaternion)
 * @return corresponding rotation vector (angle times axis).
 */
template <class T, int Options>
Vector3<T> logSO3(const SO3<T, Options>& so3) {
  // SU(2) is a double cover of SO(3), thus we have to multiply the tangent
  // vector delta by two
  return T(2.) * details::logSU2Impl(so3.quaternion());
}

// Jacobian of logarithmic map for SO(3).
//
// Computes the Jacobian for the logarithmic map of the special orthogonal group
// SO(3). This is the derivative with respect to all four coefficients of the
// given quaternion.
//
// The q parameters is the location on manifold where to compute the Jacobian.
// Returns a matrix with all first-order derivatives.
template <typename Scalar, int Options>
Matrix<Scalar, 3, 4> logSO3DerivativeManifold(
    const Quaternion<Scalar, Options>& q) {
  using std::abs;
  using std::atan2;
  using std::sqrt;
  const Scalar w = abs(q.w());
  const Scalar w_sign = (q.w() > 0 ? Scalar(1) : Scalar(-1));
  const Scalar x = w_sign * q.x();
  const Scalar y = w_sign * q.y();
  const Scalar z = w_sign * q.z();
  const Scalar vec_len_2 = q.vec().squaredNorm();
  const Scalar vec_len = sqrt(vec_len_2);
  constexpr Scalar two = Scalar(2);
  if (vec_len < Eigen::NumTraits<Scalar>::dummy_precision()) {
    // This is the Taylor expansion using terms where the combined order of x,
    // y, z is smaller than four.
    constexpr Scalar b = Scalar(1) / Scalar(3);
    const Scalar xx = x * x;
    const Scalar yy = y * y;
    const Scalar zz = z * z;
    constexpr Scalar c = -Scalar(4) / Scalar(3);
    const Scalar cxy = c * x * y;
    const Scalar cxz = c * x * z;
    const Scalar cyz = c * y * z;
    return Matrix<Scalar, 3, 4>{{-two * x, two - xx + b * (yy + zz), cxy, cxz},
                                {-two * y, cxy, two - yy + b * (xx + zz), cyz},
                                {-two * z, cxz, cyz, two - zz + b * (xx + yy)}};
  } else {
    const Scalar b = two * atan2(vec_len, w) / vec_len;
    const Scalar c = (two * w - b) / vec_len_2;
    const Scalar cx = c * x;
    const Scalar cy = c * y;
    const Scalar cz = c * z;
    const Scalar cxx = cx * x;
    const Scalar cyy = cy * y;
    const Scalar czz = cz * z;
    const Scalar cxy = cx * y;
    const Scalar cxz = cx * z;
    const Scalar cyz = cy * z;
    return w_sign * Matrix<Scalar, 3, 4>{
                        {-two * x, b + cxx, cxy, cxz},
                        {-two * y, cxy, b + cyy, cyz},
                        {-two * z, cxz, cyz, b + czz},
                    };
  }
}

// Exponential map of the group of unit quaternions.
//
// It maps a tangent element `delta` of the group at the Identity, which is
// equivalent to the rotation vector axis-times-half-angle), to the
// corresponding unit quaternion.
template <class T, int Options = Eigen::AutoAlign>
Quaternion<T, Options> expQuaternion(const Vector3<T, Options>& delta) {
  return details::expSU2Impl(delta);
}

// Logarithmic map of the group of unit quaternions.
//
// It maps a unit `quaternion` into a tangent element of the group at the
// Identity, which is equivalent to the rotation vector axis-times-half-angle).
template <class T, int Options>
Vector3<T> logQuaternion(const Quaternion<T, Options>& quaternion) {
  return details::logSU2Impl(quaternion);
}

/**
 * \brief Riemannian manifold exponential map on the composite manifold R3xS0(3)
 * at the identity
 *
 * Convention for tangent space is: (delta_x, delta_y, delta_z, delta_r1,
 * delta_r2, delta_r3).
 *
 * @param tangent Direction in tangent space
 * @return 3D transformation
 */
template <typename Scalar, int Options = Eigen::AutoAlign>
Pose3<Scalar, Eigen::AutoAlign> expRiemann(
    const Vector6<Scalar, Options>& tangent) {
  return Pose3<Scalar, Eigen::AutoAlign>{
      intrinsic::eigenmath::expSO3(
          Vector3<Scalar>(tangent.template bottomLeftCorner<3, 1>())),
      tangent.template topLeftCorner<3, 1>()};
}

/**
 * \brief Riemannian manifold log on the composite manifold R3xSO(3) at the
 * identity
 *
 * Convention for tangent space is: (delta_x, delta_y, delta_z, delta_r1,
 * delta_r2, delta_r3).
 *
 * @param pose A 3D transformation
 * @return Direction in tangent space
 */
template <typename Scalar, int Options>
Vector6<Scalar, Eigen::AutoAlign> logRiemann(
    const Pose3<Scalar, Options>& pose) {
  Vector6<Scalar, Eigen::AutoAlign> delta;
  delta.template topLeftCorner<3, 1>() = pose.translation();
  delta.template bottomLeftCorner<3, 1>() =
      intrinsic::eigenmath::logSO3(pose.so3());
  return delta;
}

/**
 * \brief Riemannian manifold log on the composite manifold R3xSO(3) around a
 * frame
 *
 * Convention for tangent space is: (delta_x, delta_y, delta_z, delta_r1,
 * delta_r2, delta_r3).
 *
 * @param ref_pose_t pose of tangent space in reference frame
 * @param ref_pose_a transformation from reference frame to frame A
 * @return tangent pointing from tangent space origin to corresponding reference
 * frame A
 */
template <typename Scalar, int Options_t, int Options_a>
Vector6<Scalar, Eigen::AutoAlign> logRiemann(
    const Pose3<Scalar, Options_t>& ref_pose_t,
    const Pose3<Scalar, Options_a>& ref_pose_a) {
  return intrinsic::eigenmath::logRiemann(ref_pose_t.inverse() * ref_pose_a);
}

}  // namespace eigenmath
}  // namespace intrinsic

#endif  // INTRINSIC_EIGENMATH_MANIFOLDS_H_
