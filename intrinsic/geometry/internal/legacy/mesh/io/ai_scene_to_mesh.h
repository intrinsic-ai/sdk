// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_AI_SCENE_TO_MESH_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_AI_SCENE_TO_MESH_H_

#include "absl/status/statusor.h"
#include "assimp/scene.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"

namespace intrinsic {
namespace geometry_legacy {

absl::StatusOr<Mesh> AiSceneToMesh(const aiScene& scene);

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_AI_SCENE_TO_MESH_H_
