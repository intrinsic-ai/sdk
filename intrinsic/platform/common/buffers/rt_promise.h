// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_PLATFORM_COMMON_BUFFERS_RT_PROMISE_H_
#define INTRINSIC_PLATFORM_COMMON_BUFFERS_RT_PROMISE_H_

#include <atomic>
#include <memory>
#include <optional>
#include <utility>

#include "absl/base/thread_annotations.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "intrinsic/icon/interprocess/binary_futex.h"
#include "intrinsic/icon/interprocess/lockable_binary_futex.h"
#include "intrinsic/icon/testing/realtime_annotations.h"
#include "intrinsic/icon/utils/log.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_macro.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/platform/common/buffers/rt_queue.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {

// Implementation of future and promise for single-producer, single-consumer
// message passing (between threads). Use, i.e., if a non-rt thread is waiting
// for an rt thread to eventually generate a value.
//
// Example:
//
// RealtimeFuture<bool> rt_job_result;
// Thread rt_thread;
// ThreadOptions rt_thread_options;
// ASSIGN_OR_RETURN(RealtimePromise<bool> promise, rt_job_result.GetPromise());
// RETURN_IF_ERROR(rt_thread.Start(rt_thread_options,
//                           [promise = std::move(promise)]() mutable {
//                             // Do fancy real time stuff that'll set the value
//                             // of `result`.
//                             // ...
//                             auto status = promise.Set(result);
//                           }));
// // Wait for the `rt_thread` to set the value on the promise.
// ASSIGN_OR_RETURN(bool job_result,
//   rt_job_result.WaitUntilAndGet(absl::InfiniteFuture()));
//
// General notes:
// * The shared data lives in `FuturePromiseContext`.
// * The promise can only be moved, but not copied.
// * The future can only be moved, but not copied.
// * The future and promise are not reentrant, nor can they be used from
//   different threads at the same time.
// * Calling `Get()`, `WaitForAndGet()`, or `WaitUntilAndGet()` consumes the
//   value of the future, enabling a new promise to set a new value.
// * Calling `Peek()`, `WaitForAndPeek()`, or `WaitUntilAndPeek()` does not
//   consume the value of the future. The value can be peeked at multiple times.
// * Calling `Set()` on a promise or `Get()` on a future will make the future or
//   promise detach from the context, allowing a new future or promise to be
//   created even while the old one is still alive.
// * To prevent the future or promise from detaching from the context, the
//   future or promise can be created with `is_reusable = true`. This will allow
//   the future or promise to be reused multiple times, i.e. `Set()` can be
//   called multiple times as long as the value is consumed by a `Get()`
//   inbetween. Similarly, `Get()` can be called multiple times as long as a new
//   value is provided by a `Set()` inbetween each `Get()`. In order to destroy
//   a reusable future or promise, you might want to assign a `DetachedFuture`
//   or `DetachedPromise` to it respectively.
// * If `std::is_copy_constructible<T>` is `false`, a copy functor needs to be
//   passed to `Peek()`, `WaitForAndPeek()`, and `WaitUntilAndPeek()`.
// * The `FuturePromiseContext` must always outlive the promise and future
//   handlesthat may still be attached to it.
// * A default-constructed `RealtimeFuture` or `RealtimePromise` will own their
//   `FuturePromiseContext`. Thus, the destruction of a `RealtimeFuture` will
//   wait for any pending promise to be destroyed, first. Similarly, the
//   destruction of a `RealtimePromise` will wait for any pending future to be
//   destroyed.
// * Users can use `IsWaitFreeDestructible()` to check whether the future or
//   promise can be destroyed without waiting for a potentially attached promise
//   or future.
// * For full control over the wait-free destruction, we recommend creating the
//   future and promise from a `FuturePromiseContext` instance, in which case
//   the destruction of the future or promise is guaranteed to be wait-free.
// * The `FuturePromiseContext` provides a `Reset` method to reset the state of
//   the context. This method will wait for any attached future or promise to be
//   destroyed first.

// Forward declaration for RealtimePromise.
template <typename T>
class RealtimePromise;

// Forward declaration for RealtimeFuture.
template <typename T>
class RealtimeFuture;

// FuturePromiseContext is the shared state between a RealtimePromise and a
// RealtimeFuture. It manages the underlying buffer, synchronization primitives,
// and state flags.
template <typename T>
class FuturePromiseContext {
 public:
  explicit FuturePromiseContext(
      absl::Duration detach_timeout = absl::Seconds(1))
      : buffer_(/*capacity=*/1, /*init_function=*/nullptr),
        // By default, neither promise nor future is attached.
        promise_attached_(false),
        future_attached_(false),
        is_cancelled_(false),
        read_value_futex_(/*private_futex=*/true),
        write_value_futex_(/*private_futex=*/true),
        cancel_futex_(/*private_futex=*/true),
        get_future_futex_(/*private_futex=*/true),
        get_promise_futex_(/*private_futex=*/true),
        has_value_(/*posted=*/false, /*private_futex=*/true),
        future_detached_(/*posted=*/false, /*private_futex=*/true),
        promise_detached_(/*posted=*/false, /*private_futex=*/true),
        detach_timeout_(detach_timeout) {};

