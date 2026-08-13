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
