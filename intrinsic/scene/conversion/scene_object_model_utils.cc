// Copyright 2023 Intrinsic Innovation LLC

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