  // Destructor waits until detach_timeout_ for attached promise and future
  // handles to be detached. This ensures that the FuturePromiseContext outlives
  // any active handles unless they are explicitly detached or their owners are
  // destroyed.
  ~FuturePromiseContext() {
    // Signal cancellation to the promise and future.
    is_cancelled_.store(true, std::memory_order_release);
    if (future_attached_.load(std::memory_order_acquire)) {
      // Wait for the future to be destroyed.
      auto status = future_detached_.WaitFor(detach_timeout_);
      INTRINSIC_RT_LOG_IF(ERROR, !status.ok())
          << "Future was not destroyed: " << status.ToString();
    }
    if (promise_attached_.load(std::memory_order_acquire)) {
      // Wait for the promise to be destroyed.
      auto status = promise_detached_.WaitFor(detach_timeout_);
      INTRINSIC_RT_LOG_IF(ERROR, !status.ok())
          << "Promise was not destroyed: " << status.ToString();
    }
  };

  FuturePromiseContext(const FuturePromiseContext&) = delete;
  FuturePromiseContext& operator=(const FuturePromiseContext&) = delete;
  FuturePromiseContext(FuturePromiseContext&&) = delete;
  FuturePromiseContext& operator=(FuturePromiseContext&&) = delete;

  // Returns a RealtimePromise associated with this context.
  // Fails if a promise is already attached or if a reset is in progress.
  // If `is_reusable` is true, the promise can be reused after setting a value.
  icon::RealtimeStatusOr<RealtimePromise<T>> GetPromise(
      bool is_reusable = false) INTRINSIC_CHECK_REALTIME_SAFE {
    if (!get_promise_futex_.TryLock()) {
      return icon::UnavailableError(
          "Cannot create promise while reset is in progress.");
    }
    bool expected_promise_attached = false;
    if (!promise_attached_.compare_exchange_strong(
            expected_promise_attached, true, std::memory_order_acq_rel,
            std::memory_order_acquire)) {
      // The relevant error here is the one about the promise already being
      // attached, so we ignore any unlock errors.
      (void)get_promise_futex_.Unlock();
      return icon::AlreadyExistsError("Promise is already attached.");
    }
    INTRINSIC_RT_RETURN_IF_ERROR(get_promise_futex_.Unlock());
    return std::move(RealtimePromise<T>(this, is_reusable));
  }

  // Returns true if a promise is attached to this context.
  bool IsPromiseAttached() const INTRINSIC_CHECK_REALTIME_SAFE {
    return promise_attached_.load(std::memory_order_acquire);
  }

  // Returns a RealtimeFuture associated with this context.
  // Fails if a future is already attached or if a reset is in progress.
  // If `is_reusable` is true, the future can be reused after getting a value.
  icon::RealtimeStatusOr<RealtimeFuture<T>> GetFuture(bool is_reusable = false)
      INTRINSIC_CHECK_REALTIME_SAFE {
    if (!get_future_futex_.TryLock()) {
      return icon::UnavailableError(
          "Cannot create future while reset is in progress.");
    }
    bool expected_future_attached = false;
    if (!future_attached_.compare_exchange_strong(
            expected_future_attached, true, std::memory_order_acq_rel,
            std::memory_order_acquire)) {
      // The relevant error here is the one about the future already being
      // attached, so we ignore any unlock errors.
      (void)get_future_futex_.Unlock();
      return icon::AlreadyExistsError("Future is already attached.");
    }
    INTRINSIC_RT_RETURN_IF_ERROR(get_future_futex_.Unlock());
    return std::move(RealtimeFuture<T>(this, is_reusable));
  }

  // Returns true if a future is attached to this context.
  bool IsFutureAttached() const INTRINSIC_CHECK_REALTIME_SAFE {
    return future_attached_.load(std::memory_order_acquire);
  }

  // Resets the state of the FuturePromiseContext, allowing it to be reused.
  // This operation cancels any pending operations, clears the stored value (if
  // any), and resets the internal state flags. The function will wait for any
  // existing future/promise to be detached first.
  // While the reset is in progress, no new future/promise handles can be
  // created and `Cancel()` will return an error.
  absl::Status Reset(absl::Duration timeout) INTRINSIC_NON_REALTIME_ONLY {
    // The locks below are used to prevent any future/promise handle to be
    // created while the reset is in progress.
    icon::BinaryFutexLock get_future_lock(&get_future_futex_);
    icon::BinaryFutexLock get_promise_lock(&get_promise_futex_);
    // The cancel lock prevents any cancel operations from happening while the
    // reset is in progress.
    icon::BinaryFutexLock cancel_lock(&cancel_futex_);

    // Wait for the future and promise to be detached.
    if (future_attached_.load(std::memory_order_acquire)) {
      INTR_RETURN_IF_ERROR(future_detached_.WaitFor(timeout))
          << "Future was not detached within timeout.";
    }
    if (promise_attached_.load(std::memory_order_acquire)) {
      INTR_RETURN_IF_ERROR(promise_detached_.WaitFor(timeout))
          << "Promise was not detached within timeout.";
    }

    // Prevent any parallel reads and writes to the buffer while we reset it,
    // which shouldn't be possible because there are no more attached handles.
    icon::BinaryFutexLock read_lock(&read_value_futex_);
    icon::BinaryFutexLock write_lock(&write_value_futex_);

    // Clear the buffer if there is any value.
    while (!buffer_.Empty()) {
      auto element = buffer_.reader()->MoveFront();
      if (element.has_value()) {
        element.reset();
        buffer_.reader()->DropFront();
      }
    }

    // Reset the cancelled flag.
    is_cancelled_.store(false, std::memory_order_release);
    return icon::OkStatus();
  }

