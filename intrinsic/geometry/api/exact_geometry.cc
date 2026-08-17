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

#include "intrinsic/geometry/api/exact_geometry.h"

#include <memory>
#include <optional>
#include <set>
#include <utility>
#include <variant>
#include <vector>

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/geometry/api/affine_transform_of.h"
#include "intrinsic/geometry/api/geometry_options.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh_from_primitives.h"
#include "intrinsic/geometry/internal/legacy/mesh/remove_duplicate_vertices_from_mesh.h"
#include "intrinsic/geometry/internal/legacy/point_cloud/point_cloud_riegeli_coder.h"
#include "intrinsic/geometry/proto/mesh.pb.h"
#include "intrinsic/geometry/shapes/box.h"
#include "intrinsic/geometry/shapes/capsule.h"
#include "intrinsic/geometry/shapes/cylinder.h"
#include "intrinsic/geometry/shapes/ellipsoid.h"
#include "intrinsic/geometry/shapes/frustum.h"
#include "intrinsic/geometry/shapes/point_cloud.h"
#include "intrinsic/geometry/shapes/shape_base.h"
#include "intrinsic/geometry/shapes/shapes.h"
#include "intrinsic/geometry/shapes/sphere.h"
#include "intrinsic/util/macros.h"
#include "intrinsic/util/object_store/object_ref.h"
#include "intrinsic/util/object_store/object_store.h"
#include "intrinsic/util/status/status_builder.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {

namespace {

absl::StatusOr<ObjectRef<geometry_legacy::Mesh>> MeshFromBlueShape(
    const std::vector<TransformedPrimitiveShapePtr>& shapes) {
  INTR_ASSIGN_OR_RETURN(auto mesh, geometry_legacy::MeshFromPrimitives(shapes));
  return DeDuplicate(
      geometry_legacy::RemoveDuplicateVerticesFromMesh(*mesh, std::nullopt));
}

ObjectRef<geometry_legacy::Mesh> MeshFromBlueShapeOrDie(
    const shapes::ShapeBase& shape) {
  ASSIGN_OR_DIE(auto mesh, geometry_legacy::MeshFromPrimitive(shape));
  return DeDuplicate(
      geometry_legacy::RemoveDuplicateVerticesFromMesh(mesh, std::nullopt));
}

absl::Status ValidatePrimitiveShapes(
    const std::vector<TransformedPrimitiveShapePtr>& primitive_shapes) {
  // We intentionally do not add point cloud to the list of supported shapes,
  // even though blue shapes has a point cloud type. This is because we don't
  // really consider it a primitive in the same way we do with something like
  // sphere, or box. With a sphere or box we revert to mesh when doing a lot of
  // checks, but with a point cloud we would not revert to mesh, we would use
  // the point cloud directly in those scenarios.
  static const std::set<shapes::ShapeType>* allowed_shapes = new std::set{
      shapes::ShapeType::BOX,     shapes::ShapeType::CYLINDER,
      shapes::ShapeType::SPHERE,  shapes::ShapeType::ELLIPSOID,
      shapes::ShapeType::CAPSULE, shapes::ShapeType::FRUSTUM,
  };

  for (const auto& shape : primitive_shapes) {
    if (shape.shape() == nullptr) {
      return absl::InvalidArgumentError("Primitive shape is null");
    }

    if (allowed_shapes->count(shape.shape()->getType()) == 0) {
      return ::intrinsic::InvalidArgumentErrorBuilder()
             << "Blue shape not supported for ExactGeometry: "
             << static_cast<int>(shape.shape()->getType());
    }
  }

  return absl::OkStatus();
}

}  // namespace

ExactGeometry ExactGeometry::CreateEmpty() {
  return ExactGeometry(geometry_legacy::Mesh());
}

ExactGeometry::ExactGeometry(const shapes::Box& shape, GeometryOptions options)
    : primitive_shapes_(
          {TransformedPrimitiveShapePtr(std::make_shared<shapes::Box>(shape))}),
      shape_(MeshFromBlueShapeOrDie(shape)),
      options_(std::move(options)) {}

ExactGeometry::ExactGeometry(const shapes::Capsule& shape,
                             GeometryOptions options)
    : primitive_shapes_({TransformedPrimitiveShapePtr(
          std::make_shared<shapes::Capsule>(shape))}),
      shape_(MeshFromBlueShapeOrDie(shape)),
      options_(std::move(options)) {}

ExactGeometry::ExactGeometry(const shapes::Cylinder& shape,
                             GeometryOptions options)
    : primitive_shapes_({TransformedPrimitiveShapePtr(
          std::make_shared<shapes::Cylinder>(shape))}),
      shape_(MeshFromBlueShapeOrDie(shape)),
      options_(std::move(options)) {}

ExactGeometry::ExactGeometry(const shapes::Ellipsoid& shape,
                             GeometryOptions options)
    : primitive_shapes_({TransformedPrimitiveShapePtr(
          std::make_shared<shapes::Ellipsoid>(shape))}),
      shape_(MeshFromBlueShapeOrDie(shape)),
      options_(std::move(options)) {}

ExactGeometry::ExactGeometry(const shapes::Sphere& shape,
                             GeometryOptions options)
    : primitive_shapes_({TransformedPrimitiveShapePtr(
          std::make_shared<shapes::Sphere>(shape))}),
      shape_(MeshFromBlueShapeOrDie(shape)),
      options_(std::move(options)) {}

ExactGeometry::ExactGeometry(const shapes::Frustum& shape,
                             GeometryOptions options)
    : primitive_shapes_({TransformedPrimitiveShapePtr(
          std::make_shared<shapes::Frustum>(shape))}),
      shape_(MeshFromBlueShapeOrDie(shape)),
      options_(std::move(options)) {}

