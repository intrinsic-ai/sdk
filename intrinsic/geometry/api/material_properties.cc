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

#include "intrinsic/geometry/api/material_properties.h"

namespace intrinsic {

intrinsic_proto::geometry::v1::MaterialProperties MaterialProperties::ToProto()
    const {
  intrinsic_proto::geometry::v1::MaterialProperties proto;
  auto* color_proto = proto.mutable_base_color();
  color_proto->set_red(color.x());
  color_proto->set_green(color.y());
  color_proto->set_blue(color.z());
  proto.set_metalness(metalness);
  proto.set_roughness(roughness);
  proto.set_transmission(transmission);
  return proto;
}

MaterialProperties MaterialProperties::FromProto(
    const intrinsic_proto::geometry::v1::MaterialProperties& proto) {
  MaterialProperties properties;
  if (proto.has_base_color()) {
    properties.color = eigenmath::Vector3d{proto.base_color().red(),
                                           proto.base_color().green(),
                                           proto.base_color().blue()};
  }
  if (proto.has_metalness()) {
    properties.metalness = proto.metalness();
  }
  if (proto.has_roughness()) {
    properties.roughness = proto.roughness();
  }
  if (proto.has_transmission()) {
    properties.transmission = proto.transmission();
  }
  return properties;
}

bool operator==(const MaterialProperties& a, const MaterialProperties& b) {
  return a.color == b.color && a.metalness == b.metalness &&
         a.roughness == b.roughness && a.transmission == b.transmission;
}

}  // namespace intrinsic
