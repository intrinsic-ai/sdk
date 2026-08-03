// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/api/axis_aligned_bounding_box_3d.h"

#include <algorithm>
#include <array>
#include <cstddef>
#include <ostream>
#include <vector>

#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/proto/axis_aligned_bounding_box.pb.h"
#include "intrinsic/geometry/proto/v1/axis_aligned_bounding_box.pb.h"
#include "intrinsic/math/proto/vector3.pb.h"

namespace intrinsic {

bool AxisAlignedBoundingBox3d::DoesContain(const eigenmath::Vector3d& p) const {
  if (p.x() < eigen_aabb_.min().x()) return false;
  if (p.y() < eigen_aabb_.min().y()) return false;
  if (p.z() < eigen_aabb_.min().z()) return false;
  if (p.x() > eigen_aabb_.max().x()) return false;
  if (p.y() > eigen_aabb_.max().y()) return false;
  if (p.z() > eigen_aabb_.max().z()) return false;
  return true;
}

// Two bounding boxes do overlap if their intersection is not empty.
bool AxisAlignedBoundingBox3d::DoOverlap(
    const AxisAlignedBoundingBox3d& other) const {
  return eigen_aabb_.intersects(other.eigen_aabb_);
}

bool AxisAlignedBoundingBox3d::IsDistanceSmaller(
    const AxisAlignedBoundingBox3d& other, double distance) const {
  for (int i = 0; i < 3; ++i) {
    if (GetMin(i) > other.GetMax(i) + distance) {
      return false;
    }
    if (GetMax(i) + distance < other.GetMin(i)) {
      return false;
    }
  }
  return true;
}

bool AxisAlignedBoundingBox3d::IsDistanceSmaller(
    const eigenmath::Vector3d& point, double distance) const {
  return SquareDistanceFromPoint(point) <= (distance * distance);
}

// Return the squared distance from a point to the closest point on the bbox.
double AxisAlignedBoundingBox3d::SquareDistanceFromPoint(
    const eigenmath::Vector3d& p) const {
  DCHECK(!IsEmpty()) << "Should not perform distance check on empty bbox";
  return eigen_aabb_.squaredExteriorDistance(p);
}

std::ostream& operator<<(std::ostream& os,
                         const AxisAlignedBoundingBox3d& bbox) {
  return os << "[(" << bbox.GetMin()[0] << ", " << bbox.GetMin()[1] << ", "
            << bbox.GetMin()[2] << "), (" << bbox.GetMax()[0] << ", "
            << bbox.GetMax()[1] << ", " << bbox.GetMax()[2] << ")]";
}

bool AxisAlignedBoundingBox3d::IsEmpty() const { return eigen_aabb_.isEmpty(); }

std::array<eigenmath::Vector3d, 8> AxisAlignedBoundingBox3d::GetCorners()
    const {
  // We spell this out instead of calling GetCorner() in a loop because it has
  // produced a ~2% speed up on the distance benchmarks. This improvement is
  // significant as the distance checks are a hot path during path planning.
  // For the benchmark results see cl/327874335.
  const auto& min = eigen_aabb_.min();
  const auto& max = eigen_aabb_.max();
  return {eigenmath::Vector3d(min[0], min[1], min[2]),
          eigenmath::Vector3d(max[0], min[1], min[2]),
          eigenmath::Vector3d(max[0], max[1], min[2]),
          eigenmath::Vector3d(min[0], max[1], min[2]),
          eigenmath::Vector3d(min[0], min[1], max[2]),
          eigenmath::Vector3d(max[0], min[1], max[2]),
          eigenmath::Vector3d(max[0], max[1], max[2]),
          eigenmath::Vector3d(min[0], max[1], max[2])};
}

std::vector<Triangle> AxisAlignedBoundingBox3d::GetBoundaryTriangles() const {
  auto corners = this->GetCorners();
  std::vector<Triangle> result;
  result.reserve(12);
  result.push_back(Triangle{corners[0], corners[1], corners[2]});
  result.push_back(Triangle{corners[2], corners[3], corners[0]});

  result.push_back(Triangle{corners[0], corners[1], corners[4]});
  result.push_back(Triangle{corners[4], corners[5], corners[1]});

  result.push_back(Triangle{corners[1], corners[2], corners[5]});
  result.push_back(Triangle{corners[5], corners[6], corners[2]});

  result.push_back(Triangle{corners[2], corners[3], corners[6]});
  result.push_back(Triangle{corners[6], corners[7], corners[3]});

  result.push_back(Triangle{corners[3], corners[0], corners[7]});
  result.push_back(Triangle{corners[7], corners[4], corners[0]});

  result.push_back(Triangle{corners[4], corners[5], corners[6]});
  result.push_back(Triangle{corners[4], corners[6], corners[7]});

  return result;
}

// Creates a spatial partitioning of the bbox and returns a vector of size
// (2 ** 3) including the partitions. The partitioning is done by
// dividing the bbox along every dimension at its centroid.
AxisAlignedBoundingBox3d AxisAlignedBoundingBox3d::GetOctant(
    const eigenmath::Vector3d& split, int i) const {
  // We use thread_local here to avoid creating extra memory and to speed up the
  // performance of this very hot path used during path planning and collision
  // checking. See cl/353083446 and cl/524072502 for more details.
  thread_local std::array<double, 6> all;
  for (int dim = 0; dim < 3; ++dim) {
    if ((i >> dim) & 1) {
      // Upper half
      all[dim] = split[dim];
      all[3 + dim] = GetMax(dim);
    } else {
      // Lower half
      all[dim] = GetMin(dim);
      all[3 + dim] = split[dim];
    }
  }

  return AxisAlignedBoundingBox3d(all[0], all[1], all[2], all[3], all[4],
                                  all[5]);
}

std::array<AxisAlignedBoundingBox3d, 8> AxisAlignedBoundingBox3d::GetOctants()
    const {
  // We use thread_local here to avoid creating extra memory and to speed up the
  // performance of this very hot path used during path planning and collision
  // checking. See cl/524077992 for more details.
  thread_local std::array<AxisAlignedBoundingBox3d, 8> octants;
  eigenmath::Vector3d center = GetCenter();
  for (size_t i = 0; i < 8; ++i) {
    octants[i] = GetOctant(center, i);
  }
  return octants;
}

absl::StatusOr<intrinsic_proto::geometry::AxisAlignedBoundingBox3> ToProto(
    const AxisAlignedBoundingBox3d& bbox) {
  intrinsic_proto::geometry::AxisAlignedBoundingBox3 result;
  if (!bbox.IsEmpty()) {
    result.mutable_min()->set_x(bbox.GetMin(0));
    result.mutable_min()->set_y(bbox.GetMin(1));
    result.mutable_min()->set_z(bbox.GetMin(2));
    result.mutable_max()->set_x(bbox.GetMax(0));
    result.mutable_max()->set_y(bbox.GetMax(1));
    result.mutable_max()->set_z(bbox.GetMax(2));
  }
  return result;
}

absl::StatusOr<AxisAlignedBoundingBox3d> FromProto(
    const intrinsic_proto::geometry::AxisAlignedBoundingBox3& proto) {
  if (proto.has_min() && proto.has_max()) {
    return AxisAlignedBoundingBox3d(proto.min().x(), proto.min().y(),
                                    proto.min().z(), proto.max().x(),
                                    proto.max().y(), proto.max().z());
  } else if (!proto.has_min() && !proto.has_max()) {
    return AxisAlignedBoundingBox3d();
  } else if (!proto.has_min()) {
    return absl::InvalidArgumentError(
        "Axis aligned bounding box does not have any min point.");
  } else {  //  if (!proto.has_max()) {
    return absl::InvalidArgumentError(
        "Axis aligned bounding box does not have any max point.");
  }
}

absl::StatusOr<intrinsic_proto::geometry::v1::AxisAlignedBoundingBox3>
ToProtoV1(const AxisAlignedBoundingBox3d& bbox) {
  intrinsic_proto::geometry::v1::AxisAlignedBoundingBox3 result;
  if (!bbox.IsEmpty()) {
    result.mutable_min()->set_x(bbox.GetMin(0));
    result.mutable_min()->set_y(bbox.GetMin(1));
    result.mutable_min()->set_z(bbox.GetMin(2));
    result.mutable_max()->set_x(bbox.GetMax(0));
    result.mutable_max()->set_y(bbox.GetMax(1));
    result.mutable_max()->set_z(bbox.GetMax(2));
  }
  return result;
}

absl::StatusOr<AxisAlignedBoundingBox3d> FromProto(
    const intrinsic_proto::geometry::v1::AxisAlignedBoundingBox3& proto) {
  if (proto.has_min() && proto.has_max()) {
    return AxisAlignedBoundingBox3d(proto.min().x(), proto.min().y(),
                                    proto.min().z(), proto.max().x(),
                                    proto.max().y(), proto.max().z());
  } else if (!proto.has_min() && !proto.has_max()) {
    return AxisAlignedBoundingBox3d();
  } else if (!proto.has_min()) {
    return absl::InvalidArgumentError(
        "Axis aligned bounding box does not have any min point.");
  } else {  //  if (!proto.has_max()) {
    return absl::InvalidArgumentError(
        "Axis aligned bounding box does not have any max point.");
  }
}

}  // namespace intrinsic
