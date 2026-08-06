// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/api/compute_axis_aligned_bounding_box_3d.h"

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/geometry/api/apply_transform.h"
#include "intrinsic/geometry/api/axis_aligned_bounding_box_3d.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/geometry/internal/legacy/point_cloud/get_bounding_box_from_point_cloud.h"
#include "intrinsic/geometry/shapes/point_cloud.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/util/object_store/object_ref.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {
namespace {

using geometry_legacy::GetBoundingBoxFromPointCloud;
using geometry_legacy::Mesh;
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
      const ObjectRef<shapes::PointCloud>& geo) const {
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

}  // namespace intrinsic
