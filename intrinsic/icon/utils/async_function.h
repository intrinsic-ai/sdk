// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_UTILS_ASYNC_FUNCTION_H_
#define INTRINSIC_ICON_UTILS_ASYNC_FUNCTION_H_

#include <memory>
#include <tuple>
#include <type_traits>
#include <utility>

#include "absl/functional/any_invocable.h"
#include "absl/status/status.h"
#include "absl/time/time.h"
#include "intrinsic/icon/testing/realtime_annotations.h"
#include "intrinsic/icon/utils/log.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_macro.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/platform/common/buffers/rt_promise.h"
#include "intrinsic/util/thread/stop_token.h"
#include "intrinsic/util/thread/thread.h"

namespace intrinsic::icon {
// `AsyncFunction` allows to execute a function asynchronously.
// It supports two modes of execution:
//   * Realtime: The function is executed on the calling thread. This is
//     suitable for short, real-time safe functions.
//   * Non-Realtime: The function is executed on a separate thread. This is
//     suitable for long-running or blocking operations that should not be on a
//     real-time thread.
//
// Both modes of execution return a `RealtimeFuture` to provide a consistent
// interface. The AsyncFunction can be called multiple times (not in parallel,
// though).
//
// Usage example:
//
// // Define a function that you want to execute asynchronously.
// int my_function(int a, float b) {
//    // ... some work ...
//    return a * static_cast<int>(b);
// }
//
// // Create an AsyncFunction for the function.
// auto async_function_rt =
//   AsyncFunction<int(int, float)>::CreateFromRealtimeFunction(my_function);
//
// // Call the function with arguments. This will execute on the calling thread,
// // since it was defined as a real-time safe function.
// INTR_ASSIGN_OR_RETURN(auto future, (*async_function_rt)(5, 2.3f));
//
// // Get the result when it is available. Since we know that the function is
// // executed on the calling thread, we could have used `Get` instead to get
// // the result without blocking.
// INTR_ASSIGN_OR_RETURN(int result, future.WaitForAndGet());
//
// // Alternatively, create an AsyncFunction for the same function, but
// // executing on a separate thread.
// auto async_function_non_rt =
//   AsyncFunction<int(int, float)>::CreateFromNonRealtimeFunction(my_function);
//
// // Call the function with arguments.
// INTR_ASSIGN_OR_RETURN(auto future_non_rt, (*async_function_non_rt)(5, 2.3f));
//
// // Get the result when it is available.
// INTR_ASSIGN_OR_RETURN(int result_non_rt, future_non_rt.WaitForAndGet());
//
// // You can call the same AsyncFunction again.
// INTR_ASSIGN_OR_RETURN(future_non_rt, (*async_function_non_rt)(5, 2.3f));
//
// // Get the result when it is available.
// INTR_ASSIGN_OR_RETURN(result_non_rt, future_non_rt.WaitForAndGet());

// Forward declaration for AsyncRealtimeFunction.
template <typename Signature>
class AsyncRealtimeFunction;

// Forward declaration for AsyncNonRealtimeFunction.
template <typename Signature>
class AsyncNonRealtimeFunction;

namespace internal {
// A helper struct to extract properties from a function signature.
// The primary template is left undefined.
template <typename T>
struct FunctionTraits;

// Partial specialization of FunctionTraits for function types.
// This allows extracting the return type and argument types from a function
// signature like `R(Args...)`.
template <typename R, typename... Args>
struct FunctionTraits<R(Args...)> {
  // The return type of the function.
  using ResultType = R;
  // A std::tuple containing all argument types.
  using ArgsTuple = std::tuple<Args...>;
};
}  // namespace internal

// An abstract base class for functions that can be executed asynchronously.
// This class provides a common interface for different asynchronous execution
// strategies, such as running on a separate thread or running on the calling
// thread but returning a future. Since `RealtimeFuture<void>` is not
// supported,the return type of the function must not be void.
//
// The class is templated on a function signature `Signature`, e.g., `int(int,
// float)`.
//
// Template parameters:
//   Signature: The function signature of the function to be executed.
template <typename Signature>
class AsyncFunction {
  using Traits = internal::FunctionTraits<Signature>;

