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

#ifndef INTRINSIC_GEOMETRY_SHAPES_ELLIPSOID_H_
#define INTRINSIC_GEOMETRY_SHAPES_ELLIPSOID_H_

#include "absl/log/check.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic::geo {

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

}  // namespace intrinsic::geo

namespace intrinsic::shapes {
using ::intrinsic::geo::Ellipsoid;
}  // namespace intrinsic::shapes

#endif  // INTRINSIC_GEOMETRY_SHAPES_ELLIPSOID_H_
