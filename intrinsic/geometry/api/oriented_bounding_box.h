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

#ifndef INTRINSIC_GEOMETRY_API_ORIENTED_BOUNDING_BOX_H_
#define INTRINSIC_GEOMETRY_API_ORIENTED_BOUNDING_BOX_H_

#include <ostream>

#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/proto/oriented_bounding_box.pb.h"
#include "intrinsic/marshal/riegeli_proto_coder.h"
#include "intrinsic/math/pose3.h"

namespace intrinsic::geo {
template <typename Scalar>
class OrientedBoundingBox3 {
 public:
  OrientedBoundingBox3(const eigenmath::Vector3<Scalar> dimensions,
                       const Pose3<Scalar>& ref_t_box)
      : dimensions_(dimensions), ref_t_box_(ref_t_box) {}

  // Returns the dimensions of the bounding box.
  const eigenmath::Vector3<Scalar>& Dimensions() const { return dimensions_; }
  // Returns the dimensions of the bounding box.
  eigenmath::Vector3<Scalar>& Dimensions() { return dimensions_; }
  // Returns the pose of the center of the bounding box in the reference frame.
  const Pose3<Scalar>& RefTBox() const { return ref_t_box_; }
  // Returns the pose of the center of the bounding box in the reference frame.
  Pose3<Scalar>& RefTBox() { return ref_t_box_; }
  // Returns the pose of the center of the bounding box in the reference frame.
  template <typename OtherScalar>
  OrientedBoundingBox3<OtherScalar> Cast() const {
    return OrientedBoundingBox3<OtherScalar>(
        dimensions_.template cast<OtherScalar>(),
        ref_t_box_.template cast<OtherScalar>());
  }

 private:
  // Full-length dimensions of the bounding box.
  eigenmath::Vector3<Scalar> dimensions_;
  // Pose of the center of the box in the reference frame.
  Pose3<Scalar> ref_t_box_;
};

template <typename Scalar>
std::ostream& operator<<(std::ostream& os,
                         const OrientedBoundingBox3<Scalar>& bbox) {
  return os << "OrientedBoundingBox3" << typeid(Scalar).name()
            << "(dimensions: " << bbox.Dimensions().transpose()
            << ", ref_t_box: " << bbox.RefTBox() << ")";
}

using OrientedBoundingBox3d = OrientedBoundingBox3<double>;
using OrientedBoundingBox3f = OrientedBoundingBox3<float>;

absl::StatusOr<intrinsic_proto::geometry::OrientedBoundingBox3> ToProto(
    const OrientedBoundingBox3d& bbox);
absl::StatusOr<OrientedBoundingBox3d> FromProto(
    const intrinsic_proto::geometry::OrientedBoundingBox3& proto);
}  // namespace intrinsic::geo

REGISTER_RIEGELI_PROTO_CODER_EXPLICIT(
    intrinsic::geo::OrientedBoundingBox3d,
    intrinsic_proto::geometry::OrientedBoundingBox3, intrinsic::geo::ToProto,
    intrinsic::geo::FromProto);

namespace intrinsic {
using ::intrinsic::geo::FromProto;
using ::intrinsic::geo::OrientedBoundingBox3f;
using ::intrinsic::geo::ToProto;
}  // namespace intrinsic
#endif  // INTRINSIC_GEOMETRY_API_ORIENTED_BOUNDING_BOX_H_
