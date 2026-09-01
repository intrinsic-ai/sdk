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

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_MESH_PRIMITIVES_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_MESH_PRIMITIVES_H_

#include <cstdint>

#include "Eigen/Core"
#include "absl/status/statusor.h"
#include "intrinsic/geometry/internal/mesh/mesh.h"

namespace intrinsic::geo {

// Creates a capsule of a total length 'length'.
// Both the cylinder and the hemispherical caps have radius 'radius'.
// 'resolution' is the number of vertices per ring and also increases the
// number of rings.
// The capsule axis is z-aligned.
Mesh CreateCapsule(int resolution, double radius, double height);

// The capped cylinder axis is z-aligned.
Mesh CreateCylinder(int N, double radius, double length);

// recursion_level is the number of times that the triangles are subdivided for
// the geometry to be created: https://en.wikipedia.org/wiki/Icosphere
Mesh CreateSphere(int recursion_level, double radius);

Mesh CreateEllipsoid(int recursion_level, double radius_x, double radius_y,
                     double radius_z);

// Create a cuboid centered at the origin. The dimension is specified by
// positive half edge.
Mesh CreateCuboid(const Eigen::Vector3d& half_edge);

// Creates a cuboid using the given half_edge for size with a select set of
// sides.
// The sides param should contain a bitmask for which sides are activated.  The
// bits correspond to the following faces:
//
// 0 -X
// 1 +X
// 2 -Y
// 3 +Y
// 4 -Z
// 5 +Z
//
// Setting bits 6 or 7 will result in an assert.
Mesh CreatePartialCuboid(const Eigen::Vector3d& half_edge, uint8_t sides);

// Creates a frustum using the given angles and near far distances.
// The frustum extend from the origin towards the +z direction with its tip at
// the origin.
// x_angle represents the angle in radians between the x-z plane and the frustum
// faces rotating around the x-axis.
// y_angle represents the angle in radians between the y-z plane and the frustum
// faces rotating around the y-axis.
// min_z_distance is the distance between the origin and the near frustum face
// cutting the tip of the pyramid.
// max_z_distance is the distance between the origin and the far frustum face
// representing the base of the pyramid.
absl::StatusOr<Mesh> CreateFrustumMesh(double x_angle, double y_angle,
                                       double min_z_distance,
                                       double max_z_distance);

}  // namespace intrinsic::geo
namespace intrinsic {
using ::intrinsic::geo::CreateSphere;
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_MESH_PRIMITIVES_H_
