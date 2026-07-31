// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_SHAPE_MESH_FROM_PRIMITIVES_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_SHAPE_MESH_FROM_PRIMITIVES_H_

#include <memory>
#include <vector>

#include "absl/status/statusor.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/geometry/shapes/shape_base.h"
#include "intrinsic/geometry/shapes/shape_ptrs.h"

namespace intrinsic {
namespace geometry_legacy {

// Creates a mesh from the given shape.
absl::StatusOr<Mesh> MeshFromPrimitive(
    const intrinsic::shapes::ShapeBase& shape);

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

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_SHAPE_MESH_FROM_PRIMITIVES_H_
