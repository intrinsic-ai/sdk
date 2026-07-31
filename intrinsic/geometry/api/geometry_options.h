// Copyright 2023 Intrinsic Innovation LLC

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
        .fill_inside_for_distance_queries = true,
        .simulation_convex_decomposition_resolution = std::nullopt,
    };
    return options;
  }

  void MergeWith(const GeometryOptions& other) {
    // If either option requires filling on the inside, then we need to fill on
    // the inside.
    fill_inside_for_distance_queries |= other.fill_inside_for_distance_queries;

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
    return fill_inside_for_distance_queries ==
               other.fill_inside_for_distance_queries &&
           simulation_convex_decomposition_resolution ==
               other.simulation_convex_decomposition_resolution;
  }

  // Used to determine if the octree needs to be filled on the inside when
  // constructing the octree for distance queries.
  bool fill_inside_for_distance_queries = true;

  // Optional parameter to specify the voxel resolution for convex
  // decomposition. This parameter only affects how the geometry is simulated
  // in Gazebo and is typically used for optimizing collision meshes.
  // If the associated geometry is not a mesh type, this field has no effect.
  std::optional<uint32_t> simulation_convex_decomposition_resolution =
      std::nullopt;
};

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_GEOMETRY_OPTIONS_H_
