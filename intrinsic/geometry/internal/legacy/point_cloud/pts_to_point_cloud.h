// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_PTS_TO_POINT_CLOUD_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_PTS_TO_POINT_CLOUD_H_

#include <string>

#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/point_cloud.h"

namespace intrinsic {
namespace geometry_legacy {

absl::StatusOr<shapes::PointCloud> PtsFileToPointCloud(
    const std::string& file_content);

absl::StatusOr<shapes::PointCloud> PtsFileToPointCloud(
    const std::string& file_content, const eigenmath::Vector3d& scale);

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_PTS_TO_POINT_CLOUD_H_
