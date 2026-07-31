// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_SHAPES_LINE_STRIP_H_
#define INTRINSIC_GEOMETRY_SHAPES_LINE_STRIP_H_

#include <memory>
#include <vector>

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace shapes {

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

}  //  namespace shapes
}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_LINE_STRIP_H_
