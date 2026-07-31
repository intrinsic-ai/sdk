// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_SHAPES_SPHERE_H_
#define INTRINSIC_GEOMETRY_SHAPES_SPHERE_H_

#include "absl/log/check.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace shapes {

/** A Sphere shape */
class Sphere : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::SPHERE;

  /**
   * Constructs a sphere shape
   * A sphere is represented by it's radius
   * @param radius the radius of the sphere
   */
  explicit Sphere(double radius) : ShapeBase(type), radius_(radius) {
    CHECK_GT(radius, 0) << "Sphere radius must be positive";
  }
  ~Sphere() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::make_unique<Sphere>(radius_);
  }

  /**
   * Gets the radius of the sphere
   * @returns the radius of the sphere
   */
  double getRadius() const { return radius_; }

 private:
  double radius_;
};

}  //  namespace shapes
}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_SPHERE_H_
