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

#include "intrinsic/scene/processing/orient_scene_object.h"

#include <numbers>
#include <type_traits>
#include <utility>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/scene/processing/transform_scene_object.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::scene_object {
namespace {

using ::intrinsic::Pose3d;
using ::intrinsic::eigenmath::kDoNotNormalize;
using ::intrinsic::eigenmath::Quaterniond;
using ::intrinsic::eigenmath::Vector3d;
using ::intrinsic_proto::scene_object::v1::SceneObject;

// Function-local static `Pose3d` constants below need no destructor call at
// exit.
static_assert(std::is_trivially_destructible_v<Pose3d>);

constexpr double kInvSqrt2 = std::numbers::sqrt2 / 2.0;

// Returns the rotation that maps the `up_axis` unit vector to +Z. The returned
// pointer points to a static constant built once from an exact unit
// quaternion, so callers do not construct or orthogonalize a rotation per call.
// Returns INVALID_ARGUMENT if `up_axis` is not a valid `AxisDirection`.
absl::StatusOr<const Pose3d*> UpAxisToPlusZRotation(AxisDirection up_axis) {
  switch (up_axis) {
    case AxisDirection::kPosX: {
      // -90 deg about +Y.
      static const Pose3d kRotation(
          Quaterniond(kInvSqrt2, 0.0, -kInvSqrt2, 0.0), Vector3d::Zero(),
          kDoNotNormalize);
      return &kRotation;
    }
    case AxisDirection::kNegX: {
      // +90 deg about +Y.
      static const Pose3d kRotation(Quaterniond(kInvSqrt2, 0.0, kInvSqrt2, 0.0),
                                    Vector3d::Zero(), kDoNotNormalize);
      return &kRotation;
    }
    case AxisDirection::kPosY: {
      // +90 deg about +X.
      static const Pose3d kRotation(Quaterniond(kInvSqrt2, kInvSqrt2, 0.0, 0.0),
                                    Vector3d::Zero(), kDoNotNormalize);
      return &kRotation;
    }
    case AxisDirection::kNegY: {
      // -90 deg about +X.
      static const Pose3d kRotation(
          Quaterniond(kInvSqrt2, -kInvSqrt2, 0.0, 0.0), Vector3d::Zero(),
          kDoNotNormalize);
      return &kRotation;
    }
    case AxisDirection::kPosZ: {
      static const Pose3d kRotation = Pose3d::Identity();
      return &kRotation;
    }
    case AxisDirection::kNegZ: {
      // 180 deg about +X.
      static const Pose3d kRotation(Quaterniond(0.0, 1.0, 0.0, 0.0),
                                    Vector3d::Zero(), kDoNotNormalize);
      return &kRotation;
    }
  }
  return absl::InvalidArgumentError(
      absl::StrCat("Invalid up_axis value: ", static_cast<int>(up_axis), "."));
}

}  // namespace

absl::StatusOr<OrientAndRecenterResult> OrientAndRecenterSceneObject(
    const SceneObject& scene_object, AxisDirection up_axis,
    const Vector3d& origin_offset) {
  INTR_ASSIGN_OR_RETURN(const Pose3d* const rot,
                        UpAxisToPlusZRotation(up_axis));
  if (!origin_offset.allFinite()) {
    return absl::InvalidArgumentError(
        absl::StrCat("origin_offset must be finite, got (", origin_offset.x(),
                     ", ", origin_offset.y(), ", ", origin_offset.z(), ")."));
  }

  const Pose3d applied_transform(rot->so3(), -origin_offset);
  INTR_ASSIGN_OR_RETURN(SceneObject transformed_scene_object,
                        TransformSceneObject(scene_object, applied_transform));
  return OrientAndRecenterResult{
      .scene_object = std::move(transformed_scene_object),
      .applied_transform = applied_transform,
  };
}

}  // namespace intrinsic::scene_object
