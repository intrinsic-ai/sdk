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

#include "intrinsic/geometry/api/shapes.h"

#include <memory>
#include <optional>
#include <string>
#include <utility>

#include "absl/base/attributes.h"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/ascii.h"
#include "absl/strings/str_cat.h"
#include "assimp/material.h"
#include "assimp/scene.h"
#include "assimp/types.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/affine_transform_of_geometry.h"
#include "intrinsic/geometry/api/apply_material_properties.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/geometry_options.h"
#include "intrinsic/geometry/api/material.h"
#include "intrinsic/geometry/api/renderable.h"
#include "intrinsic/geometry/internal/legacy/mesh/io/load_ai_scene_from_file.h"
#include "intrinsic/geometry/internal/legacy/mesh/load_mesh_file_with_scale.h"
#include "intrinsic/geometry/internal/legacy/mesh/memoized_mesh_from_triangle_mesh.h"
#include "intrinsic/geometry/internal/legacy/point_cloud/pts_to_point_cloud.h"
#include "intrinsic/geometry/internal/util/export_as_gltf.h"
#include "intrinsic/geometry/proto/v1/material.pb.h"
#include "intrinsic/geometry/shapes/shape_base.h"
#include "intrinsic/geometry/shapes/shapes.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/util/status/status_macros.h"
#include "ortools/base/helpers.h"
#include "ortools/base/options.h"
#include "ortools/base/path.h"

