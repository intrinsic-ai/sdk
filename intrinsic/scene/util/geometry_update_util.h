// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_UTIL_GEOMETRY_UPDATE_UTIL_H_
#define INTRINSIC_SCENE_UTIL_GEOMETRY_UPDATE_UTIL_H_

#include "absl/status/status.h"
#include "intrinsic/scene/proto/v1/scene_object_updates.pb.h"
#include "intrinsic/world/proto/geometry_component.pb.h"

namespace intrinsic::scene_object {

// Applies a GeometryUpdate to a GeometryComponent.
//
// This function modifies `geometry_component` in-place based on the
// instructions in `update`.
absl::Status ApplyGeometryUpdate(
    const intrinsic_proto::scene_object::v1::GeometryUpdate& update,
    intrinsic_proto::world::GeometryComponent& geometry_component);

}  // namespace intrinsic::scene_object

#endif  // INTRINSIC_SCENE_UTIL_GEOMETRY_UPDATE_UTIL_H_