  // Returns true if the context has been cancelled.
  bool IsCancelled() const INTRINSIC_CHECK_REALTIME_SAFE {
    return is_cancelled_.load(std::memory_order_acquire);
  }

  // Cancels the context.
  // This informs the associated future and promise that no value will be (or
  // can be) set.
  // Returns a ResourceExhaustedError if a reset is in progress.
  // Returns a CancelledError if the context was already cancelled.
  icon::RealtimeStatus Cancel() INTRINSIC_CHECK_REALTIME_SAFE {
    if (!cancel_futex_.TryLock()) {
      return icon::ResourceExhaustedError(
          "Cannot cancel while reset is in progress.");
    }

    // Signal cancellation to the promise and future.
    is_cancelled_.store(true, std::memory_order_release);
    INTRINSIC_RT_RETURN_IF_ERROR(cancel_futex_.Unlock());
    return has_value_.Post();
  }

  // Returns true if the context has a value.
  // Returns a CancelledError if the context was cancelled.
  // Note that the value returned by this function might change anytime after
  // this function returns if a future or promise is attached.
  icon::RealtimeStatusOr<bool> HasValue() const INTRINSIC_CHECK_REALTIME_SAFE {
    if (IsCancelled()) {
      return icon::CancelledError("Context was cancelled.");
    }
    return !buffer_.Empty();
  }

  // Test only method to acquire a lock for all the futexes used to test the
  // behavior of the context while a reset is in progress.
  icon::RealtimeStatus TestOnlyResetLock() INTRINSIC_NON_REALTIME_ONLY
      ABSL_EXCLUSIVE_LOCK_FUNCTION(read_value_futex_, write_value_futex_,
                                   cancel_futex_, get_future_futex_,
                                   get_promise_futex_) {
    auto status = get_future_futex_.Lock();
    status = icon::OverwriteIfNotInError(status, get_promise_futex_.Lock());
    status = icon::OverwriteIfNotInError(status, read_value_futex_.Lock());
    status = icon::OverwriteIfNotInError(status, write_value_futex_.Lock());
    status = icon::OverwriteIfNotInError(status, cancel_futex_.Lock());
    return status;
  }

  // Test only method to release a lock for all the futexes used to test the
  // behavior of the context while a reset is in progress.
  icon::RealtimeStatus TestOnlyResetUnlock() INTRINSIC_CHECK_REALTIME_SAFE
      ABSL_UNLOCK_FUNCTION(read_value_futex_, write_value_futex_, cancel_futex_,
                           get_future_futex_, get_promise_futex_) {
    auto status = cancel_futex_.Unlock();
    status = icon::OverwriteIfNotInError(status, write_value_futex_.Unlock());
    status = icon::OverwriteIfNotInError(status, read_value_futex_.Unlock());
    status = icon::OverwriteIfNotInError(status, get_promise_futex_.Unlock());
    status = icon::OverwriteIfNotInError(status, get_future_futex_.Unlock());
    return status;
  }

 private:
  // Grant access to the Detach* methods and the buffer to the promise and
  // future.
  friend class RealtimePromise<T>;
  friend class RealtimeFuture<T>;

  // Detaches the promise from this context. Called by RealtimePromise
  // destructor.
  icon::RealtimeStatus DetachPromise() INTRINSIC_CHECK_REALTIME_SAFE {
    bool previously_attached =
        promise_attached_.exchange(false, std::memory_order_acq_rel);
    if (previously_attached) {
      INTRINSIC_RT_RETURN_IF_ERROR(promise_detached_.Post());
      return icon::OkStatus();
    }
    return icon::InternalError("Promise was not attached.");
  }

  // Detaches the future from this context. Called by RealtimeFuture
  // destructor.
  icon::RealtimeStatus DetachFuture() INTRINSIC_CHECK_REALTIME_SAFE {
    bool previously_attached =
        future_attached_.exchange(false, std::memory_order_acq_rel);
    if (previously_attached) {
      INTRINSIC_RT_RETURN_IF_ERROR(future_detached_.Post());
      return icon::OkStatus();
    }
    return icon::InternalError("Future was not attached.");
  }

