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

#ifndef INTRINSIC_EIGENMATH_EIGEN_MATRIX_HASH_H_
#define INTRINSIC_EIGENMATH_EIGEN_MATRIX_HASH_H_

#include "Eigen/Core"

namespace intrinsic {
namespace eigenmath {

// Simple hash for Eigen-derived types, e.g., using an Eigen::Vector as key in a
// absl hash_set or hash_map. The hash is computed by combining both the matrix
// dimensions and the matrix coefficients.
template <class Derived>
class EigenMatrixHash {
 public:
  using MatrixType =
      Eigen::Matrix<typename Derived::Scalar, Derived::RowsAtCompileTime,
                    Derived::ColsAtCompileTime>;
  using IndexType = typename MatrixType::Index;

  template <class AnyDerived>
  explicit EigenMatrixHash(const Eigen::MatrixBase<AnyDerived>& matrix)
      : matrix_(matrix) {}
  ~EigenMatrixHash() = default;

  friend bool operator==(const EigenMatrixHash& lhs,
                         const EigenMatrixHash& rhs) {
    if (lhs.matrix_.rows() != rhs.matrix_.rows()) return false;
    if (lhs.matrix_.cols() != rhs.matrix_.cols()) return false;
    for (IndexType row = 0; row < lhs.matrix_.rows(); row++) {
      for (IndexType col = 0; col < lhs.matrix_.cols(); col++) {
        if (lhs.matrix_(row, col) != rhs.matrix_(row, col)) {
          return false;
        }
      }
    }
    return true;
  }

  template <typename H>
  friend H AbslHashValue(H state, const EigenMatrixHash& obj) {
    state = H::combine(std::move(state), obj.matrix_.rows());
    state = H::combine(std::move(state), obj.matrix_.cols());
    for (IndexType row = 0; row < obj.matrix_.rows(); row++) {
      for (IndexType col = 0; col < obj.matrix_.cols(); col++) {
        state = H::combine(std::move(state), obj.matrix_(row, col));
      }
    }
    return state;
  }

 private:
  MatrixType matrix_;
};

}  // namespace eigenmath
}  // namespace intrinsic

#endif  // INTRINSIC_EIGENMATH_EIGEN_MATRIX_HASH_H_
