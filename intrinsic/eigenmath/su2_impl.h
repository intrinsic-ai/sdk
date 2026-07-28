// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_EIGENMATH_SU2_IMPL_H_
#define INTRINSIC_EIGENMATH_SU2_IMPL_H_

#include <cmath>

#include "Eigen/Core"
#include "absl/log/absl_check.h"
#include "intrinsic/eigenmath/types.h"

namespace intrinsic {
namespace eigenmath {
namespace details {

/**
 * @brief Exponential map of the special unitary group SU(2). This function is
 * an implementation detail for expSO3.
 *
 * Let \f$ \exp \f$ be the matrix exponential and \f$ \hat{\cdot} \f$ be the
 * function which maps a tangent vector in SU(2) to its corresponding (2x2)
 * matrix representation. It holds hat:
 *
 * \f$ \exp_{SU(2)}(\delta) = exp(\hat{delta}) \f$
 *
 * @param delta tangent vector of SU(2)
 * @return unit quaternion, a member of SU(2)
 */
template <class T, int Options = Eigen::AutoAlign>
Quaternion<T, Options> expSU2Impl(const eigenmath::Vector3<T, Options>& delta) {
  // Employ "using" statement here, since we need to find "sin" and friends
  // using ADL if they are not defined in the std namespace for type T.
  using std::cos;
  using std::sin;
  using std::sqrt;

  Quaternion<T, Options> q_delta;
  T theta_squared = delta.squaredNorm();
  if (theta_squared > Eigen::NumTraits<T>::dummy_precision()) {
    T theta = sqrt(theta_squared);
    q_delta.w() = cos(theta);
    q_delta.vec() = (sin(theta) / theta) * delta;
  } else {
    // taylor expansions around theta_squared==0
    q_delta.w() = T(1.) - T(0.5) * theta_squared;
    q_delta.vec() = (T(1.) - T(1. / 6.) * theta_squared) * delta;
  }
  return q_delta;
}

/**
 * @brief Logarithmic map of the special unitary group SU(2). This function is
 * an implementation detail for logSO3.
 *
 * This is the inverse of expSU2Impl.
 *
 * @param q quaternion as element of SU(2)
 * @pre quaternion q must of unit length
 * @return corresponding vector in the tangent space of SU(2)
 */
template <class T, int Options = Eigen::AutoAlign>
eigenmath::Vector3<T> logSU2Impl(const Quaternion<T, Options>& q) {
  // Employ using statement here, since we need to find "abs" and friends using
  // ADL if they are not defined in the std namespace for type T.
  using std::abs;
  using std::atan2;
  using std::pow;
  using std::sqrt;
  ABSL_CHECK(abs(T(1.) - q.squaredNorm()) <
             Eigen::NumTraits<T>::dummy_precision())
      << "quaternion q must be approx. of unit length";

  // Implementation of the logarithmic map of SU(2) using atan.
  // This follows Hertzberg et al. "Integrating Generic Sensor Fusion Algorithms
  // with Sound State Representations through Encapsulation of Manifolds", Eq.
  // (31)
  // We use atan2 instead of atan to enable the use of Eigen Autodiff with SU2:
  // atan2(y,x) is equivalent to atan(y/x) for x > 0. In our case x = w, the
  // real part of the quaternion. With q = -q we chose the quaternion with
  // positive real part.
  const T sign_of_w = q.w() < T{0.0} ? T{-1.0} : T{1.0};
  const T abs_w = sign_of_w * q.w();
  const eigenmath::Vector3<T> v = sign_of_w * q.vec();
  const T squared_norm_of_v = v.squaredNorm();
  eigenmath::Vector3<T> delta;

  if (squared_norm_of_v > Eigen::NumTraits<T>::dummy_precision()) {
    const T norm_of_v = sqrt(squared_norm_of_v);
    delta = (atan2(norm_of_v, abs_w) / norm_of_v) * v;
  } else {
    // taylor expansion at squared_norm_of_v == 0
    delta = (T(1.) / abs_w - squared_norm_of_v / (T(3.) * pow(abs_w, 3))) * v;
  }
  return delta;
}

}  // namespace details
}  // namespace eigenmath
}  // namespace intrinsic

#endif  // INTRINSIC_EIGENMATH_SU2_IMPL_H_
