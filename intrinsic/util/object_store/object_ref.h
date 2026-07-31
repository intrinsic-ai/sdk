// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_OBJECT_STORE_OBJECT_REF_H_
#define INTRINSIC_UTIL_OBJECT_STORE_OBJECT_REF_H_

#include <cstdint>
#include <memory>
#include <string>
#include <utility>

#include "absl/strings/string_view.h"
#include "intrinsic/util/hash.h"
#include "intrinsic/util/object_store/memoize_fingerprint.h"
#include "intrinsic/util/object_store/object_cache.h"

namespace intrinsic {

namespace object_store_internal {
class ObjectStore;  // Used for forward declaration of friendship
}  // namespace  object_store_internal

// A reference to an object residing in the ObjectStore.
//
// This is a lightweight object, which behaves very similarly to
// shared_ptr<const T>, with the additional functionality of keeping track of
// the object's key (ID) and of keeping the object pinned to cache as long as
// the reference is alive.
template <typename T>
class ObjectRef {
 public:
  // The underlying object. Valid as long as the ObjectRef is alive.
  const T& Value() const { return data_->obj->Value(); }

  template <typename H>
  friend H AbslHashValue(H h, const ObjectRef<T>& obj) {
    return H::combine(std::move(h), obj.Id());
  }

 private:
  template <typename T2>
  friend bool operator==(const ObjectRef<T2>& a, const ObjectRef<T2>& b);

  template <typename T2>
  friend bool operator<(const ObjectRef<T2>& a, const ObjectRef<T2>& b);

  friend struct memoize::MemoizeFingerprint<ObjectRef<T>>;

  // The ID of the object. Valid as long as the ObjectRef is alive.
  const std::string& Id() const { return data_->id; }

  // We're putting all the actual members of the object in a shared_ptr in order
  // to make copying cheap.
  struct Data {
    std::unique_ptr<const ObjectCache::CachedObject<T>> obj;
    std::string id;
  };

  friend class object_store_internal::ObjectStore;
  ObjectRef(absl::string_view id,
            std::unique_ptr<const ObjectCache::CachedObject<T>> obj)
      : data_(new Data{std::move(obj), std::string(id)}) {}

  std::shared_ptr<const Data> data_;
};

template <typename T>
bool operator==(const ObjectRef<T>& a, const ObjectRef<T>& b) {
  return a.Id() == b.Id();
}

template <typename T>
bool operator!=(const ObjectRef<T>& a, const ObjectRef<T>& b) {
  return !(a == b);
}

// Defined so that ObjectRef can be placed into a set.
template <typename T>
bool operator<(const ObjectRef<T>& a, const ObjectRef<T>& b) {
  return a.Id() < b.Id();
}

// Specialization for fingerprinting ObjectRef.
namespace memoize {
template <typename T>
struct MemoizeFingerprint<ObjectRef<T>> {
  static uint64_t Fingerprint(const ObjectRef<T>& value) {
    return intrinsic::Fingerprint(value.Id());
  }
};
}  // namespace memoize

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_OBJECT_STORE_OBJECT_REF_H_
