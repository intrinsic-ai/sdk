// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_THREAD_CONCURRENT_QUEUE_H_
#define INTRINSIC_UTIL_THREAD_CONCURRENT_QUEUE_H_

#include <deque>

#include "absl/base/thread_annotations.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/synchronization/mutex.h"
#include "absl/time/time.h"

namespace intrinsic {

template <typename T>
class ConcurrentQueue {
 public:
  // Constructs a ConcurrentQueue with the given maximum size.
  // If the max_size is zero, the queue will be unbounded.
  explicit ConcurrentQueue(int max_size) : max_size_(max_size) {}

  absl::Status Enqueue(const T& item,
                       absl::Duration timeout = absl::InfiniteDuration())
      ABSL_LOCKS_EXCLUDED(mutex_) {
    return Enqueue(T(item), timeout);
  }

  // Enqueues the given item into the queue. If the queue is full, waits for
  // space to become available until the given timeout. If the timeout is
  // exceeded, returns a ResourceExhaustedError. If the user specifies a timeout
  // of zero (or negative), this function will return immediately.
  // The function returns an UnavailableError if the queue is closed.
  absl::Status Enqueue(T&& item,
                       absl::Duration timeout = absl::InfiniteDuration())
      ABSL_LOCKS_EXCLUDED(mutex_) {
    auto queue_not_at_max_size_or_closed =
        [this]() ABSL_EXCLUSIVE_LOCKS_REQUIRED(mutex_) {
          return max_size_ == 0 || queue_.size() < max_size_ || closed_;
        };
    absl::MutexLock lock(mutex_);
    if (mutex_.AwaitWithTimeout(
            absl::Condition(&queue_not_at_max_size_or_closed), timeout)) {
      if (closed_) {
        // The await condition is true because the queue is closed.
        return absl::UnavailableError("Queue is closed.");
      }
      queue_.push_back(std::move(item));
      return absl::OkStatus();
    }
    return absl::ResourceExhaustedError("Queue is full.");
  }

  // Dequeues the first item in the queue. If the queue is empty, waits for an
  // item to become available until the given timeout. If the timeout is
  // exceeded, returns a DeadlineExceededError. If the user specifies a timeout
  // of zero (or negative), this function will return immediately.
  // The function returns an UnavailableError if the queue is closed and all
  // values have been dequeued.
  absl::StatusOr<T> Dequeue(absl::Duration timeout = absl::InfiniteDuration())
      ABSL_LOCKS_EXCLUDED(mutex_) {
    auto queue_not_empty_or_closed = [this]()
                                         ABSL_EXCLUSIVE_LOCKS_REQUIRED(mutex_) {
                                           return !queue_.empty() || closed_;
                                         };
    absl::MutexLock lock(mutex_);
    if (mutex_.AwaitWithTimeout(absl::Condition(&queue_not_empty_or_closed),
                                timeout)) {
      if (closed_ && queue_.empty()) {
        // The await condition is true because the queue is closed and all
        // values have been already dequeued.
        return absl::UnavailableError("Queue is closed.");
      }
      T element = std::move(queue_.front());
      queue_.pop_front();
      return element;
    }
    return absl::DeadlineExceededError("Queue is empty.");
  }

  bool IsEmpty() const ABSL_LOCKS_EXCLUDED(mutex_) {
    absl::MutexLock lock(mutex_);
    return queue_.empty();
  }

  // Returns the current size of the queue.
  int GetSize() const ABSL_LOCKS_EXCLUDED(mutex_) {
    absl::MutexLock lock(mutex_);
    return queue_.size();
  }

  // Closes the queue. No more items can be enqueued or dequeued after this
  // call and all subsequent calls to Enqueue and Dequeue will return an
  // UnavailableError.
  void Close() ABSL_LOCKS_EXCLUDED(mutex_) {
    absl::MutexLock lock(mutex_);
    closed_ = true;
  }

  void Clear() ABSL_LOCKS_EXCLUDED(mutex_) {
    absl::MutexLock lock(mutex_);
    queue_.clear();
  }

 private:
  mutable absl::Mutex mutex_;
  int max_size_;
  bool closed_ ABSL_GUARDED_BY(mutex_) = false;
  std::deque<T> queue_ ABSL_GUARDED_BY(mutex_);
};

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_THREAD_CONCURRENT_QUEUE_H_
