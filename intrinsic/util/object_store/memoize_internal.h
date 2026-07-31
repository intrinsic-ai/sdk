// Copyright 2023 Intrinsic Innovation LLC

// This file provides a helper function for working with the object store
// memoization function. It provides an automated way of converting a function
// and its arguments into a key to be used with the object store memoization
// calls.
//
// See the Memoize function documentation for further details on usage.
#ifndef INTRINSIC_UTIL_OBJECT_STORE_MEMOIZE_INTERNAL_H_
#define INTRINSIC_UTIL_OBJECT_STORE_MEMOIZE_INTERNAL_H_

#include <memory>
#include <type_traits>
#include <utility>

#include "absl/status/statusor.h"
#include "intrinsic/util/object_store/object_ref.h"
#include "intrinsic/util/object_store/object_store.h"

namespace intrinsic {
namespace object_store_internal {

// Some helper definitions for defining the Memoization functions return type.
// This has to be at the top of the file because its evaluated for the return
// argument of the call below.
template <typename F, typename... Args>
struct MemoizeReturnTypeHelper {
  template <typename T>
  struct GetBaseType {
    using Type = ObjectRef<T>;
  };

  template <typename T>
  struct GetBaseType<std::unique_ptr<const T>> {
    using Type = ObjectRef<T>;
  };

  template <typename T>
  struct GetBaseType<std::unique_ptr<T>> {
    using Type = ObjectRef<T>;
  };

  template <typename T>
  struct GetBaseType<ObjectRef<T>> {
    using Type = ObjectRef<T>;
  };

  template <typename T>
  struct GetBaseType<absl::StatusOr<T>> {
    using Type = absl::StatusOr<typename GetBaseType<T>::Type>;
  };

  // The return type is the result of F after removing any unique pointers and
  // decaying const/ref from the type. Ex: T -> T, unique_ptr<const T> -> T,
  // ObjectRef<T> -> T.
  using Type = typename GetBaseType<std::invoke_result_t<F, Args...>>::Type;
};

template <typename F, typename... Args>
using MemoizeReturnType = typename MemoizeReturnTypeHelper<F, Args...>::Type;

template <typename>
struct IsStatusOr : std::false_type {};

template <typename T>
struct IsStatusOr<absl::StatusOr<T>> : std::true_type {};

// Converts a different kinds of value types to ObjectRefs.
// Supports T&&, unique_ptr<const T> and ObjectRef<T>. In the two former
// cases, the object will be put in the object store. In the latter, it is
// already there, so conversion is a no-op.
template <typename T>
ObjectRef<T> ToRef(T&& val) {
  return DeDuplicate(std::forward<T>(val));
}

template <typename T>
ObjectRef<T> ToRef(std::unique_ptr<const T> val) {
  return DeDuplicate(std::move(val));
}

template <typename T>
ObjectRef<T> ToRef(std::unique_ptr<T> val) {
  return DeDuplicate(std::unique_ptr<const T>(std::move(val)));
}

template <typename T>
ObjectRef<T> ToRef(ObjectRef<T> val) {
  return val;
}

}  // namespace object_store_internal
}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_OBJECT_STORE_MEMOIZE_INTERNAL_H_
