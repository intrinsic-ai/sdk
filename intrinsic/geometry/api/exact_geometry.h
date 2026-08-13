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

#ifndef INTRINSIC_GEOMETRY_API_EXACT_GEOMETRY_H_
#define INTRINSIC_GEOMETRY_API_EXACT_GEOMETRY_H_

// Within the intrinsic codebase we support a few different shape types. These
// types don't always have the same interface but we need to be able to support
// them all. The way we have achieved this is by using std::variant.
//
// There is a type defined that is a collection of the common geometry types we
// expect to use and this is named ExactGeometry.

#include <memory>
#include <variant>
#include <vector>

#include "absl/log/log.h"
#include "absl/status/statusor.h"
#include "intrinsic/geometry/api/affine_transform_of.h"
#include "intrinsic/geometry/api/geometry_options.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/geometry/shapes/box.h"
#include "intrinsic/geometry/shapes/capsule.h"
#include "intrinsic/geometry/shapes/cylinder.h"
#include "intrinsic/geometry/shapes/ellipsoid.h"
#include "intrinsic/geometry/shapes/frustum.h"
#include "intrinsic/geometry/shapes/point_cloud.h"
#include "intrinsic/geometry/shapes/shape_base.h"
#include "intrinsic/geometry/shapes/shape_ptrs.h"
#include "intrinsic/geometry/shapes/shapes.h"
#include "intrinsic/util/object_store/object_ref.h"

namespace intrinsic {

// ExactGeometry is a collection of supported shapes that are common within the
// intrinsic codebase.
class ExactGeometry : public std::enable_shared_from_this<ExactGeometry> {
 public:
  using Shape = std::variant<std::vector<TransformedPrimitiveShapePtr>,
                             ObjectRef<geometry_legacy::Mesh>,
                             ObjectRef<shapes::PointCloud>>;

  using ComputedShape = std::variant<ObjectRef<geometry_legacy::Mesh>,
                                     ObjectRef<shapes::PointCloud>>;

  // Default constructor with an empty mesh.
  ExactGeometry() = delete;
  static ExactGeometry CreateEmpty();

  // Standard Move/Copy constructors
  ExactGeometry(const ExactGeometry& other) = default;
  ExactGeometry(ExactGeometry&& other) = default;

  ExactGeometry& operator=(const ExactGeometry& other) = default;
  ExactGeometry& operator=(ExactGeometry&& other) = default;

  explicit ExactGeometry(const shapes::Box& shape,
                         GeometryOptions options = GeometryOptions::Default());
  explicit ExactGeometry(const shapes::Capsule& shape,
                         GeometryOptions options = GeometryOptions::Default());
  explicit ExactGeometry(const shapes::Cylinder& shape,
                         GeometryOptions options = GeometryOptions::Default());
  explicit ExactGeometry(const shapes::Ellipsoid& shape,
                         GeometryOptions options = GeometryOptions::Default());
  explicit ExactGeometry(const shapes::Sphere& shape,
                         GeometryOptions options = GeometryOptions::Default());
  explicit ExactGeometry(const shapes::Frustum& shape,
                         GeometryOptions options = GeometryOptions::Default());
  explicit ExactGeometry(geometry_legacy::Mesh mesh,
                         GeometryOptions options = GeometryOptions::Default());
  explicit ExactGeometry(shapes::PointCloud point_cloud,
                         GeometryOptions options = GeometryOptions::Default());
  explicit ExactGeometry(ObjectRef<geometry_legacy::Mesh> mesh,
                         GeometryOptions options = GeometryOptions::Default());
  explicit ExactGeometry(ObjectRef<shapes::PointCloud> point_cloud,
                         GeometryOptions options = GeometryOptions::Default());

  // Create an ExactGeometry from the given transformed shape. A mesh will be
  // computed from the shape (accounting for the transform).
  static absl::StatusOr<ExactGeometry> Create(
      TransformedPrimitiveShapePtr shape,
      GeometryOptions options = GeometryOptions::Default());

  static absl::StatusOr<ExactGeometry> Create(
      Shape shape, GeometryOptions options = GeometryOptions::Default());

  static absl::StatusOr<ExactGeometry> Create(
      std::vector<TransformedPrimitiveShapePtr> primitive_shapes,
      ComputedShape computed_shape,
      GeometryOptions options = GeometryOptions::Default());

  // Helper function for accessing information about the any shape internals.
  // Use of this function is discouraged outside of the geometry package.
  template <typename Func>
  auto visit(Func func) const {
    return std::visit(func, shape_);
  }

  bool HasMesh() const;
  bool HasPointCloud() const;
  bool HasPrimitiveShapes() const;

  const std::vector<TransformedPrimitiveShapePtr>& GetPrimitiveShapes() const;

  absl::StatusOr<ObjectRef<geometry_legacy::Mesh>> GetMesh() const;
  absl::StatusOr<ObjectRef<shapes::PointCloud>> GetPointCloud() const;

  // Returns the runtime options for this geometry.
  const GeometryOptions& options() const;

  bool operator==(const ExactGeometry& other) const;
  bool operator!=(const ExactGeometry& other) const;

  // A copy constructor that allows for changing the options.
  ExactGeometry(const ExactGeometry& other, GeometryOptions options);

 private:
  ExactGeometry(std::vector<TransformedPrimitiveShapePtr> primitive_shapes,
                ComputedShape computed_shape, GeometryOptions options);

  std::vector<TransformedPrimitiveShapePtr> primitive_shapes_;

  // We have a mesh or point cloud and an optional primitive set.
  // If the primivite set is present, we use the mesh as a fallback.
  ComputedShape shape_;

  // Runtime options for the shape.
  GeometryOptions options_;
};

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_EXACT_GEOMETRY_H_
