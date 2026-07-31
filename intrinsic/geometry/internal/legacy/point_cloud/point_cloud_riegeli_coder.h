// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_POINT_CLOUD_RIEGELI_CODER_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_POINT_CLOUD_RIEGELI_CODER_H_

#include "absl/status/statusor.h"
#include "intrinsic/geometry/proto/v1/point_cloud.pb.h"
#include "intrinsic/geometry/shapes/point_cloud.h"
#include "intrinsic/marshal/riegeli_proto_coder.h"

namespace intrinsic {

absl::StatusOr<intrinsic_proto::geometry::v1::PointCloud> ToProto(
    const shapes::PointCloud& point_cloud);

// This can't be `FromProto` since it conflicts with existing `FromProto`
// function that returns a different type.
absl::StatusOr<shapes::PointCloud> ToShape(
    const intrinsic_proto::geometry::v1::PointCloud& proto);

REGISTER_RIEGELI_PROTO_CODER_EXPLICIT(shapes::PointCloud,
                                      intrinsic_proto::geometry::v1::PointCloud,
                                      ToProto, ToShape);

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_POINT_CLOUD_POINT_CLOUD_RIEGELI_CODER_H_
