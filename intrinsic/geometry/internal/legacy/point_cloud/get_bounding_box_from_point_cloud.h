// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_GET_BOUNDING_BOX_FROM_POINT_CLOUD_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_GET_BOUNDING_BOX_FROM_POINT_CLOUD_H_

#include "intrinsic/geometry/api/axis_aligned_bounding_box_3d.h"
#include "intrinsic/geometry/shapes/point_cloud.h"

namespace intrinsic {
namespace geometry_legacy {

// Returns the bounding box that correspond to the given triangle.
AxisAlignedBoundingBox3d GetBoundingBoxFromPointCloud(
    const shapes::PointCloud& point_cloud);

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_GET_BOUNDING_BOX_FROM_POINT_CLOUD_H_
