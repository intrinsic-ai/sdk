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

#include "intrinsic/world/component/geometry_component.h"

#include <algorithm>
#include <cstddef>
#include <memory>
#include <ranges>
#include <string>
#include <utility>
#include <vector>

#include "absl/algorithm/container.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/strings/substitute.h"
#include "google/protobuf/repeated_ptr_field.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/affine_transform_of.h"
#include "intrinsic/geometry/api/affine_transform_of_geometry.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/geometry_options.h"
#include "intrinsic/geometry/api/io.h"
#include "intrinsic/geometry/proto/v1/geometry.pb.h"
#include "intrinsic/geometry/proto/v1/transformed_geometry.pb.h"
#include "intrinsic/geometry/storage/geometry_deserializer.h"
#include "intrinsic/geometry/storage/geometry_serializer.h"
#include "intrinsic/math/proto/matrix.pb.h"
#include "intrinsic/math/proto_conversion.h"
#include "intrinsic/util/eigen.h"
#include "intrinsic/util/status/status_builder.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/hashing/hashing.h"
#include "intrinsic/world/proto/geometry_component.pb.h"

namespace intrinsic {
namespace {

// This is an implementation detail while we transition the geometry component
// to a new API that provides deserialization at the point of requesting the
// geometry (see b/460531662). We will replace uses of this type with just the
// proto once the migration is complete.
struct MigrationGeometry {
  std::optional<intrinsic_proto::geometry::v1::TransformedGeometry> proto;
  std::optional<TransformedGeometry> geo;
};

using NamedMigrationGeometrySet = WorldHashMap<std::string, MigrationGeometry>;
using MigrationGeometryMap =
    WorldHashMap<std::string, NamedMigrationGeometrySet>;

// Returns a sorted list of keys for the given named geometry set so that keys
// created from std::vector<TransformedGeometry> are sorted as the original
// vector at the beginning of the list
template <typename T>
std::vector<std::string> GetSortedGeometrySetKeys(
    const WorldHashMap<std::string, T>& geo_set) {
  std::vector<size_t> numbered_keys;
  std::vector<std::string> non_numbered_keys;
  for (const auto& [key, _] : geo_set) {
    size_t key_int;
    if (absl::SimpleAtoi(key, &key_int)) {
      numbered_keys.push_back(key_int);
    } else {
      non_numbered_keys.push_back(key);
    }
  }
  std::sort(numbered_keys.begin(), numbered_keys.end());
  std::sort(non_numbered_keys.begin(), non_numbered_keys.end());

  std::vector<std::string> sorted_geo_set_keys;
  sorted_geo_set_keys.reserve(numbered_keys.size() + non_numbered_keys.size());
  for (const auto& key : numbered_keys) {
    sorted_geo_set_keys.push_back(absl::StrCat(key));
  }
  for (auto&& key : non_numbered_keys) {
    sorted_geo_set_keys.push_back(std::move(key));
  }
  return sorted_geo_set_keys;
}

absl::StatusOr<MigrationGeometry> UpdateGeometryOptions(
    MigrationGeometry geo, const GeometryOptions& options) {
  MigrationGeometry result = std::move(geo);
  bool updated_proto = false;
  if (result.proto.has_value() &&
      result.proto->geometry().has_inline_geometry_data()) {
    intrinsic_proto::geometry::v1::GeometryOptions* result_options =
        result.proto->mutable_geometry()
            ->mutable_inline_geometry_data()
            ->mutable_exact_geometry()
            ->mutable_options();
    if (options.simulation_convex_decomposition_resolution.has_value()) {
      result_options->set_simulation_convex_decomposition_resolution(
          options.simulation_convex_decomposition_resolution.value());
    }
    updated_proto = true;
  } else if (!result.geo.has_value()) {
    return FailedPreconditionErrorBuilder()
           << "Cannot set geometry options for proto (likely with remote ref) "
           << "that has not been deserialized yet.";
  }

  if (result.geo.has_value()) {
    auto geo_shape = result.geo->shape();
    if (geo_shape.GetExactGeometry().options() != options) {
      result.geo = TransformedGeometry(
          Geometry(ExactGeometry(geo_shape.GetExactGeometry(), options),
                   geo_shape.GetRenderable(),
                   geo_shape.KeepRenderableForSerialization(),
                   geo_shape.material_properties(), geo_shape.provenance()),
          result.geo->ref_t_shape());
      if (!updated_proto) {
        // Reset the proto if we've only updated options in a deserialized geo.
        result.proto = std::nullopt;
      }
    }
  }

  return result;
}

class GeometryComponentImpl : public GeometryComponent {
 public:
  GeometryComponentImpl() = default;
  explicit GeometryComponentImpl(MigrationGeometryMap&& geometries)
      : geometries_(std::move(geometries)) {}

