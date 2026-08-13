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

#ifndef INTRINSIC_SCENE_VALIDATE_SCENE_OBJECT_VALIDATION_H_
#define INTRINSIC_SCENE_VALIDATE_SCENE_OBJECT_VALIDATION_H_

#include "absl/container/flat_hash_set.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/scene/proto/v1/collision_rules.pb.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/world/proto/kinematics_component.pb.h"

namespace intrinsic {
namespace scene_object {

// Returns OK if this is a well formed scene object.
absl::Status ValidateSceneObject(
    const intrinsic_proto::scene_object::v1::SceneObject& object);

// Returns OK if this component has valid joint values and limits.
absl::Status ValidateJointLimits(
    const intrinsic_proto::world::KinematicsComponent& component);

// Returns OK if the set of collision rules are well formed.
absl::Status ValidateCollisionRules(
    const intrinsic_proto::scene_object::v1::CollisionRules& collision_rules,
    const absl::flat_hash_set<absl::string_view>& entity_names);

// Returns OK if the given of collision entity pair is well formed.
absl::Status ValidateCollisionEntityPair(
    const intrinsic_proto::scene_object::v1::CollisionEntityPair& entity_pair,
    const absl::flat_hash_set<absl::string_view>& entity_names);

}  // namespace scene_object
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_VALIDATE_SCENE_OBJECT_VALIDATION_H_
