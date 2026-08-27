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

#include "intrinsic/geometry/internal/util/scale_shape.h"

#include <memory>
#include <utility>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/substitute.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/shapes/box.h"
#include "intrinsic/geometry/shapes/capsule.h"
#include "intrinsic/geometry/shapes/cylinder.h"
#include "intrinsic/geometry/shapes/ellipsoid.h"
#include "intrinsic/geometry/shapes/shape_base.h"
#include "intrinsic/geometry/shapes/sphere.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::geo {

absl::StatusOr<std::unique_ptr<ShapeBase>> ScaleShape(
    const ShapeBase& shape, const eigenmath::Vector3d& scale) {
  switch (shape.getType()) {
    case ShapeType::BOX: {
      const auto& box = shape.get<Box>();
      return std::make_unique<Box>(box.getSize().array() * scale.array());
    }
    case ShapeType::ELLIPSOID: {
      const auto& ellipsoid = shape.get<Ellipsoid>();
      return std::make_unique<Ellipsoid>(ellipsoid.getRadii().array() *
                                         scale.array());
    }
    case ShapeType::SPHERE: {
      const auto& sphere = shape.get<Sphere>();
      if (scale.x() == scale.y() && scale.y() == scale.z()) {
        return std::make_unique<Sphere>(scale.x() * sphere.getRadius());
      } else {
        return std::make_unique<Ellipsoid>(scale * sphere.getRadius());
      }
    }
    case ShapeType::CYLINDER: {
      // We only scale if the scale vector is uniform across the axis.
      if (scale.x() == scale.y() && scale.y() == scale.z()) {
        const auto scale_value = scale.x();
        const auto& cylinder = shape.get<Cylinder>();
        return std::make_unique<Cylinder>(cylinder.getLength() * scale_value,
                                          cylinder.getRadius() * scale_value);
      }
      break;
    }
    case ShapeType::CAPSULE: {
      // We only scale if the scale vector is uniform across the axis.
      if (scale.x() == scale.y() && scale.y() == scale.z()) {
        const auto scale_value = scale.x();
        const auto& capsule = shape.get<Capsule>();
        return std::make_unique<Capsule>(capsule.getLength() * scale_value,
                                         capsule.getRadius() * scale_value);
      }
      break;
    }
    default:
      return absl::InvalidArgumentError(absl::Substitute(
          "Cannot scale blue shape type $0", ToString(shape.getType())));
  }
  return absl::InvalidArgumentError(absl::Substitute(
      "Cannot scale blue shape type $0", ToString(shape.getType())));
}

absl::StatusOr<PrimitiveShapePtr> ScaleShape(PrimitiveShapePtr primitive,
                                             const eigenmath::Vector3d& scale) {
  if (primitive == nullptr) {
    return absl::InvalidArgumentError("Primitive shape is null");
  }

  INTR_ASSIGN_OR_RETURN(auto scaled_shape, ScaleShape(*primitive, scale));
  return std::move(scaled_shape);
}

absl::StatusOr<TransformedPrimitiveShapePtr> ScaleShape(
    TransformedPrimitiveShapePtr primitive, const eigenmath::Vector3d& scale) {
  if (primitive.shape() == nullptr) {
    return absl::InvalidArgumentError("Primitive shape is null");
  }

  INTR_ASSIGN_OR_RETURN(auto scaled_shape,
                        ScaleShape(*primitive.shape(), scale));
  return TransformedPrimitiveShapePtr(std::move(scaled_shape),
                                      primitive.ref_t_shape());
}

}  // namespace intrinsic::geo
