// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_STORAGE_GEOMETRY_SERIALIZER_H_
#define INTRINSIC_GEOMETRY_STORAGE_GEOMETRY_SERIALIZER_H_

#include "absl/status/statusor.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/proto/geometry_storage_refs.pb.h"
#include "intrinsic/geometry/proto/v1/geometry_storage_refs.pb.h"

namespace intrinsic {

// The GeometrySerializer is responsible for writing geometry to some storage
// backend for a given geometry.
// The exact details of where the data is stored has been abstracted. It is up
// to the implementation to decide how much deduplication to support.
class GeometrySerializer {
 public:
  virtual ~GeometrySerializer() = default;

  // DEPRECATED: this interface exists only for backwards compatibility. Use
  // `v1::GeometryStorageRefs` API below instead.
  virtual absl::StatusOr<intrinsic_proto::geometry::GeometryStorageRefs>
  SaveGeometry(const Geometry& geometry) = 0;

  // Stores the given geometry object and returns its storage refs. The
  // storage refs can be used with a GeometryDeserializer configured to the same
  // store to fetch the geometry later.
  virtual absl::StatusOr<intrinsic_proto::geometry::v1::GeometryStorageRefs>
  SaveGeometryV1(const Geometry& geometry) = 0;
};

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_STORAGE_GEOMETRY_SERIALIZER_H_
