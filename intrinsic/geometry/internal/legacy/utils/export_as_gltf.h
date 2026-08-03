// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_UTILS_EXPORT_AS_GLTF_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_UTILS_EXPORT_AS_GLTF_H_

#include <string>

#include "Eigen/Core"
#include "absl/status/statusor.h"
#include "assimp/scene.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/material.h"
#include "intrinsic/geometry/internal/legacy/renderable_info/renderable_info.h"

namespace intrinsic {
namespace geometry_legacy {

// Returns a string representing the serialized glb data present in the given
// aiScene.
absl::StatusOr<std::string> ExportAiSceneAsGltf(const aiScene* scene,
                                                const Eigen::Matrix4d& trans);

// Returns a string representing the serialized glb data present in the given
// RenderableInfoData.
absl::StatusOr<std::string> ExportAsGltf(
    const RenderableInfoData& renderable_data);

// Returns a string representing the serialized glb given data with the applied
// transform.
absl::StatusOr<std::string> ExportAsGltf(std::string glb_bytes,
                                         const Eigen::Matrix4d& trans);

// Returns a string representing the serialized glb data from the given
// ExactGeometry.
absl::StatusOr<std::string> ExportAsGltf(const ExactGeometry& geometry,
                                         const Material& material = Material());

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_UTILS_EXPORT_AS_GLTF_H_
