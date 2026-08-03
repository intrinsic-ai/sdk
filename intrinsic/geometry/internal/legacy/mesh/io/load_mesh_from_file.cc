// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/internal/legacy/mesh/io/load_mesh_from_file.h"

#include <memory>
#include <string>

#include "absl/log/check.h"
#include "absl/status/statusor.h"
#include "assimp/scene.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/internal/legacy/mesh/io/ai_scene_to_mesh.h"
#include "intrinsic/geometry/internal/legacy/mesh/io/load_ai_scene_from_file.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {
namespace geometry_legacy {

absl::StatusOr<Mesh> LoadMeshFromFile(const std::string& filename,
                                      const eigenmath::Vector3d& scale) {
  INTR_ASSIGN_OR_RETURN(
      std::unique_ptr<aiScene> scene,
      LoadAiSceneFromFile(filename, eigenmath::Vector3d::Ones()));
  INTR_ASSIGN_OR_RETURN(Mesh mesh, AiSceneToMesh(*scene));
  mesh.Scale(scale);
  return mesh;
}

}  // namespace geometry_legacy
}  // namespace intrinsic
