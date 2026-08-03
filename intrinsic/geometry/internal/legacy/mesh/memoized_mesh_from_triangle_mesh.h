// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_MEMOIZED_MESH_FROM_TRIANGLE_MESH_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_MEMOIZED_MESH_FROM_TRIANGLE_MESH_H_

#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/geometry/shapes/triangle_mesh.h"
#include "intrinsic/util/object_store/object_ref.h"

namespace intrinsic {
namespace geometry_legacy {

// Memoizes a geometry triangle mesh into the object store (as mesh).
ObjectRef<Mesh> MemoizedMeshFromTriangleMesh(
    const intrinsic::shapes::TriangleMesh& triangle_mesh);

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_MEMOIZED_MESH_FROM_TRIANGLE_MESH_H_
