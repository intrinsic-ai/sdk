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

#ifndef INTRINSIC_SCENE_VALIDATE_LARGE_MESH_H_
#define INTRINSIC_SCENE_VALIDATE_LARGE_MESH_H_

#include <string>

#include "absl/status/status.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/world/hashing/hashing.h"

namespace intrinsic::scene_object {

// Returns an error if the length of the bounding box diagonal for any of the
// mesh geometries exceeds `max_aabb_diag`.
absl::Status CheckForLargeMeshes(
    const WorldHashMap<std::string, Geometry>& geometry, double max_aabb_diag);

}  // namespace intrinsic::scene_object

#endif  // INTRINSIC_SCENE_VALIDATE_LARGE_MESH_H_
