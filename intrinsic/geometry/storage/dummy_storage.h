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

#ifndef INTRINSIC_GEOMETRY_STORAGE_DUMMY_STORAGE_H_
#define INTRINSIC_GEOMETRY_STORAGE_DUMMY_STORAGE_H_

#include <memory>

#include "intrinsic/geometry/storage/geometry_library.h"

namespace intrinsic::geo {

// Returns a geometry deserializer instance that always returns an empty Mesh.
// Returns a GeometryLibrary that doesn't do anything:
//
// * Its Deserializer() always returns an empty Mesh.
// * Its Serializer() does not actually save anything.
std::unique_ptr<GeometryLibrary> GetDummyGeometryLibrary();

}  // namespace intrinsic::geo

namespace intrinsic {
using ::intrinsic::geo::GetDummyGeometryLibrary;
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_STORAGE_DUMMY_STORAGE_H_
