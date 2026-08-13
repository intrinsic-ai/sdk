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

#include "intrinsic/scene/conversion/scene_object_model_utils.h"

#include <string>
#include <vector>

#include "absl/container/flat_hash_set.h"
#include "intrinsic/scene/proto/v1/entity.pb.h"
#include "intrinsic/world/proto/kinematics_component.pb.h"

namespace intrinsic {
namespace scene_object {

namespace {
using Entity = ::intrinsic_proto::scene_object::v1::Entity;
}

absl::flat_hash_set<std::string> GetNonFixedJointNames(
    const std::vector<Entity>& entities) {
  absl::flat_hash_set<std::string> joint_names;
  for (const auto& entity : entities) {
    if (entity.has_joint() &&
        entity.joint().kinematics_component().motion_type() !=
            intrinsic_proto::world::KinematicsComponent::MOTION_TYPE_FIXED)
      joint_names.insert(entity.name());
  }
  return joint_names;
}

}  // namespace scene_object
}  // namespace intrinsic
