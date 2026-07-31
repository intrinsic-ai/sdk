// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_REMOVE_DUPLICATE_VERTICES_FROM_MESH_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_REMOVE_DUPLICATE_VERTICES_FROM_MESH_H_

#include <optional>

#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"

namespace intrinsic {
namespace geometry_legacy {

// Removes the duplicate vertices from the given 'mesh' and returns a new mesh.
// For considering two vertices as duplicates the function allows two variants:
// 1. Approximate (using a 'maybe_epsilon' > 0): vertices within a radius of
//    epsilon are considered duplicates, and replaced by the first vertex in the
//    cluster (default to have the previously used behavior).
// 2. Exact ('maybe_epsilon' = nullopt): only vertices being exactly equal are
//    considered duplicates.
Mesh RemoveDuplicateVerticesFromMesh(
    const Mesh& mesh, std::optional<double> maybe_epsilon = 1e-6);

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_REMOVE_DUPLICATE_VERTICES_FROM_MESH_H_
