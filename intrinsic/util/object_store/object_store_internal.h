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

#ifndef INTRINSIC_UTIL_OBJECT_STORE_OBJECT_STORE_INTERNAL_H_
#define INTRINSIC_UTIL_OBJECT_STORE_OBJECT_STORE_INTERNAL_H_

#include <cstddef>
#include <cstdint>
#include <functional>
#include <initializer_list>
#include <memory>
#include <optional>
#include <string>
#include <type_traits>
#include <typeinfo>
#include <utility>

#include "absl/log/check.h"
#include "absl/meta/type_traits.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "intrinsic/marshal/riegeli_coder.h"
#include "intrinsic/util/object_store/multi_mutex.h"
#include "intrinsic/util/object_store/object_cache.h"
#include "intrinsic/util/object_store/object_ref.h"
#include "riegeli/base/maker.h"
#include "riegeli/bytes/string_writer.h"
#include "riegeli/records/record_writer.h"

namespace intrinsic {
namespace object_store_internal {

template <typename>
struct ObjectStoreIsStatusOr : std::false_type {};

template <typename T>
struct ObjectStoreIsStatusOr<absl::StatusOr<T>> : std::true_type {};

// Generate a (almost certainly) unique ID for an object, given a list of
// strings that uniquely identify it.
std::string GenerateId(std::initializer_list<absl::string_view> args);

// Generate a (almost certainly) unique ID for an object, given a list of
// uint64 that uniquely identify it.
std::string GenerateId(std::initializer_list<uint64_t> args);

// Helper for getting the global instance of the ObjectStore. Used by the
// external facing helpers for memoize and deduplicate.
ObjectStore& GlobalObjectStore();

// A persistent, cached store for immutable objects, in the form <id, obj>.
//
// Objects are put into the store by calling Put(). It is assumed that the
// <id> key of that object is unique and will always map to the
// same object for the lifetime of the store. It is, however, valid to try to
// put *the same* object into the store with the same key.
//
// The object store maintains an in-memory object cache, so recently used
// objects (via Put()) are likely to be cheaply accessible.
//
// Objects residing in the ObjectStore are owned by the store and referencing
// them is achieved by holding on to lightweight ObjectRef<T> objects, which
// behave very similarly to shared_ptr<const T>, with the additional
// functionality of keeping track of the object's key (ID) and of keeping the
// object pinned to cache as long as the reference is alive.
//
// This class is thread-safe.
class ObjectStore {
 public:
  // Remove CV and ref, as in C++20 remove_cvref_t.
  template <typename T>
  using remove_cvref = absl::remove_cv_t<absl::remove_reference_t<T>>;

  ObjectStore();

  // Puts an object in the object store and returns a reference to it. If an
  // object with the same hashed value already exists in the store, that object
  // will be returned instead.
  template <typename T>
  ObjectRef<T> Put(std::unique_ptr<const T> obj);

  // Puts an object in the object store and returns a reference to it. The
  // ID will be internally assigned and is a function of the serialized value
  // and the Coder type.
  template <typename T>
  ObjectRef<T> Put(T&& obj);

  template <typename T>
  ObjectRef<ObjectStore::remove_cvref<T>> GetOrGenerate(
      absl::string_view cache_key,
      const std::function<std::unique_ptr<const T>()>& generator,
      bool save_memoize_on_cancellation);

  template <typename T>
  ObjectRef<ObjectStore::remove_cvref<T>> GetOrGenerate(
      absl::string_view cache_key, const std::function<T()>& generator,
      bool save_memoize_on_cancellation);

 private:
  template <typename T>
  ObjectRef<T> PutInCache(absl::string_view cache_key,
                          std::unique_ptr<const T> obj,
                          size_t size = sizeof(T));

  template <typename T>
  std::optional<ObjectRef<T>> GetFromCache(absl::string_view cache_key);

