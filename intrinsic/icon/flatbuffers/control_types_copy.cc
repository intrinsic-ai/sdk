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


#include "intrinsic/icon/flatbuffers/control_types_copy.h"

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/flatbuffers/transform_types.fbs.h"
#include "intrinsic/icon/flatbuffers/transform_view.h"
#include "intrinsic/math/pose3.h"

namespace intrinsic_fbs {

// Copies the raw data memory from Pose3d to memory in
// Transform.
void CopyTo(Transform* transform, const intrinsic::Pose3d& pose) {
  transform->mutable_position() = *ViewAs<Point>(pose.translation());
  transform->mutable_rotation() = *ViewAs<Rotation>(pose.quaternion());
}

// Copies the raw data memory from Transform to memory in Pose3d.
void CopyTo(intrinsic::Pose3d* pose, const Transform& transform) {
  pose->translation() = ViewAs<Eigen::Vector3d>(transform.position());
  pose->so3().quaternion() = ViewAs<Eigen::Quaterniond>(transform.rotation());
}

}  // namespace intrinsic_fbs
