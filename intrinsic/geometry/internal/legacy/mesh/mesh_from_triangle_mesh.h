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
