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

#include "intrinsic/geometry/internal/point_cloud/point_cloud_riegeli_coder.h"

#include <utility>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/proto/v1/point_cloud.pb.h"
#include "intrinsic/geometry/shapes/point_cloud.h"

namespace intrinsic::geo {

absl::StatusOr<intrinsic_proto::geometry::v1::PointCloud> ToProto(
    const PointCloud& point_cloud) {
  intrinsic_proto::geometry::v1::PointCloud result;

  const std::vector<eigenmath::Vector3d>& points = point_cloud.getPoints();
  const std::vector<eigenmath::Vector3d>& normals = point_cloud.getNormals();

  if (!normals.empty() && points.size() != normals.size()) {
    return absl::InvalidArgumentError("Points and normals length mismatch");
  }

  result.mutable_points()->Reserve(points.size() * 3);
  for (const eigenmath::Vector3d& p : points) {
    result.add_points(p.x());
    result.add_points(p.y());
    result.add_points(p.z());
  }

  result.mutable_normals()->Reserve(normals.size() * 3);
  for (const eigenmath::Vector3d& n : normals) {
    result.add_normals(n.x());
    result.add_normals(n.y());
    result.add_normals(n.z());
  }

  return result;
}

absl::StatusOr<PointCloud> ToShape(
    const intrinsic_proto::geometry::v1::PointCloud& point_cloud) {
  if (point_cloud.points_size() % 3 != 0) {
    return absl::DataLossError(
        "Point cloud does not have the right number of points values");
  }

  if (point_cloud.normals_size() % 3 != 0) {
    return absl::DataLossError(
        "Point cloud does not have the right number of normals values");
  }

  if (point_cloud.normals_size() != 0 &&
      point_cloud.points_size() != point_cloud.normals_size()) {
    return absl::InvalidArgumentError("Points and normals length mismatch");
  }

  std::vector<eigenmath::Vector3d> points;
  points.reserve(point_cloud.points_size() / 3);
  for (int i = 0; i < point_cloud.points_size(); i += 3) {
    eigenmath::Vector3d point = {point_cloud.points(i),
                                 point_cloud.points(i + 1),
                                 point_cloud.points(i + 2)};
    points.push_back(std::move(point));
  }

  std::vector<eigenmath::Vector3d> normals;
  normals.reserve(point_cloud.normals_size() / 3);
  for (int i = 0; i < point_cloud.normals_size(); i += 3) {
    eigenmath::Vector3d normal = {point_cloud.normals(i),
                                  point_cloud.normals(i + 1),
                                  point_cloud.normals(i + 2)};
    normals.push_back(std::move(normal));
  }

  return PointCloud(std::move(points), std::move(normals));
}

}  // namespace intrinsic::geo
