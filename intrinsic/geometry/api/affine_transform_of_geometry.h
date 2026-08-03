// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_API_AFFINE_TRANSFORM_OF_GEOMETRY_H_
#define INTRINSIC_GEOMETRY_API_AFFINE_TRANSFORM_OF_GEOMETRY_H_

#include "intrinsic/geometry/api/affine_transform_of.h"
#include "intrinsic/geometry/api/geometry.h"

namespace intrinsic {

// A AffineTransformedGeometry is essentially a pair of Geometry and a affine
// transformation. This allows explicit separation of the Geometry while
// keeping different or changing transforms for it. For example two spheres with
// different world coordinates can share their internal Geometry while having
// different transforms.
using TransformedGeometry = AffineTransformOf<Geometry>;

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_AFFINE_TRANSFORM_OF_GEOMETRY_H_