 public:
  using ResultType = typename Traits::ResultType;
  using ArgsTuple = typename Traits::ArgsTuple;
  using FutureType = intrinsic::RealtimeFuture<ResultType>;
  static_assert(!std::is_void_v<ResultType>, "ResultType must not be void.");

  // Creates an `AsyncFunction` that executes a real-time safe function on the
  // calling thread.
  // It's the caller's responsibility to ensure that the function is
  // real-time safe.
  // @param func The real-time safe function to wrap.
  // @return A unique_ptr to an `AsyncFunction` instance.
  static std::unique_ptr<AsyncFunction<Signature>> CreateFromRealtimeFunction(
      absl::AnyInvocable<Signature> func) INTRINSIC_NON_REALTIME_ONLY {
    return std::make_unique<AsyncRealtimeFunction<Signature>>(std::move(func));
  }

  // Creates an `AsyncFunction` that executes a non-real-time function on a
  // background thread.
  // @param func The non-real-time function to wrap.
  // @return A unique_ptr to an `AsyncFunction` instance.
  static std::unique_ptr<AsyncFunction<Signature>>
  CreateFromNonRealtimeFunction(absl::AnyInvocable<Signature> func)
      INTRINSIC_NON_REALTIME_ONLY {
    return std::make_unique<AsyncNonRealtimeFunction<Signature>>(
        std::move(func));
  }

  AsyncFunction() = default;
  virtual ~AsyncFunction() = default;

  AsyncFunction(AsyncFunction&& other) = default;
  AsyncFunction& operator=(AsyncFunction&& other) = default;

  // Invokes the asynchronous function with the given arguments.
  // This is a convenience operator that forwards the call to the virtual
  // `CallImpl` method.
  //
  // Template parameters:
  //   CallArgs: The types of the arguments to pass to the function.
  // @param args The arguments to pass to the asynchronous function.
  // @return A `RealtimeStatusOr` containing a `RealtimeFuture` for the result,
  // or an error status if the invocation fails. The specific error conditions
  // depend on the derived class implementation.
  template <typename... CallArgs>
  auto operator()(CallArgs&&... args)
    requires(std::is_invocable_v<absl::AnyInvocable<Signature>&, CallArgs...>)
  INTRINSIC_CHECK_REALTIME_SAFE {
    // Pack arguments into a tuple and call the virtual function.
    return this->CallImpl(std::make_tuple(std::forward<CallArgs>(args)...));
  }

 protected:
  // The internal implementation of the function call.
  // Derived classes must implement this method to define the execution
  // strategy.
  // @param args A tuple containing the arguments for the function call.
  // @return A `RealtimeStatusOr` containing a `RealtimeFuture` for the result.
  // Returns an error status if the invocation fails. For example, it might
  // return a `ResourceExhaustedError` if the function is already running, or an
  // `InternalError` if the object has been moved from.
  virtual intrinsic::icon::RealtimeStatusOr<FutureType> CallImpl(
      std::add_rvalue_reference_t<ArgsTuple> args) = 0;
};

// An `AsyncFunction` implementation that executes a function on the calling
// thread. The execution is synchronous, but it returns a `RealtimeFuture` to
// provide a consistent asynchronous-style interface. This is suitable for
// short, real-time safe functions. It's the caller's responsibility to ensure
// that the function is real-time safe.
template <typename Signature>
class AsyncRealtimeFunction : public AsyncFunction<Signature> {
  using Traits = internal::FunctionTraits<Signature>;

 public:
  using ResultType = typename Traits::ResultType;
  using ArgsTuple = typename Traits::ArgsTuple;
  static_assert(!std::is_void_v<ResultType>, "ResultType must not be void.");

