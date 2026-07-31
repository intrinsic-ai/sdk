// Copyright 2023 Intrinsic Innovation LLC

// This file provides a helper function for working with the object store
// memoization function. It provides an automated way of converting a function
// and its arguments into a key to be used with the object store memoization
// calls.
//
// See the Memoize function documentation for further details on usage.
#ifndef INTRINSIC_UTIL_OBJECT_STORE_MEMOIZE_H_
#define INTRINSIC_UTIL_OBJECT_STORE_MEMOIZE_H_

#include <functional>
#include <string>
#include <utility>

#include "absl/log/check.h"
#include "absl/meta/type_traits.h"
#include "absl/strings/string_view.h"
#include "intrinsic/marshal/riegeli_coder.h"
#include "intrinsic/util/hash.h"
#include "intrinsic/util/object_store/memoize_fingerprint.h"
#include "intrinsic/util/object_store/memoize_internal.h"
#include "intrinsic/util/object_store/object_store_internal.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {

// Options to control the behavior of Memoize.
struct MemoizeOptions {
  // The name of the function to memoize.
  std::string fn_name;

  // If true, the memoization will be saved in the object store even if the
  // function returns a cancelled error.
  bool save_memoize_on_cancellation = false;

  // Deletes the default constructor to require function name.
  MemoizeOptions() = delete;

  explicit MemoizeOptions(absl::string_view fn_name)
      : fn_name(fn_name), save_memoize_on_cancellation(false) {}

  // Enables saving the memoization in the object store even if the function
  // returns a cancelled error.
  MemoizeOptions& WithSaveMemoizeOnCancellation() {
    this->save_memoize_on_cancellation = true;
    return *this;
  }
};

// Memoize the result of calling F with the given args.
//
// function_name should be a name unique to the object store for the given
// function. This name is mixed into the key such that all memoization calls
// for different functions with matching parameters will not match. (This can
// also be used to version functions if the output from the same arguments
// will change the result).
//
// F must return either:
// - std::unique_ptr<T>
// - std::unique_ptr<const T>
// - T, where T is move-constructible.
// - ObjectRef<T>.
//
// Each passed argument must be passable to memoize::MemoizeFingerprint<T>.
// There is a default implementation provided for string types and for types
// that support Coder<T>.
//
// The result is an object ref for the base type returned by F.
//
// If an object made with the same function name and args was previously
// inserted to the store, the store returns it without calling the provided
// function. Otherwise, calls the provided function in order to produce the
// object and puts it in the store.
//
// If an object was previously stored but fails to parse correctly we will call
// the generation function F anyway and override the value.
//
// The provided function is expected to have "const" semantics (i.e. its
// output is always the same and it produces no side-effects). If called from
// multiple threads with the same key at about the same time, there is no
// guarantee that the function will not be invoked more than once, but
// behavior will still remain correct.
//
// Memoization example:
//
// SomeObject VeryExpensiveFunction(int arg1, float arg2) { /* ... */ }
//
// ObjectRef<SomeObject> obj =
// Memoize(MemoizeOptions("VeryExpensiveFunction"),
//       &VeryExpensiveFunction, 15, 16);   // Takes a long time
//
// ObjectRef<SomeObject> obj2 =
// Memoize(MemoizeOptions("VeryExpensiveFunction"),
//       &VeryExpensiveFunction, 15, 16);   // Super-fast.
//
// The above example demonstrates using "normal" key generation for
// memoization, which requires that the arguments can be passed to
// memoize::MemoizeFingerprint.
//
// By default, the memoization will *not* be saved in the object store if the
// function returns absl::CancelledError(). This enables the caller to use
// Memoize with cancellation without caching the cancelled result, which is not
// desirable behavior.
//
// Use `WithSaveMemoizeOnCancellation()` to opt-out of this behavior.
// Example:
//
// INTR_ASSIGN_OR_RETURN(
//     ObjectRef<SomeObject> obj,
//     Memoize(
//         MemoizeOptions("ExpensiveFunction").WithSaveMemoizeOnCancellation(),
//     &ExpensiveFunction, stop_token, 15));
//
// Note that the `StopToken` is not included in the memoization key.
//
template <typename F, typename... Args>
object_store_internal::MemoizeReturnType<F, Args...> Memoize(
    const MemoizeOptions& options, const F& func, Args&&... args) {
  CHECK(!options.fn_name.empty()) << "Function name must be set.";
  static constexpr absl::string_view kMemoizeExtraCacheKey = "__memoize__";
  using T = object_store_internal::MemoizeReturnType<F, Args...>;

  // Here we move create the result using an indirect helper, which means that F
  // is allowed to return anything that can be passed as a constructor to T, or
  // a unique_ptr<const T>
  std::function<T()> func_to_call = [&]() -> T {
    if constexpr (object_store_internal::IsStatusOr<T>::value) {
      INTR_ASSIGN_OR_RETURN(auto val, func(std::forward<Args>(args)...));
      return object_store_internal::ToRef(std::move(val));
    } else {
      return object_store_internal::ToRef(func(std::forward<Args>(args)...));
    }
  };

  // The memoization key is generated from
  // - extra cache key, to distinguish from normal Put calls.
  // - the function name
  // - function argument types
  // - function argument fingerprints
  const std::string cache_key = object_store_internal::GenerateId({
      Fingerprint(kMemoizeExtraCacheKey),
      Fingerprint(options.fn_name),
      Fingerprint(RiegeliCoder<absl::decay_t<Args>>::TypeName())...,
      memoize::MemoizeFingerprint<absl::decay_t<Args>>::Fingerprint(args)...,
  });

  return object_store_internal::GlobalObjectStore()
      .GetOrGenerate(cache_key, func_to_call,
                     options.save_memoize_on_cancellation)
      .Value();
}

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_OBJECT_STORE_MEMOIZE_H_
