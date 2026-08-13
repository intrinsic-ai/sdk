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

#ifndef INTRINSIC_GEOMETRY_API_GEOMETRY_OPTIONS_H_
#define INTRINSIC_GEOMETRY_API_GEOMETRY_OPTIONS_H_

#include <algorithm>
#include <cstdint>
#include <optional>

namespace intrinsic {

// A set of options that can be used when processing a Geometry instance.
struct GeometryOptions {
  static const GeometryOptions& Default() {
    static GeometryOptions options = {
        .simulation_convex_decomposition_resolution = std::nullopt,
    };
    return options;
  }

  void MergeWith(const GeometryOptions& other) {
    // If convex decomposition resolution is set for either option, then use
    // that in the merged option. If the parameter is set for both options, then
    // use the higher value.
    if (!simulation_convex_decomposition_resolution.has_value()) {
      simulation_convex_decomposition_resolution =
          other.simulation_convex_decomposition_resolution;
    } else if (other.simulation_convex_decomposition_resolution.has_value()) {
      simulation_convex_decomposition_resolution =
          std::max(*other.simulation_convex_decomposition_resolution,
                   *simulation_convex_decomposition_resolution);
    }
  }

  bool operator==(const GeometryOptions& other) const {
    return simulation_convex_decomposition_resolution ==
           other.simulation_convex_decomposition_resolution;
  }

  // Optional parameter to specify the voxel resolution for convex
  // decomposition. This parameter only affects how the geometry is simulated
  // in Gazebo and is typically used for optimizing collision meshes.
  // If the associated geometry is not a mesh type, this field has no effect.
  std::optional<uint32_t> simulation_convex_decomposition_resolution =
      std::nullopt;
};

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_GEOMETRY_OPTIONS_H_
