// Copyright 2023 Intrinsic Innovation LLC

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
