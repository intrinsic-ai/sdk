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

#include "intrinsic/geometry/api/compute_axis_aligned_bounding_box_3d.h"

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/geometry/api/apply_transform.h"
#include "intrinsic/geometry/api/axis_aligned_bounding_box_3d.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/internal/mesh/mesh.h"
#include "intrinsic/geometry/internal/point_cloud/get_bounding_box_from_point_cloud.h"
#include "intrinsic/geometry/shapes/point_cloud.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/util/object_store/object_ref.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::geo {
namespace {

class ComputeAxisAlignedBoundingBox3dFunctor {
 public:
  template <typename Geo>
  absl::StatusOr<AxisAlignedBoundingBox3d> operator()(const Geo& geo) const {
    LOG(WARNING)
        << "WARNING: intrinsic::ComputeAxisAlignedBoundingBox3dFunctor "
           "does not support\n"
        << "operator()(const Geo&)) const, with \n"
        << "Geo = " << Demangle<Geo>();
    return absl::InvalidArgumentError(
        "ComputeAxisAlignedBoundingBox3dFunctor does not support argument "
        "type.");
  }

  template <>
  absl::StatusOr<AxisAlignedBoundingBox3d> operator()(
      const ObjectRef<Mesh>& geo) const {
    AxisAlignedBoundingBox3d result;
    for (const auto& v : geo.Value().vertices()) {
      result.ExtendBy(v);
    }
    return result;
  }

  template <>
  absl::StatusOr<AxisAlignedBoundingBox3d> operator()(
      const ObjectRef<PointCloud>& geo) const {
    return GetBoundingBoxFromPointCloud(geo.Value());
  }
};

}  // namespace

// Compute the AxisAlignedBoundingBox3d of the given Geometry.
absl::StatusOr<AxisAlignedBoundingBox3d> ComputeAxisAlignedBoundingBox3d(
    const Geometry& geo) {
  return geo.GetExactGeometry().visit(ComputeAxisAlignedBoundingBox3dFunctor());
}

// Compute the AxisAlignedBoundingBox3d of the given geometry. This is
// equivalent to, but more efficient than,
// ComputeAxisAlignedBoundingBox3d(ApplyTransform(geo, ref_t_geo));
absl::StatusOr<AxisAlignedBoundingBox3d> ComputeAxisAlignedBoundingBox3d(
    const Geometry& geo, const Pose3d& ref_t_geo) {
  INTR_ASSIGN_OR_RETURN(auto transformed_hull, ApplyTransform(geo, ref_t_geo));
  return ComputeAxisAlignedBoundingBox3d(transformed_hull);
}

}  // namespace intrinsic::geo