  bool HasGeometry(absl::string_view name) const override;
  absl::StatusOr<NamedGeometrySet> GetGeometry(
      absl::string_view name) const override;
  absl::StatusOr<NamedGeometrySet> GetGeometry(
      absl::string_view name, const GeometryDeserializer& geolib) override;
  void SetGeometry(absl::string_view name,
                   const NamedGeometrySet& geometry) override;
  void SetGeometry(absl::string_view name,
                   const NamedGeometryProtoSet& geometry) override;

  absl::Status SetGeometryOptions(const GeometryOptions& options,
                                  absl::string_view geometry_set_name) override;
  absl::Status SetGeometryOptions(const GeometryOptions& options) override;
  absl::Status SetGeometryOptionsImpl(const GeometryOptions& options,
                                      absl::string_view geometry_set_name);

  absl::Status ApplyGeometryOptionOverrides(
      const intrinsic_proto::world::GeometryOptions& options) override;
  absl::Status ApplyGeometryOptionOverrides(
      const intrinsic_proto::world::GeometryOptions& options,
      absl::string_view geometry_set_name) override;
  absl::Status ApplyGeometryOptionOverridesImpl(
      const intrinsic_proto::world::GeometryOptions& options,
      absl::string_view geometry_set_name);

  WorldHashSet<std::string> GetGeometryNames() const override;
  absl::StatusOr<intrinsic_proto::world::GeometryComponent> ToProto(
      GeometrySerializer* geolib) const override;
  absl::StatusOr<intrinsic_proto::world::GeometryComponent> ToProto()
      const override;
  std::unique_ptr<GeometryComponent> Clone() const override;

  absl::Status UpdateFromProto(
      const intrinsic_proto::world::GeometryComponent& proto) override;
  absl::Status UpdateFromProto(
      const intrinsic_proto::world::GeometryComponent& proto,
      const GeometryDeserializer& geolib) override;

