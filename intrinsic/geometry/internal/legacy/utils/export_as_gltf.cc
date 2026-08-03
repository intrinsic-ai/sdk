// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/internal/legacy/utils/export_as_gltf.h"

#include <memory>
#include <sstream>
#include <string>
#include <vector>

#include "absl/log/check.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "assimp/Exporter.hpp"
#include "assimp/Importer.hpp"
#include "assimp/cexport.h"
#include "assimp/color4.h"
#include "assimp/config.h"
#include "assimp/defs.h"
#include "assimp/material.h"
#include "assimp/matrix4x4.h"
#include "assimp/scene.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/material.h"
#include "intrinsic/geometry/internal/legacy/mesh/io/mesh_to_ai_scene.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/geometry/internal/legacy/point_cloud/pts_to_ai_scene.h"
#include "intrinsic/geometry/internal/legacy/renderable_info/renderable_info.h"
#include "intrinsic/util/object_store/object_ref.h"
#include "intrinsic/util/status/status_macros.h"
#include "ortools/base/path.h"
#include "tiny_gltf.h"

namespace intrinsic {
namespace geometry_legacy {
namespace {

absl::StatusOr<std::string> ScaleGltfString(absl::string_view gltf_string,
                                            const eigenmath::Vector3d& scale) {
  tinygltf::TinyGLTF loader;
  std::string err;
  std::string warn;

  tinygltf::Model model;
  if (!loader.LoadBinaryFromMemory(
          &model, &err, &warn,
          reinterpret_cast<const unsigned char*>(gltf_string.data()),
          gltf_string.size())) {
    return absl::InternalError(absl::StrCat(
        "Failed to load glTF string to glTF model. Error(s): ", err));
  }
  if (!warn.empty()) {
    LOG(WARNING) << "Warning(s) during the loading of glTF string: " << warn;
  }
  for (const auto& scene : model.scenes) {
    for (int node_index : scene.nodes) {
      if (model.nodes[node_index].matrix.size() == 16) {
        eigenmath::Matrix4d matrix;
        matrix(0, 0) = model.nodes[node_index].matrix[0];
        matrix(1, 0) = model.nodes[node_index].matrix[1];
        matrix(2, 0) = model.nodes[node_index].matrix[2];
        matrix(3, 0) = model.nodes[node_index].matrix[3];
        matrix(0, 1) = model.nodes[node_index].matrix[4];
        matrix(1, 1) = model.nodes[node_index].matrix[5];
        matrix(2, 1) = model.nodes[node_index].matrix[6];
        matrix(3, 1) = model.nodes[node_index].matrix[7];
        matrix(0, 2) = model.nodes[node_index].matrix[8];
        matrix(1, 2) = model.nodes[node_index].matrix[9];
        matrix(2, 2) = model.nodes[node_index].matrix[10];
        matrix(3, 2) = model.nodes[node_index].matrix[11];
        matrix(0, 3) = model.nodes[node_index].matrix[12];
        matrix(1, 3) = model.nodes[node_index].matrix[13];
        matrix(2, 3) = model.nodes[node_index].matrix[14];
        matrix(3, 3) = model.nodes[node_index].matrix[15];
        model.nodes[node_index].matrix[0] = matrix(0, 0) * scale[0];
        model.nodes[node_index].matrix[1] = matrix(1, 0) * scale[0];
        model.nodes[node_index].matrix[2] = matrix(2, 0) * scale[0];
        model.nodes[node_index].matrix[3] = matrix(3, 0);
        model.nodes[node_index].matrix[4] = matrix(0, 1) * scale[1];
        model.nodes[node_index].matrix[5] = matrix(1, 1) * scale[1];
        model.nodes[node_index].matrix[6] = matrix(2, 1) * scale[1];
        model.nodes[node_index].matrix[7] = matrix(3, 1);
        model.nodes[node_index].matrix[8] = matrix(0, 2) * scale[2];
        model.nodes[node_index].matrix[9] = matrix(1, 2) * scale[2];
        model.nodes[node_index].matrix[10] = matrix(2, 2) * scale[2];
        model.nodes[node_index].matrix[11] = matrix(3, 2);
        model.nodes[node_index].matrix[12] = matrix(0, 3);
        model.nodes[node_index].matrix[13] = matrix(1, 3);
        model.nodes[node_index].matrix[14] = matrix(2, 3);
        model.nodes[node_index].matrix[15] = matrix(3, 3);
      } else if (model.nodes[node_index].scale.size() == 3) {
        model.nodes[node_index].scale[0] *= scale[0];
        model.nodes[node_index].scale[1] *= scale[1];
        model.nodes[node_index].scale[2] *= scale[2];
      } else {
        model.nodes[node_index].scale.resize(3);
        model.nodes[node_index].scale[0] = scale[0];
        model.nodes[node_index].scale[1] = scale[1];
        model.nodes[node_index].scale[2] = scale[2];
      }
    }
  }

  std::stringstream ss;
  if (!loader.WriteGltfSceneToStream(&model, ss, /*prettyPrint=*/false,
                                     /*writeBinary=*/true)) {
    return absl::InternalError("Failed to save glTF string from glTF model.");
  }
  return ss.str();
}

void SetMaterialToOpaqueWhite(const aiScene& scene) {
  if (scene.HasMaterials()) {
    for (int i = 0; i < scene.mNumMaterials; ++i) {
      aiColor4D opaque_white(ai_real(1.0), ai_real(1.0), ai_real(1.0),
                             ai_real(1.0));

      aiMaterial* original_mat = scene.mMaterials[i];
      original_mat->AddProperty(&opaque_white, 1, AI_MATKEY_COLOR_DIFFUSE);
      original_mat->AddProperty(&opaque_white, 1, AI_MATKEY_COLOR_SPECULAR);
    }
  }
}

}  // namespace

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

absl::StatusOr<std::string> ExportAsGltf(
    const RenderableInfoData& renderable_data) {
  absl::string_view format = file::Extension(renderable_data.filename);
  if (format == "glb") {
    if (renderable_data.scale.isOnes()) {
      return renderable_data.content;
    } else {
      return ScaleGltfString(renderable_data.content, renderable_data.scale);
    }
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
  if (format == "pts") {
    INTR_ASSIGN_OR_RETURN(scene_ptr, PtsFileToAiScene(renderable_data.content));
    scene = scene_ptr.get();
  } else {
    // Extract a Mesh, apply scale, convert to GLTF.
    importer.SetPropertyInteger(AI_CONFIG_IMPORT_COLLADA_IGNORE_UP_DIRECTION,
                                1);
    scene = importer.ReadFileFromMemory(renderable_data.content.c_str(),
                                        renderable_data.content.size(), 0,
                                        format.data());

    if (format == "stl") {
      // Override any material in the scene to be opaque white. Assimp parses
      // non standard materialise materials for STL files. These materials are
      // not recognised by other softwares causing confusion when loading the
      // same model.
      SetMaterialToOpaqueWhite(*scene);
    }
  }

  if (scene == nullptr) {
    return absl::InvalidArgumentError("Input geometry can't be loaded!");
  }

  // Note that the renderableInfo can have a non-uniform scale.
  Eigen::Affine3d scaleTransform(Eigen::Scaling(renderable_data.scale));

  // Post apply the scaling matrix.
  eigenmath::Matrix4d rend_trans = scaleTransform.matrix();

  // Finally export the scene with the scale transform.
  return ExportAiSceneAsGltf(scene, rend_trans);
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

}  // namespace geometry_legacy
}  // namespace intrinsic