 private:
  // Object cache.
  ObjectCache cache_;
  // A multi-mutex used to coordinate multi-threaded access to the same object.
  // We lock every object using its unique cache key on every access, which
  // helps us avoid multiple saves or loads of the same object in multi-threaded
  // scenarios.
  MultiMutex object_mutexes_;
};

// Specializations of serializer / deserializer for ObjectRef types.
template <typename T>
ObjectRef<T> ObjectStore::Put(std::unique_ptr<const T> obj) {
  static constexpr absl::string_view kObjectExtraCacheKey = "__objects__";

  std::string serialized;
  riegeli::RecordWriter writer(
      riegeli::Maker<riegeli::StringWriter>(&serialized),
      riegeli::RecordWriterBase::Options().set_uncompressed());
  CHECK_OK(RiegeliCoder<T>::Encode(*obj, writer));
  CHECK(writer.Close());

  const auto cache_key = object_store_internal::GenerateId(
      {kObjectExtraCacheKey, typeid(T).name(), RiegeliCoder<T>::TypeName(),
       serialized});

  auto lock = object_mutexes_.Acquire(cache_key);

  // First, check if the object is already in a cache. This is a simple
  // heuristic to avoid redundant saves in case an object gets Put() multiple
  // times.
  std::optional<ObjectRef<T>> from_cache = GetFromCache<T>(cache_key);
  if (from_cache.has_value()) {
    return std::move(from_cache).value();
  }

  return PutInCache(cache_key, std::move(obj), serialized.size());
}

template <typename T>
ObjectRef<T> ObjectStore::Put(T&& obj) {
  return ObjectStore::Put(std::make_unique<const T>(std::forward<T>(obj)));
}

template <typename T>
ObjectRef<ObjectStore::remove_cvref<T>> ObjectStore::GetOrGenerate(
    absl::string_view cache_key,
    const std::function<std::unique_ptr<const T>()>& generator,
    const bool save_memoize_on_cancellation) {
  auto lock = object_mutexes_.Acquire(cache_key);

  // First try in cache.
  std::optional<ObjectRef<T>> from_cache = GetFromCache<T>(cache_key);
  if (from_cache.has_value()) {
    return std::move(from_cache).value();
  }

  // Finally, generate, cache and save.
  std::unique_ptr<const T> obj = generator();
  CHECK(obj);

  // Handles cancelled result for StatusOr<T>.
  if constexpr (ObjectStoreIsStatusOr<T>::value) {
    if (!save_memoize_on_cancellation &&
        obj->status().code() == absl::StatusCode::kCancelled) {
      // For a cancelled result, creates an ObjectRef that holds the result
      // but is not findable via the original cache key. This is done by using
      // a unique key that includes the address of the result object.
      // The returned ObjectRef will remain valid and keep the object alive, but
      // subsequent lookups for the original `cache_key` will not find this
      // entry. Ideally, this entry from the cache should be removed after the
      // caller has destroyed the returned ObjectRef.
      const std::string unique_key = absl::StrCat(
          cache_key, "/cancelled/", reinterpret_cast<uintptr_t>(obj.get()));
      auto lock2 = object_mutexes_.Acquire(unique_key);
      return PutInCache(unique_key, std::move(obj));
    }
  }

  return PutInCache(cache_key, std::move(obj));
}

template <typename T>
ObjectRef<ObjectStore::remove_cvref<T>> ObjectStore::GetOrGenerate(
    absl::string_view cache_key, const std::function<T()>& generator,
    const bool save_memoize_on_cancellation) {
  return GetOrGenerate<T>(
      cache_key,
      [&generator]() {
        return std::make_unique<const remove_cvref<T>>(generator());
      },
      save_memoize_on_cancellation);
}

template <typename T>
ObjectRef<T> ObjectStore::PutInCache(absl::string_view cache_key,
                                     std::unique_ptr<const T> obj,
                                     size_t size) {
  return ObjectRef<T>(
      cache_key, cache_.InsertIfNotPresent(cache_key, std::move(obj), size));
}

template <typename T>
std::optional<ObjectRef<T>> ObjectStore::GetFromCache(
    absl::string_view cache_key) {
  std::unique_ptr<const ObjectCache::CachedObject<T>> obj =
      cache_.Lookup<T>(cache_key);
  if (!obj) {
    return std::nullopt;
  }
  return ObjectRef<T>(cache_key, std::move(obj));
}

}  // namespace  object_store_internal
}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_OBJECT_STORE_OBJECT_STORE_INTERNAL_H_
