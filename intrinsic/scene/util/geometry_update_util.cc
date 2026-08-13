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

#include "intrinsic/scene/util/geometry_update_util.h"

#include <string>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/geometry_types.h"

namespace intrinsic::scene_object {

namespace {

using ::intrinsic_proto::scene_object::v1::GeometryUpdate;

absl::StatusOr<std::string> GetGeometryTypeKey(
    GeometryUpdate::GeometryType type) {
  switch (type) {
    case GeometryUpdate::GEOMETRY_TYPE_VISUAL:
      return intrinsic::kKindVisualGeometry;
    case GeometryUpdate::GEOMETRY_TYPE_COLLISION:
      return intrinsic::kKindCollisionGeometry;
    default:
      return absl::InvalidArgumentError("Invalid geometry type");
  }
}

absl::Status RemoveGeometry(
    ::google::protobuf::Map<
        std::string, intrinsic_proto::world::GeometryComponent::GeometrySet>*
        top_named_geometries,
    absl::string_view top_key, absl::string_view geo_name) {
  std::string top_key_str(top_key);
  if (!top_named_geometries->contains(top_key_str)) {
    return absl::InternalError(absl::StrCat(
        "Unable to remove unknown geometry from '", top_key, "'."));
  }
  intrinsic_proto::world::GeometryComponent::GeometrySet& geo_set =
      (*top_named_geometries)[top_key_str];

  // Always clear v0 geometries if we're editing the v1 geometry.
  geo_set.mutable_geometries()->Clear();

  if (geo_name.empty()) {
    geo_set.mutable_named_geometries()->clear();
  } else {
    geo_set.mutable_named_geometries()->erase(std::string(geo_name));
  }

  if (geo_set.named_geometries().empty() && geo_set.geometries().empty()) {
    top_named_geometries->erase(top_key_str);
  }
  return absl::OkStatus();
}

}  // namespace

absl::Status ApplyGeometryUpdate(
    const GeometryUpdate& update,
    intrinsic_proto::world::GeometryComponent& geometry_component) {
  ::google::protobuf::Map<
      std::string, intrinsic_proto::world::GeometryComponent::GeometrySet>*
      top_named_geometries = geometry_component.mutable_named_geometries();

  for (const GeometryUpdate::GeometryToRemove& remove :
       update.geometries_to_remove()) {
    std::vector<std::string> types_to_remove;
    if (remove.type() == GeometryUpdate::GEOMETRY_TYPE_UNSPECIFIED) {
      types_to_remove = {intrinsic::kKindVisualGeometry,
                         intrinsic::kKindCollisionGeometry};
    } else {
      INTR_ASSIGN_OR_RETURN(std::string type_key,
                            GetGeometryTypeKey(remove.type()));
      types_to_remove.push_back(type_key);
    }

    for (const auto& type_key : types_to_remove) {
      if (top_named_geometries->contains(type_key)) {
        INTR_RETURN_IF_ERROR(RemoveGeometry(top_named_geometries, type_key,
                                            remove.geometry_name()));
      }
    }
  }

  for (const GeometryUpdate::GeometryToSet& set : update.geometries_to_set()) {
    if (set.type() == GeometryUpdate::GEOMETRY_TYPE_UNSPECIFIED) {
      return absl::InvalidArgumentError(absl::StrCat(
          "While updating entity '", update.entity_name(),
          "': GeometryType must be specified for geometries_to_set."));
    }
    if (set.geometry_name().empty()) {
      return absl::InvalidArgumentError(absl::StrCat(
          "While updating entity '", update.entity_name(),
          "': geometry_name must be specified for geometries_to_set."));
    }

    INTR_ASSIGN_OR_RETURN(std::string type_key, GetGeometryTypeKey(set.type()));
    if (!set.has_geometry()) {
      INTR_RETURN_IF_ERROR(
          RemoveGeometry(top_named_geometries, type_key, set.geometry_name()));
      continue;
    }

    intrinsic_proto::world::GeometryComponent::GeometrySet& geo_set =
        (*top_named_geometries)[type_key];
    // Always clear v0 geometries if we're editing the v1 geometry.
    geo_set.mutable_geometries()->Clear();
    (*geo_set.mutable_named_geometries())[set.geometry_name()] = set.geometry();
  }

  return absl::OkStatus();
}

}  // namespace intrinsic::scene_object