  bool ReadTryLock() INTRINSIC_CHECK_REALTIME_SAFE
      ABSL_EXCLUSIVE_TRYLOCK_FUNCTION(true, read_value_futex_) {
    return read_value_futex_.TryLock();
  }

  icon::RealtimeStatus ReadUnlock() INTRINSIC_CHECK_REALTIME_SAFE
      ABSL_UNLOCK_FUNCTION(read_value_futex_) {
    return read_value_futex_.Unlock();
  }

  bool WriteTryLock() INTRINSIC_CHECK_REALTIME_SAFE
      ABSL_EXCLUSIVE_TRYLOCK_FUNCTION(true, write_value_futex_) {
    return write_value_futex_.TryLock();
  }

  icon::RealtimeStatus WriteUnlock() INTRINSIC_CHECK_REALTIME_SAFE
      ABSL_UNLOCK_FUNCTION(write_value_futex_) {
    return write_value_futex_.Unlock();
  }

  // The underlying queue used to store the value.
  RealtimeQueue<T> buffer_;

  // True if a RealtimePromise handle is currently active/attached.
  std::atomic<bool> promise_attached_{false};
  // True if a RealtimeFuture handle is currently active/attached.
  std::atomic<bool> future_attached_{false};
  // True if the promise has cancelled the contract to provide a value.
  std::atomic<bool> is_cancelled_{false};

  // The futexes below are used to protect against race conditions between
  // this class's methods and the Reset() method. The goal is that the Reset()
  // method (being the non-realtime method) is the only one requiring all the
  // locks, while the other methods can use the individual locks. All other
  // members of this class are atomic or otherwise thread-safe and do not
  // require this lock. Therefore, annotating the entire class or its methods
  // with thread-safety analysis annotations is not necessary.
  icon::LockableBinaryFutex read_value_futex_;
  icon::LockableBinaryFutex write_value_futex_;
  icon::LockableBinaryFutex cancel_futex_;
  icon::LockableBinaryFutex get_future_futex_;
  icon::LockableBinaryFutex get_promise_futex_;
  // Futex to signal when a value has been set or an operation cancelled.
  icon::BinaryFutex has_value_;
  // Futex to signal when the future has been detached.
  icon::BinaryFutex future_detached_;
  // Futex to signal when the promise has been detached.
  icon::BinaryFutex promise_detached_;
  // The timeout to wait for the future and promise to be detached during
  // destruction.
  const absl::Duration detach_timeout_;
};

// RealtimePromise is the producer side of a single-producer, single-consumer
// value transfer. It allows setting a value or cancelling the operation.
// This class is designed for real-time contexts where allocations are
// restricted.
template <typename T>
class RealtimePromise {
 private:
  // Alias for the shared context type.
  using Context = FuturePromiseContext<T>;

 public:
  // Default constructor.
  // Creates an internal context object and attaches the promise to it.
  // Note that the destructor will wait for any attached future to be
  // detached before destroying the promise. This is to prevent dangling
  // references to the context by a still attached future.
  RealtimePromise() INTRINSIC_NON_REALTIME_ONLY : is_reusable_(true) {
    internal_context_ = std::make_unique<Context>();
    context_ = internal_context_.get();
    context_->promise_attached_.store(true, std::memory_order_release);
  }

  // Destructor.
  // Detaches the promise from the context and waits for any attached future
  // to be detached before destroying the promise. This is to prevent dangling
  // references to the context by a still attached future.
  ~RealtimePromise() {
    if (context_) {
      auto status = Detach();
      INTRINSIC_RT_LOG_IF(WARNING, !status.ok())
          << "Promise was not detached: " << status.ToString();
    }
    // The destructor of `internal_context_` will wait for any attached future
    // to be detached before destroying the context.
    internal_context_ = nullptr;
  }

  // Returns a promise that is not attached to any context.
  // This can be used to create a temporary promise in a real-time thread, i.e.
  // to cause the destruction of the assigned-to promise.
  // Relies on RVO of the compiler to instantiate the promise in-place.
  static RealtimePromise<T> GetDetachedPromise() INTRINSIC_CHECK_REALTIME_SAFE {
    return RealtimePromise<T>(/*context=*/nullptr, /*is_reusable=*/false);
  }

  // The promise is not copyable or copy-assignable.
  RealtimePromise(const RealtimePromise&) = delete;
  RealtimePromise& operator=(const RealtimePromise&) = delete;

  // The promise is movable and move-assignable.
  RealtimePromise(RealtimePromise&& other) noexcept
      INTRINSIC_CHECK_REALTIME_SAFE {
    // Call the move assignment operator for consistent logic.
    *this = std::move(other);
  }
  RealtimePromise& operator=(RealtimePromise&& other)
      INTRINSIC_CHECK_REALTIME_SAFE {
    if (this == &other) {
      return *this;
    }
    // Detach from the current context.
    if (context_) {
      auto status = Detach();
      INTRINSIC_RT_LOG_IF(WARNING, !status.ok())
          << "Promise was not detached: " << status.ToString();
    }
    std::swap(context_, other.context_);
    std::swap(internal_context_, other.internal_context_);
    std::swap(is_reusable_, other.is_reusable_);
    return *this;
  }

