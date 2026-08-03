// Copyright 2023 Intrinsic Innovation LLC

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