namespace intrinsic {
namespace {

using ::intrinsic::shapes::Box;
using ::intrinsic::shapes::Capsule;
using ::intrinsic::shapes::Cylinder;
using ::intrinsic::shapes::Ellipsoid;
using ::intrinsic::shapes::Frustum;
using ::intrinsic::shapes::MeshFile;
using ::intrinsic::shapes::ShapeBase;
using ::intrinsic::shapes::ShapeType;
using ::intrinsic::shapes::TriangleMesh;

absl::StatusOr<ExactGeometry> FromShapeBaseImpl(const ShapeBase& shape,
                                                GeometryOptions options) {
  switch (shape.getType()) {
    case ShapeType::BOX: {
      return ExactGeometry{shape.get<Box>(), std::move(options)};
    }
    case ShapeType::CAPSULE: {
      return ExactGeometry{shape.get<Capsule>(), std::move(options)};
    }
    case ShapeType::CYLINDER: {
      return ExactGeometry{shape.get<Cylinder>(), std::move(options)};
    }
    case ShapeType::ELLIPSOID: {
      return ExactGeometry{shape.get<Ellipsoid>(), std::move(options)};
    }
    case ShapeType::SPHERE: {
      return ExactGeometry{shape.get<intrinsic::shapes::Sphere>(),
                           std::move(options)};
    }
    case ShapeType::FRUSTUM: {
      return ExactGeometry{shape.get<Frustum>(), std::move(options)};
    }
    case ShapeType::MESHFILE: {
      auto& mesh_file = shape.get<MeshFile>();
      const std::string extension =
          absl::AsciiStrToLower(file::Extension(mesh_file.getFilename()));
      if (extension == "pts") {
        INTR_ASSIGN_OR_RETURN(
            const auto file_content,
            file::GetContents(mesh_file.getFilename(), file::Defaults()));
        INTR_ASSIGN_OR_RETURN(auto points,
                              geometry_legacy::PtsFileToPointCloud(
                                  file_content, mesh_file.getScale()));
        return ExactGeometry{std::move(points), std::move(options)};
      } else {
        INTR_ASSIGN_OR_RETURN(
            auto mesh, geometry_legacy::LoadMeshFileWithScale(
                           mesh_file.getFilename(), mesh_file.getScale()));
        return ExactGeometry{std::move(mesh), std::move(options)};
      }
    }
    case ShapeType::TRIANGLE_MESH: {
      const auto& triangle_mesh = shape.get<TriangleMesh>();
      return ExactGeometry{
          geometry_legacy::MemoizedMeshFromTriangleMesh(triangle_mesh),
          std::move(options)};
    }
    case ShapeType::LINE_STRIP:
      ABSL_FALLTHROUGH_INTENDED;
    case ShapeType::POINT_CLOUD:
      ABSL_FALLTHROUGH_INTENDED;
    case ShapeType::AXES:
      ABSL_FALLTHROUGH_INTENDED;
    case ShapeType::CONVEX_HULL:
      ABSL_FALLTHROUGH_INTENDED;
    case ShapeType::SPHERES:
      ABSL_FALLTHROUGH_INTENDED;
    case ShapeType::DISCS:
      return absl::UnimplementedError(
          absl::StrCat("Cannot convert a shape of type: ", shape.getType(),
                       " to ExactGeometry"));
  }
}

// Check if the passed assimp scene is fully transparent. The check is
// optimistic in the sense that if there is a single partially/ fully opaque
// material defined in the scene, the check returns false, even if there may not
// be a shape in the scene that actually uses that material.
bool IsAiSceneFullyTransparent(const aiScene* scene) {
  CHECK(scene != nullptr);

  bool found_opacity = false;
  for (int i = 0; i < scene->mNumMaterials; ++i) {
    float opacity = 1;
    if (scene->mMaterials[i]->Get(AI_MATKEY_OPACITY, opacity) !=
        aiReturn_SUCCESS) {
      continue;
    }
    found_opacity = true;

    // Ignore the default material if there are other materials in the scene
    // since it won't be used.
    if (scene->mNumMaterials > 1 &&
        scene->mMaterials[i]->GetName() == aiString(AI_DEFAULT_MATERIAL_NAME)) {
      continue;
    }

    if (opacity > 0) {
      return false;
    }
  }
  // If we did not find a material with a specified opacity value, assume scene
  // is fully opaque.
  return found_opacity;
}
}  // namespace

absl::StatusOr<TransformedGeometry> ToGeometry(
    const shapes::ShapeBase& shape, std::optional<Material> material_opt,
    bool check_for_transparency, GeometryOptions options) {
  INTR_ASSIGN_OR_RETURN(ExactGeometry updated_shape,
                        FromShapeBaseImpl(shape, std::move(options)));
  std::shared_ptr<const Renderable> renderable;

  if (material_opt.has_value() && check_for_transparency &&
      material_opt->ambient[3] == 0 && material_opt->diffuse[3] == 0) {
    return absl::InvalidArgumentError(
        "Specified material is fully transparent.");
  }

  if (shape.getType() == shapes::ShapeType::MESHFILE) {
    // Since we have a mesh file we want to read the contents of the file and
    // store that as part of the renderable info.
    const auto& mesh_file = shape.get<shapes::MeshFile>();
    INTR_ASSIGN_OR_RETURN(auto ai_scene,
                          geometry_legacy::LoadAiSceneFromFile(
                              mesh_file.getFilename(), mesh_file.getScale()));
    if (check_for_transparency && IsAiSceneFullyTransparent(ai_scene.get())) {
      return absl::InvalidArgumentError(absl::StrCat(
          "Following mesh is fully transparent and cannot be rendered: ",
          mesh_file.getFilename()));
    }
    if (material_opt.has_value()) {
      INTR_RETURN_IF_ERROR(ApplyMaterialToAiScene(*ai_scene, *material_opt));
    }

    INTR_ASSIGN_OR_RETURN(
        std::string gltf_string,
        ExportAiSceneAsGltf(ai_scene.get(), Eigen::Matrix4d::Identity()));

    renderable = std::make_shared<Renderable>(gltf_string);
  } else if (material_opt.has_value()) {
    INTR_ASSIGN_OR_RETURN(std::string gltf_string,
                          ExportAsGltf(updated_shape, *material_opt));
    renderable = std::make_shared<Renderable>(gltf_string);
  }

  std::optional<intrinsic_proto::geometry::v1::MaterialProperties>
      material_properties;
  if (material_opt.has_value()) {
  }

  const bool keep_renderable = renderable != nullptr;
  return TransformedGeometry{
      Geometry(std::move(updated_shape), std::move(renderable), keep_renderable,
               std::move(material_properties), /*provenance=*/std::nullopt)};
}

}  // namespace intrinsic
