// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_STORAGE_GEOMETRY_DESERIALIZER_H_
#define INTRINSIC_GEOMETRY_STORAGE_GEOMETRY_DESERIALIZER_H_

#include <optional>

#include "absl/status/statusor.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/geometry_options.h"
#include "intrinsic/geometry/proto/geometry_storage_refs.pb.h"
#include "intrinsic/geometry/proto/v1/geometry_storage_refs.pb.h"

namespace intrinsic {

// The GeometryDeserializer is responsible for reading geometry data from some
// source based on the given id. The exact details of where the data is stored
// has been abstracted but once the data has been fetched it is placed into a
// Geometry instance.
class GeometryDeserializer {
 public:
  virtual ~GeometryDeserializer() = default;

  // DEPRECATED: `GetGeometry` with
  // `intrinsic_proto::geometry::GeometryStorageRefs` exists primarily for
  // backwards compatibility. Use v1::GeometryStorageRefs version below instead.
  virtual absl::StatusOr<Geometry> GetGeometry(
      const intrinsic_proto::geometry::GeometryStorageRefs& geo_storage_refs,
      const GeometryOptions& options) const = 0;

  // Retrieves a shape and renderable for the given `geo_storage_refs`. One way
  // to get a valid `GeometryStorageRef` is to save the return value of
  // `GeometrySerializer::SaveGeometry()` when saving a piece of geometry.
  virtual absl::StatusOr<Geometry> GetGeometry(
      const intrinsic_proto::geometry::v1::GeometryStorageRefs&
          geo_storage_refs,
      std::optional<intrinsic_proto::geometry::v1::MaterialProperties>
          material_properties) const = 0;
};

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_STORAGE_GEOMETRY_DESERIALIZER_H_
