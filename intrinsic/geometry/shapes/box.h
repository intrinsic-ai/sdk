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
