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

#include "intrinsic/geometry/internal/mesh/io/mesh_to_ai_scene.h"

#include <cstdint>
#include <optional>
#include <span>

#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "assimp/color4.h"
#include "assimp/material.h"
#include "assimp/mesh.h"
#include "assimp/scene.h"
#include "assimp/vector3.h"
#include "intrinsic/geometry/api/material.h"
#include "intrinsic/geometry/internal/mesh/mesh.h"

namespace intrinsic::geo {

absl::Status MeshToAiScene(const Mesh& mesh, const Material& material,
                           std::optional<std::span<const uint8_t>> colors,
                           aiScene& aiscene) {
  Material effective_material = material;
  const size_t vertex_count = mesh.vertices().size();

  if (colors.has_value()) {
    if (colors->size() != 3 && colors->size() != 3 * vertex_count) {
      return absl::InvalidArgumentError(absl::StrCat(
          "Mesh colors size (", colors->size(),
          ") is invalid. Expected 3 for uniform RGB or 3 * vertex_count (",
          3 * vertex_count, ") for per-vertex RGB."));
    }
    if (colors->size() == 3) {
      effective_material.diffuse = {
          static_cast<float>((*colors)[0]) / 255.0f,
          static_cast<float>((*colors)[1]) / 255.0f,
          static_cast<float>((*colors)[2]) / 255.0f,
          1.0f,
      };
    }
  }

  // Converts mesh to aiscene.
  aiscene.mRootNode = new aiNode();

  aiscene.mMaterials = new aiMaterial*[1];
  aiscene.mMaterials[0] = new aiMaterial();
  aiscene.mNumMaterials = 1;

  aiColor4D ambient(
      effective_material.ambient[0], effective_material.ambient[1],
      effective_material.ambient[2], effective_material.ambient[3]);
  aiscene.mMaterials[0]->AddProperty(&ambient, 1, AI_MATKEY_COLOR_AMBIENT);
  aiColor4D diffuse(
      effective_material.diffuse[0], effective_material.diffuse[1],
      effective_material.diffuse[2], effective_material.diffuse[3]);
  aiscene.mMaterials[0]->AddProperty(&diffuse, 1, AI_MATKEY_COLOR_DIFFUSE);
  aiColor4D specular(
      effective_material.specular[0], effective_material.specular[1],
      effective_material.specular[2], effective_material.specular[3]);
  aiscene.mMaterials[0]->AddProperty(&specular, 1, AI_MATKEY_COLOR_SPECULAR);
  aiColor4D emission(
      effective_material.emission[0], effective_material.emission[1],
      effective_material.emission[2], effective_material.emission[3]);
  aiscene.mMaterials[0]->AddProperty(&emission, 1, AI_MATKEY_COLOR_EMISSIVE);
  aiscene.mMaterials[0]->AddProperty(&effective_material.shininess, 1,
                                     AI_MATKEY_SHININESS);

  aiscene.mMeshes = new aiMesh*[1];
  aiscene.mMeshes[0] = new aiMesh();
  aiscene.mNumMeshes = 1;
  aiscene.mMeshes[0]->mMaterialIndex = 0;

  aiscene.mRootNode->mMeshes = new unsigned int[1];
  aiscene.mRootNode->mMeshes[0] = 0;
  aiscene.mRootNode->mNumMeshes = 1;

  aiMesh* exported_mesh = aiscene.mMeshes[0];
  exported_mesh->mPrimitiveTypes |= aiPrimitiveType_TRIANGLE;

  // Copies vertices and per-vertex colors if available.
  exported_mesh->mNumVertices = vertex_count;
  exported_mesh->mVertices = new aiVector3D[vertex_count];
  const bool has_per_vertex_colors = colors.has_value() &&
                                     colors->size() == 3 * vertex_count &&
                                     vertex_count > 0;
  if (has_per_vertex_colors) {
    exported_mesh->mColors[0] = new aiColor4D[vertex_count];
  }

  for (size_t vdx = 0; vdx < vertex_count; ++vdx) {
    const Mesh::Vertex& v = mesh.vertices()[vdx];
    exported_mesh->mVertices[vdx].Set(v[0], v[1], v[2]);
    if (has_per_vertex_colors) {
      exported_mesh->mColors[0][vdx] =
          aiColor4D(static_cast<float>((*colors)[3 * vdx]) / 255.0f,
                    static_cast<float>((*colors)[3 * vdx + 1]) / 255.0f,
                    static_cast<float>((*colors)[3 * vdx + 2]) / 255.0f, 1.0f);
    }
  }

  // Copies faces.
  exported_mesh->mNumFaces = mesh.faces().size();
  exported_mesh->mFaces = new aiFace[exported_mesh->mNumFaces];
  for (size_t fdx = 0; fdx < mesh.faces().size(); ++fdx) {
    exported_mesh->mFaces[fdx].mNumIndices = 3;
    exported_mesh->mFaces[fdx].mIndices = new unsigned int[3];
    exported_mesh->mFaces[fdx].mIndices[0] = mesh.faces()[fdx][0];
    exported_mesh->mFaces[fdx].mIndices[1] = mesh.faces()[fdx][1];
    exported_mesh->mFaces[fdx].mIndices[2] = mesh.faces()[fdx][2];
  }

  return absl::OkStatus();
}

}  // namespace intrinsic::geo
