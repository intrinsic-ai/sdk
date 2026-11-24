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

#include "intrinsic/scene/processing/joint_editing.h"

#include <cmath>
#include <optional>
#include <string>
#include <utility>

#include "absl/container/flat_hash_map.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/math/proto_conversion.h"
#include "intrinsic/scene/proto/v1/entity.pb.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/scene/proto/v1/scene_object_updates.pb.h"
#include "intrinsic/scene/scene_object/scene_object_view.h"
#include "intrinsic/scene/util/scene_object_updates.h"
#include "intrinsic/scene/validate/scene_object_validation.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/proto/kinematics_component.pb.h"

namespace intrinsic::scene_object {
namespace {

using ::intrinsic::Pose3d;
using ::intrinsic::ToProto;
using ::intrinsic::ToVectorProto;
using ::intrinsic::eigenmath::Quaterniond;
using ::intrinsic::eigenmath::Vector3d;
using ::intrinsic_proto::FromProto;
using ::intrinsic_proto::FromProtoNormalized;
using ::intrinsic_proto::scene_object::v1::CreateJointUpdate;
using ::intrinsic_proto::scene_object::v1::DeleteEntityUpdate;
using ::intrinsic_proto::scene_object::v1::Entity;
using ::intrinsic_proto::scene_object::v1::ReparentEntityUpdate;
using ::intrinsic_proto::scene_object::v1::SceneObject;
using ::intrinsic_proto::scene_object::v1::SceneObjectUpdates;
using ::intrinsic_proto::world::KinematicsComponent;

using RootPoseMap = ::absl::flat_hash_map<absl::string_view, Pose3d>;

constexpr double kMinAxisSquaredNorm = 1e-12;

// Applies SceneObjectUpdates via ProcessSceneObjectUpdates, verifies post-edit
// SceneObject invariants, and returns the updated SceneObject.
absl::StatusOr<SceneObject> ApplyUpdates(const SceneObject& scene_object,
                                         const SceneObjectUpdates& updates) {
  INTR_ASSIGN_OR_RETURN(SceneObjectUpdateResult update_result,
                        ProcessSceneObjectUpdates(scene_object, updates));
  INTR_RETURN_IF_ERROR(ValidateSceneObject(update_result.result));
  return std::move(update_result).result;
}

// Validates that `joint_name` identifies an existing joint entity in `view` and
// returns a non-null pointer to that entity.
absl::StatusOr<const Entity*> FindJointEntity(const SceneObjectView& view,
                                              absl::string_view joint_name) {
  if (joint_name.empty()) {
    return absl::InvalidArgumentError("Missing joint_name");
  }
  const Entity* const entity = view.FindEntity(joint_name);
  if (entity == nullptr) {
    return absl::NotFoundError(
        absl::StrCat("Joint '", joint_name, "' not found"));
  }
  if (!entity->has_joint()) {
    return absl::InvalidArgumentError(
        absl::StrCat("Entity '", joint_name, "' is not a joint"));
  }
  return entity;
}

// Validates that joint origin, axis, and optional limits in `params` are finite
// and compatible with `params.motion_type`.
absl::Status ValidateAddJointParams(const AddJointParams& params) {
  if (!params.origin_root_frame.allFinite()) {
    return absl::InvalidArgumentError("Joint origin must be finite");
  }
  if (!params.axis_root_frame.allFinite() ||
      params.axis_root_frame.squaredNorm() < kMinAxisSquaredNorm) {
    return absl::InvalidArgumentError(
        absl::StrCat("Joint axis must be finite and non-zero for joint '",
                     params.joint_name, "'"));
  }
  if (params.limits.has_value()) {
    if (params.motion_type == KinematicsComponent::MOTION_TYPE_FIXED) {
      return absl::InvalidArgumentError(
          "Limits cannot be specified for a fixed joint");
    }
    if (!std::isfinite(params.limits->first) ||
        !std::isfinite(params.limits->second)) {
      return absl::InvalidArgumentError("Joint position limits must be finite");
    }
  }
  return absl::OkStatus();
}

// Returns true when `child_link` is currently parented to a link other than
// `parent_link` and both `parent_link` and `child_link` are distinct links.
bool ShouldReparentChildLink(const SceneObjectView& view,
                             absl::string_view parent_link,
                             absl::string_view child_link) {
  if (parent_link == child_link) {
    return false;
  }
  const Entity* const parent_entity = view.FindEntity(parent_link);
  const Entity* const child_entity = view.FindEntity(child_link);
  const Entity* const current_child_parent = view.GetParent(child_link);
  return parent_entity != nullptr && parent_entity->has_link() &&
         child_entity != nullptr && child_entity->has_link() &&
         current_child_parent != nullptr && current_child_parent->has_link() &&
         current_child_parent->name() != parent_link;
}

// Builds SceneObjectUpdates (optional child-link reparent + CreateJointUpdate)
// from root-frame origin (meters) and axis coordinates.
absl::StatusOr<SceneObjectUpdates> BuildCreateJointUpdate(
    const SceneObject& scene_object, const AddJointParams& params) {
  INTR_RETURN_IF_ERROR(ValidateAddJointParams(params));
  INTR_ASSIGN_OR_RETURN(const SceneObjectView view,
                        SceneObjectView::Create(scene_object));

  SceneObjectUpdates updates;
  if (ShouldReparentChildLink(view, params.parent_link, params.child_link)) {
    ReparentEntityUpdate* const reparent =
        updates.add_updates()->mutable_reparent_entity();
    reparent->set_entity_name(params.child_link);
    reparent->set_new_parent_name(params.parent_link);
  }

  CreateJointUpdate* const create =
      updates.add_updates()->mutable_create_joint();
  create->set_parent_link_name(params.parent_link);
  create->set_child_link_name(params.child_link);
  create->set_new_joint_name(params.joint_name);
  create->set_parent_t_child_policy(
      CreateJointUpdate::PARENT_T_CHILD_POLICY_PRESERVE);

  KinematicsComponent* const kc =
      create->mutable_joint()->mutable_kinematics_component();
  kc->set_motion_type(params.motion_type);
  *kc->mutable_axis() = ToVectorProto(params.axis_root_frame);

  INTR_ASSIGN_OR_RETURN(const RootPoseMap root_t_entity, view.GetRootPoses());
  if (const auto parent_it = root_t_entity.find(params.parent_link);
      parent_it != root_t_entity.end()) {
    const Pose3d root_t_joint(Quaterniond::Identity(),
                              params.origin_root_frame);
    *kc->mutable_parent_t_inboard() =
        ToProto(parent_it->second.inverse() * root_t_joint);
  }

  if (params.limits.has_value()) {
    KinematicsComponent::Limits::LimitPair* const limits =
        kc->mutable_system_limits()->mutable_fixed_limits();
    limits->set_lower(params.limits->first);
    limits->set_upper(params.limits->second);
  }
  return updates;
}

// Builds a DeleteEntityUpdate that reparents child entities of `joint_name`.
SceneObjectUpdates BuildDeleteJointUpdate(absl::string_view joint_name) {
  SceneObjectUpdates updates;
  DeleteEntityUpdate* const del =
      updates.add_updates()->mutable_delete_entity();
  del->set_entity_name(joint_name);
  del->set_child_policy(DeleteEntityUpdate::CHILD_POLICY_REPARENT_CHILDREN);
  return updates;
}

// Computes the existing joint's inboard frame pose in root coordinates.
absl::StatusOr<Pose3d> ResolveExistingRootTInboard(
    const SceneObjectView& view, const Entity& existing_joint) {
  const KinematicsComponent& existing_kc =
      existing_joint.joint().kinematics_component();
  INTR_ASSIGN_OR_RETURN(const RootPoseMap root_poses, view.GetRootPoses());
  const auto parent_it = root_poses.find(existing_joint.parent_name());
  if (parent_it == root_poses.end()) {
    return absl::NotFoundError(
        absl::StrCat("Parent link '", existing_joint.parent_name(),
                     "' of joint '", existing_joint.name(), "' not found"));
  }
  if (!existing_kc.has_parent_t_inboard()) {
    return parent_it->second;
  }
  INTR_ASSIGN_OR_RETURN(const Pose3d existing_parent_t_inboard,
                        FromProtoNormalized(existing_kc.parent_t_inboard()));
  return parent_it->second * existing_parent_t_inboard;
}

// Finds the name of the first child link entity under `joint_name`, or "".
std::string FindExistingChildLinkName(const SceneObjectView& view,
                                      absl::string_view joint_name) {
  for (const Entity* const child : view.GetChildren(joint_name)) {
    if (child->has_link()) {
      return child->name();
    }
  }
  return "";
}

// Resolves joint limits for UpdateJoint, preserving existing fixed limits on
// non-fixed joints when `limits_override` is unset.
std::optional<std::pair<double, double>> ResolveJointLimits(
    const KinematicsComponent& existing_kc,
    KinematicsComponent::MotionType resolved_motion_type,
    const std::optional<std::pair<double, double>>& limits_override) {
  const KinematicsComponent::Limits& system_limits =
      existing_kc.system_limits();
  if (limits_override.has_value() ||
      resolved_motion_type == KinematicsComponent::MOTION_TYPE_FIXED ||
      !system_limits.has_fixed_limits()) {
    return limits_override;
  }
  const KinematicsComponent::Limits::LimitPair& fixed =
      system_limits.fixed_limits();
  return std::make_pair(fixed.lower(), fixed.upper());
}

}  // namespace

absl::StatusOr<SceneObject> AddJoint(const SceneObject& scene_object,
                                     const AddJointParams& params) {
  INTR_ASSIGN_OR_RETURN(const SceneObjectUpdates updates,
                        BuildCreateJointUpdate(scene_object, params));
  return ApplyUpdates(scene_object, updates);
}

absl::StatusOr<SceneObject> UpdateJoint(const SceneObject& scene_object,
                                        const UpdateJointParams& params) {
  INTR_ASSIGN_OR_RETURN(const SceneObjectView view,
                        SceneObjectView::Create(scene_object));
  INTR_ASSIGN_OR_RETURN(const Entity* const existing_joint,
                        FindJointEntity(view, params.joint_name));
  INTR_ASSIGN_OR_RETURN(const Pose3d root_t_inboard,
                        ResolveExistingRootTInboard(view, *existing_joint));

  const KinematicsComponent& existing_kc =
      existing_joint->joint().kinematics_component();
  const KinematicsComponent::MotionType resolved_motion_type =
      params.motion_type.value_or(existing_kc.motion_type());
  const Vector3d existing_axis_root =
      root_t_inboard.quaternion() * (existing_kc.has_axis()
                                         ? FromProto(existing_kc.axis())
                                         : Vector3d::UnitZ());
  const AddJointParams resolved_params{
      .parent_link = params.parent_link.value_or(existing_joint->parent_name()),
      .child_link = params.child_link.value_or(
          FindExistingChildLinkName(view, params.joint_name)),
      .joint_name = params.joint_name,
      .motion_type = resolved_motion_type,
      .origin_root_frame =
          params.origin_root_frame.value_or(root_t_inboard.translation()),
      .axis_root_frame = params.axis_root_frame.value_or(existing_axis_root),
      .limits =
          ResolveJointLimits(existing_kc, resolved_motion_type, params.limits),
  };

  INTR_ASSIGN_OR_RETURN(
      const SceneObject without_joint,
      ApplyUpdates(scene_object, BuildDeleteJointUpdate(params.joint_name)));
  return AddJoint(without_joint, resolved_params);
}

absl::StatusOr<SceneObject> DeleteJoint(const SceneObject& scene_object,
                                        absl::string_view joint_name) {
  INTR_ASSIGN_OR_RETURN(const SceneObjectView view,
                        SceneObjectView::Create(scene_object));
  INTR_RETURN_IF_ERROR(FindJointEntity(view, joint_name).status());
  return ApplyUpdates(scene_object, BuildDeleteJointUpdate(joint_name));
}

}  // namespace intrinsic::scene_object
