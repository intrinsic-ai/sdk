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

#include "intrinsic/geometry/internal/mesh/triangle_mesh_riegeli_coder.h"

#include <utility>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"

namespace intrinsic::geo {

absl::StatusOr<intrinsic_proto::geometry::v1::TriangleMesh> ToProto(
    const TriangleMesh& triangle_mesh) {
  intrinsic_proto::geometry::v1::TriangleMesh result;
  const auto& vertices = triangle_mesh.getVertices();
  const auto& triangles = triangle_mesh.getTriangles();

  result.mutable_vertices()->Reserve(vertices.size() * 3);
  for (const auto& v : vertices) {
    result.add_vertices(v.x());
    result.add_vertices(v.y());
    result.add_vertices(v.z());
  }

  result.mutable_faces()->Reserve(triangles.size());
  for (unsigned int t : triangles) {
    result.add_faces(t);
  }

  return result;
}

absl::StatusOr<TriangleMesh> ToShape(
    const intrinsic_proto::geometry::v1::TriangleMesh& proto) {
  if (proto.vertices_size() % 3 != 0) {
    return absl::DataLossError("Vertices size must be multiple of 3");
  }
  if (proto.faces_size() % 3 != 0) {
    return absl::DataLossError("Faces size must be multiple of 3");
  }

  std::vector<eigenmath::Vector3d> vertices;
  vertices.reserve(proto.vertices_size() / 3);
  for (int i = 0; i < proto.vertices_size(); i += 3) {
    vertices.push_back(
        {proto.vertices(i), proto.vertices(i + 1), proto.vertices(i + 2)});
  }

  std::vector<unsigned int> triangles;
  triangles.reserve(proto.faces_size());
  for (int i = 0; i < proto.faces_size(); ++i) {
    triangles.push_back(proto.faces(i));
  }

  return TriangleMesh(std::move(vertices), std::move(triangles));
}

}  // namespace intrinsic::geo
