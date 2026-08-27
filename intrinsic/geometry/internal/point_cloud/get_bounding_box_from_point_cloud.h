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

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_GET_BOUNDING_BOX_FROM_POINT_CLOUD_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_GET_BOUNDING_BOX_FROM_POINT_CLOUD_H_

#include "intrinsic/geometry/api/axis_aligned_bounding_box_3d.h"
#include "intrinsic/geometry/shapes/point_cloud.h"

namespace intrinsic::geo {

// Returns the bounding box that correspond to the given triangle.
AxisAlignedBoundingBox3d GetBoundingBoxFromPointCloud(
    const PointCloud& point_cloud);

}  // namespace intrinsic::geo

namespace intrinsic {
using ::intrinsic::geo::GetBoundingBoxFromPointCloud;
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_GET_BOUNDING_BOX_FROM_POINT_CLOUD_H_
