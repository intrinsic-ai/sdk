// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_SHAPES_ELLIPSOID_H_
#define INTRINSIC_GEOMETRY_SHAPES_ELLIPSOID_H_

#include "absl/log/check.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace shapes {

/** A Ellipsoid shape */
class Ellipsoid : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::ELLIPSOID;

  /**
   * Constructs an ellipsoid
   * @param radii The radii along the x, y, z axes.
   */
  explicit Ellipsoid(const eigenmath::Vector3d& radii)
      : ShapeBase(type), radii_(radii) {
    CHECK_GT(radii.minCoeff(), 0) << "Ellipsoid radii must be positive";
  }
  ~Ellipsoid() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::make_unique<Ellipsoid>(radii_);
  }

  /**
   * Gets the ellipsoid radii
   * @return The radii along the x, y, z axes.
   */
  const eigenmath::Vector3d& getRadii() const { return radii_; }

 private:
  eigenmath::Vector3d radii_;
};

}  //  namespace shapes
}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_ELLIPSOID_H_
