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

#include "intrinsic/geometry/api/validate_mesh.h"

#include <string>

#include "absl/container/flat_hash_set.h"
#include "absl/status/status.h"
#include "absl/strings/ascii.h"
#include "absl/strings/str_join.h"
#include "absl/strings/substitute.h"
#include "assimp/material.h"
#include "assimp/scene.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/internal/mesh/io/load_ai_scene_from_buffer.h"
#include "intrinsic/geometry/internal/mesh/io/load_ai_scene_from_file.h"
#include "intrinsic/util/status/status_macros.h"
#include "ortools/base/filesystem.h"
#include "ortools/base/helpers.h"
#include "ortools/base/options.h"
#include "ortools/base/path.h"

namespace intrinsic::geo {

namespace {

absl::Status ValidateAiSceneTextures(const aiScene& ai_scene,
                                     absl::string_view filename) {
  absl::string_view mesh_dir = file::Dirname(filename);

  absl::flat_hash_set<std::string> missing_textures;
  for (unsigned int i = 0; i < ai_scene.mNumMaterials; ++i) {
    const aiMaterial* material = ai_scene.mMaterials[i];
    // Checks common texture types.
    static constexpr aiTextureType kTextureTypes[] = {
        aiTextureType_DIFFUSE,  aiTextureType_SPECULAR,  aiTextureType_AMBIENT,
        aiTextureType_EMISSIVE, aiTextureType_HEIGHT,    aiTextureType_NORMALS,
        aiTextureType_OPACITY,  aiTextureType_BASE_COLOR};

    for (aiTextureType type : kTextureTypes) {
      unsigned int count = material->GetTextureCount(type);
      for (unsigned int j = 0; j < count; ++j) {
        aiString path;
        if (material->GetTexture(type, j, &path) == AI_SUCCESS) {
          // NULL path is valid.
          if (path.C_Str() == nullptr) {
            continue;
          }

          // From upstream documentation:
          // ```
          // To test if it is an embedded texture use
          // `aiScene::GetEmbeddedTexture`. If the returned pointer is not null,
          // it is embedded and can be loaded from the data structure. If it is
          // null, search for a separate file.
          // ```
          // References:
          // https://the-asset-importer-lib-documentation.readthedocs.io/en/latest/usage/use_the_lib.html#textures
          // https://github.com/assimp/assimp/blob/9c10e7d2b57747b6576c8230bea3c2a2e9770f11/include/assimp/material.h#L825-L829
          if (ai_scene.GetEmbeddedTexture(path.C_Str()) != nullptr) {
            continue;
          }

          if (!mesh_dir.empty()) {
            // Checks if the texture exists externally.
            std::string texture_path = file::JoinPath(mesh_dir, path.C_Str());
            if (file::Exists(texture_path, file::Defaults()).ok()) {
              continue;
            }
          }

          // Invalid texture: not embedded, not present externally.
          missing_textures.insert(path.C_Str());
        }
      }
    }
  }

  if (missing_textures.empty()) {
    return absl::OkStatus();
  }

  return absl::NotFoundError(
      absl::Substitute("Missing texture files referenced from $0: $1",
                       filename.empty() ? "mesh data" : filename,
                       absl::StrJoin(missing_textures, ",")));
}

}  // namespace

absl::Status ValidateMeshFile(absl::string_view filename,
                              const eigenmath::Vector3d& scale) {
  INTR_ASSIGN_OR_RETURN(auto ai_scene,
                        LoadAiSceneFromFile(std::string(filename), scale));

  return ValidateAiSceneTextures(*ai_scene, filename);
}

absl::Status ValidateMeshData(absl::string_view mesh_data,
                              absl::string_view extension,
                              const eigenmath::Vector3d& scale) {
  INTR_ASSIGN_OR_RETURN(auto ai_scene,
                        LoadAiSceneFromBuffer(std::string(mesh_data),
                                              std::string(extension), scale));
  return ValidateAiSceneTextures(*ai_scene, "");
}

}  // namespace intrinsic::geo
