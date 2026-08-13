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

#include "intrinsic/geometry/api/renderable_generation.h"

#include <memory>
#include <string>
#include <utility>

#include "absl/status/statusor.h"
#include "assimp/color4.h"
#include "assimp/material.h"
#include "assimp/mesh.h"
#include "assimp/scene.h"
#include "assimp/vector3.h"
#include "intrinsic/geometry/api/apply_material_properties.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/renderable.h"
#include "intrinsic/geometry/internal/legacy/utils/export_as_gltf.h"
#include "intrinsic/util/status/ret_check.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {

absl::StatusOr<std::shared_ptr<const Renderable>> GetOrGenerateRenderable(
    const Geometry& geo) {
  std::shared_ptr<const Renderable> renderable = geo.GetRenderable();
  if (renderable == nullptr) {
    INTR_ASSIGN_OR_RETURN(std::string glb_string, geometry_legacy::ExportAsGltf(
                                                      geo.GetExactGeometry()));
    renderable = std::make_shared<const Renderable>(std::move(glb_string));
  }

  return renderable;
}

absl::StatusOr<std::shared_ptr<const Renderable>>
GenerateRenderableWithMaterialOverrides(const Geometry& geo) {
  INTR_ASSIGN_OR_RETURN(std::shared_ptr<const Renderable> renderable,
                        GetOrGenerateRenderable(geo));
  INTR_RET_CHECK(renderable != nullptr);

  // If we have material overrides apply them here.
  if (geo.material_properties().has_value()) {
    INTR_ASSIGN_OR_RETURN(
        auto updated_glb,
        ApplyMaterialPropertiesToGlb(renderable->GetGLBString(),
                                     geo.material_properties().value()));
    renderable = std::make_shared<const Renderable>(std::move(updated_glb));
  }

  return renderable;
}

absl::StatusOr<std::shared_ptr<const Renderable>>
GenerateRenderableForExactGeometry(
    const ExactGeometry& exact_geometry,
    const std::optional<intrinsic_proto::geometry::v1::MaterialProperties>&
        material_properties) {
  INTR_ASSIGN_OR_RETURN(std::string glb_string,
                        geometry_legacy::ExportAsGltf(exact_geometry));
  if (material_properties.has_value()) {
    INTR_ASSIGN_OR_RETURN(
        glb_string,
        ApplyMaterialPropertiesToGlb(glb_string, material_properties.value()));
  }
  return std::make_shared<const Renderable>(std::move(glb_string));
}

absl::StatusOr<Geometry> EnsureRenderableIsAvailable(Geometry geo) {
  if (geo.GetRenderable() != nullptr) {
    return std::move(geo);
  }

  INTR_ASSIGN_OR_RETURN(std::string glb_string,
                        geometry_legacy::ExportAsGltf(geo.GetExactGeometry()));
  auto renderable = std::make_shared<const Renderable>(std::move(glb_string));

  // Because we are just generating the renderable, we don't need to change the
  // provenance.
  return Geometry(geo.GetExactGeometry(), std::move(renderable),
                  /*keep_renderable=*/false, geo.material_properties(),
                  geo.provenance());
}

absl::Status RenderableGenerator::AddGeometry(
    const ExactGeometry& exact_geometry,
    const std::optional<MaterialProperties>& material_properties) {
  INTR_ASSIGN_OR_RETURN(auto mesh_ref, exact_geometry.GetMesh());
  if (const geometry_legacy::Mesh& mesh = mesh_ref.Value(); !mesh.empty()) {
    AddMesh(mesh, material_properties);
  } else {
    return absl::InvalidArgumentError("The given mesh is empty");
  }
  return absl::OkStatus();
}

