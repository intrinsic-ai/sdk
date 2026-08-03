// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_LOAD_AI_SCENE_FROM_BUFFER_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_LOAD_AI_SCENE_FROM_BUFFER_H_

#include <memory>
#include <string>

#include "absl/log/log.h"
#include "absl/status/statusor.h"
#include "assimp/scene.h"
#include "intrinsic/eigenmath/types.h"

namespace intrinsic {
namespace geometry_legacy {

// Reads an aiScene from the given file content.
//
// File type must be readable by assimp.
absl::StatusOr<std::unique_ptr<aiScene>> LoadAiSceneFromBuffer(
    const std::string& file_content, const std::string& extension,
    const eigenmath::Vector3d& scale = eigenmath::Vector3d::Ones());

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_LOAD_AI_SCENE_FROM_BUFFER_H_
