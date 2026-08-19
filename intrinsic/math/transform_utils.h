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

#ifndef INTRINSIC_MATH_TRANSFORM_UTILS_H_
#define INTRINSIC_MATH_TRANSFORM_UTILS_H_

#include "intrinsic/math/pose3.h"
#include "intrinsic/math/twist.h"

namespace intrinsic {

/**
 * Transforms a wrench from frame B to frame A.
 *
 * @param a_T_b the position and orientation of B relative to A
 * @param b_W   wrench (fx,fy,fz,tx,ty,tz) at the origin of frame B and
 *              expressed in B coordinates.
 *
 * @return a_W the same wrench sitting at the origin of A and is expressed in A
 * coordinates.
 */
Wrench TransformWrench(const Pose3d& a_T_b, const Wrench& b_W);

}  // namespace intrinsic

#endif  // INTRINSIC_MATH_TRANSFORM_UTILS_H_
