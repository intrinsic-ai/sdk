// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_STORAGE_GEOMETRY_LIBRARY_H_
#define INTRINSIC_GEOMETRY_STORAGE_GEOMETRY_LIBRARY_H_

#include "intrinsic/geometry/storage/geometry_deserializer.h"
#include "intrinsic/geometry/storage/geometry_serializer.h"

namespace intrinsic {

// Provides read and write access to a single form of geometry storage.
//
// It should be uncommon for code to have *values of*, rather than *references
// to*, a Geometry(De)serializer directly. Use a GeometryLibrary for storage, to
// ensure that the Serializer and Deserializer match (i.e. actually use the same
// underlying storage).
//
// Implementations of GeometryLibrary can make optimizations across both
// serialization and deserialization. For example, a GeometryLibrary can contain
// a cache where reasonable, and write to that cache when its GeometrySerializer
// saves a piece of geometry, so that its GeometryDeserializer does not have to
// hit the underlying data store.
class GeometryLibrary {
 public:
  virtual ~GeometryLibrary() = default;

  virtual const GeometryDeserializer& Deserializer() const = 0;
  virtual GeometrySerializer& Serializer() = 0;
};

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_STORAGE_GEOMETRY_LIBRARY_H_