  // Sets the value of the promise to `value`. This will unblock a waiting
  // future. This operation can only be called once until the value has been
  // consumed (through any of the Get or Wait*AndGet ethods). Returns a
  // ResourceExhaustedError if the value has already been set. Returns a
  // CancelledError if the promise was cancelled. Returns a
  // FailedPreconditionError if the promise is not attached.
  icon::RealtimeStatus Set(T value) INTRINSIC_CHECK_REALTIME_SAFE {
    if (!context_) {
      return icon::FailedPreconditionError("Promise is not attached.");
    }

    if (!context_->WriteTryLock()) {
      // If the try lock fails, the context is being cancelled.
      return icon::CancelledError("The promise is being cancelled.");
    }

    if (context_->IsCancelled()) {
      INTRINSIC_RT_RETURN_IF_ERROR(context_->WriteUnlock());
      return icon::CancelledError("Promise was cancelled.");
    }

    auto writer = context_->buffer_.writer();
    if (writer == nullptr) {
      INTRINSIC_RT_RETURN_IF_ERROR(context_->WriteUnlock());
      return icon::InternalError("Received nullptr writer.");
    }
    T* element = writer->PrepareInsert();
    if (element == nullptr) {
      // No need to call `FinishInsert()` here, since the insertion failed.
      INTRINSIC_RT_RETURN_IF_ERROR(context_->WriteUnlock());
      return icon::ResourceExhaustedError("Value has already been set.");
    }
    *element = std::move(value);
    writer->FinishInsert();
    // Signal that a value has been set. This will unblock a waiting future.
    auto status = context_->has_value_.Post();
    INTRINSIC_RT_RETURN_IF_ERROR(context_->WriteUnlock());
    if (!is_reusable_) {
      status = icon::OverwriteIfNotInError(status, Detach());
    }
    return status;
  }

  // Cancels the promise.
  // This informs the associated future that no value will be provided.
  // If the promise is already cancelled, an error is returned.
  // A Cancel call won't interrupt an ongoing RealtimePromise::Set() call nor
  // an ongoing RealtimeFuture::Get() call. Both will be allowed to complete
  // successfully.
  // Returns a FailedPreconditionError if the promise is not attached.
  icon::RealtimeStatus Cancel() INTRINSIC_CHECK_REALTIME_SAFE {
    if (!context_) {
      return icon::FailedPreconditionError("Promise is not attached.");
    }
    // Call the underlying method to check for cancellation, since this has
    // all the logic to handle a concurrent `Reset` call.
    return context_->Cancel();
  }

  // Checks if the promise has been cancelled.
  // Returns a FailedPreconditionError if the promise is not attached.
  icon::RealtimeStatusOr<bool> IsCancelled() const
      INTRINSIC_CHECK_REALTIME_SAFE {
    if (!context_) {
      return icon::FailedPreconditionError("Promise is not attached.");
    }
    // Call the underlying method to check for cancellation, since this has
    // all the logic to handle a concurrent `Reset` call.
    return context_->IsCancelled();
  }

  // Checks if a value is available in the promise. Returns an error if the
  // promise is not attached or has been cancelled.
  icon::RealtimeStatusOr<bool> HasValue() const INTRINSIC_CHECK_REALTIME_SAFE {
    if (!context_) {
      return icon::FailedPreconditionError("Promise is not attached.");
    }
    return context_->HasValue();
  }

  // Gets the future associated with this promise's context.
  // Returns an error if the promise is not attached or if a future has
  // already been returned.
  icon::RealtimeStatusOr<RealtimeFuture<T>> GetFuture(bool is_reusable = false)
      INTRINSIC_CHECK_REALTIME_SAFE {
    if (!context_) {
      return icon::FailedPreconditionError("Promise is not attached.");
    }
    return context_->GetFuture(is_reusable);
  }

  // Returns true, if the promise can be destroyed without waiting for a
  // future to be detached.
  // This is always true, if the promise doesn't own the context and false
  // only if it owns the context and a future is attached. Note, that calling
  // this method right before destroying the promise doesn't guarantee a
  // wait-free destruction, since a new future might have been created in the
  // meantime.
  bool IsWaitFreeDestructible() const INTRINSIC_CHECK_REALTIME_SAFE {
    // If the promise is not attached, it is wait-free.
    if (!context_) {
      // This should never happen, but we check anyway.
      return true;
    }
    // If the promise doesn't own the context, it is wait-free destructible.
    if (internal_context_.get() == nullptr) {
      return true;
    }
    // If a future is attached, the promise is not wait-free destructible.
    return context_->future_attached_.load(std::memory_order_acquire) == false;
  }

 protected:
  // Friend declaration so that the FuturePromiseContext can access the
  // private constructor of RealtimePromise.
  friend class FuturePromiseContext<T>;

