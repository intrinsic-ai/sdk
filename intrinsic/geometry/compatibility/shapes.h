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