ExactGeometry::ExactGeometry(geometry_legacy::Mesh mesh,
                             GeometryOptions options)
    : shape_(DeDuplicate(std::move(mesh))), options_(std::move(options)) {}

ExactGeometry::ExactGeometry(shapes::PointCloud point_cloud,
                             GeometryOptions options)
    : shape_(DeDuplicate(std::move(point_cloud))),
      options_(std::move(options)) {}

ExactGeometry::ExactGeometry(ObjectRef<geometry_legacy::Mesh> mesh,
                             GeometryOptions options)
    : shape_(std::move(mesh)), options_(std::move(options)) {}

ExactGeometry::ExactGeometry(ObjectRef<shapes::PointCloud> point_cloud,
                             GeometryOptions options)
    : shape_(std::move(point_cloud)), options_(std::move(options)) {}

ExactGeometry::ExactGeometry(
    std::vector<TransformedPrimitiveShapePtr> primitive_shapes,
    ExactGeometry::ComputedShape computed_shape, GeometryOptions options)
    : primitive_shapes_(std::move(primitive_shapes)),
      shape_(std::move(computed_shape)),
      options_(std::move(options)) {
  // Nothing to do here.
}

absl::StatusOr<ExactGeometry> ExactGeometry::Create(
    TransformedPrimitiveShapePtr shape, GeometryOptions options) {
  std::vector<TransformedPrimitiveShapePtr> shapes = {shape};
  return ExactGeometry::Create(shapes, options);
}

absl::StatusOr<ExactGeometry> ExactGeometry::Create(ExactGeometry::Shape shape,
                                                    GeometryOptions options) {
  if (std::holds_alternative<std::vector<TransformedPrimitiveShapePtr>>(
          shape)) {
    auto primitive_shapes =
        std::get<std::vector<TransformedPrimitiveShapePtr>>(shape);
    INTR_RETURN_IF_ERROR(ValidatePrimitiveShapes(primitive_shapes));
    INTR_ASSIGN_OR_RETURN(auto computed_shape,
                          MeshFromBlueShape(primitive_shapes));
    return ExactGeometry(std::move(primitive_shapes), std::move(computed_shape),
                         std::move(options));
  } else if (std::holds_alternative<ObjectRef<geometry_legacy::Mesh>>(shape)) {
    return ExactGeometry(std::get<ObjectRef<geometry_legacy::Mesh>>(shape),
                         std::move(options));
  } else if (std::holds_alternative<ObjectRef<shapes::PointCloud>>(shape)) {
    return ExactGeometry(std::get<ObjectRef<shapes::PointCloud>>(shape),
                         std::move(options));
  }

  return absl::UnimplementedError("Unknown variant type for ExactGeometry");
}

absl::StatusOr<ExactGeometry> ExactGeometry::Create(
    std::vector<TransformedPrimitiveShapePtr> primitive_shapes,
    ExactGeometry::ComputedShape computed_shape, GeometryOptions options) {
  const bool has_point_cloud =
      std::holds_alternative<ObjectRef<shapes::PointCloud>>(computed_shape);

  if (has_point_cloud && !primitive_shapes.empty()) {
    return absl::InvalidArgumentError(
        "Cannot have a point cloud with a primitive shapes.");
  }

  INTR_RETURN_IF_ERROR(ValidatePrimitiveShapes(primitive_shapes));

  return ExactGeometry(std::move(primitive_shapes), std::move(computed_shape),
                       std::move(options));
}

bool ExactGeometry::HasMesh() const {
  return std::holds_alternative<ObjectRef<geometry_legacy::Mesh>>(shape_);
}

bool ExactGeometry::HasPointCloud() const {
  return std::holds_alternative<ObjectRef<shapes::PointCloud>>(shape_);
}

bool ExactGeometry::HasPrimitiveShapes() const {
  return !primitive_shapes_.empty();
}

const std::vector<TransformedPrimitiveShapePtr>&
ExactGeometry::GetPrimitiveShapes() const {
  return primitive_shapes_;
}

absl::StatusOr<ObjectRef<geometry_legacy::Mesh>> ExactGeometry::GetMesh()
    const {
  if (std::holds_alternative<ObjectRef<geometry_legacy::Mesh>>(shape_)) {
    return std::get<ObjectRef<geometry_legacy::Mesh>>(shape_);
  }
  return intrinsic::NotFoundErrorBuilder()
         << "Can't find a Mesh in ExactGeometry with internal type index:"
         << shape_.index();
}

absl::StatusOr<ObjectRef<shapes::PointCloud>> ExactGeometry::GetPointCloud()
    const {
  if (std::holds_alternative<ObjectRef<shapes::PointCloud>>(shape_)) {
    return std::get<ObjectRef<shapes::PointCloud>>(shape_);
  }
  return intrinsic::NotFoundErrorBuilder()
         << "Can't find a PointCloud in ExactGeometry with internal type "
            "index:"
         << shape_.index();
}

const GeometryOptions& ExactGeometry::options() const { return options_; }

bool ExactGeometry::operator==(const ExactGeometry& other) const {
  return shape_ == other.shape_;
}

bool ExactGeometry::operator!=(const ExactGeometry& other) const {
  return !(*this == other);
}

ExactGeometry::ExactGeometry(const ExactGeometry& other,
                             GeometryOptions options)
    : ExactGeometry(other.primitive_shapes_, other.shape_, std::move(options)) {
}

}  // namespace intrinsic
