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

#include "intrinsic/geometry/internal/mesh/concatenate_meshes.h"

#include <vector>

#include "absl/log/check.h"
#include "intrinsic/geometry/internal/mesh/mesh.h"
#include "intrinsic/geometry/internal/mesh/remove_duplicate_vertices_from_mesh.h"

namespace intrinsic::geo {

namespace {

// Returns a stable hash for deterministic sorting.
template <class T>
size_t ComputeHash(const std::vector<T>& v) {
  return std::hash<std::string_view>{}(std::string_view(
      reinterpret_cast<const char*>(v.data()), v.size() * sizeof(T)));
}

}  // namespace

Mesh ConcatenateMeshes(const std::vector<Mesh>& meshes) {
  std::vector<const Mesh*> mesh_references;
  mesh_references.reserve(meshes.size());
  for (const auto& mesh : meshes) {
    mesh_references.push_back(&mesh);
  }

  // Sorts for deterministic order. Returning the same concatenated mesh
  // irrespective of the order of input meshes enables us to use memoized
  // results in other parts of the stack (e.g., octree computation).
  std::sort(mesh_references.begin(), mesh_references.end(),
            [](const Mesh* a, const Mesh* b) -> bool {
              if (a->vertex_count() != b->vertex_count()) {
                return a->vertex_count() < b->vertex_count();
              }
              if (a->face_count() != b->face_count()) {
                return a->face_count() < b->face_count();
              }

              const size_t hash_a = ComputeHash<>(a->vertices());
              const size_t hash_b = ComputeHash<>(b->vertices());
              if (hash_a != hash_b) {
                return hash_a < hash_b;
              }

              const size_t hash_faces_a = ComputeHash<>(a->faces());
              const size_t hash_faces_b = ComputeHash<>(b->faces());
              // If face hashes are same, then the meshes are identical.
              return hash_faces_a < hash_faces_b;
            });

  return ConcatenateMeshes(mesh_references);
}

Mesh ConcatenateMeshes(const std::vector<const Mesh*>& meshes) {
  if (meshes.empty()) {
    return Mesh();
  }
  CHECK(meshes[0] != nullptr) << "mesh at index 0 is NULL";
  Mesh result = RemoveDuplicateVerticesFromMesh(*meshes[0], std::nullopt);
  for (size_t i = 1; i < meshes.size(); ++i) {
    CHECK(meshes[i] != nullptr) << "mesh at index " << i << " is NULL";
    result.Append(RemoveDuplicateVerticesFromMesh(*meshes[i], std::nullopt));
  }
  return result;
}

}  // namespace intrinsic::geo
