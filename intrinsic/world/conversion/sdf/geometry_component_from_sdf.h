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

#ifndef INTRINSIC_WORLD_CONVERSION_SDF_GEOMETRY_COMPONENT_FROM_SDF_H_
#define INTRINSIC_WORLD_CONVERSION_SDF_GEOMETRY_COMPONENT_FROM_SDF_H_

#include <optional>

#include "absl/status/statusor.h"
#include "intrinsic/geometry/api/affine_transform_of_geometry.h"
#include "intrinsic/geometry/api/material.h"
#include "intrinsic/geometry/storage/geometry_serializer.h"
#include "intrinsic/scene/sdf/sdf_path_resolver.h"
#include "intrinsic/world/proto/geometry_component.pb.h"
#include "sdf/Collision.hh"
#include "sdf/Link.hh"
#include "sdf/Visual.hh"

// Utility functions that create GeometryComponent and related data types from
// SDF <visual> and <collision> elements.
namespace intrinsic {
namespace sdf {

// Construct a GeometryComponent from the given `sdf_link`.
// A `uri_resolver` is needed to resolve geometry file paths referenced by
// `sdf_link`'s children. A `geometry_serializer` is needed to store the
// geometry separately from the GeometryComponent proto.
absl::StatusOr<intrinsic_proto::world::GeometryComponent>
GeometryComponentFromSdfLink(const ::sdf::Link& sdf_link,
                             const UriResolver& uri_resolver,
                             GeometrySerializer& geometry_serializer);

// Construct a optional `TransformedGeometry` from the given `sdf_collision`.
// A `uri_resolver` is needed to resolve geometry file paths referenced by
// `sdf_collision`'s children.
absl::StatusOr<std::optional<TransformedGeometry>> GeometryFromSdfCollision(
    const ::sdf::Collision& sdf_collision, const UriResolver& uri_resolver);

// Construct a optional `TransformedGeometry` from the given `sdf_visual`.
// A `uri_resolver` is needed to resolve geometry file paths referenced by
// `sdf_visual`'s children.
absl::StatusOr<std::optional<TransformedGeometry>> GeometryFromSdfVisual(
    const ::sdf::Visual& sdf_visual, const UriResolver& uri_resolver);
}  // namespace sdf
}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_CONVERSION_SDF_GEOMETRY_COMPONENT_FROM_SDF_H_
