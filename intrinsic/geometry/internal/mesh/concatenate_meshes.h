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

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_CONCATENATE_MESHES_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_CONCATENATE_MESHES_H_

#include <vector>

#include "intrinsic/geometry/internal/mesh/mesh.h"

namespace intrinsic::geo {

// Create a single mesh by simple appending of vertex and triangles.
// Note that the meshes are only concatenated, duplicate vertices or triangles
// are not merged.
Mesh ConcatenateMeshes(const std::vector<const Mesh*>& meshes);

// Create a single mesh by simple appending of vertex and triangles.
// Note that the meshes are only concatenated, duplicate vertices or triangles
// are not merged.
Mesh ConcatenateMeshes(const std::vector<Mesh>& meshes);

}  // namespace intrinsic::geo

namespace intrinsic {
using ::intrinsic::geo::ConcatenateMeshes;
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_CONCATENATE_MESHES_H_
