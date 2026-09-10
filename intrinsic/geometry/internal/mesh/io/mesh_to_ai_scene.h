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

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_MESH_TO_AI_SCENE_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_MESH_TO_AI_SCENE_H_

#include <cstdint>
#include <optional>
#include <span>

#include "absl/status/status.h"
#include "assimp/scene.h"
#include "intrinsic/geometry/api/material.h"
#include "intrinsic/geometry/internal/mesh/mesh.h"

namespace intrinsic::geo {

// Converts the given mesh and material to an assimp scene.
//
// `colors` is an optional contiguous span of uint8 RGB values formatted in
// [R0, G0, B0, R1, G1, B1, ...] order.
// - If size is 3*V (where V is mesh.vertices().size()), per-vertex colors are
//   assigned.
// - If size is 3, a uniform color is applied to the material's diffuse color.
//
// Pre-conditions:
// - If `colors` has a value, its size must be either 3 (for a uniform RGB
// color) or 3 * mesh.vertices().size() (for per-vertex RGB colors).
//
// Returns InvalidArgumentError if pre-conditions fail.
absl::Status MeshToAiScene(const Mesh& mesh, const Material& material,
                           std::optional<std::span<const uint8_t>> colors,
                           aiScene& aiscene);

}  // namespace intrinsic::geo
#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_MESH_TO_AI_SCENE_H_
