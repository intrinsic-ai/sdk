// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/internal/legacy/utils/scale_shape.h"

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

namespace intrinsic {
namespace geometry_legacy {

absl::StatusOr<std::unique_ptr<shapes::ShapeBase>> ScaleShape(
    const shapes::ShapeBase& shape, const eigenmath::Vector3d& scale) {
  switch (shape.getType()) {
    case shapes::ShapeType::BOX: {
      const auto& box = shape.get<shapes::Box>();
      return std::make_unique<shapes::Box>(box.getSize().array() *
                                           scale.array());
    }
    case shapes::ShapeType::ELLIPSOID: {
      const auto& ellipsoid = shape.get<shapes::Ellipsoid>();
      return std::make_unique<shapes::Ellipsoid>(ellipsoid.getRadii().array() *
                                                 scale.array());
    }
    case shapes::ShapeType::SPHERE: {
      const auto& sphere = shape.get<shapes::Sphere>();
      if (scale.x() == scale.y() && scale.y() == scale.z()) {
        return std::make_unique<shapes::Sphere>(scale.x() * sphere.getRadius());
      } else {
        return std::make_unique<shapes::Ellipsoid>(scale * sphere.getRadius());
      }
    }
    case shapes::ShapeType::CYLINDER: {
      // We only scale if the scale vector is uniform across the axis.
      if (scale.x() == scale.y() && scale.y() == scale.z()) {
        const auto scale_value = scale.x();
        const auto& cylinder = shape.get<shapes::Cylinder>();
        return std::make_unique<shapes::Cylinder>(
            cylinder.getLength() * scale_value,
            cylinder.getRadius() * scale_value);
      }
      break;
    }
    case shapes::ShapeType::CAPSULE: {
      // We only scale if the scale vector is uniform across the axis.
      if (scale.x() == scale.y() && scale.y() == scale.z()) {
        const auto scale_value = scale.x();
        const auto& capsule = shape.get<shapes::Capsule>();
        return std::make_unique<shapes::Capsule>(
            capsule.getLength() * scale_value,
            capsule.getRadius() * scale_value);
      }
      break;
    }
    default:
      return absl::InvalidArgumentError(
          absl::Substitute("Cannot scale blue shape type $0",
                           shapes::ToString(shape.getType())));
  }
  return absl::InvalidArgumentError(absl::Substitute(
      "Cannot scale blue shape type $0", shapes::ToString(shape.getType())));
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

}  // namespace geometry_legacy
}  // namespace intrinsic
