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

#include "intrinsic/geometry/internal/point_cloud/pts_to_ai_scene.h"

#include <array>
#include <cstddef>
#include <memory>
#include <optional>
#include <sstream>
#include <string>
#include <utility>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "assimp/color4.h"
#include "assimp/color4.inl"
#include "assimp/defs.h"
#include "assimp/material.h"
#include "assimp/mesh.h"
#include "assimp/metadata.h"
#include "assimp/scene.h"
#include "assimp/vector3.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/point_cloud.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::geo {
namespace {

template <typename PointFunc, typename ColorFunc>
absl::StatusOr<std::unique_ptr<aiScene>> ToAiScene(
    const size_t num_vertices, const PointFunc& point_func,
    const ColorFunc& color_func) {
  std::unique_ptr<aiScene> scene = std::make_unique<aiScene>();
  scene->mNumMeshes = 1;
  scene->mMeshes = new aiMesh* [] { new aiMesh() };

  auto mesh = scene->mMeshes[0];
  mesh->mNumVertices = num_vertices;
  mesh->mVertices = new aiVector3D[num_vertices];
  mesh->mNumFaces = num_vertices;
  mesh->mFaces = new aiFace[num_vertices];
  mesh->mColors[0] = new aiColor4D[num_vertices];
  for (int i = 0; i < num_vertices; ++i) {
    INTR_ASSIGN_OR_RETURN(const auto& p, point_func(i));

    mesh->mVertices[i] = {static_cast<ai_real>(p[0]),
                          static_cast<ai_real>(p[1]),
                          static_cast<ai_real>(p[2])};
    mesh->mFaces[i].mNumIndices = 1;
    mesh->mFaces[i].mIndices = new unsigned[]{static_cast<unsigned int>(i)};
    INTR_ASSIGN_OR_RETURN(std::optional<aiColor4D> color, color_func(i));
    if (color.has_value()) {
      mesh->mColors[0][i] = std::move(color).value();
    }
  }

  // workaround, https://github.com/assimp/assimp/issues/3778
  mesh->mPrimitiveTypes = aiPrimitiveType_POINT;

  scene->mNumMaterials = 1;
  scene->mMaterials = new aiMaterial* [] { new aiMaterial() };
  scene->mRootNode = new aiNode();
  scene->mRootNode->mNumMeshes = 1;
  scene->mRootNode->mMeshes = new unsigned[]{0};
  // workaround, https://github.com/assimp/assimp/issues/3781
  scene->mMetaData = new aiMetadata();
  return scene;
}

}  // namespace

absl::StatusOr<std::unique_ptr<aiScene>> PtsFileToAiScene(
    const std::string& file_content, eigenmath::Vector3d scale) {
  std::stringstream stream(file_content);
  // Get the vertex count from the first line;
  int num_vertices = 0;
  stream >> num_vertices;
  if (stream.bad() || stream.fail()) {
    return absl::InvalidArgumentError("Bad point count");
  }
  if (num_vertices <= 0) {
    return absl::InvalidArgumentError(
        absl::StrCat("Invalid point count: ", num_vertices));
  }
  return ToAiScene(
      num_vertices,
      [&stream, &scale](int i) -> absl::StatusOr<std::array<double, 3>> {
        double x, y, z;
        int intensity;
        stream >> x >> y >> z;
        stream >> intensity;
        if (stream.bad() || stream.fail()) {
          return absl::InvalidArgumentError(
              absl::StrCat("Bad point cloud data on line ", i + 1));
        }
        return std::array<double, 3>{x * scale[0], y * scale[1], z * scale[2]};
      },
      [&stream](int i) -> absl::StatusOr<std::optional<aiColor4D>> {
        // We have the implicit assumption that we call the other lambda first
        // per line and so this one expects to see r,g,b and not anything else.
        int r, g, b;
        stream >> r >> g >> b;
        if (stream.bad() || stream.fail()) {
          return absl::InvalidArgumentError(
              absl::StrCat("Bad point cloud data on line ", i + 1));
        }

        return aiColor4D(r, g, b, 255) / aiColor4D(255);
      });
}

absl::StatusOr<std::unique_ptr<aiScene>> PointCloudToAiScene(
    const PointCloud& point_cloud, eigenmath::Vector3d scale) {
  return ToAiScene(
      point_cloud.getPoints().size(),
      [&point_cloud, &scale](int i) -> absl::StatusOr<eigenmath::Vector3d> {
        return point_cloud.getPoints()[i].cwiseProduct(scale);
      },
      [](int) -> absl::StatusOr<std::optional<aiColor4D>> {
        return std::nullopt;
      });
}

}  // namespace intrinsic::geo
