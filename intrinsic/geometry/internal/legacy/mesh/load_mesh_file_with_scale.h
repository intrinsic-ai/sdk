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

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_LOAD_MESH_FILE_WITH_SCALE_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_LOAD_MESH_FILE_WITH_SCALE_H_

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"

namespace intrinsic {
namespace geometry_legacy {

// Load a mesh object from a file, scale according to the given vector.
absl::StatusOr<Mesh> LoadMeshFileWithScale(absl::string_view path,
                                           const eigenmath::Vector3d& scale);

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_LOAD_MESH_FILE_WITH_SCALE_H_
