// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_SHAPES_AXES_H_
#define INTRINSIC_GEOMETRY_SHAPES_AXES_H_

#include "absl/log/check.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace shapes {

/** Represents a set of a coordinate axes: X, Y, and Z. */
class Axes : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::AXES;

  /**
   * Constructs the axes
   * @param length The length of each axis.
   */
  explicit Axes(double length) : ShapeBase(type), length_(length) {
    CHECK_GT(length, 0) << "Length must be positive";
  }
  ~Axes() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::make_unique<Axes>(length_);
  }

  double getLength() const { return length_; }

 private:
  // The length of each axis.
  double length_;
};

}  //  namespace shapes
}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_AXES_H_
