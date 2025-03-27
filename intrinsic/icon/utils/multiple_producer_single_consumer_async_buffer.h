// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_UTILS_MULTIPLE_PRODUCER_SINGLE_CONSUMER_ASYNC_BUFFER_H_
#define INTRINSIC_ICON_UTILS_MULTIPLE_PRODUCER_SINGLE_CONSUMER_ASYNC_BUFFER_H_

#include <functional>

#include "absl/base/thread_annotations.h"
#include "absl/status/status.h"
#include "absl/synchronization/mutex.h"
#include "intrinsic/icon/testing/realtime_annotations.h"
#include "intrinsic/icon/utils/async_buffer.h"

namespace intrinsic {

// MultipleProducerSingleConsumerAsyncBuffer is a buffer that supports non
// blocking reads for the single consumer, and blocking thread-safe writes for
// multiple producers.
// It is based on `AsyncBuffer`, but adds a mutex to allow multiple producers.
//
// Example:
//
//   MultipleProducerSingleConsumerAsyncBuffer<int> buffer;
//
//   // Thread 1:
//   buffer.UpdateBuffer(42);
//
//   // Thread 2:
//   buffer.UpdateBuffer([](int& buffer) {
//     buffer = 1337;
//   });
//
//   // Thread 3:
//   int* value;
//   buffer.GetActiveBuffer(&value);
//   EXPECT_EQ(*value, 42);
template <typename T>
class MultipleProducerSingleConsumerAsyncBuffer {
 public:
  using AsyncBuffer = AsyncBuffer<T>;

  // Updates and commits the free buffer in an atomic operation.
  //
  // This function is not realtime safe.
  //
  // Returns an error if the commit fails. This should never happen since we
  // just got the free buffer and have a mutex on async_buffer_.
  absl::Status UpdateBuffer(const T& value) INTRINSIC_NON_REALTIME_ONLY {
    absl::MutexLock lock(&mutex_);
    *async_buffer_.GetFreeBuffer() = value;
    // Commit should always return true since we just got the free buffer and
    // have a mutex on async_buffer_.
    if (!async_buffer_.CommitFreeBuffer()) {
      return absl::InternalError(
          "Failed to commit free buffer in "
          "MultipleProducerSingleConsumerAsyncBuffer. This is a bug.");
    }
    return absl::OkStatus();
  }

  // Updates and commits the free buffer in an atomic operation.
  //
  // This function is not realtime safe.
  //
  // Returns an error if the commit fails. This should never happen since we
  // just got the free buffer and have a mutex on async_buffer_.
  absl::Status UpdateBuffer(T&& value) INTRINSIC_NON_REALTIME_ONLY {
    absl::MutexLock lock(&mutex_);
    *async_buffer_.GetFreeBuffer() = std::move(value);
    // Commit should always return true since we just got the free buffer and
    // have a mutex on async_buffer_.
    if (!async_buffer_.CommitFreeBuffer()) {
      return absl::InternalError(
          "Failed to commit free buffer in "
          "MultipleProducerSingleConsumerAsyncBuffer. This is a bug.");
    }
    return absl::OkStatus();
  }

  // Updates and commits the free buffer in an atomic operation.
  // `callback` should update the given value to perform the update and gives
  // more flexibility how the value is updated. This overload is useful if the
  // value is not copyable nor movable.
  //
  // This function is not realtime safe.
  //
  // Returns an error if the commit fails. This should never happen since we
  // just got the free buffer and have a mutex on async_buffer_.
  //
  // Example:
  //
  //   MultipleProducerSingleConsumerAsyncBuffer<int> buffer;
  //   buffer.UpdateBuffer([](int& buffer) {
  //     buffer = 42;
  //   });
  absl::Status UpdateBuffer(std::function<void(T&)> callback)
      INTRINSIC_NON_REALTIME_ONLY {
    absl::MutexLock lock(&mutex_);
    callback(*async_buffer_.GetFreeBuffer());
    // Commit should always return true since we just got the free buffer and
    // have a mutex on async_buffer_.
    if (!async_buffer_.CommitFreeBuffer()) {
      return absl::InternalError(
          "Failed to commit free buffer in "
          "MultipleProducerSingleConsumerAsyncBuffer. This is a bug.");
    }
    return absl::OkStatus();
  }

  // Returns the active buffer to the consumer. This buffer is guaranteed not
  // to be modified until the next call to GetActiveBuffer() by the consumer.
  // This may be the same buffer which was returned to last call to
  // GetActiveBuffer().
  //
  // This function is realtime safe.
  //
  // Returns true if the mailbox buffer was full and the returned pointer
  // points to it. Returns false if the mailbox buffer was empty and the
  // returned pointer points to the buffer that was already active at the time
  // of this call.
  bool GetActiveBuffer(T** buffer) INTRINSIC_CHECK_REALTIME_SAFE {
    return ABSL_TS_UNCHECKED_READ(async_buffer_).GetActiveBuffer(buffer);
  }

 private:
  absl::Mutex mutex_;
  AsyncBuffer async_buffer_;
};
}  // namespace intrinsic

#endif  // INTRINSIC_ICON_UTILS_MULTIPLE_PRODUCER_SINGLE_CONSUMER_ASYNC_BUFFER_H_
