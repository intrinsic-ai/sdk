// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_STORAGE_DUMMY_STORAGE_H_
#define INTRINSIC_GEOMETRY_STORAGE_DUMMY_STORAGE_H_

#include <memory>

#include "intrinsic/geometry/storage/geometry_library.h"

namespace intrinsic {

// Returns a geometry deserializer instance that always returns an empty Mesh.
// Returns a GeometryLibrary that doesn't do anything:
//
// * Its Deserializer() always returns an empty Mesh.
// * Its Serializer() does not actually save anything.
std::unique_ptr<GeometryLibrary> GetDummyGeometryLibrary();

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_STORAGE_DUMMY_STORAGE_H_
