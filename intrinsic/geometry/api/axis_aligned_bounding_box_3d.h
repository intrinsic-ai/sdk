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

#ifndef INTRINSIC_GEOMETRY_API_AXIS_ALIGNED_BOUNDING_BOX_3D_H_
#define INTRINSIC_GEOMETRY_API_AXIS_ALIGNED_BOUNDING_BOX_3D_H_

#include <array>
#include <ostream>
#include <type_traits>
#include <utility>
#include <vector>

#include "absl/log/check.h"
#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/triangle.h"
#include "intrinsic/geometry/proto/axis_aligned_bounding_box.pb.h"
#include "intrinsic/geometry/proto/v1/axis_aligned_bounding_box.pb.h"
#include "intrinsic/marshal/riegeli_proto_coder.h"

namespace intrinsic::geo {

class AxisAlignedBoundingBox3d {
 public:
  AxisAlignedBoundingBox3d() = default;

  AxisAlignedBoundingBox3d(double xmin, double ymin, double zmin, double xmax,
                           double ymax, double zmax)
      : eigen_aabb_(eigenmath::Vector3d(xmin, ymin, zmin),
                    eigenmath::Vector3d(xmax, ymax, zmax)) {}

  // Using two template parameters allows constructing the bounding box even
  // when min and max are passed as different point-like types (e.g., mixed
  // Eigen vectors and std::array), preventing compilation failures from strict
  // type matching.
  template <class Point1, class Point2>
  AxisAlignedBoundingBox3d(const Point1& min, const Point2& max)
      : eigen_aabb_(InitBox(min, max)) {}

  AxisAlignedBoundingBox3d(const std::array<double, 3>& min,
                           const std::array<double, 3>& max)
      : eigen_aabb_(eigenmath::Vector3d(min[0], min[1], min[2]),
                    eigenmath::Vector3d(max[0], max[1], max[2])) {}

  void ExtendBy(const AxisAlignedBoundingBox3d& bbox) {
    eigen_aabb_.extend(bbox.eigen_aabb_);
  }
  void ExtendBy(const std::array<double, 3>& p) {
    eigen_aabb_.extend(eigenmath::Vector3d(p[0], p[1], p[2]));
  }
  template <typename Point>
  void ExtendBy(const Point& p) {
    if constexpr (std::is_same_v<Point, eigenmath::Vector3d>) {
      eigen_aabb_.extend(p);
    } else {
      eigen_aabb_.extend(eigenmath::Vector3d(p[0], p[1], p[2]));
    }
  }

  bool DoesContain(const eigenmath::Vector3d& p) const;

  // Two bounding boxes do overlap if their intersection is not empty.
  bool DoOverlap(const AxisAlignedBoundingBox3d& other) const;

  // Returns true if the distance between this bounding box and the given
  // bounding box is below the given distance. If distance == 0 then this
  // returns the same result as DoOverlap.
  bool IsDistanceSmaller(const AxisAlignedBoundingBox3d& other,
                         double distance) const;
  // Returns true if the distance between this bounding box and the given
  // point is below the given distance. If distance == 0 then this returns the
  // same result as DoesContain.
  bool IsDistanceSmaller(const eigenmath::Vector3d& point,
                         double distance) const;

  double GetMin(int i) const { return eigen_aabb_.min()[i]; }
  double GetMax(int i) const { return eigen_aabb_.max()[i]; }
  eigenmath::Vector3d GetMin() const { return eigen_aabb_.min(); }
  eigenmath::Vector3d GetMax() const { return eigen_aabb_.max(); }
  eigenmath::Vector3d GetDiagonal() const {
    CHECK(!this->IsEmpty())
        << "Diagonal of empty AxisAlignedBoundingBox3d requested!";
    return eigen_aabb_.sizes();
  }

  eigenmath::Vector3d GetCenter() const {
    CHECK(!this->IsEmpty())
        << "Center of empty AxisAlignedBoundingBox3d requested!";
    return GetCenterNoEmptyCheck();
  }

