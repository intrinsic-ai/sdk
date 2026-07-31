// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_SHAPES_CAPSULE_H_
#define INTRINSIC_GEOMETRY_SHAPES_CAPSULE_H_

#include "absl/log/check.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace shapes {

// A Capsule shape
// The capsule is centered at the origin and aligned with the Z-axis
class Capsule : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::CAPSULE;

  /**
   * Constructs a capsule shape
   * A capsule is represented by it's length and radius
   * @param length the length of the capsule
   * @param radius the radius of the capsule
   */
  Capsule(double length, double radius)
      : ShapeBase(type), length_(length), radius_(radius) {
    CHECK_GT(length, 0) << "Capsule length must be positive";
    CHECK_GT(radius, 0) << "Capsule radius must be positive";
  }
  ~Capsule() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::make_unique<Capsule>(length_, radius_);
  }

  /**
   * Gets the length of the capsule
   * @returns the length of the capsule
   */
  double getLength() const { return length_; }

  /**
   * Gets the radius of the capsule
   * @returns the radius of the capsule
   */
  double getRadius() const { return radius_; }

 private:
  double length_;
  double radius_;
};

}  //  namespace shapes
}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_CAPSULE_H_
