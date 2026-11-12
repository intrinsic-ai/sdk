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

#include "intrinsic/geometry/api/oriented_bounding_box.h"

#include "absl/status/statusor.h"
#include "intrinsic/geometry/proto/oriented_bounding_box.pb.h"
#include "intrinsic/math/proto_conversion.h"
#include "intrinsic/util/status/ret_check.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::geo {
absl::StatusOr<intrinsic_proto::geometry::OrientedBoundingBox3> ToProto(
    const OrientedBoundingBox3d& bbox) {
  intrinsic_proto::geometry::OrientedBoundingBox3 proto;
  *proto.mutable_dimensions() = ToVectorProto(bbox.Dimensions());
  *proto.mutable_ref_t_box() = ToProto(bbox.RefTBox());
  return proto;
}

absl::StatusOr<OrientedBoundingBox3d> FromProto(
    const intrinsic_proto::geometry::OrientedBoundingBox3& proto) {
  const auto dimensions = FromProto(proto.dimensions());
  INTR_RET_CHECK(dimensions.array().isFinite().all())
      << "Dimensions are not finite.";
  INTR_ASSIGN_OR_RETURN(const auto ref_t_box, FromProto(proto.ref_t_box()));
  return OrientedBoundingBox3d(dimensions, ref_t_box);
}

}  // namespace intrinsic::geo