  eigenmath::Vector3d GetCenterNoEmptyCheck() const {
    return eigen_aabb_.center();
  }

  bool operator==(const AxisAlignedBoundingBox3d& other) const {
    return eigen_aabb_.min() == other.eigen_aabb_.min() &&
           eigen_aabb_.max() == other.eigen_aabb_.max();
  }
  bool operator!=(const AxisAlignedBoundingBox3d& other) const {
    return !(*this == other);
  }

  // Return the corners of this AxisAlignedBoundingBox3d in order: 0...3 ->
  // Lower SW, SE, NE, NW 4...7 -> Upper SW, SE, NE, NW This matches an
  // enumeration of the Rect corners for the lower face of the
  // AxisAlignedBoundingBox3d and then the upper face of the
  // AxisAlignedBoundingBox3d.
  //
  // NOTE: Changing the order of enum items will lead to errors in code that
  // uses enum order as iteration order, so do not do it.
  enum Corner {
    kLowerSW,
    kLowerSE,
    kLowerNE,
    kLowerNW,
    kUpperSW,
    kUpperSE,
    kUpperNE,
    kUpperNW
  };
  std::vector<Triangle> GetBoundaryTriangles() const;
  std::array<eigenmath::Vector3d, 8> GetCorners() const;
  AxisAlignedBoundingBox3d GetOctant(const eigenmath::Vector3d& split,
                                     int index) const;
  std::array<AxisAlignedBoundingBox3d, 8> GetOctants() const;

  // Return the squared distance from a point to the closest point on the bbox.
  double SquareDistanceFromPoint(const eigenmath::Vector3d& p) const;

  // Gets the radius of a minimum bounding sphere around this box.
  double GetBoundingSphereRadius() const { return GetDiagonal().norm() * 0.5; }

  friend std::ostream& operator<<(std::ostream& os,
                                  const AxisAlignedBoundingBox3d& bbox);

  bool IsEmpty() const;

 private:
  template <class Point>
  static eigenmath::Vector3d AsVector3d(const Point& p) {
    if constexpr (std::is_same_v<Point, eigenmath::Vector3d>) {
      return p;
    } else {
      return eigenmath::Vector3d(p[0], p[1], p[2]);
    }
  }

  template <class Point1, class Point2>
  static Eigen::AlignedBox3d InitBox(const Point1& min, const Point2& max) {
    return Eigen::AlignedBox3d(AsVector3d(min), AsVector3d(max));
  }

  Eigen::AlignedBox3d eigen_aabb_;
};
absl::StatusOr<intrinsic_proto::geometry::AxisAlignedBoundingBox3> ToProto(
    const AxisAlignedBoundingBox3d& bbox);
absl::StatusOr<AxisAlignedBoundingBox3d> FromProto(
    const intrinsic_proto::geometry::AxisAlignedBoundingBox3& proto);

absl::StatusOr<intrinsic_proto::geometry::v1::AxisAlignedBoundingBox3>
ToProtoV1(const AxisAlignedBoundingBox3d& bbox);
absl::StatusOr<AxisAlignedBoundingBox3d> FromProto(
    const intrinsic_proto::geometry::v1::AxisAlignedBoundingBox3& proto);

}  // namespace intrinsic::geo

namespace intrinsic {
REGISTER_RIEGELI_PROTO_CODER_EXPLICIT(
    geo::AxisAlignedBoundingBox3d,
    intrinsic_proto::geometry::AxisAlignedBoundingBox3, geo::ToProto,
    geo::FromProto);

using ::intrinsic::geo::AxisAlignedBoundingBox3d;
using ::intrinsic::geo::FromProto;
using ::intrinsic::geo::ToProto;
using ::intrinsic::geo::ToProtoV1;
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_AXIS_ALIGNED_BOUNDING_BOX_3D_H_