  // Protected constructor used by FuturePromiseContext to create a promise.
  explicit RealtimePromise(Context* context, bool is_reusable)
      : context_(context),
        internal_context_(nullptr),
        is_reusable_(is_reusable) {}

  icon::RealtimeStatus Detach() INTRINSIC_CHECK_REALTIME_SAFE {
    if (!context_) {
      return icon::FailedPreconditionError("Promise is not attached.");
    }
    auto status = context_->DetachPromise();
    context_ = nullptr;
    return status;
  }

  // Pointer to the shared context.
  Context* context_{nullptr};
  // The optional internal context. Exists only if default constructed, in
  // which case the context_ pointer points to the internal context. Will be
  // nullptr if created from a FuturePromiseContext or RealtimeFuture.
  std::unique_ptr<Context> internal_context_{nullptr};

  // Indicates if the promise is reusable.
  // A reusable promise won't detach from the context after its `Set()` method
  // is called and thus can be reused again to set a new value (after the
  // previous value has been consumed by a future).
  bool is_reusable_{false};
};

// RealtimeFuture is the consumer side of a single-producer, single-consumer
// value transfer. It allows getting a value (with or without waiting) or
// checking for cancellation. This class is designed for real-time contexts.
template <typename T>
class RealtimeFuture {
 private:
  // Alias for the shared context type.
  using Context = FuturePromiseContext<T>;

 public:
  // Default constructor.
  // Creates an internal context object and attaches the future to it.
  // Note that the destructor will wait for any attached promise to be
  // detached before destroying the future. This is to prevent dangling
  // references to the context by a still attached promise.
  RealtimeFuture() INTRINSIC_NON_REALTIME_ONLY : is_reusable_(true) {
    internal_context_ = std::make_unique<Context>();
    context_ = internal_context_.get();
    context_->future_attached_.store(true, std::memory_order_release);
  }

  // Destructor.
  // Detaches the future from the context and waits for any attached promise
  // to be detached before destroying the future. This is to prevent dangling
  // references to the context by a still attached promise.
  ~RealtimeFuture() {
    if (context_) {
      auto status = Detach();
      INTRINSIC_RT_LOG_IF(WARNING, !status.ok())
          << "Future was not detached: " << status.ToString();
    }
    // The destructor of `internal_context_` will wait for any attached promise
    // to be detached before destroying the context.
    internal_context_ = nullptr;
  }

  // Returns a future that is not attached to any context.
  // This can be used to create a temporary future in a real-time thread, i.e.
  // to cause the destruction of the assigned-to future.
  // Relies on RVO of the compiler to instantiate the future in-place.
  static RealtimeFuture<T> GetDetachedFuture() INTRINSIC_CHECK_REALTIME_SAFE {
    return RealtimeFuture<T>(/*context=*/nullptr, /*is_reusable=*/false);
  }

  // The future is not copyable or copy-assignable.
  RealtimeFuture(const RealtimeFuture&) = delete;
  RealtimeFuture& operator=(const RealtimeFuture&) = delete;

  // The future is movable and move-assignable.
  RealtimeFuture(RealtimeFuture&& other) noexcept
      INTRINSIC_CHECK_REALTIME_SAFE {
    // Call the move assignment operator for consistent logic.
    *this = std::move(other);
  }
  RealtimeFuture& operator=(RealtimeFuture&& other)
      INTRINSIC_CHECK_REALTIME_SAFE {
    if (this == &other) {
      return *this;
    }
    // Detach from the current context.
    if (context_) {
      auto status = Detach();
      INTRINSIC_RT_LOG_IF(WARNING, !status.ok())
          << "Future was not detached: " << status.ToString();
    }
    std::swap(context_, other.context_);
    std::swap(internal_context_, other.internal_context_);
    std::swap(is_reusable_, other.is_reusable_);
    return *this;
  }

  // Gets the value from the future [non-blocking].
  // Returns an UnavailableError if the value is not available yet.
  // Returns a CancelledError if the future was cancelled.
  // Returns a FailedPreconditionError if the future is not attached.
  icon::RealtimeStatusOr<T> Get() INTRINSIC_CHECK_REALTIME_SAFE {
    if (!context_) {
      return icon::FailedPreconditionError("Future is not attached.");
    }

    if (!context_->ReadTryLock()) {
      // If the try lock fails, the context is being cancelled.
      return icon::CancelledError("Future is being cancelled.");
    }

    if (context_->IsCancelled()) {
      INTRINSIC_RT_RETURN_IF_ERROR(context_->ReadUnlock());
      return icon::CancelledError("Future was cancelled.");
    }

    auto reader = context_->buffer_.reader();
    if (reader == nullptr) {
      INTRINSIC_RT_RETURN_IF_ERROR(context_->ReadUnlock());
      return icon::InternalError("Received nullptr reader.");
    }

    std::optional<T> element = reader->MoveFront();
    if (element.has_value()) {
      reader->DropFront();
      // Call TryWait() so that we can wait for the next value to be
      // available.
      context_->has_value_.TryWait();
      INTRINSIC_RT_RETURN_IF_ERROR(context_->ReadUnlock());
      if (!is_reusable_) {
        INTRINSIC_RT_RETURN_IF_ERROR(Detach());
      }
      return std::move(element.value());
    } else {
      // No need to call `DropFront()`, since the value retrieval failed.
      INTRINSIC_RT_RETURN_IF_ERROR(context_->ReadUnlock());
      return icon::UnavailableError(
          "Value is not available yet or has already been retrieved.");
    }
  }

