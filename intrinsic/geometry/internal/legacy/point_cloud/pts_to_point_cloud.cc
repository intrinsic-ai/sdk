// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/internal/legacy/point_cloud/pts_to_point_cloud.h"

#include <sstream>
#include <string>
#include <utility>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/shapes/point_cloud.h"

namespace intrinsic {
namespace geometry_legacy {

absl::StatusOr<shapes::PointCloud> PtsFileToPointCloud(
    const std::string& file_content) {
  return PtsFileToPointCloud(file_content, eigenmath::Vector3d::Ones());
}

absl::StatusOr<shapes::PointCloud> PtsFileToPointCloud(
    const std::string& file_content, const eigenmath::Vector3d& scale) {
  std::stringstream stream(file_content);
  // Get the vertex count from the first line;
  int num_vertices = 0;
  stream >> num_vertices;
  if (stream.bad() || stream.fail()) {
    return absl::InvalidArgumentError("Bad point count");
  }
  if (num_vertices <= 0) {
    return absl::InvalidArgumentError(
        absl::StrCat("Invalid point count: ", num_vertices));
  }
  std::vector<eigenmath::Vector3d> points;
  points.reserve(num_vertices);

  for (int i = 0; i < num_vertices; ++i) {
    float x, y, z;
    int r, g, b;
    int intensity;
    stream >> x >> y >> z;
    stream >> intensity;
    stream >> r >> g >> b;
    if (stream.bad() || stream.fail()) {
      return absl::InvalidArgumentError(
          absl::StrCat("Bad point cloud data on line ", i + 1));
    }
    points.emplace_back(x * scale.x(), y * scale.y(), z * scale.z());
  }

  return shapes::PointCloud(std::move(points), {});
}

}  // namespace geometry_legacy
}  // namespace intrinsic
