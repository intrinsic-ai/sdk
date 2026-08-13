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

#ifndef INTRINSIC_ICON_FLATBUFFERS_MATRIX_TYPES_UTILS_H_
#define INTRINSIC_ICON_FLATBUFFERS_MATRIX_TYPES_UTILS_H_

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/flatbuffers/matrix_types.fbs.h"
namespace intrinsic_fbs {

// Copies the raw data memory from eigenmath::Matrix3d to memory in
// Matrix3d.
void CopyTo(const intrinsic::eigenmath::Matrix3d& matrix, Matrix3d& matrix_fbs);

// Returns a reference to the memory pointed at by a `intrinsic_fbs::Matrix3d`.
// The returned result is Eigen::Map<const eigenmath::Matrix3d>
template <typename T>
Eigen::Map<const T> ViewAs(const Matrix3d& matrix);
template <>
Eigen::Map<const intrinsic::eigenmath::Matrix3d> ViewAs(const Matrix3d& matrix);

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_FLATBUFFERS_MATRIX_TYPES_UTILS_H_
