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

#ifndef INTRINSIC_GEOMETRY_SHAPES_SPHERE_H_
#define INTRINSIC_GEOMETRY_SHAPES_SPHERE_H_

#include "absl/log/check.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic::geo {

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

}  // namespace intrinsic::geo

namespace intrinsic::shapes {
using ::intrinsic::geo::Sphere;
}  // namespace intrinsic::shapes

#endif  // INTRINSIC_GEOMETRY_SHAPES_SPHERE_H_
