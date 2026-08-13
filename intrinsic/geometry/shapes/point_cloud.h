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

#ifndef INTRINSIC_GEOMETRY_SHAPES_POINT_CLOUD_H_
#define INTRINSIC_GEOMETRY_SHAPES_POINT_CLOUD_H_

#include <memory>
#include <utility>
#include <vector>

#include "absl/log/check.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace shapes {

/** A Point cloud shape */
class PointCloud : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::POINT_CLOUD;

  /**
   * Constructs a point cloud shape
   * A point cloud is defined by a set of 3D points, with (possibly) their 3D
   * normals
   * @param points a vector with the point cloud 3D points
   * @param normals a vector with a normal for each point cloud point
   */
  explicit PointCloud(std::vector<eigenmath::Vector3d>&& points,
                      std::vector<eigenmath::Vector3d>&& normals =
                          std::vector<eigenmath::Vector3d>())
      : ShapeBase(type),
        points_(std::move(points)),
        normals_(std::move(normals)) {
    CHECK(normals_.empty() || points_.size() == normals_.size())
        << "Points and normals vectors should have the same dimensions";
  }

  explicit PointCloud(const std::vector<eigenmath::Vector3d>& points,
                      const std::vector<eigenmath::Vector3d>& normals =
                          std::vector<eigenmath::Vector3d>())
      : PointCloud(std::vector<eigenmath::Vector3d>(points),
                   std::vector<eigenmath::Vector3d>(normals)) {}

  ~PointCloud() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::make_unique<PointCloud>(points_, normals_);
  }

  void Scale(const eigenmath::Vector3d& scale) {
    // Update the points based on the scale.
    for (auto& point : points_) {
      point = eigenmath::Vector3d(point[0] * scale[0], point[1] * scale[1],
                                  point[2] * scale[2]);
    }

    // Update the normals based on the scale and ensure they are still
    // normalized.
    for (auto& normal : normals_) {
      normal = eigenmath::Vector3d(normal[0] * scale[0], normal[1] * scale[1],
                                   normal[2] * scale[2]);
      normal.normalize();
    }
  }

  /**
   * Gets the point cloud 3D points
   * @returns a vector with the point cloud 3D points
   */

  const std::vector<eigenmath::Vector3d>& getPoints() const { return points_; }

  /**
   * Gets point normals
   * @returns a vector with a normal for each 3D point
   */
  const std::vector<eigenmath::Vector3d>& getNormals() const {
    return normals_;
  }

  /** Checks if the shape defines normals.
   * @returns true if normals have been defined
   */
  bool hasNormals() const { return (normals_.size() == points_.size()); }

 private:
  std::vector<eigenmath::Vector3d> points_;
  std::vector<eigenmath::Vector3d> normals_;
};

}  //  namespace shapes
}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_POINT_CLOUD_H_
