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

#include "intrinsic/scene/validate/scene_object_validate_geo.h"

#include <string>

#include "absl/status/status.h"
#include "absl/strings/substitute.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/renderable.h"
#include "intrinsic/geometry/api/validate_mesh.h"
#include "intrinsic/geometry/proto/v1/geometric_transform.pb.h"
#include "intrinsic/geometry/proto/v1/geometry.pb.h"
#include "intrinsic/geometry/proto/v1/transformed_geometry.pb.h"
#include "intrinsic/math/proto_conversion.h"
#include "intrinsic/scene/proto/v1/entity.pb.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/proto/geometry_component.pb.h"

namespace intrinsic {
namespace scene_object {

namespace {

using DataCase = intrinsic_proto::geometry::v1::Geometry::DataCase;
using intrinsic_proto::scene_object::v1::Entity;
using intrinsic_proto::scene_object::v1::SceneObject;

eigenmath::Vector3d GetScale(
    const intrinsic_proto::geometry::v1::GeometricTransform& transform) {
  if (transform.has_trs() && transform.trs().has_scale()) {
    return FromProto(transform.trs().scale());
  }
  return eigenmath::Vector3d::Ones();
}

absl::Status ValidateTransformedGeometry(
    const intrinsic_proto::geometry::v1::TransformedGeometry& transformed_geo,
    const GeometryDeserializer& geolib) {
  const auto& geo = transformed_geo.geometry();
  const eigenmath::Vector3d scale = GetScale(transformed_geo.ref_t_shape());

  switch (geo.data_case()) {
    case DataCase::kGeoRef: {
      INTR_ASSIGN_OR_RETURN(
          const auto geometry,
          geolib.GetGeometry(geo.geo_ref(), geo.material_overrides()));
      // Only the renderables are validated.
      if (auto renderable = geometry.GetRenderable(); renderable) {
        INTR_RETURN_IF_ERROR(
            ValidateMeshData(renderable->GetGLBString(), "glb", scale));
      }
    } break;

    case DataCase::kInlineGeometryData: {
      const auto& inline_geo = geo.inline_geometry_data();
      // Only the renderables are validated.
      if (inline_geo.has_renderable()) {
        INTR_RETURN_IF_ERROR(ValidateMeshData(
            inline_geo.renderable().glb_bytes(), "glb", scale));
      } else if (inline_geo.has_generated_renderable()) {
        INTR_RETURN_IF_ERROR(ValidateMeshData(
            inline_geo.generated_renderable().glb_bytes(), "glb", scale));
      }
    } break;

    case DataCase::DATA_NOT_SET:
    default:
      break;
  }

  return absl::OkStatus();
}

}  // namespace

absl::Status ValidateReferencedGeos(const SceneObject& object,
                                    const GeometryDeserializer& geolib) {
  for (const auto& entity : object.entities()) {
    if (!entity.has_link()) {
      continue;
    }

    // Validates `named_geometries` only as `geometries` field has been
    // deprecated.
    for (const auto& [_, geos] :
         entity.link().geometry_component().named_geometries()) {
      for (const auto& [geo_name, transformed_geo] : geos.named_geometries()) {
        // Fails fast and returns the first error.
        INTR_RETURN_IF_ERROR(
            ValidateTransformedGeometry(transformed_geo, geolib))
            << absl::Substitute("Geometry '$0' for link '$1' failed validation",
                                geo_name, entity.name());
      }
    }
  }
  return absl::OkStatus();
}

}  // namespace scene_object
}  // namespace intrinsic
