// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/internal/legacy/mesh/mesh_from_triangle_mesh.h"

#include <vector>

#include "Eigen/Core"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/geometry/shapes/triangle_mesh.h"

namespace intrinsic {
namespace geometry_legacy {

Mesh MeshFromTriangleMesh(
    const intrinsic::shapes::TriangleMesh& triangle_mesh) {
  const auto& tm_vertices = triangle_mesh.getVertices();
  const auto& tm_triangles = triangle_mesh.getTriangles();

  Mesh::VertexCollection vertices(tm_vertices.begin(), tm_vertices.end());

  Mesh::FaceCollection faces;
  if (tm_triangles.size() % 3 == 0) {
    faces.reserve(tm_triangles.size() / 3);
    for (size_t i = 0; i < tm_triangles.size(); i += 3) {
      faces.push_back(Mesh::Face(tm_triangles[i], tm_triangles[i + 1],
                                 tm_triangles[i + 2]));
    }
  }

  return Mesh(std::move(vertices), std::move(faces));
}

}  // namespace geometry_legacy
}  // namespace intrinsic
