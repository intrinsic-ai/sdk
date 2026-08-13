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

#ifndef INTRINSIC_GEOMETRY_SHAPES_MESH_FILE_H_
#define INTRINSIC_GEOMETRY_SHAPES_MESH_FILE_H_

#include "absl/log/check.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/shape_base.h"

namespace intrinsic {
namespace shapes {

class MeshFile : public ShapeBase {
 public:
  static constexpr const ShapeType type = ShapeType::MESHFILE;

  // Constructs a shape defined by a mesh file.
  //
  // Note that the value of is_convex is abstracted from whether the actual mesh
  // is convex, and is more about how the collision engine should handle this
  // object.
  //
  // Setting convex_ = true when the mesh is not actually convex will likely
  // lead to the collision engine failing in an opaque way (for instance
  // hanging during the importing of that object).
  explicit MeshFile(
      const std::string& filename, bool is_convex = false,
      const eigenmath::Vector3d& scale = eigenmath::Vector3d::Constant(1.0))
      : ShapeBase(type),
        filename_(filename),
        convex_(is_convex),
        scale_(scale) {
    CHECK_GT(scale.minCoeff(), 0) << "Mesh scale elements must be all positive";
  }
  ~MeshFile() override = default;

  std::unique_ptr<ShapeBase> clone() const override {
    return std::make_unique<MeshFile>(filename_, convex_, scale_);
  }

  // Returns the mesh file name.
  const std::string& getFilename() const { return filename_; }

  // Returns the mesh scale.
  const eigenmath::Vector3d& getScale() const { return scale_; }

  // Returns true if the mesh mesh should be processed as convex.
  bool isConvex() const { return convex_; }

 private:
  std::string filename_;
  bool convex_;
  eigenmath::Vector3d scale_;
};

}  //  namespace shapes
}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_MESH_FILE_H_
