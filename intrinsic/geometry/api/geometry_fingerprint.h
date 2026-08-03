// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_API_GEOMETRY_FINGERPRINT_H_
#define INTRINSIC_GEOMETRY_API_GEOMETRY_FINGERPRINT_H_

#include <string>

#include "absl/status/statusor.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/renderable.h"

namespace intrinsic {

// Computes the fingerprint of the given input geometry.
absl::StatusOr<std::string> GenerateFingerprint(const Geometry& geometry);

// Computes the fingerprint of the given input geometry.
absl::StatusOr<std::string> GenerateFingerprint(const ExactGeometry& exact_geo);

// Computes the fingerprint of the given input geometry, or dies.
std::string GenerateFingerprintOrDie(const Geometry& geometry);

// Computes the fingerprint of the given input geometry, or dies.
std::string GenerateFingerprintOrDie(const ExactGeometry& exact_geo);

// Computes the fingerprint of the given input geometry.
std::string GenerateFingerprint(const Renderable& renderable);

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_GEOMETRY_FINGERPRINT_H_
