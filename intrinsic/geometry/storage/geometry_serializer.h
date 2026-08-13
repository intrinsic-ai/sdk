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

#ifndef INTRINSIC_GEOMETRY_STORAGE_GEOMETRY_SERIALIZER_H_
#define INTRINSIC_GEOMETRY_STORAGE_GEOMETRY_SERIALIZER_H_

#include "absl/status/statusor.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/proto/v1/geometry_storage_refs.pb.h"

namespace intrinsic {

// The GeometrySerializer is responsible for writing geometry to some storage
// backend for a given geometry.
// The exact details of where the data is stored has been abstracted. It is up
// to the implementation to decide how much deduplication to support.
class GeometrySerializer {
 public:
  virtual ~GeometrySerializer() = default;

  // Stores the given geometry object and returns its storage refs. The
  // storage refs can be used with a GeometryDeserializer configured to the same
  // store to fetch the geometry later.
  virtual absl::StatusOr<intrinsic_proto::geometry::v1::GeometryStorageRefs>
  SaveGeometryV1(const Geometry& geometry) = 0;
};

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_STORAGE_GEOMETRY_SERIALIZER_H_
