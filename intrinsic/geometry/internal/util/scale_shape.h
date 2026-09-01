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

#ifndef INTRINSIC_GEOMETRY_INTERNAL_UTIL_SCALE_SHAPE_H_
#define INTRINSIC_GEOMETRY_INTERNAL_UTIL_SCALE_SHAPE_H_

#include <memory>

#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic::geo {

// Applies the given scale factor to the input blue_shape and return a scaled
// blue_shape (which could be different types). If the scale request can't
// result in a valid blue shape, return an error.
absl::StatusOr<std::unique_ptr<ShapeBase>> ScaleShape(
    const ShapeBase& blue_shape, const eigenmath::Vector3d& scale);

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

}  // namespace intrinsic::geo
#endif  // INTRINSIC_GEOMETRY_INTERNAL_UTIL_SCALE_SHAPE_H_
