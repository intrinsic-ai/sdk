// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_API_IO_H_
#define INTRINSIC_GEOMETRY_API_IO_H_

#include <optional>
#include <vector>

#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/affine_transform_of_geometry.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/renderable.h"
#include "intrinsic/geometry/proto/v1/exact_geometry.pb.h"
#include "intrinsic/geometry/proto/v1/geometric_transform.pb.h"
#include "intrinsic/geometry/proto/v1/geometry.pb.h"
#include "intrinsic/geometry/proto/v1/geometry_storage_refs.pb.h"
#include "intrinsic/geometry/proto/v1/inline_geometry.pb.h"
#include "intrinsic/geometry/proto/v1/primitive_shape.pb.h"
#include "intrinsic/geometry/proto/v1/primitives.pb.h"
#include "intrinsic/geometry/proto/v1/renderable.pb.h"
#include "intrinsic/geometry/proto/v1/transformed_geometry.pb.h"
#include "intrinsic/geometry/proto/v1/transformed_primitive_shape.pb.h"
#include "intrinsic/geometry/proto/v1/transformed_primitive_shape_set.pb.h"
#include "intrinsic/geometry/storage/geometry_deserializer.h"
#include "intrinsic/geometry/storage/geometry_serializer.h"

namespace intrinsic {

// Geometry proto conversions.
absl::StatusOr<intrinsic_proto::geometry::v1::Geometry> ToProto(
    const Geometry& geo, GeometrySerializer* serializer = nullptr);
absl::StatusOr<intrinsic_proto::geometry::v1::Geometry> ToInlinedProto(
    const Geometry& geo);

absl::StatusOr<Geometry> ToGeometry(
    const intrinsic_proto::geometry::v1::Geometry& proto,
    const GeometryDeserializer* deserializer = nullptr);
absl::StatusOr<Geometry> ToGeometry(
    const intrinsic_proto::geometry::v1::InlineGeometry& proto);
absl::StatusOr<Geometry> ToGeometry(
    const intrinsic_proto::geometry::v1::InlineGeometry& proto,
    std::optional<intrinsic_proto::geometry::v1::MaterialProperties>
        material_properties);

// TransformedGeometry proto conversions.
absl::StatusOr<intrinsic_proto::geometry::v1::TransformedGeometry> ToProto(
    const TransformedGeometry& geo, GeometrySerializer* serializer = nullptr);
absl::StatusOr<TransformedGeometry> ToGeometry(
    const intrinsic_proto::geometry::v1::TransformedGeometry& proto,
    const GeometryDeserializer* deserializer);

// Renderable proto conversions.
absl::StatusOr<intrinsic_proto::geometry::v1::Renderable> ToProto(
    const Renderable& geo);
absl::StatusOr<Renderable> ToRenderable(
    const intrinsic_proto::geometry::v1::Renderable& proto);

// GeometricTransform proto conversions.
absl::StatusOr<intrinsic_proto::geometry::v1::GeometricTransform>
ToGeometricTransform(const eigenmath::Matrix4d& matrix);
absl::StatusOr<eigenmath::Matrix4d> ToAffineTransform(
    const intrinsic_proto::geometry::v1::GeometricTransform& proto);

// ExactGeometry proto conversions.
absl::StatusOr<intrinsic_proto::geometry::v1::ExactGeometry> ToProto(
    const ExactGeometry& geo);
absl::StatusOr<ExactGeometry> ToGeometry(
    const intrinsic_proto::geometry::v1::ExactGeometry& proto);

namespace geometry_details {

// PrimitiveShape proto conversions.
absl::StatusOr<intrinsic_proto::geometry::v1::PrimitiveShape> ToProto(
    const PrimitiveShapePtr& shape);
absl::StatusOr<PrimitiveShapePtr> ToPrimitive(
    const intrinsic_proto::geometry::v1::PrimitiveShape& proto);

// TransformedPrimitiveShape proto conversions.
absl::StatusOr<intrinsic_proto::geometry::v1::TransformedPrimitiveShape>
ToProto(const TransformedPrimitiveShapePtr& shape);
absl::StatusOr<TransformedPrimitiveShapePtr> ToPrimitive(
    const intrinsic_proto::geometry::v1::TransformedPrimitiveShape& proto);

// TransformedPrimitiveShapeSet proto conversions.
absl::StatusOr<intrinsic_proto::geometry::v1::TransformedPrimitiveShapeSet>
ToProto(const std::vector<TransformedPrimitiveShapePtr>& shape);
absl::StatusOr<std::vector<TransformedPrimitiveShapePtr>> ToPrimitiveSet(
    const intrinsic_proto::geometry::v1::TransformedPrimitiveShapeSet& proto);

// Converts Geometry into an InlineGeometry proto.
absl::StatusOr<intrinsic_proto::geometry::v1::InlineGeometry>
ToInlineGeometryProto(const Geometry& geo);

absl::StatusOr<Geometry::Provenance> ToGeometryProvenance(
    const intrinsic_proto::geometry::v1::GeometryProvenance& proto,
    const GeometryDeserializer* deserializer = nullptr);
absl::StatusOr<intrinsic_proto::geometry::v1::GeometryProvenance> ToProto(
    const Geometry::Provenance& geo_provenance,
    GeometrySerializer* serializer = nullptr);

// The maximum size of the geometry proto that we will inline into the
// provenance proto. If the size is larger than this, we will store the proto in
// CAS and instead store only the uri into the provenance proto.
inline constexpr int kMaxGeometryProtoSize = 100 * 1024 * 1024;  // 100 MB;

}  // namespace geometry_details

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_IO_H_