  // Default copy function for Peek.
  // Peek allows providing a custom copy function to copy the value from the
  // future in cases where T is not copyable. For types that are copyable,
  // this function simply returns the value causing a normal copy.
  template <typename V>
  static V DefaultCopyFn(const V& t) {
    return t;
  }

  // Peeks the value from the future [non-blocking].
  // This will not remove the value from the context's buffer.
  // Use `copy_fn` to provide a custom copy function in case the value type T
  // is not (trivially) copyable. Returns an UnavailableError if the value is
  // not available yet. Returns a CancelledError if the future was cancelled.
  // Returns a FailedPreconditionError if the future is not attached.
  template <typename CopyFn = decltype(DefaultCopyFn<T>)>
  icon::RealtimeStatusOr<T> Peek(CopyFn copy_fn = DefaultCopyFn<T>)
      INTRINSIC_CHECK_REALTIME_SAFE {
    if (!context_) {
      return icon::FailedPreconditionError("Future is not attached.");
    }

    if (context_->IsCancelled()) {
      return icon::CancelledError("Future was cancelled.");
    }

    auto reader = context_->buffer_.reader();
    if (reader == nullptr) {
      return icon::InternalError("Received nullptr reader.");
    }

    auto front = reader->Front();
    if (front == nullptr) {
      // No need to call `KeepFront()`, since the value retrieval failed.
      return icon::UnavailableError("Value is not available yet.");
    }
    reader->KeepFront();
    return copy_fn(*front);
  }

  // Waits the specified duration for the value to become available or if the
  // future is cancelled.
  // Does not return the value. If successful, any call to `*Get*` or `*Peek*`
  // will immediately return the value.
  // Returns an error if the timeout occurs, if the future is cancelled or if
  // the future is not attached.
  absl::Status WaitFor(absl::Duration duration) INTRINSIC_NON_REALTIME_ONLY {
    return WaitUntil(absl::Now() + duration);
  }

  // Waits until the specified deadline for the value to become available.
  // Does not return the value. If successful, any call to `*Get*` or `*Peek*`
  // will immediately return the value.
  // Returns an error if the timeout occurs, if the future is cancelled or if
  // the future is not attached.
  absl::Status WaitUntil(absl::Time deadline) INTRINSIC_NON_REALTIME_ONLY {
    if (!context_) {
      return icon::FailedPreconditionError("Future is not attached.");
    }

    // `HasValue()` will check for cancellation and returns true if a value is
    // available in the buffer.
    // We're not using `has_value_.WaitUntil()` here, since has_value_ might
    // have been consumed already by a previous Wait...AndPeek call.
    INTR_ASSIGN_OR_RETURN(bool has_value, context_->HasValue());
    if (has_value) {
      return absl::OkStatus();
    }

    // It's safe to use `WaitUntil` here, since both `Cancel()` and `Set()`
    // will post to `has_value_`.
    INTR_RETURN_IF_ERROR(context_->has_value_.WaitUntil(deadline))
        << "New value is not available yet.";

    // If the buffer is empty after the wait, it must have been a cancellation.
    if (context_->IsCancelled()) {
      return icon::CancelledError("Future was cancelled.");
    }

    return absl::OkStatus();
  }

  // Waits until the specified deadline for the value to become available.
  // Returns the value if successful, or an error if the deadline is exceeded,
  // the future is cancelled or the future is not attached.
  absl::StatusOr<T> WaitUntilAndGet(absl::Time deadline)
      INTRINSIC_NON_REALTIME_ONLY {
    INTR_RETURN_IF_ERROR(WaitUntil(deadline));
    INTRINSIC_RT_ASSIGN_OR_RETURN(auto value, Get());
    return std::move(value);
  }

  // Waits for the specified duration for the value to become available.
  // Returns the value if successful, or an error if the timeout occurs,
  // the future is cancelled or the future is not attached.
  absl::StatusOr<T> WaitForAndGet(absl::Duration duration)
      INTRINSIC_NON_REALTIME_ONLY {
    INTR_RETURN_IF_ERROR(WaitFor(duration));
    INTRINSIC_RT_ASSIGN_OR_RETURN(auto value, Get());
    return std::move(value);
  }

  // Waits until the specified deadline for the value to become available and
  // peeks it.
  // Returns a copy of the value if successful, or an error if the deadline is
  // exceeded, the future is cancelled or the future is not attached.
  // Use `copy_fn` to provide a custom copy function in case the value type T
  // is not (trivially) copyable.
  template <typename CopyFn = decltype(DefaultCopyFn<T>)>
  absl::StatusOr<T> WaitUntilAndPeek(absl::Time deadline,
                                     CopyFn copy_fn = DefaultCopyFn<T>)
      INTRINSIC_NON_REALTIME_ONLY {
    INTR_RETURN_IF_ERROR(WaitUntil(deadline));
    INTRINSIC_RT_ASSIGN_OR_RETURN(auto value, Peek(copy_fn));
    return std::move(value);
  }

