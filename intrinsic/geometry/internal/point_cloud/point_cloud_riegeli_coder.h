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

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_POINT_CLOUD_RIEGELI_CODER_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_POINT_CLOUD_RIEGELI_CODER_H_

#include "absl/status/statusor.h"
#include "intrinsic/geometry/proto/v1/point_cloud.pb.h"
#include "intrinsic/geometry/shapes/point_cloud.h"
#include "intrinsic/marshal/riegeli_proto_coder.h"

namespace intrinsic::geo {

absl::StatusOr<intrinsic_proto::geometry::v1::PointCloud> ToProto(
    const PointCloud& point_cloud);

// This can't be `FromProto` since it conflicts with existing `FromProto`
// function that returns a different type.
absl::StatusOr<PointCloud> ToShape(
    const intrinsic_proto::geometry::v1::PointCloud& proto);

}  // namespace intrinsic::geo

namespace intrinsic {
REGISTER_RIEGELI_PROTO_CODER_EXPLICIT(geo::PointCloud,
                                      intrinsic_proto::geometry::v1::PointCloud,
                                      geo::ToProto, geo::ToShape);

using ::intrinsic::geo::ToProto;
using ::intrinsic::geo::ToShape;
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_POINT_CLOUD_RIEGELI_CODER_H_
