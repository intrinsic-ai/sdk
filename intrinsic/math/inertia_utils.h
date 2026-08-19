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

#ifndef INTRINSIC_MATH_INERTIA_UTILS_H_
#define INTRINSIC_MATH_INERTIA_UTILS_H_

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/world/proto/robot_payload.pb.h"

namespace intrinsic {

constexpr double kMatrixDifferenceThreshold = 1e-6;

bool IsSymmetric(const eigenmath::Matrix3d& matrix,
                 double max_difference_threshold = kMatrixDifferenceThreshold);

// Creates an inertia tensor from the given inertia moments and the products of
// inertia moments.
//
// `i_xx`, `i_yy`, `i_zz`: The inertia moments about the x, y, and z axis
// respectively. The inertia moments must be positive.
//
// `i_xy`, `i_xz`, `i_yz`: The products of the inertia moments
// about the x, y, and z axis respectively.
icon::RealtimeStatusOr<eigenmath::Matrix3d> CreateInertiaTensor(
    double i_xx, double i_yy, double i_zz, double i_xy, double i_xz,
    double i_yz);

// Validates that the link inertia expressed at the center of gravity is
// positive definite (symmetric and with positive eigenvalues) and that its
// eigenvalues fulfill the triangle inequalities.
icon::RealtimeStatus ValidateInertia(const eigenmath::Matrix3d& inertia);

// Rotates the inertia tensor by the given rotation with I' = R * I *
// R^-1.
eigenmath::Matrix3d RotateInertiaTensor(const eigenmath::Matrix3d& inertia,
                                        const eigenmath::Quaterniond& rotation);

struct PrincipalInertiaMoments {
  // The principal inertia moments. These are the diagonal elements of the
  // principal inertia tensor. Since these are the principal inertia moments,
  // the off-diagonal elements are zero and therefore the moments are
  // represented as a vector.
  eigenmath::Vector3d moments;
  // The principal inertia axes represented as a quaternion. Rotating around
  // this quaternion transforms the principal inertia moments back to the
  // original inertia. To construct the principal inertia axes, convert the
  // quaternion to a rotation matrix.
  eigenmath::Quaterniond rotation;
};

// Transforms the inertia tensor to the principal inertia moments and the
// principal inertia axes and returns them. The given `inertia` must be
// expressed at the center of gravity.
// The returned principal inertia axes represent the rotation matrix that
// transforms the returned principal inertia moments back to the original
// `inertia`. The `inertia` tensor must be positive definite (symmetric and with
// positive eigenvalues) and its eigenvalues must fulfill the triangle
// inequalities. Since the principal inertia moments result in a inertia tensor
// as a diagonal matrix, they are represented as a vector. The principal inertia
// axes are the eigenvectors of the inertia tensor arranged as a proper rotation
// matrix and converted to a quaternion.
icon::RealtimeStatusOr<PrincipalInertiaMoments>
TransformToPrincipalInertiaMoments(const eigenmath::Matrix3d& inertia);

}  // namespace intrinsic

#endif  // INTRINSIC_MATH_INERTIA_UTILS_H_
