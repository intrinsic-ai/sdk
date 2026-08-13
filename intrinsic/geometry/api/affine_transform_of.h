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

#ifndef INTRINSIC_GEOMETRY_API_AFFINE_TRANSFORM_OF_H_
#define INTRINSIC_GEOMETRY_API_AFFINE_TRANSFORM_OF_H_

#include <utility>

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"

// AffineTransformOf attaches a 3D affine transform to a geometric object.
namespace intrinsic {

template <typename Shape>
class AffineTransformOf {
 public:
  explicit AffineTransformOf(Shape shape)
      : shape_(std::move(shape)),
        ref_t_shape_(eigenmath::Matrix4d::Identity()) {}

  AffineTransformOf(Shape shape, const Pose3d& ref_t_shape)
      : shape_(std::move(shape)), ref_t_shape_(ref_t_shape.matrix()) {}

  AffineTransformOf(Shape shape, const Pose3d& ref_t_shape,
                    const eigenmath::Vector3d& scale)
      : shape_(std::move(shape)),
        ref_t_shape_(ref_t_shape.matrix() *
                     eigenmath::Vector4d(scale(0), scale(1), scale(2), 1.0)
                         .asDiagonal()) {}

  AffineTransformOf(Shape shape, const eigenmath::Matrix4d& ref_t_shape)
      : shape_(std::move(shape)), ref_t_shape_(ref_t_shape) {}

  AffineTransformOf(const AffineTransformOf& other)
      : shape_(other.shape_), ref_t_shape_(other.ref_t_shape_) {}

  AffineTransformOf(AffineTransformOf&& other) noexcept
      : shape_(std::move(other.shape_)),
        ref_t_shape_(std::move(other.ref_t_shape_)) {}

  AffineTransformOf& operator=(AffineTransformOf&& other) noexcept {
    shape_ = std::move(other.shape_);
    ref_t_shape_ = std::move(other.ref_t_shape_);
    return *this;
  }

  AffineTransformOf& operator=(const AffineTransformOf& other) = default;

  const Shape& shape() const { return shape_; }
  const eigenmath::Matrix4d& ref_t_shape() const { return ref_t_shape_; }

  bool operator==(const AffineTransformOf<Shape>& other) const {
    return ref_t_shape_.isApprox(other.ref_t_shape()) && shape_ == other.shape_;
  }

  bool operator!=(const AffineTransformOf<Shape>& other) const {
    return !(*this == other);
  }

 private:
  Shape shape_;
  eigenmath::Matrix4d ref_t_shape_;
};

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_AFFINE_TRANSFORM_OF_H_
