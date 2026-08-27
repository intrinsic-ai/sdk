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

#include <array>

#include "Eigen/Core"
#include "assimp/color4.h"
#include "assimp/material.h"
#include "assimp/mesh.h"
#include "assimp/scene.h"
#include "assimp/vector3.h"
#include "intrinsic/geometry/api/material.h"
#include "intrinsic/geometry/internal/mesh/mesh.h"

namespace intrinsic::geo {

void MeshToAiScene(const Mesh& mesh, const Material& material,
                   aiScene* aiscene) {
  // convert mesh to aiscene
  aiscene->mRootNode = new aiNode();

  aiscene->mMaterials = new aiMaterial*[1];
  aiscene->mMaterials[0] = new aiMaterial();
  aiscene->mNumMaterials = 1;

  aiColor4D ambient(material.ambient[0], material.ambient[1],
                    material.ambient[2], material.ambient[3]);
  aiscene->mMaterials[0]->AddProperty(&ambient, 1, AI_MATKEY_COLOR_AMBIENT);
  aiColor4D diffuse(material.diffuse[0], material.diffuse[1],
                    material.diffuse[2], material.diffuse[3]);
  aiscene->mMaterials[0]->AddProperty(&diffuse, 1, AI_MATKEY_COLOR_DIFFUSE);
  aiColor4D specular(material.specular[0], material.specular[1],
                     material.specular[2], material.specular[3]);
  aiscene->mMaterials[0]->AddProperty(&specular, 1, AI_MATKEY_COLOR_SPECULAR);
  aiColor4D emission(material.emission[0], material.emission[1],
                     material.emission[2], material.emission[3]);
  aiscene->mMaterials[0]->AddProperty(&emission, 1, AI_MATKEY_COLOR_EMISSIVE);
  aiscene->mMaterials[0]->AddProperty(&material.shininess, 1,
                                      AI_MATKEY_SHININESS);

  aiscene->mMeshes = new aiMesh*[1];
  aiscene->mMeshes[0] = nullptr;
  aiscene->mNumMeshes = 1;

  aiscene->mMeshes[0] = new aiMesh();
  aiscene->mMeshes[0]->mMaterialIndex = 0;

  aiscene->mRootNode->mMeshes = new unsigned int[1];
  aiscene->mRootNode->mMeshes[0] = 0;
  aiscene->mRootNode->mNumMeshes = 1;

  aiMesh* exported_mesh = aiscene->mMeshes[0];
  exported_mesh->mPrimitiveTypes |= aiPrimitiveType_TRIANGLE;

  // copy vertices
  exported_mesh->mNumVertices = mesh.vertices().size();
  exported_mesh->mVertices = new aiVector3D[exported_mesh->mNumVertices];
  for (int vdx = 0; vdx < mesh.vertices().size(); ++vdx) {
    Mesh::Vertex v = mesh.vertices()[vdx];
    exported_mesh->mVertices[vdx].Set(v[0], v[1], v[2]);
  }
  // copy faces
  exported_mesh->mNumFaces = mesh.faces().size();
  exported_mesh->mFaces = new aiFace[exported_mesh->mNumFaces];
  for (int fdx = 0; fdx < mesh.faces().size(); ++fdx) {
    exported_mesh->mFaces[fdx].mNumIndices = 3;
    exported_mesh->mFaces[fdx].mIndices = new unsigned int[3];
    exported_mesh->mFaces[fdx].mIndices[0] = mesh.faces()[fdx][0];
    exported_mesh->mFaces[fdx].mIndices[1] = mesh.faces()[fdx][1];
    exported_mesh->mFaces[fdx].mIndices[2] = mesh.faces()[fdx][2];
  }
}

}  // namespace intrinsic::geo
