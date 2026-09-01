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

#ifndef INTRINSIC_GEOMETRY_SHAPES_AXES_H_
#define INTRINSIC_GEOMETRY_SHAPES_AXES_H_

#include "absl/log/check.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic::geo {

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

}  // namespace intrinsic::geo
#endif  // INTRINSIC_GEOMETRY_SHAPES_AXES_H_
