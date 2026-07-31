// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_SHAPES_SPHERES_H_
#define INTRINSIC_GEOMETRY_SHAPES_SPHERES_H_

#include <memory>
#include <vector>

#include "intrinsic/geometry/shapes/shape_base.h"
#include "intrinsic/geometry/shapes/sphere.h"

namespace intrinsic {
namespace shapes {

/** A shape composed of a set of spheres */
class Spheres : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::SPHERES;

  /**
   * Constructs a shape from a set of spheres
   * @param spheres a vector of spheres
   */
  explicit Spheres(const std::vector<Sphere>& spheres)
      : ShapeBase(type), spheres_(spheres) {}
  ~Spheres() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::make_unique<Spheres>(spheres_);
  }

  const std::vector<Sphere>& getSpheres() const { return spheres_; }

 private:
  std::vector<Sphere> spheres_;
};

}  //  namespace shapes
}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_SPHERES_H_
