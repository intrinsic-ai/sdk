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

#ifndef INTRINSIC_GEOMETRY_SHAPES_SPHERES_H_
#define INTRINSIC_GEOMETRY_SHAPES_SPHERES_H_

#include <memory>
#include <vector>

#include "intrinsic/geometry/shapes/shape_base.h"
#include "intrinsic/geometry/shapes/sphere.h"

namespace intrinsic::geo {

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

}  // namespace intrinsic::geo
#endif  // INTRINSIC_GEOMETRY_SHAPES_SPHERES_H_
