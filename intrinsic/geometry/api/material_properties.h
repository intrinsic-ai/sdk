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

#ifndef INTRINSIC_GEOMETRY_API_MATERIAL_PROPERTIES_H_
#define INTRINSIC_GEOMETRY_API_MATERIAL_PROPERTIES_H_

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/proto/v1/material.pb.h"

namespace intrinsic {

// Simple PBR MaterialProperties struct, meant to make working with
// the v1 MaterialProperties proto easier. See the MaterialProperties proto
// message for detailed explanations of the fields.
struct MaterialProperties {
  // The surface's RGB base color, where each component is [0, 1]
  eigenmath::Vector3d color{0.5, 0.5, 0.5};
  // The surface metalness, [0, 1]
  double metalness{1.0};
  // The surface roughness, [0, 1]
  double roughness{1.0};
  // The surface transmission [0, 1]
  double transmission{0.0};

  intrinsic_proto::geometry::v1::MaterialProperties ToProto() const;
  static MaterialProperties FromProto(
      const intrinsic_proto::geometry::v1::MaterialProperties& proto);
};

bool operator==(const MaterialProperties& a, const MaterialProperties& b);

template <typename H>
H AbslHashValue(H h, const MaterialProperties& m) {
  return H::combine(std::move(h), m.color.x(), m.color.y(), m.color.z(),
                    m.metalness, m.roughness, m.transmission);
}

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_MATERIAL_PROPERTIES_H_
