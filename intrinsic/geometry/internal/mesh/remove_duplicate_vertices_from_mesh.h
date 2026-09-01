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

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_REMOVE_DUPLICATE_VERTICES_FROM_MESH_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_REMOVE_DUPLICATE_VERTICES_FROM_MESH_H_

#include <optional>

#include "intrinsic/geometry/internal/mesh/mesh.h"

namespace intrinsic::geo {

// Removes the duplicate vertices from the given 'mesh' and returns a new mesh.
// For considering two vertices as duplicates the function allows two variants:
// 1. Approximate (using a 'maybe_epsilon' > 0): vertices within a radius of
//    epsilon are considered duplicates, and replaced by the first vertex in the
//    cluster (default to have the previously used behavior).
// 2. Exact ('maybe_epsilon' = nullopt): only vertices being exactly equal are
//    considered duplicates.
Mesh RemoveDuplicateVerticesFromMesh(
    const Mesh& mesh, std::optional<double> maybe_epsilon = 1e-6);

}  // namespace intrinsic::geo
#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_REMOVE_DUPLICATE_VERTICES_FROM_MESH_H_
