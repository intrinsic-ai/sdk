// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_CONCATENATE_MESHES_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_CONCATENATE_MESHES_H_

#include <vector>

#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"

namespace intrinsic {
namespace geometry_legacy {

// Create a single mesh by simple appending of vertex and triangles.
// Note that the meshes are only concatenated, duplicate vertices or triangles
// are not merged.
Mesh ConcatenateMeshes(const std::vector<const Mesh*>& meshes);

// Create a single mesh by simple appending of vertex and triangles.
// Note that the meshes are only concatenated, duplicate vertices or triangles
// are not merged.
Mesh ConcatenateMeshes(const std::vector<Mesh>& meshes);

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_CONCATENATE_MESHES_H_
