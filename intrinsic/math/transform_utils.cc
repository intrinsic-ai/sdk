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

#include "intrinsic/math/transform_utils.h"

#include "intrinsic/math/pose3.h"
#include "intrinsic/math/twist.h"

namespace intrinsic {

Wrench TransformWrench(const Pose3d& a_T_b, const Wrench& b_W) {
  Wrench a_W;
  a_W.head<3>() = a_T_b.so3() * b_W.head<3>();
  a_W.tail<3>() =
      a_T_b.so3() * b_W.tail<3>() + a_T_b.translation().cross(a_W.head<3>());
  return Wrench(a_W);
}

}  // namespace intrinsic