 private:
  // The implementation details of AsyncRealtimeFunction, using the PIMPL
  // idiom to allow for cheap, non-allocating and non-blocking move operations.
  struct Impl {
    explicit Impl(absl::AnyInvocable<Signature> f)
        : return_data(), func(std::move(f)) {}

    // Delete all 5 constructors and assignment operators to prevent the
    // accidental move or copy of the underlying data.
    Impl() = delete;
    Impl(const Impl&) = delete;
    Impl(Impl&& other) = delete;
    Impl& operator=(const Impl&) = delete;
    Impl& operator=(Impl&& other) = delete;

    ~Impl() = default;

    // Context for managing the promise and future for the function's return
    // value.
    intrinsic::FuturePromiseContext<ResultType> return_data;
    // The function to be executed.
    absl::AnyInvocable<Signature> func;
  };

 public:
  // Constructs an `AsyncRealtimeFunction`.
  // @param func The function to be executed. If `nullptr`, calls to this
  // object will fail.
  explicit AsyncRealtimeFunction(absl::AnyInvocable<Signature> func = nullptr)
      INTRINSIC_NON_REALTIME_ONLY
      : AsyncFunction<Signature>(),
        pimpl_(std::make_unique<Impl>(std::move(func))) {}

 protected:
  // Executes the function on the calling thread.
  // @param args A tuple of arguments to pass to the function.
  // @return A `RealtimeStatusOr` containing a `RealtimeFuture` with the result
  // of the function. Returns a `FailedPreconditionError` if the wrapped
  // function is null. Returns an `InternalError` if the `AsyncRealtimeFunction`
  // object has been moved from. The status may also contain errors from
  // promise/future operations.
  intrinsic::icon::RealtimeStatusOr<intrinsic::RealtimeFuture<ResultType>>
  CallImpl(std::add_rvalue_reference_t<ArgsTuple> args)
      INTRINSIC_CHECK_REALTIME_SAFE override {
    if (pimpl_ == nullptr) {
      return intrinsic::icon::InternalError(
          "AsyncRealtimeFunction is not initialized. Was it invalidated by a "
          "move operation?");
    }
    if (!pimpl_->func) {
      return intrinsic::icon::FailedPreconditionError(
          "AsyncRealtimeFunction holds a null function.");
    }

    // Check if no previous result future is still attached.
    if (pimpl_->return_data.IsFutureAttached()) {
      return intrinsic::icon::ResourceExhaustedError(
          "A previous result future is still attached.");
    }
    // Consume any previous result.
    auto has_value = pimpl_->return_data.HasValue();
    if (!has_value.ok()) {
      return intrinsic::icon::ResourceExhaustedError(
          "Failed to query if a previous result is still available.");
    }
    if (has_value.value()) {
      // Get the future and consume the result.
      auto future = pimpl_->return_data.GetFuture();
      if (!future.ok()) {
        return intrinsic::icon::ResourceExhaustedError(
            "Failed to get a future to consume a previous result.");
      }
      auto result = future.value().Get();
      if (!result.ok()) {
        return intrinsic::icon::ResourceExhaustedError(
            "Failed to consume a previous result.");
      }
    }

    // Execute the function in-place.
    INTRINSIC_RT_ASSIGN_OR_RETURN(auto future, pimpl_->return_data.GetFuture());
    INTRINSIC_RT_ASSIGN_OR_RETURN(auto promise,
                                  pimpl_->return_data.GetPromise());
    INTRINSIC_RT_RETURN_IF_ERROR(
        promise.Set(std::apply(pimpl_->func, std::move(args))));
    // The function was executed on the calling thread, so the result is
    // immediately available.
    return std::move(future);
  }

  std::unique_ptr<Impl> pimpl_{nullptr};
};

// An `AsyncFunction` implementation that runs a function in a background
// thread. This is suitable for long-running or blocking operations that should
// not be on a real-time thread.
//
// The `operator()` is real-time safe and sends the arguments to the background
// thread, which then executes the function. The result is returned via a
// `RealtimeFuture`.
//
// Template parameters:
//   Signature: The function signature, e.g., `int(int)`.
template <typename Signature>
class AsyncNonRealtimeFunction : public AsyncFunction<Signature> {
  using Traits = internal::FunctionTraits<Signature>;

