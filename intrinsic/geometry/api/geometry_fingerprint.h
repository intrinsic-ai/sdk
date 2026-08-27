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

#ifndef INTRINSIC_GEOMETRY_API_GEOMETRY_FINGERPRINT_H_
#define INTRINSIC_GEOMETRY_API_GEOMETRY_FINGERPRINT_H_

#include <string>

#include "absl/status/statusor.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/renderable.h"

namespace intrinsic::geo {

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

}  // namespace intrinsic::geo

namespace intrinsic {
using ::intrinsic::geo::GenerateFingerprint;
using ::intrinsic::geo::GenerateFingerprintOrDie;
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_GEOMETRY_FINGERPRINT_H_
