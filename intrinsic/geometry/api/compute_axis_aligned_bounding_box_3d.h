// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_API_COMPUTE_AXIS_ALIGNED_BOUNDING_BOX_3D_H_
#define INTRINSIC_GEOMETRY_API_COMPUTE_AXIS_ALIGNED_BOUNDING_BOX_3D_H_

#include "absl/status/statusor.h"
#include "intrinsic/geometry/api/axis_aligned_bounding_box_3d.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/math/pose3.h"

namespace intrinsic {

// Compute the AxisAlignedBoundingBox3d of the given Geometry.
absl::StatusOr<AxisAlignedBoundingBox3d> ComputeAxisAlignedBoundingBox3d(
    const Geometry& geo);

// Compute the AxisAlignedBoundingBox3d of the given geometry. This is
// equiavalent to, but more efficient than,
// ComputeAxisAlignedBoundingBox3d(ApplyTransform(geo, ref_t_geo));
absl::StatusOr<AxisAlignedBoundingBox3d> ComputeAxisAlignedBoundingBox3d(
    const Geometry& geo, const Pose3d& ref_t_geo);

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_COMPUTE_AXIS_ALIGNED_BOUNDING_BOX_3D_H_
