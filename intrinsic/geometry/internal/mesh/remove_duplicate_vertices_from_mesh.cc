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

#include "intrinsic/geometry/internal/mesh/remove_duplicate_vertices_from_mesh.h"

#include <math.h>
#include <stdint.h>

#include <optional>
#include <set>
#include <utility>
#include <vector>

#include "Eigen/Core"
#include "absl/container/flat_hash_map.h"
#include "absl/meta/type_traits.h"
#include "intrinsic/eigenmath/eigen_matrix_hash.h"
#include "intrinsic/geometry/internal/mesh/mesh.h"

namespace intrinsic::geo {

namespace {

// Removal of duplicate vertices with exact equality of vertex coefficients.
Mesh RemoveDuplicateVerticesFromMeshExact(const Mesh& mesh) {
  std::vector<Mesh::Vertex> unique_vertices;
  unique_vertices.reserve(mesh.vertex_count());

  using HashType = eigenmath::EigenMatrixHash<eigenmath::Vector3d>;
  absl::flat_hash_map<HashType, int> vertex_to_index;

  std::vector<int> old_to_new(mesh.vertex_count());
  for (int i = 0; i < mesh.vertex_count(); ++i) {
    const auto& v = mesh.vertex(i);
    auto [it, inserted] =
        vertex_to_index.try_emplace(HashType(v), unique_vertices.size());
    if (inserted) {
      unique_vertices.push_back(v);
    }
    old_to_new[i] = it->second;
  }

  Mesh::FaceCollection new_faces;
  new_faces.reserve(mesh.face_count());
  for (const auto& f : mesh.faces()) {
    new_faces.emplace_back(old_to_new[f[0]], old_to_new[f[1]],
                           old_to_new[f[2]]);
  }
  return Mesh(std::move(unique_vertices), std::move(new_faces));
}

// Removal of duplicate vertices where vertices within the given precision of
// 'epsilon' are merged. That precision is used for both hashing the vertex
// coefficients (after rounding to 'epsilon' precision) and finding matching
// vertices within a radius 'epsilon'.
Mesh RemoveDuplicateVerticesFromMeshApproximate(const Mesh& mesh,
                                                const double epsilon) {
  using VertexMatrix = Eigen::Matrix<double, Eigen::Dynamic, Eigen::Dynamic>;
  VertexMatrix old_vertices(mesh.vertex_count(), 3);
  for (int row = 0; row < mesh.vertex_count(); ++row) {
    old_vertices.row(row)[0] = mesh.vertices()[row][0];
    old_vertices.row(row)[1] = mesh.vertices()[row][1];
    old_vertices.row(row)[2] = mesh.vertices()[row][2];
  }

  using MatrixHashingType =
      Eigen::Matrix<int64_t, Eigen::Dynamic, Eigen::Dynamic>;
  using HashType = eigenmath::EigenMatrixHash<MatrixHashingType>;
  using IndexType = VertexMatrix::Scalar;

  // Rounding to multiples of epsilon.
  const double inverse_epsilon = 1.0 / epsilon;
  Eigen::Matrix<double_t, Eigen::Dynamic, Eigen::Dynamic> rounded_vertices =
      ((inverse_epsilon * old_vertices).array().floor());
  MatrixHashingType multiples =
      rounded_vertices.cast<MatrixHashingType::Scalar>();

  // Extracting all non-duplicate vertices.

  std::vector<IndexType> new_indices(mesh.vertex_count());
  std::vector<Mesh::Vertex> new_vertices;
  new_vertices.reserve(mesh.vertex_count());
  absl::flat_hash_map<HashType, std::set<IndexType>> buckets;
  for (IndexType i = 0; i < mesh.vertex_count(); i++) {
    // Computing hashes for all relevant buckets.
    std::vector<HashType> vertex_hashes{
        HashType(multiples.row(i)),
        HashType(multiples.row(i) + Eigen::Matrix<int64_t, 1, 3>{0, 0, 1}),
        HashType(multiples.row(i) + Eigen::Matrix<int64_t, 1, 3>{0, 1, 0}),
        HashType(multiples.row(i) + Eigen::Matrix<int64_t, 1, 3>{0, 1, 1}),
        HashType(multiples.row(i) + Eigen::Matrix<int64_t, 1, 3>{1, 0, 0}),
        HashType(multiples.row(i) + Eigen::Matrix<int64_t, 1, 3>{1, 0, 1}),
        HashType(multiples.row(i) + Eigen::Matrix<int64_t, 1, 3>{1, 1, 0}),
        HashType(multiples.row(i) + Eigen::Matrix<int64_t, 1, 3>{1, 1, 1}),
    };

    // Finding matches within the found buckets (if present).
    const double squared_epsilon = epsilon * epsilon;
    bool found_matching_vertex = false;
    for (auto hash_it = vertex_hashes.begin();
         hash_it != vertex_hashes.end() && !found_matching_vertex; hash_it++) {
      if (auto bucket_it = buckets.find(*hash_it); bucket_it != buckets.end()) {
        for (auto match : bucket_it->second) {
          if ((mesh.vertex(i) - new_vertices[match]).squaredNorm() <
              squared_epsilon) {
            new_indices[i] = match;
            found_matching_vertex = true;
            break;
          }
        }
      }
      if (found_matching_vertex) break;
    }

    // Adding the vertex to the output and the buckets iff no match was found.
    if (!found_matching_vertex) {
      const IndexType new_index = new_vertices.size();
      new_vertices.push_back(mesh.vertex(i));
      new_indices[i] = new_index;
      for (const auto& vertex_hash : vertex_hashes)
        buckets[vertex_hash].insert(new_index);
    }
  }
  new_vertices.shrink_to_fit();

  // Copy and correct vertex indices.
  Mesh::FaceCollection new_faces;
  new_faces.reserve(mesh.face_count());
  for (const Mesh::Face& f : mesh.faces()) {
    new_faces.emplace_back(new_indices[f[0]], new_indices[f[1]],
                           new_indices[f[2]]);
  }
  return Mesh(std::move(new_vertices), std::move(new_faces));
}

}  // namespace

Mesh RemoveDuplicateVerticesFromMesh(const Mesh& mesh,
                                     std::optional<double> maybe_epsilon) {
  // Check for and return empty or otherwise invalid meshes.
  if (mesh.empty()) return mesh.Clone();

  Mesh result = mesh.Clone();
  if (maybe_epsilon && maybe_epsilon > 0) {
    // Vertices within epsilon are considered duplicates.
    result =
        RemoveDuplicateVerticesFromMeshApproximate(mesh, maybe_epsilon.value());
  }
  return RemoveDuplicateVerticesFromMeshExact(result);
}

}  // namespace intrinsic::geo