  // Waits for the specified duration for the value to become available and
  // peeks it.
  // Returns the value if successful, or an error if the timeout occurs,
  // the future is cancelled or the future is not attached.
  // Use `copy_fn` to provide a custom copy function in case the value type T
  // is not (trivially) copyable.
  template <typename CopyFn = decltype(DefaultCopyFn<T>)>
  absl::StatusOr<T> WaitForAndPeek(absl::Duration duration,
                                   CopyFn copy_fn = DefaultCopyFn<T>)
      INTRINSIC_NON_REALTIME_ONLY {
    INTR_RETURN_IF_ERROR(WaitFor(duration));
    INTRINSIC_RT_ASSIGN_OR_RETURN(auto value, Peek(copy_fn));
    return std::move(value);
  }

  // Cancels the future.
  // This informs the associated promise that no value will be consumed.
  // Returns a ResourceExhaustedError if a reset is in progress.
  // Returns a CancelledError if the context was already cancelled.
  // Returns a FailedPreconditionError if the future is not attached.
  icon::RealtimeStatus Cancel() INTRINSIC_CHECK_REALTIME_SAFE {
    if (!context_) {
      return icon::FailedPreconditionError("Future is not attached.");
    }
    // Call the underlying method to check for cancellation, since this has
    // all the logic to handle a concurrent `Reset` call.
    return context_->Cancel();
  }

  // Checks if the future has been cancelled.
  // Returns a FailedPreconditionError if the future is not attached.
  icon::RealtimeStatusOr<bool> IsCancelled() const
      INTRINSIC_CHECK_REALTIME_SAFE {
    if (!context_) {
      return icon::FailedPreconditionError("Future is not attached.");
    }
    // Call the underlying method to check for cancellation, since this has
    // all the logic to handle a concurrent `Reset` call.
    return context_->IsCancelled();
  }

  // Checks if a value is available in the future. Returns an error if the
  // future is not attached or has been cancelled.
  icon::RealtimeStatusOr<bool> HasValue() const INTRINSIC_CHECK_REALTIME_SAFE {
    if (!context_) {
      return icon::FailedPreconditionError("Future is not attached.");
    }
    return context_->HasValue();
  }

  // Gets the promise associated with this future's context.
  // Returns an error if the future is not attached or if a promise has
  // already been returned.
  icon::RealtimeStatusOr<RealtimePromise<T>> GetPromise(
      bool is_reusable = false) INTRINSIC_CHECK_REALTIME_SAFE {
    if (!context_) {
      return icon::FailedPreconditionError("Future is not attached.");
    }
    return context_->GetPromise(is_reusable);
  }

  // Returns true, if the future can be destroyed without waiting for a
  // promise to be detached.
  // This is always true, if the future doesn't own the context and false only
  // if it owns the context and a promise is attached. Note, that calling this
  // method right before destroying the future doesn't guarantee a wait-free
  // destruction, since a new promise might have been created in the meantime.
  bool IsWaitFreeDestructible() const INTRINSIC_CHECK_REALTIME_SAFE {
    // If the future is not attached, it is wait-free.
    if (!context_) {
      // This should never happen, but we check anyway.
      return true;
    }
    // If the future doesn't own the context, it is wait-free destructible.
    if (internal_context_.get() == nullptr) {
      return true;
    }
    // If a promise is attached, the future is not wait-free destructible.
    return context_->promise_attached_.load(std::memory_order_acquire) == false;
  }

 protected:
  // Friend declaration so that the FuturePromiseContext can access the
  // private constructor of RealtimeFuture.
  friend class FuturePromiseContext<T>;

  // Protected constructor used by FuturePromiseContext to create a future.
  explicit RealtimeFuture(Context* context, bool is_reusable)
      : context_(context),
        internal_context_(nullptr),
        is_reusable_(is_reusable) {}

  icon::RealtimeStatus Detach() INTRINSIC_CHECK_REALTIME_SAFE {
    if (!context_) {
      return icon::FailedPreconditionError("Future is not attached.");
    }
    auto status = context_->DetachFuture();
    context_ = nullptr;
    return status;
  }

  // Pointer to the shared context.
  Context* context_{nullptr};
  // The optional internal context. Exists only if default constructed, in
  // which case the context_ pointer points to the internal context. Will be
  // nullptr if created from a FuturePromiseContext or RealtimePromise.
  std::unique_ptr<Context> internal_context_{nullptr};

  // Indicates if the future is reusable.
  // A reusable future won't detach from the context after its `Get()` method is
  // called and thus can be reused again to get a new value.
  bool is_reusable_{false};
};

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_COMMON_BUFFERS_RT_PROMISE_H_
