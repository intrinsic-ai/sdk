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

#ifndef INTRINSIC_GEOMETRY_STORAGE_GZF_STORAGE_H_
#define INTRINSIC_GEOMETRY_STORAGE_GZF_STORAGE_H_

#include <memory>

#include "absl/base/attributes.h"
#include "intrinsic/geometry/storage/geometry_library.h"
#include "intrinsic/world/gzfile/gzfile.h"

namespace intrinsic::geo {

// Returns a GeometryLibrary instance that is backed by `gzfile`.
//
// DOES NOT take ownership of `gzfile`. `gzfile` MUST outlive the
// GeometryLibrary.
std::unique_ptr<GeometryLibrary> GetGzfGeometryLibrary(
    GZFile& gzfile ABSL_ATTRIBUTE_LIFETIME_BOUND);

// Returns a read-only GeometryLibrary instance that is backed by `gzfile`.
//
// DOES NOT take ownership of `gzfile`. `gzfile` MUST outlive the
// GeometryLibrary.
//
// The Library's Serializer() returns UnimplementedError.
std::unique_ptr<GeometryLibrary> GetReadOnlyGzfGeometryLibrary(
    const GZFile& gzfile ABSL_ATTRIBUTE_LIFETIME_BOUND);

}  // namespace intrinsic::geo

namespace intrinsic {
using ::intrinsic::geo::GetGzfGeometryLibrary;
using ::intrinsic::geo::GetReadOnlyGzfGeometryLibrary;
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_STORAGE_GZF_STORAGE_H_
