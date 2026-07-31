// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/object_store/object_store_internal.h"

#include <cstddef>
#include <cstdint>
#include <initializer_list>
#include <string>

#include "absl/flags/flag.h"
#include "absl/log/check.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "intrinsic/util/hash.h"
#include "intrinsic/util/object_store/object_cache.h"
#include "third_party/thorough_hash/thorough_hash.h"

ABSL_FLAG(size_t, object_store_cache_size_mb, 500,
          "The size of the object store cache in MB.");

namespace intrinsic {
namespace object_store_internal {

std::string GenerateId(std::initializer_list<absl::string_view> args) {
  CHECK_GT(args.size(), 0);
  auto iter = args.begin();
  uint64_t fp = intrinsic::Fingerprint(*iter++);
  while (iter != args.end()) {
    fp = MixTwoUInt64(fp, intrinsic::Fingerprint(*iter++));
  }
  return absl::StrCat(absl::Hex(fp));
}

std::string GenerateId(std::initializer_list<uint64_t> args) {
  CHECK_GT(args.size(), 0);
  auto iter = args.begin();
  uint64_t fp = *iter++;
  while (iter != args.end()) {
    fp = MixTwoUInt64(fp, *iter++);
  }
  return absl::StrCat(absl::Hex(fp));
}

ObjectStore& GlobalObjectStore() {
  static ObjectStore* current = new ObjectStore();
  return *current;
}

inline constexpr size_t kBytesPerMegaByte = 10e6;
ObjectStore::ObjectStore()
    : cache_(absl::GetFlag(FLAGS_object_store_cache_size_mb) *
             kBytesPerMegaByte) {}

}  // namespace  object_store_internal
}  // namespace intrinsic
