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

#ifndef INTRINSIC_WORLD_COMPONENT_GEOMETRY_COMPONENT_H_
#define INTRINSIC_WORLD_COMPONENT_GEOMETRY_COMPONENT_H_

#include <memory>
#include <string>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/geometry/api/affine_transform_of_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/geometry_options.h"
#include "intrinsic/geometry/storage/geometry_deserializer.h"
#include "intrinsic/geometry/storage/geometry_serializer.h"
#include "intrinsic/world/hashing/hashing.h"
#include "intrinsic/world/proto/geometry_component.pb.h"

namespace intrinsic {

using NamedGeometrySet = WorldHashMap<std::string, TransformedGeometry>;
using NamedGeometryProtoSet =
    WorldHashMap<std::string,
                 intrinsic_proto::geometry::v1::TransformedGeometry>;

// A component to hold the geometric representations of an entity.
class GeometryComponent {
 public:
  virtual ~GeometryComponent() = default;

  // Returns a new GeometryComponent instance.
  static std::unique_ptr<GeometryComponent> Create();

  // Returns a new GeometryComponent instance derived from the given proto.
  ABSL_DEPRECATED("Use the non-geometry-deserializing deserializer instead")
  static absl::StatusOr<std::unique_ptr<GeometryComponent>> FromProto(
      const intrinsic_proto::world::GeometryComponent& proto,
      const GeometryDeserializer& geolib);

  // Returns a new GeometryComponent instance derived from the given proto.
  static absl::StatusOr<std::unique_ptr<GeometryComponent>> FromProto(
      const intrinsic_proto::world::GeometryComponent& proto);

  // Returns a new GeometryComponent instance that is a copy of this one.
  virtual std::unique_ptr<GeometryComponent> Clone() const = 0;

  // Returns a proto representation of this component.
  ABSL_DEPRECATED("Use ToProto() instead (without arguments)")
  virtual absl::StatusOr<intrinsic_proto::world::GeometryComponent> ToProto(
      GeometrySerializer* geolib) const = 0;

  // Returns a proto representation of this component.
  virtual absl::StatusOr<intrinsic_proto::world::GeometryComponent> ToProto()
      const = 0;

  // Returns true if this object contains a geometry for the given name.
  virtual bool HasGeometry(absl::string_view name) const = 0;

  // Returns the set of geometry used to represent this entity given some
  // identifier or not found error if the geometry does not exist for this
  // object. There are two names with special semantics, which are defined in
  // geometry_types.h
  ABSL_DEPRECATED("Pass a geometry deserializer to get your geometry.")
  virtual absl::StatusOr<NamedGeometrySet> GetGeometry(
      absl::string_view name) const = 0;

  // Returns the set of geometry used to represent this entity given some
  // identifier or not found error if the geometry does not exist for this
  // object. There are two names with special semantics, which are defined in
  // geometry_types.h
  virtual absl::StatusOr<NamedGeometrySet> GetGeometry(
      absl::string_view name, const GeometryDeserializer& geolib) = 0;
  ABSL_DEPRECATED("Use the version that provides the serialized proto")
  virtual void SetGeometry(absl::string_view name,
                           const NamedGeometrySet& geometry) = 0;

  // Sets the geometry for this entity for the given name.
  virtual void SetGeometry(absl::string_view name,
                           const NamedGeometryProtoSet& geometry) = 0;

  // Updates the geometry options for all the geometries of this component.
  virtual absl::Status SetGeometryOptions(const GeometryOptions& options) = 0;

  // Updates the geometry options for the geometries of that are part of the
  // given set.
  virtual absl::Status SetGeometryOptions(
      const GeometryOptions& options, absl::string_view geometry_set_name) = 0;

  // Applies the geometry option overrides for all the geometries of this
  // component. If a field is not set in the overrides, it will not modify the
  // equivalent field in the existing (or default) options for the geometry.
  virtual absl::Status ApplyGeometryOptionOverrides(
      const intrinsic_proto::world::GeometryOptions& options) = 0;

  // Applies the geometry option overrides for all the geometries of the given
  // set. If a field is not set in the overrides, it will not modify the
  // equivalent field in the existing (or default) options for the geometry.
  virtual absl::Status ApplyGeometryOptionOverrides(
      const intrinsic_proto::world::GeometryOptions& options,
      absl::string_view geometry_set_name) = 0;

  // Returns a list of geometries used to represent this entity.
  virtual WorldHashSet<std::string> GetGeometryNames() const = 0;

  // Updates the component based on the given proto. This will do a full
  // override, if something is missing from this proto it will override any
  // existing values with the defaults.
  ABSL_DEPRECATED(
      "Please use the non-geometry-deserializing version of this function.")
  virtual absl::Status UpdateFromProto(
      const intrinsic_proto::world::GeometryComponent& proto,
      const GeometryDeserializer& geolib) = 0;

  // Updates the component based on the given proto. This will do a full
  // override, if something is missing from this proto it will override any
  // existing values with the defaults.
  virtual absl::Status UpdateFromProto(
      const intrinsic_proto::world::GeometryComponent& proto) = 0;
};

// Applies the overrides from the proto to the C++ struct, if the fields are
// present in the proto.
void ApplyOverrides(GeometryOptions& options,
                    const intrinsic_proto::world::GeometryOptions& proto);

// Converts the proto representation of GeometryOptions to the C++ struct.
GeometryOptions FromProto(const intrinsic_proto::world::GeometryOptions& proto);

// Converts the C++ struct representation of GeometryOptions to the proto.
intrinsic_proto::world::GeometryOptions ToProto(const GeometryOptions& shape);

}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_COMPONENT_GEOMETRY_COMPONENT_H_
