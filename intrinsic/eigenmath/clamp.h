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

#ifndef INTRINSIC_EIGENMATH_CLAMP_H_
#define INTRINSIC_EIGENMATH_CLAMP_H_

#include <algorithm>

#include "absl/base/attributes.h"

namespace intrinsic::eigenmath {

// Clamps vector to boundaries and returns true on success. Returns false on
// dimension mismatch. Expects lower <= upper bound, will return meaningless
// result if lower > upper.
template <typename T, typename LB, typename UB>
inline ABSL_MUST_USE_RESULT bool ClampVector(const LB& lower, const UB& upper,
                                             T& v) {
  // Check correctness of dimensions.
  if (lower.size() != v.size() || upper.size() != v.size()) return false;

  for (auto i = 0; i < v.size(); i++) {
    v(i) = std::clamp(v(i), lower(i), upper(i));
  }
  return true;
}

}  // namespace intrinsic::eigenmath

#endif  // INTRINSIC_EIGENMATH_CLAMP_H_
