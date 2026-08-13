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

// This file provides a helper function for working with the object store
// memoization function. It provides an automated way of converting a function
// and its arguments into a key to be used with the object store memoization
// calls.
//
// See the Memoize function documentation for further details on usage.
#ifndef INTRINSIC_UTIL_OBJECT_STORE_MEMOIZE_FINGERPRINT_H_
#define INTRINSIC_UTIL_OBJECT_STORE_MEMOIZE_FINGERPRINT_H_

#include <cstdint>
#include <set>
#include <string>
#include <utility>
#include <vector>

#include "absl/log/check.h"
#include "absl/meta/type_traits.h"
#include "absl/strings/string_view.h"
#include "intrinsic/marshal/riegeli_coder.h"
#include "intrinsic/util/hash.h"
#include "intrinsic/util/string_type.h"
#include "intrinsic/util/thread/stop_token.h"
#include "riegeli/base/maker.h"
#include "riegeli/bytes/string_writer.h"
#include "riegeli/records/record_writer.h"

namespace intrinsic {
namespace memoize {

// Struct for computing the fingerprint of a value. Defaults to using Coder
// encoding and fingerprinting the resulting string. More specialization can be
// added for types that want to support more efficient fingerprinting.
template <typename T>
struct MemoizeFingerprint {
  static uint64_t Fingerprint(const T& value) {
    std::string encoded;
    riegeli::RecordWriter writer(
        riegeli::Maker<riegeli::StringWriter>(&encoded),
        riegeli::RecordWriterBase::Options().set_uncompressed());
    CHECK_OK(RiegeliCoder<absl::decay_t<T>>::Encode(value, writer));
    CHECK(writer.Close()) << "Failed to close writer: " << writer.status();
    return intrinsic::Fingerprint(encoded);
  }
};

// Specialization for fingerprinting strings, it skips the Encode step as it is
// already a string.
template <>
struct MemoizeFingerprint<std::string> {
  static uint64_t Fingerprint(const std::string& value) {
    return intrinsic::Fingerprint(value);
  }
};

// Specialization for fingerprinting string_view, it skips the Encode step as it
// is already a string.
template <>
struct MemoizeFingerprint<absl::string_view> {
  static uint64_t Fingerprint(absl::string_view value) {
    return intrinsic::Fingerprint(value);
  }
};

// Specialization for fingerprinting std::vector, deferring to the underlying T
// for the actual fingerprinting.
template <typename T>
struct MemoizeFingerprint<std::vector<T>> {
  static uint64_t Fingerprint(const std::vector<T>& value) {
    if (value.empty()) return 0;

    auto fp = MemoizeFingerprint<T>::Fingerprint(value[0]);
    for (int i = 1; i < value.size(); ++i) {
      fp = MixTwoUInt64(
          fp, MemoizeFingerprint<absl::decay_t<T>>::Fingerprint(value[i]));
    }
    return fp;
  }
};

// Specialization for fingerprinting std::set, deferring to the underlying T
// for the actual fingerprinting.
template <typename T>
struct MemoizeFingerprint<std::set<T>> {
  static uint64_t Fingerprint(const std::set<T>& value) {
    if (value.empty()) return 0;

    uint64_t fp = 7;
    for (const auto& value_i : value) {
      fp = MixTwoUInt64(
          fp, MemoizeFingerprint<absl::decay_t<T>>::Fingerprint(value_i));
    }
    return fp;
  }
};

// Specialization for fingerprinting std::pair, deferring to the underlying T1
// and T2 for the actual fingerprinting.
template <typename T1, typename T2>
struct MemoizeFingerprint<std::pair<T1, T2>> {
  static uint64_t Fingerprint(const std::pair<T1, T2>& value) {
    const auto fp1 =
        MemoizeFingerprint<absl::decay_t<T1>>::Fingerprint(value.first);
    const auto fp2 =
        MemoizeFingerprint<absl::decay_t<T2>>::Fingerprint(value.second);
    return MixTwoUInt64(fp1, fp2);
  }
};

// Specialization for fingerprinting intrinsic::StringType, it skips the Encode
// step as it is already a string.
template <typename T>
struct MemoizeFingerprint<intrinsic::StringType<T>> {
  static uint64_t Fingerprint(const intrinsic::StringType<T>& value) {
    return intrinsic::Fingerprint(value.value());
  }
};

// Specialization for fingerprinting intrinsic::StopToken, it returns a static
// fingerprint. We don't want to include the stop token in the memoization key
// as it will be unique each time.
template <>
struct MemoizeFingerprint<intrinsic::StopToken> {
  static auto Fingerprint(const intrinsic::StopToken& /*token*/) {
    return intrinsic::Fingerprint("StopToken");
  }
};

}  // namespace memoize
}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_OBJECT_STORE_MEMOIZE_FINGERPRINT_H_
