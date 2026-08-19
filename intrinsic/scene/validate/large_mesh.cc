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

#include "intrinsic/scene/validate/large_mesh.h"

#include <cmath>
#include <string>

#include "absl/container/flat_hash_map.h"
#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/substitute.h"
#include "intrinsic/geometry/api/axis_aligned_bounding_box_3d.h"
#include "intrinsic/geometry/api/compute_axis_aligned_bounding_box_3d.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::scene_object {

absl::Status CheckForLargeMeshes(
    const WorldHashMap<std::string, Geometry>& geometry, double max_aabb_diag) {
  const double max_mesh_diagonal_sq = max_aabb_diag * max_aabb_diag;
  absl::flat_hash_map<std::string, double> geo_id_to_large_diagonals;

  for (const auto& [geo_id, geo] : geometry) {
    INTR_ASSIGN_OR_RETURN(AxisAlignedBoundingBox3d mesh_bounding_box,
                          ComputeAxisAlignedBoundingBox3d(geo));
    const eigenmath::Vector3d diag = mesh_bounding_box.GetDiagonal();
    if (diag.squaredNorm() > max_mesh_diagonal_sq) {
      geo_id_to_large_diagonals.emplace(geo_id, std::sqrt(diag.squaredNorm()));
    }
  }

  if (!geo_id_to_large_diagonals.empty()) {
    std::string err_msg = absl::Substitute(
        "Meshes were found that have a diagonal larger than allowed size $0. "
        "Please double check that this mesh is provided in meter scale (as "
        "opposed to millimeters). If this is intentional, increase the allowed "
        "maximum diagonal or disable the check. The objects that have meshes "
        "that violate this diagonal are:",
        max_aabb_diag);
    for (const auto& [geo_id, diag] : geo_id_to_large_diagonals) {
      absl::StrAppend(
          &err_msg, absl::Substitute("\ngeo_id: '$0' diag: $1", geo_id, diag));
    }
    return absl::InvalidArgumentError(err_msg);
  }

  return absl::OkStatus();
}

}  // namespace intrinsic::scene_object
