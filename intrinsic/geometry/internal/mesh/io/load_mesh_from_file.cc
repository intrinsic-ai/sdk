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

#include "intrinsic/geometry/internal/mesh/io/load_mesh_from_file.h"

#include <memory>
#include <string>

#include "absl/log/check.h"
#include "absl/status/statusor.h"
#include "assimp/scene.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/internal/mesh/io/ai_scene_to_mesh.h"
#include "intrinsic/geometry/internal/mesh/io/load_ai_scene_from_file.h"
#include "intrinsic/geometry/internal/mesh/mesh.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::geo {

absl::StatusOr<Mesh> LoadMeshFromFile(const std::string& filename,
                                      const eigenmath::Vector3d& scale) {
  INTR_ASSIGN_OR_RETURN(
      std::unique_ptr<aiScene> scene,
      LoadAiSceneFromFile(filename, eigenmath::Vector3d::Ones()));
  INTR_ASSIGN_OR_RETURN(Mesh mesh, AiSceneToMesh(*scene));
  mesh.Scale(scale);
  return mesh;
}

}  // namespace intrinsic::geo
