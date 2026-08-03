// Copyright 2023 Intrinsic Innovation LLC

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
