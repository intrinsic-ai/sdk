// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/internal/legacy/mesh/memoized_mesh_from_triangle_mesh.h"

#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh_from_triangle_mesh.h"
#include "intrinsic/geometry/internal/legacy/mesh/triangle_mesh_riegeli_coder.h"
#include "intrinsic/geometry/shapes/triangle_mesh.h"
#include "intrinsic/util/object_store/memoize.h"
#include "intrinsic/util/object_store/object_ref.h"

namespace intrinsic {
namespace geometry_legacy {

ObjectRef<Mesh> MemoizedMeshFromTriangleMesh(
    const intrinsic::shapes::TriangleMesh& triangle_mesh) {
  auto func = [](const intrinsic::shapes::TriangleMesh& triangle_mesh) {
    return MeshFromTriangleMesh(triangle_mesh);
  };
  return Memoize(MemoizeOptions("MeshFromTriangleMesh"), func, triangle_mesh);
}

}  // namespace geometry_legacy
}  // namespace intrinsic
