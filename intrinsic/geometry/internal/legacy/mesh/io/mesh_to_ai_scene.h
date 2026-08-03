// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_MESH_TO_AI_SCENE_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_MESH_TO_AI_SCENE_H_

#include "assimp/scene.h"
#include "intrinsic/geometry/api/material.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"

namespace intrinsic {
namespace geometry_legacy {

void MeshToAiScene(const Mesh& mesh, const Material& material,
                   aiScene* aiscene);

}  // namespace geometry_legacy
}  // namespace intrinsic
#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_MESH_TO_AI_SCENE_H_
