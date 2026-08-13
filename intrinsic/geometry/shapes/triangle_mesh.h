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

#ifndef INTRINSIC_GEOMETRY_SHAPES_TRIANGLE_MESH_H_
#define INTRINSIC_GEOMETRY_SHAPES_TRIANGLE_MESH_H_

#include <vector>

#include "absl/log/check.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace shapes {

/** A Triangle Mesh shape */
class TriangleMesh : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::TRIANGLE_MESH;

  /**
   * Constructs a triangle mesh shape.
   * A triangle mesh is defined by a set triangles, each one consisting of three
   * vertices
   * @param vertices the (x,y,z) points that define all the mesh vertices
   * @param triangles a sequence of indices to the vertices vector. Each three
   * elements define a triangle
   */
  explicit TriangleMesh(std::vector<eigenmath::Vector3d> vertices,
                        std::vector<unsigned int> triangles)
      : ShapeBase(type),
        vertices_(std::move(vertices)),
        triangles_(std::move(triangles)) {
    CHECK_EQ(triangles_.size() % 3, 0)
        << "The triangles vector must have a size multiple of 3";
    CHECK(triangles_.empty() ||
          *std::max_element(triangles_.begin(), triangles_.end()) <
              vertices_.size())
        << "The triangles vector must contain valid indices to the vertices "
           "vector";
  }
  ~TriangleMesh() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::make_unique<TriangleMesh>(vertices_, triangles_);
  }

  /**
   * Gets the vertices that define the triangle mesh
   * @returns a vector of the vertices (x,y,z) that define the triangle mesh
   */
  const std::vector<eigenmath::Vector3d>& getVertices() const {
    return vertices_;
  }

  /**
   * Gets the triangles that define the triangle mesh
   * @returns a vector with the mesh triangles, each element contains a vertex
   * index. Each three elements define a triangle
   */
  const std::vector<unsigned int>& getTriangles() const { return triangles_; }

 private:
  std::vector<eigenmath::Vector3d> vertices_;
  std::vector<unsigned int> triangles_;
};

}  //  namespace shapes
}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_TRIANGLE_MESH_H_
