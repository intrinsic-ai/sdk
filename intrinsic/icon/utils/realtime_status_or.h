// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_UTILS_REALTIME_STATUS_OR_H_
#define INTRINSIC_ICON_UTILS_REALTIME_STATUS_OR_H_

#include <type_traits>
#include <utility>
#include <variant>

#include "absl/base/attributes.h"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "intrinsic/icon/testing/realtime_annotations.h"
#include "intrinsic/icon/utils/realtime_status.h"

namespace intrinsic::icon {

// A variant of absl::StatusOr for realtime contexts.
//
// It contains either a usable object of type 'T' or an error of type
// RealtimeStatus that explains why the object is missing. Assuming that type
// 'T' is realtime-safe, this class is also, because RealtimeStatus does not
// allocate on the heap.
//
// Compared to absl::StatusOr, this here has fewer features.
// Among other things, it does not warn about undesirable corner cases like
// constructing with kOk, or ambiguous construction with "{}".
template <typename T>
class RealtimeStatusOr {
 private:
  static_assert(!std::is_rvalue_reference_v<T>,
                "RealtimeStatusOr<T> does not support rvalue reference types.");
  using StoredT = std::remove_reference_t<T>;
  using VariantT = std::variant<RealtimeStatus, StoredT>;
  static constexpr int kStatusIndex = 0;
  static constexpr int kDataIndex = 1;
  static constexpr RealtimeStatus kOkStatus = OkStatus();

 public:
  // Allows passing a usable object. Object type must be an r value, if it is
  // not trivially copyable.
  //
  // Trivially copyable example:
  //
  //    RealtimeStatusOr<int> GetRandomNumber() {
  //      return 4;
  //    }
  //
  // Not trivially copyable example:
  //
  //    RealtimeStatusOr<std::unique_ptr<Foo>> GetRandomNumber() {
  //      std::unique_ptr<Foo> bar = ...;
  //      return std::move(bar);
  //    }
  //
  // std::conditional templating uses the copy constructor T if the type is
  // trivially copyable, otherwise it uses T&&.
  // Checking for trivially copyable of StoredT, because T can be a reference,
  // of a non-trivially copyable type, but would still make IsCopyable true.
  template <bool IsCopyable = std::is_trivially_copy_constructible_v<StoredT>>
  RealtimeStatusOr(  // NOLINT(google-explicit-constructor)
      typename std::conditional<IsCopyable, StoredT, StoredT&&>::type data)
      : status_or_data_(std::in_place_index<kDataIndex>, std::move(data)) {}

  // Allows returning errors.
  // Example:
  //   RealtimeStatusOr<int> GetRandomNumber() {
  //     return UnimplementedError("GetRandomNumber not implemented");
  //   }
  RealtimeStatusOr(  // NOLINT(google-explicit-constructor)
      RealtimeStatus&& status)
      : status_or_data_(std::in_place_index<kStatusIndex>, std::move(status)) {}

  // Allows forwarding errors from a different RealtimeStatusOr type, for
  // instance in nested functions.
  RealtimeStatusOr(  // NOLINT(google-explicit-constructor)
      const RealtimeStatus& status)
      : status_or_data_(std::in_place_index<kStatusIndex>, status) {}

  // Allowed so users can create containers.
  RealtimeStatusOr()
      : status_or_data_(std::in_place_index<kStatusIndex>,
                        RealtimeStatus(absl::StatusCode::kUnknown, "")) {}

  // Copy, move and assignment are allowed.
  RealtimeStatusOr(const RealtimeStatusOr&) = default;
  RealtimeStatusOr& operator=(const RealtimeStatusOr&) = default;
  RealtimeStatusOr(RealtimeStatusOr&&) noexcept = default;
  RealtimeStatusOr& operator=(RealtimeStatusOr&&) noexcept = default;

  // Returns true if the status is ok, i.e. the data is present.
  bool ok() const { return status_or_data_.index() == kDataIndex; }

  // Returns the status. If a usable object is present, the status code is
  // "kOk".
  const RealtimeStatus& status() const {
    if (ok()) {
      return kOkStatus;
    }
    return std::get<kStatusIndex>(status_or_data_);
  }

  // Get the usable object.
  // Only allowed if 'ok() == true', otherwise fails a runtime assert.
  const T& value() const&;
  T& value() &;

  // Move the usable object out.
  // Only allowed if 'ok() == true', otherwise fails a runtime assert.
  const T&& value() const&&;
  T&& value() &&;