 public:
  using ResultType = typename Traits::ResultType;
  using ArgsTuple = typename Traits::ArgsTuple;
  static_assert(!std::is_void_v<ResultType>, "ResultType must not be void.");

 private:
  // The implementation details of AsyncNonRealtimeFunction, using the PIMPL
  // idiom.
  // Allows us to cheaply move the AsyncNonRealtimeFunction without
  // destructing the underlying data, which would require thread-recreation due
  // to dangling references.
  struct Impl {
    explicit Impl(absl::AnyInvocable<Signature> f)
        : return_data(),
          args_data(),
          func(std::move(f)),
          async_thread(intrinsic::Thread(
              [this](intrinsic::StopToken st) { ThreadFunc(st); })) {
      if (!func) {
        // If the function is null, cancel the data such that any returned
        // future will be in a cancelled state.
        if (!return_data.Cancel().ok()) {
          INTRINSIC_RT_LOG(ERROR)
              << "Failed to cancel return_data when created "
                 "with a null function.";
        }
        if (!args_data.Cancel().ok()) {
          INTRINSIC_RT_LOG(ERROR) << "Failed to cancel args_data when created "
                                     "with a null function.";
        }
      }
    }

    // Delete all 5 constructors and assignment operators to prevent the
    // accidental move or copy of the underlying data.
    Impl() = delete;
    Impl(const Impl&) = delete;
    Impl(Impl&& other) = delete;
    Impl& operator=(const Impl&) = delete;
    Impl& operator=(Impl&& other) = delete;

    ~Impl() {
      if (!return_data.Cancel().ok()) {
        INTRINSIC_RT_LOG(ERROR) << "Failed to cancel return_data.";
      };
      if (!args_data.Cancel().ok()) {
        INTRINSIC_RT_LOG(ERROR) << "Failed to cancel args_data.";
      }
    }

    // The function is executed by the background thread.
    // It waits for arguments, executes the stored function, and sets the
    // result promise. This function runs in a loop until the `StopToken` is
    // requested.
    // @param st A `StopToken` used to signal the thread to terminate.
    void ThreadFunc(intrinsic::StopToken st) {
      auto args_future_or = args_data.GetFuture(/*is_reusable=*/true);
      auto return_promise_or = return_data.GetPromise(/*is_reusable=*/true);
      absl::Status status = args_future_or.status();
      status = intrinsic::icon::OverwriteIfNotInError(
          status, return_promise_or.status());
      if (!status.ok()) {
        INTRINSIC_RT_LOG(ERROR)
            << "AsyncFunction failed to start due to: " << status.message();
        return;
      }
      auto return_promise = std::move(return_promise_or).value();
      auto args_future = std::move(args_future_or).value();

      while (!st.stop_requested()) {
        // Just peek the value to prevent another set of args from being set,
        // while we're still processing the current set.
        auto args = args_future.WaitUntilAndGet(absl::InfiniteFuture());
        if (!args.ok()) {
          // If this happens, the args_data was most likely cancelled due to the
          // destruction of this object. We just need to continue the loop. If
          // the stop token is set, the thread will become joinable.
          continue;
        }

        // Call the function and set the return promise.
        auto status =
            return_promise.Set(std::apply(func, std::move(args.value())));
        if (!status.ok()) {
          // It's unlikely that this happens, so we log an error and exit the
          // loop.
          INTRINSIC_RT_LOG(ERROR)
              << "return_promise.Set failed due to: " << status.message();
          break;
        }
      }
    }

    // Context for managing the promise and future for the function's return
    // value.
    intrinsic::FuturePromiseContext<ResultType> return_data;
    // Context for passing arguments from the calling thread to the background
    // thread.
    intrinsic::FuturePromiseContext<ArgsTuple> args_data;
    // The function to be executed on the background thread.
    absl::AnyInvocable<Signature> func;
    // The background thread that executes `func`.
    intrinsic::Thread async_thread;
  };

