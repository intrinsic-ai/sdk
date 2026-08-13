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

#ifndef INTRINSIC_WORLD_CONVERSION_SDF_JOINT_CONVERSION_H_
#define INTRINSIC_WORLD_CONVERSION_SDF_JOINT_CONVERSION_H_

#include <memory>
#include <optional>
#include <string>
#include <vector>

#include "absl/status/statusor.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/world/component/kinematics_component.h"
#include "intrinsic/world/conversion/sdf/sensor_conversion.h"
#include "sdf/Joint.hh"
#include "sdf/Model.hh"

namespace intrinsic {
namespace sdf {

// Parses necessary from `sdf_joint` into data that can instantiate a Intrinsic
// equivalent of a joint.
struct ParseJointResult {
  // Name of the joint.
  std::string name;

  // Intrinsic and SDF store joint's pose and parent-child relationships in a
  // different way.
  //
  // SDF decomposes the transforms into:
  // - model_t_parent (parent's <pose>)
  // - model_t_child (child's <pose>)
  // - child_t_joint (joint's <pose>)
  // See: https://sdformat.org/tutorials?tut=pose_frame_semantics for detailed
  // explanation.
  //
  // While we decompose the transforms into:
  // - parent_t_joint (a.k.a. the joint's inboard transform, stored in its
  //   KinematicsComponent.parent_t_inboard)

  // Name of the joint's parent.
  std::string parent_name;

  // Name of the joint's child
  std::string child_name;

  // The pose from child frame to joint frame.
  Pose3d child_t_joint;

  // Pose from joint's parent frame to joint's child frame.
  // If joint's parent is "world", parent_t_child will be std::nullopt because
  // there is no obvious way to find the world pose of a nested model.
  std::optional<Pose3d> parent_t_child;

  // Kinematics component contains kinematic motion data. Also contain joint's
  // inbound transform(joint's parent frame to joint) if joint's parent is not
  // "world".
  std::unique_ptr<KinematicsComponent> kinematics_component;

  // Sensors that are direct children elements of the sdf joint.
  std::vector<ParseSensorResult> sensors;
};
absl::StatusOr<ParseJointResult> ParseJoint(const ::sdf::Joint& sdf_joint);
}  // namespace sdf
}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_CONVERSION_SDF_JOINT_CONVERSION_H_
