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

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_MESH_FROM_PRIMITIVES_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_MESH_FROM_PRIMITIVES_H_

#include <memory>
#include <vector>

#include "absl/status/statusor.h"
#include "intrinsic/geometry/internal/mesh/mesh.h"
#include "intrinsic/geometry/shapes/shape_base.h"
#include "intrinsic/geometry/shapes/shape_ptrs.h"

namespace intrinsic::geo {

// Creates a mesh from the given shape.
absl::StatusOr<Mesh> MeshFromPrimitive(const ShapeBase& shape);

// Creates a mesh from the given transformed shape. The shape's affine transform
// is applied to the mesh.
absl::StatusOr<Mesh> MeshFromPrimitive(
    const TransformedPrimitiveShapePtr& shape);

// Create a mesh by concatenating the meshes of the given shapes.
absl::StatusOr<std::unique_ptr<const Mesh>> MeshFromPrimitives(
    const std::vector<PrimitiveShapePtr>& geometry);

// Create a mesh by concatenating the meshes of the given transformed shapes.
absl::StatusOr<std::unique_ptr<const Mesh>> MeshFromPrimitives(
    const std::vector<TransformedPrimitiveShapePtr>& geometry);

}  // namespace intrinsic::geo
#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_MESH_FROM_PRIMITIVES_H_
