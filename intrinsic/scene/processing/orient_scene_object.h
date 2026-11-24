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

// Copyright 2026 Intrinsic Innovation LLC. All rights reserved.

#ifndef INTRINSIC_SCENE_PROCESSING_ORIENT_SCENE_OBJECT_H_
#define INTRINSIC_SCENE_PROCESSING_ORIENT_SCENE_OBJECT_H_

#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"

namespace intrinsic::scene_object {

// A cardinal axis direction in the SceneObject's root entity frame.
enum class AxisDirection {
  kPosX,
  kNegX,
  kPosY,
  kNegY,
  kPosZ,
  kNegZ,
};

// The result of `OrientAndRecenterSceneObject`. `applied_transform` is the
// `new_root_t_old_root` pose applied to all children of the root entity and to
// the root entity's geometry.
struct OrientAndRecenterResult {
  intrinsic_proto::scene_object::v1::SceneObject scene_object;
  Pose3d applied_transform;
};

// Orients `scene_object` so that `up_axis` points along +Z, then recenters it
// so that `origin_offset` becomes the new root origin.
//
// The two steps run in this order:
//   1. Rotation first: rotates `scene_object` about the root origin so that
//      `up_axis` maps to +Z (`p_rotated = R * p`).
//   2. Recentering second, in the post-rotation +Z-up frame: shifts the
//      rotated root frame by `-origin_offset` (`p_new = (R * p) -
//      origin_offset`). The post-rotation point `origin_offset` moves to the
//      new root origin (0, 0, 0).
//
// `origin_offset` is in meters and is expressed in the post-rotation +Z-up
// frame, not in the original frame.
//
// Example: a +Y-up SceneObject (`AxisDirection::kPosY`) has a point at
// (0, 2, 0), and `origin_offset = (0, 0, 0.5)`.
//   (1) Rotating +Y to +Z (+90 deg about +X) maps (0, 2, 0) to (0, 0, 2) in the
//       +Z-up frame.
//   (2) Translating by `-origin_offset = (0, 0, -0.5)` maps (0, 0, 2) to
//       (0, 0, 1.5), and places (0, 0, 0.5) at the new root origin (0, 0, 0).
//
// Returns INVALID_ARGUMENT if `up_axis` is not a valid `AxisDirection` or if
// `origin_offset` has a non-finite component. If the applied transform is not
// the identity, returns the errors of `TransformSceneObject` (for example,
// INVALID_ARGUMENT if `scene_object` has no unique root entity).
absl::StatusOr<OrientAndRecenterResult> OrientAndRecenterSceneObject(
    const intrinsic_proto::scene_object::v1::SceneObject& scene_object,
    AxisDirection up_axis, const eigenmath::Vector3d& origin_offset);

}  // namespace intrinsic::scene_object

#endif  // INTRINSIC_SCENE_PROCESSING_ORIENT_SCENE_OBJECT_H_