 private:
  MigrationGeometryMap geometries_;
};

absl::StatusOr<intrinsic_proto::world::GeometryComponent>
GeometryComponentImpl::ToProto(GeometrySerializer* geolib) const {
  intrinsic_proto::world::GeometryComponent result;
  // Sort the keys so we have deterministic processing of the geometry data.
  auto keys = std::views::keys(geometries_);
  std::vector<std::string> sorted_keys(keys.begin(), keys.end());
  std::sort(sorted_keys.begin(), sorted_keys.end());

  for (const auto& name : sorted_keys) {
    const auto& geo_set = geometries_.at(name);

    // If the component had defined geometries by name, but didn't contain any
    // actual geometries, then the empty geometry set should be reflected in the
    // proto.
    auto& proto_named_geometries =
        *(*result.mutable_named_geometries())[name].mutable_named_geometries();

    for (const auto& [geo_name, transformed_geometry] : geo_set) {
      if (transformed_geometry.geo.has_value()) {
        INTR_ASSIGN_OR_RETURN(
            proto_named_geometries[geo_name],
            intrinsic::ToProto(transformed_geometry.geo.value(), geolib));
      } else if (transformed_geometry.proto.has_value()) {
        proto_named_geometries[geo_name] = transformed_geometry.proto.value();
      } else {
        return InternalErrorBuilder()
               << "Named geometry '" << name << "[" << geo_name
               << "]' has neither proto nor geo";
      }
    }

  }

  return result;
}

absl::StatusOr<intrinsic_proto::world::GeometryComponent>
GeometryComponentImpl::ToProto() const {
  intrinsic_proto::world::GeometryComponent result;
  // Sort the keys so we have deterministic processing of the geometry data.
  auto keys = std::views::keys(geometries_);
  std::vector<std::string> sorted_keys(keys.begin(), keys.end());
  std::sort(sorted_keys.begin(), sorted_keys.end());

  for (const auto& name : sorted_keys) {
    const auto& proto_and_geo_map = geometries_.at(name);

    // If the component had defined geometries by name, but didn't contain any
    // actual geometries, then the empty geometry set should be reflected in the
    // proto.
    auto& proto_named_geometries =
        *(*result.mutable_named_geometries())[name].mutable_named_geometries();

    for (const auto& [geo_name, proto_and_geo] : proto_and_geo_map) {
      const auto& [proto, geo] = proto_and_geo;
      if (!proto.has_value()) {
        return FailedPreconditionErrorBuilder()
               << "Cannot serialize unserialized geometry '" << geo_name
               << "'. Please call `SetGeometry` with the appropriate proto for "
               << "this geometry.";
      }

      proto_named_geometries[geo_name] = *proto;
    }
  }

  return result;
}

bool GeometryComponentImpl::HasGeometry(absl::string_view name) const {
  auto iter = geometries_.find(name);
  return iter != geometries_.end() && !iter->second.empty() && [iter]() {
    for (const auto& [_, proto_and_geo] : iter->second) {
      const auto& [proto, geo] = proto_and_geo;
      if (!proto.has_value() && !geo.has_value()) {
        return false;
      }
    }
    return true;
  }();
}

absl::StatusOr<NamedGeometrySet> GeometryComponentImpl::GetGeometry(
    absl::string_view name) const {
  if (!HasGeometry(name)) {
    return intrinsic::NotFoundErrorBuilder()
           << "Could not find any '" << name << "' geometry for this entity";
  }

  NamedGeometrySet geo_set;
  for (const auto& [geo_name, proto_and_geo] : geometries_.at(name)) {
    const auto& [_, geo] = proto_and_geo;
    if (!geo.has_value()) {
      return FailedPreconditionErrorBuilder()
             << "Geometry '" << geo_name << "' in set named '" << name
             << "' doesn't have a deserialized geometry. Please use the "
             << "version of `GetGeometry` that passes a geometry deserializer.";
    }

    geo_set.emplace(geo_name, *geo);
  }
  return geo_set;
}

absl::StatusOr<NamedGeometrySet> GeometryComponentImpl::GetGeometry(
    absl::string_view name, const GeometryDeserializer& geolib) {
  if (!HasGeometry(name)) {
    return intrinsic::NotFoundErrorBuilder()
           << "Could not find any '" << name << "' geometry for this entity";
  }

  NamedGeometrySet geo_set;
  for (auto& [geo_name, proto_and_geo] : geometries_.at(name)) {
    auto& [proto, geo] = proto_and_geo;
    if (!geo.has_value()) {
      if (!proto.has_value()) {
        LOG(WARNING)
            << "Geometry '" << geo_name << "' in set named '" << name
            << "' doesn't have either a deserialized geometry nor a geometry "
            << "proto. Skipping.";
        continue;
      }

      INTR_ASSIGN_OR_RETURN(
          geo, ToGeometry(*proto, &geolib),
          _ << "could not load GeometryComponent geometry from proto");
    }

    geo_set.emplace(geo_name, *geo);
  }
  return geo_set;
}

void GeometryComponentImpl::SetGeometry(absl::string_view name,
                                        const NamedGeometrySet& geometry) {
  const std::string key{name};
  if (geometry.empty()) {
    geometries_.erase(key);
    return;
  }

  LOG(WARNING) << "Setting geometry '" << name
               << "' without providing a serialized proto. This world will "
               << "likely need a GeometrySerializer to serialize properly";

  for (const auto& [geo_name, geo] : geometry) {
    geometries_[key][geo_name] = {std::nullopt, geo};
  }
}

void GeometryComponentImpl::SetGeometry(absl::string_view name,
                                        const NamedGeometryProtoSet& geometry) {
  const std::string key{name};
  if (geometry.empty()) {
    geometries_.erase(key);
    return;
  }

  for (const auto& [geo_name, proto] : geometry) {
    geometries_[key][geo_name] = {proto, std::nullopt};
  }
}

WorldHashSet<std::string> GeometryComponentImpl::GetGeometryNames() const {
  WorldHashSet<std::string> result;
  for (const auto& name : std::views::keys(geometries_)) {
    if (HasGeometry(name)) {
      result.insert(name);
    }
  }
  return result;
}

// Same as SetGeometryOptions but skips the check for the geometry set name.
absl::Status GeometryComponentImpl::SetGeometryOptionsImpl(
    const GeometryOptions& options, absl::string_view geometry_set_name) {
  for (auto& [_, proto_and_geo] : geometries_.at(geometry_set_name)) {
    INTR_ASSIGN_OR_RETURN(proto_and_geo,
                          UpdateGeometryOptions(proto_and_geo, options));
  }

  return absl::OkStatus();
}

absl::Status GeometryComponentImpl::SetGeometryOptions(
    const GeometryOptions& options) {
  for (auto& [name, _] : geometries_) {
    INTR_RETURN_IF_ERROR(SetGeometryOptionsImpl(options, name));
  }

  return absl::OkStatus();
}

absl::Status GeometryComponentImpl::SetGeometryOptions(
    const GeometryOptions& options, absl::string_view geometry_set_name) {
  if (!geometries_.contains(geometry_set_name)) {
    return absl::NotFoundError(
        absl::StrCat("Geometry set '", geometry_set_name,
                     "' not found in geometry component"));
  }

  return SetGeometryOptionsImpl(options, geometry_set_name);
}

// Same as ApplyGeometryOptionOverrides but skips the check for the geometry
// set name.
absl::Status GeometryComponentImpl::ApplyGeometryOptionOverridesImpl(
    const intrinsic_proto::world::GeometryOptions& options,
    absl::string_view geometry_set_name) {
  for (auto& [_, proto_and_geo] : geometries_.at(geometry_set_name)) {
    GeometryOptions merged_options =
        proto_and_geo.geo->shape().GetExactGeometry().options();
    ApplyOverrides(merged_options, options);
    INTR_ASSIGN_OR_RETURN(proto_and_geo,
                          UpdateGeometryOptions(proto_and_geo, merged_options));
  }

  return absl::OkStatus();
}

absl::Status GeometryComponentImpl::ApplyGeometryOptionOverrides(
    const intrinsic_proto::world::GeometryOptions& options) {
  for (auto& [name, _] : geometries_) {
    INTR_RETURN_IF_ERROR(ApplyGeometryOptionOverridesImpl(options, name));
  }

  return absl::OkStatus();
}

absl::Status GeometryComponentImpl::ApplyGeometryOptionOverrides(
    const intrinsic_proto::world::GeometryOptions& options,
    absl::string_view geometry_set_name) {
  if (!geometries_.contains(geometry_set_name)) {
    return absl::NotFoundError(
        absl::StrCat("Geometry set '", geometry_set_name,
                     "' not found in geometry component"));
  }

  return ApplyGeometryOptionOverridesImpl(options, geometry_set_name);
}

absl::StatusOr<NamedMigrationGeometrySet> ParseGeometrySet(
    const intrinsic_proto::world::GeometryComponent::GeometrySet& geometry_set,
    const GeometryDeserializer* geolib) {
  NamedMigrationGeometrySet named_set;
  if (!geometry_set.named_geometries().empty()) {
    if (!geometry_set.geometries().empty()) {
      VLOG(1) << absl::Substitute(
          "Both named_geometries and geometries are set. Using "
          "named_geometries with "
          "size $0 and ignoring geometries with size $1.",
          geometry_set.geometries_size(), geometry_set.geometries_size());
    }

    for (const auto& [name, geo_proto] : geometry_set.named_geometries()) {
      if (geolib != nullptr) {
        INTR_ASSIGN_OR_RETURN(
            auto geometry, ToGeometry(geo_proto, geolib),
            _ << "could not load GeometryComponent geometry from proto");
        named_set.insert_or_assign(
            name, {.proto = geo_proto, .geo = std::move(geometry)});
      } else {
        named_set.insert_or_assign(name,
                                   {.proto = geo_proto, .geo = std::nullopt});
      }
    }
  }
  return named_set;
}

std::unique_ptr<GeometryComponent> GeometryComponentImpl::Clone() const {
  MigrationGeometryMap geometry_copy = geometries_;
  return std::make_unique<GeometryComponentImpl>(std::move(geometry_copy));
}

absl::StatusOr<MigrationGeometryMap> ParseGeometryMap(
    const intrinsic_proto::world::GeometryComponent& proto,
    const GeometryDeserializer* geolib) {
  MigrationGeometryMap geometries;
  for (const auto& [name, geo_set] : proto.named_geometries()) {
    INTR_ASSIGN_OR_RETURN(auto parsed_set, ParseGeometrySet(geo_set, geolib),
                          _ << "Parsing geometry named '" << name << "'");
    if (!parsed_set.empty()) {
      geometries[name] = std::move(parsed_set);
    }
  }

  return geometries;
}

absl::Status GeometryComponentImpl::UpdateFromProto(
    const intrinsic_proto::world::GeometryComponent& proto,
    const GeometryDeserializer& geolib) {

  INTR_ASSIGN_OR_RETURN(geometries_, ParseGeometryMap(proto, &geolib));
  return absl::OkStatus();
}

absl::Status GeometryComponentImpl::UpdateFromProto(
    const intrinsic_proto::world::GeometryComponent& proto) {
  INTR_ASSIGN_OR_RETURN(geometries_, ParseGeometryMap(proto, nullptr));
  return absl::OkStatus();
}

}  // namespace

std::unique_ptr<GeometryComponent> GeometryComponent::Create() {
  return std::make_unique<GeometryComponentImpl>();
}

absl::StatusOr<std::unique_ptr<GeometryComponent>> GeometryComponent::FromProto(
    const intrinsic_proto::world::GeometryComponent& proto,
    const GeometryDeserializer& geolib) {
  auto result = Create();
  INTR_RETURN_IF_ERROR(result->UpdateFromProto(proto, geolib));
  return result;
}

absl::StatusOr<std::unique_ptr<GeometryComponent>> GeometryComponent::FromProto(
    const intrinsic_proto::world::GeometryComponent& proto) {
  auto result = Create();
  INTR_RETURN_IF_ERROR(result->UpdateFromProto(proto));
  return result;
}

void ApplyOverrides(GeometryOptions& options,
                    const intrinsic_proto::world::GeometryOptions& proto) {
}

GeometryOptions FromProto(
    const intrinsic_proto::world::GeometryOptions& proto) {
  GeometryOptions options = GeometryOptions::Default();
  ApplyOverrides(options, proto);
  return options;
}

intrinsic_proto::world::GeometryOptions ToProto(const GeometryOptions& shape) {
  intrinsic_proto::world::GeometryOptions proto;
  return proto;
}

}  // namespace intrinsic
