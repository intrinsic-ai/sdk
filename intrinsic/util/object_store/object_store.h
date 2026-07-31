// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_OBJECT_STORE_OBJECT_STORE_H_
#define INTRINSIC_UTIL_OBJECT_STORE_OBJECT_STORE_H_

#include <memory>
#include <string>
#include <utility>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/marshal/riegeli_coder.h"
#include "intrinsic/util/object_store/object_ref.h"
#include "intrinsic/util/object_store/object_store_internal.h"
#include "intrinsic/util/status/status_macros.h"
#include "riegeli/records/record_reader.h"
#include "riegeli/records/record_writer.h"

namespace intrinsic {

// A set of functions for working with a cache of immutable objects of arbitrary
// types.
//
// Applications can use an object store to share and persist immutable objects.
// Additional additions of the same object from the store would de-duplicate the
// values and return the same object as the first call.
//
// Stored objects are required to be serializable. Specifically, they have to
// have a Coder<T> concept for their type (see intrinsic/marshal/coder.h),
//
// Object store users hold objects via an ObjectRef<T>, which behaves much
// like a shared_ptr<T>. An in-memory LRU cache is used to keep objects alive in
// memory for quick retrieval even when they are not being referenced.
//
// Basic usage:
//
// ObjectRef<float> my_float = DeDuplicate(3.14f);
// EXPECT_EQ(my_float.Value(), 3.14f);
//
// ObjectRef<float> fetched = DeDuplicate(3.14f);
// EXPECT_EQ(fetched.Value(), 3.14f);
//
// // We also expect that the two ObjectRef instances share their memory.
// EXPECT_EQ(&my_float.Value(), &fetched.Value());
//

// Puts an object in the object store and returns a reference to it. If an
// object with the same hashed value already exists in the store, that object
// will be returned instead.
template <typename T>
ObjectRef<T> DeDuplicate(std::unique_ptr<const T> obj) {
  return object_store_internal::GlobalObjectStore().Put(std::move(obj));
}

// Puts an object in the object store and returns a reference to it. If an
// object with the same hashed value already exists in the store, that object
// will be returned instead.
template <typename T>
ObjectRef<T> DeDuplicate(T&& obj) {
  return object_store_internal::GlobalObjectStore().Put(std::forward<T>(obj));
}

// Specializations of serializer / deserializer for ObjectRef type.
template <typename T>
struct RiegeliCoder<intrinsic::ObjectRef<T>> {
  static std::string TypeName() { return Demangle<intrinsic::ObjectRef<T>>(); }
  static absl::Status Encode(const intrinsic::ObjectRef<T>& obj,
                             riegeli::RecordWriterBase& writer) {
    return RiegeliCoder<T>::Encode(obj.Value(), writer);
  }
  static absl::StatusOr<intrinsic::ObjectRef<T>> Decode(
      riegeli::RecordReaderBase& reader) {
    INTR_ASSIGN_OR_RETURN(auto value, RiegeliCoder<T>::Decode(reader));
    return intrinsic::DeDuplicate<T>(std::move(value));
  }
};

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_OBJECT_STORE_OBJECT_STORE_H_