 public:
  // Constructs an `AsyncNonRealtimeFunction` and starts its background thread.
  // @param func The function to be executed asynchronously. If `nullptr`, calls
  // to this object will fail.
  explicit AsyncNonRealtimeFunction(
      absl::AnyInvocable<Signature> func = nullptr) INTRINSIC_NON_REALTIME_ONLY
      : AsyncFunction<Signature>(),
        pimpl_(std::make_unique<Impl>(std::move(func))) {}

  AsyncNonRealtimeFunction(AsyncNonRealtimeFunction&& other) = default;
  AsyncNonRealtimeFunction& operator=(AsyncNonRealtimeFunction&& other) =
      default;

  // Deleted copy constructor. AsyncFunction is not copyable.
  AsyncNonRealtimeFunction(const AsyncNonRealtimeFunction&) = delete;
  // Deleted copy assignment operator. AsyncFunction is not copyable.
  AsyncNonRealtimeFunction& operator=(const AsyncNonRealtimeFunction&) = delete;

  // Destructor. Stops and joins the background thread if it is running.
  ~AsyncNonRealtimeFunction() override = default;

 protected:
  // Passes arguments to the background thread to invoke the function.
  // This method is real-time safe.
  // @param args A tuple of arguments to pass to the function.
  // @return A `RealtimeStatusOr` containing a `RealtimeFuture` for the result,
  // or an error status. Returns `ResourceExhaustedError` if a previous call is
  // still pending. Returns `FailedPreconditionError` if the wrapped function
  // is null. Returns an `InternalError` if the `AsyncNonRealtimeFunction`
  // object has been moved from. The status can also contain errors from
  // promise/future operations.
  intrinsic::icon::RealtimeStatusOr<intrinsic::RealtimeFuture<ResultType>>
  CallImpl(std::add_rvalue_reference_t<ArgsTuple> args) override {
    if (pimpl_ == nullptr) {
      return intrinsic::icon::InternalError(
          "AsyncNonRealtimeFunction is not initialized. Was it moved from?");
    }
    if (!pimpl_->func) {
      return intrinsic::icon::FailedPreconditionError(
          "AsyncNonRealtimeFunction holds a null function.");
    }
    // Don't call the function, if the previous args have not been processed.
    INTRINSIC_RT_ASSIGN_OR_RETURN(auto has_args, pimpl_->args_data.HasValue());
    if (has_args) {
      return intrinsic::icon::ResourceExhaustedError(
          "AsyncFunction is already running.");
    }
    // Check if no previous result future is still attached.
    if (pimpl_->return_data.IsFutureAttached()) {
      return intrinsic::icon::ResourceExhaustedError(
          "A previous result future is still attached.");
    }
    // Consume any previous result.
    auto has_value = pimpl_->return_data.HasValue();
    if (!has_value.ok()) {
      return intrinsic::icon::ResourceExhaustedError(
          "Failed to query if a previous result is still available.");
    }
    if (has_value.value()) {
      // Get the future and consume the result.
      auto future = pimpl_->return_data.GetFuture();
      if (!future.ok()) {
        return intrinsic::icon::ResourceExhaustedError(
            "Failed to get a future to consume a previous result.");
      }
      auto result = future.value().Get();
      if (!result.ok()) {
        return intrinsic::icon::ResourceExhaustedError(
            "Failed to consume a previous result.");
      }
    }

    INTRINSIC_RT_ASSIGN_OR_RETURN(auto promise, pimpl_->args_data.GetPromise());
    INTRINSIC_RT_RETURN_IF_ERROR(promise.Set(std::move(args)));
    return pimpl_->return_data.GetFuture();
  }

  std::unique_ptr<Impl> pimpl_{nullptr};
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_UTILS_ASYNC_FUNCTION_H_
