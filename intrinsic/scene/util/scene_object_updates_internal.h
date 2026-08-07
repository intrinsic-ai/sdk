// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_UTIL_SCENE_OBJECT_UPDATES_INTERNAL_H_
#define INTRINSIC_SCENE_UTIL_SCENE_OBJECT_UPDATES_INTERNAL_H_

#include "absl/status/statusor.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/scene/proto/v1/scene_object_updates.pb.h"
#include "intrinsic/scene/util/scene_object_updates.h"

namespace intrinsic {
namespace scene_object {
namespace internal {

absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdates(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::SceneObjectUpdate& update,
    UpdateType update_type);

absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdates(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::SceneObjectInstanceUpdate& update);

// Will update the entity pose. If we set the update type to instance and we
// attempt to update a non frame entity, we will return an error as that is not
// allowed.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::EntityPoseUpdate& update,
    UpdateType update_type);

// Will add a new frame to the scene object.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::CreateFrameUpdate& update,
    UpdateType update_type);

// Will remove an entity from the scene object. Instance editing update_type is
// not allowed and will return an error.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::DeleteEntityUpdate& update,
    UpdateType update_type);

// Will update the named configurations on the scene object.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::SetNamedConfigurationsUpdate&
        update,
    UpdateType update_type);

// Will update the joints on the scene object.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::UpdateJointsRequest& update,
    UpdateType update_type);

// Will update the cartesian limits on the scene object.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::CartesianLimitsUpdate& update,
    UpdateType update_type);

// Will update the entity name of an entity within the scene object.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::RenameEntityUpdate& update,
    UpdateType update_type);

// Will reparent an entity.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::ReparentEntityUpdate& update,
    UpdateType update_type);

// Will update frame properties.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::UpdateFrameProperties& update,
    UpdateType update_type);

// Will update the collision rules within the scene object.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::UpdateCollisionRules& update);

// Will update the user data within the scene object.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::UpdateUserData& update);

// Will update the IK solvers within the scene object.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::SetIKSolversUpdate& update);

// Will update the geometry within the scene object.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::GeometryUpdate& update);

// Will update the physics properties within the scene object.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::UpdatePhysicsProperties& update);

// Will update the simulation properties within the scene object.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
ProcessSceneObjectUpdate(
    intrinsic_proto::scene_object::v1::SceneObject&& object,
    const intrinsic_proto::scene_object::v1::UpdateSimulationProperties&
        update);

}  // namespace internal
}  // namespace scene_object
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_UTIL_SCENE_OBJECT_UPDATES_INTERNAL_H_
