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

#ifndef INTRINSIC_EIGENMATH_SKEW_SYMMETRIC_MATRIX_H_
#define INTRINSIC_EIGENMATH_SKEW_SYMMETRIC_MATRIX_H_

#include "Eigen/Core"

namespace intrinsic {
namespace eigenmath {

// Returns the skew-symmetric (cross-product) matrix of the vector x.
template <typename Derived>
inline Eigen::Matrix3<typename Derived::Scalar> SkewSymmetricMatrix(
    const Eigen::MatrixBase<Derived>& x) {
  EIGEN_STATIC_ASSERT_VECTOR_SPECIFIC_SIZE(Derived, 3);
  using Scalar = typename Derived::Scalar;

  const Scalar zero = static_cast<Scalar>(0);
  Eigen::Matrix3<Scalar> s;
  // clang-format off
  s << zero, -x(2), x(1),
      x(2), zero, -x(0),
      -x(1), x(0), zero;
  // clang-format on
  return s;
}

}  // namespace eigenmath
}  // namespace intrinsic

#endif  // INTRINSIC_EIGENMATH_SKEW_SYMMETRIC_MATRIX_H_
