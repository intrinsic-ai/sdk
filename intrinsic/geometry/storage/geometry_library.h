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