void RenderableGenerator::AddMesh(
    const geometry_legacy::Mesh& mesh,
    const std::optional<MaterialProperties>& material_properties) {
  // Copy vertices
  auto new_ai_mesh = std::make_unique<aiMesh>();
  new_ai_mesh->mPrimitiveTypes = aiPrimitiveType_TRIANGLE;
  new_ai_mesh->mNumVertices = mesh.vertex_count();
  new_ai_mesh->mVertices = new aiVector3D[new_ai_mesh->mNumVertices];
  for (size_t vdx = 0; vdx < mesh.vertex_count(); ++vdx) {
    const auto& v = mesh.vertices()[vdx];
    new_ai_mesh->mVertices[vdx].Set(v[0], v[1], v[2]);
  }

  // Copy faces
  new_ai_mesh->mNumFaces = mesh.face_count();
  new_ai_mesh->mFaces = new aiFace[new_ai_mesh->mNumFaces];
  for (size_t fdx = 0; fdx < mesh.face_count(); ++fdx) {
    const auto& face = mesh.face(fdx);
    new_ai_mesh->mFaces[fdx].mNumIndices = 3;
    new_ai_mesh->mFaces[fdx].mIndices = new unsigned int[3];
    new_ai_mesh->mFaces[fdx].mIndices[0] = face[0];
    new_ai_mesh->mFaces[fdx].mIndices[1] = face[1];
    new_ai_mesh->mFaces[fdx].mIndices[2] = face[2];
  }

  new_ai_mesh->mMaterialIndex =
      GetOrAddMaterial(material_properties.value_or(MaterialProperties{}));
  meshes_.push_back(std::move(new_ai_mesh));
}

size_t RenderableGenerator::GetOrAddMaterial(
    const MaterialProperties& material) {
  auto it = material_to_index_.find(material);
  if (it != material_to_index_.end()) {
    return it->second;
  }

  auto ai_material = std::make_unique<aiMaterial>();

  // Set the ai_material properties:
  aiString name("user_specified_material");
  ai_material->AddProperty(&name, AI_MATKEY_NAME);
  aiColor4D baseColor(material.color[0], material.color[1], material.color[2],
                      1.0f);
  ai_material->AddProperty(&baseColor, 1, AI_MATKEY_COLOR_DIFFUSE);
  float metalness = material.metalness;
  ai_material->AddProperty(&metalness, 1, AI_MATKEY_METALLIC_FACTOR);
  float roughness = material.roughness;
  ai_material->AddProperty(&roughness, 1, AI_MATKEY_ROUGHNESS_FACTOR);
  // NOTE - assimp should automatically add the required GLTF extension for
  // transmission when this key is set.
  float transmission = material.transmission;
  ai_material->AddProperty(&transmission, 1, AI_MATKEY_TRANSMISSION_FACTOR);

  materials_.push_back(std::move(ai_material));
  size_t index = materials_.size() - 1;
  material_to_index_[material] = index;
  return index;
}

absl::StatusOr<std::shared_ptr<const Renderable>>
RenderableGenerator::Finish() {
  if (meshes_.empty()) {
    return absl::InvalidArgumentError(
        "No geometry added to RenderableGenerator.");
  }

  aiScene scene;
  scene.mRootNode = new aiNode();

  // Move meshes
  scene.mNumMeshes = meshes_.size();
  scene.mMeshes = new aiMesh*[scene.mNumMeshes];
  for (size_t i = 0; i < meshes_.size(); ++i) {
    scene.mMeshes[i] = meshes_[i].release();
  }
  meshes_.clear();

  // Move materials
  scene.mNumMaterials = materials_.size();
  scene.mMaterials = new aiMaterial*[scene.mNumMaterials];
  for (size_t i = 0; i < materials_.size(); ++i) {
    scene.mMaterials[i] = materials_[i].release();
  }
  materials_.clear();

  // Attach all meshes to the root node
  scene.mRootNode->mNumMeshes = scene.mNumMeshes;
  scene.mRootNode->mMeshes = new unsigned int[scene.mNumMeshes];
  for (unsigned int i = 0; i < scene.mNumMeshes; ++i) {
    scene.mRootNode->mMeshes[i] = i;
  }

  INTR_ASSIGN_OR_RETURN(auto glb_string,
                        geometry_legacy::ExportAiSceneAsGltf(
                            &scene, Eigen::Matrix4d::Identity()));
  return std::make_shared<const Renderable>(std::move(glb_string));
}

}  // namespace intrinsic
