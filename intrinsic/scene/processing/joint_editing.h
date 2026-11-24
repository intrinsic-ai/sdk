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

#ifndef INTRINSIC_SCENE_PROCESSING_JOINT_EDITING_H_
#define INTRINSIC_SCENE_PROCESSING_JOINT_EDITING_H_

#include <optional>
#include <string>
#include <utility>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/world/proto/kinematics_component.pb.h"

namespace intrinsic::scene_object {

// Parameters for creating a new articulated joint in a SceneObject.
// `origin_root_frame` (meters, default `(0, 0, 0)`) and `axis_root_frame`
// (default `(0, 0, 1)`) are specified in the SceneObject's root entity frame.
// `limits` specifies optional `(lower, upper)` bounds in radians for revolute
// joints or meters for prismatic joints.
struct AddJointParams {
  std::string parent_link;
  std::string child_link;
  std::string joint_name;
  intrinsic_proto::world::KinematicsComponent::MotionType motion_type =
      intrinsic_proto::world::KinematicsComponent::MOTION_TYPE_UNDEFINED;
  eigenmath::Vector3d origin_root_frame = eigenmath::Vector3d::Zero();
  eigenmath::Vector3d axis_root_frame = eigenmath::Vector3d::UnitZ();
  std::optional<std::pair<double, double>> limits;
};

// Creates a new joint connecting `params.parent_link` and `params.child_link`
// while preserving the child link's pose relative to the root entity.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject> AddJoint(
    const intrinsic_proto::scene_object::v1::SceneObject& scene_object,
    const AddJointParams& params);

// Parameters for updating an existing joint named `joint_name` in a
// SceneObject. Any optional field left as `std::nullopt` preserves the existing
// joint's property; any field with a value updates that property
// (`origin_root_frame` in meters and `axis_root_frame` are specified in the
// SceneObject's root entity frame, and `limits` specifies `(lower, upper)`
// bounds).
struct UpdateJointParams {
  std::string joint_name;
  std::optional<std::string> parent_link;
  std::optional<std::string> child_link;
  std::optional<intrinsic_proto::world::KinematicsComponent::MotionType>
      motion_type;
  std::optional<eigenmath::Vector3d> origin_root_frame;
  std::optional<eigenmath::Vector3d> axis_root_frame;
  std::optional<std::pair<double, double>> limits;
};

// Updates an existing joint named `params.joint_name`, preserving any
// `std::nullopt` fields (`parent_link`, `child_link`, `motion_type`,
// `origin_root_frame`, `axis_root_frame`, `limits`) from the existing joint
// definition.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject> UpdateJoint(
    const intrinsic_proto::scene_object::v1::SceneObject& scene_object,
    const UpdateJointParams& params);

// Deletes the joint named `joint_name` and reparents its child entities to the
// joint's parent link while preserving their resting poses.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject> DeleteJoint(
    const intrinsic_proto::scene_object::v1::SceneObject& scene_object,
    absl::string_view joint_name);

}  // namespace intrinsic::scene_object

#endif  // INTRINSIC_SCENE_PROCESSING_JOINT_EDITING_H_
