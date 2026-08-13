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
