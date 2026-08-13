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

// A simple, thread-safe LRU cache for arbitrary-typed, immutable objects.
//
// - The cache key is always a string.
// - The cache owns any objects inserted to it and will keep them alive as long
//   as there are live references to them.
// - The cache is constructed with a size limit (units are arbitrary, may be
//   chosen by the client to be number of bytes, number of objects, etc.) and
//   the cache will try to limit to total size of objects contained in it to be
//   at or below the specified limit (as much as it is possible, given that live
//   references may force a minimum size).
// - Before the cache can be destructed, all object references must be released.
//
// Typical usage:
//
// ObjectCache cache(kMaxSizeBytes);
//
// auto int_ref = cache.Insert("1234",
//                             make_unique<const int>(4321),
//                             sizeof(int));
// CHECK_EQ(int_ref->Id(), "1234");
// CHECK_EQ(int_ref->Value(), 4321);
//
// // Copying references is safe, can now delete the original.
// auto int_ref_copy = int_ref;
//
// cache.Insert("314",
//              make_unique<const float>(3.14f),
//              sizeof(float));
//
// // Client is responsible for knowing the right type when fetching.
// auto float_ref = cache.Lookup<float>("314");
// CHECK(float_ref);
// CHECK_EQ(float_ref->Id(), "314");
// CHECK_EQ(float_ref->Value(), 3.14f);
#ifndef INTRINSIC_UTIL_OBJECT_STORE_OBJECT_CACHE_H_
#define INTRINSIC_UTIL_OBJECT_STORE_OBJECT_CACHE_H_

#include <any>
#include <cstddef>
#include <cstdint>
#include <memory>
#include <mutex>  // NOLINT(build/c++11)
#include <string>
#include <utility>

#include "absl/strings/string_view.h"
#include "absl/types/any.h"
#include "intrinsic/util/lru_cache.h"

namespace intrinsic {

// A simple, thread-safe LRU cache for arbitrary-typed, immutable objects.
// See top-of-file comment for details.
class ObjectCache {
 private:
  // Only allow ObjectCache to construct CachedObjects.
  struct ConstructionKey {};

 public:
  // A reference to a cached object. Will keep the object alive (pinned) as long
  // as it lives. Provides access to the object and ID.
  // This type is not movable, assignable or copyable. It is intended to be
  // contained in a shared_ptr.
  template <typename T>
  class CachedObject {
   public:
    // Ctor. Only accessible through ObjectCache.
    CachedObject(absl::string_view id, const T& value, const std::any* any,
                 const ObjectCache* cache, ConstructionKey);

    // Dtor. Will make sure the object is un-pinned in the underlying cache.
    ~CachedObject();

    // The object ID.
    const std::string& Id() const;

    // The actual object.
    const T& Value() const;

    // Copy / move construction and assignment are disallowed (since this object
    // owns exactly one cache pin).
    CachedObject(const CachedObject&) = delete;
    CachedObject& operator=(const CachedObject&) = delete;

   private:
    std::string id_;
    const T& value_;
    const std::any* const any_;
    const ObjectCache* const cache_;
  };

  // Ctor. Specifies maximum total size of live objects in cache (may be
  // exceeded only in the event that objects are forced to be kept alive via
  // keeping references to them).
  explicit ObjectCache(size_t cache_size);

  // Dtor. Client must make sure all object references are deleted before
  // deleting the cache.
  ~ObjectCache();

  // Inserts a new object to cache, returns a reference to it.
  // If an item with the key already exists in cache, returns a reference to
  // the existing item and discards the given one.
  template <typename T>
  std::unique_ptr<const CachedObject<T>> InsertIfNotPresent(
      absl::string_view id, std::unique_ptr<const T> obj,
      size_t size = sizeof(T));

  // Looks up object by ID, return a reference to it or nullptr if not found.
  template <typename T>
  std::unique_ptr<const CachedObject<T>> Lookup(absl::string_view id) const;

 private:
  // We subclass the LRUCache so that we can add our own custom printing
  // mechanism for debug purposes.
  class ObjectStoreLruCache
      : public intrinsic::LruCache<std::string, const std::any> {
   public:
    explicit ObjectStoreLruCache(size_t cache_size);
  };
  mutable std::recursive_mutex mutex_;
  // Required to be mutable in order to be able to be able to pretend that
  // Lookup() is a const operation (technically it is not, since it would modify
  // the pin count).
  mutable ObjectStoreLruCache lru_cache_;

  // Thread-safe wrapper around the LRU release operation.
  void SafeRelease(absl::string_view id, const std::any* value) const;
};

////////////////////////////////////////////////////////////////////////////////
// Implementation detail below.

template <typename T>
ObjectCache::CachedObject<T>::CachedObject(absl::string_view id, const T& value,
                                           const std::any* any,
                                           const ObjectCache* cache,
                                           ConstructionKey)
    : id_(id), value_(value), any_(any), cache_(cache) {}

template <typename T>
ObjectCache::CachedObject<T>::~CachedObject() {
  cache_->SafeRelease(id_, any_);
}

template <typename T>
const std::string& ObjectCache::CachedObject<T>::Id() const {
  return id_;
}

template <typename T>
const T& ObjectCache::CachedObject<T>::Value() const {
  return value_;
}

template <typename T>
std::unique_ptr<const ObjectCache::CachedObject<T>>
ObjectCache::InsertIfNotPresent(absl::string_view id,
                                std::unique_ptr<const T> obj, size_t size) {
  const std::any* any;
  {
    const std::string key(id);
    std::lock_guard<std::recursive_mutex> lock(mutex_);
    any = lru_cache_.lookup(key);
    if (!any) {
      any = new const std::any(std::shared_ptr<const T>(std::move(obj)));
      lru_cache_.insertPinned(key, any, size);
      lru_cache_.MarkAdded(key);
    }
  }
  auto ptr = std::any_cast<const std::shared_ptr<const T>>(*any);
  const T& value = *ptr;
  return std::make_unique<const CachedObject<T>>(id, value, any, this,
                                                 ConstructionKey());
}

template <typename T>
std::unique_ptr<const ObjectCache::CachedObject<T>> ObjectCache::Lookup(
    absl::string_view id) const {
  const std::any* any;
  {
    std::lock_guard<std::recursive_mutex> lock(mutex_);
    any = lru_cache_.lookup(std::string(id));
  }
  if (!any) {
    return nullptr;
  }
  auto ptr = std::any_cast<const std::shared_ptr<const T>>(*any);
  const T& value = *ptr;
  return std::make_unique<const CachedObject<T>>(id, value, any, this,
                                                 ConstructionKey());
}

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_OBJECT_STORE_OBJECT_CACHE_H_
