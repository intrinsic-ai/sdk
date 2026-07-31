// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_SHAPES_BOX_H_
#define INTRINSIC_GEOMETRY_SHAPES_BOX_H_

#include "absl/log/check.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace shapes {

/** A Box shape */
class Box : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::BOX;

  /**
   * Constructs a box shape
   * @param size the full size of the box (lengths on x, y, z axis)
   */
  explicit Box(const eigenmath::Vector3d& size) : ShapeBase(type), size_(size) {
    CHECK_GT(size.minCoeff(), 0) << "Box dimensions must be positive";
  }
  ~Box() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::make_unique<Box>(size_);
  }

  /**
   * Gets the box size
   * @return the full size of the box (lengths on x, y, z axis)
   */
  const eigenmath::Vector3d& getSize() const { return size_; }

 private:
  eigenmath::Vector3d size_;
};

}  //  namespace shapes
}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_BOX_H_
