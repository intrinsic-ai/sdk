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

#ifndef INTRINSIC_GEOMETRY_API_VALIDATE_MESH_H_
#define INTRINSIC_GEOMETRY_API_VALIDATE_MESH_H_

#include <string>

#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/eigenmath/types.h"

namespace intrinsic::geo {

// Validates a mesh file (STL, GLB, DAE, etc.) to ensure it is valid and all its
// dependencies (e.g. textures) are present.
absl::Status ValidateMeshFile(absl::string_view filename,
                              const eigenmath::Vector3d& scale);

// Validates mesh data provided as a buffer.
// `extension` should be the file extension (e.g. "stl", "glb", "dae") to help
// identifying the format.
absl::Status ValidateMeshData(absl::string_view mesh_data,
                              absl::string_view extension,
                              const eigenmath::Vector3d& scale);

}  // namespace intrinsic::geo
namespace intrinsic {
using ::intrinsic::geo::ValidateMeshData;
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_VALIDATE_MESH_H_
