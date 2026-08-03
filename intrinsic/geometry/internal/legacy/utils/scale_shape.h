// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_UTILS_SCALE_SHAPE_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_UTILS_SCALE_SHAPE_H_

#include <memory>

#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace geometry_legacy {

// Applies the given scale factor to the input blue_shape and return a scaled
// blue_shape (which could be different types). If the scale request can't
// result in a valid blue shape, return an error.
absl::StatusOr<std::unique_ptr<shapes::ShapeBase>> ScaleShape(
    const shapes::ShapeBase& blue_shape, const eigenmath::Vector3d& scale);

// Applies the given scale factor to the input primitive and return a scaled
// primitive (which could be different types). If the scale request can't
// result in a valid blue shape, return an error.
absl::StatusOr<PrimitiveShapePtr> ScaleShape(PrimitiveShapePtr primitive,
                                             const eigenmath::Vector3d& scale);

// Applies the given scale factor to the input primitive and return a scaled
// primitive (which could be different types). If the scale request can't
// result in a valid blue shape, return an error.
absl::StatusOr<TransformedPrimitiveShapePtr> ScaleShape(
    TransformedPrimitiveShapePtr primitive, const eigenmath::Vector3d& scale);

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_UTILS_SCALE_SHAPE_H_