  // DEPRECATED: Prefer accessing the value using `operator*` or `operator->`
  // after testing that the StatusOr is OK. If program termination is desired in
  // the case of an error status, consider `CHECK_OK(status_or.status());`.
  //
  // `ValueOrDie` will be removed after all usages have been migrated.
  //
  // Returns a reference to our current value, or CHECK-fails if `!this->ok()`.
  //
  // Methods exist for legacy macros only and should not be used in user code.
  ABSL_DEPRECATED("Use operator* or operator->, after testing ok or CHECK_OK.")
  const T& ValueOrDie() const&;
  ABSL_DEPRECATED("Use operator* or operator->, after testing ok or CHECK_OK.")
  T& ValueOrDie() &;
  ABSL_DEPRECATED("Use operator* or operator->, after testing ok or CHECK_OK.")
  const T&& ValueOrDie() const&&;
  ABSL_DEPRECATED("Use operator* or operator->, after testing ok or CHECK_OK.")
  T&& ValueOrDie() &&;

  // Returns a reference to the current value.
  //
  // REQUIRES: `this->ok() == true`, otherwise the behavior is undefined.
  //
  // Use `this->ok()` to verify that there is a current value within the
  // `RealtimeStatusOr<T>`. Alternatively, see the `value()` member function for
  // a similar API that guarantees crashing or throwing an exception if there is
  // no current value.
  const T& operator*() const&;
  T& operator*() &;
  const T&& operator*() const&&;
  T&& operator*() &&;

 private:
  std::variant<RealtimeStatus, StoredT> status_or_data_;
};

template <typename T>
const T& RealtimeStatusOr<T>::value() const& INTRINSIC_SUPPRESS_REALTIME_CHECK {
  CHECK(ok()) << "RealtimeStatusOr::value() only allowed if ok() aka usable "
                 "value has been set";
  return std::get<kDataIndex>(status_or_data_);
}

template <typename T>
    T& RealtimeStatusOr<T>::value() & INTRINSIC_SUPPRESS_REALTIME_CHECK {
  CHECK(ok()) << "RealtimeStatusOr::value() only allowed if ok() aka usable "
                 "value has been set";
  return std::get<kDataIndex>(status_or_data_);
}

template <typename T>
const T&& RealtimeStatusOr<T>::value()
    const&& INTRINSIC_SUPPRESS_REALTIME_CHECK {
  CHECK(ok()) << "RealtimeStatusOr::value() only allowed if ok() aka usable "
                 "value has been set";
  return std::move(std::get<kDataIndex>(status_or_data_));
}

template <typename T>
    T&& RealtimeStatusOr<T>::value() && INTRINSIC_SUPPRESS_REALTIME_CHECK {
  CHECK(ok()) << "RealtimeStatusOr::value() only allowed if ok() aka usable "
                 "value has been set";
  return std::move(std::get<kDataIndex>(status_or_data_));
}

template <typename T>
const T& RealtimeStatusOr<T>::ValueOrDie() const& {
  CHECK(ok())
      << "RealtimeStatusOr::ValueOrDie() only allowed if ok() aka usable "
         "value has been set";
  return std::get<kDataIndex>(status_or_data_);
}

template <typename T>
T& RealtimeStatusOr<T>::ValueOrDie() & {
  CHECK(ok())
      << "RealtimeStatusOr::ValueOrDie() only allowed if ok() aka usable "
         "value has been set";
  return std::get<kDataIndex>(status_or_data_);
}

template <typename T>
const T&& RealtimeStatusOr<T>::ValueOrDie() const&& {
  CHECK(ok())
      << "RealtimeStatusOr::ValueOrDie() only allowed if ok() aka usable "
         "value has been set";
  return std::move(std::get<kDataIndex>(status_or_data_));
}

template <typename T>
T&& RealtimeStatusOr<T>::ValueOrDie() && {
  CHECK(ok())
      << "RealtimeStatusOr::ValueOrDie() only allowed if ok() aka usable "
         "value has been set";
  return std::move(std::get<kDataIndex>(status_or_data_));
}

template <typename T>
const T& RealtimeStatusOr<T>::operator*() const& {
  CHECK(ok())
      << "RealtimeStatusOr::operator*() only allowed if ok() aka usable "
         "value has been set";
  return std::get<kDataIndex>(status_or_data_);
}

template <typename T>
T& RealtimeStatusOr<T>::operator*() & {
  CHECK(ok())
      << "RealtimeStatusOr::operator*() only allowed if ok() aka usable "
         "value has been set";
  return std::get<kDataIndex>(status_or_data_);
}

template <typename T>
const T&& RealtimeStatusOr<T>::operator*() const&& {
  CHECK(ok())
      << "RealtimeStatusOr::operator*() only allowed if ok() aka usable "
         "value has been set";
  return std::move(std::get<kDataIndex>(status_or_data_));
}

template <typename T>
T&& RealtimeStatusOr<T>::operator*() && {
  CHECK(ok())
      << "RealtimeStatusOr::operator*() only allowed if ok() aka usable "
         "value has been set";
  return std::move(std::get<kDataIndex>(status_or_data_));
}

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_UTILS_REALTIME_STATUS_OR_H_
