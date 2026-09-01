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

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_PTS_TO_AI_SCENE_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_PTS_TO_AI_SCENE_H_

#include <memory>
#include <string>

#include "absl/status/statusor.h"
#include "assimp/scene.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/point_cloud.h"

namespace intrinsic::geo {

// Read through the content of the pts file and convert it to an aiScene.
absl::StatusOr<std::unique_ptr<aiScene>> PtsFileToAiScene(
    const std::string& file_content,
    eigenmath::Vector3d scale = eigenmath::Vector3d::Ones());

// Go through the point cloud and convert it to an aiScene.
absl::StatusOr<std::unique_ptr<aiScene>> PointCloudToAiScene(
    const PointCloud& point_cloud,
    eigenmath::Vector3d scale = eigenmath::Vector3d::Ones());

}  // namespace intrinsic::geo
#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_PTS_TO_AI_SCENE_H_
