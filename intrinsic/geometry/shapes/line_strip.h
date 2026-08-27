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

#ifndef INTRINSIC_GEOMETRY_SHAPES_LINE_STRIP_H_
#define INTRINSIC_GEOMETRY_SHAPES_LINE_STRIP_H_

#include <memory>
#include <vector>

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic::geo {

/**
 * A series of vertices connected by straight lines. The line is not closed
 * and not filled.
 */
class LineStrip : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::LINE_STRIP;

  /**
   * Constructs a line strip
   * @param vertices the vertices in the line strip
   */
  explicit LineStrip(const std::vector<eigenmath::Vector3d>& vertices)
      : ShapeBase(type), vertices_(vertices) {}
  ~LineStrip() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::make_unique<LineStrip>(vertices_);
  }

  /**
   * @returns The vertices in the line strip
   */
  const std::vector<eigenmath::Vector3d>& getVertices() const {
    return vertices_;
  }

 private:
  std::vector<eigenmath::Vector3d> vertices_;
};

}  // namespace intrinsic::geo

namespace intrinsic::shapes {
using ::intrinsic::geo::LineStrip;
}  // namespace intrinsic::shapes

#endif  // INTRINSIC_GEOMETRY_SHAPES_LINE_STRIP_H_
