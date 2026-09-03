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

#ifndef INTRINSIC_GEOMETRY_API_SHAPE_FACTORY_H_
#define INTRINSIC_GEOMETRY_API_SHAPE_FACTORY_H_

#include "intrinsic/geometry/api/affine_transform_of_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/math/pose3.h"

namespace intrinsic::geo {
// Generates a sphere with the given `radius` that is centered at origin.
Geometry MakeSphere(double radius);
// Generates a sphere with the given `radius` that is centered at origin.
TransformedGeometry MakeTransformedSphere(double radius);

// Generates a sphere with the given `radius`that is centered at origin and then
// transformed by the given `ref_t_shape`.
TransformedGeometry MakeTransformedSphere(double radius,
                                          const Pose3d& ref_t_shape);

// Generates a cylinder with the given `radius` and `length` that is centered at
// origin.
Geometry MakeCylinder(double length, double radius);
// Generates a cylinder with the given `radius` and `length` that is centered at
// origin.
TransformedGeometry MakeTransformedCylinder(double length, double radius);

// Generates a cylinder with the given `radius` and `length` that is centered at
// origin and then transformed by the given `ref_t_shape`.
TransformedGeometry MakeTransformedCylinder(double length, double radius,
                                            const Pose3d& ref_t_shape);

// Generates a box with the given edge lengths in `x`, `y`, and `z` that is
// centered at origin.
Geometry MakeCenteredBox(double x, double y, double z);
// Generates a box with the given edge lengths in `x`, `y`, and `z` that is
// centered at origin.
TransformedGeometry MakeTransformedCenteredBox(double x, double y, double z);

// Generates a box with the given edge lengths in `x`, `y`, and `z` that is
// centered at origin and then transformed by the given `ref_t_shape`.
TransformedGeometry MakeTransformedCenteredBox(double x, double y, double z,
                                               const Pose3d& ref_t_shape);

// Generates a Capsule with the given `radius` and `length` that is centered at
// origin.
Geometry MakeCapsule(double length, double radius);
// Generates a Capsule with the given `radius` and `length` that is centered at
// origin.
TransformedGeometry MakeTransformedCapsule(double length, double radius);

// Generates a Capsule with the given `radius` and `length` that is centered at
// origin and then transformed by the given `ref_t_shape`.
TransformedGeometry MakeTransformedCapsule(double length, double radius,
                                           const Pose3d& ref_t_shape);
}  // namespace intrinsic::geo

#endif  // INTRINSIC_GEOMETRY_API_SHAPE_FACTORY_H_
