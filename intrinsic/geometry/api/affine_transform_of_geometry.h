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

#ifndef INTRINSIC_GEOMETRY_API_AFFINE_TRANSFORM_OF_GEOMETRY_H_
#define INTRINSIC_GEOMETRY_API_AFFINE_TRANSFORM_OF_GEOMETRY_H_

#include "intrinsic/geometry/api/affine_transform_of.h"
#include "intrinsic/geometry/api/geometry.h"

namespace intrinsic::geo {

// A AffineTransformedGeometry is essentially a pair of Geometry and a affine
// transformation. This allows explicit separation of the Geometry while
// keeping different or changing transforms for it. For example two spheres with
// different world coordinates can share their internal Geometry while having
// different transforms.
using TransformedGeometry = AffineTransformOf<Geometry>;

}  // namespace intrinsic::geo

namespace intrinsic {
using ::intrinsic::geo::TransformedGeometry;
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_AFFINE_TRANSFORM_OF_GEOMETRY_H_
