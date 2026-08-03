// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_MESH_FROM_TRIANGLE_MESH_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_MESH_FROM_TRIANGLE_MESH_H_

#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/geometry/shapes/triangle_mesh.h"

namespace intrinsic {
namespace geometry_legacy {

// Creates a mesh from a geometry triangle mesh.
Mesh MeshFromTriangleMesh(const intrinsic::shapes::TriangleMesh& triangle_mesh);

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_MESH_FROM_TRIANGLE_MESH_H_
