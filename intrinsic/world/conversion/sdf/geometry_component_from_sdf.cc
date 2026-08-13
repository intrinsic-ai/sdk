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

#include "intrinsic/world/conversion/sdf/geometry_component_from_sdf.h"

#include <optional>
#include <utility>

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/substitute.h"
#include "intrinsic/geometry/api/affine_transform_of_geometry.h"
#include "intrinsic/geometry/api/io.h"
#include "intrinsic/geometry/api/material.h"
#include "intrinsic/scene/sdf/sdf_path_resolver.h"
#include "intrinsic/scene/sdf/sdf_util.h"
#include "intrinsic/util/status/status_builder.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/component/geometry_component.h"
#include "intrinsic/world/geometry_types.h"
#include "intrinsic/world/proto/geometry_component.pb.h"
#include "sdf/Collision.hh"
#include "sdf/Link.hh"
#include "sdf/Visual.hh"

namespace intrinsic {
namespace sdf {

absl::StatusOr<intrinsic_proto::world::GeometryComponent>
GeometryComponentFromSdfLink(const ::sdf::Link& sdf_link,
                             const UriResolver& uri_resolver,
                             GeometrySerializer& geometry_serializer) {
  auto geometry_component = GeometryComponent::Create();
  NamedGeometryProtoSet visual_geometry_set;
  const auto visual_count = sdf_link.VisualCount();
  for (int i = 0; i < visual_count; ++i) {
    const auto& visual = *sdf_link.VisualByIndex(i);
    INTR_ASSIGN_OR_RETURN(
        auto visual_geometry, GeometryFromSdfVisual(visual, uri_resolver),
        _ << absl::Substitute(
            "while parsing $0th visual geometry '$1' of link '$2'", i,
            visual.Name(), sdf_link.Name()));
    if (visual_geometry.has_value()) {
      INTR_ASSIGN_OR_RETURN(
          auto visual_geometry_proto,
          ToProto(*visual_geometry, &geometry_serializer),
          _ << absl::Substitute(
              "generating proto for $0th visual geometry '$1' of link '$2'", i,
              visual.Name(), sdf_link.Name()));
      visual_geometry_set.emplace(visual.Name(),
                                  std::move(visual_geometry_proto));
    }
  }
  geometry_component->SetGeometry(kKindVisualGeometry, visual_geometry_set);

  NamedGeometryProtoSet collision_geometry_set;
  const auto collision_count = sdf_link.CollisionCount();
  for (int i = 0; i < collision_count; ++i) {
    const auto& collision = *sdf_link.CollisionByIndex(i);
    INTR_ASSIGN_OR_RETURN(
        auto collision_geometry,
        GeometryFromSdfCollision(collision, uri_resolver),
        _ << absl::Substitute("while parsing $0th collision '$1' of link '$2'",
                              i, collision.Name(), sdf_link.Name()));
    if (collision_geometry.has_value()) {
      INTR_ASSIGN_OR_RETURN(
          auto collision_geometry_proto,
          ToProto(*collision_geometry, &geometry_serializer),
          _ << absl::Substitute(
              "generating proto for $0th collision geometry '$1' of link '$2'",
              i, collision.Name(), sdf_link.Name()));
      collision_geometry_set.emplace(collision.Name(),
                                     std::move(collision_geometry_proto));
    }
  }
  geometry_component->SetGeometry(kKindCollisionGeometry,
                                  collision_geometry_set);

  return geometry_component->ToProto();
}

absl::StatusOr<std::optional<TransformedGeometry>> GeometryFromSdfCollision(
    const ::sdf::Collision& sdf_collision, const UriResolver& uri_resolver) {
  if (!sdf_collision.Geom()) {
    return std::nullopt;
  }
  const auto& geometry = *sdf_collision.Geom();
  INTR_ASSIGN_OR_RETURN(
      auto transformed_geo,
      ParseGeometry(geometry, uri_resolver, /*material=*/std::nullopt,
                    /*parse_as_collision=*/true));
  if (transformed_geo == std::nullopt) {
    return std::nullopt;
  }

  // Process optional <pose>
  INTR_ASSIGN_OR_RETURN(const auto parent_link_t_this,
                        ParseSemanticPose(sdf_collision.SemanticPose()));
  return TransformedGeometry(
      transformed_geo->shape(),
      parent_link_t_this * transformed_geo->ref_t_shape());
}

absl::StatusOr<std::optional<TransformedGeometry>> GeometryFromSdfVisual(
    const ::sdf::Visual& sdf_visual, const UriResolver& uri_resolver) {
  // Process optional <material>.
  std::optional<Material> material;
  if (const auto* material_sdf = sdf_visual.Material()) {
    if (auto material_intrinsic = ParseMaterial(*material_sdf);
        material_intrinsic.ok()) {
      material = *material_intrinsic;
    } else if (!absl::IsNotFound(material_intrinsic.status())) {
      LOG(WARNING) << material_intrinsic.status();
    }
  }

  if (sdf_visual.Geom() == nullptr) {
    return std::nullopt;
  }
  // Process <geometry>.
  const auto& geometry = *sdf_visual.Geom();
  INTR_ASSIGN_OR_RETURN(auto transformed_geo,
                        ParseGeometry(geometry, uri_resolver, material,
                                      /*parse_as_collision=*/false));

  if (transformed_geo == std::nullopt) {
    return std::nullopt;
  }

  // Process optional <pose>.
  INTR_ASSIGN_OR_RETURN(const auto parent_link_t_this,
                        sdf::ParseSemanticPose(sdf_visual.SemanticPose()));
  return TransformedGeometry(
      transformed_geo->shape(),
      parent_link_t_this * transformed_geo->ref_t_shape());
}

}  // namespace sdf
}  // namespace intrinsic
