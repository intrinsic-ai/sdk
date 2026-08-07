// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_COMPATIBILITY_SHAPES_H_
#define INTRINSIC_GEOMETRY_COMPATIBILITY_SHAPES_H_

#include <optional>

#include "absl/status/statusor.h"
#include "intrinsic/geometry/api/affine_transform_of_geometry.h"
#include "intrinsic/geometry/api/geometry_options.h"
#include "intrinsic/geometry/api/material.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace compatibility {

// Returns a TransformedGeometry instance from the given ShapeBase instance.
absl::StatusOr<TransformedGeometry> ToGeometry(
    const shapes::ShapeBase& shape,
    std::optional<Material> material_opt = std::nullopt,
    bool check_for_transparency = false,
    GeometryOptions options = GeometryOptions::Default());

}  // namespace compatibility
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_COMPATIBILITY_SHAPES_H_
