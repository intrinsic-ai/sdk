// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/internal/legacy/point_cloud/get_bounding_box_from_point_cloud.h"

#include "intrinsic/geometry/api/axis_aligned_bounding_box_3d.h"
#include "intrinsic/geometry/shapes/point_cloud.h"

namespace intrinsic {
namespace geometry_legacy {

AxisAlignedBoundingBox3d GetBoundingBoxFromPointCloud(
    const shapes::PointCloud& point_cloud) {
  AxisAlignedBoundingBox3d bbox;
  for (const auto& point : point_cloud.getPoints()) {
    bbox.ExtendBy(point);
  }
  return bbox;
}

}  // namespace geometry_legacy
}  // namespace intrinsic
