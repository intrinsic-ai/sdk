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

#ifndef INTRINSIC_GEOMETRY_SHAPES_CYLINDER_H_
#define INTRINSIC_GEOMETRY_SHAPES_CYLINDER_H_

#include "absl/log/check.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic::geo {

// A Cylinder shape
// The cylinder is centered at the origin and aligned with the Z-axis
class Cylinder : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::CYLINDER;

  /**
   * Constructs a cylinder shape
   * A cylinder is represented by it's length and radius
   * @param length the length of the cylinder
   * @param radius the radius of the cylinder
   */
  Cylinder(double length, double radius)
      : ShapeBase(type), length_(length), radius_(radius) {
    CHECK_GT(length, 0) << "Cylinder length must be positive";
    CHECK_GT(radius, 0) << "Cylinder radius must be positive";
  }
  ~Cylinder() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::make_unique<Cylinder>(length_, radius_);
  }

  /**
   * Gets the length of the cylinder
   * @returns the length of the cylinder
   */
  double getLength() const { return length_; }

  /**
   * Gets the radius of the cylinder
   * @returns the radius of the cylinder
   */
  double getRadius() const { return radius_; }

 private:
  double length_;
  double radius_;
};

}  // namespace intrinsic::geo

namespace intrinsic::shapes {
using ::intrinsic::geo::Cylinder;
}  // namespace intrinsic::shapes

#endif  // INTRINSIC_GEOMETRY_SHAPES_CYLINDER_H_
