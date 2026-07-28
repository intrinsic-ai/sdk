// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/thread/thread_pool.h"

#include <utility>

#include "absl/functional/any_invocable.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/time/time.h"
#include "intrinsic/util/thread/stop_token.h"

namespace intrinsic {

ThreadPool::ThreadPool(int num_threads, int max_queue_size)
    : queue_(max_queue_size) {
  for (int i = 0; i < num_threads; ++i) {
    workers_.emplace_back([this](StopToken stop_token) {
      while (!stop_token.stop_requested()) {
        auto f = queue_.Dequeue();
        if (f.status().code() == absl::StatusCode::kUnavailable) {
          // The queue is closed, so we should stop.
          return;
        } else if (!f.ok()) {
          LOG_EVERY_N_SEC(ERROR, 1)
              << "Failed to dequeue from thread pool queue: " << f.status();
          continue;
        }
        // `std::move` is required because the `operator()` of `f` is
        // rvalue-reference qualified. This ensures that the function is
        // executed exactly once.
        (*std::move(f))(stop_token);
      }
    });
  }
}

ThreadPool::~ThreadPool() {
  // Closing the queue will ensure that ConcurrentQueue::Dequeue() immediately
  // returns UnavailableError for all subsequent calls which forces immediate
  // termination of all worker threads.
  queue_.Close();

  // The worker thread's dtor calls request_stop() and then join().
  // For this dtor to not block forever, it is essential, the scheduled
  // invocables are not blocking and instead eventually terminate. This is
  // guaranteed by the StopToken.
}

absl::Status ThreadPool::Schedule(absl::AnyInvocable<void(StopToken) &&> f,
                                  absl::Duration timeout) {
  return queue_.Enqueue(std::move(f), timeout);
}

absl::Status ThreadPool::Schedule(absl::AnyInvocable<void() &&> f,
                                  absl::Duration timeout) {
  return queue_.Enqueue(
      [f = std::move(f)](StopToken) mutable {
        // `std::move` is required because the `operator()` of `f` is
        // rvalue-reference qualified. This ensures that the function is
        // executed exactly once.
        std::move(f)();
      },
      timeout);
}

int ThreadPool::GetSize() const { return queue_.GetSize(); }

}  // namespace intrinsic
