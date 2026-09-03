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

#include "intrinsic/geometry/api/shape_factory.h"

#include <optional>

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/affine_transform_of_geometry.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/shapes/box.h"
#include "intrinsic/geometry/shapes/capsule.h"
#include "intrinsic/geometry/shapes/cylinder.h"
#include "intrinsic/geometry/shapes/sphere.h"
#include "intrinsic/math/pose3.h"

namespace intrinsic::geo {
Geometry MakeSphere(double radius) {
  return Geometry(ExactGeometry(Sphere(radius)),
                  /*provenance=*/std::nullopt);
}

TransformedGeometry MakeTransformedSphere(double radius,
                                          const Pose3d& ref_t_shape) {
  return TransformedGeometry{MakeSphere(radius), ref_t_shape};
}

TransformedGeometry MakeTransformedSphere(double radius) {
  return MakeTransformedSphere(radius, Pose3d());
}

Geometry MakeCylinder(double length, double radius) {
  return Geometry(ExactGeometry(Cylinder(length, radius)),
                  /*provenance=*/std::nullopt);
}

TransformedGeometry MakeTransformedCylinder(double length, double radius,
                                            const Pose3d& ref_t_shape) {
  return TransformedGeometry{MakeCylinder(length, radius), ref_t_shape};
}

TransformedGeometry MakeTransformedCylinder(double length, double radius) {
  return MakeTransformedCylinder(length, radius, Pose3d());
}

Geometry MakeCenteredBox(double x, double y, double z) {
  return Geometry(ExactGeometry(Box(eigenmath::Vector3d(x, y, z))),
                  /*provenance=*/std::nullopt);
}

TransformedGeometry MakeTransformedCenteredBox(double x, double y, double z,
                                               const Pose3d& ref_t_shape) {
  return TransformedGeometry{MakeCenteredBox(x, y, z), ref_t_shape};
}

TransformedGeometry MakeTransformedCenteredBox(double x, double y, double z) {
  return MakeTransformedCenteredBox(x, y, z, Pose3d());
}

Geometry MakeCapsule(double length, double radius) {
  return Geometry(ExactGeometry(Capsule(length, radius)),
                  /*provenance=*/std::nullopt);
}

TransformedGeometry MakeTransformedCapsule(double length, double radius,
                                           const Pose3d& ref_t_shape) {
  return TransformedGeometry{MakeCapsule(length, radius), ref_t_shape};
}

TransformedGeometry MakeTransformedCapsule(double length, double radius) {
  return MakeTransformedCapsule(length, radius, Pose3d());
}

}  // namespace intrinsic::geo
