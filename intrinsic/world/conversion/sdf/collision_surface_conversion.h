// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_WORLD_CONVERSION_SDF_COLLISION_SURFACE_CONVERSION_H_
#define INTRINSIC_WORLD_CONVERSION_SDF_COLLISION_SURFACE_CONVERSION_H_

#include <optional>

#include "absl/status/statusor.h"
#include "intrinsic/world/proto/physics_component.pb.h"
#include "sdf/Link.hh"
#include "sdf/Surface.hh"

// Utility functions for converting sdf <surface> element under the <collision>
// element. See http://sdformat.org/spec?elem=collision#collision_surface for
// reference.
namespace intrinsic {
namespace sdf {

struct ParseContactResult {
  std::optional<double> kp;
  std::optional<double> kd;
  std::optional<bool> collide_without_contact;
};
// Parses a <contact> element and update both the relevant `physics_component`
// and `collision_component. See
// http://sdformat.org/spec?elem=collision#surface_contact for reference.
absl::StatusOr<ParseContactResult> ParseContact(const ::sdf::Contact& contact);

struct ParseFrictionResult {
  std::optional<intrinsic_proto::world::Torsional> torsional;
  std::optional<intrinsic_proto::world::Friction> surface_friction;
};

// Parses a <friction> element and update the relevant `physics_component`. See
// http://sdformat.org/spec?elem=collision#surface_friction for reference.
absl::StatusOr<ParseFrictionResult> ParseFriction(
    const ::sdf::Friction& friction);

struct ParseSurfaceResult {
  ParseContactResult contact;
  ParseFrictionResult friction;
};

// Parses a <surface> element converted to `sdf_surface` and update
// `physics_component` and optionally `collision_component`. See
// http://sdformat.org/spec?elem=collision#collision_surface for reference.
absl::StatusOr<ParseSurfaceResult> ParseSurface(
    const ::sdf::Surface& sdf_surface);

// Gets the single identical sdf surface from all <collision> elements in
// `sdf_link`. Returns nullptr if there are no available <surface> elements.
// Returns an error multiple surfaces with conflicting xml exists.
absl::StatusOr<const ::sdf::Surface*> GetSingleSurface(
    const ::sdf::Link& sdf_link);
}  // namespace sdf
}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_CONVERSION_SDF_COLLISION_SURFACE_CONVERSION_H_
