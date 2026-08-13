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

#include "intrinsic/geometry/internal/legacy/mesh/load_mesh_file_with_scale.h"

#include <string>

#include "absl/log/check.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/internal/legacy/mesh/io/load_mesh_from_file.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/scene/sdf/sdf_path_resolver.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {
namespace geometry_legacy {

absl::StatusOr<Mesh> LoadMeshFileWithScale(absl::string_view path,
                                           const eigenmath::Vector3d& scale) {
  INTR_ASSIGN_OR_RETURN(auto full_path,
                        sdf::SdfPathResolver(std::string(path)));
  return LoadMeshFromFile(full_path, scale);
}

}  // namespace geometry_legacy
}  // namespace intrinsic
