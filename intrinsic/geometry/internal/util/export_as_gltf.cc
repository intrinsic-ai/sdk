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

#include "intrinsic/geometry/internal/util/export_as_gltf.h"

#include <memory>
#include <string>
#include <utility>

#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "assimp/Exporter.hpp"
#include "assimp/Importer.hpp"
#include "assimp/cexport.h"
#include "assimp/cimport.h"
#include "assimp/config.h"
#include "assimp/matrix4x4.h"
#include "assimp/scene.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/material.h"
#include "intrinsic/geometry/internal/mesh/io/mesh_to_ai_scene.h"
#include "intrinsic/geometry/internal/mesh/mesh.h"
#include "intrinsic/geometry/internal/point_cloud/pts_to_ai_scene.h"
#include "intrinsic/util/object_store/object_ref.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::geo {

absl::StatusOr<std::string> ExportAiSceneAsGltf(const aiScene* scene,
                                                const Eigen::Matrix4d& trans) {
  if (scene == nullptr) {
    return absl::InvalidArgumentError("Input scene is null!");
  }
  if (scene->mRootNode == nullptr) {
    return absl::InvalidArgumentError("Input scene has no root node!");
  }

  if (!trans.isIdentity()) {
    aiMatrix4x4 aiTrans(trans(0, 0), trans(0, 1), trans(0, 2), trans(0, 3),
                        trans(1, 0), trans(1, 1), trans(1, 2), trans(1, 3),
                        trans(2, 0), trans(2, 1), trans(2, 2), trans(2, 3),
                        trans(3, 0), trans(3, 1), trans(3, 2), trans(3, 3));

    // Set the extra transform to the root node's transform.
    scene->mRootNode->mTransformation = aiTrans;
  }
  Assimp::ExportProperties props;
  props.SetPropertyBool(AI_CONFIG_EXPORT_POINT_CLOUDS, true);
  Assimp::Exporter exporter;
  const aiExportDataBlob* blob =
      exporter.ExportToBlob(scene, "glb2", 0, &props);
  if (blob == nullptr) {
    return absl::InternalError(
        absl::StrCat("aiScene cannot be saved to glTF 2 binary format: ",
                     exporter.GetErrorString()));
  }

  std::string result(static_cast<char*>(blob->data), blob->size);
  return result;
}

absl::StatusOr<std::string> ExportAsGltf(std::string glb_bytes,
                                         const Eigen::Matrix4d& trans) {
  if (trans.isIdentity()) {
    return glb_bytes;
  }

  // The importer must have the same lifetime as the scene. Otherwise we will
  // get into trouble when trying to access the scene later.
  Assimp::Importer importer;

  // When we create the scene from scratch we must manage the memory ourselves,
  // but when we create it using the importer we defer the management to the
  // importer as we don't own the resulting aiScene*. This complicated duality
  // is why we have both a unique_ptr and a regular raw pointer.
  std::unique_ptr<const aiScene> scene_ptr;
  const aiScene* scene = nullptr;
  // Extract a Mesh, apply scale, convert to GLTF.
  importer.SetPropertyInteger(AI_CONFIG_IMPORT_COLLADA_IGNORE_UP_DIRECTION, 1);
  scene = importer.ReadFileFromMemory(glb_bytes.c_str(), glb_bytes.size(), 0,
                                      "glb");

  if (scene == nullptr) {
    return absl::InvalidArgumentError("Input geometry can't be loaded!");
  }

  // Finally export the scene with the scale transform.
  return ExportAiSceneAsGltf(scene, trans);
}

absl::StatusOr<std::string> ExportAsGltf(const ExactGeometry& geometry,
                                         const Material& material) {
  // We skip the conversion to mesh for point clouds.
  if (geometry.HasPointCloud()) {
    INTR_ASSIGN_OR_RETURN(auto point_cloud, geometry.GetPointCloud());
    INTR_ASSIGN_OR_RETURN(auto scene, PointCloudToAiScene(point_cloud.Value()));
    return ExportAiSceneAsGltf(scene.get(), eigenmath::Matrix4d::Identity());
  }

  aiScene scene;
  INTR_ASSIGN_OR_RETURN(auto mesh_ref, geometry.GetMesh());
  if (const Mesh& mesh = mesh_ref.Value(); !mesh.empty()) {
    MeshToAiScene(mesh, material, &scene);
  }

  return ExportAiSceneAsGltf(&scene, Eigen::Matrix4d::Identity());
}

}  // namespace intrinsic::geo
