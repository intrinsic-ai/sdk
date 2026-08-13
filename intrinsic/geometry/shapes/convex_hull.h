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

#ifndef INTRINSIC_GEOMETRY_SHAPES_CONVEX_HULL_H_
#define INTRINSIC_GEOMETRY_SHAPES_CONVEX_HULL_H_

#include <vector>

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace shapes {

/** A Convex Hull shape */
class ConvexHull : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::CONVEX_HULL;

  /**
   * Constructs a convex hull shape
   * A convex hull is represented by a set of points that need to be inside.
   * @param vertices the (x,y,z) points that define the convex hull
   */
  explicit ConvexHull(std::vector<eigenmath::Vector3d> vertices)
      : ShapeBase(type), vertices_(std::move(vertices)) {}
  ~ConvexHull() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::make_unique<ConvexHull>(vertices_);
  }

  /**
   * Gets the vertices that define the convex hull
   * @returns a vector of the vertices (x,y,z) that define the convex hull
   */
  const std::vector<eigenmath::Vector3d>& getVertices() const {
    return vertices_;
  }

 private:
  std::vector<eigenmath::Vector3d> vertices_;
};

}  //  namespace shapes
}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_CONVEX_HULL_H_
