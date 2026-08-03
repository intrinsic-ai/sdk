// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_PTS_TO_AI_SCENE_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_PTS_TO_AI_SCENE_H_

#include <memory>
#include <string>

#include "absl/status/statusor.h"
#include "assimp/scene.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/point_cloud.h"

namespace intrinsic {
namespace geometry_legacy {

// Read through the content of the pts file and convert it to an aiScene.
absl::StatusOr<std::unique_ptr<aiScene>> PtsFileToAiScene(
    const std::string& file_content,
    eigenmath::Vector3d scale = eigenmath::Vector3d::Ones());

// Go through the point cloud and convert it to an aiScene.
absl::StatusOr<std::unique_ptr<aiScene>> PointCloudToAiScene(
    const shapes::PointCloud& point_cloud,
    eigenmath::Vector3d scale = eigenmath::Vector3d::Ones());

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_PTS_TO_AI_SCENE_H_
