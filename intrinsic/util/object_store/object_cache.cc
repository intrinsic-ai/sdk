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

#include "intrinsic/util/object_store/object_cache.h"

#include <any>
#include <cstddef>
#include <cstdint>
#include <mutex>  // NOLINT Requied for std::recursive_mutex
#include <string>

#include "absl/flags/flag.h"
#include "absl/log/log.h"
#include "absl/strings/str_format.h"
#include "absl/strings/string_view.h"
#include "absl/types/any.h"
#include "intrinsic/util/lru_cache.h"

ABSL_FLAG(double, object_store_cache_max_idle_seconds, -1,
          "After an object has be unused for this many seconds, it will get "
          "evicted from the object store cache. negative means never.");

namespace intrinsic {

ObjectCache::ObjectStoreLruCache::ObjectStoreLruCache(size_t cache_size)
    : LruCache<std::string, const std::any>(cache_size) {}

ObjectCache::ObjectCache(size_t cache_size) : lru_cache_(cache_size) {
  lru_cache_.setMaxIdleSeconds(
      absl::GetFlag(FLAGS_object_store_cache_max_idle_seconds));
}

ObjectCache::~ObjectCache() {
  // Objects stored in cache might reference each other, which means that the
  // existence of an object in cache (even if not referenced from outside), may
  // cause another object to become pinned. Therefore, even if no object is
  // referenced externally, calling Clear() might be invalid since some objects
  // may still be pinned. To avoid that, we first RemoveAll(), which removes all
  // the objects from cache and causes all *unpinned* objects to be deleted.
  // Once that is done, we are expecting no objects to still be referenced, and
  // Clear() to succeed. If Clear() fails, it mean that there is a leaked
  // external reference somewhere, or a cyclic reference (which should be
  // avoided).

  lru_cache_.removeAll();

  // If we are trying to debug object store cache issues this will help by
  // giving us extra information about the pinned items that are in the cache
  // and were not removed by the previous call.

  lru_cache_.clear();
}

void ObjectCache::SafeRelease(absl::string_view id,
                              const std::any* value) const {
  std::lock_guard<std::recursive_mutex> lock(mutex_);
  lru_cache_.release(std::string(id), value);
}

}  // namespace intrinsic
