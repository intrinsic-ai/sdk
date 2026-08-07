// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_VALIDATE_LARGE_MESH_H_
#define INTRINSIC_SCENE_VALIDATE_LARGE_MESH_H_

#include <string>

#include "absl/status/status.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/world/hashing/hashing.h"

namespace intrinsic::scene_object {

// Returns an error if the length of the bounding box diagonal for any of the
// mesh geometries exceeds `max_aabb_diag`.
absl::Status CheckForLargeMeshes(
    const WorldHashMap<std::string, Geometry>& geometry, double max_aabb_diag);

}  // namespace intrinsic::scene_object

#endif  // INTRINSIC_SCENE_VALIDATE_LARGE_MESH_H_
