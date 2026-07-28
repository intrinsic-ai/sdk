// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_THREAD_THREAD_POOL_H_
#define INTRINSIC_UTIL_THREAD_THREAD_POOL_H_

#include <vector>

#include "absl/functional/any_invocable.h"
#include "absl/status/status.h"
#include "absl/time/time.h"
#include "intrinsic/util/thread/concurrent_queue.h"
#include "intrinsic/util/thread/stop_token.h"
#include "intrinsic/util/thread/thread.h"

namespace intrinsic {

// A thread pool that allows scheduling functions to be executed on one of the
// worker threads.
// The worker threads are started when the ThreadPool is constructed and
// terminated when the ThreadPool is destructed. The ThreadPool's dtor
// blocks until all worker threads have terminated. The passed function must
// return eventually (i.e. not block forever) otherwise the ThreadPool's
// dtor will block forever.
// When many hetereogeneously sized tasks are scheduled, later tasks maybe
// starved by earlier tasks.
class ThreadPool {
 public:
  // Create a thread pool with the given number of threads and a queue with the
  // given maximum size. If the max_queue_size is zero, the queue will be
  // unbounded.
  explicit ThreadPool(int num_threads, int max_queue_size = 0);

  // Closes the queue and waits for all worker threads to terminate.
  ~ThreadPool();

  // Schedules the given function to be executed on one of the worker threads.
  // If the queue is full, returns a ResourceExhaustedError if the timeout is
  // exceeded. If the timeout is zero (or negative), this function will return
  // immediately with a ResourceExhaustedError, an UnavailableError if the
  // queue is closed, or OkStatus if enqueuing was successful.
  // The passed function must return eventually (i.e. not block forever)
  // otherwise the ThreadPool's dtor will block forever.
  absl::Status Schedule(absl::AnyInvocable<void(StopToken) &&> f,
                        absl::Duration timeout = absl::ZeroDuration());

  // Schedules the given function to be executed on one of the worker threads.
  // If the queue is full, returns a ResourceExhaustedError if the timeout is
  // exceeded. If the timeout is zero (or negative), this function will return
  // immediately with a ResourceExhaustedError, an UnavailableError if the
  // queue is closed, or OkStatus if enqueuing was successful.
  // The passed function must return eventually (i.e. not block forever)
  // otherwise the ThreadPool's dtor will block forever.
  absl::Status Schedule(absl::AnyInvocable<void() &&> f,
                        absl::Duration timeout = absl::ZeroDuration());

  // Get current size of the queue.
  int GetSize() const;

 private:
  ConcurrentQueue<absl::AnyInvocable<void(StopToken) &&>> queue_;
  std::vector<Thread> workers_;
};

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_THREAD_THREAD_POOL_H_
