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

#ifndef INTRINSIC_GEOMETRY_SHAPES_DISCS_H_
#define INTRINSIC_GEOMETRY_SHAPES_DISCS_H_

#include <vector>

#include "absl/log/check.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace shapes {

/**
 * A shape representing a collection of discs
 */
class Discs : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::DISCS;

  /**
   * Constructs a shape describing a set of discs
   * A disc is defined by its centroid (3D point), a normal (3D vector) and a
   * radius
   * @param centroids a vector with a centroid for each disc
   * @param normals a vector with a normal for each disc
   * @param radiuses a vector with a radius for each disc
   */
  Discs(std::vector<eigenmath::Vector3d> centroids,
        std::vector<eigenmath::Vector3d> normals, std::vector<double> radiuses)
      : ShapeBase(type),
        centroids_(std::move(centroids)),
        normals_(std::move(normals)),
        radiuses_(std::move(radiuses)) {
    CHECK(centroids_.size() == normals_.size() &&
          centroids_.size() == radiuses_.size())
        << "Centroids, normals and radiuses vectors should have "
           "the same dimensions";
    CHECK_GT(*std::min_element(radiuses_.begin(), radiuses_.end()), 0)
        << "Discs radiuses must be all positive";
  }
  ~Discs() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::make_unique<Discs>(centroids_, normals_, radiuses_);
  }

  /**
   * Gets the centroids for all the discs
   * @returns a vector with the discs centroids
   */
  const std::vector<eigenmath::Vector3d>& getCentroids() const {
    return centroids_;
  }

  /**
   * Gets the normals for all the discs
   * @returns a vector with the discs normals
   */
  const std::vector<eigenmath::Vector3d>& getNormals() const {
    return normals_;
  }

  /**
   * Gets the radiuses for all the discs
   * @returns a vector with the discs radiuses
   */
  const std::vector<double>& getRadiuses() const { return radiuses_; }

 private:
  std::vector<eigenmath::Vector3d> centroids_;
  std::vector<eigenmath::Vector3d> normals_;
  std::vector<double> radiuses_;
};

}  //  namespace shapes
}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_DISCS_H_
