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

#include "intrinsic/geometry/api/geometry_fingerprint.h"

#include <string>

#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/synchronization/mutex.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/io.h"
#include "intrinsic/geometry/api/renderable.h"
#include "intrinsic/util/hash.h"
#include "intrinsic/util/macros.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::geo {

absl::StatusOr<std::string> GenerateFingerprint(const Geometry& geometry) {
  absl::MutexLock lock(geometry.fingerprint_mutex_);
  if (geometry.fingerprint_.empty()) {
    INTR_ASSIGN_OR_RETURN(auto proto, ToInlinedProto(geometry));
    geometry.fingerprint_ = absl::StrCat(
        absl::Hex(intrinsic::Fingerprint(proto.SerializeAsString())));
  }

  return geometry.fingerprint_;
}

absl::StatusOr<std::string> GenerateFingerprint(
    const ExactGeometry& exact_geo) {
  INTR_ASSIGN_OR_RETURN(auto proto, ToProto(exact_geo));
  return absl::StrCat(
      absl::Hex(intrinsic::Fingerprint(proto.SerializeAsString())));
}

std::string GenerateFingerprintOrDie(const Geometry& geometry) {
  ASSIGN_OR_DIE(auto fingerprint, GenerateFingerprint(geometry));
  return fingerprint;
}

std::string GenerateFingerprintOrDie(const ExactGeometry& exact_geo) {
  ASSIGN_OR_DIE(auto fingerprint, GenerateFingerprint(exact_geo));
  return fingerprint;
}

std::string GenerateFingerprint(const Renderable& renderable) {
  return absl::StrCat(
      absl::Hex(intrinsic::Fingerprint(renderable.GetGLBString())));
}

}  // namespace intrinsic::geo
