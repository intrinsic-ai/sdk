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

#ifndef INTRINSIC_GEOMETRY_SHAPES_FRUSTUM_H_
#define INTRINSIC_GEOMETRY_SHAPES_FRUSTUM_H_

#include <memory>

#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace shapes {

/** A Frustum shape */
class Frustum : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::FRUSTUM;  // NOLINT

  /**
   * Creates a pyramid frustum shape
   * A pyramid frustum is a truncated pyramid in which its tip is cut off by a
   * plane parallel to its base. This class represents a frustum that extends
   * from its tip in the +z direction with its tip at the origin.
   * @param x_angle the angle in radians between x-z plane and the frustum
   * surface planes that intersect the x-axis.
   * @param y_angle the angle in radians between y-z plane and the frustum
   * surface planes that intersect the y-axis.
   * @param min_z_distance the distance in meters between the origin and near
   * plane in the positive z direction
   * @param max_z_distance the distance in meters between the origin and far
   * plane in the positive z direction
   */
  static absl::StatusOr<Frustum> Create(double x_angle, double y_angle,
                                        double min_z_distance,
                                        double max_z_distance) {
    if (x_angle < 0) {
      return absl::InvalidArgumentError("Frustum x angle must be non-negative");
    }
    if (x_angle >= M_PI * 0.5) {
      return absl::InvalidArgumentError(
          "Frustum x angle must be less than 90 degrees / PI/2 radians");
    }
    if (y_angle < 0) {
      return absl::InvalidArgumentError("Frustum y angle must be non-negative");
    }
    if (y_angle >= M_PI * 0.5) {
      return absl::InvalidArgumentError(
          "Frustum y angle must be less than 90 degrees / PI/2 radians");
    }
    if (min_z_distance < 0) {
      return absl::InvalidArgumentError(
          "Frustum min z distance must be non-negative");
    }
    if (max_z_distance < min_z_distance) {
      return absl::InvalidArgumentError(
          "Frustum max z distance must be greater than or equal to min z "
          "distance");
    }
    return Frustum(x_angle, y_angle, min_z_distance, max_z_distance);
  }

  Frustum() : ShapeBase(type) {}
  ~Frustum() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::unique_ptr<ShapeBase>(
        new Frustum(x_angle_, y_angle_, min_z_distance_, max_z_distance_));
  }

  /**
   * Gets the angle in radians between the x-z plane and the frustum surface
   * planes that intersect the x-axis.
   * @return the x angle of the frustum in radians
   */
  double getXAngle() const { return x_angle_; }

  /**
   * Gets the angle in radians between the y-z plane and the frustum surface
   * planes that intersect the y-axis.
   * @return the y angle of the frustum in radians
   */
  double getYAngle() const { return y_angle_; }

  /**
   * Gets the distance in meters between the origin and the plane that cuts the
   * tip of the pyramid to form a surface of the frustum. The plane is parallel
   * to the x-y plane.
   * @return the min z of the frustum along the positive z axis
   */
  double getMinZDistance() const { return min_z_distance_; }

  /**
   * Gets the distance in meters between the origin and the plane that defines
   * the base of the frustum. The plane is parallel to the x-y plane.
   * @return the max z of the frustum along the positive z axis
   */
  double getMaxZDistance() const { return max_z_distance_; }

 private:
  Frustum(double x_angle, double y_angle, double min_z_distance,
          double max_z_distance)
      : ShapeBase(type),
        x_angle_(x_angle),
        y_angle_(y_angle),
        min_z_distance_(min_z_distance),
        max_z_distance_(max_z_distance) {}

  double x_angle_;
  double y_angle_;
  double min_z_distance_;
  double max_z_distance_;
};

}  //  namespace shapes
}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_FRUSTUM_H_
