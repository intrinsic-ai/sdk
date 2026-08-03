// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_TRIANGLE_MESH_RIEGELI_CODER_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_TRIANGLE_MESH_RIEGELI_CODER_H_

#include "absl/status/statusor.h"
#include "intrinsic/geometry/proto/v1/triangle_mesh.pb.h"
#include "intrinsic/geometry/shapes/triangle_mesh.h"
#include "intrinsic/marshal/riegeli_proto_coder.h"

namespace intrinsic {

absl::StatusOr<intrinsic_proto::geometry::v1::TriangleMesh> ToProto(
    const shapes::TriangleMesh& triangle_mesh);

// This can't be `FromProto` since it conflicts with existing `FromProto`
// function that returns a different type.
absl::StatusOr<shapes::TriangleMesh> ToShape(
    const intrinsic_proto::geometry::v1::TriangleMesh& proto);

REGISTER_RIEGELI_PROTO_CODER_EXPLICIT(
    shapes::TriangleMesh, intrinsic_proto::geometry::v1::TriangleMesh, ToProto,
    ToShape);

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_TRIANGLE_MESH_RIEGELI_CODER_H_
