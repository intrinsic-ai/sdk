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

#include "intrinsic/geometry/internal/legacy/mesh/mesh_from_primitives.h"

#include <cmath>
#include <memory>
#include <utility>
#include <vector>

#include "absl/base/attributes.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "intrinsic/geometry/api/affine_transform_of.h"
#include "intrinsic/geometry/internal/legacy/mesh/concatenate_meshes.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh_primitives.h"
#include "intrinsic/geometry/shapes/shape_base.h"
#include "intrinsic/geometry/shapes/shapes.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {
namespace geometry_legacy {

namespace {

constexpr int kNbCylinderEdges = 20;
constexpr int kSphereRecursionLevel = 3;

}  // namespace

absl::StatusOr<Mesh> MeshFromPrimitive(
    const intrinsic::shapes::ShapeBase& shape) {
  switch (shape.getType()) {
    case intrinsic::shapes::ShapeType::CYLINDER: {
      const auto& cylinder = shape.get<intrinsic::shapes::Cylinder>();
      return CreateCylinder(kNbCylinderEdges, cylinder.getRadius(),
                            cylinder.getLength());
    } break;
    case intrinsic::shapes::ShapeType::CAPSULE: {
      const auto& capsule = shape.get<intrinsic::shapes::Capsule>();
      return CreateCapsule(kSphereRecursionLevel, capsule.getRadius(),
                           capsule.getLength());
    } break;
    case intrinsic::shapes::ShapeType::SPHERE: {
      const auto& sphere = shape.get<intrinsic::shapes::Sphere>();
      double r = sphere.getRadius();
      return CreateSphere(kSphereRecursionLevel, r);
    } break;
    case intrinsic::shapes::ShapeType::ELLIPSOID: {
      const auto& ellipsoid = shape.get<intrinsic::shapes::Ellipsoid>();
      const auto& r = ellipsoid.getRadii();
      return CreateEllipsoid(kSphereRecursionLevel, r[0], r[1], r[2]);
    } break;
    case intrinsic::shapes::ShapeType::BOX: {
      const auto& box = shape.get<intrinsic::shapes::Box>();
      return CreateCuboid(box.getSize() / 2);
    } break;
    case intrinsic::shapes::ShapeType::FRUSTUM: {
      const auto& frustum = shape.get<intrinsic::shapes::Frustum>();
      return CreateFrustumMesh(frustum.getXAngle(), frustum.getYAngle(),
                               frustum.getMinZDistance(),
                               frustum.getMaxZDistance());
    } break;
    case intrinsic::shapes::ShapeType::MESHFILE:
      ABSL_FALLTHROUGH_INTENDED;
    case intrinsic::shapes::ShapeType::TRIANGLE_MESH:
      ABSL_FALLTHROUGH_INTENDED;
    default:
      return absl::InvalidArgumentError(
          "Unsupported shape " +
          absl::StrCat(static_cast<int>(shape.getType())));
  }
}

absl::StatusOr<Mesh> MeshFromPrimitive(
    const TransformedPrimitiveShapePtr& shape) {
  if (shape.shape() == nullptr) {
    return absl::InvalidArgumentError("Shape is null");
  }

  INTR_ASSIGN_OR_RETURN(auto mesh, MeshFromPrimitive(*shape.shape()));
  mesh.Transform(shape.ref_t_shape());
  for (const auto& vertex : mesh.vertices()) {
    if (!std::isfinite(vertex[0]) || !std::isfinite(vertex[1]) ||
        !std::isfinite(vertex[2])) {
      return absl::InvalidArgumentError(
          "Transformed mesh from primitive shape contains non-finite "
          "vertices");
    }
  }

  return mesh;
}

absl::StatusOr<std::unique_ptr<const Mesh>> MeshFromPrimitives(
    const std::vector<PrimitiveShapePtr>& geometry) {
  if (geometry.empty()) {
    return std::make_unique<const Mesh>();
  }

  std::vector<Mesh> meshes;
  meshes.reserve(geometry.size());
  for (const auto& shape : geometry) {
    if (shape == nullptr) {
      return absl::InvalidArgumentError("Shape is null");
    }

    INTR_ASSIGN_OR_RETURN(auto new_mesh, MeshFromPrimitive(*shape));
    meshes.push_back(std::move(new_mesh));
  }

  return std::make_unique<const Mesh>(ConcatenateMeshes(meshes));
}

absl::StatusOr<std::unique_ptr<const Mesh>> MeshFromPrimitives(
    const std::vector<TransformedPrimitiveShapePtr>& geometry) {
  if (geometry.empty()) {
    return std::make_unique<const Mesh>();
  }

  std::vector<Mesh> meshes;
  meshes.reserve(geometry.size());
  for (const auto& shape : geometry) {
    INTR_ASSIGN_OR_RETURN(auto new_mesh, MeshFromPrimitive(shape));
    meshes.push_back(std::move(new_mesh));
  }

  return std::make_unique<const Mesh>(ConcatenateMeshes(meshes));
}

}  // namespace geometry_legacy
}  // namespace intrinsic
