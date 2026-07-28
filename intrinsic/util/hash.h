// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_HASH_H_
#define INTRINSIC_UTIL_HASH_H_

#include <cstdint>

#include "absl/strings/string_view.h"
#include "third_party/thorough_hash/thorough_hash.h"

namespace intrinsic {

// Returns a fingerprint of the given string.
inline uint64_t Fingerprint(absl::string_view s) {
  return ThoroughHash(s.data(), s.size());
}

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_HASH_H_
