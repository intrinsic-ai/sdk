// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_SHAPES_SHAPE_PTRS_H_
#define INTRINSIC_GEOMETRY_SHAPES_SHAPE_PTRS_H_

#include <memory>

#include "intrinsic/geometry/api/affine_transform_of.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {

// Helper aliases for some common shape pointer types
using PrimitiveShapePtr = std::shared_ptr<const shapes::ShapeBase>;
using TransformedPrimitiveShapePtr = AffineTransformOf<PrimitiveShapePtr>;

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_SHAPE_PTRS_H_
